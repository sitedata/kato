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
	"time"

	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/worker/appm/store"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

//ErrWaitTimeOut wait time out
var ErrWaitTimeOut = fmt.Errorf("Wait time out")

//ErrWaitCancel wait cancel
var ErrWaitCancel = fmt.Errorf("Wait cancel")

//WaitReady wait ready
func WaitReady(store store.Storer, a *v1.AppService, timeout time.Duration, logger event.Logger, cancel chan struct{}) error {
	if timeout < 40 {
		timeout = time.Second * 40
	}
	logger.Info(fmt.Sprintf("waiting app ready timeout %ds", int(timeout.Seconds())), map[string]string{"step": "appruntime", "status": "running"})
	logrus.Debugf("waiting app ready timeout %ds", int(timeout.Seconds()))
	ticker := time.NewTicker(timeout / 10)
	timer := time.NewTimer(timeout)
	defer ticker.Stop()
	var i int
	for {
		if i > 2 {
			a = store.UpdateGetAppService(a.ServiceID)
		}
		if a.Ready() {
			return nil
		}
		printLogger(a, logger)
		select {
		case <-cancel:
			return ErrWaitCancel
		case <-timer.C:
			//if service status is waitting, the event is not timeout
			// if a.IsWaitting() {
			// 	timer.Reset(timeout)
			// }
			return ErrWaitTimeOut
		case <-ticker.C:
		}
		i++
	}
}

//WaitStop wait service stop complete
func WaitStop(store store.Storer, a *v1.AppService, timeout time.Duration, logger event.Logger, cancel chan struct{}) error {
	if a == nil {
		return nil
	}
	if timeout < 40 {
		timeout = time.Second * 40
	}
	logger.Info(fmt.Sprintf("waiting app closed timeout %ds", int(timeout.Seconds())), map[string]string{"step": "appruntime", "status": "running"})
	logrus.Debugf("waiting app ready timeout %ds", int(timeout.Seconds()))
	ticker := time.NewTicker(timeout / 10)
	timer := time.NewTimer(timeout)
	defer ticker.Stop()
	var i int
	for {
		i++
		if i > 2 {
			a = store.UpdateGetAppService(a.ServiceID)
		}
		if a.IsClosed() {
			return nil
		}
		printLogger(a, logger)
		select {
		case <-cancel:
			return ErrWaitCancel
		case <-timer.C:
			return ErrWaitTimeOut
		case <-ticker.C:
		}
	}
}

//WaitUpgradeReady wait upgrade success
func WaitUpgradeReady(store store.Storer, a *v1.AppService, timeout time.Duration, logger event.Logger, cancel chan struct{}) error {
	if a == nil {
		return nil
	}
	if timeout < 40 {
		timeout = time.Second * 40
	}
	logger.Info(fmt.Sprintf("waiting app upgrade complete timeout %ds", int(timeout.Seconds())), map[string]string{"step": "appruntime", "status": "running"})
	logrus.Debugf("waiting app upgrade complete timeout %ds", int(timeout.Seconds()))
	ticker := time.NewTicker(timeout / 10)
	timer := time.NewTimer(timeout)
	defer ticker.Stop()
	for {
		if a.UpgradeComlete() {
			return nil
		}
		printLogger(a, logger)
		select {
		case <-cancel:
			return ErrWaitCancel
		case <-timer.C:
			return ErrWaitTimeOut
		case <-ticker.C:
		}
	}
}
func printLogger(a *v1.AppService, logger event.Logger) {
	var ready int32
	if a.GetStatefulSet() != nil {
		ready = a.GetStatefulSet().Status.ReadyReplicas
	}
	if a.GetDeployment() != nil {
		ready = a.GetDeployment().Status.ReadyReplicas
	}
	logger.Info(fmt.Sprintf("current instance(count:%d ready:%d notready:%d)", len(a.GetPods(false)), ready, int32(len(a.GetPods(false)))-ready), map[string]string{"step": "appruntime", "status": "running"})
	pods := a.GetPods(false)
	for _, pod := range pods {
		for _, con := range pod.Status.Conditions {
			if con.Status == corev1.ConditionFalse {
				logger.Debug(fmt.Sprintf("instance %s %s status is %s: %s", pod.Name, con.Type, con.Status, con.Message), map[string]string{"step": "appruntime", "status": "running"})
			}
		}
	}
}
