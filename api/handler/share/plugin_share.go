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

package share

import (
	"context"
	"fmt"

	"github.com/gridworkz/kato/mq/client"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/builder/exector"
	"github.com/gridworkz/kato/db"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
)

//PluginShareHandle plugin share
type PluginShareHandle struct {
	MQClient client.MQClient
	EtcdCli  *clientv3.Client
}

//PluginResult share plugin api return
type PluginResult struct {
	EventID   string `json:"event_id"`
	ShareID   string `json:"share_id"`
	ImageName string `json:"image_name"`
}

//PluginShare
type PluginShare struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	TenantID   string
	// in: path
	// required: true
	PluginID string `json:"plugin_id"`
	//in: body
	Body struct {
		//in: body
		//application Sharing Key
		PluginKey     string `json:"plugin_key" validate:"plugin_key|required"`
		PluginVersion string `json:"plugin_version" validate:"plugin_version|required"`
		EventID       string `json:"event_id"`
		ShareUser     string `json:"share_user"`
		ShareScope    string `json:"share_scope"`
		ImageInfo     struct {
			HubURL      string `json:"hub_url"`
			HubUser     string `json:"hub_user"`
			HubPassword string `json:"hub_password"`
			Namespace   string `json:"namespace"`
			IsTrust     bool   `json:"is_trust,omitempty" validate:"is_trust"`
		} `json:"image_info,omitempty"`
	}
}

//Share share app
func (s *PluginShareHandle) Share(ss PluginShare) (*PluginResult, *util.APIHandleError) {
	_, err := db.GetManager().TenantPluginDao().GetPluginByID(ss.PluginID, ss.TenantID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("query plugin error", err)
	}
	//query new build version
	version, err := db.GetManager().TenantPluginBuildVersionDao().GetLastBuildVersionByVersionID(ss.PluginID, ss.Body.PluginVersion)
	if err != nil {
		logrus.Error("query service deploy version error", err.Error())
		return nil, util.CreateAPIHandleErrorFromDBError("query plugin version error", err)
	}
	shareID := uuid.NewV4().String()
	shareImageName, err := version.CreateShareImage(ss.Body.ImageInfo.HubURL, ss.Body.ImageInfo.Namespace)
	if err != nil {
		return nil, util.CreateAPIHandleErrorf(500, "create share image name error:%s", err.Error())
	}
	info := map[string]interface{}{
		"image_info":       ss.Body.ImageInfo,
		"event_id":         ss.Body.EventID,
		"tenant_name":      ss.TenantName,
		"image_name":       shareImageName,
		"share_id":         shareID,
		"local_image_name": version.BuildLocalImage,
	}
	err = s.MQClient.SendBuilderTopic(client.TaskStruct{
		TaskType: "share-plugin",
		TaskBody: info,
		Topic:    client.BuilderTopic,
	})
	if err != nil {
		logrus.Errorf("equque mq error, %v", err)
		return nil, util.CreateAPIHandleError(502, err)
	}
	return &PluginResult{EventID: ss.Body.EventID, ShareID: shareID, ImageName: shareImageName}, nil
}

//ShareResult - share application results query
func (s *PluginShareHandle) ShareResult(shareID string) (i exector.ShareStatus, e *util.APIHandleError) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := s.EtcdCli.Get(ctx, fmt.Sprintf("/kato/shareresult/%s", shareID))
	if err != nil {
		e = util.CreateAPIHandleError(500, err)
	} else {
		if res.Count == 0 {
			i.ShareID = shareID
		} else {
			if err := ffjson.Unmarshal(res.Kvs[0].Value, &i); err != nil {
				return i, util.CreateAPIHandleError(500, err)
			}
		}
	}
	return
}
