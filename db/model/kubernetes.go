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

//TypeStatefulSet typestateful
var TypeStatefulSet = "statefulset"

//TypeDeployment typedeployment
var TypeDeployment = "deployment"

//TypeReplicationController type rc
var TypeReplicationController = "replicationcontroller"

//LocalScheduler
type LocalScheduler struct {
	Model
	ServiceID string `gorm:"column:service_id;size:32"`
	NodeIP    string `gorm:"column:node_ip;size:32"`
	PodName   string `gorm:"column:pod_name;size:32"`
}

//TableName
func (t *LocalScheduler) TableName() string {
	return "local_scheduler"
}

//ServiceSourceConfig service source config info
//such as deployment、statefulset、configmap
type ServiceSourceConfig struct {
	Model
	ServiceID  string `gorm:"column:service_id;size:32"`
	SourceType string `gorm:"column:source_type;size:32"`
	SourceBody string `gorm:"column:source_body;size:2000"`
}

//TableName
func (t *ServiceSourceConfig) TableName() string {
	return "tenant_services_source"
}
