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

//CreatePluginStruct
//swagger:parameters createPlugin
type CreatePluginStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: body
	Body struct {
		//plugin id
		//in: body
		//required: true
		PluginID string `json:"plugin_id" validate:"plugin_id|required"`
		//in: body
		//required: true
		PluginName string `json:"plugin_name" validate:"plugin_name|required"`
		//plugin description
		//in: body
		//required: false
		PluginInfo string `json:"plugin_info" validate:"plugin_info"`
		//plugin docker address
		//in:body
		//required: false
		ImageURL string `json:"image_url" validate:"image_url"`
		//git address
		//in: body
		//required: false
		GitURL string `json:"git_url" validate:"git_url"`
		//build mode
		//in: body
		//required: false
		BuildModel string `json:"build_model" validate:"build_model"`
		//plugin mode
		//in: body
		//required: false
		PluginModel string `json:"plugin_model" validate:"plugin_model"`
		//tenant id
		//in: body
		//required: false
		TenantID string `json:"tenant_id" validate:"tenant_id"`
	}
}

//UpdatePluginStruct
//swagger:parameters updatePlugin
type UpdatePluginStruct struct {
	// tenant name
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name|required"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id" validate:"tenant_name|required"`
	// in: body
	Body struct {
		//plugin name
		//in: body
		//required: false
		PluginName string `json:"plugin_name" validate:"plugin_name"`
		//plugin description
		//in: body
		//required: false
		PluginInfo string `json:"plugin_info" validate:"plugin_info"`
		//plugin docker address
		//in: body
		//required: false
		ImageURL string `json:"image_url" validate:"image_url"`
		//git address
		//in: body
		//required: false
		GitURL string `json:"git_url" validate:"git_url"`
		//build mode
		//in: body
		//required: false
		BuildModel string `json:"build_model" validate:"build_model"`
		//plugin mode
		//in: body
		//required: false
		PluginModel string `json:"plugin_model" validate:"plugin_model"`
	}
}

//DeletePluginStruct
//swagger:parameters deletePlugin
type DeletePluginStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name|required"`
	// in: path
	// required: true
	PluginID string `json:"plugin_id" validate:"plugin_id|required"`
}

//ENVStruct
//swagger:parameters adddefaultenv updatedefaultenv
type ENVStruct struct {
	// tenant name
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

//DeleteENVstruct
//swagger:parameters deletedefaultenv
type DeleteENVstruct struct {
	// tenant name
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
	//configuration item name
	//in: path
	//required: true
	ENVName string `json:"env_name" validate:"env_name|required"`
}

//PluginDefaultENV - plug-in default environment variables
type PluginDefaultENV struct {
	//corresponding plugin id
	//in: body
	//required: true
	PluginID string `json:"plugin_id" validate:"plugin_id"`
	//build version id
	//in: body
	//required: true
	VersionID string `json:"version_id" validate:"version_id"`
	//configuration item name
	//in: body
	//required: true
	ENVName string `json:"env_name" validate:"env_name"`
	//configuration item value
	//in: body
	//required: true
	ENVValue string `json:"env_value" validate:"env_value"`
	//can be modified by the user
	//in :body
	//required: false
	IsChange bool `json:"is_change" validate:"is_change|bool"`
}

//BuildPluginStruct
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
		// plugin CPU weight, default 125
		// in: body
		// required: true
		PluginCPU int `json:"plugin_cpu" validate:"plugin_cpu|required"`
		// plugin maximum memory, default 50
		// in: body
		// required: true
		PluginMemory int `json:"plugin_memory" validate:"plugin_memory|required"`
		// plugin cmd, default 50
		// in: body
		// required: false
		PluginCMD string `json:"plugin_cmd" validate:"plugin_cmd"`
		// the version number of the plugin
		// in: body
		// required: true
		BuildVersion string `json:"build_version" validate:"build_version|required"`
		// plugin build version number
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
		// version information, assist in selecting the plug-in version
		// in:body
		// required: true
		Info string `json:"info" validate:"info"`
		// operator
		// in: body
		// required: false
		Operator string `json:"operator" validate:"operator"`
		// tenant id
		// in: body
		// required: true
		TenantID string `json:"tenant_id" validate:"tenant_id"`
		// mirror address
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

//PluginBuildVersionStruct
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

//AllPluginBuildVersionStruct
//swagger:parameters allPluginVersions
type AllPluginBuildVersionStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	//in: path
	//required: true
	PluginID string `json:"plugin_id" validate:"plugin_id"`
}

//PluginSetStruct
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

//GetPluginsStruct
//swagger:parameters getPlugins
type GetPluginsStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
}

//GetPluginSetStruct
//swagger:parameters getPluginSet
type GetPluginSetStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
	//in: path
	//required: true
	ServiceAlias string `json:"service_alias"`
}

//DeletePluginSetStruct
//swagger:parameters deletePluginRelation
type DeletePluginSetStruct struct {
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
	//in: path
	//required: true
	ServiceAlias string `json:"service_alias"`
	//plugin id
	//in: path
	//required: true
	PluginID string `json:"plugin_id"`
}

//GetPluginEnvStruct
//swagger:parameters getPluginEnv getPluginDefaultEnv
type GetPluginEnvStruct struct {
	//tenant name
	//in: path
	//required: true
	TenantName string `json:"tenant_name"`
	// plugin id
	// in: path
	// required: true
	PluginID string `json:"plugin_id"`
	// build version id
	// in: path
	// required: true
	VersionID string `json:"version_id"`
}

//GetVersionEnvStruct
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

//SetVersionEnv
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
		// environment variable
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

//VersionEnv VersionEnv
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

//TransPlugins TransPlugins
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
