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

package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	api_model "github.com/gridworkz/kato/api/model"
	ctxutil "github.com/gridworkz/kato/api/util/ctx"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	validator "github.com/gridworkz/kato/util/govalidator"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/gridworkz/kato/worker/discover/model"
	"github.com/jinzhu/gorm"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
)

//StartService StartService
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/start  v2 startService
//
// Start the app
//
// start service
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) StartService(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)

	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	sEvent := r.Context().Value(ctxutil.ContextKey("event")).(*dbmodel.ServiceEvent)
	if service.Kind != "third_party" {
		if err := handler.CheckTenantResource(r.Context(), tenant, service.Replicas*service.ContainerMemory); err != nil {
			httputil.ReturnResNotEnough(r, w, sEvent.EventID, err.Error())
			return
		}
	}

	startStopStruct := &api_model.StartStopStruct{
		TenantID:  tenantID,
		ServiceID: serviceID,
		EventID:   sEvent.EventID,
		TaskType:  "start",
	}
	if err := handler.GetServiceManager().StartStopService(startStopStruct); err != nil {
		httputil.ReturnError(r, w, 500, "get service info error.")
		return
	}
	httputil.ReturnSuccess(r, w, sEvent)
}

//StopService StopService
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/stop v2 stopService
//
// Close app
//
// stop service
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) StopService(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	sEvent := r.Context().Value(ctxutil.ContextKey("event")).(*dbmodel.ServiceEvent)
	//save event
	defer event.CloseManager()
	startStopStruct := &api_model.StartStopStruct{
		TenantID:  tenantID,
		ServiceID: serviceID,
		EventID:   sEvent.EventID,
		TaskType:  "stop",
	}
	if err := handler.GetServiceManager().StartStopService(startStopStruct); err != nil {
		httputil.ReturnError(r, w, 500, "get service info error.")
		return
	}
	httputil.ReturnSuccess(r, w, sEvent)
}

//RestartService RestartService
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/restart v2 restartService
//
// Restart application
//
// restart service
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) RestartService(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	sEvent := r.Context().Value(ctxutil.ContextKey("event")).(*dbmodel.ServiceEvent)
	defer event.CloseManager()
	startStopStruct := &api_model.StartStopStruct{
		TenantID:  tenantID,
		ServiceID: serviceID,
		EventID:   sEvent.EventID,
		TaskType:  "restart",
	}

	curStatus := t.StatusCli.GetStatus(serviceID)
	if curStatus == "closed" {
		startStopStruct.TaskType = "start"
	}

	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	if err := handler.CheckTenantResource(r.Context(), tenant, service.Replicas*service.ContainerMemory); err != nil {
		httputil.ReturnResNotEnough(r, w, sEvent.EventID, err.Error())
		return
	}

	if err := handler.GetServiceManager().StartStopService(startStopStruct); err != nil {
		httputil.ReturnError(r, w, 500, "get service info error.")
		return
	}
	httputil.ReturnSuccess(r, w, sEvent)
}

//VerticalService VerticalService
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/vertical v2 verticalService
//
// Apply vertical scaling
//
// service vertical
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) VerticalService(w http.ResponseWriter, r *http.Request) {
	rules := validator.MapData{
		"container_cpu":    []string{"required"},
		"container_memory": []string{"required"},
	}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rules, nil)
	if !ok {
		return
	}
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	sEvent := r.Context().Value(ctxutil.ContextKey("event")).(*dbmodel.ServiceEvent)
	var cpuSet, gpuSet, memorySet *int
	if cpu, ok := data["container_cpu"].(float64); ok {
		cpuInt := int(cpu)
		cpuSet = &cpuInt
	}
	if memory, ok := data["container_memory"].(float64); ok {
		memoryInt := int(memory)
		memorySet = &memoryInt
	}
	if gpu, ok := data["container_gpu"].(float64); ok {
		gpuInt := int(gpu)
		gpuSet = &gpuInt
	}
	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	if memorySet != nil {
		if err := handler.CheckTenantResource(r.Context(), tenant, service.Replicas*(*memorySet)); err != nil {
			httputil.ReturnResNotEnough(r, w, sEvent.EventID, err.Error())
			return
		}
	}
	verticalTask := &model.VerticalScalingTaskBody{
		TenantID:        tenantID,
		ServiceID:       serviceID,
		EventID:         sEvent.EventID,
		ContainerCPU:    cpuSet,
		ContainerMemory: memorySet,
		ContainerGPU:    gpuSet,
	}
	if err := handler.GetServiceManager().ServiceVertical(r.Context(), verticalTask); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("service vertical error. %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, sEvent)
}

//HorizontalService HorizontalService
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/horizontal v2 horizontalService
//
// Application level scaling
//
// service horizontal
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) HorizontalService(w http.ResponseWriter, r *http.Request) {
	rules := validator.MapData{
		"node_num": []string{"required"},
	}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rules, nil)
	if !ok {
		return
	}
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	sEvent := r.Context().Value(ctxutil.ContextKey("event")).(*dbmodel.ServiceEvent)
	replicas := int32(data["node_num"].(float64))

	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	if err := handler.CheckTenantResource(r.Context(), tenant, service.ContainerMemory*int(replicas)); err != nil {
		httputil.ReturnResNotEnough(r, w, sEvent.EventID, err.Error())
		return
	}

	horizontalTask := &model.HorizontalScalingTaskBody{
		TenantID:  tenantID,
		ServiceID: serviceID,
		EventID:   sEvent.EventID,
		Username:  sEvent.UserName,
		Replicas:  replicas,
	}

	if err := handler.GetServiceManager().ServiceHorizontal(horizontalTask); err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}
	httputil.ReturnSuccess(r, w, sEvent)
}

//BuildService BuildService
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/build v2 serviceBuild
//
// Application build
//
// service build
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) BuildService(w http.ResponseWriter, r *http.Request) {
	var build api_model.ComponentBuildReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &build, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantName := r.Context().Value(ctxutil.ContextKey("tenant_name")).(string)
	build.TenantName = tenantName
	build.EventID = r.Context().Value(ctxutil.ContextKey("event_id")).(string)
	if build.ServiceID != serviceID {
		httputil.ReturnError(r, w, 400, "build service id is failure")
		return
	}

	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	if err := handler.CheckTenantResource(r.Context(), tenant, service.Replicas*service.ContainerMemory); err != nil {
		httputil.ReturnResNotEnough(r, w, build.EventID, err.Error())
		return
	}

	res, err := handler.GetOperationHandler().Build(&build)
	if err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}
	httputil.ReturnSuccess(r, w, res)
}

//BuildList BuildList
func (t *TenantStruct) BuildList(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)

	resp, err := handler.GetServiceManager().ListVersionInfo(serviceID)

	if err != nil {
		logrus.Error("get version info error", err.Error())
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get version info erro, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, resp)
}

//BuildVersionIsExist -
func (t *TenantStruct) BuildVersionIsExist(w http.ResponseWriter, r *http.Request) {
	statusMap := make(map[string]bool)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	buildVersion := chi.URLParam(r, "build_version")
	_, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(buildVersion, serviceID)
	if err != nil && err != gorm.ErrRecordNotFound {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get build version status erro, %v", err))
		return
	}
	if err == gorm.ErrRecordNotFound {
		statusMap["status"] = false
	} else {
		statusMap["status"] = true
	}
	httputil.ReturnSuccess(r, w, statusMap)

}

//DeleteBuildVersion -
func (t *TenantStruct) DeleteBuildVersion(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	buildVersion := chi.URLParam(r, "build_version")
	val, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(buildVersion, serviceID)
	if err != nil && err != gorm.ErrRecordNotFound {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("delete build version erro, %v", err))
		return
	}
	if err == gorm.ErrRecordNotFound {

	} else {
		if val.DeliveredType == "slug" && val.FinalStatus == "success" {
			if err := os.Remove(val.DeliveredPath); err != nil {
				httputil.ReturnError(r, w, 500, fmt.Sprintf("delete build version erro, %v", err))
				return

			}
			if err := db.GetManager().VersionInfoDao().DeleteVersionInfo(val); err != nil {
				httputil.ReturnError(r, w, 500, fmt.Sprintf("delete build version erro, %v", err))
				return

			}
		}
		if val.FinalStatus == "failure" {
			if err := db.GetManager().VersionInfoDao().DeleteVersionInfo(val); err != nil {
				httputil.ReturnError(r, w, 500, fmt.Sprintf("delete build version erro, %v", err))
				return
			}
		}
		if val.DeliveredType == "image" {
			if err := db.GetManager().VersionInfoDao().DeleteVersionInfo(val); err != nil {
				httputil.ReturnError(r, w, 500, fmt.Sprintf("delete build version erro, %v", err))
				return
			}
		}
	}
	httputil.ReturnSuccess(r, w, nil)

}

//UpdateBuildVersion -
func (t *TenantStruct) UpdateBuildVersion(w http.ResponseWriter, r *http.Request) {
	var build api_model.UpdateBuildVersionReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &build, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	buildVersion := chi.URLParam(r, "build_version")
	versionInfo, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(buildVersion, serviceID)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("update build version info error, %v", err))
		return
	}
	versionInfo.PlanVersion = build.PlanVersion
	err = db.GetManager().VersionInfoDao().UpdateModel(versionInfo)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("update build version info error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//BuildVersionInfo -
func (t *TenantStruct) BuildVersionInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		t.DeleteBuildVersion(w, r)
	case "GET":
		t.BuildVersionIsExist(w, r)
	case "PUT":
		t.UpdateBuildVersion(w, r)
	}

}

//GetDeployVersion GetDeployVersion by service
func (t *TenantStruct) GetDeployVersion(w http.ResponseWriter, r *http.Request) {
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	version, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(service.DeployVersion, service.ServiceID)
	if err != nil && err != gorm.ErrRecordNotFound {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get build version status erro, %v", err))
		return
	}
	if err == gorm.ErrRecordNotFound {
		httputil.ReturnError(r, w, 404, "build version do not exist")
		return
	}
	httputil.ReturnSuccess(r, w, version)
}

//GetManyDeployVersion GetDeployVersion by some service id
func (t *TenantStruct) GetManyDeployVersion(w http.ResponseWriter, r *http.Request) {
	rules := validator.MapData{
		"service_ids": []string{"required"},
	}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rules, nil)
	if !ok {
		return
	}
	serviceIDs, ok := data["service_ids"].([]interface{})
	if !ok {
		httputil.ReturnError(r, w, 400, "service ids must be a array")
		return
	}
	var list []string
	for _, s := range serviceIDs {
		list = append(list, s.(string))
	}
	services, err := db.GetManager().TenantServiceDao().GetServiceByIDs(list)
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	var versionList []*dbmodel.VersionInfo
	for _, service := range services {
		version, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(service.DeployVersion, service.ServiceID)
		if err != nil && err != gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 500, fmt.Sprintf("get build version status erro, %v", err))
			return
		}
		versionList = append(versionList, version)
	}
	httputil.ReturnSuccess(r, w, versionList)
}

//DeployService DeployService
func (t *TenantStruct) DeployService(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("trans deploy service")
	w.Write([]byte("deploy service"))
}

//UpgradeService UpgradeService
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/upgrade v2 upgradeService
//
// Upgrade application
//
// upgrade service
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) UpgradeService(w http.ResponseWriter, r *http.Request) {
	var upgradeRequest api_model.ComponentUpgradeReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &upgradeRequest, nil)
	if !ok {
		logrus.Errorf("start operation validate request body failure")
		return
	}
	upgradeRequest.EventID = r.Context().Value(ctxutil.ContextKey("event_id")).(string)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	if upgradeRequest.ServiceID != serviceID {
		httputil.ReturnError(r, w, 400, "upgrade service id failure")
		return
	}

	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	if service.Kind != "third_party" {
		if err := handler.CheckTenantResource(r.Context(), tenant, service.Replicas*service.ContainerMemory); err != nil {
			httputil.ReturnResNotEnough(r, w, upgradeRequest.EventID, err.Error())
			return
		}
	}

	res, err := handler.GetOperationHandler().Upgrade(&upgradeRequest)
	if err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}
	httputil.ReturnSuccess(r, w, res)
}

//CheckCode CheckCode
// swagger:operation POST /v2/tenants/{tenant_name}/code-check v2 checkCode
//
// Application code detection
//
// check  code
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) CheckCode(w http.ResponseWriter, r *http.Request) {

	var ccs api_model.CheckCodeStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ccs.Body, nil)
	if !ok {
		return
	}
	if ccs.Body.TenantID == "" {
		tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
		ccs.Body.TenantID = tenantID
	}
	ccs.Body.Action = "code_check"
	if err := handler.GetServiceManager().CodeCheck(&ccs); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("task code check error,%v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//RollBack RollBack
// swagger:operation Post /v2/tenants/{tenant_name}/services/{service_alias}/rollback v2 rollback
//
// Application version rollback
//
// service rollback
//
// ---
// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/xml
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) RollBack(w http.ResponseWriter, r *http.Request) {
	var rollbackRequest api_model.RollbackInfoRequestStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &rollbackRequest, nil)
	if !ok {
		logrus.Errorf("start operation validate request body failure")
		return
	}
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	if rollbackRequest.ServiceID != serviceID {
		httputil.ReturnError(r, w, 400, "rollback service id failure")
		return
	}
	rollbackRequest.EventID = r.Context().Value(ctxutil.ContextKey("event_id")).(string)

	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)
	service := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	if err := handler.CheckTenantResource(r.Context(), tenant, service.Replicas*service.ContainerMemory); err != nil {
		httputil.ReturnResNotEnough(r, w, rollbackRequest.EventID, err.Error())
		return
	}

	re := handler.GetOperationHandler().RollBack(rollbackRequest)
	httputil.ReturnSuccess(r, w, re)
}

type limitMemory struct {
	LimitMemory int `json:"limit_memory"`
}

//LimitTenantMemory -
func (t *TenantStruct) LimitTenantMemory(w http.ResponseWriter, r *http.Request) {
	var lm limitMemory
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	err = ffjson.Unmarshal(body, &lm)
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}

	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	tenant, err := db.GetManager().TenantDao().GetTenantByUUID(tenantID)
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}
	tenant.LimitMemory = lm.LimitMemory
	if err := db.GetManager().TenantDao().UpdateModel(tenant); err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
	}
	httputil.ReturnSuccess(r, w, "success!")

}

//SourcesInfo -
type SourcesInfo struct {
	TenantID        string `json:"tenant_id"`
	AvailableMemory int    `json:"available_memory"`
	Status          bool   `json:"status"`
	MemTotal        int    `json:"mem_total"`
	MemUsed         int    `json:"mem_used"`
	CPUTotal        int    `json:"cpu_total"`
	CPUUsed         int    `json:"cpu_used"`
}

//TenantResourcesStatus tenant resources status
func (t *TenantStruct) TenantResourcesStatus(w http.ResponseWriter, r *http.Request) {

	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	tenant, err := db.GetManager().TenantDao().GetTenantByUUID(tenantID)
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}
	//11ms
	services, err := handler.GetServiceManager().GetService(tenant.UUID)
	if err != nil {
		msg := httputil.ResponseBody{
			Msg: fmt.Sprintf("get service error, %v", err),
		}
		httputil.Return(r, w, 500, msg)
		return
	}

	statsInfo, _ := handler.GetTenantManager().StatsMemCPU(services)

	if tenant.LimitMemory == 0 {
		sourcesInfo := SourcesInfo{
			TenantID:        tenantID,
			AvailableMemory: 0,
			Status:          true,
			MemTotal:        tenant.LimitMemory,
			MemUsed:         statsInfo.MEM,
			CPUTotal:        0,
			CPUUsed:         statsInfo.CPU,
		}
		httputil.ReturnSuccess(r, w, sourcesInfo)
		return
	}
	if statsInfo.MEM >= tenant.LimitMemory {
		sourcesInfo := SourcesInfo{
			TenantID:        tenantID,
			AvailableMemory: tenant.LimitMemory - statsInfo.MEM,
			Status:          false,
			MemTotal:        tenant.LimitMemory,
			MemUsed:         statsInfo.MEM,
			CPUTotal:        tenant.LimitMemory / 4,
			CPUUsed:         statsInfo.CPU,
		}
		httputil.ReturnSuccess(r, w, sourcesInfo)
	} else {
		sourcesInfo := SourcesInfo{
			TenantID:        tenantID,
			AvailableMemory: tenant.LimitMemory - statsInfo.MEM,
			Status:          true,
			MemTotal:        tenant.LimitMemory,
			MemUsed:         statsInfo.MEM,
			CPUTotal:        tenant.LimitMemory / 4,
			CPUUsed:         statsInfo.CPU,
		}
		httputil.ReturnSuccess(r, w, sourcesInfo)
	}
}

//GetServiceDeployInfo get service deploy info
func GetServiceDeployInfo(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	info, err := handler.GetServiceManager().GetServiceDeployInfo(tenantID, serviceID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, info)
}

// Log -
func (t *TenantStruct) Log(w http.ResponseWriter, r *http.Request) {
	component := r.Context().Value(ctxutil.ContextKey("service")).(*dbmodel.TenantServices)
	podName := r.URL.Query().Get("podName")
	containerName := r.URL.Query().Get("containerName")
	follow, _ := strconv.ParseBool(r.URL.Query().Get("follow"))

	err := handler.GetServiceManager().Log(w, r, component, podName, containerName, follow)
	if err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}
}
