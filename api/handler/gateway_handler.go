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
	apimodel "github.com/gridworkz/kato/api/model"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/jinzhu/gorm"
)

// ComponentIngressTask -
type ComponentIngressTask struct {
	ComponentID string `json:"service_id"`
	Action      string `json:"action"`
	Port        int    `json:"port"`
	IsInner     bool   `json:"is_inner"`
}

//GatewayHandler gateway api handler
type GatewayHandler interface {
	AddHTTPRule(req *apimodel.AddHTTPRuleStruct) error
	CreateHTTPRule(tx *gorm.DB, req *apimodel.AddHTTPRuleStruct) error
	UpdateHTTPRule(req *apimodel.UpdateHTTPRuleStruct) error
	DeleteHTTPRule(req *apimodel.DeleteHTTPRuleStruct) error
	DeleteHTTPRuleByServiceIDWithTransaction(sid string, tx *gorm.DB) error

	AddCertificate(req *apimodel.AddHTTPRuleStruct, tx *gorm.DB) error
	UpdateCertificate(req apimodel.AddHTTPRuleStruct, httpRule *dbmodel.HTTPRule, tx *gorm.DB) error

	AddTCPRule(req *apimodel.AddTCPRuleStruct) error
	CreateTCPRule(tx *gorm.DB, req *apimodel.AddTCPRuleStruct) error
	UpdateTCPRule(req *apimodel.UpdateTCPRuleStruct, minPort int) error
	DeleteTCPRule(req *apimodel.DeleteTCPRuleStruct) error
	DeleteTCPRuleByServiceIDWithTransaction(sid string, tx *gorm.DB) error
	AddRuleExtensions(ruleID string, ruleExtensions []*apimodel.RuleExtensionStruct, tx *gorm.DB) error
	GetAvailablePort(ip string, lock bool) (int, error)
	TCPIPPortExists(ip string, port int) bool
	// Deprecated.
	SendTaskDeprecated(in map[string]interface{}) error
	SendTask(task *ComponentIngressTask) error
	RuleConfig(req *apimodel.RuleConfigReq) error
	UpdCertificate(req *apimodel.UpdCertificateReq) error
	GetGatewayIPs() []IPAndAvailablePort
	ListHTTPRulesByCertID(certID string) ([]*dbmodel.HTTPRule, error)
	DeleteIngressRulesByComponentPort(tx *gorm.DB, componentID string, port int) error
	SyncHTTPRules(tx *gorm.DB, components []*apimodel.Component) error
	SyncTCPRules(tx *gorm.DB, components []*apimodel.Component) error
}
