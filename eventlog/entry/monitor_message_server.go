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
	"github.com/gridworkz/kato/eventlog/store"
	"time"

	"golang.org/x/net/context"

	"sync"

	zmq4 "github.com/pebbe/zmq4"
	"github.com/sirupsen/logrus"
)

//MonitorMessageServer - monitor real-time data acceptance service
type MonitorMessageServer struct {
	conf               conf.MonitorMessageServerConf
	log                *logrus.Entry
	cancel             func()
	context            context.Context
	server             *zmq4.Socket
	storemanager       store.Manager
	messageChan        chan [][]byte
	listenErr          chan error
	serverLock         sync.Mutex
	stopReceiveMessage bool
}

//NewMonitorMessageServer - create zmq sub
func NewMonitorMessageServer(conf conf.MonitorMessageServerConf, log *logrus.Entry, storeManager store.Manager) (*MonitorMessageServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &MonitorMessageServer{
		conf:         conf,
		log:          log,
		cancel:       cancel,
		context:      ctx,
		storemanager: storeManager,
		listenErr:    make(chan error),
	}
	server, err := zmq4.NewSocket(zmq4.SUB)
	server.SetSubscribe(conf.SubSubscribe)
	if err != nil {
		s.log.Error("create rep zmq socket error.", err.Error())
		return nil, err
	}
	for _, add := range conf.SubAddress {
		server.Connect(add)
		s.log.Infof("Monitor message server sub %s", add)
	}
	s.server = server
	s.messageChan = s.storemanager.MonitorMessageChan()
	if s.messageChan == nil {
		return nil, errors.New("receive monitor message server can not get store message chan ")
	}
	return s, nil
}

//Serve
func (s *MonitorMessageServer) Serve() {
	s.handleMessage()
}

//Stop
func (s *MonitorMessageServer) Stop() {
	s.cancel()
	s.log.Info("receive event message server stop")
}

func (s *MonitorMessageServer) handleMessage() {
	chQuit := make(chan interface{})
	chErr := make(chan error, 2)
	channel := make(chan [][]byte, s.conf.CacheMessageSize)
	newServerListen := func(sock *zmq4.Socket, channel chan [][]byte) {
		socketHandler := func(state zmq4.State) error {
			msgs, err := sock.RecvMessageBytes(0)
			if err != nil {
				s.log.Error("server receive message error.", err.Error())
				return err
			}
			s.messageChan <- msgs
			return nil
		}
		quitHandler := func(interface{}) error {
			close(channel)
			s.log.Infof("Event message receive Server quit.")
			return nil
		}
		reactor := zmq4.NewReactor()
		reactor.AddSocket(sock, zmq4.POLLIN, socketHandler)
		reactor.AddChannel(chQuit, 1, quitHandler)
		err := reactor.Run(100 * time.Millisecond)
		chErr <- err
	}
	go newServerListen(s.server, channel)

	func() {
		for !s.stopReceiveMessage {
			select {
			case msg := <-channel:
				s.messageChan <- msg
			case <-s.context.Done():
				s.log.Debug("handle message core begin close.")
				close(chQuit)
				s.stopReceiveMessage = true
				// close(s.messageChan)
			}
		}
	}()
	s.log.Info("Handle message core stop.")
}

type event struct {
	Name   string        `json:"name"`
	Data   []interface{} `json:"data"`
	Update string        `json:"update_time"`
}

//ListenError listen error chan
func (s *MonitorMessageServer) ListenError() chan error {
	return s.listenErr
}
