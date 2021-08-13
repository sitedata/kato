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
	"net/http"

	"github.com/gridworkz/kato/api/handler/share"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/util"

	api_model "github.com/gridworkz/kato/api/model"
	ctxutil "github.com/gridworkz/kato/api/util/ctx"
	httputil "github.com/gridworkz/kato/util/http"
)

//PluginAction plugin action
func (t *TenantStruct) PluginAction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		t.UpdatePlugin(w, r)
	case "DELETE":
		t.DeletePlugin(w, r)
	case "POST":
		t.CreatePlugin(w, r)
	case "GET":
		t.GetPlugins(w, r)
	}
}

//CreatePlugin add plugin
func (t *TenantStruct) CreatePlugin(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants/{tenant_name}/plugin v2 createPlugin
	//
	// Create plugin
	//
	// create plugin
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
	// description: unified return format
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	tenantName := r.Context().Value(ctxutil.ContextKey("tenant_name")).(string)
	var cps api_model.CreatePluginStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &cps.Body, nil); !ok {
		return
	}
	cps.Body.TenantID = tenantID
	cps.TenantName = tenantName
	if err := handler.GetPluginManager().CreatePluginAct(&cps); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//UpdatePlugin UpdatePlugin
func (t *TenantStruct) UpdatePlugin(w http.ResponseWriter, r *http.Request) {
	// swagger:operation PUT /v2/tenants/{tenant_name}/plugin/{plugin_id} v2 updatePlugin
	//
	// The plugin is updated in full, but the pluginID and the tenant do not provide modification
	//
	// update plugin
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
	// description: unified return format

	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	var ups api_model.UpdatePluginStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &ups.Body, nil); !ok {
		return
	}
	if err := handler.GetPluginManager().UpdatePluginAct(pluginID, tenantID, &ups); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeletePlugin DeletePlugin
func (t *TenantStruct) DeletePlugin(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /v2/tenants/{tenant_name}/plugin/{plugin_id} v2 deletePlugin
	//
	// plugin delete
	//
	// delete plugin
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
	// description: unified return format
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	if err := handler.GetPluginManager().DeletePluginAct(pluginID, tenantID); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//GetPlugins GetPlugins
func (t *TenantStruct) GetPlugins(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/plugin v2 getPlugins
	//
	// Get all available plugins under the current tenant
	//
	// get plugins
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
	// description: unified return format
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	plugins, err := handler.GetPluginManager().GetPlugins(tenantID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, plugins)
}

//PluginDefaultENV PluginDefaultENV
func (t *TenantStruct) PluginDefaultENV(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		t.AddDefatultENV(w, r)
	case "DELETE":
		t.DeleteDefaultENV(w, r)
	case "PUT":
		t.UpdateDefaultENV(w, r)
	}
}

// AddDefatultENV AddDefatultENV
func (t *TenantStruct) AddDefatultENV(w http.ResponseWriter, r *http.Request) {
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	versionID := chi.URLParam(r, "version_id")
	var est api_model.ENVStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &est.Body, nil); !ok {
		return
	}
	est.VersionID = versionID
	est.PluginID = pluginID
	if err := handler.GetPluginManager().AddDefaultEnv(&est); err != nil {
		err.Handle(r, w)
		return
	}
}

//DeleteDefaultENV DeleteDefaultENV
func (t *TenantStruct) DeleteDefaultENV(w http.ResponseWriter, r *http.Request) {
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	envName := chi.URLParam(r, "env_name")
	versionID := chi.URLParam(r, "version_id")
	if err := handler.GetPluginManager().DeleteDefaultEnv(pluginID, versionID, envName); err != nil {
		err.Handle(r, w)
		return
	}
}

//UpdateDefaultENV UpdateDefaultENV
func (t *TenantStruct) UpdateDefaultENV(w http.ResponseWriter, r *http.Request) {

	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	versionID := chi.URLParam(r, "version_id")
	var est api_model.ENVStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &est.Body, nil); !ok {
		return
	}
	est.PluginID = pluginID
	est.VersionID = versionID
	if err := handler.GetPluginManager().UpdateDefaultEnv(&est); err != nil {
		err.Handle(r, w)
		return
	}
}

//GetPluginDefaultEnvs GetPluginDefaultEnvs
func (t *TenantStruct) GetPluginDefaultEnvs(w http.ResponseWriter, r *http.Request) {
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	versionID := chi.URLParam(r, "version_id")
	envs, err := handler.GetPluginManager().GetDefaultEnv(pluginID, versionID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, envs)
}

//PluginBuild PluginBuild
// swagger:operation POST /v2/tenants/{tenant_name}/plugin/{plugin_id}/build v2 buildPlugin
//
// Build plugin
//
// build plugin
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
// description: unified return format
func (t *TenantStruct) PluginBuild(w http.ResponseWriter, r *http.Request) {
	var build api_model.BuildPluginStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &build.Body, nil)
	if !ok {
		return
	}
	tenantName := r.Context().Value(ctxutil.ContextKey("tenant_name")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	build.TenantName = tenantName
	build.PluginID = pluginID
	build.Body.TenantID = tenantID
	pbv, err := handler.GetPluginManager().BuildPluginManual(&build)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, pbv)
}

//GetAllPluginBuildVersions Get all the build versions of the plug-in
// swagger:operation GET /v2/tenants/{tenant_name}/plugin/{plugin_id}/build-version v2 allPluginVersions
//
// Get all build version information
//
// all plugin versions
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
// description: unified return format
func (t *TenantStruct) GetAllPluginBuildVersions(w http.ResponseWriter, r *http.Request) {
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	versions, err := handler.GetPluginManager().GetAllPluginBuildVersions (pluginID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, versions)
}

//GetPluginBuildVersion to obtain a build version information
// swagger:operation GET /v2/tenants/{tenant_name}/plugin/{plugin_id}/build-version/{version_id} v2 pluginVersion
//
// Get information about a build version
//
// plugin version
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
// description: unified return format
func (t *TenantStruct) GetPluginBuildVersion(w http.ResponseWriter, r *http.Request) {
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	versionID := chi.URLParam(r, "version_id")
	version, err := handler.GetPluginManager().GetPluginBuildVersion(pluginID, versionID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, version)
}

//DeletePluginBuildVersion DeletePluginBuildVersion
// swagger:operation DELETE /v2/tenants/{tenant_name}/plugin/{plugin_id}/build-version/{version_id} v2 deletePluginVersion
//
// Delete a certain build version information
//
// delete plugin version
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
// description: unified return format
func (t *TenantStruct) DeletePluginBuildVersion(w http.ResponseWriter, r *http.Request) {
	pluginID := r.Context().Value(ctxutil.ContextKey("plugin_id")).(string)
	versionID := chi.URLParam(r, "version_id")
	err := handler.GetPluginManager().DeletePluginBuildVersion(pluginID, versionID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//PluginSet PluginSet
func (t *TenantStruct) PluginSet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		t.updatePluginSet(w, r)
	case "POST":
		t.addPluginSet(w, r)
	case "GET":
		t.getPluginSet(w, r)
	}
}

// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/plugin v2 updatePluginSet
//
// Update plugin settings
//
// update plugin setting
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
// description: unified return format
func (t *TenantStruct) updatePluginSet(w http.ResponseWriter, r *http.Request) {
	var pss api_model.PluginSetStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &pss.Body, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	relation, err := handler.GetServiceManager().UpdateTenantServicePluginRelation(serviceID, &pss)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, relation)
}

// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/plugin v2 addPluginSet
//
// Add plugin settings
//
// add plugin setting
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
// description: unified return format
func (t *TenantStruct) addPluginSet(w http.ResponseWriter, r *http.Request) {
	var pss api_model.PluginSetStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &pss.Body, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	serviceAlias := r.Context().Value(ctxutil.ContextKey("service_alias")).(string)
	tenantName := r.Context().Value(ctxutil.ContextKey("tenant_name")).(string)
	pss.ServiceAlias = serviceAlias
	pss.TenantName = tenantName
	re, err := handler.GetServiceManager().SetTenantServicePluginRelation(tenantID, serviceID, &pss)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, re)
}

// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/plugin v2 getPluginSet
//
// Get plug-in settings
//
// get plugin setting
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
// description: unified return format
func (t *TenantStruct) getPluginSet(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	gps, err := handler.GetServiceManager().GetTenantServicePluginRelation(serviceID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, gps)

}

//DeletePluginRelation DeletePluginRelation
// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/plugin/{plugin_id} v2 deletePluginRelation
//
// Remove plugin dependencies
//
// delete plugin relation
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
// description: unified return format
func (t *TenantStruct) DeletePluginRelation(w http.ResponseWriter, r *http.Request) {
	pluginID := chi.URLParam(r, "plugin_id")
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	if err := handler.GetServiceManager().TenantServiceDeletePluginRelation(tenantID, serviceID, pluginID); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//GePluginEnvWhichCanBeSet GePluginEnvWhichCanBeSet
// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/plugin/{plugin_id}/envs v2 getVersionEnvs
//
// Get configurable env; take it from the service plugin correspondence, if it does not exist, return the default modifiable variable
//
// get version env
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
// description: unified return format
func (t *TenantStruct) GePluginEnvWhichCanBeSet(w http.ResponseWriter, r *http.Request) {
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	pluginID := chi.URLParam(r, "plugin_id")
	envs, err := handler.GetPluginManager().GetEnvsWhichCanBeSet(serviceID, pluginID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, envs)
}

//UpdateVersionEnv UpdateVersionEnv
// swagger:operation PUT /v2/tenants/{tenant_name}/services/{service_alias}/plugin/{plugin_id}/upenv v2 updateVersionEnv
//
// modify the app plugin config info. it will Thermal effect
//
// update version env
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
// description: unified return format
func (t *TenantStruct) UpdateVersionEnv(w http.ResponseWriter, r *http.Request) {
	var uve api_model.SetVersionEnv
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &uve.Body, nil)
	if !ok {
		return
	}
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	serviceAlias := r.Context().Value(ctxutil.ContextKey("service_alias")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	pluginID := chi.URLParam(r, "plugin_id")
	uve.PluginID = pluginID
	uve.Body.TenantID = tenantID
	uve.ServiceAlias = serviceAlias
	uve.Body.ServiceID = serviceID
	if err := handler.GetServiceManager().UpdateVersionEnv(&uve); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//SharePlugin share tenants plugin
func (t *TenantStruct) SharePlugin(w http.ResponseWriter, r *http.Request) {
	var sp share.PluginShare
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &sp.Body, nil)
	if !ok {
		return
	}
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	sp.TenantID = tenantID
	sp.PluginID = chi.URLParam(r, "plugin_id")
	if sp.Body.EventID == "" {
		sp.Body.EventID = util.NewUUID()
	}
	res, errS := handler.GetPluginShareHandle().Share(sp)
	if errS != nil {
		errS.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, res)
}

//SharePluginResult SharePluginResult
func (t *TenantStruct) SharePluginResult(w http.ResponseWriter, r *http.Request) {
	shareID := chi.URLParam(r, "share_id")
	res, errS := handler.GetPluginShareHandle().ShareResult(shareID)
	if errS != nil {
		errS.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, res)
}

//BatchInstallPlugins -
func (t *TenantStruct) BatchInstallPlugins(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	var req api_model.BatchCreatePlugins
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil); !ok {
		return
	}
	if err := handler.GetPluginManager().BatchCreatePlugins(tenantID, req.Plugins); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

// BatchBuildPlugins -
func (t *TenantStruct) BatchBuildPlugins(w http.ResponseWriter, r *http.Request) {
	var builds api_model.BatchBuildPlugins
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &builds, nil)
	if !ok {
		return
	}
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	err := handler.GetPluginManager().BatchBuildPlugins(&builds, tenantID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}
