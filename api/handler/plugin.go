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
	"fmt"
	"strings"
	"time"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/mq/client"
	core_util "github.com/gridworkz/kato/util"

	builder_model "github.com/gridworkz/kato/builder/model"

	"github.com/sirupsen/logrus"
)

//PluginAction  plugin action struct
type PluginAction struct {
	MQClient client.MQClient
}

//CreatePluginManager get plugin manager
func CreatePluginManager(mqClient client.MQClient) *PluginAction {
	return &PluginAction{
		MQClient: mqClient,
	}
}

//CreatePluginAct PluginAct
func (p *PluginAction) CreatePluginAct(cps *api_model.CreatePluginStruct) *util.APIHandleError {
	tp := &dbmodel.TenantPlugin{
		TenantID:    cps.Body.TenantID,
		PluginID:    cps.Body.PluginID,
		PluginInfo:  cps.Body.PluginInfo,
		PluginModel: cps.Body.PluginModel,
		PluginName:  cps.Body.PluginName,
		ImageURL:    cps.Body.ImageURL,
		GitURL:      cps.Body.GitURL,
		BuildModel:  cps.Body.BuildModel,
		Domain:      cps.TenantName,
	}
	if err := db.GetManager().TenantPluginDao().AddModel(tp); err != nil {
		return util.CreateAPIHandleErrorFromDBError("create plugin", err)
	}
	return nil
}

//UpdatePluginAct UpdatePluginAct
func (p *PluginAction) UpdatePluginAct(pluginID, tenantID string, cps *api_model.UpdatePluginStruct) *util.APIHandleError {
	tp, err := db.GetManager().TenantPluginDao().GetPluginByID(pluginID, tenantID)
	if err != nil {
		return util.CreateAPIHandleErrorFromDBError("get old plugin info", err)
	}
	tp.PluginInfo = cps.Body.PluginInfo
	tp.PluginModel = cps.Body.PluginModel
	tp.PluginName = cps.Body.PluginName
	tp.ImageURL = cps.Body.ImageURL
	tp.GitURL = cps.Body.GitURL
	tp.BuildModel = cps.Body.BuildModel
	err = db.GetManager().TenantPluginDao().UpdateModel(tp)
	if err != nil {
		return util.CreateAPIHandleErrorFromDBError("update plugin", err)
	}
	return nil
}

//DeletePluginAct DeletePluginAct
func (p *PluginAction) DeletePluginAct(pluginID, tenantID string) *util.APIHandleError {
	tx := db.GetManager().Begin()
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Unexpected panic occurred, rollback transaction: %v", r)
			tx.Rollback()
		}
	}()
	//step1: delete service plugin relation
	err := db.GetManager().TenantServicePluginRelationDaoTransactions(tx).DeleteALLRelationByPluginID(pluginID)
	if err != nil {
		tx.Rollback()
		return util.CreateAPIHandleErrorFromDBError("delete plugin relation", err)
	}
	//step2: delete plugin build version
	err = db.GetManager().TenantPluginBuildVersionDaoTransactions(tx).DeleteBuildVersionByPluginID(pluginID)
	if err != nil {
		tx.Rollback()
		return util.CreateAPIHandleErrorFromDBError("delete plugin build version", err)
	}
	//step3: delete plugin
	err = db.GetManager().TenantPluginDaoTransactions(tx).DeletePluginByID(pluginID, tenantID)
	if err != nil {
		tx.Rollback()
		return util.CreateAPIHandleErrorFromDBError("delete plugin", err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return util.CreateAPIHandleErrorFromDBError("commit delete plugin transactions", err)
	}
	return nil
}

//GetPlugins get all plugins by tenantID
func (p *PluginAction) GetPlugins(tenantID string) ([]*dbmodel.TenantPlugin, *util.APIHandleError) {
	plugins, err := db.GetManager().TenantPluginDao().GetPluginsByTenantID(tenantID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("get plugins by tenant id", err)
	}
	return plugins, nil
}

//AddDefaultEnv AddDefaultEnv
func (p *PluginAction) AddDefaultEnv(est *api_model.ENVStruct) *util.APIHandleError {
	tx := db.GetManager().Begin()
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Unexpected panic occurred, rollback transaction: %v", r)
			tx.Rollback()
		}
	}()
	for _, env := range est.Body.EVNInfo {
		vis := &dbmodel.TenantPluginDefaultENV{
			PluginID:  est.PluginID,
			ENVName:   env.ENVName,
			ENVValue:  env.ENVValue,
			IsChange:  env.IsChange,
			VersionID: env.VersionID,
		}
		err := db.GetManager().TenantPluginDefaultENVDaoTransactions(tx).AddModel(vis)
		if err != nil {
			tx.Rollback()
			return util.CreateAPIHandleErrorFromDBError(fmt.Sprintf("add default env %s", env.ENVName), err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return util.CreateAPIHandleErrorFromDBError("commit add default env transactions", err)
	}
	return nil
}

//UpdateDefaultEnv UpdateDefaultEnv
func (p *PluginAction) UpdateDefaultEnv(est *api_model.ENVStruct) *util.APIHandleError {
	for _, env := range est.Body.EVNInfo {
		vis := &dbmodel.TenantPluginDefaultENV{
			ENVName:   env.ENVName,
			ENVValue:  env.ENVValue,
			IsChange:  env.IsChange,
			VersionID: env.VersionID,
		}
		err := db.GetManager().TenantPluginDefaultENVDao().UpdateModel(vis)
		if err != nil {
			return util.CreateAPIHandleErrorFromDBError(fmt.Sprintf("update default env %s", env.ENVName), err)
		}
	}
	return nil
}

//DeleteDefaultEnv DeleteDefaultEnv
func (p *PluginAction) DeleteDefaultEnv(pluginID, versionID, name string) *util.APIHandleError {
	if err := db.GetManager().TenantPluginDefaultENVDao().DeleteDefaultENVByName(pluginID, name, versionID); err != nil {
		return util.CreateAPIHandleErrorFromDBError(fmt.Sprintf("delete default env %s", name), err)
	}
	return nil
}

//GetDefaultEnv GetDefaultEnv
func (p *PluginAction) GetDefaultEnv(pluginID, versionID string) ([]*dbmodel.TenantPluginDefaultENV, *util.APIHandleError) {
	envs, err := db.GetManager().TenantPluginDefaultENVDao().GetDefaultENVSByPluginID(pluginID, versionID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("get default env", err)
	}
	return envs, nil
}

//GetEnvsWhichCanBeSet GetEnvsWhichCanBeSet
func (p *PluginAction) GetEnvsWhichCanBeSet(serviceID, pluginID string) (interface{}, *util.APIHandleError) {
	relation, err := db.GetManager().TenantServicePluginRelationDao().GetRelateionByServiceIDAndPluginID(serviceID, pluginID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("get relation", err)
	}
	envs, err := db.GetManager().TenantPluginVersionENVDao().GetVersionEnvByServiceID(serviceID, pluginID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("get envs which can be set", err)
	}
	if len(envs) > 0 {
		return envs, nil
	}
	envD, errD := db.GetManager().TenantPluginDefaultENVDao().GetDefaultEnvWhichCanBeSetByPluginID(pluginID, relation.VersionID)
	if errD != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("get envs which can be set", errD)
	}
	return envD, nil
}

//BuildPluginManual BuildPluginManual
func (p *PluginAction) BuildPluginManual(bps *api_model.BuildPluginStruct) (*dbmodel.TenantPluginBuildVersion, *util.APIHandleError) {
	eventID := bps.Body.EventID
	logger := event.GetManager().GetLogger(eventID)
	defer event.CloseManager()
	plugin, err := db.GetManager().TenantPluginDao().GetPluginByID(bps.PluginID, bps.Body.TenantID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError(fmt.Sprintf("get plugin by %v", bps.PluginID), err)
	}
	switch plugin.BuildModel {
	case "image":
		pbv, err := p.buildPlugin(bps, plugin)
		if err != nil {
			logrus.Error("build plugin from image error ", err.Error())
			logger.Error("Failed to send the build plugin task from image "+err.Error(), map[string]string{"step": "callback", "status": "failure"})
			return nil, util.CreateAPIHandleError(500, fmt.Errorf("build plugin from image error"))
		}
		logger.Info("The build plugin task from the image was sent successfully ", map[string]string{"step": "image-plugin", "status": "starting"})
		return pbv, nil
	case "dockerfile":
		pbv, err := p.buildPlugin(bps, plugin)
		if err != nil {
			logrus.Error("build plugin from image error ", err.Error())
			logger.Error("Failed to send the build plugin task from dockerfile "+err.Error(), map[string]string{"step": "callback", "status": "failure"})
			return nil, util.CreateAPIHandleError(500, fmt.Errorf("build plugin from dockerfile error"))
		}
		logger.Info("The build plugin task from dockerfile was sent successfully ", map[string]string{"step": "dockerfile-plugin", "status": "starting"})
		return pbv, nil
	default:
		return nil, util.CreateAPIHandleError(400, fmt.Errorf("unexpect kind"))
	}
}

//buildPlugin buildPlugin
func (p *PluginAction) buildPlugin(b *api_model.BuildPluginStruct, plugin *dbmodel.TenantPlugin) (
	*dbmodel.TenantPluginBuildVersion, error) {
	if plugin.ImageURL == "" && plugin.BuildModel == "image" {
		return nil, fmt.Errorf("need image url")
	}
	if plugin.GitURL == "" && plugin.BuildModel == "dockerfile" {
		return nil, fmt.Errorf("need git repo url")
	}
	if b.Body.Operator == "" {
		b.Body.Operator = "define"
	}
	if b.Body.BuildVersion == "" {
		return nil, fmt.Errorf("build version can not be empty")
	}
	if b.Body.DeployVersion == "" {
		b.Body.DeployVersion = core_util.CreateVersionByTime()
	}
	pbv := &dbmodel.TenantPluginBuildVersion{
		VersionID:       b.Body.BuildVersion,
		DeployVersion:   b.Body.DeployVersion,
		PluginID:        b.PluginID,
		Kind:            plugin.BuildModel,
		Repo:            b.Body.RepoURL,
		GitURL:          plugin.GitURL,
		BaseImage:       plugin.ImageURL,
		ContainerCPU:    b.Body.PluginCPU,
		ContainerMemory: b.Body.PluginMemory,
		ContainerCMD:    b.Body.PluginCMD,
		BuildTime:       time.Now().Format(time.RFC3339),
		Info:            b.Body.Info,
		Status:          "building",
	}
	if b.Body.PluginCPU == 0 {
		pbv.ContainerCPU = 125
	}
	if b.Body.PluginMemory == 0 {
		pbv.ContainerMemory = 50
	}
	if err := db.GetManager().TenantPluginBuildVersionDao().AddModel(pbv); err != nil {
		if !strings.Contains(err.Error(), "exist") {
			logrus.Errorf("build plugin error: %s", err.Error())
			return nil, err
		}
	}
	var updateVersion = func() {
		pbv.Status = "failure"
		db.GetManager().TenantPluginBuildVersionDao().UpdateModel(pbv)
	}
	taskBody := &builder_model.BuildPluginTaskBody{
		TenantID:      b.Body.TenantID,
		PluginID:      b.PluginID,
		Operator:      b.Body.Operator,
		DeployVersion: b.Body.DeployVersion,
		ImageURL:      plugin.ImageURL,
		EventID:       b.Body.EventID,
		Kind:          plugin.BuildModel,
		PluginCMD:     b.Body.PluginCMD,
		PluginCPU:     b.Body.PluginCPU,
		PluginMemory:  b.Body.PluginMemory,
		VersionID:     b.Body.BuildVersion,
		ImageInfo:     b.Body.ImageInfo,
		Repo:          b.Body.RepoURL,
		GitURL:        plugin.GitURL,
		GitUsername:   b.Body.Username,
		GitPassword:   b.Body.Password,
	}
	taskType := "plugin_image_build"
	if plugin.BuildModel == "dockerfile" {
		taskType = "plugin_dockerfile_build"
	}
	err := p.MQClient.SendBuilderTopic(client.TaskStruct{
		TaskType: taskType,
		TaskBody: taskBody,
		Topic:    client.BuilderTopic,
	})
	if err != nil {
		updateVersion()
		logrus.Errorf("equque mq error, %v", err)
		return nil, err
	}
	logrus.Debugf("equeue mq build plugin from image success")
	return pbv, nil
}

//GetAllPluginBuildVersions GetAllPluginBuildVersions
func (p *PluginAction) GetAllPluginBuildVersions(pluginID string) ([]*dbmodel.TenantPluginBuildVersion, *util.APIHandleError) {
	versions, err := db.GetManager().TenantPluginBuildVersionDao().GetBuildVersionByPluginID(pluginID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("get all plugin build version", err)
	}
	return versions, nil
}

//GetPluginBuildVersion GetPluginBuildVersion
func (p *PluginAction) GetPluginBuildVersion(pluginID, versionID string) (*dbmodel.TenantPluginBuildVersion, *util.APIHandleError) {
	version, err := db.GetManager().TenantPluginBuildVersionDao().GetBuildVersionByVersionID(pluginID, versionID)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError(fmt.Sprintf("get plugin build version by id %v", versionID), err)
	}
	if version.Status == "building" {
		//check build whether timeout
		if buildTime, err := time.Parse(time.RFC3339, version.BuildTime); err == nil {
			if buildTime.Add(time.Minute * 5).Before(time.Now()) {
				version.Status = "timeout"
			}
		}
	}
	return version, nil
}

//DeletePluginBuildVersion DeletePluginBuildVersion
func (p *PluginAction) DeletePluginBuildVersion(pluginID, versionID string) *util.APIHandleError {

	tx := db.GetManager().Begin()
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Unexpected panic occurred, rollback transaction: %v", r)
			tx.Rollback()
		}
	}()
	err := db.GetManager().TenantPluginBuildVersionDaoTransactions(tx).DeleteBuildVersionByVersionID(versionID)
	if err != nil {
		tx.Rollback()
		return util.CreateAPIHandleErrorFromDBError(fmt.Sprintf("delete plugin build version by id %v", versionID), err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return util.CreateAPIHandleErrorFromDBError("commit delete plugin transactions", err)
	}
	return nil
}
