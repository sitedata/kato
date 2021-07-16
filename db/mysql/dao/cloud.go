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
	"time"

	"github.com/gridworkz/kato/db/model"
	"github.com/jinzhu/gorm"
)

//RegionUserInfoDaoImpl
type RegionUserInfoDaoImpl struct {
	DB *gorm.DB
}

//AddModel - add cloud information
func (t *RegionUserInfoDaoImpl) AddModel(mo model.Interface) error {
	info := mo.(*model.RegionUserInfo)
	var oldInfo model.RegionUserInfo
	if ok := t.DB.Where("eid = ?", info.EID).Find(&oldInfo).RecordNotFound(); ok {
		if err := t.DB.Create(info).Error; err != nil {
			return err
		}
	} else {
		return fmt.Errorf("eid %s is exist", info.EID)
	}
	return nil
}

//UpdateModel - update cloud information
func (t *RegionUserInfoDaoImpl) UpdateModel(mo model.Interface) error {
	info := mo.(*model.RegionUserInfo)
	if info.ID == 0 {
		return fmt.Errorf("region user info id can not be empty when update ")
	}
	if err := t.DB.Save(info).Error; err != nil {
		return err
	}
	return nil
}

//GetTokenByEid
func (t *RegionUserInfoDaoImpl) GetTokenByEid(eid string) (*model.RegionUserInfo, error) {
	var rui model.RegionUserInfo
	if err := t.DB.Where("eid=?", eid).Find(&rui).Error; err != nil {
		return nil, err
	}
	return &rui, nil
}

//GetTokenByTokenID
func (t *RegionUserInfoDaoImpl) GetTokenByTokenID(token string) (*model.RegionUserInfo, error) {
	var rui model.RegionUserInfo
	if err := t.DB.Where("token=?", token).Find(&rui).Error; err != nil {
		return nil, err
	}
	return &rui, nil
}

//GetALLTokenInValidityPeriod
func (t *RegionUserInfoDaoImpl) GetALLTokenInValidityPeriod() ([]*model.RegionUserInfo, error) {
	var ruis []*model.RegionUserInfo
	timestamp := int(time.Now().Unix())
	if err := t.DB.Select("api_range, validity_period, token").Where("validity_period > ?", timestamp).Find(&ruis).Error; err != nil {
		return nil, err
	}
	return ruis, nil
}
