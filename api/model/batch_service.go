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

package model

import (
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	wmodel "github.com/gridworkz/kato/worker/discover/model"
)

var _ ComponentOpReq = &ComponentStartReq{}
var _ ComponentOpReq = &ComponentStopReq{}
var _ ComponentOpReq = &ComponentBuildReq{}
var _ ComponentOpReq = &ComponentUpgradeReq{}

// BatchOpRequesters -
type BatchOpRequesters []ComponentOpReq

// ComponentIDs returns a list of components ids.
func (b BatchOpRequesters) ComponentIDs() []string {
	var componentIDs []string
	for _, item := range b {
		componentIDs = append(componentIDs, item.GetComponentID())
	}
	return componentIDs
}

// ComponentOpReq -
type ComponentOpReq interface {
	GetComponentID() string
	GetEventID() string
	TaskBody(component *dbmodel.TenantServices) interface{}
	BatchOpFailureItem() *ComponentOpResult
	UpdateConfig(key, value string)
	OpType() string
	SetVersion(version string)
	GetVersion() string
}

// BatchOpResult -
type BatchOpResult []*ComponentOpResult

// BatchOpResultItemStatus is the status of ComponentOpResult.
type BatchOpResultItemStatus string

// BatchOpResultItemStatus -
var (
	BatchOpResultItemStatusFailure BatchOpResultItemStatus = "failure"
	BatchOpResultItemStatusSuccess BatchOpResultItemStatus = "success"
)

// ComponentOpResult -
type ComponentOpResult struct {
	ServiceID     string                  `json:"service_id"`
	Operation     string                  `json:"operation"`
	EventID       string                  `json:"event_id"`
	Status        BatchOpResultItemStatus `json:"status"`
	ErrMsg        string                  `json:"err_message"`
	DeployVersion string                  `json:"deploy_version"`
}

// Success sets the status to success.
func (b *ComponentOpResult) Success() {
	b.Status = BatchOpResultItemStatusSuccess
}

// ComponentOpGeneralReq -
type ComponentOpGeneralReq struct {
	EventID   string            `json:"event_id"`
	ServiceID string            `json:"service_id"`
	Configs   map[string]string `json:"configs"`
	// When determining the startup sequence of services, you need to know the services they depend on
	DepServiceIDInBootSeq []string `json:"dep_service_ids_in_boot_seq"`
}

// UpdateConfig -
func (b *ComponentOpGeneralReq) UpdateConfig(key, value string) {
	if b.Configs == nil {
		b.Configs = make(map[string]string)
	}
	b.Configs[key] = value
}

// ComponentStartReq -
type ComponentStartReq struct {
	ComponentOpGeneralReq
}

// GetEventID -
func (s *ComponentStartReq) GetEventID() string {
	if s.EventID == "" {
		s.EventID = util.NewUUID()
	}
	return s.EventID
}

// GetVersion -
func (s *ComponentStartReq) GetVersion() string {
	return ""
}

// SetVersion -
func (s *ComponentStartReq) SetVersion(string) {
	// no need
}

// GetComponentID -
func (s *ComponentStartReq) GetComponentID() string {
	return s.ServiceID
}

// TaskBody -
func (s *ComponentStartReq) TaskBody(cpt *dbmodel.TenantServices) interface{} {
	return &wmodel.StartTaskBody{
		TenantID:              cpt.TenantID,
		ServiceID:             cpt.ServiceID,
		DeployVersion:         cpt.DeployVersion,
		EventID:               s.GetEventID(),
		Configs:               s.Configs,
		DepServiceIDInBootSeq: s.DepServiceIDInBootSeq,
	}
}

// OpType -
func (s *ComponentStartReq) OpType() string {
	return "start-service"
}

// BatchOpFailureItem -
func (s *ComponentStartReq) BatchOpFailureItem() *ComponentOpResult {
	return &ComponentOpResult{
		ServiceID: s.ServiceID,
		EventID:   s.GetEventID(),
		Operation: "start",
		Status:    BatchOpResultItemStatusFailure,
	}
}

// ComponentStopReq -
type ComponentStopReq struct {
	ComponentStartReq
}

// OpType -
func (s *ComponentStopReq) OpType() string {
	return "stop-service"
}

// BatchOpFailureItem -
func (s *ComponentStopReq) BatchOpFailureItem() *ComponentOpResult {
	return &ComponentOpResult{
		ServiceID: s.ServiceID,
		EventID:   s.GetEventID(),
		Operation: "stop",
		Status:    BatchOpResultItemStatusFailure,
	}
}
