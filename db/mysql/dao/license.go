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

//LicenseDaoImpl
type LicenseDaoImpl struct {
	DB *gorm.DB
}

//AddModel
func (l *LicenseDaoImpl) AddModel(mo model.Interface) error {
	license := mo.(*model.LicenseInfo)
	var oldLicense model.LicenseInfo
	if ok := l.DB.Where("license=?", license.License).Find(&oldLicense).RecordNotFound(); ok {
		if err := l.DB.Create(license).Error; err != nil {
			return err
		}
	} else {
		return fmt.Errorf("license exists")
	}
	return nil
}

//UpdateModel
func (l *LicenseDaoImpl) UpdateModel(mo model.Interface) error {
	return nil
}

//DeleteLicense
func (l *LicenseDaoImpl) DeleteLicense(token string) error {
	return nil
}

//ListLicenses
func (l *LicenseDaoImpl) ListLicenses() ([]*model.LicenseInfo, error) {
	var licenses []*model.LicenseInfo
	if err := l.DB.Find(&licenses).Error; err != nil {
		return nil, err
	}
	return licenses, nil
}
