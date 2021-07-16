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

package server

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/store"

	"github.com/gridworkz/kato/eventlog/entry/grpc/pb"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type EventLogRPCServer struct {
	conf         conf.EventLogServerConf
	log          *logrus.Entry
	cancel       func()
	context      context.Context
	storemanager store.Manager
	messageChan  chan []byte
	listenErr    chan error
	lis          net.Listener
}

//NewServer
func NewServer(conf conf.EventLogServerConf, log *logrus.Entry, storeManager store.Manager, listenErr chan error) *EventLogRPCServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventLogRPCServer{
		conf:         conf,
		log:          log,
		storemanager: storeManager,
		context:      ctx,
		cancel:       cancel,
		messageChan:  storeManager.ReceiveMessageChan(),
		listenErr:    listenErr,
	}
}

//Start grpc server
func (s *EventLogRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.conf.BindIP, s.conf.BindPort))
	if err != nil {
		logrus.Errorf("failed to listen: %v", err)
		return err
	}
	s.lis = lis
	server := grpc.NewServer()
	pb.RegisterEventLogServer(server, s)
	// Register reflection service on gRPC server.
	reflection.Register(server)
	s.log.Infof("event message grpc server listen %s:%d", s.conf.BindIP, s.conf.BindPort)
	if err := server.Serve(lis); err != nil {
		s.log.Error("event log api grpc listen error.", err.Error())
		s.listenErr <- err
	}
	return nil
}

//Stop
func (s *EventLogRPCServer) Stop() {
	s.cancel()
	// if s.lis != nil {
	// 	s.lis.Close()
	// }
}

//Log impl EventLogServerServer
func (s *EventLogRPCServer) Log(stream pb.EventLog_LogServer) error {
	for {
		select {
		case <-s.context.Done():
			if err := stream.SendAndClose(&pb.Reply{Status: "success", Message: "server closed"}); err != nil {
				return err
			}
			return nil
		default:
		}
		log, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				s.log.Error("receive log error:", err.Error())
				if err := stream.SendAndClose(&pb.Reply{Status: "success"}); err != nil {
					return err
				}
				return nil
			}
			return err
		}
		select {
		case s.messageChan <- log.Log:
		default:
		}
	}
}
