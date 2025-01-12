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
	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util/bcode"
	ctxutil "github.com/gridworkz/kato/api/util/ctx"
	dbmodel "github.com/gridworkz/kato/db/model"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/sirupsen/logrus"
)

//VolumeDependency VolumeDependency
func (t *TenantStruct) VolumeDependency(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		t.DeleteVolumeDependency(w, r)
	case "POST":
		t.AddVolumeDependency(w, r)
	}
}

//AddVolumeDependency add volume dependency
func (t *TenantStruct) AddVolumeDependency(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/volume-dependency v2 addVolumeDependency
	//
	// Increase application persistence dependency
	//
	// add volume dependency
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

	logrus.Debugf("trans add volumn dependency service ")
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	var tsr api_model.V2AddVolumeDependencyStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tsr.Body, nil); !ok {
		return
	}
	vd := &dbmodel.TenantServiceMountRelation{
		TenantID:        tenantID,
		ServiceID:       serviceID,
		DependServiceID: tsr.Body.DependServiceID,
		HostPath:        tsr.Body.MntDir,
		VolumePath:      tsr.Body.MntName,
	}
	if err := handler.GetServiceManager().VolumeDependency(vd, "add"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeleteVolumeDependency delete volume dependency
func (t *TenantStruct) DeleteVolumeDependency(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/volume-dependency v2 deleteVolumeDependency
	//
	// Delete application persistence dependencies
	//
	// delete volume dependency
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

	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	var tsr api_model.V2DelVolumeDependencyStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tsr.Body, nil); !ok {
		return
	}
	vd := &dbmodel.TenantServiceMountRelation{
		TenantID:        tenantID,
		ServiceID:       serviceID,
		DependServiceID: tsr.Body.DependServiceID,
	}
	if err := handler.GetServiceManager().VolumeDependency(vd, "delete"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

// AddVolume AddVolume
func (t *TenantStruct) AddVolume(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/volume v2 addVolume
	//
	// Add application persistent information
	//
	// add volume
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

	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	avs := &api_model.V2AddVolumeStruct{}
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &avs.Body, nil); !ok {
		return
	}
	tsv := &dbmodel.TenantServiceVolume{
		ServiceID:          serviceID,
		VolumePath:         avs.Body.VolumePath,
		HostPath:           avs.Body.HostPath,
		Category:           avs.Body.Category,
		VolumeCapacity:     avs.Body.VolumeCapacity,
		VolumeType:         dbmodel.ShareFileVolumeType.String(),
		VolumeProviderName: avs.Body.VolumeProviderName,
		AccessMode:         avs.Body.AccessMode,
		SharePolicy:        avs.Body.SharePolicy,
		BackupPolicy:       avs.Body.BackupPolicy,
		ReclaimPolicy:      avs.Body.ReclaimPolicy,
	}
	if !strings.HasPrefix(tsv.VolumePath, "/") {
		httputil.ReturnError(r, w, 400, "volume path is invalid,must begin with /")
		return
	}
	if err := handler.GetServiceManager().VolumnVar(tsv, tenantID, "", "add"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

// UpdVolume updates service volume.
func (t *TenantStruct) UpdVolume(w http.ResponseWriter, r *http.Request) {
	var req api_model.UpdVolumeReq
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil); !ok {
		return
	}

	if req.Mode != nil && (*req.Mode > 777 || *req.Mode < 0) {
		httputil.ReturnBcodeError(r, w, bcode.NewBadRequest("mode be a number between 0 and 777 (octal)"))
		return
	}

	sid := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	if err := handler.GetServiceManager().UpdVolume(sid, &req); err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
	}
	httputil.ReturnSuccess(r, w, "success")
}

//DeleteVolume DeleteVolume
func (t *TenantStruct) DeleteVolume(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/volume v2 deleteVolume
	//
	// Delete application persistent information
	//
	// delete volume
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

	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	avs := &api_model.V2DelVolumeStruct{}
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &avs.Body, nil); !ok {
		return
	}
	tsv := &dbmodel.TenantServiceVolume{
		ServiceID:  serviceID,
		VolumePath: avs.Body.VolumePath,
		Category:   avs.Body.Category,
	}
	if err := handler.GetServiceManager().VolumnVar(tsv, tenantID, "", "delete"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//The following is the V2.1 version persistence API, which supports multiple persistence modes

//AddVolumeDependency add volume dependency
func AddVolumeDependency(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/depvolumes v2 addDepVolume
	//
	// Increase application persistence dependency (V2.1 supports multiple types of storage)
	//
	// add volume dependency
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

	logrus.Debugf("trans add volumn dependency service ")
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	var tsr api_model.AddVolumeDependencyStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tsr.Body, nil); !ok {
		return
	}

	vd := &dbmodel.TenantServiceMountRelation{
		TenantID:        tenantID,
		ServiceID:       serviceID,
		DependServiceID: tsr.Body.DependServiceID,
		VolumeName:      tsr.Body.VolumeName,
		VolumePath:      tsr.Body.VolumePath,
		VolumeType:      tsr.Body.VolumeType,
	}
	if err := handler.GetServiceManager().VolumeDependency(vd, "add"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeleteVolumeDependency delete volume dependency
func DeleteVolumeDependency(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/depvolumes v2 delDepVolume
	//
	// Delete application persistence dependencies (V2.1 supports multiple types of storage)
	//
	// delete volume dependency
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

	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	var tsr api_model.DeleteVolumeDependencyStruct
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &tsr.Body, nil); !ok {
		return
	}
	vd := &dbmodel.TenantServiceMountRelation{
		TenantID:        tenantID,
		ServiceID:       serviceID,
		DependServiceID: tsr.Body.DependServiceID,
		VolumeName:      tsr.Body.VolumeName,
	}
	if err := handler.GetServiceManager().VolumeDependency(vd, "delete"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

// AddVolume AddVolume
func AddVolume(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/tenants/{tenant_name}/services/{service_alias}/volumes v2 addVolumes
	//
	// Add application persistent information (V2.1 supports multiple types of storage)
	//
	// add volume
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

	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	avs := &api_model.AddVolumeStruct{}
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &avs.Body, nil); !ok {
		return
	}

	if avs.Body.Mode != nil && (*avs.Body.Mode > 777 || *avs.Body.Mode < 0) {
		httputil.ReturnBcodeError(r, w, bcode.NewBadRequest("mode be a number between 0 and 777 (octal)"))
		return
	}

	tsv := &dbmodel.TenantServiceVolume{
		ServiceID:          serviceID,
		VolumeName:         avs.Body.VolumeName,
		VolumePath:         avs.Body.VolumePath,
		VolumeType:         avs.Body.VolumeType,
		Category:           avs.Body.Category,
		VolumeProviderName: avs.Body.VolumeProviderName,
		IsReadOnly:         avs.Body.IsReadOnly,
		VolumeCapacity:     avs.Body.VolumeCapacity,
		AccessMode:         avs.Body.AccessMode,
		SharePolicy:        avs.Body.SharePolicy,
		BackupPolicy:       avs.Body.BackupPolicy,
		ReclaimPolicy:      avs.Body.ReclaimPolicy,
		AllowExpansion:     avs.Body.AllowExpansion,
		Mode:               avs.Body.Mode,
	}

	// TODO gdevs validate VolumeCapacity  AccessMode SharePolicy BackupPolicy ReclaimPolicy AllowExpansion

	if !strings.HasPrefix(avs.Body.VolumePath, "/") {
		httputil.ReturnError(r, w, 400, "volume path is invalid,must begin with /")
		return
	}
	if err := handler.GetServiceManager().VolumnVar(tsv, tenantID, avs.Body.FileContent, "add"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//DeleteVolume DeleteVolume
func DeleteVolume(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /v2/tenants/{tenant_name}/services/{service_alias}/volumes/{volume_name} v2 deleteVolumes
	//
	// Delete application persistent information (V2.1 supports multiple types of storage)
	//
	// delete volume
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

	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	tenantID := r.Context().Value(ctxutil.ContextKey("tenant_id")).(string)
	tsv := &dbmodel.TenantServiceVolume{}
	tsv.ServiceID = serviceID
	tsv.VolumeName = chi.URLParam(r, "volume_name")
	if err := handler.GetServiceManager().VolumnVar(tsv, tenantID, "", "delete"); err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

//GetVolume Get all the storage of the application, including dependent storage
func GetVolume(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/volumes v2 getVolumes
	//
	// Get all storage of the application, including dependent storage (V2.1 supports multiple types of storage)
	//
	// get volumes
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
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	volumes, err := handler.GetServiceManager().GetVolumes(serviceID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, volumes)
}

//GetDepVolume Get all the storage that the application depends on
func GetDepVolume(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /v2/tenants/{tenant_name}/services/{service_alias}/depvolumes v2 getDepVolumes
	//
	// Get the storage that the application depends on (V2.1 supports multiple types of storage)
	//
	// get depvolumes
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
	serviceID := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	volumes, err := handler.GetServiceManager().GetDepVolumes(serviceID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, volumes)
}
