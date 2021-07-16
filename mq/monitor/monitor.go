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

package monitor

import (
	"github.com/gridworkz/kato/mq/api/mq"
	"github.com/prometheus/client_golang/prometheus"
)

// Metric name parts.
const (
	// Namespace for all metrics.
	namespace = "acp_mq"
	// Subsystem(s).
	exporter = "exporter"
)

// Metric descriptors.
var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
)

//Exporter collects entrance metrics. It implements prometheus.Collector.
type Exporter struct {
	error              prometheus.Gauge
	totalScrapes       prometheus.Counter
	scrapeErrors       *prometheus.CounterVec
	lbPluginUp         prometheus.Gauge
	queueMessageNumber *prometheus.GaugeVec
	mqm                mq.ActionMQ
}

var healthDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, exporter, "health_status"),
	"health status.",
	[]string{"service_name"}, nil,
)

//NewExporter
func NewExporter(mqm mq.ActionMQ) *Exporter {
	return &Exporter{
		mqm: mqm,
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "scrapes_total",
			Help:      "Total number of times Entrance was scraped for metrics.",
		}),
		scrapeErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "scrape_errors_total",
			Help:      "Total number of times an error occurred scraping a Entrance.",
		}, []string{"collector"}),
		error: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: exporter,
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from Entrance resulted in an error (1 for error, 0 for success).",
		}),
		lbPluginUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Whether the default lb plugin is up.",
		}),
		queueMessageNumber: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "queue_message_number",
			Help:      "Message queue enqueue total.",
		}, []string{"topic"}),
	}
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
	e.scrape(ch)
	ch <- e.totalScrapes
	ch <- e.error
	e.scrapeErrors.Collect(ch)
	for _, topic := range e.mqm.GetAllTopics() {
		e.queueMessageNumber.WithLabelValues(topic).Set(float64(e.mqm.MessageQueueSize(topic)))
	}
	e.queueMessageNumber.Collect(ch)
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	e.totalScrapes.Inc()
	ch <- prometheus.MustNewConstMetric(healthDesc, prometheus.GaugeValue, 1, "mq")
}
