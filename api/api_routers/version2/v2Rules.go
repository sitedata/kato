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

package version2

import (
	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/controller"
)

//PluginRouter plugin router
func (v2 *V2) rulesRouter() chi.Router {
	r := chi.NewRouter()
	// service rule
	// url: v2/tenant/{tenant_name}/services/{service_alias}/net-rule/xxx
	// -- --
	//downstream
	r.Mount("/downstream", v2.downstreamRouter())
	//upstream
	r.Mount("/upstream", v2.upstreamRouter())
	return r
}

func (v2 *V2) downstreamRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", controller.GetManager().SetDownStreamRule)
	//r.Get("/")
	r.Get("/{dest_service_alias}/{port}", controller.GetManager().GetDownStreamRule)
	r.Put("/{dest_service_alias}/{port}", controller.GetManager().UpdateDownStreamRule)
	//r.Delete("/{service_alias}/{port}")
	return r
}

func (v2 *V2) upstreamRouter() chi.Router {
	r := chi.NewRouter()
	return r
}
