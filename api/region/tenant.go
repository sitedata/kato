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

package region

import (
	"path"

	"github.com/gridworkz/kato/api/util"
	dbmodel "github.com/gridworkz/kato/db/model"
	utilhttp "github.com/gridworkz/kato/util/http"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type tenant struct {
	regionImpl
	tenantName string
	prefix     string
}

//TenantInterface
type TenantInterface interface {
	Get() (*dbmodel.Tenants, *util.APIHandleError)
	List() ([]*dbmodel.Tenants, *util.APIHandleError)
	Delete() *util.APIHandleError
	Services(serviceAlias string) ServiceInterface
	// DefineSources(ss *api_model.SourceSpec) DefineSourcesInterface
	// DefineCloudAuth(gt *api_model.GetUserToken) DefineCloudAuthInterface
}

func (t *tenant) Get() (*dbmodel.Tenants, *util.APIHandleError) {
	var decode utilhttp.ResponseBody
	var tenant dbmodel.Tenants
	decode.Bean = &tenant
	code, err := t.DoRequest(t.prefix, "GET", nil, &decode)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	return &tenant, nil
}
func (t *tenant) List() ([]*dbmodel.Tenants, *util.APIHandleError) {
	if t.tenantName != "" {
		return nil, util.CreateAPIHandleErrorf(400, "tenant name must be empty in this api")
	}
	var decode utilhttp.ResponseBody
	code, err := t.DoRequest(t.prefix, "GET", nil, &decode)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if decode.Bean == nil {
		return nil, nil
	}
	bean, ok := decode.Bean.(map[string]interface{})
	if !ok {
		logrus.Warningf("list tenants; wrong data: %v", decode.Bean)
		return nil, nil
	}
	list, ok := bean["list"]
	if !ok {
		return nil, nil
	}
	var tenants []*dbmodel.Tenants
	if err := mapstructure.Decode(list, &tenants); err != nil {
		logrus.Errorf("map: %+v; error decoding to map to []*dbmodel.Tenants: %v", list, err)
		return nil, util.CreateAPIHandleError(500, err)
	}
	return tenants, nil
}
func (t *tenant) Delete() *util.APIHandleError {
	return nil
}
func (t *tenant) Services(serviceAlias string) ServiceInterface {
	return &services{
		prefix: path.Join(t.prefix, "services", serviceAlias),
		tenant: *t,
	}
}
