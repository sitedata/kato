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

	"github.com/pquerna/ffjson/ffjson"

	"github.com/coreos/etcd/clientv3"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/sirupsen/logrus"
)

//NetRulesAction  rules action struct
type NetRulesAction struct {
	etcdCli *clientv3.Client
}

//CreateNetRulesManager get net rules manager
func CreateNetRulesManager(etcdCli *clientv3.Client) *NetRulesAction {
	return &NetRulesAction{
		etcdCli: etcdCli,
	}
}

//CreateDownStreamNetRules CreateDownStreamNetRules
func (n *NetRulesAction) CreateDownStreamNetRules(
	tenantID string,
	rs *api_model.SetNetDownStreamRuleStruct) *util.APIHandleError {
	k := fmt.Sprintf("/netRules/%s/%s/downstream/%s/%v",
		tenantID, rs.ServiceAlias, rs.Body.DestServiceAlias, rs.Body.Port)
	sb := &api_model.NetRulesDownStreamBody{
		DestService:      rs.Body.DestService,
		DestServiceAlias: rs.Body.DestServiceAlias,
		Port:             rs.Body.Port,
		Protocol:         rs.Body.Protocol,
		Rules:            rs.Body.Rules,
	}
	v, err := ffjson.Marshal(sb)
	if err != nil {
		logrus.Errorf("mashal etcd value error, %v", err)
		return util.CreateAPIHandleError(500, err)
	}
	_, err = n.etcdCli.Put(context.TODO(), k, string(v))
	if err != nil {
		logrus.Errorf("put k %s into etcd error, %v", k, err)
		return util.CreateAPIHandleError(500, err)
	}
	//TODO: store mysql
	return nil
}

//GetDownStreamNetRule GetDownStreamNetRule
func (n *NetRulesAction) GetDownStreamNetRule(
	tenantID,
	serviceAlias,
	destServiceAlias,
	port string) (*api_model.NetRulesDownStreamBody, *util.APIHandleError) {
	k := fmt.Sprintf(
		"/netRules/%s/%s/downstream/%s/%v",
		tenantID,
		serviceAlias,
		destServiceAlias,
		port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := n.etcdCli.Get(ctx, k)
	cancel()
	if err != nil {
		logrus.Errorf("get etcd value error, %v", err)
		return nil, util.CreateAPIHandleError(500, err)
	}
	if resp.Count != 0 {
		v := resp.Kvs[0].Value
		var sb api_model.NetRulesDownStreamBody
		if err := ffjson.Unmarshal(v, &sb); err != nil {
			logrus.Errorf("unmashal etcd v error, %v", err)
			return nil, util.CreateAPIHandleError(500, err)
		}
		return &sb, nil
	}
	//TODO: query mysql
	//TODO: create etcd record
	return nil, nil
}

//UpdateDownStreamNetRule UpdateDownStreamNetRule
func (n *NetRulesAction) UpdateDownStreamNetRule(
	tenantID string,
	urs *api_model.UpdateNetDownStreamRuleStruct) *util.APIHandleError {

	srs := &api_model.SetNetDownStreamRuleStruct{
		TenantName:   urs.TenantName,
		ServiceAlias: urs.ServiceAlias,
	}
	srs.Body.DestService = urs.Body.DestService
	srs.Body.DestServiceAlias = urs.DestServiceAlias
	srs.Body.Port = urs.Port
	srs.Body.Protocol = urs.Body.Protocol
	srs.Body.Rules = urs.Body.Rules

	//TODO: update mysql transaction
	if err := n.CreateDownStreamNetRules(tenantID, srs); err != nil {
		return err
	}
	return nil
}
