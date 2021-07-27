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

import dbmodel "github.com/gridworkz/kato/db/model"

// AutoscalerRuleReq -
type AutoscalerRuleReq struct {
	RuleID      string `json:"rule_id" validate:"rule_id|required"`
	ServiceID   string
	Enable      bool   `json:"enable" validate:"enable|required"`
	XPAType     string `json:"xpa_type" validate:"xpa_type|required"`
	MinReplicas int    `json:"min_replicas" validate:"min_replicas|required"`
	MaxReplicas int    `json:"max_replicas" validate:"min_replicas|required"`
	Metrics     []struct {
		MetricsType       string `json:"metric_type"`
		MetricsName       string `json:"metric_name"`
		MetricTargetType  string `json:"metric_target_type"`
		MetricTargetValue int    `json:"metric_target_value"`
	} `json:"metrics"`
}

// AutoscalerRuleResp -
type AutoscalerRuleResp struct {
	RuleID      string `json:"rule_id"`
	ServiceID   string `json:"service_id"`
	Enable      bool   `json:"enable"`
	XPAType     string `json:"xpa_type"`
	MinReplicas int    `json:"min_replicas"`
	MaxReplicas int    `json:"max_replicas"`
	Metrics     []struct {
		MetricsType       string `json:"metric_type"`
		MetricsName       string `json:"metric_name"`
		MetricTargetType  string `json:"metric_target_type"`
		MetricTargetValue int    `json:"metric_target_value"`
	} `json:"metrics"`
}

// AutoScalerRule -
type AutoScalerRule struct {
	RuleID      string       `json:"rule_id"`
	Enable      bool         `json:"enable"`
	XPAType     string       `json:"xpa_type"`
	MinReplicas int          `json:"min_replicas"`
	MaxReplicas int          `json:"max_replicas"`
	RuleMetrics []RuleMetric `json:"metrics"`
}

// DbModel return database model
func (a AutoScalerRule) DbModel(componentID string) *dbmodel.TenantServiceAutoscalerRules {
	return &dbmodel.TenantServiceAutoscalerRules{
		RuleID:      a.RuleID,
		ServiceID:   componentID,
		MinReplicas: a.MinReplicas,
		MaxReplicas: a.MaxReplicas,
		Enable:      a.Enable,
		XPAType:     a.XPAType,
	}
}

// RuleMetric -
type RuleMetric struct {
	MetricsType       string `json:"metric_type"`
	MetricsName       string `json:"metric_name"`
	MetricTargetType  string `json:"metric_target_type"`
	MetricTargetValue int    `json:"metric_target_value"`
}

// DbModel return database model
func (r RuleMetric) DbModel(ruleID string) *dbmodel.TenantServiceAutoscalerRuleMetrics {
	return &dbmodel.TenantServiceAutoscalerRuleMetrics{
		RuleID:            ruleID,
		MetricsType:       r.MetricsType,
		MetricsName:       r.MetricsName,
		MetricTargetType:  r.MetricTargetType,
		MetricTargetValue: r.MetricTargetValue,
	}
}
