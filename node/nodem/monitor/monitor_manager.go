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

package monitor

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/api"
	"github.com/gridworkz/kato/node/monitormessage"
	"github.com/gridworkz/kato/node/statsd"
	innerprometheus "github.com/gridworkz/kato/node/statsd/prometheus"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/node_exporter/collector"
	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Manager
type Manager interface {
	Start(errchan chan error) error
	Stop() error
	SetAPIRoute(apim *api.Manager) error
}

type manager struct {
	statsdExporter     *statsd.Exporter
	statsdRegistry     *innerprometheus.Registry
	nodeExporterRestry *prometheus.Registry
	meserver           *monitormessage.UDPServer
}

func createNodeExporterRestry() (*prometheus.Registry, error) {
	registry := prometheus.NewRegistry()
	filters := []string{"cpu", "diskstats", "filesystem",
		"ipvs", "loadavg", "meminfo", "netdev",
		"netclass", "netdev", "netstat",
		"uname", "mountstats", "nfs"}
	// init kingpin parse
	kingpin.CommandLine.Parse([]string{"--collector.mountstats=true"})
	nc, err := collector.NewNodeCollector(log.NewNopLogger(), filters...)
	if err != nil {
		return nil, err
	}
	for n := range nc.Collectors {
		logrus.Infof("node collector - %s", n)
	}
	err = registry.Register(nc)
	if err != nil {
		return nil, err
	}
	return registry, nil
}

//CreateManager
func CreateManager(ctx context.Context, c *option.Conf) (Manager, error) {
	//statsd exporter
	statsdRegistry := innerprometheus.NewRegistry()
	exporter := statsd.CreateExporter(c.StatsdConfig, statsdRegistry)
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: c.EtcdEndpoints,
		CaFile:    c.EtcdCaFile,
		CertFile:  c.EtcdCertFile,
		KeyFile:   c.EtcdKeyFile,
	}
	meserver := monitormessage.CreateUDPServer(ctx, "0.0.0.0", 6666, etcdClientArgs)
	nodeExporterRestry, err := createNodeExporterRestry()
	if err != nil {
		return nil, err
	}
	manage := &manager{
		statsdExporter:     exporter,
		statsdRegistry:     statsdRegistry,
		nodeExporterRestry: nodeExporterRestry,
		meserver:           meserver,
	}
	return manage, nil
}

func (m *manager) Start(errchan chan error) error {
	if err := m.statsdExporter.Start(); err != nil {
		logrus.Errorf("start statsd exporter server error,%s", err.Error())
		return err
	}
	if err := m.meserver.Start(); err != nil {
		return err
	}

	return nil
}

func (m *manager) Stop() error {
	return nil
}

//ReloadStatsdMappConfig
func (m *manager) ReloadStatsdMappConfig(w http.ResponseWriter, r *http.Request) {
	if err := m.statsdExporter.ReloadConfig(); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
	} else {
		w.Write([]byte("Success reload"))
		w.WriteHeader(200)
	}
}

//HandleStatsd handle
func (m *manager) HandleStatsd(w http.ResponseWriter, r *http.Request) {
	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		m.statsdRegistry,
	}
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      logrus.StandardLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	h.ServeHTTP(w, r)
}

//NodeExporter
func (m *manager) NodeExporter(w http.ResponseWriter, r *http.Request) {
	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		m.nodeExporterRestry,
	}
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      logrus.StandardLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	h.ServeHTTP(w, r)
}

//SetAPIRoute set api route rule
func (m *manager) SetAPIRoute(apim *api.Manager) error {
	apim.GetRouter().Get("/app/metrics", m.HandleStatsd)
	apim.GetRouter().Get("/-/statsdreload", m.ReloadStatsdMappConfig)
	apim.GetRouter().Get("/node/metrics", m.NodeExporter)
	return nil
}
