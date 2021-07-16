// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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

package monitor

import (
	"github.com/gridworkz/kato/node/masterserver/node"

	"github.com/prometheus/client_golang/prometheus"
)

//Manager
type Manager interface {
	Start(errchan chan error) error
	Stop() error
	GetRegistry() *prometheus.Registry
}

type manager struct {
	clusterExporterRestry *prometheus.Registry
	cluster               *node.Cluster
}

//CreateManager
func CreateManager(cluster *node.Cluster) (Manager, error) {
	clusterRegistry := prometheus.NewRegistry()
	manage := &manager{
		clusterExporterRestry: clusterRegistry,
		cluster:               cluster,
	}
	return manage, nil
}

func (m *manager) Start(errchan chan error) error {
	return m.clusterExporterRestry.Register(m.cluster)
}

func (m *manager) Stop() error {
	return nil
}

func (m *manager) GetRegistry() *prometheus.Registry {
	return m.clusterExporterRestry
}
