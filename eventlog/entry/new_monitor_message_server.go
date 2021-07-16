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
	"fmt"
	"net"
	"time"

	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/store"

	"golang.org/x/net/context"

	"sync"

	"github.com/sirupsen/logrus"
)

//NMonitorMessageServer - new performance analysis real-time data acceptance service
type NMonitorMessageServer struct {
	conf               conf.NewMonitorMessageServerConf
	log                *logrus.Entry
	cancel             func()
	context            context.Context
	storemanager       store.Manager
	messageChan        chan []byte
	listenErr          chan error
	serverLock         sync.Mutex
	stopReceiveMessage bool
	listener           *net.UDPConn
}

//NewNMonitorMessageServer - create UDP server
func NewNMonitorMessageServer(conf conf.NewMonitorMessageServerConf, log *logrus.Entry, storeManager store.Manager) (*NMonitorMessageServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	s := &NMonitorMessageServer{
		conf:         conf,
		log:          log,
		cancel:       cancel,
		context:      ctx,
		storemanager: storeManager,
		listenErr:    make(chan error),
	}
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(conf.ListenerHost), Port: conf.ListenerPort})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	log.Infof("UDP Server Listener: %s", listener.LocalAddr().String())
	s.listener = listener
	s.messageChan = s.storemanager.NewMonitorMessageChan()
	if s.messageChan == nil {
		return nil, errors.New("receive monitor message server can not get store message chan ")
	}
	return s, nil
}

//Serve
func (s *NMonitorMessageServer) Serve() {
	s.handleMessage()
}

//Stop
func (s *NMonitorMessageServer) Stop() {
	s.cancel()
	s.log.Info("receive new monitor message server stop")
}

func (s *NMonitorMessageServer) handleMessage() {
	buf := make([]byte, 65535)
	defer s.listener.Close()
	s.log.Infoln("start receive monitor message by udp")
	for {
		n, _, err := s.listener.ReadFromUDP(buf)
		if err != nil {
			logrus.Errorf("read new monitor message from udp error,%s", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}
		// fix issues https://github.com/golang/go/issues/35725
		message := make([]byte, n)
		copy(message, buf[0:n])
		s.messageChan <- message
	}
}

//ListenError listen error chan
func (s *NMonitorMessageServer) ListenError() chan error {
	return s.listenErr
}
