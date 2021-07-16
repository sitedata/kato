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

	"github.com/gridworkz/kato/util"
	"github.com/gridworkz/kato/worker/appm/store"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"k8s.io/client-go/kubernetes"
)

//Controller service operating controller interface
type Controller interface {
	Begin()
	Stop() error
}

//TypeController controller type
type TypeController string

//TypeStartController start service type
var TypeStartController TypeController = "start"

//TypeStopController start service type
var TypeStopController TypeController = "stop"

//TypeRestartController restart service type
var TypeRestartController TypeController = "restart"

//TypeUpgradeController start service type
var TypeUpgradeController TypeController = "upgrade"

//TypeScalingController start service type
var TypeScalingController TypeController = "scaling"

// TypeApplyRuleController -
var TypeApplyRuleController TypeController = "apply_rule"

// TypeApplyConfigController -
var TypeApplyConfigController TypeController = "apply_config"

// TypeControllerRefreshHPA -
var TypeControllerRefreshHPA TypeController = "refreshhpa"

//Manager controller manager
type Manager struct {
	ctx         context.Context
	cancel      context.CancelFunc
	client      kubernetes.Interface
	controllers map[string]Controller
	store       store.Storer
	lock        sync.Mutex
}

//NewManager new manager
func NewManager(store store.Storer, client kubernetes.Interface) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:         ctx,
		cancel:      cancel,
		client:      client,
		controllers: make(map[string]Controller),
		store:       store,
	}
}

//Stop stop all controller
func (m *Manager) Stop() error {
	m.cancel()
	return nil
}

//GetControllerSize get running controller number
func (m *Manager) GetControllerSize() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	return len(m.controllers)
}

//StartController create and start service controller
func (m *Manager) StartController(controllerType TypeController, apps ...v1.AppService) error {
	var controller Controller
	controllerID := util.NewUUID()
	switch controllerType {
	case TypeStartController:
		controller = &startController{
			controllerID: controllerID,
			appService:   apps,
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	case TypeStopController:
		controller = &stopController{
			controllerID: controllerID,
			appService:   apps,
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	case TypeScalingController:
		controller = &scalingController{
			controllerID: controllerID,
			appService:   apps,
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	case TypeUpgradeController:
		controller = &upgradeController{
			controllerID: controllerID,
			appService:   apps,
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	case TypeRestartController:
		controller = &restartController{
			controllerID: controllerID,
			appService:   apps,
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	case TypeApplyRuleController:
		controller = &applyRuleController{
			controllerID: controllerID,
			appService:   apps,
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	case TypeApplyConfigController:
		controller = &applyConfigController{
			controllerID: controllerID,
			appService:   apps[0],
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	case TypeControllerRefreshHPA:
		controller = &refreshXPAController{
			controllerID: controllerID,
			appService:   apps,
			manager:      m,
			stopChan:     make(chan struct{}),
		}
	default:
		return fmt.Errorf("No support controller")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.controllers[controllerID] = controller
	go controller.Begin()
	return nil
}

func (m *Manager) callback(controllerID string, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.controllers, controllerID)
}

type sequencelist []sequence
type sequence []*v1.AppService

func (s *sequencelist) Contains(id string) bool {
	for _, l := range *s {
		for _, l2 := range l {
			if l2.ServiceID == id {
				return true
			}
		}
	}
	return false
}
func (s *sequencelist) Add(ids []*v1.AppService) {
	*s = append(*s, ids)
}

func foundsequence(source map[string]*v1.AppService, sl *sequencelist) {
	if len(source) == 0 {
		return
	}
	var deleteKey []string
source:
	for _, s := range source {
		for _, d := range s.Dependces {
			if !sl.Contains(d) {
				continue source
			}
		}
		deleteKey = append(deleteKey, s.ServiceID)
	}
	var list []*v1.AppService
	for _, d := range deleteKey {
		list = append(list, source[d])
		delete(source, d)
	}
	sl.Add(list)
	foundsequence(source, sl)
}

func decisionSequence(appService []*v1.AppService) sequencelist {
	var sourceIDs = make(map[string]*v1.AppService, len(appService))
	for _, a := range appService {
		sourceIDs[a.ServiceID] = a
	}
	var sl sequencelist
	foundsequence(sourceIDs, &sl)
	return sl
}
