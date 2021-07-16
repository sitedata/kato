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

package node

import (
	"strconv"
	"time"

	"github.com/gridworkz/kato/node/nodem/client"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace          = "kato"
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster", "collector_duration_seconds"),
		"cluster_exporter: Duration of a collector scrape.",
		[]string{},
		nil,
	)
	nodeStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster", "node_health"),
		"node_health: Kato node health status.",
		[]string{"node_id", "node_ip", "status", "healthy"},
		nil,
	)
	componentStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster", "component_health"),
		"component_health: Kato node component health status.",
		[]string{"node_id", "node_ip", "component"},
		nil,
	)
)

//Collect prometheus
func (n *Cluster) Collect(ch chan<- prometheus.Metric) {
	begin := time.Now()
	for _, node := range n.GetAllNode() {
		ch <- prometheus.MustNewConstMetric(nodeStatus, prometheus.GaugeValue, func() float64 {
			if node.Status == client.Running && node.NodeStatus.NodeHealth {
				return 0
			}
			return 1
		}(), node.ID, node.InternalIP, node.Status, strconv.FormatBool(node.NodeStatus.NodeHealth))
		for _, con := range node.NodeStatus.Conditions {
			ch <- prometheus.MustNewConstMetric(componentStatus, prometheus.GaugeValue, func() float64 {
				if con.Status == client.ConditionTrue {
					return 0
				}
				return 1
			}(), node.ID, node.InternalIP, string(con.Type))
		}
	}
	duration := time.Since(begin)
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds())
}

//Describe prometheus
func (n *Cluster) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- nodeStatus
	ch <- componentStatus
}
