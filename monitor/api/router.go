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

	"github.com/gridworkz/kato/util"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/monitor/api/controller"
	httputil "github.com/gridworkz/kato/util/http"
)

//Server api server
func Server(c *controller.RuleControllerManager) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/monitor", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			bean := map[string]string{"status": "health", "info": "monitor service health"}
			httputil.ReturnSuccess(r, w, bean)
		})
	})
	r.Route("/v2/rules", func(r chi.Router) {
		r.Post("/", c.AddRules)
		r.Put("/{rules_name}", c.RegRules)
		r.Delete("/{rules_name}", c.DelRules)
		r.Get("/{rules_name}", c.GetRules)
		r.Get("/all", c.GetAllRules)
	})
	util.ProfilerSetup(r)
	return r
}
