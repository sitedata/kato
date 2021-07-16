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

package discover

import (
	"errors"
	"sync"

	"github.com/gridworkz/kato/api/proxy"
	corediscover "github.com/gridworkz/kato/discover"
	corediscoverconfig "github.com/gridworkz/kato/discover/config"
	etcdutil "github.com/gridworkz/kato/util/etcd"

	"github.com/sirupsen/logrus"
)

//EndpointDiscover Automatic discovery of back-end services
type EndpointDiscover interface {
	AddProject(name string, proxy proxy.Proxy)
	Remove(name string)
	Stop()
}

var defaultEndpointDiscover EndpointDiscover

//CreateEndpointDiscover create endpoint discover
func CreateEndpointDiscover(etcdClientArgs *etcdutil.ClientArgs) (EndpointDiscover, error) {
	if defaultEndpointDiscover == nil {
		if etcdClientArgs == nil {
			return nil, errors.New("etcd args is nil")
		}
		dis, err := corediscover.GetDiscover(corediscoverconfig.DiscoverConfig{EtcdClientArgs: etcdClientArgs})
		if err != nil {
			return nil, err
		}
		defaultEndpointDiscover = &endpointDiscover{
			dis:      dis,
			projects: make(map[string]*defalt),
		}
	}
	return defaultEndpointDiscover, nil
}

//GetEndpointDiscover get endpoints discover
func GetEndpointDiscover() EndpointDiscover {
	return defaultEndpointDiscover
}

type endpointDiscover struct {
	projects map[string]*defalt
	lock     sync.Mutex
	dis      corediscover.Discover
}

func (e *endpointDiscover) AddProject(name string, pro proxy.Proxy) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if def, ok := e.projects[name]; !ok {
		e.projects[name] = &defalt{name: name, proxys: []proxy.Proxy{pro}}
		e.dis.AddProject(name, e.projects[name])
	} else {
		def.proxys = append(def.proxys, pro)
		// add proxy after update endpoint first,must initialize endpoint by cache data
		if len(def.cacheEndpointURL) > 0 {
			pro.UpdateEndpoints(def.cacheEndpointURL...)
		}
	}

}
func (e *endpointDiscover) Remove(name string) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if _, ok := e.projects[name]; ok {
		delete(e.projects, name)
	}
}
func (e *endpointDiscover) Stop() {
	e.dis.Stop()
}

type defalt struct {
	name             string
	proxys           []proxy.Proxy
	cacheEndpointURL []string
}

func (e *defalt) Error(err error) {
	logrus.Errorf("%s project auto discover occurred error.%s", e.name, err.Error())
	defaultEndpointDiscover.Remove(e.name)
}

func (e *defalt) UpdateEndpoints(endpoints ...*corediscoverconfig.Endpoint) {
	var endStr []string
	for _, end := range endpoints {
		if end.URL != "" {
			endStr = append(endStr, end.Name+"=>"+end.URL)
		}
	}
	logrus.Debugf("endstr is %v, name is %v", endStr, e.name)
	for _, p := range e.proxys {
		p.UpdateEndpoints(endStr...)
	}
	e.cacheEndpointURL = endStr
}

//when watch occurred error,will exec this method
func (e *endpointDiscover) Error(err error) {
	logrus.Errorf("discover error %s", err.Error())
}
