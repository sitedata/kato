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

package discover

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gridworkz/kato/builder/exector"
	"github.com/gridworkz/kato/cmd/builder/option"
	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/gridworkz/kato/mq/client"
	"github.com/sirupsen/logrus"
	grpc1 "google.golang.org/grpc"
)

//WTOPIC is builder
const WTOPIC string = "builder"

var healthStatus = make(map[string]string, 1)

//TaskManager task
type TaskManager struct {
	ctx, discoverCtx       context.Context
	cancel, discoverCancel context.CancelFunc
	config                 option.Config
	client                 client.MQClient
	exec                   exector.Manager
	callbackChan           chan *pb.TaskMessage
}

//NewTaskManager return *TaskManager
func NewTaskManager(c option.Config, client client.MQClient, exec exector.Manager) *TaskManager {
	ctx, cancel := context.WithCancel(context.Background())
	discoverCtx, discoverCancel := context.WithCancel(ctx)
	healthStatus["status"] = "health"
	healthStatus["info"] = "builder service health"
	callbackChan := make(chan *pb.TaskMessage, 100)
	taskManager := &TaskManager{
		discoverCtx:    discoverCtx,
		discoverCancel: discoverCancel,
		ctx:            ctx,
		cancel:         cancel,
		config:         c,
		client:         client,
		exec:           exec,
		callbackChan:   callbackChan,
	}
	exec.SetReturnTaskChan(taskManager.callback)
	return taskManager
}

//Start
func (t *TaskManager) Start(errChan chan error) error {
	go t.Do(errChan)
	logrus.Info("start discover success.")
	return nil
}
func (t *TaskManager) callback(task *pb.TaskMessage) {
	ctx, cancel := context.WithCancel(t.ctx)
	defer cancel()
	_, err := t.client.Enqueue(ctx, &pb.EnqueueRequest{
		Topic:   client.BuilderTopic,
		Message: task,
	})
	if err != nil {
		logrus.Errorf("callback task to mq failure %s", err.Error())
	}
	logrus.Infof("The build controller returns an indigestible task(%s) to the messaging system", task.TaskId)
}

//Do it
func (t *TaskManager) Do(errChan chan error) {
	hostName, _ := os.Hostname()
	for {
		select {
		case <-t.discoverCtx.Done():
			return
		default:
			ctx, cancel := context.WithCancel(t.discoverCtx)
			data, err := t.client.Dequeue(ctx, &pb.DequeueRequest{Topic: t.config.Topic, ClientHost: hostName + "-builder"})
			cancel()
			if err != nil {
				if grpc1.ErrorDesc(err) == context.DeadlineExceeded.Error() {
					logrus.Warn(err.Error())
					continue
				}
				if grpc1.ErrorDesc(err) == "context canceled" {
					logrus.Warn("grpc dequeue context canceled")
					healthStatus["status"] = "unusual"
					healthStatus["info"] = "grpc dequeue context canceled"
					return
				}
				if grpc1.ErrorDesc(err) == "context timeout" {
					logrus.Warn(err.Error())
					continue
				}
				if strings.Contains(err.Error(), "there is no connection available") {
					errChan <- fmt.Errorf("message dequeue failure %s", err.Error())
					return
				}
				logrus.Errorf("message dequeue failure %s, will retry", err.Error())
				time.Sleep(time.Second * 2)
				continue
			}
			err = t.exec.AddTask(data)
			if err != nil {
				t.callbackChan <- data
				logrus.Error("add task error:", err.Error())
			}
		}
	}
}

//Stop
func (t *TaskManager) Stop() error {
	t.discoverCancel()
	if err := t.exec.Stop(); err != nil {
		logrus.Errorf("stop task exec manager failure %s", err.Error())
	}
	for len(t.callbackChan) > 0 {
		logrus.Infof("waiting callback chan empty")
		time.Sleep(time.Second * 2)
	}
	logrus.Info("discover manager is stoping.")
	t.cancel()
	if t.client != nil {
		t.client.Close()
	}
	return nil
}

//Component healthCheck
func HealthCheck() map[string]string {
	return healthStatus
}
