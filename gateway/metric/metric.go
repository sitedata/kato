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

package metric

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/gridworkz/kato/gateway/metric/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

// Collector defines the interface for a metric collector
type Collector interface {
	Start()
	Stop()
	SetHosts(sets.String)
	SetServerNum(httpNum, tcpNum int)
	RemoveHostMetric([]string)
}

type collector struct {
	registry          *prometheus.Registry
	socket            *collectors.SocketCollector
	gatewayController *collectors.Controller
	nginxCmd          *collectors.NginxCmdMetric
}

// NewCollector creates a new metric collector the for ingress controller
func NewCollector(gatewayHost string, registry *prometheus.Registry) (Collector, error) {
	ic := collectors.NewController()
	socketCollector, err := collectors.NewSocketCollector(gatewayHost, true)
	if err != nil {
		return nil, fmt.Errorf("create socket collector failure %s", err.Error())
	}
	return Collector(&collector{
		gatewayController: ic,
		socket:            socketCollector,
		registry:          registry,
		nginxCmd:          &collectors.NginxCmdMetric{},
	}), nil
}

func (c *collector) Start() {
	c.registry.MustRegister(c.gatewayController)
	c.registry.MustRegister(c.socket)
	c.registry.MustRegister(c.nginxCmd)
	go c.socket.Start()
}

func (c *collector) Stop() {
	c.registry.Unregister(c.gatewayController)
	c.registry.Unregister(c.socket)
	c.registry.Unregister(c.nginxCmd)
}

func (c *collector) SetServerNum(httpNum, tcpNum int) {
	c.gatewayController.SetServerNum(httpNum, tcpNum)
}

func (c *collector) SetHosts(hosts sets.String) {
	c.socket.SetHosts(hosts)
}

//RemoveHostMetric -
func (c *collector) RemoveHostMetric(hosts []string) {
	c.socket.RemoveMetrics(hosts, c.registry)
}
