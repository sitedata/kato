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

package cluster

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"

	"github.com/gridworkz/kato/eventlog/cluster/connect"
	"github.com/gridworkz/kato/eventlog/cluster/discover"
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/db"

	"golang.org/x/net/context"

	"github.com/gridworkz/kato/eventlog/store"

	"github.com/gridworkz/kato/eventlog/cluster/distribution"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

//Cluster
type Cluster interface {
	//Get a node that accepts logs
	GetSuitableInstance(serviceID string) *discover.Instance
	//Cluster message broadcast
	MessageRadio(...db.ClusterMessage)
	Start() error
	Stop()
	GetInstanceID() string
	GetInstanceHost() string
	Scrape(ch chan<- prometheus.Metric, namespace, exporter string) error
}

//ClusterManager
type ClusterManager struct {
	discover     discover.Manager
	zmqPub       *connect.Pub
	zmqSub       *connect.Sub
	distribution *distribution.Distribution
	Conf         conf.ClusterConf
	log          *logrus.Entry
	storeManager store.Manager
	cancel       func()
	context      context.Context
	etcdClient   *clientv3.Client
}

//NewCluster
func NewCluster(etcdClient *clientv3.Client, conf conf.ClusterConf, log *logrus.Entry, storeManager store.Manager) Cluster {
	ctx, cancel := context.WithCancel(context.Background())
	discover := discover.New(etcdClient, conf.Discover, log.WithField("module", "Discover"))
	distribution := distribution.NewDistribution(etcdClient, conf.Discover, discover, log.WithField("Module", "Distribution"))
	sub := connect.NewSub(conf.PubSub, log.WithField("module", "MessageSubManager"), storeManager, discover, distribution)
	pub := connect.NewPub(conf.PubSub, log.WithField("module", "MessagePubServer"), storeManager, discover)

	return &ClusterManager{
		discover:     discover,
		zmqSub:       sub,
		zmqPub:       pub,
		distribution: distribution,
		Conf:         conf,
		log:          log,
		storeManager: storeManager,
		cancel:       cancel,
		context:      ctx,
		etcdClient:   etcdClient,
	}
}

//Start
func (s *ClusterManager) Start() error {
	if err := s.discover.Run(); err != nil {
		return err
	}
	if err := s.zmqPub.Run(); err != nil {
		return err
	}
	if err := s.zmqSub.Run(); err != nil {
		return err
	}
	if err := s.distribution.Start(); err != nil {
		return err
	}
	go s.monitor()
	return nil
}

//Stop
func (s *ClusterManager) Stop() {
	s.cancel()
	s.distribution.Stop()
	s.zmqPub.Stop()
	s.zmqSub.Stop()
	s.discover.Stop()
}

//GetSuitableInstance
func (s *ClusterManager) GetSuitableInstance(serviceID string) *discover.Instance {
	return s.distribution.GetSuitableInstance(serviceID)
}

//MessageRadio
func (s *ClusterManager) MessageRadio(mes ...db.ClusterMessage) {
	for _, m := range mes {
		s.zmqPub.RadioChan <- m
	}
}

func (s *ClusterManager) GetInstanceID() string {
	return s.discover.GetCurrentInstance().HostID
}
func (s *ClusterManager) GetInstanceHost() string {
	return s.discover.GetCurrentInstance().HostIP.String()
}

func (s *ClusterManager) monitor() {
	tiker := time.Tick(5 * time.Second)
	for {
		messages := s.storeManager.Monitor()
		for _, m := range messages {
			me := db.ClusterMessage{Mode: db.MonitorMessage,
				Data: []byte(fmt.Sprintf("%s,%d,%d", s.GetInstanceID(), m.ServiceSize, m.LogSizePeerM))}
			s.MessageRadio(me)
			m.InstanceID = s.GetInstanceID()
			s.distribution.Update(m)
		}
		select {
		case <-s.context.Done():
			return
		case <-tiker:

		}
	}
}

//Scrape prometheus monitor metrics
func (s *ClusterManager) Scrape(ch chan<- prometheus.Metric, namespace, exporter string) error {
	s.discover.Scrape(ch, namespace, exporter)
	return nil
}
