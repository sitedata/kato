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

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
	"github.com/gridworkz/kato/worker/appm/conversion"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
)

type restartController struct {
	stopChan     chan struct{}
	controllerID string
	appService   []v1.AppService
	manager      *Manager
	ctx          context.Context
}

func (s *restartController) Begin() {
	var wait sync.WaitGroup
	for _, service := range s.appService {
		wait.Add(1)
		go func(service v1.AppService) {
			defer wait.Done()
			service.Logger.Info("App runtime begin restart app service "+service.ServiceAlias, event.GetLoggerOption("starting"))
			if err := s.restartOne(service); err != nil {
				logrus.Errorf("restart service %s failure %s", service.ServiceAlias, err.Error())
			} else {
				service.Logger.Info(fmt.Sprintf("restart service %s success", service.ServiceAlias), event.GetLastLoggerOption())
			}
		}(service)
	}
	wait.Wait()
	s.manager.callback(s.controllerID, nil)
}
func (s *restartController) restartOne(app v1.AppService) error {
	//Restart the control set timeout interval is 5m
	stopController := &stopController{
		manager:      s.manager,
		waiting:      time.Minute * 5,
		ctx:          s.ctx,
		controllerID: s.controllerID,
	}
	if err := stopController.stopOne(app); err != nil {
		if err != ErrWaitTimeOut {
			app.Logger.Error(util.Translation("(restart)stop service error"), event.GetCallbackLoggerOption())
			return err
		}
		//waiting app closed,max wait 40 second
		var waiting = 20
		for waiting > 0 {
			storeAppService := s.manager.store.GetAppService(app.ServiceID)
			if storeAppService == nil || storeAppService.IsClosed() {
				break
			}
			waiting--
			time.Sleep(time.Second * 2)
		}
	}
	startController := startController{
		manager:      s.manager,
		ctx:          s.ctx,
		controllerID: s.controllerID,
	}
	newAppService, err := conversion.InitAppService(db.GetManager(), app.ServiceID, app.ExtensionSet)
	if err != nil {
		logrus.Errorf("Application model init create failure:%s", err.Error())
		app.Logger.Error(util.Translation("(restart)Application model init create failure"), event.GetCallbackLoggerOption())
		return fmt.Errorf("application model init create failure,%s", err.Error())
	}
	newAppService.Logger = app.Logger
	//regist new app service
	s.manager.store.RegistAppService(newAppService)
	if err := startController.startOne(*newAppService); err != nil {
		if err != ErrWaitTimeOut {
			app.Logger.Error(util.Translation("start service error"), event.GetCallbackLoggerOption())
			return err
		}
	}
	return nil
}
func (s *restartController) Stop() error {
	close(s.stopChan)
	return nil
}
