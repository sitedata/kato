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
	"time"
)

// Plugin -
type Plugin struct {
	PluginID    string `json:"plugin_id" validate:"plugin_id|required"`
	PluginName  string `json:"plugin_name" validate:"plugin_name|required"`
	PluginInfo  string `json:"plugin_info" validate:"plugin_info"`
	ImageURL    string `json:"image_url" validate:"image_url"`
	GitURL      string `json:"git_url" validate:"git_url"`
	BuildModel  string `json:"build_model" validate:"build_model"`
	PluginModel string `json:"plugin_model" validate:"plugin_model"`
	TenantID    string `json:"tenant_id" validate:"tenant_id"`
}

// DbModel return database model
func (p *Plugin) DbModel(tenantID string) *dbmodel.TenantPlugin {
	return &dbmodel.TenantPlugin{
		PluginID:    p.PluginID,
		PluginName: p.PluginName,
		PluginInfo:  p.PluginInfo,
		ImageURL:    p.ImageURL,
		GitURL:      p.GitURL,
		BuildModel:  p.BuildModel,
		PluginModel: p.PluginModel,
		TenantID:    tenantID,
	}
}

// BatchCreatePlugins -
type BatchCreatePlugins struct {
	Plugins []*Plugin `json:"plugins"`
}

//CreatePluginStruct CreatePluginStruct
//swagger:parameters createPlugin
type CreatePluginStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: body
	Body struct {
		//Plugin id
		//in: body
		//required: true
		PluginID string `json:"plugin_id" validate:"plugin_id|required"`
		//in: body
		//required: true
		PluginName string `json:"plugin_name" validate:"plugin_name|required"`
		//Plugin usage description
		//in: body
		//required: false
		PluginInfo string `json:"plugin_info" validate:"plugin_info"`
		// Plug-in docker address
		// in:body
		// required: false
		ImageURL string `json:"image_url" validate:"image_url"`
		//git address
		//in: body
		//required: false
		GitURL string `json:"git_url" validate:"git_url"`
		//Build mode
		//in: body
		//required: false
		BuildModel string `json:"build_model" validate:"build_model"`
		//Plugin mode
		//in: body
		//required: false
		PluginModel string `json:"plugin_model" validate:"plugin_model"`
		//Tenant id
		//in: body
		//required: false
		TenantID string `json:"tenant_id" validate:"tenant_id"`
	}
}

//UpdatePluginStruct UpdatePluginStruct
//swagger:parameters updatePlugin
type UpdatePluginStruct struct {
	// Tenant name
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name|required"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id" validate:"tenant_name|required"`
	// in: body
	Body struct {
		//Plugin name
		//in: body
		//required: false
		PluginName string `json:"plugin_name" validate:"plugin_name"`
		//Plugin usage description
		//in: body
		//required: false
		PluginInfo string `json:"plugin_info" validate:"plugin_info"`
		//Plugin docker address
		//in: body
		//required: false
		ImageURL string `json:"image_url" validate:"image_url"`
		//git address
		//in: body
		//required: false
		GitURL string `json:"git_url" validate:"git_url"`
		//Build mode
		//in: body
		//required: false
		BuildModel string `json:"build_model" validate:"build_model"`
		//Plugin mode
		//in: body
		//required: false
		PluginModel string `json:"plugin_model" validate:"plugin_model"`
	}
}

//DeletePluginStruct deletePluginStruct
//swagger:parameters deletePlugin
type DeletePluginStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name|required"`
	// in: path
	// required: true
	PluginID string `json:"plugin_id" validate:"plugin_id|required"`
}

//ENVStruct ENVStruct
//swagger:parameters adddefaultenv updatedefaultenv
type ENVStruct struct {
	// Tenant name
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id" validate:"plugin_id"`
	// build version
	// in: path
	// required; true
	VersionID string `json:"version_id" validate:"version_id"`
	//in : body
	Body struct {
		//in: body
		//required: true
		EVNInfo []*PluginDefaultENV
	}
}

//DeleteENVstruct DeleteENVstruct
//swagger:parameters deletedefaultenv
type DeleteENVstruct struct {
	// Tenant name
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name|required"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id" validate:"plugin_id|required"`
	// build version
	// in: path
	// required; true
	VersionID string `json:"version_id" validate:"version_id|required"`
	//Configuration item name
	//in: path
	//required: true
	ENVName string `json:"env_name" validate:"env_name|required"`
}

//PluginDefaultENV plugin default environment variable
type PluginDefaultENV struct {
	//Corresponding plug-in id
	//in: body
	//required: true
	PluginID string `json:"plugin_id" validate:"plugin_id"`
	//Build version id
	//in: body
	//required: true
	VersionID string `json:"version_id" validate:"version_id"`
	//Configuration item name
	//in: body
	//required: true
	ENVName string `json:"env_name" validate:"env_name"`
	//Configuration item value
	//in: body
	//required: true
	ENVValue string `json:"env_value" validate:"env_value"`
	//Can be modified by the user
	//in :body
	//required: false
	IsChange bool `json:"is_change" validate:"is_change|bool"`
}

//BuildPluginStruct BuildPluginStruct
//swagger:parameters buildPlugin
type BuildPluginStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	PluginID string `json:"plugin_id" validate:"plugin_id"`
	//in: body
	Body struct {
		// the event id
		// in: body
		// required: false
		EventID string `json:"event_id" validate:"event_id"`
		// Plug-in CPU weight, default 125
		// in: body
		// required: true
		PluginCPU int `json:"plugin_cpu" validate:"plugin_cpu|required"`
		// plug-in maximum memory, default 50
		// in: body
		// required: true
		PluginMemory int `json:"plugin_memory" validate:"plugin_memory|required"`
		// plugin cmd, default 50
		// in: body
		// required: false
		PluginCMD string `json:"plugin_cmd" validate:"plugin_cmd"`
		// The version number of the plugin
		// in: body
		// required: true
		BuildVersion string `json:"build_version" validate:"build_version|required"`
		// Plug-in build version number
		// in: body
		// required: true
		DeployVersion string `json:"deploy_version" validate:"deploy_version"`
		// git address branch information, the default is master
		// in: body
		// required: false
		RepoURL string `json:"repo_url" validate:"repo_url"`
		// git username
		// in: body
		// required: false
		Username string `json:"username"`
		// git password
		// in: body
		// required: false
		Password string `json:"password"`
		// Version information, assist in choosing the plug-in version
		// in:body
		// required: true
		Info string `json:"info" validate:"info"`
		// Operator
		// in: body
		// required: false
		Operator string `json:"operator" validate:"operator"`
		//Tenant id
		// in: body
		// required: true
		TenantID string `json:"tenant_id" validate:"tenant_id"`
		// Mirror address
		// in: body
		// required: false
		BuildImage string `json:"build_image" validate:"build_image"`
		//ImageInfo
		ImageInfo struct {
			HubURL      string `json:"hub_url"`
			HubUser     string `json:"hub_user"`
			HubPassword string `json:"hub_password"`
			Namespace   string `json:"namespace"`
			IsTrust     bool   `json:"is_trust,omitempty"`
		} `json:"ImageInfo" validate:"ImageInfo"`
	}
}

// BuildPluginReq -
type BuildPluginReq struct {
	PluginID      string `json:"plugin_id" validate:"plugin_id"`
	EventID       string `json:"event_id" validate:"event_id"`
	PluginCPU     int    `json:"plugin_cpu" validate:"plugin_cpu|required"`
	PluginMemory  int    `json:"plugin_memory" validate:"plugin_memory|required"`
	PluginCMD     string `json:"plugin_cmd" validate:"plugin_cmd"`
	BuildVersion  string `json:"build_version" validate:"build_version|required"`
	DeployVersion string `json:"deploy_version" validate:"deploy_version"`
	RepoURL       string `json:"repo_url" validate:"repo_url"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Info          string `json:"info" validate:"info"`
	Operator      string `json:"operator" validate:"operator"`
	TenantID      string `json:"tenant_id" validate:"tenant_id"`
	BuildImage    string `json:"build_image" validate:"build_image"`
	ImageInfo     struct {
		HubURL      string `json:"hub_url"`
		HubUser     string `json:"hub_user"`
		HubPassword string `json:"hub_password"`
		Namespace   string `json:"namespace"`
		IsTrust     bool   `json:"is_trust,omitempty"`
	} `json:"ImageInfo" validate:"ImageInfo"`
}

// DbModel return database model
func (b BuildPluginReq) DbModel(plugin *dbmodel.TenantPlugin) *dbmodel.TenantPluginBuildVersion {
	buildVersion := &dbmodel.TenantPluginBuildVersion{
		VersionID:       b.BuildVersion,
		DeployVersion:   b.DeployVersion,
		PluginID:        b.PluginID,
		Kind:            plugin.BuildModel,
		Repo:            b.RepoURL,
		GitURL:          plugin.GitURL,
		BaseImage:       plugin.ImageURL,
		ContainerCPU:    b.PluginCPU,
		ContainerMemory: b.PluginMemory,
		ContainerCMD:    b.PluginCMD,
		BuildTime:       time.Now().Format(time.RFC3339),
		Info:            b.Info,
		Status:          "building",
	}
	if b.PluginCPU == 0 {
		buildVersion.ContainerCPU = 125
	}
	if b.PluginMemory == 0 {
		buildVersion.ContainerMemory = 50
	}
	return buildVersion
}

// BatchBuildPlugins -
type BatchBuildPlugins struct {
	Plugins []*BuildPluginReq `json:"plugins"`
}

//PluginBuildVersionStruct PluginBuildVersionStruct
//swagger:parameters deletePluginVersion pluginVersion
type PluginBuildVersionStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	//in: path
	//required: true
	PluginID string `json:"plugin_id" validate:"plugin_id"`
	//in: path
	//required: true
	VersionID string `json:"version_id" validate:"version_id"`
}

//AllPluginBuildVersionStruct AllPluginBuildVersionStruct
//swagger:parameters allPluginVersions
type AllPluginBuildVersionStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	//in: path
	//required: true
	PluginID string `json:"plugin_id" validate:"plugin_id"`
}

//PluginSetStruct PluginSetStruct
//swagger:parameters updatePluginSet addPluginSet
type PluginSetStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
	//in: path
	//required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		//plugin id
		//in: body
		//required: true
		PluginID string `json:"plugin_id" validate:"plugin_id"`
		// plugin version
		//in: body
		//required: true
		VersionID string `json:"version_id" validate:"version_id"`
		// plugin is uesd
		//in: body
		//required: false
		Switch bool `json:"switch" validate:"switch|bool"`
		// plugin cpu size default 125
		// in: body
		// required: false
		PluginCPU int `json:"plugin_cpu" validate:"plugin_cpu"`
		// plugin memory size default 64
		// in: body
		// required: false
		PluginMemory int `json:"plugin_memory" validate:"plugin_memory"`
		// app plugin config
		// in: body
		// required: true
		ConfigEnvs ConfigEnvs `json:"config_envs" validate:"config_envs"`
	}
}

//GetPluginsStruct GetPluginsStruct
//swagger:parameters getPlugins
type GetPluginsStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
}

//GetPluginSetStruct GetPluginSetStruct
//swagger:parameters getPluginSet
type GetPluginSetStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
	//in: path
	//required: true
	ServiceAlias string `json:"service_alias"`
}

//DeletePluginSetStruct DeletePluginSetStruct
//swagger:parameters deletePluginRelation
type DeletePluginSetStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
	//in: path
	//required: true
	ServiceAlias string `json:"service_alias"`
	//Plugin id
	//in: path
	//required: true
	PluginID string `json:"plugin_id"`
}

//GetPluginEnvStruct GetPluginEnvStruct
//swagger:parameters getPluginEnv getPluginDefaultEnv
type GetPluginEnvStruct struct {
	//Tenant name
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id"`
	// Build version id
	// in: path
	// required: true
	VersionID string `json:"version_id"`
}

//GetVersionEnvStruct GetVersionEnvStruct
//swagger:parameters getVersionEnvs
type GetVersionEnvStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id"`
}

//SetVersionEnv SetVersionEnv
//swagger:parameters setVersionEnv updateVersionEnv
type SetVersionEnv struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id"`
	//in: body
	Body struct {
		TenantID  string `json:"tenant_id"`
		ServiceID string `json:"service_id"`
		// environment variables
		// in: body
		// required: true
		ConfigEnvs ConfigEnvs `json:"config_envs" validate:"config_envs"`
	}
}

//ConfigEnvs Config
type ConfigEnvs struct {
	NormalEnvs  []*VersionEnv `json:"normal_envs" validate:"normal_envs"`
	ComplexEnvs *ResourceSpec `json:"complex_envs" validate:"complex_envs"`
}

// VersionEnv VersionEnv
type VersionEnv struct {
	//variable name
	//in:body
	//required: true
	EnvName string `json:"env_name" validate:"env_name"`
	//variable
	//in:body
	//required: true
	EnvValue string `json:"env_value" validate:"env_value"`
}

// DbModel return database model
func (v *VersionEnv) DbModel(componentID, pluginID string) *dbmodel.TenantPluginVersionEnv {
	return &dbmodel.TenantPluginVersionEnv{
		ServiceID: componentID,
		PluginID:  pluginID,
		EnvName: v.EnvName,
		EnvValue: v.EnvValue,
	}
}

// TransPlugins TransPlugins
type TransPlugins struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	//in: body
	Body struct {
		// install from this tenant
		// in: body
		// required: true
		FromTenantName string `json:"from_tenant_name" validate:"from_tenant_name"`
		// plugin id
		// in: body
		// required: true
		PluginsID []string `json:"plugins_id" validate:"plugins_id"`
	}
}

// PluginVersionEnv -
type PluginVersionEnv struct {
	EnvName  string `json:"env_name" validate:"env_name"`
	EnvValue string `json:"env_value" validate:"env_value"`
}

// DbModel return database model
func (p *PluginVersionEnv) DbModel(componentID, pluginID string) *dbmodel.TenantPluginVersionEnv {
	return &dbmodel.TenantPluginVersionEnv{
		ServiceID: componentID,
		PluginID:  pluginID,
		EnvName:   p.EnvName,
		EnvValue:  p.EnvValue,
	}
}

// TenantPluginVersionConfig -
type TenantPluginVersionConfig struct {
	ConfigStr string `json:"config_str" validate:"config_str"`
}

// DbModel return database model
func (p *TenantPluginVersionConfig) DbModel(componentID, pluginID string) *dbmodel.TenantPluginVersionDiscoverConfig {
	return &dbmodel.TenantPluginVersionDiscoverConfig{
		ServiceID: componentID,
		PluginID:  pluginID,
		ConfigStr: p.ConfigStr,
	}
}

// ComponentPlugin -
type ComponentPlugin struct {
	PluginID        string     `json:"plugin_id"`
	VersionID       string     `json:"version_id"`
	PluginModel     string     `json:"plugin_model"`
	ContainerCPU    int        `json:"container_cpu"`
	ContainerMemory int        `json:"container_memory"`
	Switch          bool       `json:"switch"`
	ConfigEnvs      ConfigEnvs `json:"config_envs" validate:"config_envs"`
}

// DbModel return database model
func (p *ComponentPlugin) DbModel(componentID string) *dbmodel.TenantServicePluginRelation {
	return &dbmodel.TenantServicePluginRelation{
		VersionID:       p.VersionID,
		ServiceID:       componentID,
		PluginID:        p.PluginID,
		Switch:          p.Switch,
		PluginModel:     p.PluginModel,
		ContainerCPU:    p.ContainerCPU,
		ContainerMemory: p.ContainerMemory,
	}
}
