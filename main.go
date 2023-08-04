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
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//"github.com/labstack/gommon/log"

	"github.com/joanjib/fullness-server/internal/config"
	"github.com/joanjib/fullness-server/internal/executor"
	"github.com/joanjib/fullness-server/internal/sqlExecutor"
	"github.com/joanjib/fullness-server/internal/yaml"
)

var upgrader = websocket.Upgrader{}

type Request struct {
	Action string
	Input  []interface{}
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// actions the server can handle.

var executors = sync.Pool{
	New: func() interface{} { return executor.New() },
}

func wsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	log.Println("Client connected to WS")
	for {
		// code for testing sql execution generic.
		var r Request
		err = ws.ReadJSON(&r)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				return nil
			} else {
				log.Panic(err)
			}
		}

		ex := executors.Get().(*executor.Executor)
		defer executors.Put(ex)
		ex.SetInput(r.Input)
		ex.IsQuery = strings.HasPrefix(r.Action, "query:")
		if ex.IsQuery {
			ex.Action = r.Action[6:]
		} else {
			ex.Action = r.Action
		}

		sqlExecutor.ExecutionQueue <- ex

		<-ex.DataAvailable
		err = ws.WriteJSON(&ex.Output)
		if err != nil {
			c.Logger().Error(err)
		}

		<-ex.CommitDone
		// Write
		err = ws.WriteJSON("ACK")
		if err != nil {
			c.Logger().Error(err)
		}

		/*
			// Read
			err = ws.ReadJSON(&user)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Println(user.Name, "<->", user.Email)
		*/
	}
}

func main() {
	// environtment variables loading.
	// FULLNESS_SERVER_CONFIG  <- environtment variable for the general configuration file for the server.
	serverConfigFile := os.Getenv("FULLNESS_SERVER_CONFIG")
	var config config.Config
	yaml.LoadYml(serverConfigFile, &config)
	// configuration loaded into the variable config

	// SQL executor initialization:
	sqlExecutor.Init(config.ExecutionQueueSize)
	go sqlExecutor.SQLExecutor(config.DatabaseFile, config.SqlStatementsDir, config.CommitInterval, config.SchemaFile)

	e := echo.New()
	//e.Logger.SetLevel(log.INFO)
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./public")
	e.GET("/ws", wsHandler)

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sqlExecutor.StopSQLExecutor = true
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	} else {
		time.Sleep(4 * time.Second)
		log.Println("Server stopped gracefully")
	}

}
