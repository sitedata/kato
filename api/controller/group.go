// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/handler/group"
	"github.com/gridworkz/kato/api/middleware"
	httputil "github.com/gridworkz/kato/util/http"
)

//Backups list all backup history by group app
func Backups(w http.ResponseWriter, r *http.Request) {
	groupID := r.FormValue("group_id")
	if groupID == "" {
		httputil.ReturnError(r, w, 400, "group id can not be empty")
		return
	}
	list, err := handler.GetAPPBackupHandler().GetBackupByGroupID(groupID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, list)
}

//NewBackups new group app backup
func NewBackups(w http.ResponseWriter, r *http.Request) {
	var gb group.Backup
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &gb.Body, nil)
	if !ok {
		return
	}
	bean, err := handler.GetAPPBackupHandler().NewBackup(gb)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, bean)
}

//BackupCopy backup copy
func BackupCopy(w http.ResponseWriter, r *http.Request) {
	var gb group.BackupCopy
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &gb.Body, nil)
	if !ok {
		return
	}
	bean, err := handler.GetAPPBackupHandler().BackupCopy(gb)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, bean)
}

//Restore restore group app
func Restore(w http.ResponseWriter, r *http.Request) {
	var br group.BackupRestore
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &br.Body, nil)
	if !ok {
		return
	}
	br.BackupID = chi.URLParam(r, "backup_id")
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)
	br.Body.TenantID = tenantID
	bean, err := handler.GetAPPBackupHandler().RestoreBackup(br)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, bean)
}

//RestoreResult restore group app result
func RestoreResult(w http.ResponseWriter, r *http.Request) {
	restoreID := chi.URLParam(r, "restore_id")
	bean, err := handler.GetAPPBackupHandler().RestoreBackupResult(restoreID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, bean)
}

//GetBackup get one backup status
func GetBackup(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "backup_id")
	bean, err := handler.GetAPPBackupHandler().GetBackup(backupID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnSuccess(r, w, bean)
}

//DeleteBackup delete backup
func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "backup_id")

	err := handler.GetAPPBackupHandler().DeleteBackup(backupID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 404, "not found")
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}
