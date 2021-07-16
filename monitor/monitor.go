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
	"time"

	v3 "github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/cmd/monitor/option"
	discoverv1 "github.com/gridworkz/kato/discover"
	discoverv2 "github.com/gridworkz/kato/discover.v2"
	"github.com/gridworkz/kato/discover/config"
	"github.com/gridworkz/kato/monitor/callback"
	"github.com/gridworkz/kato/monitor/prometheus"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	k8sutil "github.com/gridworkz/kato/util/k8s"
	"github.com/gridworkz/kato/util/watch"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

//Monitor
type Monitor struct {
	config         *option.Config
	ctx            context.Context
	cancel         context.CancelFunc
	client         *v3.Client
	timeout        time.Duration
	manager        *prometheus.Manager
	discoverv1     discoverv1.Discover
	discoverv2     discoverv2.Discover
	serviceMonitor *prometheus.ServiceMonitorController
	stopCh         chan struct{}
}

//Start
func (d *Monitor) Start() {
	d.discoverv1.AddProject("prometheus", &callback.Prometheus{Prometheus: d.manager})
	d.discoverv1.AddProject("event_log_event_http", &callback.EventLog{Prometheus: d.manager})
	d.discoverv1.AddProject("acp_webcli", &callback.Webcli{Prometheus: d.manager})
	d.discoverv1.AddProject("gateway", &callback.GatewayNode{Prometheus: d.manager})
	d.discoverv2.AddProject("builder", &callback.Builder{Prometheus: d.manager})
	d.discoverv2.AddProject("mq", &callback.Mq{Prometheus: d.manager})
	d.discoverv2.AddProject("app_sync_runtime_server", &callback.Worker{Prometheus: d.manager})

	// node and app runtime metrics needs to be monitored separately
	go d.discoverNodes(&callback.Node{Prometheus: d.manager}, &callback.App{Prometheus: d.manager}, d.ctx.Done())

	// monitor etcd members
	go d.discoverEtcd(&callback.Etcd{
		Prometheus: d.manager,
		Scheme: func() string {
			if d.config.EtcdCertFile != "" {
				return "https"
			}
			return "http"
		}(),
		TLSConfig: prometheus.TLSConfig{
			CAFile:   d.config.EtcdCaFile,
			CertFile: d.config.EtcdCertFile,
			KeyFile:  d.config.EtcdKeyFile,
		},
	}, d.ctx.Done())

	// monitor Cadvisor
	go d.discoverCadvisor(&callback.Cadvisor{
		Prometheus: d.manager,
		ListenPort: d.config.CadvisorListenPort,
	}, d.ctx.Done())

	// kubernetes service discovery
	rbdapi := callback.RbdAPI{Prometheus: d.manager}
	rbdapi.UpdateEndpoints(nil)

	// service monitor
	d.serviceMonitor.Run(d.stopCh)
}

func (d *Monitor) discoverNodes(node *callback.Node, app *callback.App, done <-chan struct{}) {
	// start listen node modified
	watcher := watch.New(d.client, "")
	w, err := watcher.WatchList(d.ctx, "/kato/nodes", "")
	if err != nil {
		logrus.Error("failed to watch list for discover all nodes: ", err)
		return
	}
	defer w.Stop()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				logrus.Warn("the events channel is closed.")
				return
			}

			switch event.Type {
			case watch.Added:
				node.Add(&event)

				isSlave := gjson.Get(event.GetValueString(), "labels.kato_node_rule_compute").String()
				if isSlave == "true" {
					app.Add(&event)
				}
			case watch.Modified:
				node.Modify(&event)

				isSlave := gjson.Get(event.GetValueString(), "labels.kato_node_rule_compute").String()
				if isSlave == "true" {
					app.Modify(&event)
				}
			case watch.Deleted:
				node.Delete(&event)

				isSlave := gjson.Get(event.GetPreValueString(), "labels.kato_node_rule_compute").String()
				if isSlave == "true" {
					app.Delete(&event)
				}
			case watch.Error:
				logrus.Error("error when read a event from result chan for discover all nodes: ", event.Error)
			}
		case <-done:
			logrus.Info("stop discover nodes because received stop signal.")
			return
		}

	}

}

func (d *Monitor) discoverCadvisor(c *callback.Cadvisor, done <-chan struct{}) {
	// start listen node modified
	watcher := watch.New(d.client, "")
	w, err := watcher.WatchList(d.ctx, "/kato/nodes", "")
	if err != nil {
		logrus.Error("failed to watch list for discover all nodes: ", err)
		return
	}
	defer w.Stop()

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				logrus.Warn("the events channel is closed.")
				return
			}
			switch event.Type {
			case watch.Added:
				isSlave := gjson.Get(event.GetValueString(), "labels.kato_node_rule_compute").String()
				if isSlave == "true" {
					c.Add(&event)
				}
			case watch.Modified:
				isSlave := gjson.Get(event.GetValueString(), "labels.kato_node_rule_compute").String()
				if isSlave == "true" {
					c.Modify(&event)
				}
			case watch.Deleted:
				isSlave := gjson.Get(event.GetPreValueString(), "labels.kato_node_rule_compute").String()
				if isSlave == "true" {
					c.Delete(&event)
				}
			case watch.Error:
				logrus.Error("error when read a event from result chan for discover all nodes: ", event.Error)
			}
		case <-done:
			logrus.Info("stop discover nodes because received stop signal.")
			return
		}

	}

}

func (d *Monitor) discoverEtcd(e *callback.Etcd, done <-chan struct{}) {
	t := time.Tick(time.Minute)
	for {
		select {
		case <-done:
			logrus.Info("stop discover etcd because received stop signal.")
			return
		case <-t:
			resp, err := d.client.MemberList(d.ctx)
			if err != nil {
				logrus.Error("Failed to list etcd members for discover etcd.")
				continue
			}

			endpoints := make([]*config.Endpoint, 0, 5)
			for _, member := range resp.Members {
				if len(member.ClientURLs) >= 1 {
					url := member.ClientURLs[0]
					end := &config.Endpoint{
						URL: url,
					}
					endpoints = append(endpoints, end)
				}
			}
			logrus.Debugf("etcd endpoints: %+v", endpoints)
			e.UpdateEndpoints(endpoints...)
		}
	}
}

// Stop monitor
func (d *Monitor) Stop() {
	logrus.Info("Stopping all child process for monitor")
	d.cancel()
	d.discoverv1.Stop()
	d.discoverv2.Stop()
	d.client.Close()
	close(d.stopCh)
}

// NewMonitor
func NewMonitor(opt *option.Config, p *prometheus.Manager) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	defaultTimeout := time.Second * 3

	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints:   opt.EtcdEndpoints,
		DialTimeout: defaultTimeout,
		CaFile:      opt.EtcdCaFile,
		CertFile:    opt.EtcdCertFile,
		KeyFile:     opt.EtcdKeyFile,
	}

	cli, err := etcdutil.NewClient(ctx, etcdClientArgs)
	v3.New(v3.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	dc1, err := discoverv1.GetDiscover(config.DiscoverConfig{EtcdClientArgs: etcdClientArgs})
	if err != nil {
		logrus.Fatal(err)
	}

	dc3, err := discoverv2.GetDiscover(config.DiscoverConfig{EtcdClientArgs: etcdClientArgs})
	if err != nil {
		logrus.Fatal(err)
	}
	restConfig, err := k8sutil.NewRestConfig(opt.KubeConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	d := &Monitor{
		config:     opt,
		ctx:        ctx,
		cancel:     cancel,
		manager:    p,
		client:     cli,
		discoverv1: dc1,
		discoverv2: dc3,
		timeout:    defaultTimeout,
		stopCh:     make(chan struct{}),
	}
	d.serviceMonitor, err = prometheus.NewServiceMonitorController(ctx, restConfig, p)
	if err != nil {
		logrus.Fatal(err)
	}
	return d
}
