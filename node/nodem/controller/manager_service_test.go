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
	"context"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/nodem/client"
	"github.com/gridworkz/kato/node/nodem/service"
)

func TestManagerService_SetEndpoints(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	key := "/kato/endpoint/foobar"
	defer cli.Delete(ctx, key, clientv3.WithPrefix())

	m := &ManagerService{}
	srvs := []*service.Service{
		{
			Endpoints: []*service.Endpoint{
				{
					Name:     "foobar",
					Protocol: "http",
					Port:     "6442",
				},
			},
		},
	}
	m.services = srvs
	c := client.NewClusterClient(
		&option.Conf{
			EtcdCli: cli,
		},
	)
	m.cluster = c

	data := []string{
		"192.168.8.229",
		"192.168.8.230",
		"192.168.8.231",
	}

	m.SetEndpoints(data[0])
	m.SetEndpoints(data[1])
	m.SetEndpoints(data[2])

	edps := c.GetEndpoints("foobar")
	for _, d := range data {
		flag := false
		for _, edp := range edps {
			if d+":6442" == strings.Replace(edp, "http://", "", -1) {
				flag = true
			}
		}
		if !flag {
			t.Fatalf("Can not find \"%s\" in %v", d, edps)
		}
	}
}
