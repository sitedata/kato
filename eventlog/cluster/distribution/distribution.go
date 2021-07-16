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

package distribution

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/eventlog/cluster/discover"
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/db"
	"sync"
	"time"

	"golang.org/x/net/context"

	"sort"

	"github.com/sirupsen/logrus"
)

//Distribution - data partition
type Distribution struct {
	monitorDatas map[string]*db.MonitorData
	updateTime   map[string]time.Time
	abnormalNode map[string]int
	lock         sync.Mutex
	cancel       func()
	context      context.Context
	discover     discover.Manager
	log          *logrus.Entry
	etcdClient   *clientv3.Client
	conf         conf.DiscoverConf
}

func NewDistribution(etcdClient *clientv3.Client, conf conf.DiscoverConf, dis discover.Manager, log *logrus.Entry) *Distribution {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Distribution{
		cancel:       cancel,
		context:      ctx,
		discover:     dis,
		monitorDatas: make(map[string]*db.MonitorData),
		updateTime:   make(map[string]time.Time),
		abnormalNode: make(map[string]int),
		log:          log,
		etcdClient:   etcdClient,
		conf:         conf,
	}
	return d
}

//Start
func (d *Distribution) Start() error {
	go d.checkHealth()
	return nil
}

//Stop
func (d *Distribution) Stop() {
	d.cancel()
}

//Update
func (d *Distribution) Update(m db.MonitorData) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if m.InstanceID == "" {
		d.log.Warning("update monitor data but instance id is empty.")
		return
	}
	if md, ok := d.monitorDatas[m.InstanceID]; ok {
		md.LogSizePeerM = m.LogSizePeerM
		md.ServiceSize = m.ServiceSize
		if _, ok := d.abnormalNode[m.InstanceID]; ok {
			delete(d.abnormalNode, m.InstanceID)
		}
	} else {
		d.monitorDatas[m.InstanceID] = &m
	}
	d.updateTime[m.InstanceID] = time.Now()
}

func (d *Distribution) checkHealth() {
	tike := time.Tick(time.Second * 5)
	for {
		select {
		case <-tike:
		case <-d.context.Done():
			return
		}
		d.lock.Lock()
		for k, v := range d.updateTime {
			if v.Add(time.Second * 10).Before(time.Now()) { //Node offline or node failure
				status := d.discover.InstanceCheckHealth(k)
				if status == "delete" {
					delete(d.monitorDatas, k)
					delete(d.updateTime, k)
					d.log.Warnf("instance (%s) health is delete.", k)
				}
				if status == "abnormal" {
					d.abnormalNode[k] = 1
					d.log.Warnf("instance (%s) health is abnormal.", k)
				}
			}
		}
		d.lock.Unlock()
	}
}

//GetSuitableInstance - get recommended nodes
func (d *Distribution) GetSuitableInstance(serviceID string) *discover.Instance {
	d.lock.Lock()
	defer d.lock.Unlock()
	var suitableInstance *discover.Instance

	instanceID, err := discover.GetDokerLogInInstance(d.etcdClient, d.conf, serviceID)
	if err != nil {
		d.log.Error("Get docker log in instance id error ", err.Error())
	}
	if instanceID != "" {
		if _, ok := d.abnormalNode[instanceID]; !ok {
			if _, ok := d.monitorDatas[instanceID]; ok {
				suitableInstance = d.discover.GetInstance(instanceID)
				if suitableInstance != nil {
					return suitableInstance
				}
			}
		}
	}
	if len(d.monitorDatas) < 1 {
		ins := d.discover.GetCurrentInstance()
		d.log.Debug("monitor data length <1 return self")
		return &ins
	}
	d.log.Debug("start select suitable Instance")
	var flags []int
	var instances = make(map[int]*discover.Instance)
	for k, v := range d.monitorDatas {
		if _, ok := d.abnormalNode[k]; !ok {
			if ins := d.discover.GetInstance(k); ins != nil {
				flag := int(v.LogSizePeerM) + 20*v.ServiceSize
				flags = append(flags, flag)
				instances[flag] = ins
			} else {
				d.log.Debugf("instance %s stat is delete", k)
			}
		} else {
			d.log.Debugf("instance %s stat is abnormal", k)
		}
	}

	if len(flags) > 0 {
		sort.Ints(flags)
		suitableInstance = instances[flags[0]]
	}
	if suitableInstance == nil {
		d.log.Debug("suitableInstance is nil return self")
		ins := d.discover.GetCurrentInstance()
		return &ins
	}
	return suitableInstance
}
