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

package exector

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/ghodss/yaml"

	"github.com/gridworkz/kato/builder/parser"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
)

//ServiceCheckInput Task input data
type ServiceCheckInput struct {
	CheckUUID string `json:"uuid"`
	//Detection source type
	SourceType string `json:"source_type"`

	// Definition of detection source,
	// Code: https://github.com/shurcooL/githubql.git master
	// docker-run: docker run --name xxx nginx:latest nginx
	// docker-compose: compose full text
	SourceBody string `json:"source_body"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	TenantID   string
	EventID    string `json:"event_id"`
}

//ServiceCheckResult Application test results
type ServiceCheckResult struct {
	//Detection status Success Failure
	CheckStatus string                `json:"check_status"`
	ErrorInfos  parser.ParseErrorList `json:"error_infos"`
	ServiceInfo []parser.ServiceInfo  `json:"service_info"`
}

//CreateResult Create test results
func CreateResult(ErrorInfos parser.ParseErrorList, ServiceInfo []parser.ServiceInfo) (ServiceCheckResult, error) {
	var sr ServiceCheckResult
	if ErrorInfos != nil && ErrorInfos.IsFatalError() {
		sr = ServiceCheckResult{
			CheckStatus: "Failure",
			ErrorInfos:  ErrorInfos,
			ServiceInfo: ServiceInfo,
		}
	} else {
		sr = ServiceCheckResult{
			CheckStatus: "Success",
			ErrorInfos:  ErrorInfos,
			ServiceInfo: ServiceInfo,
		}
	}
	//save result
	return sr, nil
}

//serviceCheck Application creation source detection
func (e *exectorManager) serviceCheck(task *pb.TaskMessage) {
	//step1 Determine the application source type
	//step2 Obtain application source media, mirror or source code
	//step3 Analyze and judge application source specifications
	//Finish
	var input ServiceCheckInput
	if err := ffjson.Unmarshal(task.TaskBody, &input); err != nil {
		logrus.Error("Unmarshal service check input data error.", err.Error())
		return
	}
	logger := event.GetManager().GetLogger(input.EventID)
	defer event.GetManager().ReleaseLogger(logger)
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("service check error: %v", r)
			debug.PrintStack()
			logger.Error("The back-end service has deserted, please try again.", map[string]string{"step": "callback", "status": "failure"})
		}
	}()
	logger.Info("Start component deploy source check.", map[string]string{"step": "starting"})
	logrus.Infof("start check service by type: %s ", input.SourceType)
	var pr parser.Parser
	switch input.SourceType {
	case "docker-run":
		pr = parser.CreateDockerRunOrImageParse(input.Username, input.Password, input.SourceBody, e.DockerClient, logger)
	case "docker-compose":
		var yamlbody = input.SourceBody
		if input.SourceBody[0] == '{' {
			yamlbyte, err := yaml.JSONToYAML([]byte(input.SourceBody))
			if err != nil {
				logrus.Errorf("json bytes format is error, %s", input.SourceBody)
				logger.Error("The dockercompose file is not in the correct format", map[string]string{"step": "callback", "status": "failure"})
				return
			}
			yamlbody = string(yamlbyte)
		}
		pr = parser.CreateDockerComposeParse(yamlbody, e.DockerClient, input.Username, input.Password, logger)
	case "sourcecode":
		pr = parser.CreateSourceCodeParse(input.SourceBody, logger)
	case "third-party-service":
		pr = parser.CreateThirdPartyServiceParse(input.SourceBody, logger)
	}
	if pr == nil {
		logger.Error("Creating component source types is not supported", map[string]string{"step": "callback", "status": "failure"})
		return
	}
	errList := pr.Parse()
	for i, err := range errList {
		if err.SolveAdvice == "" && input.SourceType != "sourcecode" {
			errList[i].SolveAdvice = fmt.Sprintf("The parser thinks that the image name is: %s, Please confirm whether it is correct or whether the mirror exists", pr.GetImage())
		}
		if err.SolveAdvice == "" && input.SourceType == "sourcecode" {
			errList[i].SolveAdvice = "Source code intelligent parsing failed"
		}
	}
	serviceInfos := pr.GetServiceInfo()
	sr, err := CreateResult(errList, serviceInfos)
	if err != nil {
		logrus.Errorf("create check result error,%s", err.Error())
		logger.Error("Failed to create test result.", map[string]string{"step": "callback", "status": "failure"})
	}
	k := fmt.Sprintf("/servicecheck/%s", input.CheckUUID)
	v := sr
	vj, err := ffjson.Marshal(&v)
	if err != nil {
		logrus.Errorf("mashal servicecheck value error, %v", err)
		logger.Error("Failed to format test result.", map[string]string{"step": "callback", "status": "failure"})
	}
	ctx, cancel := context.WithCancel(context.Background())
	_, err = e.EtcdCli.Put(ctx, k, string(vj))
	cancel()
	if err != nil {
		logrus.Errorf("put servicecheck k %s into etcd error, %v", k, err)
		logger.Error("Failed to store test results.", map[string]string{"step": "callback", "status": "failure"})
	}
	logrus.Infof("check service by type: %s  success", input.SourceType)
	logger.Info("Successfully created the test result。", map[string]string{"step": "last", "status": "success"})
}
