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
	"time"

	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/discover/config"
	"github.com/gridworkz/kato/monitor/prometheus"
	"github.com/gridworkz/kato/monitor/utils"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

//Webcli
type Webcli struct {
	discover.Callback
	Prometheus      *prometheus.Manager
	sortedEndpoints []string
}

//UpdateEndpoints
func (w *Webcli) UpdateEndpoints(endpoints ...*config.Endpoint) {
	newArr := utils.TrimAndSort(endpoints)

	if utils.ArrCompare(w.sortedEndpoints, newArr) {
		logrus.Debugf("The endpoints is not modify: %s", w.Name())
		return
	}

	w.sortedEndpoints = newArr

	scrape := w.toScrape()
	w.Prometheus.UpdateScrape(scrape)
}

//Error
func (w *Webcli) Error(err error) {
	logrus.Error(err)
}

//Name
func (w *Webcli) Name() string {
	return "webcli"
}

func (w *Webcli) toScrape() *prometheus.ScrapeConfig {
	ts := make([]string, 0, len(w.sortedEndpoints))
	for _, end := range w.sortedEndpoints {
		ts = append(ts, end)
	}
	return &prometheus.ScrapeConfig{
		JobName:        w.Name(),
		ScrapeInterval: model.Duration(time.Minute),
		ScrapeTimeout:  model.Duration(30 * time.Second),
		MetricsPath:    "/metrics",
		HonorLabels:    true,
		ServiceDiscoveryConfig: prometheus.ServiceDiscoveryConfig{
			StaticConfigs: []*prometheus.Group{
				{
					Targets: ts,
					Labels: map[model.LabelName]model.LabelValue{
						"service_name": model.LabelValue(w.Name()),
						"component":    model.LabelValue(w.Name()),
					},
				},
			},
		},
	}
}
