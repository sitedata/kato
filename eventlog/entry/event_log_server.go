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

package entry

import (
	"errors"

	"github.com/gridworkz/kato/eventlog/conf"
	grpcserver "github.com/gridworkz/kato/eventlog/entry/grpc/server"
	"github.com/gridworkz/kato/eventlog/store"

	"golang.org/x/net/context"

	"sync"

	"github.com/sirupsen/logrus"
)

//EventLogServer - log acceptance service
type EventLogServer struct {
	conf               conf.EventLogServerConf
	log                *logrus.Entry
	cancel             func()
	context            context.Context
	storemanager       store.Manager
	messageChan        chan []byte
	listenErr          chan error
	serverLock         sync.Mutex
	stopReceiveMessage bool
	eventRPCServer     *grpcserver.EventLogRPCServer
}

//NewEventLogServer - create zmq server 
func NewEventLogServer(conf conf.EventLogServerConf, log *logrus.Entry, storeManager store.Manager) (*EventLogServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &EventLogServer{
		conf:         conf,
		log:          log,
		cancel:       cancel,
		context:      ctx,
		storemanager: storeManager,
		listenErr:    make(chan error),
	}

	//grpc service
	eventRPCServer := grpcserver.NewServer(conf, log, storeManager, s.listenErr)
	s.messageChan = s.storemanager.ReceiveMessageChan()
	if s.messageChan == nil {
		return nil, errors.New("receive log message server can not get store message chan ")
	}
	s.eventRPCServer = eventRPCServer
	return s, nil
}

//Serve
func (s *EventLogServer) Serve() {
	s.eventRPCServer.Start()
}

//Stop
func (s *EventLogServer) Stop() {
	s.cancel()
	s.eventRPCServer.Stop()
	s.log.Info("receive event message server stop")
}

//ListenError listen error chan
func (s *EventLogServer) ListenError() chan error {
	return s.listenErr
}
