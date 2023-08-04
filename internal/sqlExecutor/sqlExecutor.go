// Mcserver
// Copyright (C) 2023  JUAN JOSÉ IGLESIAS BLANCH

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package sqlExecutor

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"

	"github.com/joanjib/fullness-server/internal/executor"
	"github.com/joanjib/fullness-server/internal/files"
	"github.com/joanjib/fullness-server/internal/sqlLoader"
)

var ExecutionQueue chan executor.ExecutorImpl
var StopSQLExecutor bool = false
var stMap map[string]*sql.Stmt
var stMapCache map[string]*sql.Stmt = make(map[string]*sql.Stmt) // Stmt cache for the current transaction

func GetStmt(st string, tx *sql.Tx) *sql.Stmt {
	if stMapCache[st] == nil {
		stMapCache[st] = tx.Stmt(stMap[st])
	}
	return stMapCache[st]
}

func Init(queueSize int) {
	ExecutionQueue = make(chan executor.ExecutorImpl, queueSize)
}
func initStmtMap(db *sql.DB, strStatements map[string]string) map[string]*sql.Stmt {

	psp := make(map[string]*sql.Stmt, len(strStatements))
	for k, s := range strStatements {
		stmt, err := db.Prepare(s)
		if err != nil {
			log.Panic(err)
		}
		psp[k] = stmt
	}
	return psp
}

func resetStCache() {
	for k := range stMapCache {
		stMapCache[k] = nil
	}
}

// main function of the goroutine.
func SQLExecutor(dbFile string, sqlStatementsDir string, shotTimeOut int, schemaFile string) {
	var tx *sql.Tx

	db, err := sql.Open("sqlite3", "file:"+dbFile+"?_journal_mode=WAL&_synchronous=FULL")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	//creating the database if not exists
	schema := string(files.LoadFile(schemaFile))
	_, err = db.Exec(schema)
	if err != nil {
		log.Panic(err)
	}

	strStatements := sqlLoader.LoadStrStatements(sqlStatementsDir)
	stMap = initStmtMap(db, strStatements)
	// closing the prep statements before exiting the go routine.
	for _, st := range stMap {
		defer st.Close()
	}
	// executed actions on each cicle.
	executed := make([]executor.ExecutorImpl, 0, 100)

	t := time.Duration(shotTimeOut) * time.Millisecond
	newTimer := time.NewTimer(t)
	accumulatedExecutionTime := time.Duration(0) * time.Millisecond
	i := 0
	sp := "" //savepoint sequence
outer:
	for {
		select {
		case toExec := <-ExecutionQueue:
			// here goes the execution:
			if len(executed) == 0 {
				tx, err = db.Begin()
				if err != nil {
					log.Panic(err)
				}
				i = 0
				sp = "A" + strconv.Itoa(i)
			} else {
				i++
				sp = "A" + strconv.Itoa(i)
			}
			now := time.Now()

			toExec.Execute(GetStmt(toExec.GetAction(), tx), tx, sp)
			accumulatedExecutionTime += time.Since(now)
			toExec.SetDataAvailable()
			executed = append(executed, toExec)

			if accumulatedExecutionTime >= t { // reached max execution time, time to commit.
				if err := tx.Commit(); err != nil {
					log.Panic(err)
				}
				resetStCache()
				// sending all ACKs:
				for _, e := range executed {
					e.SetCommitDone()
				}
				executed = executed[:0] //empting the execution slice
				newTimer.Reset(t)
				if StopSQLExecutor {
					log.Println("Stopping the SqlExecutor in time accumulation")
					break outer
				}
			}

		case <-newTimer.C: // timeout of the shot
			if len(executed) > 0 {
				if err := tx.Commit(); err != nil {
					log.Panic(err)
				}
				resetStCache()
				// sending all ACKs:
				for _, e := range executed {
					e.SetCommitDone()
				}
				executed = executed[:0] //empting the execution slice
			}
			newTimer.Reset(t)
			if StopSQLExecutor {
				log.Println("Stopping the SqlExecutor in timeout")
				break outer
			}
		}
	}

}
