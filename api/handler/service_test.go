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
	"testing"

	"github.com/gridworkz/kato/api/client/prometheus"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
)

func TestABCService(t *testing.T) {
	mm := `{
		"comment":"",
		"container_env":"",
		"domain":"lichao",
		"deploy_version":"",
		"ports_info":[
			{
				"port_alias":"GR45068C5000",
				"protocol":"http",
				"mapping_port":0,
				"container_port":5000,
				"is_outer_service":true,
				"is_inner_service":false
			}
		],
		"dep_sids":null,
		"volumes_info":[
	
		],
		"extend_method":"stateless",
		"operator":"lichao",
		"container_memory":512,
		"service_key":"application",
		"category":"application",
		"service_version":"81701",
		"event_id":"e5bd1926254b447ea97817566b2d71bf",
		"container_cpu":80,
		"namespace":"gridworkz",
		"extend_info":{
			"envs":[
	
			],
			"ports":[
	
			]
		},
		"service_type":"application",
		"status":0,
		"node_label":"",
		"replicas":1,
		"image_name":"gridworkz/runner",
		"service_alias":"gr45068c",
		"service_id":"55c60b74a506261608f5c36f0f45068c",
		"code_from":"gitlab_manual",
		"volume_mount_path":"/data",
		"tenant_id":"3000bf47672b40c19529504651697b29",
		"container_cmd":"start web",
		"host_path":"/grdata/tenant/3000bf47672b40c19529504651697b29/service/55c60b74a506261608f5c36f0f45068c",
		"envs_info":[
	
		],
		"volume_path":"vol55c60b74a5",
		"port_type":"multi_outer"
	}`

	var s api_model.ServiceStruct
	err := ffjson.Unmarshal([]byte(mm), &s)
	if err != nil {
		fmt.Printf("err is %v", err)
	}
	fmt.Printf("json is \n %v", s)
}

func TestUUID(t *testing.T) {
	id := fmt.Sprintf("%s", uuid.NewV4())
	uid := strings.Replace(id, "-", "", -1)
	logrus.Debugf("uuid is %v", uid)
	name := strings.Split(id, "-")[0]
	fmt.Printf("id is %s, uid is %s, name is %v", id, uid, name)
}

func TestGetServicesDisk(t *testing.T) {
	prometheusCli, err := prometheus.NewPrometheus(&prometheus.Options{
		Endpoint: "39.96.189.166:9999",
	})
	if err != nil {
		t.Fatal(err)
	}
	disk := GetServicesDiskDeprecated([]string{"ef75e1d5e3df412a8af06129dae42869"}, prometheusCli)
	t.Log(disk)
}
