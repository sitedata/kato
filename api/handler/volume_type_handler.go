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
	"encoding/json"
	"strings"

	"fmt"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/worker/client"
	"github.com/gridworkz/kato/worker/server/pb"
	"github.com/sirupsen/logrus"
	// pb "github.com/gridworkz/rainibond/worker/server/pb"
)

//VolumeTypeHandler LicenseAction
type VolumeTypeHandler interface {
	VolumeTypeVar(action string, vtm *dbmodel.TenantServiceVolumeType) error
	GetAllVolumeTypes() ([]*api_model.VolumeTypeStruct, error)
	GetAllVolumeTypesByPage(page int, pageSize int) ([]*api_model.VolumeTypeStruct, error)
	GetVolumeTypeByType(volumeType string) (*dbmodel.TenantServiceVolumeType, error)
	GetAllStorageClasses() ([]*pb.StorageClassDetail, error)
	VolumeTypeAction(action, volumeTypeID string) error
	DeleteVolumeType(volumeTypeID string) error
	SetVolumeType(vtm *api_model.VolumeTypeStruct) error
	UpdateVolumeType(dbVolume *dbmodel.TenantServiceVolumeType, vol *api_model.VolumeTypeStruct) error
}

var defaultVolumeTypeHandler VolumeTypeHandler

//CreateVolumeTypeManger create VolumeType manager
func CreateVolumeTypeManger(statusCli *client.AppRuntimeSyncClient) *VolumeTypeAction {
	return &VolumeTypeAction{statusCli: statusCli}
}

//GetVolumeTypeHandler get volumeType handler
func GetVolumeTypeHandler() VolumeTypeHandler {
	return defaultVolumeTypeHandler
}

// VolumeTypeAction action
type VolumeTypeAction struct {
	statusCli *client.AppRuntimeSyncClient
}

// VolumeTypeVar volume type crud
func (vta *VolumeTypeAction) VolumeTypeVar(action string, vtm *dbmodel.TenantServiceVolumeType) error {
	switch action {
	case "add":
		logrus.Debug("add volumeType")
	case "update":
		logrus.Debug("update volumeType")
	}
	return nil
}

// GetAllVolumeTypes get all volume types
func (vta *VolumeTypeAction) GetAllVolumeTypes() ([]*api_model.VolumeTypeStruct, error) {
	var optionList []*api_model.VolumeTypeStruct
	volumeTypeMap := make(map[string]*dbmodel.TenantServiceVolumeType)
	volumeTypes, err := db.GetManager().VolumeTypeDao().GetAllVolumeTypes()
	if err != nil {
		logrus.Errorf("get all volumeTypes error: %s", err.Error())
		return nil, err
	}

	for _, vt := range volumeTypes {
		volumeTypeMap[vt.VolumeType] = vt
		capacityValidation := make(map[string]interface{})
		if vt.CapacityValidation != "" {
			err := json.Unmarshal([]byte(vt.CapacityValidation), &capacityValidation)
			if err != nil {
				logrus.Error(err.Error())
				return nil, fmt.Errorf("format volume type capacity validation error")
			}
		}

		storageClassDetail := make(map[string]interface{})
		if vt.StorageClassDetail != "" {
			err := json.Unmarshal([]byte(vt.StorageClassDetail), &storageClassDetail)
			if err != nil {
				logrus.Error(err.Error())
				return nil, fmt.Errorf("format storageclass detail error")
			}
		}
		accessMode := strings.Split(vt.AccessMode, ",")
		sharePolicy := strings.Split(vt.SharePolicy, ",")
		backupPolicy := strings.Split(vt.BackupPolicy, ",")
		optionList = append(optionList, &api_model.VolumeTypeStruct{
			VolumeType:         vt.VolumeType,
			NameShow:           vt.NameShow,
			Provisioner:        vt.Provisioner,
			CapacityValidation: capacityValidation,
			Description:        vt.Description,
			AccessMode:         accessMode,
			SharePolicy:        sharePolicy,
			BackupPolicy:       backupPolicy,
			ReclaimPolicy:      vt.ReclaimPolicy,
			StorageClassDetail: storageClassDetail,
			Sort:               vt.Sort,
			Enable:             vt.Enable,
		})
	}

	return optionList, nil
}

// GetAllVolumeTypesByPage get all volume types by page
func (vta *VolumeTypeAction) GetAllVolumeTypesByPage(page int, pageSize int) ([]*api_model.VolumeTypeStruct, error) {

	var optionList []*api_model.VolumeTypeStruct
	volumeTypeMap := make(map[string]*dbmodel.TenantServiceVolumeType)
	volumeTypes, err := db.GetManager().VolumeTypeDao().GetAllVolumeTypesByPage(page, pageSize)
	if err != nil {
		logrus.Errorf("get all volumeTypes error: %s", err.Error())
		return nil, err
	}

	for _, vt := range volumeTypes {
		volumeTypeMap[vt.VolumeType] = vt
		capacityValidation := make(map[string]interface{})
		if vt.CapacityValidation != "" {
			err := json.Unmarshal([]byte(vt.CapacityValidation), &capacityValidation)
			if err != nil {
				logrus.Error(err.Error())
				return nil, fmt.Errorf("format volume type capacity validation error")
			}
		}

		storageClassDetail := make(map[string]interface{})
		if vt.StorageClassDetail != "" {
			err := json.Unmarshal([]byte(vt.StorageClassDetail), &storageClassDetail)
			if err != nil {
				logrus.Error(err.Error())
				return nil, fmt.Errorf("format storageclass detail error")
			}
		}
		accessMode := strings.Split(vt.AccessMode, ",")
		sharePolicy := strings.Split(vt.SharePolicy, ",")
		backupPolicy := strings.Split(vt.BackupPolicy, ",")
		optionList = append(optionList, &api_model.VolumeTypeStruct{
			VolumeType:         vt.VolumeType,
			NameShow:           vt.NameShow,
			CapacityValidation: capacityValidation,
			Description:        vt.Description,
			AccessMode:         accessMode,
			SharePolicy:        sharePolicy,
			BackupPolicy:       backupPolicy,
			ReclaimPolicy:      vt.ReclaimPolicy,
			StorageClassDetail: storageClassDetail,
			Sort:               vt.Sort,
			Enable:             vt.Enable,
		})
	}

	return optionList, nil
}

// GetVolumeTypeByType get volume type by type
func (vta *VolumeTypeAction) GetVolumeTypeByType(volumtType string) (*dbmodel.TenantServiceVolumeType, error) {
	return db.GetManager().VolumeTypeDao().GetVolumeTypeByType(volumtType)
}

// GetAllStorageClasses get all storage class
func (vta *VolumeTypeAction) GetAllStorageClasses() ([]*pb.StorageClassDetail, error) {
	sces, err := vta.statusCli.GetStorageClasses()
	if err != nil {
		return nil, err
	}
	return sces.List, nil
}

// VolumeTypeAction open volme type or close it
func (vta *VolumeTypeAction) VolumeTypeAction(action, volumeTypeID string) error {
	return nil
}

// DeleteVolumeType delte volume type
func (vta *VolumeTypeAction) DeleteVolumeType(volumeType string) error {
	db.GetManager().VolumeTypeDao().DeleteModelByVolumeTypes(volumeType)
	return nil
}

// SetVolumeType set volume type
func (vta *VolumeTypeAction) SetVolumeType(vol *api_model.VolumeTypeStruct) error {
	var accessMode []string
	var sharePolicy []string
	var backupPolicy []string
	jsonCapacityValidationStr, _ := json.Marshal(vol.CapacityValidation)
	jsonStorageClassDetailStr, _ := json.Marshal(vol.StorageClassDetail)
	if vol.AccessMode == nil {
		accessMode[1] = "RWO"
	} else {
		accessMode = vol.AccessMode
	}
	if vol.SharePolicy == nil {
		sharePolicy[1] = "exclusive"
	} else {
		sharePolicy = vol.SharePolicy
	}

	if vol.BackupPolicy == nil {
		backupPolicy[1] = "exclusive"
	} else {
		backupPolicy = vol.BackupPolicy
	}

	dbVolume := dbmodel.TenantServiceVolumeType{}
	dbVolume.VolumeType = vol.VolumeType
	dbVolume.NameShow = vol.NameShow
	dbVolume.CapacityValidation = string(jsonCapacityValidationStr)
	dbVolume.Description = vol.Description
	dbVolume.AccessMode = strings.Join(accessMode, ",")
	dbVolume.SharePolicy = strings.Join(sharePolicy, ",")
	dbVolume.BackupPolicy = strings.Join(backupPolicy, ",")
	dbVolume.ReclaimPolicy = vol.ReclaimPolicy
	dbVolume.StorageClassDetail = string(jsonStorageClassDetailStr) // TODO StorageClass normative verification, and return the correct structure, assign the provisoner in the structure
	dbVolume.Provisioner = "provisioner"                            // TODO According to StorageClass
	dbVolume.Sort = vol.Sort
	dbVolume.Enable = vol.Enable

	err := db.GetManager().VolumeTypeDao().AddModel(&dbVolume)
	return err
}

// UpdateVolumeType
func (vta *VolumeTypeAction) UpdateVolumeType(dbVolume *dbmodel.TenantServiceVolumeType, vol *api_model.VolumeTypeStruct) error {
	var accessMode []string
	var sharePolicy []string
	var backupPolicy []string
	jsonCapacityValidationStr, _ := json.Marshal(vol.CapacityValidation)
	jsonStorageClassDetailStr, _ := json.Marshal(vol.StorageClassDetail)
	if vol.AccessMode == nil {
		accessMode[1] = "RWO"
	} else {
		accessMode = vol.AccessMode
	}
	if vol.SharePolicy == nil {
		sharePolicy[1] = "exclusive"
	} else {
		sharePolicy = vol.SharePolicy
	}

	if vol.BackupPolicy == nil {
		backupPolicy[1] = "exclusive"
	} else {
		backupPolicy = vol.BackupPolicy
	}

	dbVolume.VolumeType = vol.VolumeType
	dbVolume.NameShow = vol.NameShow
	dbVolume.CapacityValidation = string(jsonCapacityValidationStr)
	dbVolume.Description = vol.Description
	dbVolume.AccessMode = strings.Join(accessMode, ",")
	dbVolume.SharePolicy = strings.Join(sharePolicy, ",")
	dbVolume.BackupPolicy = strings.Join(backupPolicy, ",")
	dbVolume.ReclaimPolicy = vol.ReclaimPolicy
	dbVolume.StorageClassDetail = string(jsonStorageClassDetailStr)
	dbVolume.Sort = vol.Sort
	dbVolume.Enable = vol.Enable

	err := db.GetManager().VolumeTypeDao().UpdateModel(dbVolume)
	return err
}
