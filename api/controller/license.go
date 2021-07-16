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

package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gridworkz/kato/api/handler"

	httputil "github.com/gridworkz/kato/util/http"

	validator "github.com/gridworkz/kato/util/govalidator"
	"github.com/sirupsen/logrus"
)

//LicenseManager
type LicenseManager struct{}

var licenseManager *LicenseManager

//GetLicenseManager
func GetLicenseManager() *LicenseManager {
	if licenseManager != nil {
		return licenseManager
	}
	licenseManager = &LicenseManager{}
	return licenseManager
}

//AnalystLicense
// swagger:operation POST /license license SendLicense
//
// Submit license
//
// post license & get token
//
// ---
// produces:
// - application/json
// - application/xml
// parameters:
// - name: license
//   in: form
//   description: license
//   required: true
//   type: string
//   format: string
//
// Responses:
//   '200':
//	   description: '{"bean":"{\"token\": \"Q3E5OXdoZDZDX3drN0QtV2gtVmpRaGtlcHJQYmFK\"}"}'
func (l *LicenseManager) AnalystLicense(w http.ResponseWriter, r *http.Request) {
	rule := validator.MapData{
		"license": []string{"required"},
	}
	data, ok := httputil.ValidatorRequestMapAndErrorResponse(r, w, rule, nil)
	if !ok {
		return
	}
	license := data["license"].(string)
	logrus.Debugf("license is %v", license)
	/*
		text, err := handler.GetLicenseHandler().PackLicense(license)
		if err != nil {
			httputil.ReturnError(r, w, 500, fmt.Sprintf("%v", err))
			return
		}
	*/
	token, errT := handler.BasePack([]byte(license))
	if errT != nil {
		httputil.ReturnError(r, w, 500, "pack license error")
		return
	}
	logrus.Debugf("token is %v", token)
	if err := handler.GetLicenseHandler().StoreLicense(license, token); err != nil {
		logrus.Debugf("%s", err)
		logrus.Debugf("%s", fmt.Errorf("license exists"))
		if err == errors.New("license is exist") {
			//err  license is exist
			httputil.ReturnError(r, w, 400, fmt.Sprintf("storage token error, %v", err))
		}
		httputil.ReturnError(r, w, 500, fmt.Sprintf("storage token error, %v", err))
		return
	}
	rc := fmt.Sprintf(`{"token": "%v"}`, token)
	httputil.ReturnSuccess(r, w, &rc)
}
