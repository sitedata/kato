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
	"github.com/docker/distribution/reference"
	"github.com/sirupsen/logrus"
)

//TenantPlugin model
type TenantPlugin struct {
	Model
	PluginID string `gorm:"column:plugin_id;size:32"`
	//plugin name
	PluginName string `gorm:"column:plugin_name;size:32" json:"plugin_name"`
	//plugin describe
	PluginInfo string `gorm:"column:plugin_info;size:255" json:"plugin_info"`
	//plugin build by docker image name
	ImageURL string `gorm:"column:image_url" json:"image_url"`
	//plugin build by git code url
	GitURL string `gorm:"column:git_url" json:"git_url"`
	//build mode
	BuildModel string `gorm:"column:build_model" json:"build_model"`
	//plugin model InitPlugin,InBoundNetPlugin,OutBoundNetPlugin
	PluginModel string `gorm:"column:plugin_model" json:"plugin_model"`
	//tenant id
	TenantID string `gorm:"column:tenant_id" json:"tenant_id"`
	//tenant_name Used to calculate CPU and Memory.
	Domain string `gorm:"column:domain" json:"domain"`
	//gitlab; github The deprecated
	CodeFrom string `gorm:"column:code_from" json:"code_from"`
}

//TableName
func (t *TenantPlugin) TableName() string {
	return "tenant_plugin"
}

//TenantPluginDefaultENV plugin default env config
type TenantPluginDefaultENV struct {
	Model
	//plugin id
	PluginID string `gorm:"column:plugin_id" json:"plugin_id"`
	//plugin version
	VersionID string `gorm:"column:version_id;size:32" json:"version_id"`
	//env name
	ENVName string `gorm:"column:env_name" json:"env_name"`
	//env value
	ENVValue string `gorm:"column:env_value" json:"env_value"`
	//value is change
	IsChange bool `gorm:"column:is_change;default:false" json:"is_change"`
}

//TableName
func (t *TenantPluginDefaultENV) TableName() string {
	return "tenant_plugin_default_env"
}

//TenantPluginBuildVersion plugin build version
type TenantPluginBuildVersion struct {
	Model
	//plugin version eg v1.0.0
	VersionID string `gorm:"column:version_id;size:32" json:"version_id"`
	//deploy version eg 20180528071717
	DeployVersion   string `gorm:"column:deploy_version;size:32" json:"deploy_version"`
	PluginID        string `gorm:"column:plugin_id;size:32" json:"plugin_id"`
	Kind            string `gorm:"column:kind;size:24" json:"kind"`
	BaseImage       string `gorm:"column:base_image;size:200" json:"base_image"`
	BuildLocalImage string `gorm:"column:build_local_image;size:200" json:"build_local_image"`
	BuildTime       string `gorm:"column:build_time" json:"build_time"`
	Repo            string `gorm:"column:repo" json:"repo"`
	GitURL          string `gorm:"column:git_url" json:"git_url"`
	Info            string `gorm:"column:info" json:"info"`
	Status          string `gorm:"column:status;size:24" json:"status"`
	// container default cpu
	ContainerCPU int `gorm:"column:container_cpu;default:125" json:"container_cpu"`
	// container default memory
	ContainerMemory int `gorm:"column:container_memory;default:64" json:"container_memory"`
	// container args
	ContainerCMD string `gorm:"column:container_cmd;size:2048" json:"container_cmd"`
}

//TableName
func (t *TenantPluginBuildVersion) TableName() string {
	return "tenant_plugin_build_version"
}

//CreateShareImage
func (t *TenantPluginBuildVersion) CreateShareImage(hubURL, namespace string) (string, error) {
	_, err := reference.ParseAnyReference(t.BuildLocalImage)
	if err != nil {
		logrus.Errorf("reference image error: %s", err.Error())
		return "", err
	}
	image := ParseImage(t.BuildLocalImage)
	if hubURL != "" {
		image.Host = hubURL
	}
	if namespace != "" {
		image.Namespace = namespace
	}
	image.Name = image.Name + "_" + t.VersionID
	return image.String(), nil
}

//TenantPluginVersionEnv
type TenantPluginVersionEnv struct {
	Model
	//VersionID string `gorm:"column:version_id;size:32"`
	PluginID  string `gorm:"column:plugin_id;size:32" json:"plugin_id"`
	EnvName   string `gorm:"column:env_name" json:"env_name"`
	EnvValue  string `gorm:"column:env_value" json:"env_value"`
	ServiceID string `gorm:"column:service_id" json:"service_id"`
}

//TableName
func (t *TenantPluginVersionEnv) TableName() string {
	return "tenant_plugin_version_env"
}

//TenantPluginVersionDiscoverConfig service plugin config that can be dynamic discovery
type TenantPluginVersionDiscoverConfig struct {
	Model
	PluginID  string `gorm:"column:plugin_id;size:32" json:"plugin_id"`
	ServiceID string `gorm:"column:service_id;size:32" json:"service_id"`
	ConfigStr string `gorm:"column:config_str;" sql:"type:text;" json:"config_str"`
}

//TableName
func (t *TenantPluginVersionDiscoverConfig) TableName() string {
	return "tenant_plugin_version_config"
}

//TenantServicePluginRelation
type TenantServicePluginRelation struct {
	Model
	VersionID   string `gorm:"column:version_id;size:32" json:"version_id"`
	PluginID    string `gorm:"column:plugin_id;size:32" json:"plugin_id"`
	ServiceID   string `gorm:"column:service_id;size:32" json:"service_id"`
	PluginModel string `gorm:"column:plugin_model;size:24" json:"plugin_model"`
	// container default cpu  v3.5.1 add
	ContainerCPU int `gorm:"column:container_cpu;default:125" json:"container_cpu"`
	// container default memory  v3.5.1 add
	ContainerMemory int  `gorm:"column:container_memory;default:64" json:"container_memory"`
	Switch          bool `gorm:"column:switch;default:false" json:"switch"`
}

//TableName
func (t *TenantServicePluginRelation) TableName() string {
	return "tenant_service_plugin_relation"
}

//TenantServicesStreamPluginPort - port mapping information after binding the stream type plug-in
type TenantServicesStreamPluginPort struct {
	Model
	TenantID      string `gorm:"column:tenant_id;size:32" validate:"tenant_id|between:30,33" json:"tenant_id"`
	ServiceID     string `gorm:"column:service_id;size:32" validate:"service_id|between:30,33" json:"service_id"`
	PluginModel   string `gorm:"column:plugin_model;size:24" json:"plugin_model"`
	ContainerPort int    `gorm:"column:container_port" validate:"container_port|required|numeric_between:1,65535" json:"container_port"`
	PluginPort    int    `gorm:"column:plugin_port" json:"plugin_port"`
}

//TableName
func (t *TenantServicesStreamPluginPort) TableName() string {
	return "tenant_services_stream_plugin_port"
}

//Plugin model - plug-in tag

//TODO: Plug-in type name regulations
//@ 1. Plug-in category xxx-plugin
//@ 2. Major subdivision colon + subdivision xxx-plugin:up or xxx-plugin:down

//InitPlugin
var InitPlugin = "init-plugin"

//InBoundNetPlugin
var InBoundNetPlugin = "net-plugin:up"

//OutBoundNetPlugin
var OutBoundNetPlugin = "net-plugin:down"

//InBoundAndOutBoundNetPlugin
var InBoundAndOutBoundNetPlugin = "net-plugin:in-and-out"

//GeneralPlugin - general plugin, default classification, lowest priority
var GeneralPlugin = "general-plugin"
