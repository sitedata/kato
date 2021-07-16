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

	"github.com/gridworkz/kato/api/util"

	"github.com/go-chi/chi"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/sirupsen/logrus"
)

//ServiceDiscover
func ServiceDiscover(w http.ResponseWriter, r *http.Request) {
	serviceInfo := chi.URLParam(r, "service_name")
	sds, err := discoverService.DiscoverService(serviceInfo)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnNoFomart(r, w, 200, sds)
}

//ListenerDiscover
func ListenerDiscover(w http.ResponseWriter, r *http.Request) {
	tenantService := chi.URLParam(r, "tenant_service")
	serviceNodes := chi.URLParam(r, "service_nodes")
	lds, err := discoverService.DiscoverListeners(tenantService, serviceNodes)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnNoFomart(r, w, 200, lds)
}

//ClusterDiscover
func ClusterDiscover(w http.ResponseWriter, r *http.Request) {
	tenantService := chi.URLParam(r, "tenant_service")
	serviceNodes := chi.URLParam(r, "service_nodes")
	cds, err := discoverService.DiscoverClusters(tenantService, serviceNodes)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnNoFomart(r, w, 200, cds)
}

//RoutesDiscover
//no impl
func RoutesDiscover(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "tenant_id")
	serviceNodes := chi.URLParam(r, "service_nodes")
	routeConfig := chi.URLParam(r, "route_config")
	logrus.Debugf("route_config is %s, namespace %s, serviceNodes %s", routeConfig, namespace, serviceNodes)
	w.WriteHeader(200)
}

//PluginResourcesConfig
func PluginResourcesConfig(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "tenant_id")
	serviceAlias := chi.URLParam(r, "service_alias")
	pluginID := chi.URLParam(r, "plugin_id")
	ss, err := discoverService.GetPluginConfigs(namespace, serviceAlias, pluginID)
	if err != nil {
		util.CreateAPIHandleError(500, err).Handle(r, w)
		return
	}
	if ss == nil {
		util.CreateAPIHandleError(404, err).Handle(r, w)
	}
	httputil.ReturnNoFomart(r, w, 200, ss)
}
