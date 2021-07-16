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
	"github.com/gridworkz/kato/api/discover"
	"github.com/gridworkz/kato/api/proxy"
	"github.com/gridworkz/kato/cmd/api/option"
)

var nodeProxy proxy.Proxy
var builderProxy proxy.Proxy
var prometheusProxy proxy.Proxy
var monitorProxy proxy.Proxy
var kubernetesDashboard proxy.Proxy

//InitProxy
func InitProxy(conf option.Config) {
	if nodeProxy == nil {
		nodeProxy = proxy.CreateProxy("acp_node", "http", conf.NodeAPI)
		discover.GetEndpointDiscover().AddProject("acp_node", nodeProxy)
	}
	if builderProxy == nil {
		builderProxy = proxy.CreateProxy("builder", "http", conf.BuilderAPI)
	}
	if prometheusProxy == nil {
		prometheusProxy = proxy.CreateProxy("prometheus", "http", []string{conf.PrometheusEndpoint})
	}
	if monitorProxy == nil {
		monitorProxy = proxy.CreateProxy("monitor", "http", []string{"127.0.0.1:3329"})
		discover.GetEndpointDiscover().AddProject("monitor", monitorProxy)
	}
	if kubernetesDashboard == nil {
		kubernetesDashboard = proxy.CreateProxy("kubernetesdashboard", "http", []string{conf.KuberentesDashboardAPI})
	}
}

//GetNodeProxy
func GetNodeProxy() proxy.Proxy {
	return nodeProxy
}

//GetBuilderProxy
func GetBuilderProxy() proxy.Proxy {
	return builderProxy
}

//GetPrometheusProxy
func GetPrometheusProxy() proxy.Proxy {
	return prometheusProxy
}

//GetMonitorProxy
func GetMonitorProxy() proxy.Proxy {
	return monitorProxy
}

// GetKubernetesDashboardProxy returns the kubernetes dashboard proxy.
func GetKubernetesDashboardProxy() proxy.Proxy {
	return kubernetesDashboard
}
