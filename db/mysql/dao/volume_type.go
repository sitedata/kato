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

package dao

import (
	"fmt"

	"github.com/gridworkz/kato/db/model"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

//VolumeTypeDaoImpl - license model management
type VolumeTypeDaoImpl struct {
	DB *gorm.DB
}

// CreateOrUpdateVolumeType find or create volumeType, !!! attention：just for store sync storageclass from k8s
func (vtd *VolumeTypeDaoImpl) CreateOrUpdateVolumeType(vt *model.TenantServiceVolumeType) (*model.TenantServiceVolumeType, error) {
	if vt.VolumeType == model.ShareFileVolumeType.String() || vt.VolumeType == model.LocalVolumeType.String() || vt.VolumeType == model.MemoryFSVolumeType.String() {
		return vt, nil
	}
	volumeType, err := vtd.GetVolumeTypeByType(vt.VolumeType)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound || volumeType == nil {
		logrus.Debugf("volume type[%s] do not exists, create it", vt.VolumeType)
		err = vtd.AddModel(vt)
	} else {
		logrus.Debugf("volume type[%s] already exists, update it", vt.VolumeType)
		volumeType.Provisioner = vt.Provisioner
		volumeType.StorageClassDetail = vt.StorageClassDetail
		volumeType.NameShow = vt.NameShow
		err = vtd.UpdateModel(volumeType)
	}
	return volumeType, err
}

//AddModel
func (vtd *VolumeTypeDaoImpl) AddModel(mo model.Interface) error {
	volumeType := mo.(*model.TenantServiceVolumeType)
	var oldVolumeType model.TenantServiceVolumeType
	if ok := vtd.DB.Where("volume_type=?", volumeType.VolumeType).Find(&oldVolumeType).RecordNotFound(); ok {
		if err := vtd.DB.Create(volumeType).Error; err != nil {
			return err
		}
	} else {
		return fmt.Errorf("volumeType is exist")
	}
	return nil
}

// UpdateModel
func (vtd *VolumeTypeDaoImpl) UpdateModel(mo model.Interface) error {
	volumeType := mo.(*model.TenantServiceVolumeType)
	if err := vtd.DB.Save(volumeType).Error; err != nil {
		return err
	}
	return nil
}

// GetAllVolumeTypes
func (vtd *VolumeTypeDaoImpl) GetAllVolumeTypes() ([]*model.TenantServiceVolumeType, error) {
	var volumeTypes []*model.TenantServiceVolumeType
	if err := vtd.DB.Find(&volumeTypes).Error; err != nil {
		return nil, err
	}
	return volumeTypes, nil
}

// GetAllVolumeTypesByPage
func (vtd *VolumeTypeDaoImpl) GetAllVolumeTypesByPage(page int, pageSize int) ([]*model.TenantServiceVolumeType, error) {
	var volumeTypes []*model.TenantServiceVolumeType
	if err := vtd.DB.Limit(pageSize).Offset((page - 1) * pageSize).Find(&volumeTypes).Error; err != nil {
		return nil, err
	}
	return volumeTypes, nil
}

// GetVolumeTypeByType
func (vtd *VolumeTypeDaoImpl) GetVolumeTypeByType(vt string) (*model.TenantServiceVolumeType, error) {
	var volumeType model.TenantServiceVolumeType
	if err := vtd.DB.Where("volume_type=?", vt).Find(&volumeType).Error; err != nil {
		return nil, err
	}
	return &volumeType, nil
}

// DeleteModelByVolumeTypes
func (vtd *VolumeTypeDaoImpl) DeleteModelByVolumeTypes(volumeType string) error {
	if err := vtd.DB.Where("volume_type=?", volumeType).Delete(&model.TenantServiceVolumeType{}).Error; err != nil {
		return err
	}
	return nil
}
