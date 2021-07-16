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
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/node/api/controller"
	"github.com/gridworkz/kato/util/log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//Routers
func Routers(mode string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID) //Register an id in the context of each request
	//Sets a http.Request's RemoteAddr to either X-Forwarded-For or X-Real-IP
	r.Use(middleware.RealIP)
	//Logs the start and end of each request with the elapsed processing time
	logger := logrus.New()
	logger.SetLevel(logrus.GetLevel())
	r.Use(log.NewStructuredLogger(logger))
	//Gracefully absorb panics and prints the stack trace
	r.Use(middleware.Recoverer)
	//request time out
	r.Use(middleware.Timeout(time.Second * 5))
	r.Mount("/v1", DisconverRoutes())
	r.Route("/v2", func(r chi.Router) {
		r.Get("/ping", controller.Ping)
		r.Route("/apps", func(r chi.Router) {
			r.Get("/{app_name}/register", controller.APPRegister)
			r.Get("/{app_name}/discover", controller.APPDiscover)
			r.Get("/", controller.APPList)
		})
		r.Route("/localvolumes", func(r chi.Router) {
			r.Post("/create", controller.CreateLocalVolume)
			r.Delete("/", controller.DeleteLocalVolume)
		})
		//The following APIs are only available for management nodes
		if mode == "master" {
			r.Route("/configs", func(r chi.Router) {
				r.Get("/datacenter", controller.GetDatacenterConfig)
				r.Put("/datacenter", controller.PutDatacenterConfig)
			})
			r.Route("/cluster", func(r chi.Router) {
				r.Get("/", controller.ClusterInfo)
				r.Get("/service-health", controller.GetServicesHealthy)
			})
			r.Route("/nodes", func(r chi.Router) {
				// abandoned
				r.Get("/fullres", controller.ClusterInfo)
				r.Get("/{node_id}/node_resource", controller.GetNodeResource)
				r.Get("/resources", controller.Resources)
				r.Get("/capres", controller.CapRes)
				r.Get("/", controller.GetNodes)
				r.Get("/all_node_health", controller.GetAllNodeHealth)
				r.Get("/rule/{rule}", controller.GetRuleNodes)
				r.Get("/{node_id}", controller.GetNode)
				r.Put("/{node_id}/status", controller.UpdateNodeStatus)
				r.Put("/{node_id}/unschedulable", controller.Cordon)
				r.Put("/{node_id}/reschedulable", controller.UnCordon)
				r.Post("/{node_id}/labels", controller.PutLabel)
				r.Get("/{node_id}/labels", controller.GetLabel)
				r.Delete("/{node_id}/labels", controller.DeleteLabel)
				r.Post("/{node_id}/down", controller.DownNode)
				r.Post("/{node_id}/up", controller.UpNode)
				r.Get("/{node_id}/instance", controller.Instances)
				r.Get("/{node_id}/check", controller.CheckNode)
				r.Get("/{node_id}/resource", controller.Resource)
				r.Get("/{node_id}/conditions", controller.ListNodeCondition)
				r.Delete("/{node_id}/conditions/{condition}", controller.DeleteNodeCondition)
				// about node install
				r.Post("/{node_id}/install", controller.InstallNode)  //install node
				r.Post("/", controller.AddNode)                       //add node
				r.Delete("/{node_id}", controller.DeleteKatoNode) //delete node
			})
		}
	})
	return r
}
