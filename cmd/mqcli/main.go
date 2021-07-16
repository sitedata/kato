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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/gridworkz/kato/mq/client"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var server string
var topic string
var taskbody string
var taskfile string
var tasktype string
var mode string

func main() {
	AddFlags(pflag.CommandLine)
	pflag.Parse()
	c, err := client.NewMqClient(nil, server)
	if err != nil {
		logrus.Error("new mq client error.", err.Error())
		os.Exit(1)
	}
	defer c.Close()
	if mode == "enqueue" {
		if taskbody == "" && taskfile != "" {
			body, _ := ioutil.ReadFile(taskfile)
			taskbody = string(body)
		}
		fmt.Println("taskbody:" + taskbody)
		re, err := c.Enqueue(context.Background(), &pb.EnqueueRequest{
			Topic: topic,
			Message: &pb.TaskMessage{
				TaskType:   tasktype,
				CreateTime: time.Now().Format(time.RFC3339),
				TaskBody:   []byte(taskbody),
				User:       "gridworkz",
			},
		})
		if err != nil {
			logrus.Error("enqueue error.", err.Error())
			os.Exit(1)
		}
		logrus.Info(re.String())
	}
	if mode == "dequeue" {
		re, err := c.Dequeue(context.Background(), &pb.DequeueRequest{
			Topic:      topic,
			ClientHost: "cli",
		})
		if err != nil {
			logrus.Error("dequeue error.", err.Error())
			os.Exit(1)
		}
		logrus.Info(re.String())
	}

}

func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&server, "server", "127.0.0.1:6300", "mq server")
	fs.StringVar(&topic, "topic", "builder", "mq topic")
	fs.StringVar(&taskbody, "task-body", "", "mq task body")
	fs.StringVar(&taskfile, "task-file", "", "mq task body file")
	fs.StringVar(&tasktype, "task-type", "", "mq task type")
	fs.StringVar(&mode, "mode", "enqueue", "mq task type")
}
