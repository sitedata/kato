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
	"strings"

	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/model"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/pquerna/ffjson/ffjson"

	"github.com/sirupsen/logrus"
)

func (e *exectorManager) pluginImageBuild(task *pb.TaskMessage) {
	var tb model.BuildPluginTaskBody
	if err := ffjson.Unmarshal(task.TaskBody, &tb); err != nil {
		logrus.Errorf("unmarshal taskbody error, %v", err)
		return
	}
	eventID := tb.EventID
	logger := event.GetManager().GetLogger(eventID)
	logger.Info("Start execution from the image build plugin task", map[string]string{"step": "builder-exector", "status": "starting"})

	logrus.Info("start exec build plugin from image worker")
	defer event.GetManager().ReleaseLogger(logger)
	for retry := 0; retry < 2; retry++ {
		err := e.run(&tb, logger)
		if err != nil {
			logrus.Errorf("exec plugin build from image error:%s", err.Error())
			logger.Info("Mirror build plugin task execution failed, start to retry", map[string]string{"step": "builder-exector", "status": "failure"})
		} else {
			return
		}
	}
	version, err := db.GetManager().TenantPluginBuildVersionDao().GetBuildVersionByDeployVersion(tb.PluginID, tb.VersionID, tb.DeployVersion)
	if err != nil {
		logrus.Errorf("get version error, %v", err)
		return
	}
	version.Status = "failure"
	if err := db.GetManager().TenantPluginBuildVersionDao().UpdateModel(version); err != nil {
		logrus.Errorf("update version error, %v", err)
	}
	MetricErrorTaskNum++
	logger.Info("Failed to execute the image build plugin task", map[string]string{"step": "callback", "status": "failure"})
}

func (e *exectorManager) run(t *model.BuildPluginTaskBody, logger event.Logger) error {
	hubUser, hubPass := builder.GetImageUserInfoV2(t.ImageURL, t.ImageInfo.HubUser, t.ImageInfo.HubPassword)
	if _, err := sources.ImagePull(e.DockerClient, t.ImageURL, hubUser, hubPass, logger, 10); err != nil {
		logrus.Errorf("pull image %v error, %v", t.ImageURL, err)
		logger.Error("Failed to pull image", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	logger.Info("Mirror pull completed", map[string]string{"step": "build-exector", "status": "complete"})
	newTag := createPluginImageTag(t.ImageURL, t.PluginID, t.DeployVersion)
	err := sources.ImageTag(e.DockerClient, t.ImageURL, newTag, logger, 1)
	if err != nil {
		logrus.Errorf("set plugin image tag error, %v", err)
		logger.Error("Failed to modify the image tag", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	logger.Info("Modify the mirror tag complete", map[string]string{"step": "build-exector", "status": "complete"})
	if err := sources.ImagePush(e.DockerClient, newTag, builder.REGISTRYUSER, builder.REGISTRYPASS, logger, 10); err != nil {
		logrus.Errorf("push image %s error, %v", newTag, err)
		logger.Error("Failed to push image", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	version, err := db.GetManager().TenantPluginBuildVersionDao().GetBuildVersionByDeployVersion(t.PluginID, t.VersionID, t.DeployVersion)
	if err != nil {
		logger.Error("Update plug-in version information error", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	version.BuildLocalImage = newTag
	version.Status = "complete"
	if err := db.GetManager().TenantPluginBuildVersionDao().UpdateModel(version); err != nil {
		logger.Error("Update plug-in version information error", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	logger.Info("Build the plugin from the mirror image", map[string]string{"step": "last", "status": "success"})
	return nil
}

func createPluginImageTag(image string, pluginid, version string) string {
	//alias is pluginID
	mm := strings.Split(image, "/")
	tag := "latest"
	iName := ""
	if strings.Contains(mm[len(mm)-1], ":") {
		nn := strings.Split(mm[len(mm)-1], ":")
		tag = nn[1]
		iName = nn[0]
	} else {
		iName = image
	}
	if strings.HasPrefix(iName, "plugin") {
		return fmt.Sprintf("%s/%s:%s_%s", builder.REGISTRYDOMAIN, iName, pluginid, version)
	}
	return fmt.Sprintf("%s/plugin_%s_%s:%s_%s", builder.REGISTRYDOMAIN, iName, pluginid, tag, version)
}
