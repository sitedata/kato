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
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	dbmodel "github.com/gridworkz/kato/db/model"
)

//PluginHandler plugin handler
type PluginHandler interface {
	CreatePluginAct(cps *api_model.CreatePluginStruct) *util.APIHandleError
	UpdatePluginAct(pluginID, tenantID string, cps *api_model.UpdatePluginStruct) *util.APIHandleError
	DeletePluginAct(pluginID, tenantID string) *util.APIHandleError
	GetPlugins(tenantID string) ([]*dbmodel.TenantPlugin, *util.APIHandleError)
	AddDefaultEnv(est *api_model.ENVStruct) *util.APIHandleError
	UpdateDefaultEnv(est *api_model.ENVStruct) *util.APIHandleError
	DeleteDefaultEnv(pluginID, versionID, envName string) *util.APIHandleError
	BuildPluginManual(bps *api_model.BuildPluginStruct) (*dbmodel.TenantPluginBuildVersion, *util.APIHandleError)
	GetAllPluginBuildVersions(pluginID string) ([]*dbmodel.TenantPluginBuildVersion, *util.APIHandleError)
	GetPluginBuildVersion(pluginID, versionID string) (*dbmodel.TenantPluginBuildVersion, *util.APIHandleError)
	DeletePluginBuildVersion(pluginID, versionID string) *util.APIHandleError
	GetDefaultEnv(pluginID, versionID string) ([]*dbmodel.TenantPluginDefaultENV, *util.APIHandleError)
	GetEnvsWhichCanBeSet(serviceID, pluginID string) (interface{}, *util.APIHandleError)
}
