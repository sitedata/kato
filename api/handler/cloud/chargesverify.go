// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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

package cloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/model"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
)

//PubChargeSverify service Charge Sverify
func PubChargeSverify(tenant *model.Tenants, quantity int, reason string) *util.APIHandleError {
	cloudAPI := os.Getenv("CLOUD_API")
	if cloudAPI == "" {
		cloudAPI = "http://api.gridworkz.com"
	}
	regionName := os.Getenv("REGION_NAME")
	if regionName == "" {
		return util.CreateAPIHandleError(500, fmt.Errorf("region name must define in api by env REGION_NAME"))
	}
	reason = strings.Replace(reason, " ", "%20", -1)
	api := fmt.Sprintf("%s/openapi/console/v1/enterprises/%s/memory-apply?quantity=%d&tid=%s&reason=%s&region=%s", cloudAPI, tenant.EID, quantity, tenant.UUID, reason, regionName)
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		logrus.Error("create request cloud api error", err.Error())
		return util.CreateAPIHandleError(400, fmt.Errorf("create request cloud api error"))
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("create request cloud api error", err.Error())
		return util.CreateAPIHandleError(400, fmt.Errorf("create request cloud api error"))
	}
	if res.Body != nil {
		defer res.Body.Close()
		rebody, _ := ioutil.ReadAll(res.Body)
		logrus.Debugf("read memory-apply response (%s)", string(rebody))
		var re = make(map[string]interface{})
		if err := ffjson.Unmarshal(rebody, &re); err == nil {
			if msg, ok := re["msg"]; ok {
				return util.CreateAPIHandleError(res.StatusCode, fmt.Errorf("%s", msg))
			}
		}
	}
	return util.CreateAPIHandleError(res.StatusCode, fmt.Errorf("none"))
}

// PriChargeSverify verifies that the resources requested in the private cloud are legal
func PriChargeSverify(tenant *model.Tenants, quantity int) *util.APIHandleError {
	t, err := db.GetManager().TenantDao().GetTenantByUUID(tenant.UUID)
	if err != nil {
		logrus.Errorf("error getting tenant: %v", err)
		return util.CreateAPIHandleError(500, fmt.Errorf("error getting tenant: %v", err))
	}
	if t.LimitMemory == 0 {
		clusterStats, err := handler.GetTenantManager().GetAllocatableResources()
		if err != nil {
			logrus.Errorf("error getting allocatable resources: %v", err)
			return util.CreateAPIHandleError(500, fmt.Errorf("error getting allocatable resources: %v", err))
		}
		availMem := clusterStats.AllMemory - clusterStats.RequestMemory
		if availMem >= int64(quantity) {
			return util.CreateAPIHandleError(200, fmt.Errorf("success"))
		}
		return util.CreateAPIHandleError(200, fmt.Errorf("cluster_lack_of_memory"))
	}
	tenantStas, err := handler.GetTenantManager().GetTenantResource(tenant.UUID)
	// TODO: it should be limit, not request
	availMem := int64(t.LimitMemory) - (tenantStas.MemoryRequest + tenantStas.UnscdMemoryReq)
	if availMem >= int64(quantity) {
		return util.CreateAPIHandleError(200, fmt.Errorf("success"))
	}
	return util.CreateAPIHandleError(200, fmt.Errorf("lack_of_memory"))
}
