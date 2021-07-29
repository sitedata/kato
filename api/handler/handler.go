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

package handler

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/api/client/prometheus"
	api_db "github.com/gridworkz/kato/api/db"
	"github.com/gridworkz/kato/api/handler/group"
	"github.com/gridworkz/kato/api/handler/share"
	"github.com/gridworkz/kato/cmd/api/option"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/gridworkz/kato/worker/client"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

//InitHandle Initialize the handle
func InitHandle(conf option.Config,
	etcdClientArgs *etcdutil.ClientArgs,
	statusCli *client.AppRuntimeSyncClient,
	etcdcli *clientv3.Client,
	kubeClient *kubernetes.Clientset,
	katoClient versioned.Interface,
	k8sClient k8sclient.Client,
) error {
	mq := api_db.MQManager{
		EtcdClientArgs: etcdClientArgs,
		DefaultServer:  conf.MQAPI,
	}
	mqClient, errMQ := mq.NewMQManager()
	if errMQ != nil {
		logrus.Errorf("new MQ manager failed, %v", errMQ)
		return errMQ
	}
	prometheusCli, err := prometheus.NewPrometheus(&prometheus.Options{
		Endpoint: conf.PrometheusEndpoint,
	})
	if err != nil {
		logrus.Errorf("new prometheus client failure, %v", err)
		return err
	}
	dbmanager := db.GetManager()
	defaultServieHandler = CreateManager(conf, mqClient, etcdcli, statusCli, prometheusCli, katoClient)
	defaultPluginHandler = CreatePluginManager(mqClient)
	defaultAppHandler = CreateAppManager(mqClient)
	defaultTenantHandler = CreateTenManager(mqClient, statusCli, &conf, kubeClient, prometheusCli, k8sClient)
	defaultNetRulesHandler = CreateNetRulesManager(etcdcli)
	defaultCloudHandler = CreateCloudManager(conf)
	defaultAPPBackupHandler = group.CreateBackupHandle(mqClient, statusCli, etcdcli)
	defaultEventHandler = CreateLogManager(etcdcli)
	shareHandler = &share.ServiceShareHandle{MQClient: mqClient, EtcdCli: etcdcli}
	pluginShareHandler = &share.PluginShareHandle{MQClient: mqClient, EtcdCli: etcdcli}
	if err := CreateTokenIdenHandler(conf); err != nil {
		logrus.Errorf("create token identification mannager error, %v", err)
		return err
	}
	defaultGatewayHandler = CreateGatewayManager(dbmanager, mqClient, etcdcli)
	def3rdPartySvcHandler = Create3rdPartySvcHandler(dbmanager, statusCli)
	operationHandler = CreateOperationHandler(mqClient)
	batchOperationHandler = CreateBatchOperationHandler(mqClient, statusCli, operationHandler)
	defaultAppRestoreHandler = NewAppRestoreHandler()
	defPodHandler = NewPodHandler(statusCli)
	defClusterHandler = NewClusterHandler(kubeClient, conf.RbdNamespace)
	defaultVolumeTypeHandler = CreateVolumeTypeManger(statusCli)
	defaultEtcdHandler = NewEtcdHandler(etcdcli)
	defaultmonitorHandler = NewMonitorHandler(prometheusCli)
	defServiceEventHandler = NewServiceEventHandler()
	defApplicationHandler = NewApplicationHandler(statusCli, prometheusCli, katoClient, kubeClient)
	return nil
}

var defaultServieHandler ServiceHandler
var shareHandler *share.ServiceShareHandle
var pluginShareHandler *share.PluginShareHandle
var defaultmonitorHandler MonitorHandler

//GetMonitorHandle get monitor handler
func GetMonitorHandle() MonitorHandler {
	return defaultmonitorHandler
}

//GetShareHandle get share handle
func GetShareHandle() *share.ServiceShareHandle {
	return shareHandler
}

//GetPluginShareHandle get plugin share handle
func GetPluginShareHandle() *share.PluginShareHandle {
	return pluginShareHandler
}

//GetServiceManager get manager
func GetServiceManager() ServiceHandler {
	return defaultServieHandler
}

var defaultPluginHandler PluginHandler

//GetPluginManager get manager
func GetPluginManager() PluginHandler {
	return defaultPluginHandler
}

var defaultTenantHandler TenantHandler

//GetTenantManager get manager
func GetTenantManager() TenantHandler {
	return defaultTenantHandler
}

var defaultNetRulesHandler NetRulesHandler

//GetRulesManager get manager
func GetRulesManager() NetRulesHandler {
	return defaultNetRulesHandler
}

var defaultCloudHandler CloudHandler

//GetCloudManager get manager
func GetCloudManager() CloudHandler {
	return defaultCloudHandler
}

var defaultEventHandler EventHandler

//GetEventHandler get event handler
func GetEventHandler() EventHandler {
	return defaultEventHandler
}

var defaultAppHandler *AppAction

//GetAppHandler GetAppHandler
func GetAppHandler() *AppAction {
	return defaultAppHandler
}

var defaultAPPBackupHandler *group.BackupHandle

//GetAPPBackupHandler GetAPPBackupHandler
func GetAPPBackupHandler() *group.BackupHandle {
	return defaultAPPBackupHandler
}

var defaultGatewayHandler GatewayHandler

// GetGatewayHandler returns a default GatewayHandler
func GetGatewayHandler() GatewayHandler {
	return defaultGatewayHandler
}

var def3rdPartySvcHandler *ThirdPartyServiceHanlder

// Get3rdPartySvcHandler returns the defalut ThirdParthServiceHanlder
func Get3rdPartySvcHandler() *ThirdPartyServiceHanlder {
	return def3rdPartySvcHandler
}

var batchOperationHandler *BatchOperationHandler

//GetBatchOperationHandler get handler
func GetBatchOperationHandler() *BatchOperationHandler {
	return batchOperationHandler
}

var operationHandler *OperationHandler

//GetOperationHandler get handler
func GetOperationHandler() *OperationHandler {
	return operationHandler
}

var defaultAppRestoreHandler AppRestoreHandler

// GetAppRestoreHandler returns a default AppRestoreHandler
func GetAppRestoreHandler() AppRestoreHandler {
	return defaultAppRestoreHandler
}

var defPodHandler PodHandler

// GetPodHandler returns the defalut PodHandler
func GetPodHandler() PodHandler {
	return defPodHandler
}

var defaultEtcdHandler *EtcdHandler

// GetEtcdHandler returns the default etcd handler.
func GetEtcdHandler() *EtcdHandler {
	return defaultEtcdHandler
}

var defClusterHandler ClusterHandler

// GetClusterHandler returns the default cluster handler.
func GetClusterHandler() ClusterHandler {
	return defClusterHandler
}

var defApplicationHandler ApplicationHandler

// GetApplicationHandler  returns the default tenant application handler.
func GetApplicationHandler() ApplicationHandler {
	return defApplicationHandler
}

var defServiceEventHandler *ServiceEventHandler

// GetServiceEventHandler -
func GetServiceEventHandler() *ServiceEventHandler {
	return defServiceEventHandler
}
