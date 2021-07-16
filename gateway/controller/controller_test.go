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

package controller

import (
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	v1 "github.com/gridworkz/kato/gateway/v1"
)

func TestGWController_WatchRbdEndpoints(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	gwc := GWController{
		EtcdCli: cli,
	}
	go gwc.watchRbdEndpoints(0)
}

func poolsIsEqual(old []*v1.Pool, new []*v1.Pool) bool {
	if len(old) != len(new) {
		return false
	}
	for _, rp := range old {
		flag := false
		for _, cp := range new {
			if rp.Equals(cp) {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func newFakePoolWithoutNodes(name string) *v1.Pool {
	return &v1.Pool{
		Meta: v1.Meta{
			Index:      888,
			Name:       name,
			Namespace:  "gateway",
			PluginName: "Nginx",
		},
		ServiceID:         "foo-service-id",
		ServiceVersion:    "1.0.0",
		ServicePort:       80,
		Note:              "foo",
		NodeNumber:        8,
		LoadBalancingType: v1.RoundRobin,
		Monitors: []v1.Monitor{
			"monitor-a",
			"monitor-b",
			"monitor-c",
		},
	}
}
