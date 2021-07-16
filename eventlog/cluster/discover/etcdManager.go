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
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/sirupsen/logrus"
	"time"

	"golang.org/x/net/context"
)

//SaveDockerLogInInstance - store the correspondence between service and node
func SaveDockerLogInInstance(etcdClient *clientv3.Client, conf conf.DiscoverConf, serviceID, instanceID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := etcdClient.Put(ctx, conf.HomePath+"/dockerloginstacne/"+serviceID, instanceID)
	if err != nil {
		logrus.Errorf("Failed to put dockerlog instance %v", err)
		return err
	}
	return nil
}

//GetDokerLogInInstance - get application log receiving node
func GetDokerLogInInstance(etcdClient *clientv3.Client, conf conf.DiscoverConf, serviceID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := etcdClient.Get(ctx, conf.HomePath+"/dockerloginstacne/"+serviceID)
	if err != nil {
		return "", err
	}
	if len(res.Kvs) == 0 {
		return "", fmt.Errorf("get docker log instance failed")
	}
	return string(res.Kvs[0].Value), nil
}
