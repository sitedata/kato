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

package connect

import (
	"errors"
	"fmt"
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/db"
	"github.com/gridworkz/kato/eventlog/store"

	"golang.org/x/net/context"

	"sync"

	"github.com/gridworkz/kato/eventlog/cluster/discover"

	"github.com/pebbe/zmq4"
	"github.com/sirupsen/logrus"
)

type Pub struct {
	conf           conf.PubSubConf
	log            *logrus.Entry
	cancel         func()
	context        context.Context
	pubServer      *zmq4.Socket
	pubLock        sync.Mutex
	storemanager   store.Manager
	messageChan    chan [][]byte
	listenErr      chan error
	Closed         chan struct{}
	stopPubMessage bool
	discover       discover.Manager
	instance       *discover.Instance
	RadioChan      chan db.ClusterMessage
}

//NewPub - create zmq pub server
func NewPub(conf conf.PubSubConf, log *logrus.Entry, storeManager store.Manager, discover discover.Manager) *Pub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pub{
		conf:         conf,
		log:          log,
		cancel:       cancel,
		context:      ctx,
		storemanager: storeManager,
		listenErr:    make(chan error),
		Closed:       make(chan struct{}),
		discover:     discover,
		RadioChan:    make(chan db.ClusterMessage, 5),
	}
}

//Run
func (s *Pub) Run() error {
	s.log.Info("message receive server start.")
	pub, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		s.log.Error("create pub zmq socket error.", err.Error())
		return err
	}
	address := fmt.Sprintf("tcp://%s:%d", s.conf.PubBindIP, s.conf.PubBindPort)
	pub.Bind(address)
	s.log.Infof("Message pub server listen %s", address)
	s.pubServer = pub
	s.messageChan = s.storemanager.PubMessageChan()
	if s.messageChan == nil {
		return errors.New("pub log message server can not get store message chan ")
	}
	go s.handleMessage()
	s.registInstance()
	return nil
}

//Stop
func (s *Pub) Stop() {
	if s.instance != nil {
		s.discover.CancellationInstance(s.instance)
	}
	s.cancel()
	<-s.Closed
	s.log.Info("Stop pub message server")
}

func (s *Pub) handleMessage() {
	for !s.stopPubMessage {
		select {
		case msg := <-s.messageChan:
			//s.log.Debugf("Message Pub Server PUB a message %s", string(msg.Content))
			s.pubServer.SendBytes(msg[0], zmq4.SNDMORE)
			s.pubServer.SendBytes(msg[1], 0)
		case m := <-s.RadioChan:
			s.pubServer.SendBytes([]byte(m.Mode), zmq4.SNDMORE)
			s.pubServer.SendBytes(m.Data, 0)
		case <-s.context.Done():
			s.log.Debug("pub message core begin close.")
			s.stopPubMessage = true
			if err := s.pubServer.Close(); err != nil {
				s.log.Warn("Close message pub server error.", err.Error())
			}
			close(s.Closed)
		}
	}
}

func (s *Pub) registInstance() {
	s.instance = s.discover.RegisteredInstance(s.conf.PubBindIP, s.conf.PubBindPort, &s.stopPubMessage)
}
