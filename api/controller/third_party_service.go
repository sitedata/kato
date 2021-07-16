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
	"encoding/json"
	"net/http"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	"github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/errors"
	validation "github.com/gridworkz/kato/util/endpoint"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/sirupsen/logrus"
)

// ThirdPartyServiceController implements ThirdPartyServicer
type ThirdPartyServiceController struct{}

// Endpoints POST->add endpoints, PUT->update endpoints, DELETE->delete endpoints
func (t *ThirdPartyServiceController) Endpoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		t.addEndpoints(w, r)
	case "PUT":
		t.updEndpoints(w, r)
	case "DELETE":
		t.delEndpoints(w, r)
	case "GET":
		t.listEndpoints(w, r)
	}
}

func (t *ThirdPartyServiceController) addEndpoints(w http.ResponseWriter, r *http.Request) {
	var data model.AddEndpiontsReq
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &data, nil) {
		return
	}
	// if address is not ip, and then it is domain
	address := validation.SplitEndpointAddress(data.Address)
	sid := r.Context().Value(middleware.ContextKey("service_id")).(string)
	if validation.IsDomainNotIP(address) {
		// handle domain, check can add new endpoint or not
		if !canAddDomainEndpoint(sid, true) {
			logrus.Warningf("new endpoint addres[%s] is domian", address)
			httputil.ReturnError(r, w, 400, "do not support multi domain endpoints")
			return
		}
	}
	if !canAddDomainEndpoint(sid, false) {
		// handle ip, check can add new endpoint or not
		logrus.Warningf("new endpoint address[%s] is ip, but already has domain endpoint", address)
		httputil.ReturnError(r, w, 400, "do not support multi domain endpoints")
		return
	}

	if err := handler.Get3rdPartySvcHandler().AddEndpoints(sid, &data); err != nil {
		if err == errors.ErrRecordAlreadyExist {
			httputil.ReturnError(r, w, 400, err.Error())
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, "success")
}

func canAddDomainEndpoint(sid string, isDomain bool) bool {
	endpoints, err := db.GetManager().EndpointsDao().List(sid)
	if err != nil {
		logrus.Errorf("find endpoints by sid[%s], error: %s", sid, err.Error())
		return false
	}

	if len(endpoints) > 0 && isDomain {
		return false
	}
	if !isDomain {
		for _, ep := range endpoints {
			address := validation.SplitEndpointAddress(ep.IP)
			if validation.IsDomainNotIP(address) {
				return false
			}
		}
	}
	return true
}

func (t *ThirdPartyServiceController) updEndpoints(w http.ResponseWriter, r *http.Request) {
	var data model.UpdEndpiontsReq
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &data, nil) {
		return
	}

	if err := handler.Get3rdPartySvcHandler().UpdEndpoints(&data); err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, "success")
}

func (t *ThirdPartyServiceController) delEndpoints(w http.ResponseWriter, r *http.Request) {
	var data model.DelEndpiontsReq
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &data, nil) {
		return
	}
	sid := r.Context().Value(middleware.ContextKey("service_id")).(string)
	if err := handler.Get3rdPartySvcHandler().DelEndpoints(data.EpID, sid); err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, "success")
}

func (t *ThirdPartyServiceController) listEndpoints(w http.ResponseWriter, r *http.Request) {
	sid := r.Context().Value(middleware.ContextKey("service_id")).(string)
	res, err := handler.Get3rdPartySvcHandler().ListEndpoints(sid)
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	b, _ := json.Marshal(res)
	logrus.Debugf("response endpoints: %s", string(b))
	if res == nil || len(res) == 0 {
		httputil.ReturnSuccess(r, w, []*model.EndpointResp{})
		return
	}
	httputil.ReturnSuccess(r, w, res)
}
