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

package version2

import (
	"github.com/gridworkz/kato/api/controller"
	"github.com/gridworkz/kato/api/middleware"
	dbmodel "github.com/gridworkz/kato/db/model"

	"github.com/go-chi/chi"
)

//PluginRouter
func (v2 *V2) pluginRouter() chi.Router {
	r := chi.NewRouter()
	//Initialize application information
	r.Use(middleware.InitPlugin)
	//plugin uri
	//update/delete plugin
	r.Put("/", controller.GetManager().PluginAction)
	r.Delete("/", controller.GetManager().PluginAction)
	r.Post("/build", controller.GetManager().PluginBuild)
	//get this plugin all build version
	r.Get("/build-version", controller.GetManager().GetAllPluginBuildVersions)
	r.Get("/build-version/{version_id}", controller.GetManager().GetPluginBuildVersion)
	r.Delete("/build-version/{version_id}", controller.GetManager().DeletePluginBuildVersion)
	return r
}

func (v2 *V2) serviceRelatePluginRouter() chi.Router {
	r := chi.NewRouter()
	//service related plugin
	// v2/tenant/tenant_name/services/service_alias/plugin/xxx
	r.Post("/", middleware.WrapEL(controller.GetManager().PluginSet, dbmodel.TargetTypeService, "create-service-plugin", dbmodel.SYNEVENTTYPE))
	r.Put("/", middleware.WrapEL(controller.GetManager().PluginSet, dbmodel.TargetTypeService, "update-service-plugin", dbmodel.SYNEVENTTYPE))
	r.Get("/", controller.GetManager().PluginSet)
	r.Delete("/{plugin_id}", middleware.WrapEL(controller.GetManager().DeletePluginRelation, dbmodel.TargetTypeService, "delete-service-plugin", dbmodel.SYNEVENTTYPE))
	// app plugin config supdate
	r.Post("/{plugin_id}/setenv", middleware.WrapEL(controller.GetManager().UpdateVersionEnv, dbmodel.TargetTypeService, "update-service-plugin-config", dbmodel.SYNEVENTTYPE))
	r.Put("/{plugin_id}/upenv", middleware.WrapEL(controller.GetManager().UpdateVersionEnv, dbmodel.TargetTypeService, "update-service-plugin-config", dbmodel.SYNEVENTTYPE))
	//deprecated
	r.Get("/{plugin_id}/envs", controller.GetManager().GePluginEnvWhichCanBeSet)
	return r
}
