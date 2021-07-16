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

package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gridworkz/kato/gateway/cluster"

	client "github.com/coreos/etcd/clientv3"
	"github.com/eapache/channels"
	"github.com/gridworkz/kato/cmd/gateway/option"
	"github.com/gridworkz/kato/gateway/controller/openresty"
	"github.com/gridworkz/kato/gateway/metric"
	"github.com/gridworkz/kato/gateway/store"
	v1 "github.com/gridworkz/kato/gateway/v1"
	"github.com/gridworkz/kato/util/ingress-nginx/task"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/flowcontrol"
)

// kato endpoints map
var rbdemap = make(map[string]struct{})

func init() {
	rbdemap["APISERVER_ENDPOINTS"] = struct{}{}
	rbdemap["HUB_ENDPOINTS"] = struct{}{}
	rbdemap["REPO_ENDPOINTS"] = struct{}{}
}

// GWController -
type GWController struct {
	GWS   GWServicer
	store store.Storer

	syncQueue       *task.Queue
	syncRateLimiter flowcontrol.RateLimiter
	isShuttingDown  bool

	// stopLock is used to enforce that only a single call to Stop send at
	// a given time. We allow stopping through an HTTP endpoint and
	// allowing concurrent stoppers leads to stack traces.
	stopLock *sync.Mutex

	ocfg *option.Config
	rcfg *v1.Config // running configuration
	rrhp []*v1.Pool // running kato http pools
	rrtp []*v1.Pool // running kato tcp or udp pools

	stopCh   chan struct{}
	updateCh *channels.RingChannel

	EtcdCli *client.Client
	ctx     context.Context

	metricCollector metric.Collector
}

// Start Gateway
func (gwc *GWController) Start(errCh chan error) error {
	// start plugin(eg: nginx, zeus and etc)
	if err := gwc.GWS.Start(errCh); err != nil {
		return err
	}
	// start informer
	gwc.store.Run(gwc.stopCh)

	// start task queue
	go gwc.syncQueue.Run(1*time.Second, gwc.stopCh)

	// force initial sync
	gwc.syncQueue.EnqueueTask(task.GetDummyObject("initial-sync"))

	go gwc.handleEvent()

	return nil
}

// Close stops Gateway
func (gwc *GWController) Close() error {
	gwc.isShuttingDown = true
	if gwc.EtcdCli != nil {
		gwc.EtcdCli.Close()
	}
	gwc.stopLock.Lock()
	defer gwc.stopLock.Unlock()

	if gwc.syncQueue.IsShuttingDown() {
		return fmt.Errorf("shutdown already in progress")
	}

	logrus.Infof("Shutting down controller queues")
	close(gwc.stopCh) // stop the loop in *GWController#Start()
	go gwc.syncQueue.Shutdown()

	return gwc.GWS.Stop()
}

func (gwc *GWController) handleEvent() {
	for {
		select {
		case event := <-gwc.updateCh.Out(): // received k8s events
			if gwc.isShuttingDown {
				break
			}
			if evt, ok := event.(store.Event); ok {
				gwc.syncQueue.EnqueueSkippableTask(evt.Obj)
			} else {
				logrus.Warningf("Unexpected event type received %T", event)
			}
		case <-gwc.stopCh:
			break
		}
	}
}

func (gwc *GWController) syncGateway(key interface{}) error {
	gwc.syncRateLimiter.Accept()

	if gwc.syncQueue.IsShuttingDown() {
		return nil
	}
	l7sv, l4sv := gwc.store.ListVirtualService()
	httpPools, tcpPools := gwc.store.ListPool()
	currentConfig := &v1.Config{
		HTTPPools: httpPools,
		TCPPools:  tcpPools,
		L7VS:      l7sv,
		L4VS:      l4sv,
	}
	// refresh http tcp and udp pools dynamically
	httpPools = append(httpPools, gwc.rrhp...)
	tcpPools = append(tcpPools, gwc.rrtp...)
	if err := gwc.GWS.UpdatePools(httpPools, tcpPools); err != nil {
		logrus.Warningf("error updating pools: %v", err)
	}
	if gwc.rcfg.Equals(currentConfig) {
		logrus.Debug("No need to update running configuration.")
		return nil
	}
	logrus.Infof("update nginx server config file.")
	err := gwc.GWS.PersistConfig(currentConfig)
	if err != nil {
		// TODO: if nginx is not ready, then stop gateway
		logrus.Errorf("Fail to persist Nginx config: %v\n", err)
		return nil
	}

	//set metric
	remove, hosts := getHosts(gwc.rcfg, currentConfig)
	gwc.metricCollector.SetHosts(hosts)
	gwc.metricCollector.RemoveHostMetric(remove)
	gwc.metricCollector.SetServerNum(len(httpPools), len(tcpPools))

	gwc.rcfg = currentConfig
	return nil
}

//NewGWController new Gateway controller
func NewGWController(ctx context.Context, clientset kubernetes.Interface, cfg *option.Config, mc metric.Collector, node *cluster.NodeManager) (*GWController, error) {
	gwc := &GWController{
		updateCh:        channels.NewRingChannel(1024),
		syncRateLimiter: flowcontrol.NewTokenBucketRateLimiter(cfg.SyncRateLimit, 1),
		stopLock:        &sync.Mutex{},
		stopCh:          make(chan struct{}),
		ocfg:            cfg,
		ctx:             ctx,
		metricCollector: mc,
	}

	gwc.GWS = openresty.CreateOpenrestyService(cfg, &gwc.isShuttingDown)

	gwc.store = store.New(
		clientset,
		gwc.updateCh,
		cfg, node)
	gwc.syncQueue = task.NewTaskQueue(gwc.syncGateway)

	return gwc, nil
}

func poolsEqual(a []*v1.Pool, b []*v1.Pool) bool {
	if len(a) != len(b) {
		return false
	}
	for _, ap := range a {
		flag := false
		for _, bp := range b {
			if ap.Equals(bp) {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

// getHosts returns a list of the hostsnames and tobe remove hostname
// that are not associated anymore to the NGINX configuration.
func getHosts(rucfg, newcfg *v1.Config) (remove []string, current sets.String) {
	old := sets.NewString()
	new := sets.NewString()
	if rucfg != nil {
		for _, s := range rucfg.L7VS {
			if !old.Has(s.ServerName) {
				old.Insert(s.ServerName)
			}
		}
	}
	if newcfg != nil {
		for _, s := range newcfg.L7VS {
			if !new.Has(s.ServerName) {
				new.Insert(s.ServerName)
			}
		}
	}
	return old.Difference(new).List(), new
}
