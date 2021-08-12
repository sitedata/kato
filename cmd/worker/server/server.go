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

	"github.com/eapache/channels"
	"github.com/gridworkz/kato/cmd/worker/option"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/pkg/common"
	"github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	k8sutil "github.com/gridworkz/kato/util/k8s"
	"github.com/gridworkz/kato/worker/appm/componentdefinition"
	"github.com/gridworkz/kato/worker/appm/controller"
	"github.com/gridworkz/kato/worker/appm/store"
	"github.com/gridworkz/kato/worker/discover"
	"github.com/gridworkz/kato/worker/gc"
	"github.com/gridworkz/kato/worker/master"
	"github.com/gridworkz/kato/worker/monitor"
	"github.com/gridworkz/kato/worker/server"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/flowcontrol"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//Run start run
func Run(s *option.Worker) error {
	errChan := make(chan error, 2)
	dbconfig := config.Config{
		DBType:              s.Config.DBType,
		MysqlConnectionInfo: s.Config.MysqlConnectionInfo,
		EtcdEndPoints:       s.Config.EtcdEndPoints,
		EtcdTimeout:         s.Config.EtcdTimeout,
	}
	//step 1:db manager init ,event log client init
	if err := db.CreateManager(dbconfig); err != nil {
		return err
	}
	defer db.CloseManager()
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

	//step 2 : create kube client and etcd client
	restConfig, err := k8sutil.NewRestConfig(s.Config.KubeConfig)
	if err != nil {
		logrus.Errorf("create kube rest config error: %s", err.Error())
		return err
	}
	restConfig.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(float32(s.Config.KubeAPIQPS), s.Config.KubeAPIBurst)
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		logrus.Errorf("create kube client error: %s", err.Error())
		return err
	}
	s.Config.KubeClient = clientset
	runtimeClient, err := client.New(restConfig, client.Options{Scheme: common.Scheme})
	if err != nil {
		logrus.Errorf("create kube runtime client error: %s", err.Error())
		return err
	}
	katoClient := versioned.NewForConfigOrDie(restConfig)
	//step 3: create componentdefinition builder factory
	componentdefinition.NewComponentDefinitionBuilder(s.Config.RBDNamespace)

	//step 4: create component resource store
	updateCh := channels.NewRingChannel(1024)
	cachestore := store.NewStore(restConfig, clientset, katoClient, db.GetManager(), s.Config)
	if err := cachestore.Start(); err != nil {
		logrus.Error("start kube cache store error", err)
		return err
	}

	//step 5: create controller manager
	controllerManager := controller.NewManager(cachestore, clientset, runtimeClient)
	defer controllerManager.Stop()

	//step 6 : start runtime master
	masterCon, err := master.NewMasterController(s.Config, cachestore, clientset, katoClient, restConfig)
	if err != nil {
		return err
	}
	if err := masterCon.Start(); err != nil {
		return err
	}
	defer masterCon.Stop()

	//step 7 : create discover module
	garbageCollector := gc.NewGarbageCollector(clientset)
	taskManager := discover.NewTaskManager(s.Config, cachestore, controllerManager, garbageCollector)
	if err := taskManager.Start(); err != nil {
		return err
	}
	defer taskManager.Stop()

	//step 8: start app runtimer server
	runtimeServer := server.CreaterRuntimeServer(s.Config, cachestore, clientset, updateCh)
	runtimeServer.Start(errChan)

	//step 9: create application use resource exporter.
	exporterManager := monitor.NewManager(s.Config, masterCon, controllerManager)
	if err := exporterManager.Start(); err != nil {
		return err
	}
	defer exporterManager.Stop()

	logrus.Info("worker begin running...")

	//step finally: listen Signal
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		logrus.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-errChan:
		if err != nil {
			logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
		}
	}
	logrus.Info("See you next time!")
	return nil
}
