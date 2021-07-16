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

	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/core/config"
	"github.com/gridworkz/kato/node/core/service"
	"github.com/gridworkz/kato/node/kubecache"
	"github.com/gridworkz/kato/node/masterserver"
)

var datacenterConfig *config.DataCenterConfig
var prometheusService *service.PrometheusService
var appService *service.AppService
var nodeService *service.NodeService
var discoverService *service.DiscoverAction
var kubecli kubecache.KubeClient

//Init
func Init(c *option.Conf, ms *masterserver.MasterServer, kube kubecache.KubeClient) {
	if ms != nil {
		prometheusService = service.CreatePrometheusService(c, ms)
		datacenterConfig = config.GetDataCenterConfig()
		nodeService = service.CreateNodeService(c, ms.Cluster, kube)
	}
	appService = service.CreateAppService(c)
	discoverService = service.CreateDiscoverActionManager(c, kube)
	kubecli = kube
}

//Exist
func Exist(i interface{}) {
	if datacenterConfig != nil {
		datacenterConfig.Stop()
	}
}

//Ping
func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
