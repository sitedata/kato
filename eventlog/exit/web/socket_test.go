// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons 
// to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or 
// substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package web

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"fmt"

	"github.com/gorilla/websocket"
)

func TestWebSocket(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:1233", Path: "/event_log"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer c.Close()
	err = c.WriteMessage(websocket.TextMessage, []byte("event_id=qwertyuiuiosadfkbjasdv"))
	if err != nil {
		t.Fatal(err)
		return
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		t.Log("read:", err)
		return
	}
	if string(message) != "ok" {
		t.Fatal("请求失败")
		return
	}
	done := make(chan struct{})

	defer c.Close()
	defer close(done)

	go func() {
		for {
			select {
			case <-interrupt:
				log.Println("interrupt")
				// To cleanly close a connection, a client should send a close
				// frame and wait for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				c.Close()
				return
			}
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			t.Log("read:", err)
			return
		}
		fmt.Printf("recv: %s \n", message)
	}
}
