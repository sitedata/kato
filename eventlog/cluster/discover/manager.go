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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/util"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//Manager - dynamic node discovery manager
type Manager interface {
	RegisteredInstance(host string, port int, stopRegister *bool) *Instance
	CancellationInstance(instance *Instance)
	MonitorAddInstances() chan *Instance
	MonitorDelInstances() chan *Instance
	MonitorUpdateInstances() chan *Instance
	GetInstance(string) *Instance
	InstanceCheckHealth(string) string
	Run() error
	GetCurrentInstance() Instance
	Stop()
	Scrape(ch chan<- prometheus.Metric, namespace, exporter string) error
}

//EtcdDiscoverManager - automatic discovery based on ETCD
type EtcdDiscoverManager struct {
	cancel         func()
	context        context.Context
	addChan        chan *Instance
	delChan        chan *Instance
	updateChan     chan *Instance
	log            *logrus.Entry
	conf           conf.DiscoverConf
	etcdclientv3   *clientv3.Client
	selfInstance   *Instance
	othersInstance []*Instance
	stopDiscover   bool
}

//New
func New(etcdClient *clientv3.Client, conf conf.DiscoverConf, log *logrus.Entry) Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &EtcdDiscoverManager{
		conf:           conf,
		cancel:         cancel,
		context:        ctx,
		log:            log,
		addChan:        make(chan *Instance, 2),
		delChan:        make(chan *Instance, 2),
		updateChan:     make(chan *Instance, 2),
		othersInstance: make([]*Instance, 0),
		etcdclientv3:   etcdClient,
	}
}

//GetCurrentInstance
func (d *EtcdDiscoverManager) GetCurrentInstance() Instance {
	return *d.selfInstance
}

//RegisteredInstance
func (d *EtcdDiscoverManager) RegisteredInstance(host string, port int, stopRegister *bool) *Instance {
	instance := &Instance{}
	for !*stopRegister {
		if host == "0.0.0.0" || host == "127.0.0.1" || host == "localhost" {
			if d.conf.InstanceIP != "" {
				ip := net.ParseIP(d.conf.InstanceIP)
				if ip != nil {
					instance.HostIP = ip
				}
			}
		} else {
			ip := net.ParseIP(host)
			if ip != nil {
				instance.HostIP = ip
			}
		}
		if instance.HostIP == nil {
			ip, err := util.ExternalIP()
			if err != nil {
				d.log.Error("Can not get host ip for the instance.")
				time.Sleep(time.Second * 10)
				continue
			} else {
				instance.HostIP = ip
			}
		}
		instance.PubPort = port
		instance.DockerLogPort = d.conf.DockerLogPort
		instance.WebPort = d.conf.WebPort
		instance.HostID = d.conf.NodeID
		instance.HostName, _ = os.Hostname()
		instance.Status = "create"
		data, err := json.Marshal(instance)
		if err != nil {
			d.log.Error("Create register instance data error.", err.Error())
			time.Sleep(time.Second * 10)
			continue
		}
		ctx, cancel := context.WithCancel(d.context)
		_, err = d.etcdclientv3.Put(ctx, fmt.Sprintf("%s/instance/%s:%d", d.conf.HomePath, instance.HostIP, instance.PubPort), string(data))
		cancel()
		if err != nil {
			d.log.Error("Register instance data to etcd error.", err.Error())
			time.Sleep(time.Second * 10)
			continue
		}

		d.selfInstance = instance
		go d.discover()
		d.log.Infof("Register instance in cluster success. HostID:%s HostIP:%s PubPort:%d", instance.HostID, instance.HostIP, instance.PubPort)
		return instance
	}
	return nil
}

//MonitorAddInstances
func (d *EtcdDiscoverManager) MonitorAddInstances() chan *Instance {
	return d.addChan
}

//MonitorDelInstances
func (d *EtcdDiscoverManager) MonitorDelInstances() chan *Instance {
	return d.delChan
}

//MonitorUpdateInstances
func (d *EtcdDiscoverManager) MonitorUpdateInstances() chan *Instance {
	return d.updateChan
}

//Run
func (d *EtcdDiscoverManager) Run() error {
	d.log.Info("Discover manager start ")
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: d.conf.EtcdAddr,
		CaFile:    d.conf.EtcdCaFile,
		CertFile:  d.conf.EtcdCertFile,
		KeyFile:   d.conf.EtcdKeyFile,
	}
	var err error
	d.etcdclientv3, err = etcdutil.NewClient(d.context, etcdClientArgs)
	if err != nil {
		d.log.Error("Create etcd v3 client error.", err.Error())
		return err
	}
	return nil
}

//Discover
func (d *EtcdDiscoverManager) discover() {
	tike := time.NewTicker(time.Second * 5)
	defer tike.Stop()
	for {
		ctx, cancel := context.WithCancel(d.context)
		res, err := d.etcdclientv3.Get(ctx, d.conf.HomePath+"/instance/", clientv3.WithPrefix())
		cancel()
		if err != nil {
			d.log.Error("Get instance info from etcd error.", err.Error())
		} else {
			for _, kv := range res.Kvs {
				node := &Node{
					Key:   string(kv.Key),
					Value: string(kv.Value),
				}
				d.add(node)
			}
			break
		}
		select {
		case <-tike.C:
		case <-d.context.Done():
			return
		}
	}
	ctx, cancel := context.WithCancel(d.context)
	defer cancel()
	watcher := d.etcdclientv3.Watch(ctx, d.conf.HomePath+"/instance/", clientv3.WithPrefix())

	for !d.stopDiscover {
		res, ok := <-watcher
		if !ok {
			break
		}

		for _, event := range res.Events {
			node := &Node{
				Key:   string(event.Kv.Key),
				Value: string(event.Kv.Value),
			}
			switch event.Type {
			case mvccpb.PUT:
				d.add(node)
			case mvccpb.DELETE:
				//Ignore yourself
				if strings.HasSuffix(node.Key, fmt.Sprintf("/%s:%d", d.selfInstance.HostIP, d.selfInstance.PubPort)) {
					continue
				}
				keys := strings.Split(node.Key, "/")
				hostPort := keys[len(keys)-1]
				d.log.Infof("Delete an instance.%s", hostPort)
				var removeIndex int
				var have bool
				for i, ins := range d.othersInstance {
					if fmt.Sprintf("%s:%d", ins.HostIP, ins.PubPort) == hostPort {
						removeIndex = i
						have = true
						break
					}
				}
				if have {
					instance := d.othersInstance[removeIndex]
					d.othersInstance = DeleteSlice(d.othersInstance, removeIndex)
					d.MonitorDelInstances() <- instance
					d.log.Infof("A instance offline %s", instance.HostName)
				}
			}
		}

	}
	d.log.Debug("discover manager discover core stop")
}

type Node struct {
	Key   string
	Value string
}

func (d *EtcdDiscoverManager) add(node *Node) {

	//Ignore yourself
	if strings.HasSuffix(node.Key, fmt.Sprintf("/%s:%d", d.selfInstance.HostIP, d.selfInstance.PubPort)) {
		return
	}
	var instance Instance
	if err := json.Unmarshal([]byte(node.Value), &instance); err != nil {
		d.log.Error("Unmarshal instance data that from etcd error.", err.Error())
	} else {
		if strings.HasSuffix(node.Key, fmt.Sprintf("/%s:%d", d.selfInstance.HostIP, d.selfInstance.PubPort)) {
			d.selfInstance = &instance
		}

		isExist := false
		for _, i := range d.othersInstance {
			if i.HostID == instance.HostID {
				*i = instance
				d.log.Debug("update the instance " + i.HostID)
				isExist = true
				break
			}
		}

		if !isExist {
			d.log.Infof("Find an instance.IP:%s, Port:%d, NodeID:%s HostID: %s", instance.HostIP.String(), instance.PubPort, instance.HostName, instance.HostID)
			d.MonitorAddInstances() <- &instance
			d.othersInstance = append(d.othersInstance, &instance)
		}
	}

}

func (d *EtcdDiscoverManager) update(node *Node) {

	var instance Instance
	if err := json.Unmarshal([]byte(node.Value), &instance); err != nil {
		d.log.Error("Unmarshal instance data that from etcd error.", err.Error())
	} else {
		if strings.HasSuffix(node.Key, fmt.Sprintf("/%s:%d", d.selfInstance.HostIP, d.selfInstance.PubPort)) {
			d.selfInstance = &instance
		}
		for _, i := range d.othersInstance {
			if i.HostID == instance.HostID {
				*i = instance
				d.log.Debug("update the instance " + i.HostID)
			}
		}
	}

}

//DeleteSlice - delete an element from the array
func DeleteSlice(source []*Instance, index int) []*Instance {
	if len(source) == 1 {
		return make([]*Instance, 0)
	}
	if index == 0 {
		return source[1:]
	}
	if index == len(source)-1 {
		return source[:len(source)-2]
	}
	return append(source[0:index-1], source[index+1:]...)
}

//Stop
func (d *EtcdDiscoverManager) Stop() {
	d.stopDiscover = true
	d.cancel()
	d.log.Info("Stop the discover manager.")
}

//CancellationInstance
func (d *EtcdDiscoverManager) CancellationInstance(instance *Instance) {
	ctx, cancel := context.WithTimeout(d.context, time.Second*5)
	defer cancel()
	_, err := d.etcdclientv3.Delete(ctx, fmt.Sprintf("%s/instance/%s:%d", d.conf.HomePath, instance.HostIP, instance.PubPort))
	if err != nil {
		d.log.Error("Cancellation Instance from etcd error.", err.Error())
	} else {
		d.log.Info("Cancellation Instance from etcd")
	}
}

//UpdateInstance
func (d *EtcdDiscoverManager) UpdateInstance(instance *Instance) {
	instance.Status = "update"
	data, err := json.Marshal(instance)
	if err != nil {
		d.log.Error("Create update instance data error.", err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(d.context, time.Second*5)
	defer cancel()
	_, err = d.etcdclientv3.Put(ctx, fmt.Sprintf("%s/instance/%s:%d", d.conf.HomePath, instance.HostIP, instance.PubPort), string(data))
	if err != nil {
		d.log.Error(" Update Instance from etcd error.", err.Error())
	}
}

//InstanceCheckHealth - will be called by distribution, when the node is found to be abnormal
//Check here, if the node is offline, return delete
//If the node is not offline, mark it as abnormal, return abnormal
//If the node is judged to be faulty by the cluster, return delete
func (d *EtcdDiscoverManager) InstanceCheckHealth(instanceID string) string {
	d.log.Info("Start check instance health.")
	if d.selfInstance.HostID == instanceID {
		d.log.Error("The current node condition monitoring.")
		return "abnormal"
	}
	for _, i := range d.othersInstance {
		if i.HostID == instanceID {
			d.log.Errorf("Instance (%s) is abnormal.", instanceID)
			i.Status = "abnormal"
			i.TagNumber++
			if i.TagNumber > ((len(d.othersInstance) + 1) / 2) { //Node mark greater than half
				d.log.Warn("Instance (%s) is abnormal. tag number more than half of all instance number. will cancellation.", instanceID)
				d.CancellationInstance(i)
				return "delete"
			}
			d.UpdateInstance(i)
			return "abnormal"
		}
	}
	return "delete"
}

//GetInstance
func (d *EtcdDiscoverManager) GetInstance(id string) *Instance {
	if id == d.selfInstance.HostID {
		return d.selfInstance
	}
	for _, i := range d.othersInstance {
		if i.HostID == id {
			return i
		}
	}
	return nil
}

//Scrape prometheus monitor metrics
func (d *EtcdDiscoverManager) Scrape(ch chan<- prometheus.Metric, namespace, exporter string) error {
	instanceDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "instance_up"),
		"the instance in cluster status.",
		[]string{"from", "instance", "status"}, nil,
	)
	for _, i := range d.othersInstance {
		if i.Status == "delete" || i.Status == "abnormal" {
			ch <- prometheus.MustNewConstMetric(instanceDesc, prometheus.GaugeValue, 0, d.selfInstance.HostIP.String(), i.HostIP.String(), i.Status)
		} else {
			ch <- prometheus.MustNewConstMetric(instanceDesc, prometheus.GaugeValue, 1, d.selfInstance.HostIP.String(), i.HostIP.String(), i.Status)
		}
	}
	return nil
}
