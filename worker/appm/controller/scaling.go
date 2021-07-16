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
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
	types "k8s.io/apimachinery/pkg/types"
)

type scalingController struct {
	controllerID string
	appService   []v1.AppService
	manager      *Manager
	stopChan     chan struct{}
}

//Begin start handle service scaling
func (s *scalingController) Begin() {
	var wait sync.WaitGroup
	for _, service := range s.appService {
		go func(service v1.AppService) {
			wait.Add(1)
			defer wait.Done()
			service.Logger.Info("App runtime begin horizontal scaling app service "+service.ServiceAlias, event.GetLoggerOption("starting"))
			if err := s.scalingOne(service); err != nil {
				if err != ErrWaitTimeOut {
					service.Logger.Error(util.Translation("horizontal scaling service error"), event.GetCallbackLoggerOption())
					logrus.Errorf("horizontal scaling service %s failure %s", service.ServiceAlias, err.Error())
				} else {
					service.Logger.Error(util.Translation("horizontal scaling service timeout"), event.GetTimeoutLoggerOption())
				}
			} else {
				service.Logger.Info(fmt.Sprintf("horizontal scaling service %s success", service.ServiceAlias), event.GetLastLoggerOption())
			}
		}(service)
	}
	wait.Wait()
	s.manager.callback(s.controllerID, nil)
}

//Replicas fetch replicas to n
func Replicas(n int) []byte {
	return []byte(fmt.Sprintf(`{"spec":{"replicas":%d}}`, n))
}

func (s *scalingController) scalingOne(service v1.AppService) error {
	if statefulset := service.GetStatefulSet(); statefulset != nil {
		_, err := s.manager.client.AppsV1().StatefulSets(statefulset.Namespace).Patch(statefulset.Name, types.StrategicMergePatchType, Replicas(int(service.Replicas)))
		if err != nil {
			logrus.Error("patch statefulset info error.", err.Error())
			return err
		}
	}
	if deployment := service.GetDeployment(); deployment != nil {
		_, err := s.manager.client.AppsV1().Deployments(deployment.Namespace).Patch(deployment.Name, types.StrategicMergePatchType, Replicas(int(service.Replicas)))
		if err != nil {
			logrus.Error("patch deployment info error.", err.Error())
			return err
		}
	}
	return s.WaitingReady(service)
}

//WaitingReady wait app start or upgrade ready
func (s *scalingController) WaitingReady(app v1.AppService) error {
	storeAppService := s.manager.store.GetAppService(app.ServiceID)
	var initTime int32
	if podt := app.GetPodTemplate(); podt != nil {
		if probe := podt.Spec.Containers[0].ReadinessProbe; probe != nil {
			initTime = probe.InitialDelaySeconds
		}
	}
	//at least waiting time is 40 second
	initTime += 40
	waitingReplicas := math.Abs(float64(storeAppService.Replicas) - float64(storeAppService.GetReadyReplicas()))
	timeout := time.Duration(initTime * int32(waitingReplicas))
	if timeout.Seconds() < 40 {
		timeout = time.Duration(time.Second * 40)
	}
	if err := WaitReady(s.manager.store, storeAppService, timeout, app.Logger, s.stopChan); err != nil {
		return err
	}
	return nil
}
func (s *scalingController) Stop() error {
	close(s.stopChan)

	return nil
}
