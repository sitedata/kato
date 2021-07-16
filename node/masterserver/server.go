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

package masterserver

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gridworkz/kato/node/masterserver/monitor"

	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/node/kubecache"
	"github.com/gridworkz/kato/node/nodem/client"

	"github.com/gridworkz/kato/node/core/config"
	"github.com/gridworkz/kato/node/core/store"
	"github.com/gridworkz/kato/node/masterserver/node"
)

//MasterServer
type MasterServer struct {
	*store.Client
	*client.HostNode
	Cluster          *node.Cluster
	ctx              context.Context
	cancel           context.CancelFunc
	datacenterConfig *config.DataCenterConfig
	clusterMonitor   monitor.Manager
}

//NewMasterServer
func NewMasterServer(modelnode *client.HostNode, kubecli kubecache.KubeClient) (*MasterServer, error) {
	datacenterConfig := config.GetDataCenterConfig()
	ctx, cancel := context.WithCancel(context.Background())
	nodecluster := node.CreateCluster(kubecli, modelnode, datacenterConfig)
	clusterMonitor, err := monitor.CreateManager(nodecluster)
	if err != nil {
		cancel()
		return nil, err
	}
	ms := &MasterServer{
		Client:           store.DefalutClient,
		HostNode:         modelnode,
		Cluster:          nodecluster,
		ctx:              ctx,
		cancel:           cancel,
		datacenterConfig: datacenterConfig,
		clusterMonitor:   clusterMonitor,
	}
	return ms, nil
}

//Start master node
func (m *MasterServer) Start(errchan chan error) error {
	m.datacenterConfig.Start()
	if err := m.Cluster.Start(errchan); err != nil {
		logrus.Error("node cluster start error,", err.Error())
		return err
	}
	return m.clusterMonitor.Start(errchan)
}

//Stop
func (m *MasterServer) Stop(i interface{}) {
	if m.Cluster != nil {
		m.Cluster.Stop(i)
	}
	if m.clusterMonitor != nil {
		m.clusterMonitor.Stop()
	}
	m.cancel()
}

//GetRegistry
func (m *MasterServer) GetRegistry() *prometheus.Registry {
	return m.clusterMonitor.GetRegistry()
}
