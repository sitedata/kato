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

package router

import (
	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/node/api/controller"
)

//DisconverRoutes envoy discover api
// v1 api will abandoned in 5.2
func DisconverRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/ping", controller.Ping)
	r.Mount("/listeners", ListenersRoutes())
	r.Mount("/clusters", ClustersRoutes())
	r.Mount("/registration", RegistrationRoutes())
	r.Mount("/routes", RoutesRouters())
	r.Mount("/resources", SourcesRoutes())
	return r
}

//ListenersRoutes listeners routes lds
//GET /v1/listeners/(string: service_cluster)/(string: service_node)
func ListenersRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/ping", controller.Ping)
	r.Get("/{tenant_service}/{service_nodes}", controller.ListenerDiscover)
	return r
}

//ClustersRoutes cds
//GET /v1/clusters/(string: service_cluster)/(string: service_node)
func ClustersRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/ping", controller.Ping)
	r.Get("/{tenant_service}/{service_nodes}", controller.ClusterDiscover)
	return r
}

//RegistrationRoutes sds
//GET /v1/registration/(string: service_name)
func RegistrationRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/ping", controller.Ping)
	r.Get("/{service_name}", controller.ServiceDiscover)
	return r
}

//RoutesRouters rds
//GET /v1/routes/(string: route_config_name)/(string: service_cluster)/(string: service_node)
func RoutesRouters() chi.Router {
	r := chi.NewRouter()
	r.Get("/ping", controller.Ping)
	r.Get("/{route_config}/{tenant_service}/{service_nodes}", controller.RoutesDiscover)
	return r
}

//SourcesRoutes SourcesRoutes
//GET /v1/resources/(string: tenant_id)/(string: service_alias)/(string: plugin_id)
func SourcesRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{tenant_id}/{service_alias}/{plugin_id}", controller.PluginResourcesConfig)
	return r
}
