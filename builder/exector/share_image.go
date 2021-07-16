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

	"github.com/gridworkz/kato/builder"

	"github.com/coreos/etcd/clientv3"
	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/event"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
)

//ImageShareItem
type ImageShareItem struct {
	Namespace          string `json:"namespace"`
	TenantName         string `json:"tenant_name"`
	ServiceID          string `json:"service_id"`
	ServiceAlias       string `json:"service_alias"`
	ImageName          string `json:"image_name"`
	LocalImageName     string `json:"local_image_name"`
	LocalImageUsername string `json:"-"`
	LocalImagePassword string `json:"-"`
	ShareID            string `json:"share_id"`
	Logger             event.Logger
	ShareInfo          struct {
		ServiceKey string `json:"service_key" `
		AppVersion string `json:"app_version" `
		EventID    string `json:"event_id"`
		ShareUser  string `json:"share_user"`
		ShareScope string `json:"share_scope"`
		ImageInfo  struct {
			HubURL      string `json:"hub_url"`
			HubUser     string `json:"hub_user"`
			HubPassword string `json:"hub_password"`
			Namespace   string `json:"namespace"`
			IsTrust     bool   `json:"is_trust,omitempty"`
		} `json:"image_info,omitempty"`
	} `json:"share_info"`
	DockerClient *client.Client
	EtcdCli      *clientv3.Client
}

//NewImageShareItem
func NewImageShareItem(in []byte, DockerClient *client.Client, EtcdCli *clientv3.Client) (*ImageShareItem, error) {
	var isi ImageShareItem
	if err := ffjson.Unmarshal(in, &isi); err != nil {
		return nil, err
	}
	isi.LocalImageUsername = builder.REGISTRYUSER
	isi.LocalImagePassword = builder.REGISTRYPASS
	eventID := isi.ShareInfo.EventID
	isi.Logger = event.GetManager().GetLogger(eventID)
	isi.DockerClient = DockerClient
	isi.EtcdCli = EtcdCli
	return &isi, nil
}

//ShareService
func (i *ImageShareItem) ShareService() error {
	hubuser, hubpass := builder.GetImageUserInfoV2(i.LocalImageName, i.LocalImageUsername, i.LocalImagePassword)
	_, err := sources.ImagePull(i.DockerClient, i.LocalImageName, hubuser, hubpass, i.Logger, 20)
	if err != nil {
		logrus.Errorf("pull image %s error: %s", i.LocalImageName, err.Error())
		i.Logger.Error(fmt.Sprintf("Pull application image: %s failure", i.LocalImageName), map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	if err := sources.ImageTag(i.DockerClient, i.LocalImageName, i.ImageName, i.Logger, 1); err != nil {
		logrus.Errorf("change image tag error: %s", err.Error())
		i.Logger.Error(fmt.Sprintf("Modify the mirror tag: %s -> %s failure", i.LocalImageName, i.ImageName), map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	user, pass := builder.GetImageUserInfoV2(i.ImageName, i.ShareInfo.ImageInfo.HubUser, i.ShareInfo.ImageInfo.HubPassword)
	if i.ShareInfo.ImageInfo.IsTrust {
		err = sources.TrustedImagePush(i.DockerClient, i.ImageName, user, pass, i.Logger, 30)
	} else {
		err = sources.ImagePush(i.DockerClient, i.ImageName, user, pass, i.Logger, 30)
	}
	if err != nil {
		if err.Error() == "authentication required" {
			i.Logger.Error("Mirror warehouse authorization failed", map[string]string{"step": "builder-exector", "status": "failure"})
			return err
		}
		logrus.Errorf("push image into registry error: %s", err.Error())
		i.Logger.Error("Failed to push the image to the mirror warehouse", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	return nil
}

//ShareStatus
type ShareStatus struct {
	ShareID string `json:"share_id,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (s ShareStatus) String() string {
	b, _ := ffjson.Marshal(s)
	return string(b)
}

//UpdateShareStatus
func (i *ImageShareItem) UpdateShareStatus(status string) error {
	var ss = ShareStatus{
		ShareID: i.ShareID,
		Status:  status,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := i.EtcdCli.Put(ctx, fmt.Sprintf("/kato/shareresult/%s", i.ShareID), ss.String())
	if err != nil {
		logrus.Errorf("put shareresult  %s into etcd error, %v", i.ShareID, err)
		i.Logger.Error("Failed to save sharing results.", map[string]string{"step": "callback", "status": "failure"})
	}
	if status == "success" {
		i.Logger.Info("Create and share results are successful, share successfully", map[string]string{"step": "last", "status": "success"})
	} else {
		i.Logger.Info("Create and share results succeeded, share failed", map[string]string{"step": "callback", "status": "failure"})
	}
	return nil
}
