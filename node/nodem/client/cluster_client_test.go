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

package client

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/core/store"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
)

func TestEtcdClusterClient_GetEndpoints(t *testing.T) {
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

	type data struct {
		hostID string
		values []string
	}

	testCase := []string{
		"192.168.8.229:8081",
		"192.168.8.230:8081",
		"192.168.8.231:6443",
	}
	datas := []data{
		{
			hostID: util.NewUUID(),
			values: []string{
				testCase[0],
			},
		},
		{
			hostID: util.NewUUID(),
			values: []string{
				testCase[1],
			},
		},
		{
			hostID: util.NewUUID(),
			values: []string{
				testCase[2],
			},
		},
	}
	for _, d := range datas {
		s, err := json.Marshal(d.values)
		if err != nil {
			logrus.Errorf("Can not marshal %s endpoints to json.", "foobar")
			return
		}
		_, err = cli.Put(ctx, key+"/"+d.hostID, string(s))
		if err != nil {
			t.Fatal(err)
		}
	}

	c := etcdClusterClient{
		conf: &option.Conf{
			EtcdCli: cli,
		},
	}

	edps := c.GetEndpoints("foobar")
	for _, tc := range testCase {
		flag := false
		for _, edp := range edps {
			if tc == edp {
				flag = true
			}
		}
		if !flag {
			t.Fatalf("Can not find \"%s\" in %v", tc, edps)
		}
	}
}

func TestSetEndpoints(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	c := NewClusterClient(&option.Conf{EtcdCli: cli})
	c.SetEndpoints("etcd", "DSASD", []string{"http://:8080"})
	c.SetEndpoints("etcd", "192.168.1.1", []string{"http://:8080"})
	c.SetEndpoints("etcd", "192.168.1.1", []string{"http://192.168.1.1:8080"})
	c.SetEndpoints("node", "192.168.2.137", []string{"192.168.2.137:10252"})
	t.Logf("check: %v", checkURL("192.168.2.137:10252"))
}

func TestGetEndpoints(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	c := NewClusterClient(&option.Conf{EtcdCli: cli})
	t.Log(c.GetEndpoints("/etcd/"))
}
func TestEtcdClusterClient_ListEndpointKeys(t *testing.T) {
	cfg := &option.Conf{
		EtcdEndpoints:   []string{"192.168.3.3:2379"},
		EtcdDialTimeout: 5 * time.Second,
	}

	if err := store.NewClient(context.Background(), cfg); err != nil {
		t.Fatalf("error create etcd client: %v", err)
	}

	hostNode := HostNode{
		InternalIP: "192.168.2.76",
	}

	keys, err := hostNode.listEndpointKeys()
	if err != nil {
		t.Errorf("unexperted error: %v", err)
	}
	t.Logf("keys: %#v", keys)
}
