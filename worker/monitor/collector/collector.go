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

package collector

import (
	"github.com/gridworkz/kato/worker/master"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/worker/appm/controller"
	"github.com/gridworkz/kato/worker/discover"
	"github.com/prometheus/client_golang/prometheus"
)

//Exporter collector
type Exporter struct {
	error                     prometheus.Gauge
	totalScrapes              prometheus.Counter
	scrapeErrors              *prometheus.CounterVec
	workerUp                  prometheus.Gauge
	dbmanager                 db.Manager
	masterController          *master.Controller
	controllermanager         *controller.Manager
	taskNum                   prometheus.Counter
	taskUpNum                 prometheus.Gauge
	taskError prometheus.Counter
	storeComponentNum         prometheus.Gauge
	thirdComponentDiscoverNum prometheus.Gauge
}

var scrapeDurationDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "exporter", "collector_duration_seconds"),
	"Collector time duration.",
	[]string{"collector"}, nil,
)

var healthDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "exporter", "health_status"),
	"health status.",
	[]string{"service_name"}, nil,
)

//Describe Describe
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
	e.scrape (ch)
	ch <- e.totalScrapes
	ch <- e.error
	e.scrapeErrors.Collect(ch)
	ch <- e.workerUp
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	e.totalScrapes.Inc()
	e.masterController.Scrape(ch, scrapeDurationDesc)
	healthInfo := discover.HealthCheck()
	healthStatus := healthInfo["status"]
	var val float64
	if healthStatus == "health" {
		val = 1
	} else {
		val = 0
	}
	ch <- prometheus.MustNewConstMetric(healthDesc, prometheus.GaugeValue, val, "worker")
	ch <- prometheus.MustNewConstMetric(e.taskUpNum.Desc(),
		prometheus.GaugeValue,
		float64(e.controllermanager.GetControllerSize()))
	ch <- prometheus.MustNewConstMetric(e.taskNum.Desc(), prometheus.CounterValue, discover.TaskNum)
	ch <- prometheus.MustNewConstMetric(e.taskError.Desc(), prometheus.CounterValue, discover.TaskError)
	ch <- prometheus.MustNewConstMetric(e.storeComponentNum.Desc(), prometheus.GaugeValue, float64(len(e.masterController.GetStore().GetAllAppServices())))
}

var namespace = "worker"

//New Create a collector
func New(masterController *master.Controller, controllermanager *controller.Manager) *Exporter {
	return &Exporter{
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "scrapes_total",
			Help:      "Total number of times Worker was scraped for metrics.",
		}),
		scrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "scrape_errors_total",
			Help:      "Total number of times an error occurred scraping a Worker.",
		}, []string{"collector"}),
		error: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from Worker resulted in an error (1 for error, 0 for success).",
		}),
		workerUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Whether the Worker server is up.",
		}),
		taskUpNum: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "task_up_number",
			Help:      "Number of tasks being performed",
		}),
		taskNum: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "worker_task_number",
			Help:      "worker total number of tasks.",
		}),
		taskError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "exporter",
			Name:      "worker_task_error",
			Help:      "worker number of task errors.",
		}),
		storeComponentNum: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "store_component_number",
			Help:      "Number of components in the store cache.",
		}),
		dbmanager:         db.GetManager(),
		masterController:  masterController,
		controllermanager: controllermanager,
	}
}
