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
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	httputil "github.com/gridworkz/kato/util/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/proxy"
)

//EventLogStruct eventlog struct
type EventLogStruct struct {
	EventlogServerProxy proxy.Proxy
}

//HistoryLogs get service history logs
//proxy
func (e *EventLogStruct) HistoryLogs(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	serviceAlias := r.Context().Value(middleware.ContextKey("service_alias")).(string)
	name, _ := handler.GetEventHandler().GetLogInstance(serviceID)
	if name != "" {
		r.URL.Query().Add("host_id", name)
		r = r.WithContext(context.WithValue(r.Context(), proxy.ContextKey("host_id"), name))
	}
	//Replace service alias to service id in path
	r.URL.Path = strings.Replace(r.URL.Path, serviceAlias, serviceID, 1)
	r.URL.Path = strings.Replace(r.URL.Path, "/v2/", "/", 1)
	e.EventlogServerProxy.Proxy(w, r)
}

//LogList GetLogList
func (e *EventLogStruct) LogList(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET  /v2/tenants/{tenant_name}/services/{service_alias}/log-file v2 logList
	//
	// Get application log list
	//
	// get log list
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	fileList, err := handler.GetEventHandler().GetLogList(GetServiceAliasID(serviceID))
	if err != nil {
		if os.IsNotExist(err) {
			httputil.ReturnError(r, w, 404, err.Error())
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, fileList)
	return
}

//LogFile GetLogFile
func (e *EventLogStruct) LogFile(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/log-file/{file_name} v2 logFile
	//
	// Download application specific log
	//
	// get log file
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format

	fileName := chi.URLParam(r, "file_name")
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	logPath, _, err := handler.GetEventHandler().GetLogFile(GetServiceAliasID(serviceID), fileName)
	if err != nil {
		if os.IsNotExist(err) {
			httputil.ReturnError(r, w, 404, err.Error())
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	http.StripPrefix(fileName, http.FileServer(http.Dir(logPath)))
	//fs.ServeHTTP(w, r)
}

//LogSocket GetLogSocket
func (e *EventLogStruct) LogSocket(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/log-instance v2 logSocket
	//
	// Get application log web-socket instance
	//
	// get log socket
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	value, err := handler.GetEventHandler().GetLogInstance(serviceID)
	if err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			httputil.ReturnError(r, w, 404, err.Error())
			return
		}
		logrus.Errorf("get docker log instance error. %s", err.Error())
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	rc := make(map[string]string)
	rc["host_id"] = value
	httputil.ReturnSuccess(r, w, rc)
	return
}

//LogByAction GetLogByAction
func (e *EventLogStruct) LogByAction(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/event-log v2 logByAction
	//
	// Obtain the operation log of the specified operation
	//
	// get log by level
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format
	var elog api_model.LogByLevelStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &elog.Body, nil)
	if !ok {
		return
	}
	dl, err := handler.GetEventHandler().GetLevelLog(elog.Body.EventID, elog.Body.Level)
	if err != nil {
		logrus.Errorf("get event log error, %v", err)
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, dl.Data)
	return
}

//TenantLogByAction GetTenantLogByAction
// swagger:operation POST /v2/tenants/{tenant_name}/event-log v2 tenantLogByAction
//
// Obtain the operation log of the specified operation
//
// get tenant log by level
//
// ---
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (e *EventLogStruct) TenantLogByAction(w http.ResponseWriter, r *http.Request) {
	var elog api_model.TenantLogByLevelStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &elog.Body, nil)
	if !ok {
		return
	}
	logrus.Info(elog.Body.Level)
	dl, err := handler.GetEventHandler().GetLevelLog(elog.Body.EventID, elog.Body.Level)
	if err != nil {
		logrus.Errorf("get tenant event log error, %v", err)
		httputil.ReturnError(r, w, 200, "success")
		return
	}
	httputil.ReturnSuccess(r, w, dl.Data)
	return
}

//Events get log by target
func (e *EventLogStruct) Events(w http.ResponseWriter, r *http.Request) {
	target := r.FormValue("target")
	targetID := r.FormValue("target-id")
	var page, size int
	var err error
	if page, err = strconv.Atoi(r.FormValue("page")); err != nil || page <= 0 {
		page = 1
	}
	if size, err = strconv.Atoi(r.FormValue("size")); err != nil || size <= 0 {
		size = 10
	}
	logrus.Debugf("get event page param[target:%s id:%s page:%d, page_size:%d]", target, targetID, page, size)
	list, total, err := handler.GetEventHandler().GetEvents(target, targetID, page, size)
	if err != nil {
		logrus.Errorf("get event log error, %v", err)
		httputil.ReturnError(r, w, 500, "get log error")
		return
	}
	// format start and end time
	for i := range list {
		if list[i].EndTime != "" && len(list[i].EndTime) > 20 {
			list[i].EndTime = strings.Replace(list[i].EndTime[0:19]+"+08:00", " ", "T", 1)
		}
	}
	httputil.ReturnList(r, w, total, page, list)
}

//EventLog get event log by eventID
func (e *EventLogStruct) EventLog(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventID")
	if strings.TrimSpace(eventID) == "" {
		httputil.ReturnError(r, w, 400, "eventID is request")
		return
	}

	dl, err := handler.GetEventHandler().GetLevelLog(eventID, "debug")
	if err != nil {
		logrus.Errorf("get event log error, %v", err)
		httputil.ReturnError(r, w, 500, "read event log error: "+err.Error())
		return
	}

	httputil.ReturnSuccess(r, w, dl.Data)
	return
}
