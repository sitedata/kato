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
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/build"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/event"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

//ImageBuildItem ImageBuildItem
type ImageBuildItem struct {
	Namespace     string       `json:"namespace"`
	TenantName    string       `json:"tenant_name"`
	ServiceAlias  string       `json:"service_alias"`
	Image         string       `json:"image"`
	DestImage     string       `json:"dest_image"`
	Logger        event.Logger `json:"logger"`
	EventID       string       `json:"event_id"`
	DockerClient  *client.Client
	TenantID      string
	ServiceID     string
	DeployVersion string
	HubUser       string
	HubPassword   string
	Action        string
	Configs       map[string]gjson.Result `json:"configs"`
}

//NewImageBuildItem
func NewImageBuildItem(in []byte) *ImageBuildItem {
	eventID := gjson.GetBytes(in, "event_id").String()
	logger := event.GetManager().GetLogger(eventID)
	return &ImageBuildItem{
		Namespace:     gjson.GetBytes(in, "namespace").String(),
		TenantName:    gjson.GetBytes(in, "tenant_name").String(),
		ServiceAlias:  gjson.GetBytes(in, "service_alias").String(),
		ServiceID:     gjson.GetBytes(in, "service_id").String(),
		Image:         gjson.GetBytes(in, "image").String(),
		DeployVersion: gjson.GetBytes(in, "deploy_version").String(),
		Action:        gjson.GetBytes(in, "action").String(),
		HubUser:       gjson.GetBytes(in, "user").String(),
		HubPassword:   gjson.GetBytes(in, "password").String(),
		Configs:       gjson.GetBytes(in, "configs").Map(),
		Logger:        logger,
		EventID:       eventID,
	}
}

//Run
func (i *ImageBuildItem) Run(timeout time.Duration) error {
	user, pass := builder.GetImageUserInfoV2(i.Image, i.HubUser, i.HubPassword)
	_, err := sources.ImagePull(i.DockerClient, i.Image, user, pass, i.Logger, 30)
	if err != nil {
		logrus.Errorf("pull image %s error: %s", i.Image, err.Error())
		i.Logger.Error(fmt.Sprintf("get specified image: %s failure", i.Image), map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	localImageURL := build.CreateImageName(i.ServiceID, i.DeployVersion)
	if err := sources.ImageTag(i.DockerClient, i.Image, localImageURL, i.Logger, 1); err != nil {
		logrus.Errorf("change image tag error: %s", err.Error())
		i.Logger.Error(fmt.Sprintf("modify mirror tag: %s -> %s failure", i.Image, localImageURL), map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	err = sources.ImagePush(i.DockerClient, localImageURL, builder.REGISTRYUSER, builder.REGISTRYPASS, i.Logger, 30)
	if err != nil {
		logrus.Errorf("push image into registry error: %s", err.Error())
		i.Logger.Error("failed to push the image to the mirror warehouse", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}

	if err := sources.ImageRemove(i.DockerClient, localImageURL); err != nil {
		logrus.Errorf("remove image %s failure %s", localImageURL, err.Error())
	}

	if os.Getenv("DISABLE_IMAGE_CACHE") == "true" {
		if err := sources.ImageRemove(i.DockerClient, i.Image); err != nil {
			logrus.Errorf("remove image %s failure %s", i.Image, err.Error())
		}
	}
	if err := i.StorageVersionInfo(localImageURL); err != nil {
		logrus.Errorf("storage version info error, ignor it: %s", err.Error())
		i.Logger.Error("failed to update app version information", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	return nil
}

//StorageVersionInfo
func (i *ImageBuildItem) StorageVersionInfo(imageURL string) error {
	version, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(i.DeployVersion, i.ServiceID)
	if err != nil {
		return err
	}
	version.DeliveredType = "image"
	version.DeliveredPath = imageURL
	version.ImageName = imageURL
	version.RepoURL = i.Image
	version.FinalStatus = "success"
	version.FinishTime = time.Now()
	if err := db.GetManager().VersionInfoDao().UpdateModel(version); err != nil {
		return err
	}
	return nil
}

//UpdateVersionInfo
func (i *ImageBuildItem) UpdateVersionInfo(status string) error {
	version, err := db.GetManager().VersionInfoDao().GetVersionByEventID(i.EventID)
	if err != nil {
		return err
	}
	version.FinalStatus = status
	version.RepoURL = i.Image
	version.FinishTime = time.Now()
	if err := db.GetManager().VersionInfoDao().UpdateModel(version); err != nil {
		return err
	}
	return nil
}
