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

package api

import (
	"net/http"
)

// ClusterInterface -
type ClusterInterface interface {
	GetClusterInfo(w http.ResponseWriter, r *http.Request)
	MavenSettingList(w http.ResponseWriter, r *http.Request)
	MavenSettingAdd(w http.ResponseWriter, r *http.Request)
	MavenSettingUpdate(w http.ResponseWriter, r *http.Request)
	MavenSettingDelete(w http.ResponseWriter, r *http.Request)
	MavenSettingDetail(w http.ResponseWriter, r *http.Request)
}

//TenantInterface interface
type TenantInterface interface {
	TenantInterfaceWithV1
	AllTenantResources(w http.ResponseWriter, r *http.Request)
	TenantResources(w http.ResponseWriter, r *http.Request)
	ServiceResources(w http.ResponseWriter, r *http.Request)
	Tenant(w http.ResponseWriter, r *http.Request)
	Tenants(w http.ResponseWriter, r *http.Request)
	ServicesInfo(w http.ResponseWriter, r *http.Request)
	TenantsWithResource(w http.ResponseWriter, r *http.Request)
	TenantsQuery(w http.ResponseWriter, r *http.Request)
	TenantsGetByName(w http.ResponseWriter, r *http.Request)
	SumTenants(w http.ResponseWriter, r *http.Request)
	SingleTenantResources(w http.ResponseWriter, r *http.Request)
	GetSupportProtocols(w http.ResponseWriter, r *http.Request)
	TransPlugins(w http.ResponseWriter, r *http.Request)
	ServicesCount(w http.ResponseWriter, r *http.Request)
	GetManyDeployVersion(w http.ResponseWriter, r *http.Request)
	LimitTenantMemory(w http.ResponseWriter, r *http.Request)
	TenantResourcesStatus(w http.ResponseWriter, r *http.Request)
	CheckResourceName(w http.ResponseWriter, r *http.Request)
	Log(w http.ResponseWriter, r *http.Request)
}

//ServiceInterface ServiceInterface
type ServiceInterface interface {
	SetLanguage(w http.ResponseWriter, r *http.Request)
	SingleServiceInfo(w http.ResponseWriter, r *http.Request)
	CheckCode(w http.ResponseWriter, r *http.Request)
	Event(w http.ResponseWriter, r *http.Request)
	BuildList(w http.ResponseWriter, r *http.Request)
	CreateService(w http.ResponseWriter, r *http.Request)
	UpdateService(w http.ResponseWriter, r *http.Request)
	Dependency(w http.ResponseWriter, r *http.Request)
	Env(w http.ResponseWriter, r *http.Request)
	Ports(w http.ResponseWriter, r *http.Request)
	PutPorts(w http.ResponseWriter, r *http.Request)
	PortOuterController(w http.ResponseWriter, r *http.Request)
	PortInnerController(w http.ResponseWriter, r *http.Request)
	RollBack(w http.ResponseWriter, r *http.Request)
	AddVolume(w http.ResponseWriter, r *http.Request)
	UpdVolume(w http.ResponseWriter, r *http.Request)
	DeleteVolume(w http.ResponseWriter, r *http.Request)
	Pods(w http.ResponseWriter, r *http.Request)
	VolumeDependency(w http.ResponseWriter, r *http.Request)
	Probe(w http.ResponseWriter, r *http.Request)
	Label(w http.ResponseWriter, r *http.Request)
	Share(w http.ResponseWriter, r *http.Request)
	ShareResult(w http.ResponseWriter, r *http.Request)
	BuildVersionInfo(w http.ResponseWriter, r *http.Request)
	GetDeployVersion(w http.ResponseWriter, r *http.Request)
	AutoscalerRules(w http.ResponseWriter, r *http.Request)
	ScalingRecords(w http.ResponseWriter, r *http.Request)
	AddServiceMonitors(w http.ResponseWriter, r *http.Request)
	DeleteServiceMonitors(w http.ResponseWriter, r *http.Request)
	UpdateServiceMonitors(w http.ResponseWriter, r *http.Request)
}

//TenantInterfaceWithV1 funcs for both v2 and v1
type TenantInterfaceWithV1 interface {
	StartService(w http.ResponseWriter, r *http.Request)
	StopService(w http.ResponseWriter, r *http.Request)
	RestartService(w http.ResponseWriter, r *http.Request)
	VerticalService(w http.ResponseWriter, r *http.Request)
	HorizontalService(w http.ResponseWriter, r *http.Request)
	BuildService(w http.ResponseWriter, r *http.Request)
	DeployService(w http.ResponseWriter, r *http.Request)
	UpgradeService(w http.ResponseWriter, r *http.Request)
	StatusService(w http.ResponseWriter, r *http.Request)
	StatusServiceList(w http.ResponseWriter, r *http.Request)
	StatusContainerID(w http.ResponseWriter, r *http.Request)
}

//LogInterface log interface
type LogInterface interface {
	HistoryLogs(w http.ResponseWriter, r *http.Request)
	LogList(w http.ResponseWriter, r *http.Request)
	LogFile(w http.ResponseWriter, r *http.Request)
	LogSocket(w http.ResponseWriter, r *http.Request)
	LogByAction(w http.ResponseWriter, r *http.Request)
	TenantLogByAction(w http.ResponseWriter, r *http.Request)
	Events(w http.ResponseWriter, r *http.Request)
	EventLog(w http.ResponseWriter, r *http.Request)
}

//PluginInterface plugin interface
type PluginInterface interface {
	PluginAction(w http.ResponseWriter, r *http.Request)
	PluginDefaultENV(w http.ResponseWriter, r *http.Request)
	PluginBuild(w http.ResponseWriter, r *http.Request)
	GetAllPluginBuildVersions(w http.ResponseWriter, r *http.Request)
	GetPluginBuildVersion(w http.ResponseWriter, r *http.Request)
	DeletePluginBuildVersion(w http.ResponseWriter, r *http.Request)
	//plugin
	PluginSet(w http.ResponseWriter, r *http.Request)
	DeletePluginRelation(w http.ResponseWriter, r *http.Request)
	GePluginEnvWhichCanBeSet(w http.ResponseWriter, r *http.Request)
	UpdateVersionEnv(w http.ResponseWriter, r *http.Request)
	GetPluginDefaultEnvs(w http.ResponseWriter, r *http.Request)
	SharePlugin(w http.ResponseWriter, r *http.Request)
	SharePluginResult(w http.ResponseWriter, r *http.Request)
	BatchInstallPlugins(w http.ResponseWriter, r *http.Request)
	BatchBuildPlugins(w http.ResponseWriter, r *http.Request)
}

//RulesInterface RulesInterface
type RulesInterface interface {
	SetDownStreamRule(w http.ResponseWriter, r *http.Request)
	GetDownStreamRule(w http.ResponseWriter, r *http.Request)
	DeleteDownStreamRule(w http.ResponseWriter, r *http.Request)
	UpdateDownStreamRule(w http.ResponseWriter, r *http.Request)
}

//AppInterface app handle interface
type AppInterface interface {
	ExportApp(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
	NewUpload(w http.ResponseWriter, r *http.Request)
	ImportID(w http.ResponseWriter, r *http.Request)
	ImportApp(w http.ResponseWriter, r *http.Request)
}

// ApplicationInterface tenant application interface
type ApplicationInterface interface {
	CreateApp(w http.ResponseWriter, r *http.Request)
	BatchCreateApp(w http.ResponseWriter, r *http.Request)
	UpdateApp(w http.ResponseWriter, r *http.Request)
	ListApps(w http.ResponseWriter, r *http.Request)
	ListComponents(w http.ResponseWriter, r *http.Request)
	BatchBindService(w http.ResponseWriter, r *http.Request)
	DeleteApp(w http.ResponseWriter, r *http.Request)
	AddConfigGroup(w http.ResponseWriter, r *http.Request)
	UpdateConfigGroup(w http.ResponseWriter, r *http.Request)

	BatchUpdateComponentPorts(w http.ResponseWriter, r *http.Request)
	GetAppStatus(w http.ResponseWriter, r *http.Request)
	Install(w http.ResponseWriter, r *http.Request)
	ListServices(w http.ResponseWriter, r *http.Request)
	ListHelmAppReleases(w http.ResponseWriter, r *http.Request)

	DeleteConfigGroup(w http.ResponseWriter, r *http.Request)
	ListConfigGroups(w http.ResponseWriter, r *http.Request)
	SyncComponents(w http.ResponseWriter, r *http.Request)
	SyncAppConfigGroups(w http.ResponseWriter, r *http.Request)
	ListAppStatuses(w http.ResponseWriter, r *http.Request)
}

//Gatewayer gateway api interface
type Gatewayer interface {
	HTTPRule(w http.ResponseWriter, r *http.Request)
	TCPRule(w http.ResponseWriter, r *http.Request)
	GetAvailablePort(w http.ResponseWriter, r *http.Request)
	RuleConfig(w http.ResponseWriter, r *http.Request)
	Certificate(w http.ResponseWriter, r *http.Request)
}

// ThirdPartyServicer is an interface for defining methods for third-party service.
type ThirdPartyServicer interface {
	Endpoints(w http.ResponseWriter, r *http.Request)
}

// Labeler is an interface for defining methods to get information of labels.
type Labeler interface {
	Labels(w http.ResponseWriter, r *http.Request)
}

// AppRestoreInterface defines api methods to restore app.
// app means market service.
type AppRestoreInterface interface {
	RestoreEnvs(w http.ResponseWriter, r *http.Request)
	RestorePorts(w http.ResponseWriter, r *http.Request)
	RestoreVolumes(w http.ResponseWriter, r *http.Request)
	RestoreProbe(w http.ResponseWriter, r *http.Request)
	RestoreDeps(w http.ResponseWriter, r *http.Request)
	RestoreDepVols(w http.ResponseWriter, r *http.Request)
	RestorePlugins(w http.ResponseWriter, r *http.Request)
}

// PodInterface defines api methods about k8s pods.
type PodInterface interface {
	PodDetail(w http.ResponseWriter, r *http.Request)
}
