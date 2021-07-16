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
	api_model "github.com/gridworkz/kato/api/model"
	httputil "github.com/gridworkz/kato/util/http"
)

//CloudManager
type CloudManager struct{}

var defaultCloudManager *CloudManager

//GetCloudRouterManager
func GetCloudRouterManager() *CloudManager {
	if defaultCloudManager != nil {
		return defaultCloudManager
	}
	defaultCloudManager = &CloudManager{}
	return defaultCloudManager
}

//Show
func (c *CloudManager) Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cloud urls"))
}

//CreateToken CreateToken
// swagger:operation POST /cloud/auth cloud createToken
//
// Generate token
//
// create token
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
func (c *CloudManager) CreateToken(w http.ResponseWriter, r *http.Request) {
	var gt api_model.GetUserToken
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &gt.Body, nil); !ok {
		return
	}
	ti, err := handler.GetCloudManager().TokenDispatcher(&gt)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, ti)
}

//GetTokenInfo
// swagger:operation GET /cloud/auth/{eid} cloud getTokenInfo
//
// Get tokeninfo
//
// get token info
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
func (c *CloudManager) GetTokenInfo(w http.ResponseWriter, r *http.Request) {
	eid := chi.URLParam(r, "eid")
	ti, err := handler.GetCloudManager().GetTokenInfo(eid)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, ti)
}

//UpdateToken
// swagger:operation PUT /cloud/auth/{eid} cloud updateToken
//
// Update token
//
// update token info
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
func (c *CloudManager) UpdateToken(w http.ResponseWriter, r *http.Request) {
	var ut api_model.UpdateToken
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ut.Body, nil); !ok {
		return
	}
	eid := chi.URLParam(r, "eid")
	err := handler.GetCloudManager().UpdateTokenTime(eid, ut.Body.ValidityPeriod)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//GetAPIManager
// swagger:operation GET /cloud/api/manager cloud GetAPIManager
//
// Get api manager
//
// get api manager
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
func (c *CloudManager) GetAPIManager(w http.ResponseWriter, r *http.Request) {
	apiMap := handler.GetTokenIdenHandler().GetAPIManager()
	httputil.ReturnSuccess(r, w, apiMap)
}

//AddAPIManager
// swagger:operation POST /cloud/api/manager cloud addAPIManager
//
// Add api manager
//
// get api manager
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
func (c *CloudManager) AddAPIManager(w http.ResponseWriter, r *http.Request) {
	var am api_model.APIManager
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &am.Body, nil); !ok {
		return
	}
	err := handler.GetTokenIdenHandler().AddAPIManager(&am)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeleteAPIManager
// swagger:operation DELETE /cloud/api/manager cloud deleteAPIManager
//
// Delete api manager
//
// delete api manager
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
func (c *CloudManager) DeleteAPIManager(w http.ResponseWriter, r *http.Request) {
	var am api_model.APIManager
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &am.Body, nil); !ok {
		return
	}
	err := handler.GetTokenIdenHandler().DeleteAPIManager(&am)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}
