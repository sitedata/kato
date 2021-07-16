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

//SetDefineSourcesStruct
//swagger:parameters setDefineSource
type SetDefineSourcesStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	SourceAlias string `json:"source_alias" validate:"source_alias"`
	// in: body
	Body struct {
		//in: body
		//required: true
		SourceSpec *SourceSpec `json:"source_spec" validate:"source_spec"`
	}
}

//DeleteDefineSourcesStruct DeleteDefineSourcesStruct
//swagger:parameters deleteDefineSource getDefineSource
type DeleteDefineSourcesStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	SourceAlias string `json:"source_alias" validate:"source_alias"`
	// in: path
	// required: true
	EnvName string `json:"env_name" validate:"env_name"`
}

//UpdateDefineSourcesStruct UpdateDefineSourcesStruct
//swagger:parameters deleteDefineSource updateDefineSource
type UpdateDefineSourcesStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	SourceAlias string `json:"source_alias" validate:"source_alias"`
	// in: path
	// required: true
	EnvName string `json:"env_name" validate:"env_name"`
	// in: body
	Body struct {
		//in: body
		//required: true
		SourceSpec *SourceSpec `json:"source_spec" validate:"source_spec"`
	}
}

//SourceSpec SourceSpec
type SourceSpec struct {
	Alias      string               `json:"source_alias" validate:"source_alias"`
	Info       string               `json:"source_info" validate:"source_info"`
	CreateTime string               `json:"create_time" validate:"create_time"`
	Operator   string               `json:"operator" validate:"operator"`
	SourceBody *SoureBody           `json:"source_body" validate:"source_body"`
	Additions  map[string]*Addition `json:"additons" validate:"additions"`
}

//SoureBody SoureBody
type SoureBody struct {
	EnvName string      `json:"env_name" validate:"env_name"`
	EnvVal  interface{} `json:"env_value" validate:"env_value"`
	//json format
}

//ResourceSpec - resource structure
type ResourceSpec struct {
	BasePorts    []*BasePort    `json:"base_ports"`
	BaseServices []*BaseService `json:"base_services"`
	BaseNormal   BaseEnv        `json:"base_normal"`
}

//BasePort base of current app ports
type BasePort struct {
	ServiceAlias string `json:"service_alias"`
	ServiceID    string `json:"service_id"`
	//Port is the real app port
	Port int `json:"port"`
	//ListenPort is mesh listen port, proxy connetion to real app port
	ListenPort int                    `json:"listen_port"`
	Protocol   string                 `json:"protocol"`
	Options    map[string]interface{} `json:"options"`
}

//BaseService - based on dependent application and port structure
type BaseService struct {
	ServiceAlias       string                 `json:"service_alias"`
	ServiceID          string                 `json:"service_id"`
	DependServiceAlias string                 `json:"depend_service_alias"`
	DependServiceID    string                 `json:"depend_service_id"`
	Port               int                    `json:"port"`
	Protocol           string                 `json:"protocol"`
	Options            map[string]interface{} `json:"options"`
}

//BaseEnv - no platform-defined type, ordinary kv
type BaseEnv struct {
	Options map[string]interface{} `json:"options"`
}

//Item source value, key-value pair form
type Item struct {
	Key   string      `json:"key" validate:"key"`
	Value interface{} `json:"value" validate:"value"`
}

//Addition - store additional information
type Addition struct {
	Desc  string  `json:"desc" validate:"desc"`
	Items []*Item `json:"items" validate:"items"`
}
