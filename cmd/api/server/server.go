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
	"os"
	"os/signal"
	"syscall"

	"github.com/gridworkz/kato/api/controller"
	"github.com/gridworkz/kato/api/db"
	"github.com/gridworkz/kato/api/discover"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/server"
	"github.com/gridworkz/kato/cmd/api/option"
	"github.com/gridworkz/kato/event"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	k8sutil "github.com/gridworkz/kato/util/k8s"
	"github.com/gridworkz/kato/worker/client"
	"k8s.io/client-go/kubernetes"

	"github.com/sirupsen/logrus"
)

//Run
func Run(s *option.APIServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error)
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: s.Config.EtcdEndpoint,
		CaFile:    s.Config.EtcdCaFile,
		CertFile:  s.Config.EtcdCertFile,
		KeyFile:   s.Config.EtcdKeyFile,
	}
	//Start service discovery
	if _, err := discover.CreateEndpointDiscover(etcdClientArgs); err != nil {
		return err
	}
	//Create db manager
	if err := db.CreateDBManager(s.Config); err != nil {
		logrus.Debugf("create db manager error, %v", err)
		return err
	}
	//Create event manager
	if err := db.CreateEventManager(s.Config); err != nil {
		logrus.Debugf("create event manager error, %v", err)
	}
	config, err := k8sutil.NewRestConfig(s.KubeConfigPath)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	if err := event.NewManager(event.EventConfig{
		EventLogServers: s.Config.EventLogServers,
		DiscoverArgs:    etcdClientArgs,
	}); err != nil {
		return err
	}
	defer event.CloseManager()
	//create app status client
	cli, err := client.NewClient(ctx, client.AppRuntimeSyncClientConf{
		EtcdEndpoints: s.Config.EtcdEndpoint,
		EtcdCaFile:    s.Config.EtcdCaFile,
		EtcdCertFile:  s.Config.EtcdCertFile,
		EtcdKeyFile:   s.Config.EtcdKeyFile,
	})
	if err != nil {
		logrus.Errorf("create app status client error, %v", err)
		return err
	}

	etcdcli, err := etcdutil.NewClient(ctx, etcdClientArgs)
	if err != nil {
		logrus.Errorf("create etcd client v3 error, %v", err)
		return err
	}

	//middleware initialization
	handler.InitProxy(s.Config)
	//Create handle
	if err := handler.InitHandle(s.Config, etcdClientArgs, cli, etcdcli, clientset); err != nil {
		logrus.Errorf("init all handle error, %v", err)
		return err
	}
	//Create v2Router manager
	if err := controller.CreateV2RouterManager(s.Config, cli); err != nil {
		logrus.Errorf("create v2 route manager error, %v", err)
	}
	// Start api
	apiManager := server.NewManager(s.Config, etcdcli)
	if err := apiManager.Start(); err != nil {
		return err
	}
	defer apiManager.Stop()
	logrus.Info("api router is running...")

	//final step: listen Signal
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-term:
		logrus.Infof("Received a Signal  %s, exiting gracefully...", s.String())
	case err := <-errChan:
		logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
	}
	logrus.Info("See you next time!")
	return nil
}
