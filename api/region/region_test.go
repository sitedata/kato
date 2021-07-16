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

package region

import (
	"testing"

	dbmodel "github.com/gridworkz/kato/db/model"
	utilhttp "github.com/gridworkz/kato/util/http"
)

func TestListTenant(t *testing.T) {
	region, _ := NewRegion(APIConf{
		Endpoints: []string{"http://kubeapi.gridworkz:8888"},
	})
	tenants, err := region.Tenants("").List()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", tenants)
}

func TestListServices(t *testing.T) {
	region, _ := NewRegion(APIConf{
		Endpoints: []string{"http://kubeapi.gridworkz:8888"},
	})
	services, err := region.Tenants("n93lkp7t").Services("").List()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		t.Logf("%+v", s)
	}
}

func TestDoRequest(t *testing.T) {
	region, _ := NewRegion(APIConf{
		Endpoints: []string{"http://kubeapi.gridworkz:8888"},
	})
	var decode utilhttp.ResponseBody
	var tenants []*dbmodel.Tenants
	decode.List = &tenants
	code, err := region.DoRequest("/v2/tenants", "GET", nil, &decode)
	if err != nil {
		t.Fatal(err, code)
	}
	t.Logf("%+v", tenants)
}

func TestListNodes(t *testing.T) {
	region, _ := NewRegion(APIConf{
		Endpoints: []string{"http://kubeapi.gridworkz:8888"},
	})
	services, err := region.Nodes().List()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		t.Logf("%+v", s)
	}
}

func TestGetNodes(t *testing.T) {
	region, _ := NewRegion(APIConf{
		Endpoints: []string{"http://kubeapi.gridworkz:8888"},
	})
	node, err := region.Nodes().Get("a134eab8-3d42-40f5-84a5-fcf2b7a44b31")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", node)
}

func TestGetTenantsBySSL(t *testing.T) {
	region, _ := NewRegion(APIConf{
		Endpoints: []string{"https://127.0.0.1:8443"},
		Cacert:    "/Users/devs/gopath/src/github.com/gridworkz/kato/test/ssl/ca.pem",
		Cert:      "/Users/devs/gopath/src/github.com/gridworkz/kato/test/ssl/client.pem",
		CertKey:   "/Users/devs/gopath/src/github.com/gridworkz/kato/test/ssl/client.key.pem",
	})
	tenants, err := region.Tenants("").List()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", tenants)
}
