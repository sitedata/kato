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

	"github.com/jinzhu/gorm"
)

//ServiceProbeDaoImpl
type ServiceProbeDaoImpl struct {
	DB *gorm.DB
}

//AddModel - add application Probe
func (t *ServiceProbeDaoImpl) AddModel(mo model.Interface) error {
	probe := mo.(*model.TenantServiceProbe)
	var oldProbe model.TenantServiceProbe
	if ok := t.DB.Where("service_id=? and mode=?", probe.ServiceID, probe.Mode).Find(&oldProbe).RecordNotFound(); ok {
		if err := t.DB.Create(probe).Error; err != nil {
			return err
		}
	} else {
		return fmt.Errorf("probe mode %s of service %s is exist", probe.Mode, probe.ServiceID)
	}
	return nil
}

//UpdateModel - update application Probe
func (t *ServiceProbeDaoImpl) UpdateModel(mo model.Interface) error {
	probe := mo.(*model.TenantServiceProbe)
	if probe.ID == 0 {
		var oldProbe model.TenantServiceProbe
		if err := t.DB.Where("service_id = ? and probe_id=?", probe.ServiceID,
			probe.ProbeID).Find(&oldProbe).Error; err != nil {
			return err
		}
		if oldProbe.ID == 0 {
			return gorm.ErrRecordNotFound
		}
		probe.ID = oldProbe.ID
		probe.CreatedAt = oldProbe.CreatedAt
	}
	return t.DB.Save(probe).Error
}

//DeleteModel
func (t *ServiceProbeDaoImpl) DeleteModel(serviceID string, args ...interface{}) error {
	probeID := args[0].(string)
	relation := &model.TenantServiceProbe{
		ServiceID: serviceID,
		ProbeID:   probeID,
	}
	if err := t.DB.Where("service_id=? and probe_id=?", serviceID, probeID).Delete(relation).Error; err != nil {
		return err
	}
	return nil
}

// DelByServiceID deletes TenantServiceProbe based on sid(service_id)
func (t *ServiceProbeDaoImpl) DelByServiceID(sid string) error {
	return t.DB.Where("service_id=?", sid).Delete(&model.TenantServiceProbe{}).Error
}

//GetServiceProbes
func (t *ServiceProbeDaoImpl) GetServiceProbes(serviceID string) ([]*model.TenantServiceProbe, error) {
	var probes []*model.TenantServiceProbe
	if err := t.DB.Where("service_id=?", serviceID).Find(&probes).Error; err != nil {
		return nil, err
	}
	return probes, nil
}

//GetServiceUsedProbe
func (t *ServiceProbeDaoImpl) GetServiceUsedProbe(serviceID, mode string) (*model.TenantServiceProbe, error) {
	var probe model.TenantServiceProbe
	if err := t.DB.Where("service_id=? and mode=? and is_used=?", serviceID, mode, 1).Find(&probe).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &probe, nil
}

//DELServiceProbesByServiceID
func (t *ServiceProbeDaoImpl) DELServiceProbesByServiceID(serviceID string) error {
	probes := &model.TenantServiceProbe{
		ServiceID: serviceID,
	}
	if err := t.DB.Where("service_id=?", serviceID).Delete(probes).Error; err != nil {
		return err
	}
	return nil
}

//LocalSchedulerDaoImpl - local scheduling storage mysql implementation
type LocalSchedulerDaoImpl struct {
	DB *gorm.DB
}

//AddModel
func (t *LocalSchedulerDaoImpl) AddModel(mo model.Interface) error {
	ls := mo.(*model.LocalScheduler)
	var oldLs model.LocalScheduler
	if ok := t.DB.Where("service_id=? and pod_name=?", ls.ServiceID, ls.PodName).Find(&oldLs).RecordNotFound(); ok {
		if err := t.DB.Create(ls).Error; err != nil {
			return err
		}
	} else {
		return fmt.Errorf("service %s local scheduler of pod  %s is exist", ls.ServiceID, ls.PodName)
	}
	return nil
}

//UpdateModel
func (t *LocalSchedulerDaoImpl) UpdateModel(mo model.Interface) error {
	ls := mo.(*model.LocalScheduler)
	if ls.ID == 0 {
		return fmt.Errorf("LocalScheduler id can not be empty when update ")
	}
	if err := t.DB.Save(ls).Error; err != nil {
		return err
	}
	return nil
}

//GetLocalScheduler
func (t *LocalSchedulerDaoImpl) GetLocalScheduler(serviceID string) ([]*model.LocalScheduler, error) {
	var ls []*model.LocalScheduler
	if err := t.DB.Where("service_id=?", serviceID).Find(&ls).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ls, nil
}

//ServiceSourceImpl
type ServiceSourceImpl struct {
	DB *gorm.DB
}

//AddModel
func (t *ServiceSourceImpl) AddModel(mo model.Interface) error {
	ls := mo.(*model.ServiceSourceConfig)
	var oldLs model.ServiceSourceConfig
	if ok := t.DB.Where("service_id=? and source_type=?", ls.ServiceID, ls.SourceType).Find(&oldLs).RecordNotFound(); ok {
		if err := t.DB.Create(ls).Error; err != nil {
			return err
		}
	} else {
		oldLs.SourceBody = ls.SourceBody
		t.DB.Save(oldLs)
	}
	return nil
}

//UpdateModel
func (t *ServiceSourceImpl) UpdateModel(mo model.Interface) error {
	ls := mo.(*model.LocalScheduler)
	if ls.ID == 0 {
		return fmt.Errorf("ServiceSourceImpl id can not be empty when update ")
	}
	if err := t.DB.Save(ls).Error; err != nil {
		return err
	}
	return nil
}

//GetServiceSource
func (t *ServiceSourceImpl) GetServiceSource(serviceID string) ([]*model.ServiceSourceConfig, error) {
	var serviceSources []*model.ServiceSourceConfig
	if err := t.DB.Where("service_id=?", serviceID).Find(&serviceSources).Error; err != nil {
		return nil, err
	}
	return serviceSources, nil
}
