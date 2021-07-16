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

package region

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/node/api/model"
	utilhttp "github.com/gridworkz/kato/util/http"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

//ClusterInterface cluster api
type MonitorInterface interface {
	GetRule(name string) (*model.AlertingNameConfig, *util.APIHandleError)
	GetAllRule() (*model.AlertingRulesConfig, *util.APIHandleError)
	DelRule(name string) (*utilhttp.ResponseBody, *util.APIHandleError)
	AddRule(path string) (*utilhttp.ResponseBody, *util.APIHandleError)
	RegRule(ruleName string, path string) (*utilhttp.ResponseBody, *util.APIHandleError)
}

func (r *regionImpl) Monitor() MonitorInterface {
	return &monitor{prefix: "/v2/rules", regionImpl: *r}
}

type monitor struct {
	regionImpl
	prefix string
}

func (m *monitor) GetRule(name string) (*model.AlertingNameConfig, *util.APIHandleError) {
	var ac model.AlertingNameConfig
	var decode utilhttp.ResponseBody
	decode.Bean = &ac
	code, err := m.DoRequest(m.prefix+"/"+name, "GET", nil, &decode)
	if err != nil {
		return nil, handleErrAndCode(err, code)
	}
	if code != 200 {
		logrus.Error("Return failure message ", decode.Bean)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("get alerting rules error code %d", code))
	}
	return &ac, nil
}

func (m *monitor) GetAllRule() (*model.AlertingRulesConfig, *util.APIHandleError) {
	var ac model.AlertingRulesConfig
	var decode utilhttp.ResponseBody
	decode.Bean = &ac
	code, err := m.DoRequest(m.prefix+"/all", "GET", nil, &decode)
	if err != nil {
		return nil, handleErrAndCode(err, code)
	}
	if code != 200 {
		logrus.Error("Return failure message ", decode.Bean)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("get alerting rules error code %d", code))
	}
	return &ac, nil
}

func (m *monitor) DelRule(name string) (*utilhttp.ResponseBody, *util.APIHandleError) {
	var decode utilhttp.ResponseBody
	code, err := m.DoRequest(m.prefix+"/"+name, "DELETE", nil, &decode)
	if err != nil {
		return nil, handleErrAndCode(err, code)
	}
	if code != 200 {
		logrus.Error("Return failure message ", decode.Bean)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("del alerting rules error code %d", code))
	}
	return &decode, nil
}

func (m *monitor) AddRule(path string) (*utilhttp.ResponseBody, *util.APIHandleError) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			return nil, util.CreateAPIHandleError(400, errors.New("file does not exist"))
		}
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Error("Failed to read AlertingRules config file: ", err.Error())
		return nil, util.CreateAPIHandleError(400, err)
	}
	var rulesConfig model.AlertingNameConfig
	if err := yaml.Unmarshal(content, &rulesConfig); err != nil {
		logrus.Error("Unmarshal AlertingRulesConfig config string to object error.", err.Error())
		return nil, util.CreateAPIHandleError(400, err)

	}
	var decode utilhttp.ResponseBody
	body, err := json.Marshal(rulesConfig)
	if err != nil {
		return nil, util.CreateAPIHandleError(400, err)
	}
	code, err := m.DoRequest(m.prefix, "POST", bytes.NewBuffer(body), &decode)
	if err != nil {
		return nil, handleErrAndCode(err, code)
	}
	if code != 200 {
		logrus.Error("Return failure message ", decode.Bean)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("add alerting rules error code %d", code))
	}
	return &decode, nil
}

func (m *monitor) RegRule(ruleName string, path string) (*utilhttp.ResponseBody, *util.APIHandleError) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			return nil, util.CreateAPIHandleError(400, errors.New("file does not exist"))
		}
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Error("Failed to read AlertingRules config file: ", err.Error())
		return nil, util.CreateAPIHandleError(400, err)
	}
	var rulesConfig model.AlertingNameConfig
	if err := yaml.Unmarshal(content, &rulesConfig); err != nil {
		logrus.Error("Unmarshal AlertingRulesConfig config string to object error.", err.Error())
		return nil, util.CreateAPIHandleError(400, err)

	}
	var decode utilhttp.ResponseBody
	body, err := json.Marshal(rulesConfig)
	if err != nil {
		return nil, util.CreateAPIHandleError(400, err)
	}
	code, err := m.DoRequest(m.prefix+"/"+ruleName, "PUT", bytes.NewBuffer(body), &decode)
	if err != nil {
		return nil, handleErrAndCode(err, code)
	}
	if code != 200 {
		logrus.Error("Return failure message ", decode.Bean)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("add alerting rules error code %d", code))
	}
	return &decode, nil
}
