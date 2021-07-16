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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pebbe/zmq4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
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
	fs.StringVar(&endpoint, "endpoint", "tcp://127.0.0.1:6363", "docker log server url")
	fs.IntVar(&coreNum, "core", 1, "core number")
	fs.StringVar(&t, "t", "1s", "time interval")
}

func main() {
	AddFlags(pflag.CommandLine)
	pflag.Parse()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	re, _ := http.NewRequest("GET", "http://127.0.0.1:6363/docker-instance?service_id=asdasdadsasdasdassd", nil)
	res, err := http.DefaultClient.Do(re)
	if err != nil {
		logrus.Error(err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return
	}
	status := gjson.Get(string(body), "status")
	host := gjson.Get(string(body), "host")
	if status.String() == "success" {
		endpoint = host.String()
	} else {
		logrus.Error("Failed to get log receiving node." + gjson.Get(string(body), "host").String())
		return
	}
	var wait sync.WaitGroup
	d, _ := time.ParseDuration(t)
	for i := 0; i < coreNum; i++ {
		wait.Add(1)
		go func(en string) {
			client, _ := zmq4.NewSocket(zmq4.PUB)
			client.Monitor("inproc://monitor.rep", zmq4.EVENT_ALL)
			go monitor()
			client.Connect(en)
			defer client.Close()
			id := uuid.NewV4()
		Continuous:
			for {
				request := fmt.Sprintf(`{"event_id":"%s","message":"hello word2","time":"%s"}`, id, time.Now().Format(time.RFC3339))
				client.Send("servasd223123123123", zmq4.SNDMORE)
				_, err := client.Send(request, zmq4.DONTWAIT)
				if err != nil {
					logrus.Error("Send Error:", err)
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

func monitor() {
	mo, _ := zmq4.NewSocket(zmq4.PAIR)
	mo.Connect("inproc://monitor.rep")
	for {
		a, b, c, err := mo.RecvEvent(0)
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Infof("A:%s B:%s C:%d", a, b, c)
	}

}
