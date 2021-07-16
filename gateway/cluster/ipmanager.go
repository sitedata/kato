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

package cluster

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/coreos/etcd/clientv3"

	"github.com/gridworkz/kato/cmd/gateway/option"
	"github.com/gridworkz/kato/util"
)

//IPManager ip manager
//Gets all available IP addresses for synchronizing the current node
type IPManager interface {
	//Whether the IP address belongs to the current node
	IPInCurrentHost(net.IP) bool
	Start() error
	//An IP pool change triggers a forced update of the gateway policy
	NeedUpdateGatewayPolicy() <-chan util.IPEVENT
	Stop()
}

type ipManager struct {
	ctx     context.Context
	cancel  context.CancelFunc
	IPPool  *util.IPPool
	ipLease map[string]clientv3.LeaseID
	lock    sync.Mutex
	etcdCli *clientv3.Client
	config  option.Config
	//An IP pool change triggers a forced update of the gateway policy
	needUpdate chan util.IPEVENT
}

//CreateIPManager create ip manage
func CreateIPManager(ctx context.Context, config option.Config, etcdcli *clientv3.Client) (IPManager, error) {
	newCtx, cancel := context.WithCancel(ctx)
	IPPool := util.NewIPPool(config.IgnoreInterface)
	return &ipManager{
		ctx:        newCtx,
		cancel:     cancel,
		IPPool:     IPPool,
		config:     config,
		etcdCli:    etcdcli,
		ipLease:    make(map[string]clientv3.LeaseID),
		needUpdate: make(chan util.IPEVENT, 10),
	}, nil
}

func (i *ipManager) NeedUpdateGatewayPolicy() <-chan util.IPEVENT {
	return i.needUpdate
}

//IPInCurrentHost Whether the IP address belongs to the current node
func (i *ipManager) IPInCurrentHost(in net.IP) bool {
	for _, exit := range i.IPPool.GetHostIPs() {
		if exit.Equal(in) {
			return true
		}
	}
	return false
}

func (i *ipManager) Start() error {
	logrus.Info("start ip manager.")
	go i.IPPool.LoopCheckIPs()
	i.IPPool.Ready()
	logrus.Info("ip manager is ready.")
	go i.syncIP()
	return nil
}

func (i *ipManager) syncIP() {
	logrus.Debugf("start syncronizing ip.")
	ips := i.IPPool.GetHostIPs()
	i.updateIP(ips...)
	for ipevent := range i.IPPool.GetWatchIPChan() {
		switch ipevent.Type {
		case util.ADD:
			i.updateIP(ipevent.IP)
		case util.UPDATE:
			i.updateIP(ipevent.IP)
		case util.DEL:
			i.deleteIP(ipevent.IP)
		}
		i.needUpdate <- ipevent
	}
}

func (i *ipManager) updateIP(ips ...net.IP) error {
	ctx, cancel := context.WithTimeout(i.ctx, time.Second*30)
	defer cancel()
	i.lock.Lock()
	defer i.lock.Unlock()
	leaseClient := clientv3.NewLease(i.etcdCli)
	for in := range ips {
		ip := ips[in]
		if id, ok := i.ipLease[ip.String()]; ok {
			if _, err := leaseClient.KeepAliveOnce(ctx, id); err == nil {
				continue
			} else {
				logrus.Warningf("keep alive ip key failure %s", err.Error())
			}
		}
		res, err := leaseClient.Grant(ctx, 10)
		if err != nil {
			logrus.Errorf("put gateway ip to etcd failure %s", err.Error())
			continue
		}
		_, err = i.etcdCli.Put(ctx, fmt.Sprintf("/kato/gateway/ips/%s", ip.String()), ip.String(), clientv3.WithLease(res.ID))
		if err != nil {
			logrus.Errorf("put gateway ip to etcd failure %s", err.Error())
			continue
		}
		logrus.Infof("gateway init add ip %s", ip.String())
		i.ipLease[ip.String()] = res.ID
	}
	return nil
}

func (i *ipManager) deleteIP(ips ...net.IP) {
	ctx, cancel := context.WithTimeout(i.ctx, time.Second*30)
	defer cancel()
	i.lock.Lock()
	defer i.lock.Unlock()
	for _, ip := range ips {
		_, err := i.etcdCli.Delete(ctx, fmt.Sprintf("/kato/gateway/ips/%s", ip.String()))
		if err != nil {
			logrus.Errorf("put gateway ip to etcd failure %s", err.Error())
		}
		delete(i.ipLease, ip.String())
	}
}

func (i *ipManager) Stop() {
	i.cancel()
}
