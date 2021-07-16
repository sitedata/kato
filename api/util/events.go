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

package util

import (
	"time"

	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// CanDoEvent check can do event or not
func CanDoEvent(optType string, synType int, target, targetID string) bool {
	if synType == dbmodel.SYNEVENTTYPE {
		return true
	}
	event, err := db.GetManager().ServiceEventDao().GetLastASyncEvent(target, targetID)
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return true
		}
		logrus.Error("get event by targetID error:", err)
		return false
	}
	if event == nil || event.FinalStatus != "" {
		return true
	}
	if !checkTimeout(event) {
		return false
	}
	return true
}

func checkTimeout(event *dbmodel.ServiceEvent) bool {
	if event.SynType == dbmodel.ASYNEVENTTYPE {
		if event.FinalStatus == "" {
			startTime := event.StartTime
			start, err := time.ParseInLocation(time.RFC3339, startTime, time.Local)
			if err != nil {
				return false
			}
			var end time.Time
			end = start.Add(3 * time.Minute)
			if time.Now().After(end) {
				event.FinalStatus = "timeout"
				err = db.GetManager().ServiceEventDao().UpdateModel(event)
				if err != nil {
					logrus.Error("check event timeout error : ", err.Error())
					return false
				}
				return true
			}
			// latest event is still processing on
			return false
		}
	}
	return true
}

// CreateEvent save event
func CreateEvent(target, optType, targetID, tenantID, reqBody, userName string, synType int) (*dbmodel.ServiceEvent, error) {
	if len(reqBody) > 1024 {
		reqBody = reqBody[0:1024]
	}
	event := dbmodel.ServiceEvent{
		EventID:     util.NewUUID(),
		TenantID:    tenantID,
		Target:      target,
		TargetID:    targetID,
		RequestBody: reqBody,
		UserName:    userName,
		StartTime:   time.Now().Format(time.RFC3339),
		SynType:     synType,
		OptType:     optType,
	}
	err := db.GetManager().ServiceEventDao().AddModel(&event)
	return &event, err
}

// UpdateEvent
func UpdateEvent(eventID string, statusCode int) {
	event, err := db.GetManager().ServiceEventDao().GetEventByEventID(eventID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Errorf("find event by eventID error : %s", err.Error())
		return
	}
	if err == gorm.ErrRecordNotFound {
		logrus.Errorf("do not found event by eventID %s", eventID)
		return
	}
	event.FinalStatus = "complete"
	event.EndTime = time.Now().Format(time.RFC3339)
	if statusCode < 400 { // status code 2XX/3XX all equal to success
		event.Status = "success"
	} else {
		event.Status = "failure"
	}
	err = db.GetManager().ServiceEventDao().UpdateModel(event)
	if err != nil {
		logrus.Errorf("update event status failure %s", err.Error())
		retry := 2
		for retry > 0 {
			if err = db.GetManager().ServiceEventDao().UpdateModel(event); err != nil {
				retry--
			} else {
				break
			}
		}
	}
}
