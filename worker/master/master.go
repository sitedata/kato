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

package master

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gridworkz/kato/cmd/worker/option"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util/leader"
	"github.com/gridworkz/kato/worker/appm/store"
	"github.com/gridworkz/kato/worker/master/podevent"
	"github.com/gridworkz/kato/worker/master/volumes/provider"
	"github.com/gridworkz/kato/worker/master/volumes/provider/lib/controller"
	"github.com/gridworkz/kato/worker/master/volumes/statistical"
	"github.com/gridworkz/kato/worker/master/volumes/sync"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//Controller app runtime master controller
type Controller struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	conf                option.Config
	store               store.Storer
	dbmanager           db.Manager
	memoryUse           *prometheus.GaugeVec
	cpuUse              *prometheus.GaugeVec
	fsUse               *prometheus.GaugeVec
	diskCache           *statistical.DiskCache
	namespaceMemRequest *prometheus.GaugeVec
	namespaceMemLimit   *prometheus.GaugeVec
	namespaceCPURequest *prometheus.GaugeVec
	namespaceCPULimit   *prometheus.GaugeVec
	pc                  *controller.ProvisionController
	isLeader            bool

	kubeClient kubernetes.Interface

	stopCh          chan struct{}
	podEvent * podevent.PodEvent
	volumeTypeEvent *sync.VolumeTypeEvent

	version      *version.Info
	katosssc controller.Provisions
	katosslc controller.Provisions
}

//NewMasterController new master controller
func NewMasterController(conf option.Config, kubecfg *rest.Config, store store.Storer) (*Controller, error) {
	// kubecfg.RateLimiter = nil
	kubeClient, err := kubernetes.NewForConfig(kubecfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	// The controller needs to know what the server version is because out-of-tree
	// provisioners aren't officially supported until 1.5
	serverVersion, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		logrus.Errorf("Error getting server version: %v", err)
		cancel()
		return nil, err
	}

	// Create the provisioner: it implements the Provisioner interface expected by
	// the controller
	//statefulset share controller
	katossscProvisioner := provider.NewKatossscProvisioner()
	//statefulset local controller
	katosslcProvisioner := provider.NewKatosslcProvisioner(kubeClient, store)
	// Start the provision controller which will dynamically provision hostPath
	// PVs
	pc := controller.NewProvisionController(kubeClient, &conf, map[string]controller.Provisioner{
		katossscProvisioner.Name(): katossscProvisioner,
		katosslcProvisioner.Name(): katosslcProvisioner,
	}, serverVersion.GitVersion)
	stopCh := make(chan struct{})

	return &Controller{
		conf: conf,
		pc:        pc,
		store: store,
		stopCh:    stopCh,
		cancel:    cancel,
		ctx:       ctx,
		dbmanager: db.GetManager(),
		memoryUse: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "app_resource",
			Name:      "appmemory",
			Help:      "tenant service memory request.",
		}, []string{"tenant_id", "app_id", "service_id", "service_status"}),
		cpuUse: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "app_resource",
			Name:      "appcpu",
			Help:      "tenant service cpu request.",
		}, []string{"tenant_id", "app_id", "service_id", "service_status"}),
		fsUse: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "app_resource",
			Name:      "appfs",
			Help:      "tenant service fs used.",
		}, []string{"tenant_id", "app_id", "service_id", "volume_type"}),
		namespaceMemRequest: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "namespace_resource",
			Name:      "memory_request",
			Help:      "total memory request in namespace",
		}, []string{"namespace"}),
		namespaceMemLimit: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "namespace_resource",
			Name:      "memory_limit",
			Help:      "total memory limit in namespace",
		}, []string{"namespace"}),
		namespaceCPURequest: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "namespace_resource",
			Name:      "cpu_request",
			Help:      "total cpu request in namespace",
		}, []string{"namespace"}),
		namespaceCPULimit: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "namespace_resource",
			Name:      "cpu_limit",
			Help:      "total cpu limit in namespace",
		}, []string{"namespace"}),
		diskCache:       statistical.CreatDiskCache(ctx),
		podEvent:        podevent.New(conf.KubeClient, stopCh),
		volumeTypeEvent: sync.New(stopCh),
		kubeClient: kubeClient,
		katosssc:    katossscProvisioner,
		katosslc:    katosslcProvisioner,
		version:         serverVersion,
	}, nil
}

//IsLeader is leader
func (m *Controller) IsLeader() bool {
	return m.isLeader
}

//Start start
func (m *Controller) Start() error {
	logrus.Debug("master controller starting")
	start := func(ctx context.Context) {
		pc := controller.NewProvisionController(m.kubeClient, &m.conf, map[string]controller.Provisioner{
			m.katosslc.Name(): m.katosslc,
			m.katosssc.Name(): m.katosssc,
		}, m.version.GitVersion)

		m.isLeader = true
		defer func() {
			m.isLeader = false
		}()
		go m.diskCache.Start()
		defer m.diskCache.Stop()
		go pc.Run(ctx)
		m.store.RegistPodUpdateListener("podEvent", m.podEvent.GetChan())
		defer m.store.UnRegistPodUpdateListener("podEvent")
		go m.podEvent.Handle()
		m.store.RegisterVolumeTypeListener("volumeTypeEvent", m.volumeTypeEvent.GetChan())
		defer m.store.UnRegisterVolumeTypeListener("volumeTypeEvent")
		go m.volumeTypeEvent.Handle()

		select {
		case <-ctx.Done():
		case <-m.ctx.Done():
		}
	}
	// Leader election was requested.
	if m.conf.LeaderElectionNamespace == "" {
		return fmt.Errorf("-leader-election-namespace must not be empty")
	}
	if m.conf.LeaderElectionIdentity == "" {
		m.conf.LeaderElectionIdentity = m.conf.NodeName
	}
	if m.conf.LeaderElectionIdentity == "" {
		return fmt.Errorf("-leader-election-identity must not be empty")
	}
	// Name of config map with leader election lock
	lockName := "kato-appruntime-worker-leader"

	// Become leader again on stop leading.
	leaderCh := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-m.ctx.Done():
				return
			case <-leaderCh:
				logrus.Info("run as leader")
				ctx, cancel := context.WithCancel(m.ctx)
				defer cancel()
				leader.RunAsLeader(ctx, m.kubeClient, m.conf.LeaderElectionNamespace, m.conf.LeaderElectionIdentity, lockName, start, func() {
					leaderCh <- struct{}{}
					logrus.Info("restart leader")
				})
			}
		}
	}()

	leaderCh <- struct{}{}

	return nil
}

//Stop stop
func (m *Controller) Stop() {
	close(m.stopCh)
}

//Scrape scrape app runtime
func (m *Controller) Scrape(ch chan<- prometheus.Metric, scrapeDurationDesc *prometheus.Desc) {
	if !m.isLeader {
		return
	}
	scrapeTime := time.Now()
	services := m.store.GetAllAppServices()
	status := m.store.GetNeedBillingStatus(nil)
	//Get memory usage
	for _, service := range services {
		if _, ok := status[service.ServiceID]; ok {
			m.memoryUse.WithLabelValues(service.TenantID, service.AppID, service.ServiceID, "running").Set(float64(service.GetMemoryRequest()))
			m.cpuUse.WithLabelValues(service.TenantID, service.AppID, service.ServiceID, "running").Set(float64(service.GetMemoryRequest()))
		}
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.memory")
	scrapeTime = time.Now()
	diskcache := m.diskCache.Get()
	for k, v: = range diskcache {
		key := strings.Split(k, "_")
		if len(key) == 3 {
			m.fsUse.WithLabelValues(key[2], key[1], key[0], string(model.ShareFileVolumeType)).Set(v)
		}
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), "collect.fs")
	resources := m.store.GetTenantResourceList()
	for _, re := range resources {
		m.namespaceMemLimit.WithLabelValues(re.Namespace).Set(float64(re.MemoryLimit / 1024 / 1024))
		m.namespaceCPULimit.WithLabelValues(re.Namespace).Set(float64(re.CPULimit))
		m.namespaceMemRequest.WithLabelValues(re.Namespace).Set(float64(re.MemoryRequest / 1024 / 1024))
		m.namespaceCPURequest.WithLabelValues(re.Namespace).Set(float64(re.CPURequest))
	}
	m.fsUse.Collect(ch)
	m.memoryUse.Collect(ch)
	m.cpuUse.Collect(ch)
	m.namespaceMemLimit.Collect(ch)
	m.namespaceCPULimit.Collect(ch)
	m.namespaceMemRequest.Collect(ch)
	m.namespaceCPURequest.Collect(ch)
	logrus.Infof("success collect worker master metric")
}
