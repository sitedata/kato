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
	"time"

	"github.com/gridworkz/kato/cmd/worker/option"
	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/gridworkz/kato/mq/client"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/gridworkz/kato/worker/appm/controller"
	"github.com/gridworkz/kato/worker/appm/store"
	"github.com/gridworkz/kato/worker/discover/model"
	"github.com/gridworkz/kato/worker/gc"
	"github.com/gridworkz/kato/worker/handle"
	"github.com/sirupsen/logrus"
	grpc1 "google.golang.org/grpc"
)

var healthStatus = make(map[string]string, 1)

//TaskNum exec task number
var TaskNum float64

//TaskError exec error task number
var TaskError float64

//TaskManager task
type TaskManager struct {
	ctx           context.Context
	cancel        context.CancelFunc
	config        option.Config
	handleManager *handle.Manager
	client        client.MQClient
}

//NewTaskManager return *TaskManager
func NewTaskManager(cfg option.Config,
	store store.Storer,
	controllermanager *controller.Manager,
	garbageCollector *gc.GarbageCollector) *TaskManager {

	ctx, cancel := context.WithCancel(context.Background())
	handleManager := handle.NewManager(ctx, cfg, store, controllermanager, garbageCollector)
	healthStatus["status"] = "health"
	healthStatus["info"] = "worker service health"
	return &TaskManager{
		ctx:           ctx,
		cancel:        cancel,
		config:        cfg,
		handleManager: handleManager,
	}
}

//Start start
func (t *TaskManager) Start() error {
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: t.config.EtcdEndPoints,
		CaFile:    t.config.EtcdCaFile,
		CertFile:  t.config.EtcdCertFile,
		KeyFile:   t.config.EtcdKeyFile,
	}
	client, err := client.NewMqClient(etcdClientArgs, t.config.MQAPI)
	if err != nil {
		logrus.Errorf("new Mq client error, %v", err)
		healthStatus["status"] = "unusual"
		healthStatus["info"] = fmt.Sprintf("new Mq client error, %v", err)
		return err
	}
	t.client = client
	go t.Do()
	logrus.Info("start discover success.")
	return nil
}

//Do do
func (t *TaskManager) Do() {
	logrus.Info("start receive task from mq")
	hostname, _ := os.Hostname()
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			data, err := t.client.Dequeue(t.ctx, &pb.DequeueRequest{Topic: client.WorkerTopic, ClientHost: hostname + "-worker"})
			if err != nil {
				if grpc1.ErrorDesc(err) == context.DeadlineExceeded.Error() {
					continue
				}
				if grpc1.ErrorDesc(err) == "context canceled" {
					logrus.Info("receive task core context canceled")
					healthStatus["status"] = "unusual"
					healthStatus["info"] = "receive task core context canceled"
					return
				}
				if grpc1.ErrorDesc(err) == "context timeout" {
					continue
				}
				logrus.Error("receive task error.", err.Error())
				time.Sleep(time.Second * 2)
				continue
			}
			logrus.Debugf("receive a task: %v", data)
			transData, err := model.TransTask(data)
			if err != nil {
				logrus.Error("trans mq msg data error ", err.Error())
				continue
			}
			rc := t.handleManager.AnalystToExec(transData)
			if rc != nil && rc != handle.ErrCallback {
				logrus.Warningf("execute task: %v", rc)
				TaskError++
			} else if rc != nil && rc == handle.ErrCallback {
				logrus.Errorf("err callback; analyst to exet: %v", rc)
				ctx, cancel := context.WithCancel(t.ctx)
				reply, err := t.client.Enqueue(ctx, &pb.EnqueueRequest{
					Topic:   client.WorkerTopic,
					Message: data,
				})
				cancel()
				logrus.Debugf("retry send task to mq ,reply is %v", reply)
				if err != nil {
					logrus.Errorf("enqueue task %v to mq topic %v Error", data, client.WorkerTopic)
					continue
				}
				//if handle is waiting, sleep 3 second
				time.Sleep(time.Second * 3)
			} else {
				TaskNum++
			}
		}
	}
}

//Stop stop
func (t *TaskManager) Stop() error {
	logrus.Info("discover manager is stoping.")
	t.cancel()
	if t.client != nil {
		t.client.Close()
	}
	return nil
}

//HealthCheck health check
func HealthCheck() map[string]string {
	return healthStatus
}
