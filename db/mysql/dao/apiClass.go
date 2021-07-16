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

//RegionAPIClassDaoImpl RegionAPIClassDaoImpl
type RegionAPIClassDaoImpl struct {
	DB *gorm.DB
}

//AddModel - add api classification information
func (t *RegionAPIClassDaoImpl) AddModel(mo model.Interface) error {
	info := mo.(*model.RegionAPIClass)
	var oldInfo model.RegionAPIClass
	if ok := t.DB.Where("prefix = ? and class_level=?", info.Prefix, info.ClassLevel).Find(&oldInfo).RecordNotFound(); ok {
		if err := t.DB.Create(info).Error; err != nil {
			return err
		}
	} else {
		return fmt.Errorf("prefix %s is exist", info.Prefix)
	}
	return nil
}

//UpdateModel - update api classification information
func (t *RegionAPIClassDaoImpl) UpdateModel(mo model.Interface) error {
	info := mo.(*model.RegionAPIClass)
	if info.ID == 0 {
		return fmt.Errorf("region user info id can not be empty when update ")
	}
	if err := t.DB.Save(info).Error; err != nil {
		return err
	}
	return nil
}

//GetPrefixesByClass
func (t *RegionAPIClassDaoImpl) GetPrefixesByClass(apiClass string) ([]*model.RegionAPIClass, error) {
	var racs []*model.RegionAPIClass
	if err := t.DB.Select("prefix").Where("class_level =?", apiClass).Find(&racs).Error; err != nil {
		return nil, err
	}
	return racs, nil
}

//DeletePrefixInClass
func (t *RegionAPIClassDaoImpl) DeletePrefixInClass(apiClass, prefix string) error {
	relation := &model.RegionAPIClass{
		ClassLevel: apiClass,
		Prefix:     prefix,
	}
	if err := t.DB.Where("class_level=? and prefix=?", apiClass, prefix).Delete(relation).Error; err != nil {
		return err
	}
	return nil
}
