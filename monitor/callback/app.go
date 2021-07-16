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
	"github.com/gridworkz/kato/util/watch"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// App refers to app runtime information, which comes from node on all child nodes
// 127.0.0.1:6100/app/metrics
type App struct {
	discover.Callback
	Prometheus      *prometheus.Manager
	sortedEndpoints []string

	endpoints []*config.Endpoint
}

//UpdateEndpoints update endpoint
func (e *App) UpdateEndpoints(endpoints ...*config.Endpoint) {
	newArr := utils.TrimAndSort(endpoints)

	if utils.ArrCompare(e.sortedEndpoints, newArr) {
		logrus.Debugf("The endpoints is not modify: %s", e.Name())
		return
	}

	e.sortedEndpoints = newArr

	scrape := e.toScrape()
	e.Prometheus.UpdateScrape(scrape)
}

func (e *App) Error(err error) {
	logrus.Error(err)
}

//Name
func (e *App) Name() string {
	return "app"
}

func (e *App) toScrape() *prometheus.ScrapeConfig {
	ts := make([]string, 0, len(e.sortedEndpoints))
	for _, end := range e.sortedEndpoints {
		ts = append(ts, end)
	}

	return &prometheus.ScrapeConfig{
		JobName:        e.Name(),
		ScrapeInterval: model.Duration(5 * time.Second),
		ScrapeTimeout:  model.Duration(4 * time.Second),
		MetricsPath:    "/app/metrics",
		ServiceDiscoveryConfig: prometheus.ServiceDiscoveryConfig{
			StaticConfigs: []*prometheus.Group{
				{
					Targets: ts,
					Labels: map[model.LabelName]model.LabelValue{
						"component": model.LabelValue(e.Name()),
					},
				},
			},
		},
	}
}

//AddEndpoint
func (e *App) AddEndpoint(end *config.Endpoint) {
	e.endpoints = append(e.endpoints, end)
	e.UpdateEndpoints(e.endpoints...)
}

//Add
func (e *App) Add(event *watch.Event) {
	url := gjson.Get(event.GetValueString(), "internal_ip").String() + ":6100"
	end := &config.Endpoint{
		URL:  url,
		Name: event.GetKey(),
	}
	e.AddEndpoint(end)
}

//Modify
func (e *App) Modify(event *watch.Event) {
	var update bool
	url := gjson.Get(event.GetValueString(), "internal_ip").String() + ":6100"
	for i, end := range e.endpoints {
		if end.Name == event.GetKey() {
			e.endpoints[i].URL = url
			e.UpdateEndpoints(e.endpoints...)
			update = true
			break
		}
	}
	if !update {
		e.endpoints = append(e.endpoints, &config.Endpoint{
			Name: event.GetKey(),
			URL:  url,
		})
		e.UpdateEndpoints(e.endpoints...)
	}
}

//Delete
func (e *App) Delete(event *watch.Event) {
	for i, end := range e.endpoints {
		if end.Name == event.GetKey() {
			e.endpoints = append(e.endpoints[:i], e.endpoints[i+1:]...)
			e.UpdateEndpoints(e.endpoints...)
			break
		}
	}
}
