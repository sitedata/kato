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
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/node/api"
	httputil "github.com/gridworkz/kato/util/http"
)

//SetAPIRoute
func (m *ManagerService) SetAPIRoute(apim *api.Manager) error {
	apim.GetRouter().Post("/services/{service_name}/stop", m.StopServiceAPI)
	apim.GetRouter().Post("/services/{service_name}/start", m.StartServiceAPI)
	apim.GetRouter().Post("/services/update", m.UpdateConfigAPI)
	return nil
}

//StartServiceAPI
func (m *ManagerService) StartServiceAPI(w http.ResponseWriter, r *http.Request) {
	serviceName := strings.TrimSpace(chi.URLParam(r, "service_name"))
	if err := m.StartService(serviceName); err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
	}
	httputil.ReturnSuccess(r, w, nil)
}

//StopServiceAPI
func (m *ManagerService) StopServiceAPI(w http.ResponseWriter, r *http.Request) {
	serviceName := strings.TrimSpace(chi.URLParam(r, "service_name"))
	if err := m.StopService(serviceName); err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
	}
	httputil.ReturnSuccess(r, w, nil)
}

//UpdateConfigAPI
func (m *ManagerService) UpdateConfigAPI(w http.ResponseWriter, r *http.Request) {
	if err := m.ReLoadServices(); err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
	}
	httputil.ReturnSuccess(r, w, nil)
}
