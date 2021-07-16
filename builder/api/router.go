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

package api

import (
	"net/http"

	"github.com/gridworkz/kato/builder/sources"

	"github.com/gridworkz/kato/util"

	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/builder/api/controller"
	httputil "github.com/gridworkz/kato/util/http"
)

func APIServer() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v2/builder", func(r chi.Router) {
		r.Get("/publickey/{tenant_id}", func(w http.ResponseWriter, r *http.Request) {
			tenantId := strings.TrimSpace(chi.URLParam(r, "tenant_id"))
			key := sources.GetPublicKey(tenantId)
			bean := struct {
				Key string `json:"public_key"`
			}{
				Key: key,
			}
			httputil.ReturnSuccess(r, w, bean)
		})
		r.Route("/codecheck", func(r chi.Router) {
			r.Post("/", controller.AddCodeCheck)
			r.Put("/service/{serviceID}", controller.Update)
			r.Get("/service/{serviceID}", controller.GetCodeCheck)
		})
		r.Route("/version", func(r chi.Router) {
			r.Post("/", controller.UpdateDeliveredPath)
			r.Get("/event/{eventID}", controller.GetVersionByEventID)
			r.Post("/event/{eventID}", controller.UpdateVersionByEventID)
			r.Get("/service/{serviceID}", controller.GetVersionByServiceID)
			r.Delete("/service/{eventID}", controller.DeleteVersionByEventID)
		})
		r.Route("/event", func(r chi.Router) {
			r.Get("/", controller.GetEventsByIds)
		})
		r.Route("/health", func(r chi.Router) {
			r.Get("/", controller.CheckHalth)
		})
	})
	util.ProfilerSetup(r)
	return r
}
