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

package parser

import (
	"fmt"
	"testing"

	//"github.com/docker/docker/client"
	"github.com/docker/docker/client"
)

var dockerrun = `docker run -d -P -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
--name etcd quay.io/coreos/etcd:v2.3.8 \
-name etcd0 \
-advertise-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
-listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
-initial-advertise-peer-urls http://127.0.0.1:2380 \
-listen-peer-urls http://0.0.0.0:2380 \
-initial-cluster-token etcd-cluster-1 \
-initial-cluster etcd0=http://127.0.0.1:2380 \
-initial-cluster-state new`

var test_case = `docker run -d --restart=always --name powerjob-server -p 7700:7700 -p 10086:10086 -e TZ="America/Toronto" -e JVMOPTIONS="" -e PARAMS="--spring.profiles.active=product --spring.datasource.core.jdbc-url=jdbc:postgresql://127.0.0.1:5432/powerjob-product?useUnicode=true&characterEncoding=UTF-8&serverTimezone=America/Toronto --spring.datasource.core.username=admin --spring.datasource.core.password=d9e6c012 --oms.mongodb.enable=false --oms.mongodb.enable=false --spring.data.mongodb.uri=mongodb://127.0.0.1:27017/powerjob-product" -v ~/docker/powerjob-server:/root/powerjob-server -v ~/.m2:/root/.m2 tjqq/powerjob-server:latest`

func TestParse(t *testing.T) {
	dockerclient, err := client.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}
	p := CreateDockerRunOrImageParse("d", "", test_case, dockerclient, nil)
	p.ParseDockerun(test_case)
	fmt.Printf("ServiceInfo:%+v \n", p.GetServiceInfo())
}
