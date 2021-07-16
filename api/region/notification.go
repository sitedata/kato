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
	"fmt"

	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/node/api/model"
	utilhttp "github.com/gridworkz/kato/util/http"
	"github.com/sirupsen/logrus"
)

//NotificationInterface cluster api
type NotificationInterface interface {
	GetNotification(start string, end string) ([]*model.NotificationEvent, *util.APIHandleError)
	HandleNotification(serviceName string, message string) ([]*model.NotificationEvent, *util.APIHandleError)
}

func (r *regionImpl) Notification() NotificationInterface {
	return &notification{prefix: "/v2/notificationEvent", regionImpl: *r}
}

type notification struct {
	regionImpl
	prefix string
}

func (n *notification) GetNotification(start string, end string) ([]*model.NotificationEvent, *util.APIHandleError) {
	var ne []*model.NotificationEvent
	var decode utilhttp.ResponseBody
	decode.List = &ne
	code, err := n.DoRequest(n.prefix+"?start="+start+"&"+"end="+end, "GET", nil, &decode)
	if err != nil {
		return nil, handleErrAndCode(err, code)
	}
	if code != 200 {
		logrus.Error("Return failure message ", decode.Msg)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf(decode.Msg))
	}
	return ne, nil
}

func (n *notification) HandleNotification(serviceName string, message string) ([]*model.NotificationEvent, *util.APIHandleError) {
	var ne []*model.NotificationEvent
	var decode utilhttp.ResponseBody
	decode.List = &ne
	handleMessage, err := json.Marshal(map[string]string{"handle_message": message})
	body := bytes.NewBuffer(handleMessage)
	code, err := n.DoRequest(n.prefix+"/"+serviceName, "PUT", body, &decode)
	if err != nil {
		return nil, handleErrAndCode(err, code)
	}
	if code != 200 {
		logrus.Error("Return failure message ", decode.Msg)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf(decode.Msg))
	}
	return ne, nil
}
