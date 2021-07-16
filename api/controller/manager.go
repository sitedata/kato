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

package controller

import (
	"net/http"

	"github.com/gridworkz/kato/api/api"
	"github.com/gridworkz/kato/api/discover"
	"github.com/gridworkz/kato/api/proxy"
	"github.com/gridworkz/kato/cmd/api/option"
	mqclient "github.com/gridworkz/kato/mq/client"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/gridworkz/kato/worker/client"
)

//V2Manager
type V2Manager interface {
	Show(w http.ResponseWriter, r *http.Request)
	Health(w http.ResponseWriter, r *http.Request)
	AlertManagerWebHook(w http.ResponseWriter, r *http.Request)
	Version(w http.ResponseWriter, r *http.Request)
	api.ClusterInterface
	api.TenantInterface
	api.ServiceInterface
	api.LogInterface
	api.PluginInterface
	api.RulesInterface
	api.AppInterface
	api.Gatewayer
	api.ThirdPartyServicer
	api.Labeler
	api.AppRestoreInterface
	api.PodInterface
	api.ApplicationInterface
}

var defaultV2Manager V2Manager

//CreateV2RouterManager- Create manager
func CreateV2RouterManager(conf option.Config, statusCli *client.AppRuntimeSyncClient) (err error) {
	defaultV2Manager, err = NewManager(conf, statusCli)
	return err
}

//GetManager - Get manager
func GetManager() V2Manager {
	return defaultV2Manager
}

//NewManager - New manager
func NewManager(conf option.Config, statusCli *client.AppRuntimeSyncClient) (*V2Routes, error) {
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: conf.EtcdEndpoint,
		CaFile:    conf.EtcdCaFile,
		CertFile:  conf.EtcdCertFile,
		KeyFile:   conf.EtcdKeyFile,
	}
	mqClient, err := mqclient.NewMqClient(etcdClientArgs, conf.MQAPI)
	if err != nil {
		return nil, err
	}
	var v2r V2Routes
	v2r.TenantStruct.StatusCli = statusCli
	v2r.TenantStruct.MQClient = mqClient
	v2r.GatewayStruct.MQClient = mqClient
	v2r.GatewayStruct.cfg = &conf
	v2r.LabelController.optconfig = &conf
	eventServerProxy := proxy.CreateProxy("eventlog", "http", []string{"local=>rbd-eventlog:6363"})
	discover.GetEndpointDiscover().AddProject("event_log_event_http", eventServerProxy)
	v2r.EventLogStruct.EventlogServerProxy = eventServerProxy
	return &v2r, nil
}
