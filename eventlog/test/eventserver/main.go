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

package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pebbe/zmq4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/twinj/uuid"
)

const (
	REQUEST_TIMEOUT = 1000 * time.Millisecond
	MAX_RETRIES     = 3 //  Before we abandon
)

var endpoint string
var coreNum int
var t string

func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&endpoint, "endpoint", "tcp://127.0.0.1:6366", "event server url")
	fs.IntVar(&coreNum, "core", 10, "core number")
	fs.StringVar(&t, "t", "1s", "time interval")
}

func main() {
	AddFlags(pflag.CommandLine)
	pflag.Parse()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	var wait sync.WaitGroup
	d, _ := time.ParseDuration(t)
	for i := 0; i < coreNum; i++ {
		wait.Add(1)
		go func(en string) {
			client, _ := zmq4.NewSocket(zmq4.REQ)
			client.Connect(en)
			defer client.Close()
			id := uuid.NewV4()
		Continuous:
			for {
				request := []string{fmt.Sprintf(`{"event_id":"%s","message":"hello word2","time":"%s"}`, id, time.Now().Format(time.RFC3339))}
				_, err := client.SendMessage(request)
				if err != nil {
					logrus.Error("Send:", err)
				}
				poller := zmq4.NewPoller()
				poller.Add(client, zmq4.POLLIN)
				polled, err := poller.Poll(REQUEST_TIMEOUT)
				if err != nil {
					logrus.Error("Red:", err)
				}
				reply := []string{}
				if len(polled) > 0 {
					reply, err = client.RecvMessage(0)
				} else {
					err = errors.New("Time out")
				}
				if len(reply) > 0 {
					logrus.Info(en, ":", reply[0])
				}
				if err != nil {
					logrus.Error(err)
					break Continuous
				}
				select {
				case <-interrupt:
					break Continuous
				case <-time.Tick(d):
				}
			}
			wait.Done()
		}(endpoint)
	}
	wait.Wait()
}
