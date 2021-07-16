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
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	"github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/db/errors"
	httputil "github.com/gridworkz/kato/util/http"
)

// AutoscalerRules -
func (t *TenantStruct) AutoscalerRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		t.addAutoscalerRule(w, r)
	case "PUT":
		t.updAutoscalerRule(w, r)
	}
}

func (t *TenantStruct) addAutoscalerRule(w http.ResponseWriter, r *http.Request) {
	var req model.AutoscalerRuleReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	req.ServiceID = serviceID
	if err := handler.GetServiceManager().AddAutoscalerRule(&req); err != nil {
		if err == errors.ErrRecordAlreadyExist {
			httputil.ReturnError(r, w, 400, err.Error())
			return
		}
		logrus.Errorf("add autoscaler rule: %v", err)
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}

	httputil.ReturnSuccess(r, w, nil)
}

func (t *TenantStruct) updAutoscalerRule(w http.ResponseWriter, r *http.Request) {
	var req model.AutoscalerRuleReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	if err := handler.GetServiceManager().UpdAutoscalerRule(&req); err != nil {
		if err == errors.ErrRecordAlreadyExist {
			httputil.ReturnError(r, w, 400, err.Error())
			return
		}
		if err == gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 404, err.Error())
			return
		}
		logrus.Errorf("update autoscaler rule: %v", err)
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}

	httputil.ReturnSuccess(r, w, nil)
}

// ScalingRecords -
func (t *TenantStruct) ScalingRecords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t.listScalingRecords(w, r)
	}
}

func (t *TenantStruct) listScalingRecords(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logrus.Warningf("convert '%s(pageStr)' to int: %v", pageStr, err)
	}
	if page <= 0 {
		page = 1
	}

	pageSizeStr := r.URL.Query().Get("page_size")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		logrus.Warningf("convert '%s(pageSizeStr)' to int: %v", pageSizeStr, err)
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	records, count, err := handler.GetServiceManager().ListScalingRecords(serviceID, page, pageSize)
	if err != nil {
		logrus.Errorf("list scaling rule: %v", err)
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}

	httputil.ReturnSuccess(r, w, map[string]interface{}{
		"total": count,
		"data":  records,
	})
}
