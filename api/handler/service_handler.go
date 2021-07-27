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

package handler

import (
	"context"
	"github.com/jinzhu/gorm"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/builder/exector"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/worker/discover/model"
	"github.com/gridworkz/kato/worker/server/pb"
)

//ServiceHandler service handler
type ServiceHandler interface {
	ServiceBuild(tenantID, serviceID string, r *api_model.BuildServiceStruct) error
	AddLabel(l *api_model.LabelsStruct, serviceID string) error
	DeleteLabel(l *api_model.LabelsStruct, serviceID string) error
	UpdateLabel(l *api_model.LabelsStruct, serviceID string) error
	StartStopService(s *api_model.StartStopStruct) error
	ServiceVertical(ctx context.Context, v *model.VerticalScalingTaskBody) error
	ServiceHorizontal(h *model.HorizontalScalingTaskBody) error
	ServiceUpgrade(r *model.RollingUpgradeTaskBody) error
	ServiceCreate(ts *api_model.ServiceStruct) error
	ServiceUpdate(sc map[string]interface{}) error
	LanguageSet(langS *api_model.LanguageSet) error
	GetService(tenantID string) ([]*dbmodel.TenantServices, error)
	GetServicesByAppID(appID string, page, pageSize int) (*api_model.ListServiceResponse, error)
	GetPagedTenantRes(offset, len int) ([]*api_model.TenantResource, int, error)
	GetTenantRes(uuid string) (*api_model.TenantResource, error)
	CodeCheck(c *api_model.CheckCodeStruct) error
	ServiceDepend(action string, ds *api_model.DependService) error
	EnvAttr(action string, at *dbmodel.TenantServiceEnvVar) error
	PortVar(action string, tenantID, serviceID string, vp *api_model.ServicePorts, oldPort int) error
	CreatePorts(tenantID, serviceID string, vps *api_model.ServicePorts) error
	PortOuter(tenantName, serviceID string, containerPort int, servicePort *api_model.ServicePortInnerOrOuter) (*dbmodel.TenantServiceLBMappingPort, string, error)
	PortInner(tenantName, serviceID, operation string, port int) error
	VolumnVar(avs *dbmodel.TenantServiceVolume, tenantID, fileContent, action string) *util.APIHandleError
	UpdVolume(sid string, req *api_model.UpdVolumeReq) error
	VolumeDependency(tsr *dbmodel.TenantServiceMountRelation, action string) *util.APIHandleError
	GetDepVolumes(serviceID string) ([]*dbmodel.TenantServiceMountRelation, *util.APIHandleError)
	GetVolumes(serviceID string) ([]*api_model.VolumeWithStatusStruct, *util.APIHandleError)
	ServiceProbe(tsp *dbmodel.TenantServiceProbe, action string) error
	RollBack(rs *api_model.RollbackStruct) error
	GetStatus(serviceID string) (*api_model.StatusList, error)
	GetServicesStatus(tenantID string, services []string) []map[string]interface{}
	GetEnterpriseRunningServices(enterpriseID string) ([]string, *util.APIHandleError)
	CreateTenant(*dbmodel.Tenants) error
	CreateTenandIDAndName(eid string) (string, string, error)
	GetPods(serviceID string) (*K8sPodInfos, error)
	GetMultiServicePods(serviceIDs []string) (*K8sPodInfos, error)
	GetComponentPodNums(ctx context.Context, componentIDs []string) (map[string]int32, error)
	TransServieToDelete(ctx context.Context, tenantID, serviceID string) error
	TenantServiceDeletePluginRelation(tenantID, serviceID, pluginID string) *util.APIHandleError
	GetTenantServicePluginRelation(serviceID string) ([]*dbmodel.TenantServicePluginRelation, *util.APIHandleError)
	SetTenantServicePluginRelation(tenantID, serviceID string, pss *api_model.PluginSetStruct) (*dbmodel.TenantServicePluginRelation, *util.APIHandleError)
	UpdateTenantServicePluginRelation(serviceID string, pss *api_model.PluginSetStruct) (*dbmodel.TenantServicePluginRelation, *util.APIHandleError)
	UpdateVersionEnv(uve *api_model.SetVersionEnv) *util.APIHandleError
	DeletePluginConfig(serviceID, pluginID string) *util.APIHandleError
	ServiceCheck(*api_model.ServiceCheckStruct) (string, string, *util.APIHandleError)
	GetServiceCheckInfo(uuid string) (*exector.ServiceCheckResult, *util.APIHandleError)
	GetServiceDeployInfo(tenantID, serviceID string) (*pb.DeployInfo, *util.APIHandleError)
	ListVersionInfo(serviceID string) (*api_model.BuildListRespVO, error)

	AddAutoscalerRule(req *api_model.AutoscalerRuleReq) error
	UpdAutoscalerRule(req *api_model.AutoscalerRuleReq) error
	ListScalingRecords(serviceID string, page, pageSize int) ([]*dbmodel.TenantServiceScalingRecords, int, error)

	UpdateServiceMonitor(tenantID, serviceID, name string, update api_model.UpdateServiceMonitorRequestStruct) (*dbmodel.TenantServiceMonitor, error)
	DeleteServiceMonitor(tenantID, serviceID, name string) (*dbmodel.TenantServiceMonitor, error)
	AddServiceMonitor(tenantID, serviceID string, add api_model.AddServiceMonitorRequestStruct) (*dbmodel.TenantServiceMonitor, error)

	SyncComponentBase(tx *gorm.DB, app *dbmodel.Application, components []*api_model.Component) error
	SyncComponentMonitors(tx *gorm.DB,app *dbmodel.Application, components []*api_model.Component) error
	SyncComponentPorts(tx *gorm.DB, app *dbmodel.Application, components []*api_model.Component) error
	SyncComponentRelations(tx *gorm.DB, app *dbmodel.Application, components []*api_model.Component) error
	SyncComponentEnvs(tx *gorm.DB, app *dbmodel.Application, components []*api_model.Component) error
	SyncComponentVolumeRels(tx *gorm.DB, app *dbmodel.Application, components []*api_model.Component) error
	SyncComponentVolumes(tx *gorm.DB,  components []*api_model.Component) error
	SyncComponentConfigFiles(tx *gorm.DB,  components []*api_model.Component) error
	SyncComponentProbes(tx *gorm.DB,  components []*api_model.Component) error
	SyncComponentLabels(tx *gorm.DB,  components []*api_model.Component) error
	SyncComponentPlugins(tx *gorm.DB, app *dbmodel.Application, components []*api_model.Component) error
	SyncComponentScaleRules(tx *gorm.DB,  components []*api_model.Component) error
	SyncComponentEndpoints(tx *gorm.DB, components []*api_model.Component) error
}
