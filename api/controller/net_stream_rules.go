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

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	api_model "github.com/gridworkz/kato/api/model"
	httputil "github.com/gridworkz/kato/util/http"
)

//SetDownStreamRule Set downstream rules
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/net-rule/downstream v2 setNetDownStreamRuleStruct
//
// Set downstream network rules
//
// set NetDownStreamRuleStruct
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
func (t *TenantStruct) SetDownStreamRule(w http.ResponseWriter, r *http.Request) {
	var rs api_model.SetNetDownStreamRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &rs.Body, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	tenantName := r.Context().Value(middleware.ContextKey("tenant_name")).(string)
	serviceAlias := r.Context().Value(middleware.ContextKey("service_alias")).(string)
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	rs.TenantName = tenantName
	rs.ServiceAlias = serviceAlias
	rs.Body.Rules.ServiceID = serviceID
	if err := handler.GetRulesManager().CreateDownStreamNetRules(tenantID, &rs); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//GetDownStreamRule Get downstream rules
// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/net-rule/downstream/{dest_service_alias}/{port} v2 getNetDownStreamRuleStruct
//
// Get downstream network rules
//
// set NetDownStreamRuleStruct
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
func (t *TenantStruct) GetDownStreamRule(w http.ResponseWriter, r *http.Request) {
	serviceAlias := r.Context().Value(middleware.ContextKey("service_alias")).(string)
	destServiceAlias := r.Context().Value(middleware.ContextKey("dest_service_alias")).(string)
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	port := r.Context().Value(middleware.ContextKey("port")).(string)

	nrs, err := handler.GetRulesManager().GetDownStreamNetRule(
		tenantID,
		serviceAlias,
		destServiceAlias,
		port)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nrs)
}

//DeleteDownStreamRule Delete downstream rules
func (t *TenantStruct) DeleteDownStreamRule(w http.ResponseWriter, r *http.Request) {}

//UpdateDownStreamRule Update downstream rules
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/net-rule/downstream/{dest_service_alias}/{port} v2 updateNetDownStreamRuleStruct
//
// Update downstream network rules
//
// update NetDownStreamRuleStruct
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
func (t *TenantStruct) UpdateDownStreamRule(w http.ResponseWriter, r *http.Request) {
	var urs api_model.UpdateNetDownStreamRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &urs.Body, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	tenantName := r.Context().Value(middleware.ContextKey("tenant_name")).(string)
	serviceAlias := r.Context().Value(middleware.ContextKey("service_alias")).(string)
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	destServiceAlias := r.Context().Value(middleware.ContextKey("dest_service_alias")).(string)
	port := r.Context().Value(middleware.ContextKey("tenant_id")).(int)

	urs.DestServiceAlias = destServiceAlias
	urs.Port = port
	urs.ServiceAlias = serviceAlias
	urs.TenantName = tenantName
	urs.Body.Rules.ServiceID = serviceID

	if err := handler.GetRulesManager().UpdateDownStreamNetRule(tenantID, &urs); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}
