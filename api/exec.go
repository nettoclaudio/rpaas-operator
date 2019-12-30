// Copyright 2019 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/tsuru/rpaas-operator/config"
	"github.com/tsuru/rpaas-operator/internal/pkg/rpaas"
)

var (
	pingMessage = "@morpheu is boring" // rely on me...
	wsUpgrader  = newWebsocketUpgrader()
)

func newWebsocketUpgrader() websocket.Upgrader {
	cfg := config.Get()
	return websocket.Upgrader{
		HandshakeTimeout: cfg.WebsocketHandshakeTimeout,
		ReadBufferSize:   cfg.WebsocketReadBufferSize,
		WriteBufferSize:  cfg.WebsocketWriteBufferSize,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
}

func remoteExec(c echo.Context) error {
	conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("unable to upgrade to websocket connection: %w", err)
	}

	c.Logger().Infof("New connection with %s estabilshed", conn.RemoteAddr())
	defer conn.Close()

	quit := make(chan bool, 1)
	defer close(quit)

	cfg := config.Get()

	conn.SetCloseHandler(func(code int, text string) error {
		c.Logger().Infof("Closing connection with %s", conn.RemoteAddr())

		quit <- true
		conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""), time.Now().Add(cfg.WebsocketWriteWait))
		return nil
	})

	go func() {
		for {
			select {
			case <-quit:
				c.Logger().Infof("Shutdown the websocket connection gracefully")
				return

			case <-time.After(cfg.WebsocketPingInterval):
				conn.WriteControl(websocket.PingMessage, []byte(pingMessage), time.Now().Add(cfg.WebsocketWriteWait))
			}
		}
	}()

	conn.SetReadDeadline(time.Now().Add(cfg.WebsocketMaxIdleTime))

	conn.SetPongHandler(func(s string) error {
		if s != pingMessage {
			return nil
		}

		conn.SetReadDeadline(time.Now().Add(cfg.WebsocketMaxIdleTime))
		return nil
	})

	return wsRemoteExec(c, conn)
}

func wsRemoteExec(c echo.Context, conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("connection cannot be nil")
	}

	wsReader := &WSReader{conn}

	values := c.QueryParams()
	tty, _ := strconv.ParseBool(values.Get("tty"))

	var stdin io.Reader
	if useStdin, _ := strconv.ParseBool(values.Get("stdin")); useStdin {
		stdin = wsReader
	}

	options := rpaas.ExecOptions{
		Command: values["command"],
		Unit:    values.Get("unit"),
		TTY:     tty,
		Stdin:   stdin,
		Stdout:  wsReader,
		Stderr:  wsReader,
	}

	manager, err := getManager(c)
	if err != nil {
		return nil
	}

	err = manager.Exec(c.Request().Context(), c.Param("instance"), options)
	if err != nil {
		c.Logger().Errorf("exec may not be successful: %v", err)
		// ignoring the error retrun since the connection has been hijacked
		// See: https://github.com/labstack/echo/issues/268
	}

	return nil
}

type WSReader struct {
	*websocket.Conn
}

func (r *WSReader) Read(p []byte) (int, error) {
	messageType, re, err := r.NextReader()
	if err != nil {
		return 0, err
	}

	if messageType != websocket.TextMessage {
		return 0, nil
	}

	return re.Read(p)
}

func (r *WSReader) Write(p []byte) (int, error) {
	return len(p), r.WriteMessage(websocket.TextMessage, p)
}
