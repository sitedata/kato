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

package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gridworkz/kato/cmd"

	"github.com/gridworkz/kato/monitor/custom"

	"github.com/gridworkz/kato/cmd/monitor/option"
	"github.com/gridworkz/kato/monitor"
	"github.com/gridworkz/kato/monitor/api"
	"github.com/gridworkz/kato/monitor/api/controller"
	"github.com/gridworkz/kato/monitor/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		cmd.ShowVersion("monitor")
	}
	c := option.NewConfig()
	c.AddFlag(pflag.CommandLine)
	c.AddPrometheusFlag(pflag.CommandLine)
	pflag.Parse()

	c.CompleteConfig()

	// start prometheus daemon and watching this status all the time, exit monitor process if start failed
	a := prometheus.NewRulesManager(c)
	p := prometheus.NewManager(c, a)
	controllerManager := controller.NewControllerManager(a, p)

	monitorMysql(c, p)
	monitorKSM(c, p)

	errChan := make(chan error, 1)
	defer close(errChan)
	p.StartDaemon(errChan)
	defer p.StopDaemon()

	// register prometheus address to etcd cluster
	p.Registry.Start()
	defer p.Registry.Stop()

	// start watching components from etcd, and update modify to prometheus config
	m := monitor.NewMonitor(c, p)
	m.Start()
	defer m.Stop()

	r := api.Server(controllerManager)
	logrus.Info("monitor api listen port 3329")
	go http.ListenAndServe(":3329", r)

	//step finally: listen Signal
	term := make(chan os.Signal)
	defer close(term)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		logrus.Warn("Received SIGTERM, exiting monitor gracefully...")
	case err := <-errChan:
		if err != nil {
			logrus.Errorf("Received a error %s from prometheus, exiting monitor gracefully...", err.Error())
		}
	}
	logrus.Info("See you next time!")
}

func monitorMysql(c *option.Config, p *prometheus.Manager) {
	if strings.TrimSpace(c.MysqldExporter) != "" {
		metrics := strings.TrimSpace(c.MysqldExporter)
		logrus.Infof("add mysql metrics[%s] into prometheus", metrics)
		custom.AddMetrics(p, custom.Metrics{Name: "mysql", Path: "/metrics", Metrics: []string{metrics}, Interval: 30 * time.Second, Timeout: 15 * time.Second})
	}
}

func monitorKSM(c *option.Config, p *prometheus.Manager) {
	if strings.TrimSpace(c.KSMExporter) != "" {
		metrics := strings.TrimSpace(c.KSMExporter)
		logrus.Infof("add kube-state-metrics[%s] into prometheus", metrics)
		custom.AddMetrics(p, custom.Metrics{
			Name: "kubernetes",
			Path: "/metrics",
			Scheme: func() string {
				if strings.HasSuffix(metrics, "443") {
					return "https"
				}
				return "http"
			}(),
			Metrics: []string{metrics}, Interval: 30 * time.Second, Timeout: 10 * time.Second},
		)
	}
}
