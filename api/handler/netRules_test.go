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
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	api_model "github.com/gridworkz/kato/api/model"

	"testing"
)

func TestStoreETCD(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		t.Error(err)
	}
	nra := &NetRulesAction{
		etcdCli: cli,
	}
	rules := &api_model.NetDownStreamRules{
		Limit: 1024,
		//Header: "E1:V1,E2:V2",
		//Domain: "test.redis.com",
		//Prefix: "/redis",
	}

	srs := &api_model.SetNetDownStreamRuleStruct{
		TenantName:   "123",
		ServiceAlias: "grtest12",
	}
	srs.Body.DestService = "redis"
	srs.Body.DestServiceAlias = "grtest34"
	srs.Body.Port = 6379
	srs.Body.Protocol = "tcp"
	srs.Body.Rules = rules

	tenantID := "tenantid1b50sfadfadfafadfadfadf"

	if err := nra.CreateDownStreamNetRules(tenantID, srs); err != nil {
		t.Error(err)
	}

	k := fmt.Sprintf("/netRules/%s/%s/downstream/%s/%v",
		tenantID, srs.ServiceAlias, srs.Body.DestServiceAlias, srs.Body.Port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	v, err := cli.Get(ctx, k)
	cancel()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("v is %v\n", string(v.Kvs[0].Value))
}
