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

	"github.com/gridworkz/kato/api/handler"
	"github.com/prometheus/client_golang/prometheus"
)

// Metric name parts.
const (
	// Namespace for all metrics.
	namespace = "rbd_api"
	// Subsystem(s).
	exporter = "exporter"
)

//NewExporter new exporter
func NewExporter() *Exporter {
	return &Exporter{
		apiRequest: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "api_request",
			Help:      "kato cluster api request metric",
		}, []string{"code", "path"}),
		tenantLimit: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "tenant_memory_limit",
			Help:      "kato tenant memory limit",
		}, []string{"tenant_id", "namespace"}),
		clusterMemoryTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "cluster_memory_total",
			Help:      "kato cluster memory total",
		}),
		clusterCPUTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "cluster_cpu_total",
			Help:      "kato cluster cpu total",
		}),
	}
}

//Exporter exporter
type Exporter struct {
	apiRequest         *prometheus.CounterVec
	tenantLimit        *prometheus.GaugeVec
	clusterCPUTotal    prometheus.Gauge
	clusterMemoryTotal prometheus.Gauge
}

//RequestInc request inc
func (e *Exporter) RequestInc(code int, path string) {
	e.apiRequest.WithLabelValues(fmt.Sprintf("%d", code), path).Inc()
}

//Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.apiRequest.Collect(ch)
	// tenant limit value
	tenants, _ := handler.GetTenantManager().GetTenants("")
	for _, t := range tenants {
		e.tenantLimit.WithLabelValues(t.UUID, t.UUID).Set(float64(t.LimitMemory))
	}
	// cluster memory
	resource := handler.GetTenantManager().GetClusterResource()
	if resource != nil {
		e.clusterMemoryTotal.Set(float64(resource.AllMemory))
		e.clusterCPUTotal.Set(float64(resource.AllCPU))
	}
	e.tenantLimit.Collect(ch)
	e.clusterMemoryTotal.Collect(ch)
	e.clusterCPUTotal.Collect(ch)
}
