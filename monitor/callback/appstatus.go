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

package callback

import (
	"strings"
	"time"

	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/discover/config"
	"github.com/gridworkz/kato/monitor/prometheus"
	"github.com/gridworkz/kato/monitor/utils"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// Worker worker monitor
// 127.0.0.1:6369/metrics
type Worker struct {
	discover.Callback
	Prometheus      *prometheus.Manager
	sortedEndpoints []string
}

//UpdateEndpoints update endpoint
func (e *Worker) UpdateEndpoints(endpoints ...*config.Endpoint) {
	// Register with v3 API and return to json format test, so you have to deal with it in advance
	newEndpoints := make([]*config.Endpoint, 0, len(endpoints))
	for _, end := range endpoints {
		newEnd := *end
		newEndpoints = append(newEndpoints, &newEnd)
	}

	for i, end := range endpoints {
		newEndpoints[i].URL = gjson.Get(end.URL, "Addr").String()
	}

	newArr := utils.TrimAndSort(newEndpoints)

	// change port
	for i, end := range newArr {
		newArr[i] = strings.Split(end, ":")[0] + ":6369"
	}

	if utils.ArrCompare(e.sortedEndpoints, newArr) {
		logrus.Debugf("The endpoints is not modify: %s", e.Name())
		return
	}

	e.sortedEndpoints = newArr

	scrape := e.toScrape()
	e.Prometheus.UpdateScrape(scrape)
}

func (e *Worker) Error(err error) {
	logrus.Error(err)
}

//Name return name
func (e *Worker) Name() string {
	return "worker"
}

func (e *Worker) toScrape() *prometheus.ScrapeConfig {
	ts := make([]string, 0, len(e.sortedEndpoints))
	for _, end := range e.sortedEndpoints {
		ts = append(ts, end)
	}

	return &prometheus.ScrapeConfig{
		JobName:        e.Name(),
		ScrapeInterval: model.Duration(3 * time.Minute),
		ScrapeTimeout:  model.Duration(60 * time.Second),
		MetricsPath:    "/metrics",
		ServiceDiscoveryConfig: prometheus.ServiceDiscoveryConfig{
			StaticConfigs: []*prometheus.Group{
				{
					Targets: ts,
					Labels: map[model.LabelName]model.LabelValue{
						"component":    model.LabelValue(e.Name()),
						"service_name": model.LabelValue(e.Name()),
					},
				},
			},
		},
		MetricRelabelConfigs: []*prometheus.RelabelConfig{
			{
				SourceLabels: model.LabelNames{model.LabelName("tenant_id")},
				TargetLabel:  "namespace",
			},
		},
	}
}
