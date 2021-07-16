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

package controller

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/builder/parser"
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/node/nodem/client"
	"github.com/gridworkz/kato/node/nodem/healthy"
	"github.com/gridworkz/kato/node/nodem/service"
	"github.com/sirupsen/logrus"
)

var (
	// ArgsReg -
	ArgsReg = regexp.MustCompile(`\$\{(\w+)\|{0,1}(.{0,1})\}`)
)

//ManagerService
type ManagerService struct {
	node                 *client.HostNode
	ctx                  context.Context
	cancel               context.CancelFunc
	syncCtx              context.Context
	syncCancel           context.CancelFunc
	conf                 *option.Conf
	ctr                  Controller
	cluster              client.ClusterClient
	healthyManager       healthy.Manager
	services             []*service.Service
	allservice           []*service.Service
	etcdcli              *clientv3.Client
	autoStatusController map[string]statusController
	lock                 sync.Mutex
}

//GetAllService
func (m *ManagerService) GetAllService() ([]*service.Service, error) {
	return m.allservice, nil
}

//GetService
func (m *ManagerService) GetService(serviceName string) *service.Service {
	for _, s := range m.allservice {
		if s.Name == serviceName {
			return s
		}
	}
	return nil
}

//Start and monitor all service
func (m *ManagerService) Start(node *client.HostNode) error {
	logrus.Info("Starting node controller manager.")
	m.loadServiceConfig()
	m.node = node
	if m.conf.EnableInitStart {
		return m.ctr.InitStart(m.services)
	}
	return nil
}

func (m *ManagerService) loadServiceConfig() {
	m.allservice = service.LoadServicesFromLocal(m.conf.ServiceListFile)
	var controllerServices []*service.Service
	for _, s := range m.allservice {
		if !s.OnlyHealthCheck && !s.Disable {
			controllerServices = append(controllerServices, s)
		}
	}
	m.services = controllerServices
}

//Stop manager
func (m *ManagerService) Stop() error {
	m.cancel()
	return nil
}

//Online start all service on the node
func (m *ManagerService) Online() error {
	logrus.Info("Doing node online by node controller manager")
	if ok := m.ctr.CheckBeforeStart(); !ok {
		logrus.Debug("check before starting: false")
		return nil
	}
	go m.StartServices()
	m.SyncServiceStatusController()
	// register local services endpoint into cluster manager
	for _, s := range m.services {
		m.UpOneServiceEndpoint(s)
	}
	return nil
}

// SetEndpoints registers endpoints in etcd
func (m *ManagerService) SetEndpoints(hostIP string) {
	if hostIP == "" {
		logrus.Warningf("ignore wrong hostIP: %s", hostIP)
		return
	}
	for _, s := range m.services {
		m.UpOneServiceEndpoint(s)
	}
}

//StartServices
func (m *ManagerService) StartServices() {
	for _, service := range m.services {
		if !service.Disable {
			logrus.Infof("Begin start service %s", service.Name)
			update, err := m.ctr.WriteConfig(service)
			if err != nil {
				logrus.Errorf("write service config failure %s", err.Error())
				continue
			}
			if update {
				if err := m.ctr.RestartService(service); err != nil {
					logrus.Errorf("start service failure %s", err.Error())
					continue
				}
			} else {
				if err := m.ctr.StartService(service.Name); err != nil {
					logrus.Errorf("start service failure %s", err.Error())
					continue
				}
			}
		}
	}
}

//Offline stop all service of on the node
func (m *ManagerService) Offline() error {
	logrus.Info("Doing node offline by node controller manager")
	services, _ := m.GetAllService()
	for _, s := range services {
		m.DownOneServiceEndpoint(s)
	}
	m.StopSyncService()
	if err := m.ctr.StopList(m.services); err != nil {
		return err
	}
	return nil
}

//DownOneServiceEndpoint
func (m *ManagerService) DownOneServiceEndpoint(s *service.Service) {
	hostIP := m.cluster.GetOptions().HostIP
	for _, end := range s.Endpoints {
		logrus.Debug("Anti-registry endpoint: ", end.Name)
		key := end.Name + "/" + hostIP
		endpoint := toEndpoint(end, hostIP)
		oldEndpoints := m.cluster.GetEndpoints(key)
		if exist := isExistEndpoint(oldEndpoints, endpoint); exist {
			endpoints := rmEndpointFrom(oldEndpoints, endpoint)
			if len(endpoints) > 0 {
				m.cluster.SetEndpoints(end.Name, m.cluster.GetOptions().HostIP, endpoints)
				continue
			}
			m.cluster.DelEndpoints(key)
		}
	}
	logrus.Infof("node %s down service %s endpoints", hostIP, s.Name)
}

//UpOneServiceEndpoint
func (m *ManagerService) UpOneServiceEndpoint(s *service.Service) {
	if s.OnlyHealthCheck || s.Disable {
		return
	}
	hostIP := m.cluster.GetOptions().HostIP
	for _, end := range s.Endpoints {
		if end.Name == "" || strings.Replace(end.Port, " ", "", -1) == "" {
			continue
		}
		endpoint := toEndpoint(end, hostIP)
		m.cluster.SetEndpoints(end.Name, hostIP, []string{endpoint})
	}
}

//SyncServiceStatusController synchronize all service status to as we expect
func (m *ManagerService) SyncServiceStatusController() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.autoStatusController != nil && len(m.autoStatusController) > 0 {
		for _, v := range m.autoStatusController {
			v.Stop()
		}
	}
	m.autoStatusController = make(map[string]statusController, len(m.services))
	for _, s := range m.services {
		if s.ServiceHealth == nil {
			continue
		}
		ctx, cancel := context.WithCancel(context.Background())
		serviceStatusController := statusController{
			ctx:            ctx,
			cancel:         cancel,
			service:        s,
			healthyManager: m.healthyManager,
			watcher:        m.healthyManager.WatchServiceHealthy(s.Name),
			unhealthHandle: func(event *service.HealthStatus, w healthy.Watcher) {
				service := m.GetService(event.Name)
				if service == nil {
					logrus.Errorf("not found service %s", event.Name)
					return
				}
				m.DownOneServiceEndpoint(service)
				if service.OnlyHealthCheck {
					logrus.Warningf("service %s is only check health.so do not auto restart it", event.Name)
					return
				}
				if event.Name == "docker" {
					logrus.Errorf("service docker can not auto restart. must artificial processing")
					return
				}
				// disable check healthy status of the service
				logrus.Infof("service %s not healthy, will restart it", event.Name)
				m.healthyManager.DisableWatcher(event.Name, w.GetID())
				_, err := m.ctr.WriteConfig(service)
				if err == nil {
					if err := m.ctr.RestartService(service); err != nil {
						logrus.Errorf("restart service %s failure %s", event.Name, err.Error())
					} else {
						if !m.WaitStart(event.Name, time.Minute) {
							logrus.Errorf("Timeout restart service: %s, will recheck health", event.Name)
						}
					}
				} else {
					logrus.Errorf("update service %s systemctl config failure where restart it:%s", event.Name, err.Error())
				}
				// start check healthy status of the service
				m.healthyManager.EnableWatcher(event.Name, w.GetID())
			},
			healthHandle: func(event *service.HealthStatus, w healthy.Watcher) {
				service := m.GetService(event.Name)
				if service == nil {
					logrus.Errorf("not found service %s", event.Name)
					return
				}
				m.UpOneServiceEndpoint(service)
			},
		}
		m.autoStatusController[s.Name] = serviceStatusController
		go serviceStatusController.Run()
	}
}

type statusController struct {
	watcher        healthy.Watcher
	ctx            context.Context
	cancel         context.CancelFunc
	service        *service.Service
	unhealthHandle func(event *service.HealthStatus, w healthy.Watcher)
	healthHandle   func(event *service.HealthStatus, w healthy.Watcher)
	healthyManager healthy.Manager
}

func (s *statusController) Run() {
	logrus.Info("run status controller")
	s.healthyManager.EnableWatcher(s.service.Name, s.watcher.GetID())
	defer s.watcher.Close()
	defer s.healthyManager.DisableWatcher(s.service.Name, s.watcher.GetID())
	for {
		select {
		case event := <-s.watcher.Watch():
			switch event.Status {
			case service.Stat_healthy:
				if s.healthHandle != nil {
					s.healthHandle(event, s.watcher)
				}
				logrus.Debugf("service %s status is [%s]", event.Name, event.Status)
			case service.Stat_unhealthy:
				if s.service.ServiceHealth != nil {
					if event.ErrorNumber > s.service.ServiceHealth.MaxErrorsNum {
						logrus.Warningf("service %s status is [%s] more than %d times and restart it.", event.Name, event.Status, s.service.ServiceHealth.MaxErrorsNum)
						s.unhealthHandle(event, s.watcher)
					}
				}
			case service.Stat_death:
				logrus.Warningf("service %s status is [%s] will restart it.", event.Name, event.Status)
				s.unhealthHandle(event, s.watcher)
			}
		case <-s.ctx.Done():
			return
		}
	}
}
func (s *statusController) Stop() {
	s.cancel()
}

// StopSyncService -
func (m *ManagerService) StopSyncService() {
	if m.syncCtx != nil {
		m.syncCancel()
	}
}

//WaitStart waiting for service to be healthy
func (m *ManagerService) WaitStart(name string, duration time.Duration) bool {
	max := time.Now().Add(duration)
	t := time.Tick(time.Second * 3)
	for {
		if time.Now().After(max) {
			return false
		}
		status, err := m.healthyManager.GetCurrentServiceHealthy(name)
		if err != nil {
			logrus.Errorf("Can not get %s service current status: %v", name, err)
			<-t
			continue
		}
		logrus.Debugf("Check service %s current status: %s", name, status.Status)
		if status.Status == service.Stat_healthy {
			return true
		}
		<-t
	}
}

// ReLoadServices -
/*
1. reload services info from local file system
2. regenerate systemd config file and restart with config changes
3. start all newly added services
*/
func (m *ManagerService) ReLoadServices() error {
	logrus.Info("start reload service configs")
	services := service.LoadServicesFromLocal(m.conf.ServiceListFile)
	var controllerServices []*service.Service
	var restartCount int
	for _, ne := range services {
		if ne.OnlyHealthCheck {
			continue
		}
		if !ne.Disable {
			controllerServices = append(controllerServices, ne)
		}
		exists := false
		for _, old := range m.services {
			if ne.Name == old.Name {
				if ne.Disable {
					m.ctr.StopService(ne.Name)
					m.ctr.DisableService(ne.Name)
					restartCount++
				}
				logrus.Infof("Recreate service [%s]", ne.Name)
				update, err := m.ctr.WriteConfig(ne)
				if err == nil {
					m.ctr.EnableService(ne.Name)
					if update {
						m.ctr.RestartService(ne)
					} else {
						m.ctr.StartService(ne.Name)
					}
					restartCount++
				}
				exists = true
				break
			}
		}
		if !exists {
			logrus.Infof("Create service [%s]", ne.Name)
			update, err := m.ctr.WriteConfig(ne)
			if err == nil {
				m.ctr.EnableService(ne.Name)
				if update {
					m.ctr.RestartService(ne)
				} else {
					m.ctr.StartService(ne.Name)
				}
				restartCount++
			}
		}
	}
	m.allservice = services
	m.services = controllerServices
	m.healthyManager.AddServicesAndUpdate(m.allservice)
	m.SyncServiceStatusController()
	logrus.Infof("load service config success, start or stop %d service and total %d service", restartCount, len(services))
	return nil
}

//StartService
func (m *ManagerService) StartService(serviceName string) error {
	for _, service := range m.services {
		if service.Name == serviceName {
			if !service.Disable {
				return fmt.Errorf("service %s is running", serviceName)
			}
			return m.ctr.StartService(serviceName)
		}
	}
	return nil
}

//StopService
func (m *ManagerService) StopService(serviceName string) error {
	for i, service := range m.services {
		if service.Name == serviceName {
			if service.Disable {
				return fmt.Errorf("service %s is stoped", serviceName)
			}
			(m.services)[i].Disable = true
			m.lock.Lock()
			defer m.lock.Unlock()
			if controller, ok := m.autoStatusController[serviceName]; ok {
				controller.Stop()
			}
			return m.ctr.StopService(serviceName)
		}
	}
	return nil
}

//WriteServices
func (m *ManagerService) WriteServices() error {
	for _, s := range m.services {
		if s.OnlyHealthCheck {
			continue
		}
		if s.Name == "docker" {
			continue
		}
		_, err := m.ctr.WriteConfig(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func isExistEndpoint(etcdEndPoints []string, end string) bool {
	for _, v := range etcdEndPoints {
		if v == end {
			return true
		}
	}
	return false
}

func rmEndpointFrom(etcdEndPoints []string, end string) []string {
	endPoints := make([]string, 0, 5)
	for _, v := range etcdEndPoints {
		if v != end {
			endPoints = append(endPoints, v)
		}
	}
	return endPoints
}

func toEndpoint(reg *service.Endpoint, ip string) string {
	if reg.Protocol == "" {
		return fmt.Sprintf("%s:%s", ip, reg.Port)
	}
	return fmt.Sprintf("%s://%s:%s", reg.Protocol, ip, reg.Port)
}

//InjectConfig
func (m *ManagerService) InjectConfig(content string) string {
	for _, parantheses := range ArgsReg.FindAllString(content, -1) {
		logrus.Debugf("discover inject args template %s", parantheses)
		group := ArgsReg.FindStringSubmatch(parantheses)
		if group == nil || len(group) < 2 {
			logrus.Warnf("Not found group for %s", parantheses)
			continue
		}
		line := ""
		if group[1] == "NODE_UUID" {
			line = m.node.ID
		} else {
			endpoints := m.cluster.GetEndpoints(group[1])
			if len(endpoints) < 1 {
				logrus.Warnf("Failed to inject endpoints of key %s", group[1])
				continue
			}
			sep := ","
			if len(group) >= 3 && group[2] != "" {
				sep = group[2]
			}
			for _, end := range endpoints {
				if line == "" {
					line = end
				} else {
					line += sep
					line += end
				}
			}
		}
		content = strings.Replace(content, group[0], line, 1)
		logrus.Debugf("inject args into service %s => %v", group[1], line)
	}
	return content
}

// ListServiceImages -
func (m *ManagerService) ListServiceImages() []string {
	var images []string
	for _, svc := range m.services {
		if svc.Start == "" || svc.OnlyHealthCheck {
			continue
		}

		par := parser.CreateDockerRunOrImageParse("", "", svc.Start, nil, event.GetTestLogger())
		par.ParseDockerun(svc.Start)
		logrus.Debugf("detect image: %s", par.GetImage().String())
		if par.GetImage().String() == "" {
			continue
		}
		images = append(images, par.GetImage().String())
	}

	return images
}

//NewManagerService
func NewManagerService(conf *option.Conf, healthyManager healthy.Manager, cluster client.ClusterClient) *ManagerService {
	ctx, cancel := context.WithCancel(context.Background())
	manager := &ManagerService{
		ctx:            ctx,
		cancel:         cancel,
		conf:           conf,
		cluster:        cluster,
		healthyManager: healthyManager,
		etcdcli:        conf.EtcdCli,
	}
	manager.ctr = NewController(conf, manager)
	return manager
}
