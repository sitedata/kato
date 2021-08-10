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

package version2

import (
	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/controller"
	"github.com/gridworkz/kato/api/middleware"
	"github.com/gridworkz/kato/cmd/api/option"
	dbmodel "github.com/gridworkz/kato/db/model"
)

//V2 v2
type V2 struct {
	Cfg *option.Config
}

//Routes routes
func (v2 *V2) Routes() chi.Router {
	r := chi.NewRouter()
	license := middleware.NewLicense(v2.Cfg)
	r.Use(license.Verify)
	r.Get("/show", controller.GetManager().Show)
	r.Post("/show", controller.GetManager().Show)
	r.Mount("/tenants", v2.tenantRouter())
	r.Mount("/cluster", v2.clusterRouter())
	r.Mount("/notificationEvent", v2.notificationEventRouter())
	r.Mount("/resources", v2.resourcesRouter())
	r.Mount("/prometheus", v2.prometheusRouter())
	r.Get("/event", controller.GetManager().Event)
	r.Mount("/app", v2.appRouter())
	r.Get("/health", controller.GetManager().Health)
	r.Post("/alertmanager-webhook", controller.GetManager().AlertManagerWebHook)
	r.Get("/version", controller.GetManager().Version)
	// deprecated use /gateway/ports
	r.Mount("/port", v2.portRouter())
	// deprecated, use /events/<event_id>/log
	r.Get("/event-log", controller.GetManager().LogByAction)
	r.Mount("/events", v2.eventsRouter())
	r.Get("/gateway/ips", controller.GetGatewayIPs)
	r.Get("/gateway/ports", controller.GetManager().GetAvailablePort)
	r.Get("/volume-options", controller.VolumeOptions)
	r.Get("/volume-options/page/{page}/size/{pageSize}", controller.ListVolumeType)
	r.Post("/volume-options", controller.VolumeSetVar)
	r.Delete("/volume-options/{volume_type}", controller.DeleteVolumeType)
	r.Put("/volume-options/{volume_type}", controller.UpdateVolumeType)
	r.Mount("/enterprise/{enterprise_id}", v2.enterpriseRouter())
	r.Mount("/monitor", v2.monitorRouter())
	return r
}

func (v2 *V2) monitorRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/metrics", controller.GetMonitorMetrics)
	return r
}

func (v2 *V2) enterpriseRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/running-services", controller.GetRunningServices)
	return r
}

func (v2 *V2) eventsRouter() chi.Router {
	r := chi.NewRouter()
	// get target's event list with page
	r.Get("/", controller.GetManager().Events)
	// get target's event content
	r.Get("/{eventID}/log", controller.GetManager().EventLog)
	return r
}

func (v2 *V2) clusterRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", controller.GetManager().GetClusterInfo)
	r.Get("/builder/mavensetting", controller.GetManager().MavenSettingList)
	r.Post("/builder/mavensetting", controller.GetManager().MavenSettingAdd)
	r.Get("/builder/mavensetting/{name}", controller.GetManager().MavenSettingDetail)
	r.Put("/builder/mavensetting/{name}", controller.GetManager().MavenSettingUpdate)
	r.Delete("/builder/mavensetting/{name}", controller.GetManager().MavenSettingDelete)
	return r
}

func (v2 *V2) tenantRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", controller.GetManager().Tenants)
	r.Mount("/{tenant_name}", v2.tenantNameRouter())
	r.Get("/", controller.GetManager().Tenants)
	r.Get("/services-count", controller.GetManager().ServicesCount)
	return r
}

func (v2 *V2) tenantNameRouter() chi.Router {
	r := chi.NewRouter()
	//Initialize tenant and service letter
	r.Use(middleware.InitTenant)
	r.Put("/", controller.GetManager().Tenant)
	r.Get("/", controller.GetManager().Tenant)
	r.Delete("/", controller.GetManager().Tenant)
	//Logs in the tenant
	r.Post("/event-log", controller.GetManager().TenantLogByAction)
	r.Get("/protocols", controller.GetManager().GetSupportProtocols)
	//Plug-in pre-installation
	r.Post("/transplugins", controller.GetManager().TransPlugins)
	//Code detection
	r.Post("/code-check", controller.GetManager().CheckCode)
	r.Post("/servicecheck", controller.Check)
	r.Get("/servicecheck/{uuid}", controller.GetServiceCheckInfo)
	r.Get("/resources", controller.GetManager().SingleTenantResources)
	r.Get("/services", controller.GetManager().ServicesInfo)
	//Create application
	r.Post("/services", middleware.WrapEL(controller.GetManager().CreateService, dbmodel.TargetTypeService, "create-service", dbmodel.SYNEVENTTYPE))
	r.Post("/plugin", controller.GetManager().PluginAction)
	r.Post("/plugins/{plugin_id}/share", controller.GetManager().SharePlugin)
	r.Get("/plugins/{plugin_id}/share/{share_id}", controller.GetManager().SharePluginResult)
	r.Get("/plugin", controller.GetManager().PluginAction)
	// batch install and build plugins
	r.Post("/plugins", controller.GetManager().BatchInstallPlugins)
	r.Post("/batch-build-plugins", controller.GetManager().BatchBuildPlugins)
	r.Post("/services_status", controller.GetManager().StatusServiceList)
	r.Mount("/services/{service_alias}", v2.serviceRouter())
	r.Mount("/plugin/{plugin_id}", v2.pluginRouter())
	r.Get("/event", controller.GetManager().Event)
	r.Get("/chargesverify", controller.ChargesVerifyController)
	//tenant app
	r.Get("/pods/{pod_name}", controller.GetManager().PodDetail)
	r.Post("/apps", controller.GetManager().CreateApp)
	r.Post("/batch_create_apps", controller.GetManager().BatchCreateApp)
	r.Get("/apps", controller.GetManager().ListApps)
	r.Post("/checkResourceName", controller.GetManager().CheckResourceName)
	r.Get("/appstatuses", controller.GetManager().ListAppStatuses)
	r.Mount("/apps/{app_id}", v2.applicationRouter())
	//get some service pod info
	r.Get("/pods", controller.Pods)
	r.Get("/pod_nums", controller.PodNums)
	//app backup
	r.Get("/groupapp/backups", controller.Backups)
	r.Post("/groupapp/backups", controller.NewBackups)
	r.Post("/groupapp/backupcopy", controller.BackupCopy)
	r.Get("/groupapp/backups/{backup_id}", controller.GetBackup)
	r.Delete("/groupapp/backups/{backup_id}", controller.DeleteBackup)
	r.Post("/groupapp/backups/{backup_id}/restore", controller.Restore)
	r.Get("/groupapp/backups/{backup_id}/restore/{restore_id}", controller.RestoreResult)
	r.Post("/deployversions", controller.GetManager().GetManyDeployVersion)
	//Team resource limit
	r.Post("/limit_memory", controller.GetManager().LimitTenantMemory)
	r.Get("/limit_memory", controller.GetManager().TenantResourcesStatus)

	// Gateway
	r.Post("/http-rule", controller.GetManager().HTTPRule)
	r.Delete("/http-rule", controller.GetManager().HTTPRule)
	r.Put("/http-rule", controller.GetManager().HTTPRule)
	r.Post("/tcp-rule", controller.GetManager().TCPRule)
	r.Delete("/tcp-rule", controller.GetManager().TCPRule)
	r.Put("/tcp-rule", controller.GetManager().TCPRule)
	r.Mount("/gateway", v2.gatewayRouter())

	//batch operation
	r.Post("/batchoperation", controller.BatchOperation)

	return r
}

func (v2 *V2) gatewayRouter() chi.Router {
	r := chi.NewRouter()
	r.Put("/certificate", controller.GetManager().Certificate)

	return r
}

func (v2 *V2) serviceRouter() chi.Router {
	r := chi.NewRouter()
	//Initialize application information
	r.Use(middleware.InitService)
	r.Put("/", middleware.WrapEL(controller.GetManager().UpdateService, dbmodel.TargetTypeService, "update-service", dbmodel.SYNEVENTTYPE))
	// component build
	r.Post("/build", middleware.WrapEL(controller.GetManager().BuildService, dbmodel.TargetTypeService, "build-service", dbmodel.ASYNEVENTTYPE))
	// component start
	r.Post("/start", middleware.WrapEL(controller.GetManager().StartService, dbmodel.TargetTypeService, "start-service", dbmodel.ASYNEVENTTYPE))
	// component stop event set to synchronous event, not wait.
	r.Post("/stop", middleware.WrapEL(controller.GetManager().StopService, dbmodel.TargetTypeService, "stop-service", dbmodel.SYNEVENTTYPE))
	r.Post("/restart", middleware.WrapEL(controller.GetManager().RestartService, dbmodel.TargetTypeService, "restart-service", dbmodel.ASYNEVENTTYPE))
	//Application scaling
	r.Put("/vertical", middleware.WrapEL(controller.GetManager().VerticalService, dbmodel.TargetTypeService, "vertical-service", dbmodel.ASYNEVENTTYPE))
	r.Put("/horizontal", middleware.WrapEL(controller.GetManager().HorizontalService, dbmodel.TargetTypeService, "horizontal-service", dbmodel.ASYNEVENTTYPE))

	//Set application language (act)
	r.Post("/language", middleware.WrapEL(controller.GetManager().SetLanguage, dbmodel.TargetTypeService, "set-language", dbmodel.SYNEVENTTYPE))
	//Application information acquisition, modification and deletion (source)
	r.Get("/", controller.GetManager().SingleServiceInfo)
	r.Delete("/", middleware.WrapEL(controller.GetManager().SingleServiceInfo, dbmodel.TargetTypeService, "delete-service", dbmodel.SYNEVENTTYPE))
	//Application upgrade (act)
	r.Post("/upgrade", middleware.WrapEL(controller.GetManager().UpgradeService, dbmodel.TargetTypeService, "upgrade-service", dbmodel.ASYNEVENTTYPE))
	//Application status acquisition (act)
	r.Get("/status", controller.GetManager().StatusService)
	//Build version list
	r.Get("/build-list", controller.GetManager().BuildList)
	//Build version operation
	r.Get("/build-version/{build_version}", controller.GetManager().BuildVersionInfo)
	r.Put("/build-version/{build_version}", controller.GetManager().BuildVersionInfo)
	r.Get("/deployversion", controller.GetManager().GetDeployVersion)
	r.Delete("/build-version/{build_version}", middleware.WrapEL(controller.GetManager().BuildVersionInfo, dbmodel.TargetTypeService, "delete-buildversion", dbmodel.SYNEVENTTYPE))
	//Application sharing
	r.Post("/share", middleware.WrapEL(controller.GetManager().Share, dbmodel.TargetTypeService, "share-service", dbmodel.SYNEVENTTYPE))
	r.Get("/share/{share_id}", controller.GetManager().ShareResult)
	r.Get("/logs", controller.GetManager().HistoryLogs)
	r.Get("/log-file", controller.GetManager().LogList)
	r.Get("/log-instance", controller.GetManager().LogSocket)
	r.Post("/event-log", controller.GetManager().LogByAction)

	//Application dependency addition and deletion (source)
	r.Post("/dependency", middleware.WrapEL(controller.GetManager().Dependency, dbmodel.TargetTypeService, "add-service-dependency", dbmodel.SYNEVENTTYPE))
	r.Delete("/dependency", middleware.WrapEL(controller.GetManager().Dependency, dbmodel.TargetTypeService, "delete-service-dependency", dbmodel.SYNEVENTTYPE))
	//Addition, deletion and modification of environment variables (source)
	r.Post("/env", middleware.WrapEL(controller.GetManager().Env, dbmodel.TargetTypeService, "add-service-env", dbmodel.SYNEVENTTYPE))
	r.Put("/env", middleware.WrapEL(controller.GetManager().Env, dbmodel.TargetTypeService, "update-service-env", dbmodel.SYNEVENTTYPE))
	r.Delete("/env", middleware.WrapEL(controller.GetManager().Env, dbmodel.TargetTypeService, "delete-service-env", dbmodel.SYNEVENTTYPE))
	//Adding, deleting and modifying port variables (source)
	r.Post("/ports", middleware.WrapEL(controller.GetManager().Ports, dbmodel.TargetTypeService, "add-service-port", dbmodel.SYNEVENTTYPE))
	r.Put("/ports", middleware.WrapEL(controller.GetManager().PutPorts, dbmodel.TargetTypeService, "update-service-port-old", dbmodel.SYNEVENTTYPE))
	r.Put("/ports/{port}", middleware.WrapEL(controller.GetManager().Ports, dbmodel.TargetTypeService, "update-service-port", dbmodel.SYNEVENTTYPE))
	r.Delete("/ports/{port}", middleware.WrapEL(controller.GetManager().Ports, dbmodel.TargetTypeService, "delete-service-port", dbmodel.SYNEVENTTYPE))
	r.Put("/ports/{port}/outer", middleware.WrapEL(controller.GetManager().PortOuterController, dbmodel.TargetTypeService, "handle-service-outerport", dbmodel.SYNEVENTTYPE))
	r.Put("/ports/{port}/inner", middleware.WrapEL(controller.GetManager().PortInnerController, dbmodel.TargetTypeService, "handle-service-innerport", dbmodel.SYNEVENTTYPE))

	//Application version rollback (act)
	r.Post("/rollback", middleware.WrapEL(controller.GetManager().RollBack, dbmodel.TargetTypeService, "rollback-service", dbmodel.ASYNEVENTTYPE))

	//Persistent information API v2.1 supports multiple persistence formats
	r.Post("/volumes", middleware.WrapEL(controller.AddVolume, dbmodel.TargetTypeService, "add-service-volume", dbmodel.SYNEVENTTYPE))
	r.Put("/volumes", middleware.WrapEL(controller.GetManager().UpdVolume, dbmodel.TargetTypeService, "update-service-volume", dbmodel.SYNEVENTTYPE))
	r.Get("/volumes", controller.GetVolume)
	r.Delete("/volumes/{volume_name}", middleware.WrapEL(controller.DeleteVolume, dbmodel.TargetTypeService, "delete-service-volume", dbmodel.SYNEVENTTYPE))
	r.Post("/depvolumes", middleware.WrapEL(controller.AddVolumeDependency, dbmodel.TargetTypeService, "add-service-depvolume", dbmodel.SYNEVENTTYPE))
	r.Delete("/depvolumes", middleware.WrapEL(controller.DeleteVolumeDependency, dbmodel.TargetTypeService, "delete-service-depvolume", dbmodel.SYNEVENTTYPE))
	r.Get("/depvolumes", controller.GetDepVolume)
	//Persistent Information API v2
	r.Post("/volume-dependency", middleware.WrapEL(controller.GetManager().VolumeDependency, dbmodel.TargetTypeService, "add-service-depvolume", dbmodel.SYNEVENTTYPE))
	r.Delete("/volume-dependency", middleware.WrapEL(controller.GetManager().VolumeDependency, dbmodel.TargetTypeService, "delete-service-depvolume", dbmodel.SYNEVENTTYPE))
	r.Post("/volume", middleware.WrapEL(controller.GetManager().AddVolume, dbmodel.TargetTypeService, "add-service-volume", dbmodel.SYNEVENTTYPE))
	r.Delete("/volume", middleware.WrapEL(controller.GetManager().DeleteVolume, dbmodel.TargetTypeService, "delete-service-volume", dbmodel.SYNEVENTTYPE))

	//Get the application instance situation (source)
	r.Get("/pods", controller.GetManager().Pods)

	//Application probe addition, deletion and modification (source)
	r.Post("/probe", middleware.WrapEL(controller.GetManager().Probe, dbmodel.TargetTypeService, "add-service-probe", dbmodel.SYNEVENTTYPE))
	r.Put("/probe", middleware.WrapEL(controller.GetManager().Probe, dbmodel.TargetTypeService, "update-service-probe", dbmodel.SYNEVENTTYPE))
	r.Delete("/probe", middleware.WrapEL(controller.GetManager().Probe, dbmodel.TargetTypeService, "delete-service-probe", dbmodel.SYNEVENTTYPE))

	r.Post("/label", middleware.WrapEL(controller.GetManager().Label, dbmodel.TargetTypeService, "add-service-label", dbmodel.SYNEVENTTYPE))
	r.Put("/label", middleware.WrapEL(controller.GetManager().Label, dbmodel.TargetTypeService, "update-service-label", dbmodel.SYNEVENTTYPE))
	r.Delete("/label", middleware.WrapEL(controller.GetManager().Label, dbmodel.TargetTypeService, "delete-service-label", dbmodel.SYNEVENTTYPE))

	//Plug-in
	r.Mount("/plugin", v2.serviceRelatePluginRouter())

	//rule
	r.Mount("/net-rule", v2.rulesRouter())
	r.Get("/deploy-info", controller.GetServiceDeployInfo)

	// third-party service
	r.Post("/endpoints", middleware.WrapEL(controller.GetManager().Endpoints, dbmodel.TargetTypeService, "add-thirdpart-service", dbmodel.SYNEVENTTYPE))
	r.Put("/endpoints", middleware.WrapEL(controller.GetManager().Endpoints, dbmodel.TargetTypeService, "update-thirdpart-service", dbmodel.SYNEVENTTYPE))
	r.Delete("/endpoints", middleware.WrapEL(controller.GetManager().Endpoints, dbmodel.TargetTypeService, "delete-thirdpart-service", dbmodel.SYNEVENTTYPE))
	r.Get("/endpoints", controller.GetManager().Endpoints)

	// gateway
	r.Put("/rule-config", middleware.WrapEL(controller.GetManager().RuleConfig, dbmodel.TargetTypeService, "update-service-gateway-rule", dbmodel.SYNEVENTTYPE))

	// app restore
	r.Post("/app-restore/envs", middleware.WrapEL(controller.GetManager().RestoreEnvs, dbmodel.TargetTypeService, "app-restore-envs", dbmodel.SYNEVENTTYPE))
	r.Post("/app-restore/ports", middleware.WrapEL(controller.GetManager().RestorePorts, dbmodel.TargetTypeService, "app-restore-ports", dbmodel.SYNEVENTTYPE))
	r.Post("/app-restore/volumes", middleware.WrapEL(controller.GetManager().RestoreVolumes, dbmodel.TargetTypeService, "app-restore-volumes", dbmodel.SYNEVENTTYPE))
	r.Post("/app-restore/probe", middleware.WrapEL(controller.GetManager().RestoreProbe, dbmodel.TargetTypeService, "app-restore-probe", dbmodel.SYNEVENTTYPE))
	r.Post("/app-restore/deps", middleware.WrapEL(controller.GetManager().RestoreDeps, dbmodel.TargetTypeService, "app-restore-deps", dbmodel.SYNEVENTTYPE))
	r.Post("/app-restore/depvols", middleware.WrapEL(controller.GetManager().RestoreDepVols, dbmodel.TargetTypeService, "app-restore-depvols", dbmodel.SYNEVENTTYPE))
	r.Post("/app-restore/plugins", middleware.WrapEL(controller.GetManager().RestorePlugins, dbmodel.TargetTypeService, "app-restore-plugins", dbmodel.SYNEVENTTYPE))

	r.Get("/pods/{pod_name}/detail", controller.GetManager().PodDetail)

	// autoscaler
	r.Post("/xparules", middleware.WrapEL(controller.GetManager().AutoscalerRules, dbmodel.TargetTypeService, "add-app-autoscaler-rule", dbmodel.SYNEVENTTYPE))
	r.Put("/xparules", middleware.WrapEL(controller.GetManager().AutoscalerRules, dbmodel.TargetTypeService, "update-app-autoscaler-rule", dbmodel.SYNEVENTTYPE))
	r.Get("/xparecords", controller.GetManager().ScalingRecords)

	//service monitor
	r.Post("/service-monitors", middleware.WrapEL(controller.GetManager().AddServiceMonitors, dbmodel.TargetTypeService, "add-app-service-monitor", dbmodel.SYNEVENTTYPE))
	r.Put("/service-monitors/{name}", middleware.WrapEL(controller.GetManager().UpdateServiceMonitors, dbmodel.TargetTypeService, "update-app-service-monitor", dbmodel.SYNEVENTTYPE))
	r.Delete("/service-monitors/{name}", middleware.WrapEL(controller.GetManager().DeleteServiceMonitors, dbmodel.TargetTypeService, "delete-app-service-monitor", dbmodel.SYNEVENTTYPE))

	r.Get("/log", controller.GetManager().Log)

	return r
}

func (v2 *V2) applicationRouter() chi.Router {
	r := chi.NewRouter()
	// Init Application
	r.Use(middleware.InitApplication)
	// Operation application
	r.Put("/", controller.GetManager().UpdateApp)
	r.Delete("/", controller.GetManager().DeleteApp)
	// Get services under application
	r.Get("/services", controller.GetManager().ListServices)
	// bind components
	r.Put("/services", controller.GetManager().BatchBindService)
	// Application configuration group
	r.Post("/configgroups", controller.GetManager().AddConfigGroup)
	r.Put("/configgroups/{config_group_name}", controller.GetManager().UpdateConfigGroup)

	r.Put("/ports", controller.GetManager().BatchUpdateComponentPorts)
	r.Put("/status", controller.GetManager().GetAppStatus)
	// status
	r.Post("/install", controller.GetManager().Install)
	r.Get("/releases", controller.GetManager().ListHelmAppReleases)

	r.Delete("/configgroups/{config_group_name}", controller.GetManager().DeleteConfigGroup)
	r.Get("/configgroups", controller.GetManager().ListConfigGroups)

	// Synchronize component information, full coverage
	r.Post("/components", controller.GetManager().SyncComponents)
	r.Post("/app-config-groups", controller.GetManager().SyncAppConfigGroups)
	return r
}

func (v2 *V2) resourcesRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/labels", controller.GetManager().Labels)
	r.Post("/tenants", controller.GetManager().TenantResources)
	r.Post("/services", controller.GetManager().ServiceResources)
	r.Get("/tenants/sum", controller.GetManager().SumTenants)
	//tenants's resource
	r.Get("/tenants/res", controller.GetManager().TenantsWithResource)
	r.Get("/tenants/res/page/{curPage}/size/{pageLen}", controller.GetManager().TenantsWithResource)
	r.Get("/tenants/query/{tenant_name}", controller.GetManager().TenantsQuery)
	r.Get("/tenants/{tenant_name}/res", controller.GetManager().TenantsGetByName)
	return r
}

func (v2 *V2) prometheusRouter() chi.Router {
	r := chi.NewRouter()
	return r
}

func (v2 *V2) appRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/export", controller.GetManager().ExportApp)
	r.Get("/export/{eventID}", controller.GetManager().ExportApp)

	r.Get("/download/{format}/{fileName}", controller.GetManager().Download)
	r.Post("/upload/{eventID}", controller.GetManager().NewUpload)
	r.Options("/upload/{eventID}", controller.GetManager().NewUpload)

	r.Post("/import/ids/{eventID}", controller.GetManager().ImportID)
	r.Get("/import/ids/{eventID}", controller.GetManager().ImportID)
	r.Delete("/import/ids/{eventID}", controller.GetManager().ImportID)

	r.Post("/import", controller.GetManager().ImportApp)
	r.Get("/import/{eventID}", controller.GetManager().ImportApp)
	r.Delete("/import/{eventID}", controller.GetManager().ImportApp)
	return r
}

func (v2 *V2) notificationEventRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", controller.GetNotificationEvents)
	r.Put("/{serviceAlias}", controller.HandleNotificationEvent)
	r.Get("/{hash}", controller.GetNotificationEvent)
	return r
}

func (v2 *V2) portRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/avail-port", controller.GetManager().GetAvailablePort)
	return r
}
