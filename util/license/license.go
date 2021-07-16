// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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

package license

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"strconv"

	"github.com/sirupsen/logrus"
)

var enterprise = "false"

// LicInfo
type LicInfo struct {
	LicKey     string   `json:"license_key"`
	Code       string   `json:"code"`
	Company    string   `json:"company"`
	Node       int64    `json:"node"`
	CPU        int64    `json:"cpu"`
	Memory     int64    `json:"memory"`
	Tenant     int64    `json:"tenant"`
	EndTime    string   `json:"end_time"`
	StartTime  string   `json:"start_time"`
	DataCenter int64    `json:"data_center"`
	ModuleList []string `json:"module_list"`
}

func isEnterprise() bool {
	res, err := strconv.ParseBool(enterprise)
	if err != nil {
		logrus.Warningf("enterprise: %s; error parsing 'string' to 'bool': %v", enterprise, err)
	}
	return res
}

func readFromFile(lfile string) (string, error) {
	_, err := os.Stat(lfile)
	if err != nil {
		logrus.Errorf("license file is incorrect: %v", err)
		return "", err
	}
	bytes, err := ioutil.ReadFile(lfile)
	if err != nil {
		logrus.Errorf("license file: %s; error reading license file: %v", lfile, err)
		return "", err
	}
	return string(bytes), nil
}

// VerifyTime verifies the time in the license.
func VerifyTime(licPath, licSoPath string) bool {
	if !isEnterprise() {
		return true
	}
	lic, err := readFromFile(licPath)
	if err != nil {
		logrus.Errorf("failed to read license from file: %v", err)
		return false
	}
	p, err := plugin.Open(licSoPath)
	if err != nil {
		logrus.Errorf("license.so path: %s; error opening license.so: %v", licSoPath, err)
		return false
	}
	f, err := p.Lookup("VerifyTime")
	if err != nil {
		logrus.Errorf("method 'VerifyTime'; error looking up func: %v", err)
		return false
	}
	return f.(func(string) bool)(lic)
}

// VerifyNodes verifies the number of the nodes in the license.
func VerifyNodes(licPath, licSoPath string, nodeNums int) bool {
	if !isEnterprise() {
		return true
	}
	lic, err := readFromFile(licPath)
	if err != nil {
		logrus.Errorf("failed to read license from file: %v", err)
		return false
	}
	p, err := plugin.Open(licSoPath)
	if err != nil {
		logrus.Errorf("license.so path: %s; error opening license.so: %v", licSoPath, err)
		return false
	}
	f, err := p.Lookup("VerifyNodes")
	if err != nil {
		logrus.Errorf("method 'VerifyNodes'; error looking up func: %v", err)
		return false
	}
	return f.(func(string, int) bool)(lic, nodeNums)
}

// GetLicInfo -
func GetLicInfo(licPath, licSoPath string) (*LicInfo, error) {
	if !isEnterprise() {
		return nil, nil
	}
	lic, err := readFromFile(licPath)
	if err != nil {
		logrus.Errorf("failed to read license from file: %v", err)
		return nil, fmt.Errorf("failed to read license from file: %v", err)
	}
	p, err := plugin.Open(licSoPath)
	if err != nil {
		logrus.Errorf("license.so path: %s; error opening license.so: %v", licSoPath, err)
		return nil, fmt.Errorf("license.so path: %s; error opening license.so: %v", licSoPath, err)
	}

	f, err := p.Lookup("Decrypt")
	if err != nil {
		logrus.Errorf("method 'Decrypt'; error looking up func: %v", err)
		return nil, fmt.Errorf("method 'Decrypt'; error looking up func: %v", err)
	}
	bytes, err := f.(func(string) ([]byte, error))(lic)
	var licInfo LicInfo
	if err := json.Unmarshal(bytes, &licInfo); err != nil {
		logrus.Errorf("error unmarshalling license: %v", err)
		return nil, fmt.Errorf("error unmarshalling license: %v", err)
	}
	return &licInfo, nil
}

// GenLicKey
func GenLicKey(licSoPath string) (string, error) {
	if !isEnterprise() {
		return "", nil
	}
	p, err := plugin.Open(licSoPath)
	if err != nil {
		logrus.Errorf("license.so path: %s; error opening license.so: %v", licSoPath, err)
		return "", fmt.Errorf("license.so path: %s; error opening license.so: %v", licSoPath, err)
	}

	f, err := p.Lookup("GenLicKey")
	if err != nil {
		logrus.Errorf("method 'GenLicKey'; error looking up func: %v", err)
		return "", fmt.Errorf("method 'GenLicKey'; error looking up func: %v", err)
	}
	return f.(func() (string, error))()
}
