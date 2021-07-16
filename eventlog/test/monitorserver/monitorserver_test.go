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

package monitorserver

import (
	"testing"

	"github.com/pebbe/zmq4"
	"github.com/sirupsen/logrus"
)

var urlData = `
2017-05-19 11:33:34 APPS SumTimeByUrl [{"tenant":"o2o","service":"zzcplus","url":"/active/js/wx_share.js","avgtime":"1.453","sumtime":"1.453","counts":"1"}]
`
var newMonitorMessage = `
[{"ServiceID":"test",
	"Port":"5000",
	"MessageType":"http",
	"Key":"/test",
	"CumulativeTime":0.1,
	"AverageTime":0.1,
	"MaxTime":0.1,
	"Count":1,
	"AbnormalCount":0}
,{"ServiceID":"test",
	"Port":"5000",
	"MessageType":"http",
	"Key":"/test2",
	"CumulativeTime":0.36,
	"AverageTime":0.18,
	"MaxTime":0.2,
	"Count":2,
	"AbnormalCount":2}
]
`

func BenchmarkMonitorServer(t *testing.B) {
	client, _ := zmq4.NewSocket(zmq4.PUB)
	// client.Monitor("inproc://monitor.rep", zmq4.EVENT_ALL)
	// go monitor()
	client.Bind("tcp://0.0.0.0:9442")
	defer client.Close()
	var size int64
	for i := 0; i < t.N; i++ {
		client.Send("ceptop", zmq4.SNDMORE)
		_, err := client.Send(urlData, zmq4.DONTWAIT)
		if err != nil {
			logrus.Error("Send Error:", err)
		}
		size++
	}
	logrus.Info(size)
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
