// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/model"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/gridworkz/kato/worker/server"
	"github.com/sirupsen/logrus"
)

// PodController is an implementation of PodInterface
type PodController struct{}

//Pods get some service pods
// swagger:operation GET /v2/tenants/{tenant_name}/pods v2/tenants pods
//
// Get the Pod information of some apps
//
// get some service pods
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: get some service pods
func Pods(w http.ResponseWriter, r *http.Request) {
	serviceIDs := strings.Split(r.FormValue("service_ids"), ",")
	if serviceIDs == nil || len(serviceIDs) == 0 {
		tenant := r.Context().Value(middleware.ContextKey("tenant")).(*model.Tenants)
		services, _ := db.GetManager().TenantServiceDao().GetServicesByTenantID(tenant.UUID)
		for _, s := range services {
			serviceIDs = append(serviceIDs, s.ServiceID)
		}
	}
	var allpods []*handler.K8sPodInfo
	podinfo, err := handler.GetServiceManager().GetMultiServicePods(serviceIDs)
	if err != nil {
		logrus.Errorf("get service pod failure %s", err.Error())
	}
	if podinfo != nil {
		var pods []*handler.K8sPodInfo
		if podinfo.OldPods != nil {
			pods = append(podinfo.NewPods, podinfo.OldPods...)
		} else {
			pods = podinfo.NewPods
		}
		for _, pod := range pods {
			allpods = append(allpods, pod)
		}
	}
	httputil.ReturnSuccess(r, w, allpods)
}

// PodDetail -
func (p *PodController) PodDetail(w http.ResponseWriter, r *http.Request) {
	podName := chi.URLParam(r, "pod_name")
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	pd, err := handler.GetPodHandler().PodDetail(serviceID, podName)
	if err != nil {
		logrus.Errorf("error getting pod detail: %v", err)
		if err == server.ErrPodNotFound {
			httputil.ReturnError(r, w, 404, fmt.Sprintf("error getting pod detail: %v", err))
			return
		}
		httputil.ReturnError(r, w, 500, fmt.Sprintf("error getting pod detail: %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, pd)
}
