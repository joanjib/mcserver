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
package executor

import (
	"database/sql"
	"log"
	//"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

// Fields common to all Executor instances
type Executor struct {
	DataAvailable  chan struct{}
	CommitDone     chan struct{}
	IsServerAction bool
	IsClientAction bool
	Action         string
	IsQuery        bool
	Input          []interface{}
	Output         [][]interface{}
	OutputPtr      [][]interface{}
	LastInsertedId int64
	RowsAffected   int64
	Error          string
}

// Executor implementation Interface
type ExecutorImpl interface {

	// Done in the echo goroutine
	// receive a web socket as an input to fullfill the internal struct of the instance.
	SetInput(input []interface{})

	// Done in the sql execution goroutine
	// Passwd a reference to the current transaction, passed a reference to the dictionary with all prepared statements, and the one that acts as cache
	// sp: savepoint idenfiticator: used to rollback or release the current execution
	Execute(stmt *sql.Stmt, tx *sql.Tx, sp string)

	// When client action is requested:
	// Done in the echo goroutine
	// sends the output to the client
	//SendOutput(ws *websocket.Conn)

	// Done in the echo goroutine
	// sends transaction actually commited.
	//SendACK(ws *websocket.Conn)

	// Done in the echo goroutine : it has available all global resources
	// When server action is requered:
	//PostExecution()

	SetDataAvailable()
	SetCommitDone()
	GetAction() string
}

func New() *Executor {
	e := new(Executor)
	e.DataAvailable = make(chan struct{})
	e.CommitDone = make(chan struct{})

	// inicilizing output matrix

	e.Output = make([][]interface{}, 10, 10)
	e.OutputPtr = make([][]interface{}, 10, 10)
	for i, _ := range e.Output {
		e.Output[i] = make([]interface{}, 10, 10)
		e.OutputPtr[i] = make([]interface{}, 10, 10)
		for j, _ := range e.Output[i] {
			e.OutputPtr[i][j] = &e.Output[i][j]
		}
	}
	return e
}

func (e *Executor) SetDataAvailable() {
	e.DataAvailable <- struct{}{}
}

func (e *Executor) SetCommitDone() {
	e.CommitDone <- struct{}{}
}
func (e *Executor) SetInput(input []interface{}) {
	e.Input = input
}

func (e *Executor) GetAction() string {
	return e.Action
}
func (e *Executor) Execute(stmt *sql.Stmt, tx *sql.Tx, sp string) {

	if e.IsQuery {
		rows, err := stmt.Query(e.Input...)
		if err != nil {
			e.Error = err.Error()
		} else {
			defer rows.Close()
			cols, _ := rows.Columns()
			lenCol := len(cols)
			i := 1
			for rows.Next() {
				e.Output = e.Output[:i]
				e.OutputPtr = e.OutputPtr[:i]
				e.Output[i-1] = e.Output[i-1][:lenCol]
				e.OutputPtr[i-1] = e.OutputPtr[i-1][:lenCol]
				if err := rows.Scan(e.OutputPtr[i-1]...); err != nil {
					log.Panic(err)
				}
				i++
			}
			// Check for errors from iterating over rows.
			if err := rows.Err(); err != nil {
				log.Panic(err)
			}

		}
	} else {
		_, err := tx.Exec("savepoint " + sp)
		if err != nil {
			log.Panic(err)
		}
		res, err := stmt.Exec(e.Input...)
		if err != nil {
			e.Error = err.Error()
			_, err := tx.Exec("rollback to " + sp)
			if err != nil {
				log.Panic(err)
			}
		} else {
			//TODO: check the errors
			e.LastInsertedId, _ = res.LastInsertId()
			e.RowsAffected, _ = res.RowsAffected()
			_, err := tx.Exec("release " + sp)
			if err != nil {
				log.Panic(err)
			}
		}
	}
}
