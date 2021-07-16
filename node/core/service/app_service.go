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

package service

import (
	"fmt"
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/discover/config"
	"github.com/gridworkz/kato/node/core/store"
	"strconv"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

//AppService
type AppService struct {
	Prefix string
	c      *option.Conf
}

//CreateAppService
func CreateAppService(c *option.Conf) *AppService {
	return &AppService{
		c:      c,
		Prefix: "/traefik",
	}
}

//FindAppEndpoints
func (a *AppService) FindAppEndpoints(appName string) []*config.Endpoint {
	var ends = make(map[string]*config.Endpoint)
	res, err := store.DefalutClient.Get(fmt.Sprintf("%s/backends/%s/servers", a.Prefix, appName), clientv3.WithPrefix())
	if err != nil {
		logrus.Errorf("list all servers of %s error.%s", appName, err.Error())
		return nil
	}
	if res.Count == 0 {
		return nil
	}
	for _, kv := range res.Kvs {
		if strings.HasSuffix(string(kv.Key), "/url") { //Get service address
			kstep := strings.Split(string(kv.Key), "/")
			if len(kstep) > 2 {
				serverName := kstep[len(kstep)-2]
				serverURL := string(kv.Value)
				if en, ok := ends[serverName]; ok {
					en.URL = serverURL
				} else {
					ends[serverName] = &config.Endpoint{Name: serverName, URL: serverURL}
				}
			}
		}
		if strings.HasSuffix(string(kv.Key), "/weight") { //Get service weight
			kstep := strings.Split(string(kv.Key), "/")
			if len(kstep) > 2 {
				serverName := kstep[len(kstep)-2]
				serverWeight := string(kv.Value)
				if en, ok := ends[serverName]; ok {
					var err error
					en.Weight, err = strconv.Atoi(serverWeight)
					if err != nil {
						logrus.Error("get server weight error.", err.Error())
					}
				} else {
					weight, err := strconv.Atoi(serverWeight)
					if err != nil {
						logrus.Error("get server weight error.", err.Error())
					}
					ends[serverName] = &config.Endpoint{Name: serverName, Weight: weight}
				}
			}
		}
	}
	result := []*config.Endpoint{}
	for _, v := range ends {
		result = append(result, v)
	}
	return result
}
