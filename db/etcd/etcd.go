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

package etcd

import (
	"context"
	"sync"

	"github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/db/model"

	etcdutil "github.com/gridworkz/kato/util/etcd"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

//Manager - db manager
type Manager struct {
	client  *clientv3.Client
	config  config.Config
	initOne sync.Once
	models  []model.Interface
}

//CreateManager
func CreateManager(config config.Config) (*Manager, error) {
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: config.EtcdEndPoints,
		CaFile:    config.EtcdCaFile,
		CertFile:  config.EtcdCertFile,
		KeyFile:   config.EtcdKeyFile,
	}
	cli, err := etcdutil.NewClient(context.Background(), etcdClientArgs)
	if err != nil {
		etcdutil.HandleEtcdError(err)
		return nil, err
	}
	manager := &Manager{
		client:  cli,
		config:  config,
		initOne: sync.Once{},
	}
	logrus.Debug("etcd db driver create")
	return manager, nil
}

//CloseManager
func (m *Manager) CloseManager() error {
	return m.client.Close()
}
