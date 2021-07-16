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

//SetNetDownStreamRuleStruct
//swagger:parameters setNetDownStreamRuleStruct
type SetNetDownStreamRuleStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		//in: body
		//required: true
		DestService string `json:"dest_service" validate:"dest_service"`
		//downstream service alias
		//in: body
		//required: true
		DestServiceAlias string `json:"dest_service_alias" validate:"dest_service_alias"`
		//port
		//in: body
		//required: true
		Port int `json:"port" validate:"port"`
		//protocol
		//in: body
		//required: true
		Protocol string `json:"protocol" validate:"protocol|between:tcp,http"`
		//regular body
		//in: body
		//required: true
		Rules *NetDownStreamRules `json:"rules" validate:"rules"`
	}
}

//NetRulesDownStreamBody
type NetRulesDownStreamBody struct {
	DestService      string              `json:"dest_service"`
	DestServiceAlias string              `json:"dest_service_alias"`
	Port             int                 `json:"port"`
	Protocol         string              `json:"protocol"`
	Rules            *NetDownStreamRules `json:"rules"`
}

//NetDownStreamRules
type NetDownStreamRules struct {
	//current limit - max_connections
	Limit              int `json:"limit" validate:"limit|numeric_between:0,1024"`
	MaxPendingRequests int `json:"max_pending_requests"`
	MaxRequests        int `json:"max_requests"`
	MaxRetries         int `json:"max_retries"`
	//request header
	//in: body
	//required: false
	Header []HeaderRules `json:"header" validate:"header"`
	//domain forwarding
	//in: body
	//required: false
	Domain []string `json:"domain" validate:"domain"`
	//path rule
	//in: body
	//required: false
	Prefix       string `json:"prefix" validate:"prefix"`
	ServiceAlias string `json:"service_alias"`
	ServiceID    string `json:"service_id" validate:"service_id"`
}

//NetUpStreamRules NetUpStreamRules
type NetUpStreamRules struct {
	NetDownStreamRules
	SourcePort int32 `json:"source_port"`
	MapPort    int32 `json:"map_port"`
}

//HeaderRules HeaderRules
type HeaderRules struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//GetNetDownStreamRuleStruct GetNetDownStreamRuleStruct
//swagger:parameters getNetDownStreamRuleStruct
type GetNetDownStreamRuleStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
	// in: path
	// required: true
	DestServiceAlias string `json:"dest_service_alias" validate:"dest_service_alias"`
	// in: path
	// required: true
	Port int `json:"port" validate:"port|numeric_between:1,65535"`
}

//UpdateNetDownStreamRuleStruct UpdateNetDownStreamRuleStruct
//swagger:parameters updateNetDownStreamRuleStruct
type UpdateNetDownStreamRuleStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
	// in: path
	// required: true
	DestServiceAlias string `json:"dest_service_alias" validate:"dest_service_alias"`
	// in: path
	// required: true
	Port int `json:"port" validate:"port|numeric_between:1,65535"`
	// in: body
	Body struct {
		//in: body
		//required: true
		DestService string `json:"dest_service" validate:"dest_service"`
		//protocol
		//in: body
		//required: true
		Protocol string `json:"protocol" validate:"protocol|between:tcp,http"`
		//regular body
		//in: body
		//required: true
		Rules *NetDownStreamRules `json:"rules" validate:"rules"`
	}
}
