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

package model

// TableName - returns table name of Certificate
func (Certificate) TableName() string {
	return "gateway_certificate"
}

// Certificate contains TLS information
type Certificate struct {
	Model
	UUID            string `gorm:"column:uuid"`
	CertificateName string `gorm:"column:certificate_name;size:128"`
	Certificate     string `gorm:"column:certificate;size:65535"`
	PrivateKey      string `gorm:"column:private_key;size:65535"`
}

// TableName returns table name of RuleExtension
func (RuleExtension) TableName() string {
	return "gateway_rule_extension"
}

// RuleExtensionKey rule extension key
type RuleExtensionKey string

// HTTPToHTTPS forces http rewrite to https
var HTTPToHTTPS RuleExtensionKey = "httptohttps"

// LBType load balancer type
var LBType RuleExtensionKey = "lb-type"

// RuleExtension contains rule extensions for http rule or tcp rule
type RuleExtension struct {
	Model
	UUID   string `gorm:"column:uuid"`
	RuleID string `gorm:"column:rule_id"`
	Key    string `gorm:"column:key"`
	Value  string `gorm:"column:value"`
}

// LoadBalancerType load balancer type
type LoadBalancerType string

// RoundRobin round-robin load balancer type
var RoundRobin LoadBalancerType = "RoundRobin"

// ConsistenceHash consistence hash load balancer type
var ConsistenceHash LoadBalancerType = "ConsistentHash"

// TableName returns table name of HTTPRule
func (HTTPRule) TableName() string {
	return "gateway_http_rule"
}

// HTTPRule contains http rule
type HTTPRule struct {
	Model
	UUID          string `gorm:"column:uuid"`
	ServiceID     string `gorm:"column:service_id"`
	ContainerPort int    `gorm:"column:container_port"`
	Domain        string `gorm:"column:domain"`
	Path          string `gorm:"column:path"`
	Header        string `gorm:"column:header"`
	Cookie        string `gorm:"column:cookie"`
	Weight        int    `gorm:"column:weight"`
	IP            string `gorm:"column:ip"`
	CertificateID string `gorm:"column:certificate_id"`
}

// TableName returns table name of TCPRule
func (TCPRule) TableName() string {
	return "gateway_tcp_rule"
}

// TCPRule contain stream rule
type TCPRule struct {
	Model
	UUID          string `gorm:"column:uuid"`
	ServiceID     string `gorm:"column:service_id"`
	ContainerPort int    `gorm:"column:container_port"`
	// external access ip
	IP string `gorm:"column:ip"`
	// external access port
	Port int `gorm:"column:port"`
}

// GwRuleConfig describes a configuration of gateway rule.
type GwRuleConfig struct {
	Model
	RuleID string `gorm:"column:rule_id;size:32"`
	Key    string `gorm:"column:key"`
	Value  string `gorm:"column:value"`
}

// TableName -
func (GwRuleConfig) TableName() string {
	return "gateway_rule_config"
}
