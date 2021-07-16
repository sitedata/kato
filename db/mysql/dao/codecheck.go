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
	"github.com/gridworkz/kato/db/model"

	"github.com/jinzhu/gorm"
)

//AddModel
func (c *CodeCheckResultDaoImpl) AddModel(mo model.Interface) error {
	result := mo.(*model.CodeCheckResult)
	var oldResult model.CodeCheckResult
	if ok := c.DB.Where("service_id=?", result.ServiceID).Find(&oldResult).RecordNotFound(); ok {
		if err := c.DB.Create(result).Error; err != nil {
			return err
		}
	} else {
		update(result, &oldResult)
		if err := c.DB.Save(&oldResult).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}

//UpdateModel
func (c *CodeCheckResultDaoImpl) UpdateModel(mo model.Interface) error {
	result := mo.(*model.CodeCheckResult)
	var oldResult model.CodeCheckResult
	if ok := c.DB.Where("service_id=?", result.ServiceID).Find(&oldResult).RecordNotFound(); !ok {
		update(result, &oldResult)
		if err := c.DB.Save(&oldResult).Error; err != nil {
			return err
		}
	}
	return nil
}

//CodeCheckResultDaoImpl
type CodeCheckResultDaoImpl struct {
	DB *gorm.DB
}

func update(target, old *model.CodeCheckResult) {
	//o,_:=json.Marshal(old)
	//t,_:=json.Marshal(target)
	//logrus.Infof("before update,stared is %s,target is ",string(o),string(t))
	if target.DockerFileReady != old.DockerFileReady {

		old.DockerFileReady = !old.DockerFileReady
	}
	if target.VolumeList != "" && target.VolumeList != "null" {
		old.VolumeList = target.VolumeList
	}
	if target.PortList != "" && target.PortList != "null" {
		old.PortList = target.PortList
	}
	if target.BuildImageName != "" {
		old.BuildImageName = target.BuildImageName
	}
	if target.VolumeMountPath != "" {
		old.VolumeMountPath = target.VolumeMountPath
	}
	if target.InnerPort != "" {
		old.InnerPort = target.InnerPort
	}
	//o2,_:=json.Marshal(old)
	//t2,_:=json.Marshal(target)
	//logrus.Infof("after update,%s,%s",string(o2),string(t2))
}

//GetCodeCheckResult get event log message
func (c *CodeCheckResultDaoImpl) GetCodeCheckResult(serviceID string) (*model.CodeCheckResult, error) {
	var result model.CodeCheckResult
	if err := c.DB.Where("service_id=?", serviceID).Find(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//return messageRaw, nil
		}
		return nil, err
	}
	return &result, nil
}

// DeleteByServiceID deletes a CodeCheckResult base on serviceID.
func (c *CodeCheckResultDaoImpl) DeleteByServiceID(serviceID string) error {
	return c.DB.Where("service_id=?", serviceID).Delete(&model.CodeCheckResult{}).Error
}
