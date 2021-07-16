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

	"github.com/gridworkz/kato/builder/exector"

	"github.com/twinj/uuid"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/gridworkz/kato/db"

	"github.com/coreos/etcd/clientv3"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/sirupsen/logrus"
)

//ServiceShareHandle service share
type ServiceShareHandle struct {
	MQClient client.MQClient
	EtcdCli  *clientv3.Client
}

//APIResult - share interface return
type APIResult struct {
	EventID   string `json:"event_id"`
	ShareID   string `json:"share_id"`
	ImageName string `json:"image_name,omitempty"`
	SlugPath  string `json:"slug_path,omitempty"`
}

//Share - share app
func (s *ServiceShareHandle) Share(serviceID string, ss api_model.ServiceShare) (*APIResult, *util.APIHandleError) {
	service, err := db.GetManager().TenantServiceDao().GetServiceByID(serviceID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("查询应用出错", err)
	}
	//query deployment version
	version, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(service.DeployVersion, serviceID)
	if err != nil {
		logrus.Error("query service deploy version error", err.Error())
	}
	shareID := uuid.NewV4().String()
	var slugPath, shareImageName string
	var task client.TaskStruct
	if version.DeliveredType == "slug" {
		shareSlugInfo := ss.Body.SlugInfo
		slugPath = service.CreateShareSlug(ss.Body.ServiceKey, shareSlugInfo.Namespace, ss.Body.AppVersion)
		if ss.Body.SlugInfo.FTPHost == "" {
			slugPath = fmt.Sprintf("/grdata/build/tenant/%s", slugPath)
		}
		info := map[string]interface{}{
			"service_alias": ss.ServiceAlias,
			"service_id":    serviceID,
			"tenant_name":   ss.TenantName,
			"share_info":    ss.Body,
			"slug_path":     slugPath,
			"share_id":      shareID,
		}
		if version != nil && version.DeliveredPath != "" {
			info["local_slug_path"] = version.DeliveredPath
		} else {
			info["local_slug_path"] = fmt.Sprintf("/grdata/build/tenant/%s/slug/%s/%s.tgz", service.TenantID, service.ServiceID, service.DeployVersion)
		}
		task.TaskType = "share-slug"
		task.TaskBody = info
	} else {
		shareImageInfo := ss.Body.ImageInfo
		shareImageName, err = version.CreateShareImage(shareImageInfo.HubURL, shareImageInfo.Namespace, ss.Body.AppVersion)
		if err != nil {
			return nil, util.CreateAPIHandleError(500, err)
		}
		info := map[string]interface{}{
			"share_info":    ss.Body,
			"service_alias": ss.ServiceAlias,
			"service_id":    serviceID,
			"tenant_name":   ss.TenantName,
			"image_name":    shareImageName,
			"share_id":      shareID,
		}
		if version != nil && version.DeliveredPath != "" {
			info["local_image_name"] = version.DeliveredPath
		}
		task.TaskType = "share-image"
		task.TaskBody = info
	}
	label, err := db.GetManager().TenantServiceLabelDao().GetLabelByNodeSelectorKey(serviceID, "windows")
	if label == nil || err != nil {
		task.Topic = client.BuilderTopic
	} else {
		task.Topic = client.WindowsBuilderTopic
	}
	err = s.MQClient.SendBuilderTopic(task)
	if err != nil {
		logrus.Errorf("equque mq error, %v", err)
		return nil, util.CreateAPIHandleError(502, err)
	}
	return &APIResult{EventID: ss.Body.EventID, ShareID: shareID, ImageName: shareImageName, SlugPath: slugPath}, nil
}

//ShareResult - share application results query
func (s *ServiceShareHandle) ShareResult(shareID string) (i exector.ShareStatus, e *util.APIHandleError) {
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
