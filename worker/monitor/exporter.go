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
	"context"
	"net/http"

	"github.com/gridworkz/kato/worker/master"

	"github.com/gridworkz/kato/cmd/worker/option"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/gridworkz/kato/worker/appm/controller"
	"github.com/gridworkz/kato/worker/discover"
	"github.com/gridworkz/kato/worker/monitor/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
)

//ExporterManager app resource exporter
type ExporterManager struct {
	ctx               context.Context
	cancel            context.CancelFunc
	config            option.Config
	stopChan          chan struct{}
	masterController  *master.Controller
	controllermanager *controller.Manager
}

//NewManager return *NewManager
func NewManager(c option.Config, masterController *master.Controller, controllermanager *controller.Manager) *ExporterManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ExporterManager{
		ctx:               ctx,
		cancel:            cancel,
		config:            c,
		stopChan:          make(chan struct{}),
		masterController:  masterController,
		controllermanager: controllermanager,
	}
}
func (t *ExporterManager) handler(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.New(t.masterController, t.controllermanager))

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

//Start
func (t *ExporterManager) Start() error {
	http.HandleFunc(t.config.PrometheusMetricPath, t.handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Worker exporter</title></head>
			<body>
			<h1>Worker exporter</h1>
			<p><a href='` + t.config.PrometheusMetricPath + `'>Metrics</a></p>
			</body>
			</html>
			`))
	})
	http.HandleFunc("/worker/health", func(w http.ResponseWriter, r *http.Request) {
		healthStatus := discover.HealthCheck()
		if healthStatus["status"] != "health" {
			httputil.ReturnError(r, w, 400, "worker service unusual")
		}
		httputil.ReturnSuccess(r, w, healthStatus)
	})
	log.Infoln("Listening on", t.config.Listen)
	go func() {
		log.Fatal(http.ListenAndServe(t.config.Listen, nil))
	}()
	logrus.Info("start app resource exporter success.")
	return nil
}

//Stop
func (t *ExporterManager) Stop() {
	t.cancel()
}
