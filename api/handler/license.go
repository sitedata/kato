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
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"

	"github.com/jinzhu/gorm"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
)

//LicenseAction LicenseAction
type LicenseAction struct{}

//PackLicense PackLicense
func (l *LicenseAction) PackLicense(encrypted string) ([]byte, error) {
	return decrypt(key, encrypted)
}

//StoreLicense StoreLicense
func (l *LicenseAction) StoreLicense(license, token string) error {

	ls := &dbmodel.LicenseInfo{
		Token:   token,
		License: license,
	}
	if err := db.GetManager().LicenseDao().AddModel(ls); err != nil {
		return err
	}
	return nil
}

//LicensesInfos
//verification
type LicensesInfos struct {
	Infos map[string]*api_model.LicenseInfo
}

//ShowInfos ShowInfos
func (l *LicensesInfos) ShowInfos() (map[string]*api_model.LicenseInfo, error) {
	return l.Infos, nil
}

//ListLicense list license
func ListLicense() (map[string]*api_model.LicenseInfo, error) {
	licenses, err := db.GetManager().LicenseDao().ListLicenses()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	LMlicense := make(map[string]*api_model.LicenseInfo)
	for _, license := range licenses {
		mLicense := &api_model.LicenseInfo{}
		lc, err := GetLicenseHandler().PackLicense(license.License)
		if err != nil {
			logrus.Errorf("init license error.")
			continue
		}
		if err := ffjson.Unmarshal(lc, mLicense); err != nil {
			logrus.Errorf("unmashal license error, %v", err)
			continue
		}
		LMlicense[license.Token] = mLicense
	}
	return LMlicense, nil
}
