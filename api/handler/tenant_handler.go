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
	"context"

	"github.com/gridworkz/kato/api/model"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	dbmodel "github.com/gridworkz/kato/db/model"
)

//TenantHandler tenant handler
type TenantHandler interface {
	GetTenants(query string) ([]*dbmodel.Tenants, error)
	GetTenantsByName(name string) (*dbmodel.Tenants, error)
	GetTenantsByEid(eid, query string) ([]*dbmodel.Tenants, error)
	GetTenantsByUUID(uuid string) (*dbmodel.Tenants, error)
	GetTenantsName() ([]string, error)
	StatsMemCPU(services []*dbmodel.TenantServices) (*api_model.StatsInfo, error)
	TotalMemCPU(services []*dbmodel.TenantServices) (*api_model.StatsInfo, error)
	GetTenantsResources(ctx context.Context, tr *api_model.TenantResources) (map[string]map[string]interface{}, error)
	GetTenantResource(tenantID string) (TenantResourceStats, error)
	GetAllocatableResources(ctx context.Context) (*ClusterResourceStats, error)
	GetServicesResources(tr *api_model.ServicesResources) (map[string]map[string]interface{}, error)
	TenantsSum() (int, error)
	GetProtocols() ([]*dbmodel.RegionProcotols, *util.APIHandleError)
	TransPlugins(tenantID, tenantName, fromTenant string, pluginList []string) *util.APIHandleError
	GetServicesStatus(ids string) map[string]string
	IsClosedStatus(status string) bool
	BindTenantsResource(source []*dbmodel.Tenants) api_model.TenantList
	UpdateTenant(*dbmodel.Tenants) error
	DeleteTenant(ctx context.Context, tenantID string) error
	GetClusterResource(ctx context.Context) *ClusterResourceStats
	CheckResourceName(ctx context.Context, namespace string, req *model.CheckResourceNameReq) (*model.CheckResourceNameResp, error)
}
