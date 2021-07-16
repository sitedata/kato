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

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/kubernetes"

	"github.com/gridworkz/kato/cmd/gateway/option"
	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/gateway/cluster"
	"github.com/gridworkz/kato/gateway/controller"
	"github.com/gridworkz/kato/gateway/metric"
	"github.com/gridworkz/kato/util"

	etcdutil "github.com/gridworkz/kato/util/etcd"
	k8sutil "github.com/gridworkz/kato/util/k8s"
)

//Run
func Run(s *option.GWServer) error {
	logrus.Info("start gateway...")
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := k8sutil.NewRestConfig(s.K8SConfPath)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints:   s.Config.EtcdEndpoint,
		CaFile:      s.Config.EtcdCaFile,
		CertFile:    s.Config.EtcdCertFile,
		KeyFile:     s.Config.EtcdKeyFile,
		DialTimeout: time.Duration(s.Config.EtcdTimeout) * time.Second,
	}
	etcdCli, err := etcdutil.NewClient(ctx, etcdClientArgs)
	if err != nil {
		return err
	}

	//create cluster node manager
	logrus.Debug("start creating node manager")
	node, err := cluster.CreateNodeManager(ctx, s.Config, etcdCli)
	if err != nil {
		return fmt.Errorf("create gateway node manager failure %s", err.Error())
	}
	if err := node.Start(); err != nil {
		return fmt.Errorf("start node manager: %v", err)
	}
	defer node.Stop()

	reg := prometheus.NewRegistry()
	reg.MustRegister(prometheus.NewGoCollector())
	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	mc := metric.NewDummyCollector()
	if s.Config.EnableMetrics {
		mc, err = metric.NewCollector(s.NodeName, reg)
		if err != nil {
			logrus.Fatalf("Error creating prometheus collector:  %v", err)
		}
	}
	mc.Start()

	gwc, err := controller.NewGWController(ctx, clientset, &s.Config, mc, node)
	if err != nil {
		return err
	}
	if gwc == nil {
		return fmt.Errorf("Failed to create new GWController")
	}
	logrus.Debug("start gateway controller")
	if err := gwc.Start(errCh); err != nil {
		return fmt.Errorf("Failed to start GWController %s", err.Error())
	}
	defer gwc.Close()

	mux := chi.NewMux()
	registerHealthz(gwc, mux)
	registerMetrics(reg, mux)
	if s.Debug {
		util.ProfilerSetup(mux)
	}
	go startHTTPServer(s.ListenPorts.Health, mux)

	keepalive, err := discover.CreateKeepAlive(etcdClientArgs, "gateway", s.Config.NodeName,
		s.Config.HostIP, s.ListenPorts.Health)
	if err != nil {
		return err
	}
	logrus.Debug("start keepalive")
	if err := keepalive.Start(); err != nil {
		return err
	}
	defer keepalive.Stop()

	logrus.Info("RBD app gateway start success!")

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-term:
		logrus.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-errCh:
		logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
	}
	logrus.Info("See you next time!")

	return nil
}

func registerHealthz(gc *controller.GWController, mux *chi.Mux) {
	// expose health check endpoint (/healthz)
	healthz.InstallHandler(mux,
		healthz.PingHealthz,
		gc,
	)
}

func registerMetrics(reg *prometheus.Registry, mux *chi.Mux) {
	mux.Handle(
		"/metrics",
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	)
}

func startHTTPServer(port int, mux *chi.Mux) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	logrus.Fatal(server.ListenAndServe())
}
