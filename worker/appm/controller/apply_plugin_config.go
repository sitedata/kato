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
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type applyConfigController struct {
	controllerID string
	appService   v1.AppService
	manager      *Manager
	stopChan     chan struct{}
}

// Begin begins applying rule
func (a *applyConfigController) Begin() {
	nowApp := a.manager.store.GetAppService(a.appService.ServiceID)
	nowConfigMaps := nowApp.GetConfigMaps()
	newConfigMaps := a.appService.GetConfigMaps()
	var nowConfigMapMaps = make(map[string]*corev1.ConfigMap, len(nowConfigMaps))
	for i, now := range nowConfigMaps {
		nowConfigMapMaps[now.Name] = nowConfigMaps[i]
	}
	for _, new := range newConfigMaps {
		if nowConfig, ok := nowConfigMapMaps[new.Name]; ok {
			new.UID = nowConfig.UID
			newc, err := a.manager.client.CoreV1().ConfigMaps(nowApp.TenantID).Update(new)
			if err != nil {
				logrus.Errorf("update config map failure %s", err.Error())
			}
			nowApp.SetConfigMap(newc)
			nowConfigMapMaps[new.Name] = nil
			logrus.Debugf("update configmap %s for service %s", new.Name, a.appService.ServiceID)
		} else {
			newc, err := a.manager.client.CoreV1().ConfigMaps(nowApp.TenantID).Create(new)
			if err != nil {
				logrus.Errorf("update config map failure %s", err.Error())
			}
			nowApp.SetConfigMap(newc)
			logrus.Debugf("create configmap %s for service %s", new.Name, a.appService.ServiceID)
		}
	}
	for name, handle := range nowConfigMapMaps {
		if handle != nil {
			if err := a.manager.client.CoreV1().ConfigMaps(nowApp.TenantID).Delete(name, &metav1.DeleteOptions{}); err != nil {
				logrus.Errorf("delete config map failure %s", err.Error())
			}
			logrus.Debugf("delete configmap %s for service %s", name, a.appService.ServiceID)
		}
	}
	a.manager.callback(a.controllerID, nil)
}

func (a *applyConfigController) Stop() error {
	close(a.stopChan)
	return nil
}
