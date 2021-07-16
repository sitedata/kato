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
	"context"
	"fmt"
	"strings"
	"time"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/builder/exector"
	client "github.com/gridworkz/kato/mq/client"
	tutil "github.com/gridworkz/kato/util"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
)

//ServiceCheck check service build source
func (s *ServiceAction) ServiceCheck(scs *api_model.ServiceCheckStruct) (string, string, *util.APIHandleError) {
	checkUUID := uuid.NewV4().String()
	scs.Body.CheckUUID = checkUUID
	if scs.Body.EventID == "" {
		scs.Body.EventID = tutil.NewUUID()
	}
	topic := client.BuilderTopic
	if tutil.StringArrayContains(s.conf.EnableFeature, "windows") {
		if scs.Body.CheckOS == "windows" {
			topic = client.WindowsBuilderTopic
		}
		if scs.Body.SourceType == "docker-run" || scs.Body.SourceType == "docker-compose" {
			if maybeIsWindowsContainerImage(scs.Body.SourceBody) {
				topic = client.WindowsBuilderTopic
			}
		}
	}
	err := s.MQClient.SendBuilderTopic(client.TaskStruct{
		TaskType: "service_check",
		TaskBody: scs.Body,
		Topic:    topic,
	})
	if err != nil {
		logrus.Errorf("enqueue service check message to mq error, %v", err)
		return "", "", util.CreateAPIHandleError(500, err)
	}
	return checkUUID, scs.Body.EventID, nil
}

var windowsKeywords = []string{"windows", "asp", "microsoft", "nanoserver"}

func maybeIsWindowsContainerImage(source string) bool {
	for _, k := range windowsKeywords {
		if strings.Contains(source, k) {
			return true
		}
	}
	return false

}

//GetServiceCheckInfo get application source detection information
func (s *ServiceAction) GetServiceCheckInfo(uuid string) (*exector.ServiceCheckResult, *util.APIHandleError) {
	k := fmt.Sprintf("/servicecheck/%s", uuid)
	var si exector.ServiceCheckResult
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := s.EtcdCli.Get(ctx, k)
	if err != nil {
		logrus.Errorf("get etcd k %s error, %v", k, err)
		return nil, util.CreateAPIHandleError(500, err)
	}
	if resp.Count == 0 {
		return &si, nil
	}
	v := resp.Kvs[0].Value
	if err := ffjson.Unmarshal(v, &si); err != nil {
		return nil, util.CreateAPIHandleError(500, err)
	}
	if si.CheckStatus == "" {
		si.CheckStatus = "Checking"
		logrus.Debugf("checking is %v", si)
	}
	return &si, nil
}
