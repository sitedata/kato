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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util/bcode"
	"github.com/gridworkz/kato/cmd"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/errors"
	dbmodel "github.com/gridworkz/kato/db/model"
	mqclient "github.com/gridworkz/kato/mq/client"
	validation "github.com/gridworkz/kato/util/endpoint"
	"github.com/gridworkz/kato/util/fuzzy"
	validator "github.com/gridworkz/kato/util/govalidator"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/gridworkz/kato/worker/client"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//V2Routes v2Routes
type V2Routes struct {
	ClusterController
	TenantStruct
	EventLogStruct
	AppStruct
	GatewayStruct
	ThirdPartyServiceController
	LabelController
	AppRestoreController
	PodController
	ApplicationController
}

//Show test
func (v2 *V2Routes) Show(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/show v2 getApiVersion
	//
	// Display the current api version information
	//
	// show api version
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format
	w.Write([]byte(cmd.GetVersion()))
}

//Health show health status
func (v2 *V2Routes) Health(w http.ResponseWriter, r *http.Request) {
	httputil.ReturnSuccess(r, w, map[string]string{"status": "health", "info": "api service health"})
}

//AlertManagerWebHook -
func (v2 *V2Routes) AlertManagerWebHook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=======>webhook")
	in, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		httputil.ReturnError(r, w, 400, "")
		return
	}
	fmt.Println("=====>body")
	fmt.Println(string(in))
	httputil.ReturnSuccess(r, w, "")

}

//Version -
func (v2 *V2Routes) Version(w http.ResponseWriter, r *http.Request) {
	httputil.ReturnSuccess(r, w, map[string]string{"version": os.Getenv("RELEASE_DESC")})
}

//TenantStruct tenant struct
type TenantStruct struct {
	StatusCli *client.AppRuntimeSyncClient
	MQClient  mqclient.MQClient
}

//AllTenantResources GetResources
func (t *TenantStruct) AllTenantResources(w http.ResponseWriter, r *http.Request) {
	tenants, err := handler.GetTenantManager().GetTenants("")
	if err != nil {
		msg := httputil.ResponseBody{
			Msg: fmt.Sprintf("get tenant error, %v", err),
		}
		httputil.Return(r, w, 500, msg)
	}
	ts := &api_model.TotalStatsInfo{}
	for _, tenant := range tenants {
		services, err := handler.GetServiceManager().GetService(tenant.UUID)
		if err != nil {
			msg := httputil.ResponseBody{
				Msg: fmt.Sprintf("get service error, %v", err),
			}
			httputil.Return(r, w, 500, msg)
		}
		statsInfo, _ := handler.GetTenantManager().StatsMemCPU(services)
		statsInfo.UUID = tenant.UUID
		ts.Data = append(ts.Data, statsInfo)
	}
	httputil.ReturnSuccess(r, w, ts.Data)
}

//TenantResources TenantResources
func (t *TenantStruct) TenantResources(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/resources/tenants v2 tenantResources
	//
	// Tenant resource usage
	//
	// get tenant resources
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format
	var tr api_model.TenantResources
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tr.Body, nil)
	if !ok {
		return
	}

	rep, err := handler.GetTenantManager().GetTenantsResources(&tr)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get resources error, %v", err))
		return
	}
	var re []map[string]interface{}
	for _, v := range rep {
		if v != nil {
			re = append(re, v)
		}
	}
	httputil.ReturnSuccess(r, w, re)
	return
}

//ServiceResources ServiceResources
func (t *TenantStruct) ServiceResources(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/resources/services v2 serviceResources
	//
	// Application resource usage
	//
	// get service resources
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format
	var tr api_model.ServicesResources
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tr.Body, nil)
	if !ok {
		return
	}
	rep, err := handler.GetTenantManager().GetServicesResources(&tr)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get resources error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, rep)
	return
}

//TenantsQuery TenantsQuery
func (t *TenantStruct) TenantsQuery(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/query/{tenant_name} v2 tenants
	//
	// Tenant get resource list
	//
	// get tenant resources
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// parameters:
	// - name: tenant_name
	//   in: path
	//   description: '123'
	//   required: true
	//   type: string
	//   format: string
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format

	tenantName := strings.TrimSpace(chi.URLParam(r, "tenant_name"))

	rep, err := handler.GetTenantManager().GetTenantsName()
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get tenants names error, %v", err))
		return
	}

	result := fuzzy.Find(tenantName, rep) // [cartwheel wheel]
	httputil.ReturnSuccess(r, w, result)
	return
}

//TenantsGetByName TenantsGetByName
func (t *TenantStruct) TenantsGetByName(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/res v2 tenants
	//
	// Tenants bring a single resource
	//
	// get tenant resources
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// parameters:
	// - name: tenant_name
	//   in: path
	//   description: '123'
	//   required: true
	//   type: string
	//   format: string
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format

	tenantName := strings.TrimSpace(chi.URLParam(r, "tenant_name"))

	v, err := handler.GetTenantManager().GetTenantsByName(tenantName)
	if err != nil {
		httputil.ReturnError(r, w, 404, fmt.Sprintf("get tenants names error, %v", err))
		return
	}
	logrus.Infof("query tenant from db by name %s ,got %v", tenantName, v)

	tenantServiceRes, err := handler.GetServiceManager().GetTenantRes(v.UUID)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get tenants service total resources  error, %v", err))
		return
	}
	tenantServiceRes.UUID = v.UUID
	tenantServiceRes.Name = v.Name
	tenantServiceRes.EID = v.EID

	httputil.ReturnSuccess(r, w, tenantServiceRes)
	return
}

//TenantsWithResource TenantsWithResource
func (t *TenantStruct) TenantsWithResource(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/resources/tenants/res/page/{curPage}/size/{pageLen} v2 PagedTenantResList
	//
	// Tenant bring resource list
	//
	// get paged tenant resources
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// parameters:
	// - name: curPage
	//   in: path
	//   description: '123'
	//   required: true
	//   type: string
	//   format: string
	// - name: pageLen
	//   in: path
	//   description: '25'
	//   required: true
	//   type: string
	//   format: string
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format

	pageLenStr := strings.TrimSpace(chi.URLParam(r, "pageLen"))
	curPageStr := strings.TrimSpace(chi.URLParam(r, "curPage"))

	pageLen, err := strconv.Atoi(pageLenStr)
	if err != nil {
		httputil.ReturnError(r, w, 400, fmt.Sprintf("bad request, %v", err))
		return
	}
	curPage, err := strconv.Atoi(curPageStr)
	if err != nil {
		httputil.ReturnError(r, w, 400, fmt.Sprintf("bad request, %v", err))
		return
	}
	resource, count, err := handler.GetServiceManager().GetPagedTenantRes((curPage-1)*pageLen, pageLen)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get tenants  error, %v", err))
		return
	}
	var ret api_model.PagedTenantResList
	ret.List = resource
	ret.Length = count
	httputil.ReturnSuccess(r, w, ret)
	return
}

//SumTenants Count the number of tenants
func (t *TenantStruct) SumTenants(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/resources/tenants/sum v2 sumTenants
	//
	// Get the number of tenants
	//
	// get tenant resources
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//       "$ref": "#/responses/commandResponse"
	//     description: Unified return format

	s, err := handler.GetTenantManager().TenantsSum()
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("sum tenants error, %v", err))
		return
	}
	rc := make(map[string]int)
	rc["num"] = s
	httputil.ReturnSuccess(r, w, rc)
}

//Tenant one tenant controller
func (t *TenantStruct) Tenant(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t.GetTenant(w, r)
	case "DELETE":
		t.DeleteTenant(w, r)
	case "PUT":
		t.UpdateTenant(w, r)
	}
}

//Tenants Tenant
func (t *TenantStruct) Tenants(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		t.AddTenant(w, r)
	case "GET":
		t.GetTenants(w, r)
	}
}

//AddTenant AddTenant
func (t *TenantStruct) AddTenant(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants v2 addTenant
	//
	// Add tenant information
	//
	// add tenant
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
	var ts api_model.AddTenantStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ts.Body, nil)
	if !ok {
		return
	}
	var dbts dbmodel.Tenants
	if ts.Body.Eid != "" {
		//New interface
		//TODO:Generate tenant_id and tenant_name
		id, name, errN := handler.GetServiceManager().CreateTenandIDAndName(ts.Body.Eid)
		if errN != nil {
			httputil.ReturnError(r, w, 500, fmt.Sprintf("create tenant error, %v", errN))
			return
		}
		dbts.EID = ts.Body.Eid
		if ts.Body.TenantName == "" {
			dbts.Name = name
		} else {
			dbts.Name = ts.Body.TenantName
			name = ts.Body.TenantName
		}
		if ts.Body.TenantID == "" {
			dbts.UUID = id
		} else {
			dbts.UUID = ts.Body.TenantID
			id = ts.Body.TenantID
		}
		dbts.LimitMemory = ts.Body.LimitMemory
		if err := handler.GetServiceManager().CreateTenant(&dbts); err != nil {
			if strings.HasSuffix(err.Error(), "is exist") {
				httputil.ReturnError(r, w, 400, err.Error())
				return
			}
			httputil.ReturnError(r, w, 500, fmt.Sprintf("create tenant error, %v", err))
			return
		}
		rc := make(map[string]string)
		rc["tenant_id"] = id
		rc["tenang_name"] = name
		rc["eid"] = ts.Body.Eid
		httputil.ReturnSuccess(r, w, rc)
		return
	}
	if ts.Body.TenantID != "" && ts.Body.TenantName != "" {
		//Compatible with old interfaces
		dbts.Name = ts.Body.TenantName
		dbts.UUID = ts.Body.TenantID
		if err := handler.GetServiceManager().CreateTenant(&dbts); err != nil {
			if strings.HasSuffix(err.Error(), "is exist") {
				httputil.ReturnError(r, w, 400, err.Error())
				return
			}
			httputil.ReturnError(r, w, 500, fmt.Sprintf("create tenant error, %v", err))
			return
		}
		httputil.ReturnSuccess(r, w, nil)
		return
	}
	if ts.Body.Eid == "" && (ts.Body.TenantID == "" || ts.Body.TenantName == "") {
		httputil.ReturnError(r, w, 400, "args error")
		return
	}
	httputil.ReturnError(r, w, 400, "args error, need eid or tenatn_id / tenant_name")
	return
}

//GetTenants GetTenants
func (t *TenantStruct) GetTenants(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants v2 getTenants
	//
	// Get all tenant information
	//
	// get tenant
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
	value := r.FormValue("eid")
	page, _ := strconv.Atoi(r.FormValue("page"))
	if page == 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.FormValue("pageSize"))
	if pageSize == 0 {
		pageSize = 10
	}
	queryName := r.FormValue("query")
	var tenants []*dbmodel.Tenants
	var err error
	if len(value) == 0 {
		tenants, err = handler.GetTenantManager().GetTenants(queryName)
		if err != nil {
			httputil.ReturnError(r, w, 500, "get tenant error")
			return
		}
	} else {
		tenants, err = handler.GetTenantManager().GetTenantsByEid(value, queryName)
		if err != nil {
			httputil.ReturnError(r, w, 500, "get tenant error")
			return
		}
	}
	list := handler.GetTenantManager().BindTenantsResource(tenants)
	re := list.Paging(page, pageSize)
	httputil.ReturnSuccess(r, w, re)
}

//DeleteTenant
func (t *TenantStruct) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)

	if err := handler.GetTenantManager().DeleteTenant(tenantID); err != nil {
		if err == handler.ErrTenantStillHasServices || err == handler.ErrTenantStillHasPlugins {
			httputil.ReturnError(r, w, 400, err.Error())
			return
		}
		if err == gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 404, err.Error())
			return
		}

		httputil.ReturnError(r, w, 500, fmt.Sprintf("delete tenant: %v", err))
		return
	}

	httputil.ReturnSuccess(r, w, nil)
}

//UpdateTenant
//support update tenant limit memory
func (t *TenantStruct) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	var ts api_model.UpdateTenantStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ts.Body, nil)
	if !ok {
		return
	}
	tenant := r.Context().Value(middleware.ContextKey("tenant")).(*dbmodel.Tenants)
	tenant.LimitMemory = ts.Body.LimitMemory
	if err := handler.GetTenantManager().UpdateTenant(tenant); err != nil {
		httputil.ReturnError(r, w, 500, "update tenant error")
		return
	}
	httputil.ReturnSuccess(r, w, tenant)
}

//GetTenant get one tenant
func (t *TenantStruct) GetTenant(w http.ResponseWriter, r *http.Request) {
	tenant := r.Context().Value(middleware.ContextKey("tenant")).(*dbmodel.Tenants)
	list := handler.GetTenantManager().BindTenantsResource([]*dbmodel.Tenants{tenant})
	httputil.ReturnSuccess(r, w, list[0])
}

//ServicesCount Get all apps and status
func (t *TenantStruct) ServicesCount(w http.ResponseWriter, r *http.Request) {
	allStatus := t.StatusCli.GetAllStatus()
	var closed int
	var running int
	var abnormal int
	for _, v := range allStatus {
		switch v {
		case "closed":
			closed++
		case "running":
			running++
		case "abnormal":
			abnormal++
		}
	}
	serviceCount := map[string]int{"total": len(allStatus), "running": running, "closed": closed, "abnormal": abnormal}
	httputil.ReturnSuccess(r, w, serviceCount)
}

//ServicesInfo GetServiceInfo
func (t *TenantStruct) ServicesInfo(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/services v2 getServiceInfo
	//
	// Get all application information of the tenant
	//
	// get services info in tenant
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
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	services, err := handler.GetServiceManager().GetService(tenantID)
	if err != nil {
		httputil.ReturnError(r, w, 500, "get tenant services error")
		return
	}
	httputil.ReturnSuccess(r, w, services)
	return
}

//CreateService create Service
func (t *TenantStruct) CreateService(w http.ResponseWriter, r *http.Request) {
	var ss api_model.ServiceStruct
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &ss, nil) {
		return
	}

	// Check if the application ID exists
	if ss.AppID == "" {
		httputil.ReturnBcodeError(r, w, bcode.ErrCreateNeedCorrectAppID)
		return
	}
	_, err := handler.GetApplicationHandler().GetAppByID(ss.AppID)
	if err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}

	// clean etcd data(source check)
	handler.GetEtcdHandler().CleanServiceCheckData(ss.EtcdKey)

	values := url.Values{}
	if ss.Endpoints != nil && strings.TrimSpace(ss.Endpoints.Static) != "" {
		if strings.Contains(ss.Endpoints.Static, "127.0.0.1") {
			values["ip"] = []string{"The ip field is can't contains '127.0.0.1'"}
		}
	}
	if len(values) > 0 {
		httputil.ReturnValidationError(r, w, values)
		return
	}

	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	ss.TenantID = tenantID
	if err := handler.GetServiceManager().ServiceCreate(&ss); err != nil {
		if strings.Contains(err.Error(), "is exist in tenant") {
			httputil.ReturnError(r, w, 400, fmt.Sprintf("create service error, %v", err))
		}
		httputil.ReturnError(r, w, 500, fmt.Sprintf("create service error, %v", err))
		return
	}

	httputil.ReturnSuccess(r, w, nil)
}

//UpdateService create Service
func (t *TenantStruct) UpdateService(w http.ResponseWriter, r *http.Request) {
	// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias} v2 updateService
	//
	// App update
	//
	// update service
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
	// Currently provides three element modifications
	rules := validator.MapData{
		"container_cmd":    []string{},
		"image_name":       []string{},
		"container_memory": []string{},
		"service_name":     []string{},
		"extend_method":    []string{},
		"app_id":           []string{},
	}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rules, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	data["service_id"] = serviceID

	// Check if the application ID exists
	var appID string
	if data["app_id"] != nil && data["app_id"] != "" {
		appID = data["app_id"].(string)
		_, err := handler.GetApplicationHandler().GetAppByID(appID)
		if err != nil {
			httputil.ReturnBcodeError(r, w, err)
			return
		}
	}

	if err := handler.GetServiceManager().ServiceUpdate(data); err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//SetLanguage SetLanguage
func (t *TenantStruct) SetLanguage(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST  /v2/tenants/{tenant_name}/services/{service_alias}/language v2 setLanguage
	//
	// Set application language
	//
	// set language
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
	rules := validator.MapData{
		"language": []string{"required"},
	}
	langS := &api_model.LanguageSet{}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rules, nil)
	if !ok {
		return
	}
	langS.Language = data["language"].(string)
	langS.ServiceID = r.Context().Value(middleware.ContextKey("service_id")).(string)
	if err := handler.GetServiceManager().LanguageSet(langS); err != nil {
		httputil.ReturnError(r, w, 500, "set language error.")
		return
	}
	httputil.ReturnSuccess(r, w, nil)
	return
}

//StatusService StatusService
func (t *TenantStruct) StatusService(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/status v2 serviceStatus
	//
	// Get application status
	//
	// get service status
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
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	statusList, err := handler.GetServiceManager().GetStatus(serviceID)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get service list error,%v", err))
		return
	}
	httputil.ReturnSuccess(r, w, statusList)
	return
}

//PostStatusService PostStatusService
func (t *TenantStruct) PostStatusService(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("in status service serviceID")
}

//StatusServiceList service list status
func (t *TenantStruct) StatusServiceList(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants/{tenant_name}/services_status v2 serviceStatuslist
	//
	// Get application status
	//
	// get service statuslist
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
	var services api_model.StatusServiceListStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &services.Body, nil)
	if !ok {
		return
	}
	//logrus.Info(services.Body.ServiceIDs)
	serviceList := services.Body.ServiceIDs
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	info := handler.GetServiceManager().GetServicesStatus(tenantID, serviceList)

	httputil.ReturnSuccess(r, w, info)
}

//Label -
func (t *TenantStruct) Label(w http.ResponseWriter, r *http.Request) {
	var req api_model.LabelsStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}
	reqJSON, _ := json.Marshal(req)
	logrus.Debugf("Request is : %s", string(reqJSON))

	// verify request
	values := url.Values{}
	if req.Labels == nil || len(req.Labels) == 0 {
		values["labels"] = []string{"The labels field should have someting"}
	}
	for _, label := range req.Labels {
		if label.LabelKey == "" {
			values["label_key"] = []string{"The label_key field is required"}
		}
		if label.LabelValue == "" {
			values["label_value"] = []string{"The label_value field is required"}
		}
	}
	if len(values) != 0 {
		httputil.ReturnValidationError(r, w, values)
		return
	}

	switch r.Method {
	case "DELETE":
		t.DeleteLabel(w, r, &req)
	case "POST":
		t.AddLabel(w, r, &req)
	case "PUT":
		t.UpdateLabel(w, r, &req)
	}
}

// AddLabel adds label
func (t *TenantStruct) AddLabel(w http.ResponseWriter, r *http.Request, labels *api_model.LabelsStruct) {
	logrus.Debugf("add label")
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	if err := handler.GetServiceManager().AddLabel(labels, serviceID); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("add label error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

// DeleteLabel deletes labels
func (t *TenantStruct) DeleteLabel(w http.ResponseWriter, r *http.Request, labels *api_model.LabelsStruct) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	if err := handler.GetServiceManager().DeleteLabel(labels, serviceID); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("delete node label failure, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//UpdateLabel Update updates labels
func (t *TenantStruct) UpdateLabel(w http.ResponseWriter, r *http.Request, labels *api_model.LabelsStruct) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	if err := handler.GetServiceManager().UpdateLabel(labels, serviceID); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("error updating label: %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//StatusContainerID StatusContainerID
func (t *TenantStruct) StatusContainerID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("status container IDs list"))
}

//SingleServiceInfo SingleServiceInfo
func (t *TenantStruct) SingleServiceInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		t.DeleteSingleServiceInfo(w, r)
	case "GET":
		t.GetSingleServiceInfo(w, r)
	}
}

//GetSingleServiceInfo GetSingleServiceInfo
func (t *TenantStruct) GetSingleServiceInfo(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias} v2 getService
	//
	// Get application information
	//
	// get service info
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

	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	tenantName := r.Context().Value(middleware.ContextKey("tenant_name")).(string)
	serviceName := r.Context().Value(middleware.ContextKey("service_alias")).(string)
	result := make(map[string]string)
	result["tenantName"] = tenantName
	result["serviceAlias"] = serviceName
	result["tenantId"] = tenantID
	result["serviceId"] = serviceID
	httputil.ReturnSuccess(r, w, result)
}

//DeleteSingleServiceInfo DeleteService
// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias} v2 deleteService
//
// Delete app
//
// delete service
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
func (t *TenantStruct) DeleteSingleServiceInfo(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	var req api_model.EtcdCleanReq
	if httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil) {
		logrus.Debugf("delete service etcd keys : %+v", req.Keys)
		handler.GetEtcdHandler().CleanAllServiceData(req.Keys)
	}

	if err := handler.GetServiceManager().TransServieToDelete(tenantID, serviceID); err != nil {
		if err == handler.ErrServiceNotClosed {
			httputil.ReturnError(r, w, 400, fmt.Sprintf("Service must be closed"))
			return
		}
		httputil.ReturnError(r, w, 500, fmt.Sprintf("delete service error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//Dependency Dependency
func (t *TenantStruct) Dependency(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		t.DeleteDependency(w, r)
	case "POST":
		t.AddDependency(w, r)
	}
}

//AddDependency AddDependency
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/dependency v2 addDependency
//
// Increase application dependencies
//
// add dependency
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
func (t *TenantStruct) AddDependency(w http.ResponseWriter, r *http.Request) {
	rules := validator.MapData{
		"dep_service_id":   []string{"required"},
		"dep_service_type": []string{"required"},
		"dep_order":        []string{},
	}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rules, nil)
	if !ok {
		return
	}
	ds := &api_model.DependService{
		TenantID:       r.Context().Value(middleware.ContextKey("tenant_id")).(string),
		ServiceID:      r.Context().Value(middleware.ContextKey("service_id")).(string),
		DepServiceID:   data["dep_service_id"].(string),
		DepServiceType: data["dep_service_type"].(string),
	}
	if err := handler.GetServiceManager().ServiceDepend("add", ds); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("add dependency error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeleteDependency DeleteDependency
// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/dependency v2 deleteDependency
//
// Remove application dependencies
//
// delete dependency
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
func (t *TenantStruct) DeleteDependency(w http.ResponseWriter, r *http.Request) {
	logrus.Debugf("trans delete depend service ")
	rules := validator.MapData{
		"dep_service_id":   []string{"required"},
		"dep_service_type": []string{},
		"dep_order":        []string{},
	}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rules, nil)
	if !ok {
		return
	}
	ds := &api_model.DependService{
		TenantID:     r.Context().Value(middleware.ContextKey("tenant_id")).(string),
		ServiceID:    r.Context().Value(middleware.ContextKey("service_id")).(string),
		DepServiceID: data["dep_service_id"].(string),
	}
	if err := handler.GetServiceManager().ServiceDepend("delete", ds); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("delete dependency error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//Env Env
func (t *TenantStruct) Env(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		t.DeleteEnv(w, r)
	case "POST":
		t.AddEnv(w, r)
	case "PUT":
		t.UpdateEnv(w, r)
	}
}

//AddEnv AddEnv
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/env v2 addEnv
//
// Increase environment variables
//
// add env var
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
func (t *TenantStruct) AddEnv(w http.ResponseWriter, r *http.Request) {
	var envM api_model.AddTenantServiceEnvVar
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &envM, nil) {
		return
	}
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	var envD dbmodel.TenantServiceEnvVar
	envD.AttrName = envM.AttrName
	envD.AttrValue = envM.AttrValue
	envD.TenantID = tenantID
	envD.ServiceID = serviceID
	envD.ContainerPort = envM.ContainerPort
	envD.IsChange = envM.IsChange
	envD.Name = envM.Name
	envD.Scope = envM.Scope
	if err := handler.GetServiceManager().EnvAttr("add", &envD); err != nil {
		if err == errors.ErrRecordAlreadyExist {
			httputil.ReturnError(r, w, 400, fmt.Sprintf("%v", err))
			return
		}
		logrus.Errorf("Add env error, %v", err)
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Add env error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//UpdateEnv UpdateEnv
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/env v2 update Env
//
// Modify environment variables
//
// update env var
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
func (t *TenantStruct) UpdateEnv(w http.ResponseWriter, r *http.Request) {
	var envM api_model.AddTenantServiceEnvVar
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &envM, nil) {
		return
	}
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	var envD dbmodel.TenantServiceEnvVar
	envD.AttrName = envM.AttrName
	envD.AttrValue = envM.AttrValue
	envD.TenantID = tenantID
	envD.ServiceID = serviceID
	envD.ContainerPort = envM.ContainerPort
	envD.IsChange = envM.IsChange
	envD.Name = envM.Name
	envD.Scope = envM.Scope
	if err := handler.GetServiceManager().EnvAttr("update", &envD); err != nil {
		logrus.Errorf("update env error, %v", err)
		httputil.ReturnError(r, w, 500, fmt.Sprintf("update env error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeleteEnv DeleteEnv
// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/env v2 deleteEnv
//
// Delete environment variables
//
// delete env var
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
func (t *TenantStruct) DeleteEnv(w http.ResponseWriter, r *http.Request) {
	var envM api_model.DelTenantServiceEnvVar
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &envM, nil) {
		return
	}
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	envM.TenantID = tenantID
	envM.ServiceID = serviceID
	var envD dbmodel.TenantServiceEnvVar
	envD.AttrName = envM.AttrName
	envD.AttrValue = envM.AttrValue
	envD.TenantID = tenantID
	envD.ServiceID = serviceID
	envD.ContainerPort = envM.ContainerPort
	envD.IsChange = envM.IsChange
	envD.Name = envM.Name
	envD.Scope = envM.Scope
	if err := handler.GetServiceManager().EnvAttr("delete", &envD); err != nil {
		logrus.Errorf("delete env error, %v", err)
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			httputil.ReturnError(r, w, 404, "service port "+err.Error())
			return
		}
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Delete env error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//Ports Application Port Controller
func (t *TenantStruct) Ports(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		t.deletePortController(w, r)
	case "POST":
		t.addPortController(w, r)
	case "PUT":
		t.updatePortController(w, r)
	}
}

//PutPorts PortVar
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/ports v2 updatePort
//
// Update application port information (old)
//
// update port
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
func (t *TenantStruct) PutPorts(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	var ports api_model.ServicePorts
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ports, nil); !ok {
		return
	}
	if err := handler.GetServiceManager().PortVar("update", tenantID, serviceID, &ports, 0); err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//AddPortVar PortVar
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/ports v2 addPort
//
// Increase the application port, the internal and external options are disabled by default, and the corresponding interface needs to be enabled
//
// add port
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
func (t *TenantStruct) addPortController(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	var ports api_model.ServicePorts
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ports, nil); !ok {
		return
	}
	if err := handler.GetServiceManager().CreatePorts(tenantID, serviceID, &ports); err != nil {
		logrus.Errorf("add port error. %v", err)
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, ports.Port)
}

//UpdatePortVar PortVar
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/ports/{port} v2 updatePort
//
// Update application port information
//
// update port
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
func (t *TenantStruct) updatePortController(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	portStr := chi.URLParam(r, "port")
	oldPort, err := strconv.Atoi(portStr)
	if err != nil {
		httputil.ReturnError(r, w, 400, "port must be a number")
		return
	}
	var ports api_model.ServicePorts
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ports, nil); !ok {
		return
	}
	if err := handler.GetServiceManager().PortVar("update", tenantID, serviceID, &ports, oldPort); err != nil {
		logrus.Errorf("update port error. %v", err)
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeletePortVar PortVar
// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/ports/{port} v2 deletePort
//
// Delete port variable
//
// delete port
//
// ---
// Consumes:
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
func (t *TenantStruct) deletePortController(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	portStr := chi.URLParam(r, "port")
	oldPort, err := strconv.Atoi(portStr)
	if err != nil {
		httputil.ReturnError(r, w, 400, "port must be a number")
		return
	}
	var port = &api_model.TenantServicesPort{
		TenantID:      tenantID,
		ServiceID:     serviceID,
		ContainerPort: oldPort,
	}
	var ports api_model.ServicePorts
	ports.Port = append(ports.Port, port)
	if err := handler.GetServiceManager().PortVar("delete", tenantID, serviceID, &ports, oldPort); err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			httputil.ReturnError(r, w, 404, "port can not found")
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//PortOuterController Switch port external service
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/ports/{port}/outer v2 PortOuterController
//
// Switch ports for external services, the application will take effect automatically without restarting
//
// add port
//
// ---
// Consumes:
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
func (t *TenantStruct) PortOuterController(w http.ResponseWriter, r *http.Request) {
	var data api_model.ServicePortInnerOrOuter
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &(data.Body), nil) {
		return
	}

	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	service := r.Context().Value(middleware.ContextKey("service")).(*dbmodel.TenantServices)
	if dbmodel.ServiceKind(service.Kind) == dbmodel.ServiceKindThirdParty {
		endpoints, err := db.GetManager().EndpointsDao().List(serviceID)
		if err != nil {
			logrus.Errorf("find endpoints by sid[%s], error: %s", serviceID, err.Error())
			httputil.ReturnError(r, w, 500, "fund endpoints failure")
			return
		}
		for _, ep := range endpoints {
			address := validation.SplitEndpointAddress(ep.IP)
			if validation.IsDomainNotIP(address) {
				httputil.ReturnError(r, w, 400, "do not allow operate outer port for thirdpart domain endpoints")
				return
			}
		}
	}

	tenantName := r.Context().Value(middleware.ContextKey("tenant_name")).(string)
	portStr := chi.URLParam(r, "port")
	containerPort, err := strconv.Atoi(portStr)
	if err != nil {
		httputil.ReturnError(r, w, 400, "port must be a number")
		return
	}
	vsPort, protocol, errV := handler.GetServiceManager().PortOuter(tenantName, serviceID, containerPort, &data)
	if errV != nil {
		if strings.HasSuffix(errV.Error(), gorm.ErrRecordNotFound.Error()) {
			httputil.ReturnError(r, w, 404, errV.Error())
			return
		}
		httputil.ReturnError(r, w, 500, errV.Error())
		return
	}
	rc := make(map[string]string)
	domain := os.Getenv("EX_DOMAIN")
	if domain == "" {
		httputil.ReturnError(r, w, 500, "have no EX_DOMAIN")
		return
	}
	mm := strings.Split(domain, ":")
	if protocol == "http" || protocol == "https" {
		rc["domain"] = mm[0]
		if len(mm) == 2 {
			rc["port"] = mm[1]
		} else {
			rc["port"] = "10080"
		}
	} else if vsPort != nil {
		rc["domain"] = mm[0]
		rc["port"] = fmt.Sprintf("%v", vsPort.Port)
	}

	if err := handler.GetGatewayHandler().SendTask(map[string]interface{}{
		"service_id": serviceID,
		"action":     "port-" + data.Body.Operation,
		"port":       containerPort,
		"is_inner":   false,
	}); err != nil {
		logrus.Errorf("send runtime message about gateway failure %s", err.Error())
	}

	httputil.ReturnSuccess(r, w, rc)
}

//PortInnerController Switch port internal service
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/ports/{port}/inner v2 PortInnerController
//
// Switch internal services, the application does not need to be restarted, it takes effect automatically
//
// add port
//
// ---
// Consumes:
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
func (t *TenantStruct) PortInnerController(w http.ResponseWriter, r *http.Request) {
	var data api_model.ServicePortInnerOrOuter
	if !httputil.ValidatorRequestStructAndErrorResponse(r, w, &(data.Body), nil) {
		return
	}
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	tenantName := r.Context().Value(middleware.ContextKey("tenant_name")).(string)
	portStr := chi.URLParam(r, "port")
	containerPort, err := strconv.Atoi(portStr)
	if err != nil {
		httputil.ReturnError(r, w, 400, "port must be a number")
		return
	}
	if err := handler.GetServiceManager().PortInner(tenantName, serviceID, data.Body.Operation, containerPort); err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			httputil.ReturnError(r, w, 404, "service port "+err.Error())
			return
		} else if err.Error() == "already open" || err.Error() == "already close" {
			httputil.Return(r, w, 200, httputil.ResponseBody{Msg: err.Error()})
			return
		} else {
			httputil.ReturnError(r, w, 500, err.Error())
			return
		}
	}

	if err := handler.GetGatewayHandler().SendTask(map[string]interface{}{
		"service_id": serviceID,
		"action":     "port-" + data.Body.Operation,
		"port":       containerPort,
		"is_inner":   true,
	}); err != nil {
		logrus.Errorf("send runtime message about gateway failure %s", err.Error())
	}

	httputil.ReturnSuccess(r, w, nil)
}

//Pods pods
// swagger:operation GET  /v2/tenants/{tenant_name}/services/{service_alias}/pods v2 getPodsInfo
//
// Get pods information
//
// get pods info
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
func (t *TenantStruct) Pods(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	pods, err := handler.GetServiceManager().GetPods(serviceID)
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			logrus.Error("record not found:", err)
			httputil.ReturnError(r, w, 404, fmt.Sprintf("get pods error, %v", err))
			return
		}
		logrus.Error("get pods error:", err)
		httputil.ReturnError(r, w, 500, fmt.Sprintf("get pods error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, pods)
}

//Probe probe
func (t *TenantStruct) Probe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		t.UpdateProbe(w, r)
	case "DELETE":
		t.DeleteProbe(w, r)
	case "POST":
		t.AddProbe(w, r)
	}
}

//AddProbe add probe
// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/probe v2 addProbe
//
// Add application probe
//
// add probe
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
func (t *TenantStruct) AddProbe(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	var tsp api_model.ServiceProbe
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tsp, nil); !ok {
		return
	}
	var tspD dbmodel.TenantServiceProbe
	tspD.ServiceID = serviceID
	tspD.Cmd = tsp.Cmd
	tspD.FailureThreshold = tsp.FailureThreshold
	tspD.HTTPHeader = tsp.HTTPHeader
	tspD.InitialDelaySecond = tsp.InitialDelaySecond
	tspD.IsUsed = &tsp.IsUsed
	tspD.Mode = tsp.Mode
	tspD.Path = tsp.Path
	tspD.PeriodSecond = tsp.PeriodSecond
	tspD.Port = tsp.Port
	tspD.ProbeID = tsp.ProbeID
	tspD.Scheme = tsp.Scheme
	tspD.SuccessThreshold = tsp.SuccessThreshold
	tspD.TimeoutSecond = tsp.TimeoutSecond
	tspD.FailureAction = tsp.FailureAction
	//Pay attention to port issues
	if err := handler.GetServiceManager().ServiceProbe(&tspD, "add"); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("add service probe error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//UpdateProbe update probe
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/probe v2 updateProbe
//
// Update application probe information, *Note that this is a full update
//
// update probe
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
func (t *TenantStruct) UpdateProbe(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	var tsp api_model.ServiceProbe
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tsp, nil); !ok {
		return
	}
	var tspD dbmodel.TenantServiceProbe
	tspD.ServiceID = serviceID
	tspD.Cmd = tsp.Cmd
	tspD.FailureThreshold = tsp.FailureThreshold
	tspD.HTTPHeader = tsp.HTTPHeader
	tspD.InitialDelaySecond = tsp.InitialDelaySecond
	tspD.IsUsed = &tsp.IsUsed
	tspD.Mode = tsp.Mode
	tspD.Path = tsp.Path
	tspD.PeriodSecond = tsp.PeriodSecond
	tspD.Port = tsp.Port
	tspD.ProbeID = tsp.ProbeID
	tspD.Scheme = tsp.Scheme
	tspD.SuccessThreshold = tsp.SuccessThreshold
	tspD.TimeoutSecond = tsp.TimeoutSecond
	//Pay attention to port issues
	if err := handler.GetServiceManager().ServiceProbe(&tspD, "update"); err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			httputil.ReturnError(r, w, 404, fmt.Sprintf("update prob error, %v", err))
			return
		}
		httputil.ReturnError(r, w, 500, fmt.Sprintf("update service probe error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeleteProbe delete probe
// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/probe v2 deleteProbe
//
// Delete application probe
//
// delete probe
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
func (t *TenantStruct) DeleteProbe(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(middleware.ContextKey("service_id")).(string)
	var tsp api_model.ServiceProbe
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tsp, nil); !ok {
		return
	}
	var tspD dbmodel.TenantServiceProbe
	tspD.ServiceID = serviceID
	tspD.ProbeID = tsp.ProbeID
	//Pay attention to port issues
	if err := handler.GetServiceManager().ServiceProbe(&tspD, "delete"); err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			httputil.ReturnError(r, w, 404, fmt.Sprintf("delete prob error, %v", err))
			return
		}
		httputil.ReturnError(r, w, 500, fmt.Sprintf("delete service probe error, %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//Port Port
func (t *TenantStruct) Port(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		t.UpdatePort(w, r)
	case "DELETE":
		t.DeletePort(w, r)
	case "POST":
		t.AddPort(w, r)
	}
}

//AddPort add port
func (t *TenantStruct) AddPort(w http.ResponseWriter, r *http.Request) {
}

//DeletePort delete port
func (t *TenantStruct) DeletePort(w http.ResponseWriter, r *http.Request) {
}

//UpdatePort Update port
func (t *TenantStruct) UpdatePort(w http.ResponseWriter, r *http.Request) {
}

//SingleTenantResources SingleTenantResources
// swagger:operation GET /v2/tenants/{tenant_name}/resources v2 singletenantResources
//
// Specified tenant resource usage
//
// get tenant resources
//
// ---
// produces:
// - application/json
// - application/xml
// parameters:
// - name: tenant_name
//   in: path
//   description: tenant name
//   required: true
//   type: string
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) SingleTenantResources(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	//11ms
	services, err := handler.GetServiceManager().GetService(tenantID)
	if err != nil {
		msg := httputil.ResponseBody{
			Msg: fmt.Sprintf("get service error, %v", err),
		}
		httputil.Return(r, w, 500, msg)
	}
	//19ms
	statsInfo, _ := handler.GetTenantManager().StatsMemCPU(services)
	//900ms
	statsInfo.UUID = tenantID
	httputil.ReturnSuccess(r, w, statsInfo)
	return
}

//GetSupportProtocols GetSupportProtocols
// swagger:operation GET /v2/tenants/{tenant_name}/protocols v2 getSupportProtocols
//
// Get the protocols supported by the current data center
//
// get region protocols
//
// ---
// produces:
// - application/json
// - application/xml
// parameters:
// - name: tenant_name
//   in: path
//   description: tenant name
//   required: true
//   type: string
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) GetSupportProtocols(w http.ResponseWriter, r *http.Request) {
	rps, err := handler.GetTenantManager().GetProtocols()
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, rps)
	return
}

//TransPlugins transPlugins
// swagger:operation POST /v2/tenants/{tenant_name}/transplugins v2 transPlugins
//
// Install Gridworkz default plugins
//
// trans plugins
//
// ---
// produces:
// - application/json
// - application/xml
// parameters:
// - name: tenant_name
//   in: path
//   description: tenant name
//   required: true
//   type: string
//
// responses:
//   default:
//     schema:
//       "$ref": "#/responses/commandResponse"
//     description: Unified return format
func (t *TenantStruct) TransPlugins(w http.ResponseWriter, r *http.Request) {
	var tps api_model.TransPlugins
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tps.Body, nil)
	if !ok {
		return
	}
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	tenantName := r.Context().Value(middleware.ContextKey("tenant_name")).(string)
	rc := make(map[string]string)
	err := handler.GetTenantManager().TransPlugins(tenantID, tenantName, tps.Body.FromTenantName, tps.Body.PluginsID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	rc["result"] = "success"
	httputil.ReturnSuccess(r, w, rc)
	return
}
