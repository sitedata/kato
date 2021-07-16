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

package mq

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gridworkz/kato/cmd/mq/option"
	"github.com/gridworkz/kato/mq/client"

	"golang.org/x/net/context"

	etcdutil "github.com/gridworkz/kato/util/etcd"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

//ActionMQ
type ActionMQ interface {
	Enqueue(context.Context, string, string) error
	Dequeue(context.Context, string) (string, error)
	TopicIsExist(string) bool
	GetAllTopics() []string
	Start() error
	Stop() error
	MessageQueueSize(topic string) int64
}

// EnqueueNumber
var EnqueueNumber float64 = 0

// DequeueNumber
var DequeueNumber float64 = 0

// NewActionMQ
func NewActionMQ(ctx context.Context, c option.Config) ActionMQ {
	etcdQueue := etcdQueue{
		config: c,
		ctx:    ctx,
		queues: make(map[string]string),
	}
	return &etcdQueue
}

type etcdQueue struct {
	config     option.Config
	ctx        context.Context
	queues     map[string]string
	queuesLock sync.Mutex
	client     *clientv3.Client
}

func (e *etcdQueue) Start() error {
	logrus.Debug("etcd message queue client starting")
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints:   e.config.EtcdEndPoints,
		CaFile:      e.config.EtcdCaFile,
		CertFile:    e.config.EtcdCertFile,
		KeyFile:     e.config.EtcdKeyFile,
		DialTimeout: time.Duration(e.config.EtcdTimeout) * time.Second,
	}
	cli, err := etcdutil.NewClient(context.Background(), etcdClientArgs)
	if err != nil {
		etcdutil.HandleEtcdError(err)
		return err
	}
	e.client = cli
	topics := os.Getenv("topics")
	if topics != "" {
		ts := strings.Split(topics, ",")
		for _, t := range ts {
			e.registerTopic(t)
		}
	}
	e.registerTopic(client.BuilderTopic)
	e.registerTopic(client.WindowsBuilderTopic)
	e.registerTopic(client.WorkerTopic)
	logrus.Info("etcd message queue client started success")
	return nil
}

//registerTopic
func (e *etcdQueue) registerTopic(topic string) {
	e.queuesLock.Lock()
	defer e.queuesLock.Unlock()
	e.queues[topic] = topic
}

func (e *etcdQueue) TopicIsExist(topic string) bool {
	e.queuesLock.Lock()
	defer e.queuesLock.Unlock()
	_, ok := e.queues[topic]
	return ok
}
func (e *etcdQueue) GetAllTopics() []string {
	var topics []string
	for k := range e.queues {
		topics = append(topics, k)
	}
	return topics
}

func (e *etcdQueue) Stop() error {
	if e.client != nil {
		e.client.Close()
	}
	return nil
}
func (e *etcdQueue) queueKey(topic string) string {
	return e.config.EtcdPrefix + "/" + topic
}
func (e *etcdQueue) Enqueue(ctx context.Context, topic, value string) error {
	EnqueueNumber++
	queue := etcdutil.NewQueue(ctx, e.client, e.queueKey(topic))
	return queue.Enqueue(value)
}

func (e *etcdQueue) Dequeue(ctx context.Context, topic string) (string, error) {
	DequeueNumber++
	queue := etcdutil.NewQueue(ctx, e.client, e.queueKey(topic))
	return queue.Dequeue()
}

func (e *etcdQueue) MessageQueueSize(topic string) int64 {
	ctx, cancel := context.WithCancel(e.ctx)
	defer cancel()
	res, err := e.client.Get(ctx, e.queueKey(topic), clientv3.WithPrefix())
	if err != nil {
		logrus.Errorf("get message queue size failure %s", err.Error())
	}
	if res != nil {
		return res.Count
	}
	return 0
}
