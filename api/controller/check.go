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
	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	api_model "github.com/gridworkz/kato/api/model"
	httputil "github.com/gridworkz/kato/util/http"
)

//Check service check
// swagger:operation POST /v2/tenants/{tenant_name}/servicecheck v2 serviceCheck
//
// Application build source detection, support docker run ,docker compose, source code
//
// service check
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
func Check(w http.ResponseWriter, r *http.Request) {
	var gt api_model.ServiceCheckStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &gt.Body, nil); !ok {
		return
	}
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	gt.Body.TenantID = tenantID
	result, eventID, err := handler.GetServiceManager().ServiceCheck(&gt)
	if err != nil {
		err.Handle(r, w)
		return
	}
	re := struct {
		CheckUUID string `json:"check_uuid"`
		EventID   string `json:"event_id"`
	}{
		CheckUUID: result,
		EventID:   eventID,
	}
	httputil.ReturnSuccess(r, w, re)
}

//GetServiceCheckInfo get service check info
// swagger:operation GET /v2/tenants/{tenant_name}/servicecheck/{uuid} v2 getServiceCheckInfo
//
//	Obtain build inspection information
//
// get service check info
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
func GetServiceCheckInfo(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimSpace(chi.URLParam(r, "uuid"))
	si, err := handler.GetServiceManager().GetServiceCheckInfo(uuid)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, si)
}
