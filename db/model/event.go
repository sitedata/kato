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

import "time"

// ASYNEVENTTYPE
const ASYNEVENTTYPE = 0

// SYNEVENTTYPE
const SYNEVENTTYPE = 1

// TargetTypeService
const TargetTypeService = "service"

// TargetTypePod
const TargetTypePod = "pod"

// TargetTypeTenant
const TargetTypeTenant = "tenant"

// UsernameSystem
const UsernameSystem = "system"

// EventFinalStatus
type EventFinalStatus string

// String
func (e EventFinalStatus) String() string {
	return string(e)
}

// EventFinalStatusComplete
var EventFinalStatusComplete EventFinalStatus = "complete"

// EventFinalStatusFailure
var EventFinalStatusFailure EventFinalStatus = "failure"

// EventFinalStatusRunning
var EventFinalStatusRunning EventFinalStatus = "running"

// EventFinalStatusEmpty
var EventFinalStatusEmpty EventFinalStatus = "empty"

// EventFinalStatusEmptyComplete
var EventFinalStatusEmptyComplete EventFinalStatus = "emptycomplete"

// EventStatus
type EventStatus string

// String
func (e EventStatus) String() string {
	return string(e)
}

// EventStatusSuccess
var EventStatusSuccess EventStatus = "success"

// EventStatusFailure
var EventStatusFailure EventStatus = "failure"

//ServiceEvent event struct
type ServiceEvent struct {
	Model
	EventID     string `gorm:"column:event_id;size:40"`
	TenantID    string `gorm:"column:tenant_id;size:40;index:tenant_id"`
	ServiceID   string `gorm:"column:service_id;size:40;index:service_id"`
	Target      string `gorm:"column:target;size:40"`
	TargetID    string `gorm:"column:target_id;size:255"`
	RequestBody string `gorm:"column:request_body;size:1024"`
	UserName    string `gorm:"column:user_name;size:40"`
	StartTime   string `gorm:"column:start_time;size:40"`
	EndTime     string `gorm:"column:end_time;size:40"`
	OptType     string `gorm:"column:opt_type;size:40"`
	SynType     int    `gorm:"column:syn_type;size:1"`
	Status      string `gorm:"column:status;size:40"`
	FinalStatus string `gorm:"column:final_status;size:40"`
	Message     string `gorm:"column:message"`
	Reason      string `gorm:"column:reason"`
}

//TableName
func (t *ServiceEvent) TableName() string {
	return "tenant_services_event"
}

//NotificationEvent
type NotificationEvent struct {
	Model
	//Kind could be service, tenant, cluster, node
	Kind string `gorm:"column:kind;size:40"`
	//KindID could be service_id,tenant_id,cluster_id,node_id
	KindID string `gorm:"column:kind_id;size:40"`
	Hash   string `gorm:"column:hash;size:100"`
	//Type could be Normal UnNormal Notification
	Type          string    `gorm:"column:type;size:40"`
	Message       string    `gorm:"column:message;size:200"`
	Reason        string    `gorm:"column:reson;size:200"`
	Count         int       `gorm:"column:count;"`
	LastTime      time.Time `gorm:"column:last_time;"`
	FirstTime     time.Time `gorm:"column:first_time;"`
	IsHandle      bool      `gorm:"column:is_handle;"`
	HandleMessage string    `gorm:"column:handle_message;"`
	ServiceName   string    `gorm:"column:service_name;size:40"`
	TenantName    string    `gorm:"column:tenant_name;size:40"`
}

//TableName
func (n *NotificationEvent) TableName() string {
	return "region_notification_event"
}
