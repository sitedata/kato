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
	"github.com/gridworkz/kato/discover.v2"
	"github.com/gridworkz/kato/node/initiate"
	"github.com/gridworkz/kato/util/constants"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/signal"
	"syscall"

	"github.com/gridworkz/kato/cmd/node/option"
	eventLog "github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/node/api"
	"github.com/gridworkz/kato/node/api/controller"
	"github.com/gridworkz/kato/node/core/store"
	"github.com/gridworkz/kato/node/kubecache"
	"github.com/gridworkz/kato/node/masterserver"
	"github.com/gridworkz/kato/node/nodem"
	"github.com/gridworkz/kato/node/nodem/docker"
	"github.com/gridworkz/kato/node/nodem/envoy"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	k8sutil "github.com/gridworkz/kato/util/k8s"

	"github.com/sirupsen/logrus"
)

//Run start
func Run(cfg *option.Conf) error {
	var stoped = make(chan struct{})
	stopfunc := func() error {
		close(stoped)
		return nil
	}
	startfunc := func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		etcdClientArgs := &etcdutil.ClientArgs{
			Endpoints:   cfg.EtcdEndpoints,
			CaFile:      cfg.EtcdCaFile,
			CertFile:    cfg.EtcdCertFile,
			KeyFile:     cfg.EtcdKeyFile,
			DialTimeout: cfg.EtcdDialTimeout,
		}
		if err := cfg.ParseClient(ctx, etcdClientArgs); err != nil {
			return fmt.Errorf("config parse error:%s", err.Error())
		}

		config, err := k8sutil.NewRestConfig(cfg.K8SConfPath)
		if err != nil {
			return err
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}

		k8sDiscover := discover.NewK8sDiscover(ctx, clientset, cfg)
		defer k8sDiscover.Stop()

		nodemanager, err := nodem.NewNodeManager(ctx, cfg)
		if err != nil {
			return fmt.Errorf("create node manager failed: %s", err)
		}
		if err := nodemanager.InitStart(); err != nil {
			return err
		}

		err = eventLog.NewManager(eventLog.EventConfig{
			EventLogServers: cfg.EventLogServer,
			DiscoverArgs:    etcdClientArgs,
		})
		if err != nil {
			logrus.Errorf("error creating eventlog manager")
			return nil
		}
		defer eventLog.CloseManager()
		logrus.Debug("create and start event log client success")

		kubecli, err := kubecache.NewKubeClient(cfg, clientset)
		if err != nil {
			return err
		}
		defer kubecli.Stop()

		if cfg.ImageRepositoryHost == constants.DefImageRepository {
			hostManager, err := initiate.NewHostManager(cfg, k8sDiscover)
			if err != nil {
				return fmt.Errorf("create new host manager: %v", err)
			}
			hostManager.Start()
		}

		logrus.Debugf("rbd-namespace=%s; rbd-docker-secret=%s", os.Getenv("RBD_NAMESPACE"), os.Getenv("RBD_DOCKER_SECRET"))
		// sync docker inscure registries cert info into all kato node
		if err = docker.SyncDockerCertFromSecret(clientset, os.Getenv("RBD_NAMESPACE"), os.Getenv("RBD_DOCKER_SECRET")); err != nil { // TODO gridworkz namespace secretname
			return fmt.Errorf("sync docker cert from secret error: %s", err.Error())
		}

		// init etcd client
		if err = store.NewClient(ctx, cfg, etcdClientArgs); err != nil {
			return fmt.Errorf("Connect to ETCD %s failed: %s", cfg.EtcdEndpoints, err)
		}
		errChan := make(chan error, 3)
		if err := nodemanager.Start(errChan); err != nil {
			return fmt.Errorf("start node manager failed: %s", err)
		}
		defer nodemanager.Stop()
		logrus.Debug("create and start node manager moudle success")

		//The master service starts after the node service
		var ms *masterserver.MasterServer
		if cfg.RunMode == "master" {
			ms, err = masterserver.NewMasterServer(nodemanager.GetCurrentNode(), kubecli)
			if err != nil {
				logrus.Errorf(err.Error())
				return err
			}
			ms.Cluster.UpdateNode(nodemanager.GetCurrentNode())
			if err := ms.Start(errChan); err != nil {
				logrus.Errorf(err.Error())
				return err
			}
			defer ms.Stop(nil)
			logrus.Debug("create and start master server moudle success")
		}
		//create api manager
		apiManager := api.NewManager(*cfg, nodemanager.GetCurrentNode(), ms, kubecli)
		if err := apiManager.Start(errChan); err != nil {
			return err
		}
		if err := nodemanager.AddAPIManager(apiManager); err != nil {
			return err
		}
		defer apiManager.Stop()

		//create service mesh controller
		grpcserver, err := envoy.CreateDiscoverServerManager(clientset, *cfg)
		if err != nil {
			return err
		}
		if err := grpcserver.Start(errChan); err != nil {
			return err
		}
		defer grpcserver.Stop()

		logrus.Debug("create and start api server moudle success")

		defer controller.Exist(nil)
		//step finally: listen Signal
		term := make(chan os.Signal)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		select {
		case <-stoped:
			logrus.Infof("windows service stoped..")
		case <-term:
			logrus.Warn("Received SIGTERM, exiting gracefully...")
		case err := <-errChan:
			logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
		}
		logrus.Info("See you next time!")
		return nil
	}
	err := initService(cfg, startfunc, stopfunc)
	if err != nil {
		return err
	}
	return nil
}
