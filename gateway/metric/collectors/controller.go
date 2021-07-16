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

package collectors

import "github.com/prometheus/client_golang/prometheus"

// PrometheusNamespace default metric namespace
var PrometheusNamespace = "gateway"

// Controller defines base metrics about the rbd-gateway
type Controller struct {
	prometheus.Collector

	activeDomain *prometheus.GaugeVec

	constLabels prometheus.Labels
}

// NewController creates a new prometheus collector for the
// gateway controller operations
func NewController() *Controller {
	constLabels := prometheus.Labels{}
	cm := &Controller{
		activeDomain: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   "nginx",
				Name:        "active_server",
				Help:        "Cumulative number of active servers",
				ConstLabels: constLabels,
			},
			[]string{"type"}),
	}
	return cm
}

// Describe - implements prometheus.Collector
func (cm Controller) Describe(ch chan<- *prometheus.Desc) {
	cm.activeDomain.Describe(ch)
}

// Collect - implements the prometheus.Collector interface.
func (cm Controller) Collect(ch chan<- prometheus.Metric) {
	cm.activeDomain.Collect(ch)
}

// SetServerNum sets the number of active domains
func (cm *Controller) SetServerNum(httpNum, tcpNum int) {
	labels := make(prometheus.Labels, 1)
	labels["type"] = "http"
	cm.activeDomain.With(labels).Set(float64(httpNum))
	labels["type"] = "tcp"
	cm.activeDomain.With(labels).Set(float64(tcpNum))
}
