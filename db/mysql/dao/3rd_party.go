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
	"reflect"
	"strings"

	"github.com/gridworkz/kato/db/errors"
	"github.com/gridworkz/kato/db/model"
	"github.com/jinzhu/gorm"
)

// EndpointDaoImpl
type EndpointDaoImpl struct {
	DB *gorm.DB
}

// AddModel add one record for table 3rd_party_svc_endpoint
func (e *EndpointDaoImpl) AddModel(mo model.Interface) error {
	ep, ok := mo.(*model.Endpoint)
	if !ok {
		return fmt.Errorf("Type conversion error. From %s to *model.Endpoint", reflect.TypeOf(mo))
	}
	var o model.Endpoint
	if ok := e.DB.Where("service_id=? and ip=? and port=?", ep.ServiceID, ep.IP, ep.Port).Find(&o).RecordNotFound(); ok {
		if err := e.DB.Save(ep).Error; err != nil {
			return err
		}
	} else {
		return errors.ErrRecordAlreadyExist
	}
	return nil
}

// UpdateModel updates one record for table 3rd_party_svc_endpoint
func (e *EndpointDaoImpl) UpdateModel(mo model.Interface) error {
	ep, ok := mo.(*model.Endpoint)
	if !ok {
		return fmt.Errorf("Type conversion error. From %s to *model.Endpoint", reflect.TypeOf(mo))
	}
	if strings.Replace(ep.UUID, " ", "", -1) == "" {
		return fmt.Errorf("uuid can not be empty")
	}
	return e.DB.Save(ep).Error
}

// GetByUUID returns endpints matching the given uuid.
func (e *EndpointDaoImpl) GetByUUID(uuid string) (*model.Endpoint, error) {
	var ep model.Endpoint
	if err := e.DB.Where("uuid=?", uuid).Find(&ep).Error; err != nil {
		return nil, err
	}
	return &ep, nil
}

// List list all endpints matching the given serivce_id(sid).
func (e *EndpointDaoImpl) List(sid string) ([]*model.Endpoint, error) {
	var eps []*model.Endpoint
	if err := e.DB.Where("service_id=?", sid).Find(&eps).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return eps, nil
}

// ListIsOnline lists *model.Endpoint according to sid, and filter out the ones that are not online.
func (e *EndpointDaoImpl) ListIsOnline(sid string) ([]*model.Endpoint, error) {
	var eps []*model.Endpoint
	if err := e.DB.Where("service_id=? and is_online=1", sid).Find(&eps).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return eps, nil
}

// DelByUUID deletes endpoints matching uuid.
func (e *EndpointDaoImpl) DelByUUID(uuid string) error {
	if err := e.DB.Where("uuid=?", uuid).Delete(&model.Endpoint{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteByServiceID delete endpoints based on service id.
func (e *EndpointDaoImpl) DeleteByServiceID(sid string) error {
	return e.DB.Where("service_id=?", sid).Delete(&model.Endpoint{}).Error
}

// ThirdPartySvcDiscoveryCfgDaoImpl implements ThirdPartySvcDiscoveryCfgDao
type ThirdPartySvcDiscoveryCfgDaoImpl struct {
	DB *gorm.DB
}

// AddModel add one record for table 3rd_party_svc_discovery_cfg.
func (t *ThirdPartySvcDiscoveryCfgDaoImpl) AddModel(mo model.Interface) error {
	cfg, ok := mo.(*model.ThirdPartySvcDiscoveryCfg)
	if !ok {
		return fmt.Errorf("Type conversion error. From %s to *model.ThirdPartySvcDiscoveryCfg",
			reflect.TypeOf(mo))
	}
	var old model.ThirdPartySvcDiscoveryCfg
	if ok := t.DB.Where("service_id=?", cfg.ServiceID).Find(&old).RecordNotFound(); ok {
		if err := t.DB.Create(cfg).Error; err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Discovery configuration exists based on servicd_id(%s)", cfg.ServiceID)
	}
	return nil
}

// UpdateModel blabla
func (t *ThirdPartySvcDiscoveryCfgDaoImpl) UpdateModel(mo model.Interface) error {
	return nil
}

// GetByServiceID return third-party service discovery configuration according to service_id.
func (t *ThirdPartySvcDiscoveryCfgDaoImpl) GetByServiceID(sid string) (*model.ThirdPartySvcDiscoveryCfg, error) {
	var cfg model.ThirdPartySvcDiscoveryCfg
	if err := t.DB.Where("service_id=?", sid).Find(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cfg, nil
}

// DeleteByServiceID delete discovery config based on service id.
func (t *ThirdPartySvcDiscoveryCfgDaoImpl) DeleteByServiceID(sid string) error {
	return t.DB.Where("service_id=?", sid).Delete(&model.ThirdPartySvcDiscoveryCfg{}).Error
}
