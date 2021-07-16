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

package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/node/kubecache"
	"github.com/gridworkz/kato/node/masterserver"
	"github.com/gridworkz/kato/node/statsd"

	"github.com/gridworkz/kato/node/api/controller"
	"github.com/gridworkz/kato/node/api/router"

	"context"
	"strings"

	"github.com/gridworkz/kato/cmd/node/option"
	nodeclient "github.com/gridworkz/kato/node/nodem/client"

	_ "net/http/pprof"

	client "github.com/coreos/etcd/clientv3"
	"github.com/go-chi/chi"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/sirupsen/logrus"
)

//Manager
type Manager struct {
	ctx            context.Context
	cancel         context.CancelFunc
	conf           option.Conf
	router         *chi.Mux
	node           *nodeclient.HostNode
	lID            client.LeaseID // lease id
	ms             *masterserver.MasterServer
	keepalive      *discover.KeepAlive
	exporter       *statsd.Exporter
	etcdClientArgs *etcdutil.ClientArgs
}

//NewManager
func NewManager(c option.Conf, node *nodeclient.HostNode, ms *masterserver.MasterServer, kubecli kubecache.KubeClient) *Manager {
	r := router.Routers(c.RunMode)
	ctx, cancel := context.WithCancel(context.Background())
	controller.Init(&c, ms, kubecli)
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints:   c.EtcdEndpoints,
		CaFile:      c.EtcdCaFile,
		CertFile:    c.EtcdCertFile,
		KeyFile:     c.EtcdKeyFile,
		DialTimeout: c.EtcdDialTimeout,
	}
	m := &Manager{
		ctx:            ctx,
		cancel:         cancel,
		conf:           c,
		router:         r,
		node:           node,
		ms:             ms,
		etcdClientArgs: etcdClientArgs,
	}
	// set node cluster monitor route
	m.router.Get("/cluster/metrics", m.HandleClusterScrape)
	return m
}

//Start
func (m *Manager) Start(errChan chan error) error {
	logrus.Infof("api server start listening on %s", m.conf.APIAddr)
	go func() {
		if err := http.ListenAndServe(m.conf.APIAddr, m.router); err != nil {
			logrus.Error("kato node api listen error.", err.Error())
			errChan <- err
		}
	}()
	go func() {
		if err := http.ListenAndServe(":6102", nil); err != nil {
			logrus.Error("kato node debug api listen error.", err.Error())
			errChan <- err
		}
	}()
	if m.conf.RunMode == "master" {
		portinfo := strings.Split(m.conf.APIAddr, ":")
		var port int
		if len(portinfo) != 2 {
			port = 6100
		} else {
			var err error
			port, err = strconv.Atoi(portinfo[1])
			if err != nil {
				return fmt.Errorf("get the api port info error.%s", err.Error())
			}
		}
		keepalive, err := discover.CreateKeepAlive(m.etcdClientArgs, "acp_node", m.conf.PodIP, m.conf.PodIP, port)
		if err != nil {
			return err
		}
		if err := keepalive.Start(); err != nil {
			return err
		}
	}
	return nil
}

//Stop
func (m *Manager) Stop() error {
	logrus.Info("api server is stoping.")
	m.cancel()
	if m.keepalive != nil {
		m.keepalive.Stop()
	}
	return nil
}

//GetRouter
func (m *Manager) GetRouter() *chi.Mux {
	return m.router
}

//HandleClusterScrape prometheus handle
func (m *Manager) HandleClusterScrape(w http.ResponseWriter, r *http.Request) {
	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
	}
	if m.ms != nil {
		gatherers = append(gatherers, m.ms.GetRegistry())
	}
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      logrus.StandardLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})
	h.ServeHTTP(w, r)
}
