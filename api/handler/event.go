// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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

package handler

import (
	"time"

	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/sirupsen/logrus"
)

// ServiceEventHandler -
type ServiceEventHandler struct {
}

// NewServiceEventHandler -
func NewServiceEventHandler() *ServiceEventHandler {
	return &ServiceEventHandler{}
}

// ListByEventIDs -
func (s *ServiceEventHandler) ListByEventIDs(eventIDs []string) ([]*dbmodel.ServiceEvent, error) {
	events, err := db.GetManager().ServiceEventDao().GetEventByEventIDs(eventIDs)
	if err != nil {
		return nil, err
	}

	// timeout events
	var timeoutEvents []*dbmodel.ServiceEvent
	for _, event := range events {
		if !s.isTimeout(event) {
			continue
		}
		event.Status = "timeout"
		event.FinalStatus = "complete"
		timeoutEvents = append(timeoutEvents, event)
	}

	return events, db.GetManager().ServiceEventDao().UpdateInBatch(timeoutEvents)
}

func (s *ServiceEventHandler) isTimeout(event *dbmodel.ServiceEvent) bool {
	if event.FinalStatus != "" {
		return false
	}

	startTime, err := time.ParseInLocation(time.RFC3339, event.StartTime, time.Local)
	if err != nil {
		logrus.Errorf("[ServiceEventHandler] [isTimeout] parse start time(%s): %v", event.StartTime, err)
		return false
	}

	if event.OptType == "deploy" || event.OptType == "create" || event.OptType == "build" || event.OptType == "upgrade" {
		end := startTime.Add(3 * time.Minute)
		if time.Now().After(end) {
			return true
		}
	} else {
		end := startTime.Add(30 * time.Second)
		if time.Now().After(end) {
			return true
		}
	}

	return false
}
