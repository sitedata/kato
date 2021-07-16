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
	"os"
	"os/signal"
	"syscall"

	"github.com/gridworkz/kato/builder/discover"
	"github.com/gridworkz/kato/builder/exector"
	"github.com/gridworkz/kato/builder/monitor"
	"github.com/gridworkz/kato/cmd/builder/option"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/mq/client"

	"net/http"

	"github.com/gridworkz/kato/builder/api"
	"github.com/gridworkz/kato/builder/clean"
	discoverv2 "github.com/gridworkz/kato/discover.v2"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

//Run
func Run(s *option.Builder) error {
	errChan := make(chan error)
	//init mysql
	dbconfig := config.Config{
		DBType:              s.Config.DBType,
		MysqlConnectionInfo: s.Config.MysqlConnectionInfo,
		EtcdEndPoints:       s.Config.EtcdEndPoints,
		EtcdTimeout:         s.Config.EtcdTimeout,
	}
	if err := db.CreateManager(dbconfig); err != nil {
		return err
	}
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: s.Config.EtcdEndPoints,
		CaFile:    s.Config.EtcdCaFile,
		CertFile:  s.Config.EtcdCertFile,
		KeyFile:   s.Config.EtcdKeyFile,
	}
	if err := event.NewManager(event.EventConfig{
		EventLogServers: s.Config.EventLogServers,
		DiscoverArgs:    etcdClientArgs,
	}); err != nil {
		return err
	}
	defer event.CloseManager()
	client, err := client.NewMqClient(etcdClientArgs, s.Config.MQAPI)
	if err != nil {
		logrus.Errorf("new Mq client error, %v", err)
		return err
	}
	exec, err := exector.NewManager(s.Config, client)
	if err != nil {
		return err
	}
	if err := exec.Start(); err != nil {
		return err
	}
	//exec manage stop by discover
	dis := discover.NewTaskManager(s.Config, client, exec)
	if err := dis.Start(errChan); err != nil {
		return err
	}
	defer dis.Stop()

	if s.Config.CleanUp {
		cle, err := clean.CreateCleanManager()
		if err != nil {
			return err
		}
		if err := cle.Start(errChan); err != nil {
			return err
		}
		defer cle.Stop()
	}
	keepalive, err := discoverv2.CreateKeepAlive(etcdClientArgs, "builder",
		"", s.Config.HostIP, s.Config.APIPort)
	if err != nil {
		return err
	}
	if err := keepalive.Start(); err != nil {
		return err
	}
	defer keepalive.Stop()

	exporter := monitor.NewExporter(exec)
	prometheus.MustRegister(exporter)
	r := api.APIServer()
	r.Handle(s.Config.PrometheusMetricPath, promhttp.Handler())
	logrus.Info("builder api listen port 3228")
	go http.ListenAndServe(":3228", r)

	logrus.Info("builder begin running...")
	//final step: listen Signal
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		logrus.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-errChan:
		logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
	}
	logrus.Info("See you next time!")
	return nil
}
