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

package region

import (
	"bytes"
	"fmt"
	"github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	dbmodel "github.com/gridworkz/kato/db/model"
	coreutil "github.com/gridworkz/kato/util"
	utilhttp "github.com/gridworkz/kato/util/http"
)

type services struct {
	tenant
	prefix string
	model  model.ServiceStruct
}

//ServiceInterface ServiceInterface
type ServiceInterface interface {
	Get() (*serviceInfo, *util.APIHandleError)
	GetDeployInfo() (*ServiceDeployInfo, *util.APIHandleError)
	Pods() ([]*podInfo, *util.APIHandleError)
	List() ([]*dbmodel.TenantServices, *util.APIHandleError)
	Stop(eventID string) (string, *util.APIHandleError)
	Start(eventID string) (string, *util.APIHandleError)
	EventLog(eventID, level string) ([]*model.MessageData, *util.APIHandleError)
}

func (s *services) Pods() ([]*podInfo, *util.APIHandleError) {
	var gc []*podInfo
	var decode utilhttp.ResponseBody
	decode.List = &gc
	code, err := s.DoRequest(s.prefix+"/pods", "GET", nil, &decode)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get pods code %d", code))
	}
	return gc, nil
}
func (s *services) Get() (*serviceInfo, *util.APIHandleError) {
	var service serviceInfo
	var decode utilhttp.ResponseBody
	decode.Bean = &service
	code, err := s.DoRequest(s.prefix, "GET", nil, &decode)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get err with code %d", code))
	}
	return &service, nil
}
func (s *services) EventLog(eventID, level string) ([]*model.MessageData, *util.APIHandleError) {
	data := []byte(`{"event_id":"` + eventID + `","level":"` + level + `"}`)
	var message []*model.MessageData
	var decode utilhttp.ResponseBody
	decode.List = &message
	code, err := s.DoRequest(s.prefix+"/event-log", "POST", bytes.NewBuffer(data), &decode)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get event log code %d", code))
	}
	return message, nil
}

func (s *services) List() ([]*dbmodel.TenantServices, *util.APIHandleError) {
	var gc []*dbmodel.TenantServices
	var decode utilhttp.ResponseBody
	decode.List = &gc
	code, err := s.DoRequest(s.prefix, "GET", nil, &decode)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get with code %d", code))
	}
	return gc, nil
}
func (s *services) Stop(eventID string) (string, *util.APIHandleError) {
	if eventID == "" {
		eventID = coreutil.NewUUID()
	}
	data := []byte(`{"event_id":"` + eventID + `"}`)
	var res utilhttp.ResponseBody
	code, err := s.DoRequest(s.prefix+"/stop", "POST", bytes.NewBuffer(data), nil)
	if err != nil {
		return "", handleErrAndCode(err, code)
	}
	return eventID, handleAPIResult(code, res)
}
func (s *services) Start(eventID string) (string, *util.APIHandleError) {
	if eventID == "" {
		eventID = coreutil.NewUUID()
	}
	var res utilhttp.ResponseBody
	data := []byte(`{"event_id":"` + eventID + `"}`)
	code, err := s.DoRequest(s.prefix+"/start", "POST", bytes.NewBuffer(data), nil)
	if err != nil {
		return "", handleErrAndCode(err, code)
	}
	return eventID, handleAPIResult(code, res)
}

//GetDeployInfo get service deploy info
func (s *services) GetDeployInfo() (*ServiceDeployInfo, *util.APIHandleError) {
	var deployInfo ServiceDeployInfo
	var decode utilhttp.ResponseBody
	decode.Bean = &deployInfo
	code, err := s.DoRequest(s.prefix+"/deploy-info", "GET", nil, &decode)
	return &deployInfo, handleErrAndCode(err, code)
}
