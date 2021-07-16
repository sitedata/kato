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

package core

import (
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/core/store"
	"sync"

	"github.com/sirupsen/logrus"
)

type etcdRegistrar struct {
	//projects map[string]CallbackUpdate
	lock   sync.Mutex
	client *store.Client
	prefix string
}

//GetRegistrar
func GetRegistrar() *etcdRegistrar {
	return &etcdRegistrar{
		prefix: option.Config.Service,
		client: store.DefalutClient,
	}
}
func (r *etcdRegistrar) RegService(serviceName, hostname, url string) error {
	r.lock.Lock()
	_, err := r.client.Put(r.getPath(serviceName, hostname), url)
	if err != nil {
		logrus.Infof("reg service %s to path %s failed,details %s", serviceName, r.getPath(serviceName, hostname), err.Error())
		return err
	}
	r.lock.Unlock()
	return nil
}
func (r *etcdRegistrar) RemoveService(serviceName, hostname string) error {
	r.lock.Lock()
	_, err := r.client.Delete(r.getPath(serviceName, hostname))
	if err != nil {
		logrus.Infof("del  service %s from path %s failed,details %s", serviceName, r.getPath(serviceName, hostname), err.Error())
		return err
	}
	r.lock.Unlock()
	return nil
}
func (r *etcdRegistrar) getPath(serviceName, hostName string) string {
	return r.prefix + serviceName + "/servers/" + hostName + "/url"
}
