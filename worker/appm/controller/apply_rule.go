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
	"sync"

	"github.com/gridworkz/kato/worker/appm/f"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
)

type applyRuleController struct {
	controllerID string
	appService   []v1.AppService
	manager      *Manager
	stopChan     chan struct{}
	ctx          context.Context
}

// Begin begins applying rule
func (a *applyRuleController) Begin() {
	var wait sync.WaitGroup
	for _, service := range a.appService {
		wait.Add(1)
		go func(service v1.AppService) {
			defer wait.Done()
			if err := f.ApplyOne(a.ctx, a.manager.apply, a.manager.client, &service); err != nil {
				logrus.Errorf("apply rules for service %s failure: %s", service.ServiceAlias, err.Error())
			}
		}(service)
	}
	wait.Wait()
	a.manager.callback(a.controllerID, nil)
}

func (a *applyRuleController) Stop() error {
	close(a.stopChan)
	return nil
}
