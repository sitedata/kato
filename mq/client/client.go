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

package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gridworkz/kato/mq/api/grpc/pb"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	grpcutil "github.com/gridworkz/kato/util/grpc"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

//BuilderTopic
var BuilderTopic = "builder"

//WindowsBuilderTopic
var WindowsBuilderTopic = "windows_builder"

//WorkerTopic
var WorkerTopic = "worker"

//MQClient
type MQClient interface {
	pb.TaskQueueClient
	Close()
	SendBuilderTopic(t TaskStruct) error
}

type mqClient struct {
	pb.TaskQueueClient
	ctx    context.Context
	cancel context.CancelFunc
}

//NewMqClient
func NewMqClient(etcdClientArgs *etcdutil.ClientArgs, defaultserver string) (MQClient, error) {
	ctx, cancel := context.WithCancel(context.Background())
	var conn *grpc.ClientConn
	if etcdClientArgs != nil && etcdClientArgs.Endpoints != nil && len(defaultserver) > 1 {
		c, err := etcdutil.NewClient(ctx, etcdClientArgs)
		if err != nil {
			return nil, err
		}
		r := &grpcutil.GRPCResolver{Client: c}
		b := grpc.RoundRobin(r)
		conn, err = grpc.DialContext(ctx, "/kato/discover/kato_mq", grpc.WithBalancer(b), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		conn, err = grpc.DialContext(ctx, defaultserver, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	}
	cli := pb.NewTaskQueueClient(conn)
	client := &mqClient{
		ctx:    ctx,
		cancel: cancel,
	}
	client.TaskQueueClient = cli
	return client, nil
}

//Close mq grpc client must be closed after uesd
func (m *mqClient) Close() {
	m.cancel()
}

//TaskStruct
type TaskStruct struct {
	Topic    string
	TaskType string
	TaskBody interface{}
}

//BuildTask
func buildTask(t TaskStruct) (*pb.EnqueueRequest, error) {
	var er pb.EnqueueRequest
	taskJSON, err := json.Marshal(t.TaskBody)
	if err != nil {
		logrus.Errorf("tran task json error")
		return &er, err
	}
	er.Topic = t.Topic
	er.Message = &pb.TaskMessage{
		TaskType:   t.TaskType,
		CreateTime: time.Now().Format(time.RFC3339),
		TaskBody:   taskJSON,
		User:       "kato",
	}
	return &er, nil
}

func (m *mqClient) SendBuilderTopic(t TaskStruct) error {
	request, err := buildTask(t)
	if err != nil {
		return fmt.Errorf("create task body error %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(m.ctx, time.Second*5)
	defer cancel()
	_, err = m.TaskQueueClient.Enqueue(ctx, request)
	if err != nil {
		return fmt.Errorf("send enqueue request error %s", err.Error())
	}
	return nil
}
