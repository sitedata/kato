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

package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	c "github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

// Etcd implements Discoverier
type etcd struct {
	cli       *c.Client
	endpoints []string
	key       string
	username  string
	password  string
}

// NewEtcd creates a new Discorvery which implemeted by etcd.
func NewEtcd(info *Info) Discoverier {
	// TODO: validate endpoints
	return &etcd{
		endpoints: info.Servers,
		key:       info.Key,
		username:  info.Username,
		password:  info.Password,
	}
}

// Connect connects a etcdv3 client with a given configuration.
func (e *etcd) Connect() error {
	cli, err := c.New(c.Config{
		Endpoints:   e.endpoints,
		DialTimeout: 10 * time.Second,
		Username:    e.username,
		Password:    e.password,
	})
	if err != nil {
		logrus.Errorf("Endpoints: %s; error connecting etcd: %v", strings.Join(e.endpoints, ","), err)
		return err
	}
	e.cli = cli
	return nil
}

// Fetch fetches data from Etcd.
func (e *etcd) Fetch() ([]*Endpoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if e.cli == nil {
		return nil, fmt.Errorf("can't fetching data from etcd without etcdv3 client")
	}

	resp, err := e.cli.Get(ctx, e.key, c.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("error fetching endpoints form etcd: %v", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("error fetching endpoints form etcd: empty GetResponse")
	}

	var res []*Endpoint
	for _, kv := range resp.Kvs {
		var ep Endpoint
		if err := json.Unmarshal(kv.Value, &ep); err != nil {
			return nil, fmt.Errorf("error parsing the data from etcd: %v", err)
		}
		ep.Ep = strings.Replace(string(kv.Key), e.key+"/", "", -1)
		res = append(res, &ep)
	}
	return res, nil
}

// Close shuts down the client's etcd connections.
func (e *etcd) Close() error {
	if e.cli != nil {
		return nil
	}
	return e.cli.Close()
}
