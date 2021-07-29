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
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/db"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//Event GetLogs
func (t *TenantStruct) Event(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET  /v2/tenants/{tenant_name}/event v2 getevents
	//
	// Get detailed information about the specified event_ids
	//
	// get events
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
	// description: unified return format
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	j, err := simplejson.NewJson(b)
	if err != nil {
		logrus.Errorf("error decode json,details %s", err.Error())
		httputil.ReturnError(r, w, 400, "bad request")
		return
	}
	eventIDs, err := j.Get("event_ids").StringArray()
	if err != nil {
		logrus.Errorf("error get event_id in json,details %s", err.Error())
		httputil.ReturnError(r, w, 400, "bad request")
		return
	}

	events, err := handler.GetServiceEventHandler (). ListByEventIDs (eventIDs)
	if err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}

	httputil.ReturnSuccess(r, w, events)
}

//GetNotificationEvents GetNotificationEvent
//support query from start and end time or all
// swagger:operation GET  /v2/notificationEvent v2/notificationEvent getevents
//
// Get data center notification events
//
// get events
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
// description: unified return format
func GetNotificationEvents(w http.ResponseWriter, r *http.Request) {
	var startTime, endTime time.Time
	start := r.FormValue("start")
	end := r.FormValue("end")
	if si, err := strconv.Atoi(start); err == nil {
		startTime = time.Unix(int64(si), 0)
	}
	if ei, err := strconv.Atoi(end); err == nil {
		endTime = time.Unix(int64(ei), 0)
	}
	res, err := db.GetManager().NotificationEventDao().GetNotificationEventByTime(startTime, endTime)
	if err != nil {
		logrus.Errorf(err.Error())
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	for _, v := range res {
		service, err := db.GetManager().TenantServiceDao().GetServiceByID(v.KindID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				v.ServiceName = ""
				v.TenantName = ""
				continue
			} else {
				logrus.Errorf(err.Error())
				httputil.ReturnError(r, w, 500, err.Error())
				return
			}
		}
		tenant, err := db.GetManager().TenantDao().GetTenantByUUID(service.TenantID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				v.ServiceName = ""
				v.TenantName = ""
				continue
			} else {
				logrus.Errorf(err.Error())
				httputil.ReturnError(r, w, 500, err.Error())
				return
			}
		}
		v.ServiceName = service.ServiceAlias
		v.TenantName = tenant.Name
	}
	httputil.ReturnSuccess(r, w, res)
}

//Handle Handle
// swagger:parameters handlenotify
type Handle struct {
	Body struct {
		//in: body
		//handle message
		HandleMessage string `json:"handle_message" validate:"handle_message"`
	}
}

//HandleNotificationEvent HandleNotificationEvent
// swagger:operation PUT  /v2/notificationEvent/{hash} v2/notificationEvent handlenotify
//
// handle notification events
//
// get events
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
// description: unified return format
func HandleNotificationEvent(w http.ResponseWriter, r *http.Request) {
	serviceAlias := chi.URLParam(r, "serviceAlias")
	if serviceAlias == "" {
		httputil.ReturnError(r, w, 400, "ServiceAlias id do not empty")
		return
	}
	var handle Handle
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &handle.Body, nil)
	if !ok {
		return
	}
	service, err := db.GetManager().TenantServiceDao().GetServiceByServiceAlias(serviceAlias)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 404, "not found")
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	eventList, err := db.GetManager().NotificationEventDao().GetNotificationEventByKind("service", service.ServiceID)
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	for _, event := range eventList {
		event.IsHandle = true
		event.HandleMessage = handle.Body.HandleMessage
		err = db.GetManager().NotificationEventDao().UpdateModel(event)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				httputil.ReturnError(r, w, 404, "not found")
				return
			}
			httputil.ReturnError(r, w, 500, err.Error())
			return
		}
	}

	httputil.ReturnSuccess(r, w, nil)
}

//GetNotificationEvent GetNotificationEvent
// swagger:operation GET  /v2/notificationEvent/{hash} v2/notificationEvent getevents
//
// Get notification events
//
// get events
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
// description: unified return format
func GetNotificationEvent(w http.ResponseWriter, r *http.Request) {

	hash := chi.URLParam(r, "hash")
	if hash == "" {
		httputil.ReturnError(r, w, 400, "hash id do not empty")
		return
	}
	event, err := db.GetManager().NotificationEventDao().GetNotificationEventByHash(hash)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 404, "not found")
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, event)
}
