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
	"strings"

	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/util"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/event"

	"github.com/docker/docker/api/types"
	"github.com/pquerna/ffjson/ffjson"

	"github.com/gridworkz/kato/builder/model"

	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/sirupsen/logrus"
)

const (
	cloneTimeout    = 60
	buildingTimeout = 180
	formatSourceDir = "/cache/build/%s/source/%s"
)

func (e *exectorManager) pluginDockerfileBuild(task *pb.TaskMessage) {
	var tb model.BuildPluginTaskBody
	if err := ffjson.Unmarshal(task.TaskBody, &tb); err != nil {
		logrus.Errorf("unmarshal taskbody error, %v", err)
		return
	}
	eventID := tb.EventID
	logger := event.GetManager().GetLogger(eventID)
	logger.Info("Start execution from the dockerfile build plugin task", map[string]string{"step": "builder-exector", "status": "starting"})
	logrus.Info("start exec build plugin from image worker")
	defer event.GetManager().ReleaseLogger(logger)
	for retry := 0; retry < 2; retry++ {
		err := e.runD(&tb, logger)
		if err != nil {
			logrus.Errorf("exec plugin build from dockerfile error:%s", err.Error())
			logger.Info("Dockerfile build plugin task execution failed, start to retry", map[string]string{"step": "builder-exector", "status": "failure"})
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
	logger.Error("Dockerfile build plugin task execution failed", map[string]string{"step": "callback", "status": "failure"})
}

func (e *exectorManager) runD(t *model.BuildPluginTaskBody, logger event.Logger) error {
	logger.Info("Start pulling code", map[string]string{"step": "build-exector"})
	sourceDir := fmt.Sprintf(formatSourceDir, t.TenantID, t.VersionID)
	if t.Repo == "" {
		t.Repo = "master"
	}
	if !util.DirIsEmpty(sourceDir) {
		os.RemoveAll(sourceDir)
	}
	if err := util.CheckAndCreateDir(sourceDir); err != nil {
		return err
	}
	if _, err := sources.GitClone(sources.CodeSourceInfo{RepositoryURL: t.GitURL, Branch: t.Repo, User: t.GitUsername, Password: t.GitPassword}, sourceDir, logger, 4); err != nil {
		logger.Error("Failed to pull code", map[string]string{"step": "builder-exector", "status": "failure"})
		logrus.Errorf("[plugin]git clone code error %v", err)
		return err
	}
	if !checkDockerfile(sourceDir) {
		logger.Error("The code does not detect the dockerfile, does not support construction temporarily, the task will exit soon", map[string]string{"step": "builder-exector", "status": "failure"})
		logrus.Error("The code does not detect the dockerfile")
		return fmt.Errorf("have no dockerfile")
	}

	logger.Info("The code is detected as a dockerfile and start to compile", map[string]string{"step": "build-exector"})
	mm := strings.Split(t.GitURL, "/")
	n1 := strings.Split(mm[len(mm)-1], ".")[0]
	buildImageName := fmt.Sprintf(builder.REGISTRYDOMAIN+"/plugin_%s_%s:%s", n1, t.PluginID, t.DeployVersion)
	buildOptions := types.ImageBuildOptions{
		Tags:   []string{buildImageName},
		Remove: true,
	}
	if noCache := os.Getenv("NO_CACHE"); noCache != "" {
		buildOptions.NoCache = true
	} else {
		buildOptions.NoCache = false
	}
	logger.Info("start build image", map[string]string{"step": "builder-exector"})
	_, err := sources.ImageBuild(e.DockerClient, sourceDir, buildOptions, logger, 5)
	if err != nil {
		logger.Error(fmt.Sprintf("build image %s failure,find log in rbd-chaos", buildImageName), map[string]string{"step": "builder-exector", "status": "failure"})
		logrus.Errorf("[plugin]build image error: %s", err.Error())
		return err
	}
	logger.Info("build image success, start to push image to local image registry", map[string]string{"step": "builder-exector"})
	err = sources.ImagePush(e.DockerClient, buildImageName, builder.REGISTRYUSER, builder.REGISTRYPASS, logger, 2)
	if err != nil {
		logger.Error("push image failure, find log in rbd-chaos", map[string]string{"step": "builder-exector"})
		logrus.Errorf("push image error: %s", err.Error())
		return err
	}
	logger.Info("push image success", map[string]string{"step": "build-exector"})
	version, err := db.GetManager().TenantPluginBuildVersionDao().GetBuildVersionByDeployVersion(t.PluginID, t.VersionID, t.DeployVersion)
	if err != nil {
		logrus.Errorf("get version error, %v", err)
		return err
	}
	version.BuildLocalImage = buildImageName
	version.Status = "complete"
	if err := db.GetManager().TenantPluginBuildVersionDao().UpdateModel(version); err != nil {
		logrus.Errorf("update version error, %v", err)
		return err
	}
	logger.Info("build plugin version by dockerfile success", map[string]string{"step": "last", "status": "success"})
	return nil
}

func checkDockerfile(sourceDir string) bool {
	if _, err := os.Stat(fmt.Sprintf("%s/Dockerfile", sourceDir)); os.IsNotExist(err) {
		return false
	}
	return true
}
