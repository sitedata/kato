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

package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ContextKey ctx key type
type ContextKey string

var pool []string

func init() {
	pool = []string{
		"services_status",
	}
}

//InitTenant - implement middleware
func InitTenant(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		debugRequestBody(r)

		tenantName := chi.URLParam(r, "tenant_name")
		if tenantName == "" {
			httputil.ReturnError(r, w, 404, "cant find tenant")
			return
		}
		tenant, err := db.GetManager().TenantDao().GetTenantIDByName(tenantName)
		if err != nil {
			logrus.Errorf("get tenant by tenantName error: %s %v", tenantName, err)
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				httputil.ReturnError(r, w, 404, "cant find tenant")
				return
			}
			httputil.ReturnError(r, w, 500, "get assign tenant uuid failed")
			return
		}
		ctx := context.WithValue(r.Context(), ContextKey("tenant_name"), tenantName)
		ctx = context.WithValue(ctx, ContextKey("tenant_id"), tenant.UUID)
		ctx = context.WithValue(ctx, ContextKey("tenant"), tenant)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

//InitService - implement serviceinit middleware
func InitService(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		serviceAlias := chi.URLParam(r, "service_alias")
		if serviceAlias == "" {
			httputil.ReturnError(r, w, 404, "cant find service alias")
			return
		}
		tenantID := r.Context().Value(ContextKey("tenant_id"))
		service, err := db.GetManager().TenantServiceDao().GetServiceByTenantIDAndServiceAlias(tenantID.(string), serviceAlias)
		if err != nil {
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				httputil.ReturnError(r, w, 404, "cant find service")
				return
			}
			logrus.Errorf("get service by tenant & service alias error, %v", err)
			httputil.ReturnError(r, w, 500, "get service id error")
			return
		}
		serviceID := service.ServiceID
		ctx := context.WithValue(r.Context(), ContextKey("service_alias"), serviceAlias)
		ctx = context.WithValue(ctx, ContextKey("service_id"), serviceID)
		ctx = context.WithValue(ctx, ContextKey("service"), service)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// InitApplication -
func InitApplication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		appID := chi.URLParam(r, "app_id")
		tenantApp, err := handler.GetApplicationHandler().GetAppByID(appID)
		if err != nil {
			httputil.ReturnBcodeError(r, w, err)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("app_id"), tenantApp.AppID)
		ctx = context.WithValue(ctx, ContextKey("application"), tenantApp)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

//InitPlugin - implement plugin init middleware
func InitPlugin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		debugRequestBody(r)

		pluginID := chi.URLParam(r, "plugin_id")
		tenantID := r.Context().Value(ContextKey("tenant_id")).(string)
		if pluginID == "" {
			httputil.ReturnError(r, w, 404, "need plugin id")
			return
		}
		_, err := db.GetManager().TenantPluginDao().GetPluginByID(pluginID, tenantID)
		if err != nil {
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				httputil.ReturnError(r, w, 404, "cant find plugin")
				return
			}
			logrus.Errorf("get plugin error, %v", err)
			httputil.ReturnError(r, w, 500, "get plugin error")
			return
		}
		ctx := context.WithValue(r.Context(), ContextKey("plugin_id"), pluginID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

//SetLog SetLog
func SetLog(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "event_id")
		if eventID != "" {
			logger := event.GetManager().GetLogger(eventID)
			ctx := context.WithValue(r.Context(), ContextKey("logger"), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
	return http.HandlerFunc(fn)
}

//Proxy - reverse proxy middleware
func Proxy(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/v2/nodes") {
			handler.GetNodeProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/v2/cluster/service-health") {
			handler.GetNodeProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/v2/builder") {
			handler.GetBuilderProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/v2/tasks") {
			handler.GetNodeProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/v2/tasktemps") {
			handler.GetNodeProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/v2/taskgroups") {
			handler.GetNodeProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/v2/configs") {
			handler.GetNodeProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/v2/rules") {
			handler.GetMonitorProxy().Proxy(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/kubernetes/dashboard") {
			proxy := handler.GetKubernetesDashboardProxy()
			r.URL.Path = strings.Replace(r.URL.Path, "/kubernetes/dashboard", "", 1)
			proxy.Proxy(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func apiExclude(r *http.Request) bool {
	if r.Method == "GET" {
		return true
	}
	for _, item := range pool {
		if strings.Contains(r.RequestURI, item) {
			return true
		}
	}
	return false
}

type resWriter struct {
	origWriter http.ResponseWriter
	statusCode int
}

func (w *resWriter) Header() http.Header {
	return w.origWriter.Header()
}
func (w *resWriter) Write(p []byte) (int, error) {
	return w.origWriter.Write(p)
}
func (w *resWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.origWriter.WriteHeader(statusCode)
}

// WrapEL wrap eventlog, handle event log before and after process
func WrapEL(f http.HandlerFunc, target, optType string, synType int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logrus.Warningf("error reading request body: %v", err)
			} else {
				logrus.Debugf("method: %s; uri: %s; body: %s", r.Method, r.RequestURI, string(body))
			}
			// set a new body, which will simulate the same data we read
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			var targetID string
			var ok bool
			if targetID, ok = r.Context().Value(ContextKey("service_id")).(string); !ok {
				var reqDataMap map[string]interface{}
				if err = json.Unmarshal(body, &reqDataMap); err != nil {
					httputil.ReturnError(r, w, 400, "operation object is not specified")
					return
				}

				if targetID, ok = reqDataMap["service_id"].(string); !ok {
					httputil.ReturnError(r, w, 400, "operation object is not specified")
					return
				}
			}

			//eventLog check the latest event
			if !util.CanDoEvent(optType, synType, target, targetID) {
				logrus.Errorf("operation too frequent. uri: %s; target: %s; target id: %s", r.RequestURI, target, targetID)
				httputil.ReturnError(r, w, 409, "the operation is too frequent, please try again later") // status code 409 conflict
				return
			}

			// handle operator
			var operator string
			var reqData map[string]interface{}
			if err = json.Unmarshal(body, &reqData); err == nil {
				if operatorI := reqData["operator"]; operatorI != nil {
					operator = operatorI.(string)
				}
			}

			// tenantID cannot be null
			tenantID := r.Context().Value(ContextKey("tenant_id")).(string)
			var ctx context.Context

			event, err := util.CreateEvent(target, optType, targetID, tenantID, string(body), operator, synType)
			if err != nil {
				logrus.Error("create event error : ", err)
				httputil.ReturnError(r, w, 500, "operation failed")
				return
			}
			ctx = context.WithValue(r.Context(), ContextKey("event"), event)
			ctx = context.WithValue(ctx, ContextKey("event_id"), event.EventID)
			rw := &resWriter{origWriter: w}
			f(rw, r.WithContext(ctx))
			if synType == dbmodel.SYNEVENTTYPE || (synType == dbmodel.ASYNEVENTTYPE && rw.statusCode >= 400) { // status code 2XX/3XX all equal to success
				util.UpdateEvent(event.EventID, rw.statusCode)
			}
		}
	}
}

func debugRequestBody(r *http.Request) {
	if !apiExclude(r) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Warningf("error reading request body: %v", err)
		}
		logrus.Debugf("method: %s; uri: %s; body: %s", r.Method, r.RequestURI, string(body))

		// set a new body, which will simulate the same data we read
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
}
