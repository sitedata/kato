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

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	api_model "github.com/gridworkz/kato/api/model"
	httputil "github.com/gridworkz/kato/util/http"
)

//Share Application sharing
func (t *TenantStruct) Share(w http.ResponseWriter, r *http.Request) {
	//Share ShareService
	// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/share  v2 shareService
	//
	// Share application media
	//
	// share service
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
	//     description: Unified return format
	var ccs api_model.ServiceShare
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ccs.Body, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	ccs.Body.EventID = r.Context().Value(middleware.ContextKey("event_id")).(string)
	res, errS := handler.GetShareHandle().Share(serviceID, ccs)
	if errS != nil {
		errS.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, res)
}

//ShareResult Get shared results
func (t *TenantStruct) ShareResult(w http.ResponseWriter, r *http.Request) {
	//ShareResult ShareResult
	// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/share  v2 get_share_result
	//
	// Get results of sharing application media
	//
	// share service
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
	//     description: Unified return format
	shareID := chi.URLParam(r, "share_id")
	res, errS := handler.GetShareHandle().ShareResult(shareID)
	if errS != nil {
		errS.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, res)
}
