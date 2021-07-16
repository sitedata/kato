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

//+build windows

package controller

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/nodem/service"
	"github.com/gridworkz/kato/util/windows"
)

//NewController At the stage you want to load the configurations of all kato components
func NewController(conf *option.Conf, manager *ManagerService) Controller {
	logrus.Infof("Create windows service controller")
	return &windowsServiceController{
		conf:    conf,
		manager: manager,
	}
}

type windowsServiceController struct {
	conf    *option.Conf
	manager *ManagerService
}

func (w *windowsServiceController) InitStart(services []*service.Service) error {
	for _, s := range services {
		if s.IsInitStart && !s.Disable && !s.OnlyHealthCheck {
			if err := w.writeConfig(s, false); err != nil {
				return err
			}
			if err := w.StartService(s.Name); err != nil {
				return fmt.Errorf("start windows service %s failure %s", s.Name, err.Error())
			}
		}
	}
	return nil
}

func (w *windowsServiceController) StartService(name string) error {
	if err := windows.StartService(name); err != nil {
		logrus.Errorf("windows service controller start service %s failure %s", name, err.Error())
		return err
	}
	logrus.Infof("windows service controller start service %s success", name)
	return nil
}
func (w *windowsServiceController) StopService(name string) error {
	if err := windows.StopService(name); err != nil && !strings.Contains(err.Error(), "service has not been started") {
		logrus.Errorf("windows service controller stop service %s failure %s", name, err.Error())
		return err
	}
	logrus.Infof("windows service controller stop service %s success", name)
	return nil
}
func (w *windowsServiceController) StartList(list []*service.Service) error {
	for _, s := range list {
		w.StartService(s.Name)
	}
	return nil
}
func (w *windowsServiceController) StopList(list []*service.Service) error {
	for _, s := range list {
		w.StopService(s.Name)
	}
	return nil
}
func (w *windowsServiceController) RestartService(s *service.Service) error {
	if err := windows.RestartService(s.Name); err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			if _, err := w.WriteConfig(s); err != nil {
				return fmt.Errorf("ReWrite service config failure %s", err.Error())
			}
		}
		logrus.Errorf("windows service controller restart service %s failure %s", s.Name, err.Error())
		return err
	}
	logrus.Infof("windows service controller restart service %s success", s.Name)
	return nil
}
func (w *windowsServiceController) WriteConfig(s *service.Service) (bool, error) {
	return true, w.writeConfig(s, true)
}
func (w *windowsServiceController) writeConfig(s *service.Service, parseAndCoverOld bool) error {
	cmdstr := s.Start
	if parseAndCoverOld {
		cmdstr = w.manager.InjectConfig(s.Start)
	}
	cmds := strings.Split(cmdstr, " ")
	logrus.Debugf("write service %s config args %s", s.Name, cmds)
	if err := windows.RegisterService(s.Name, cmds[0], "Kato "+s.Name, s.Requires, cmds); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			if parseAndCoverOld {
				w.RemoveConfig(s.Name)
				err = windows.RegisterService(s.Name, cmds[0], "Kato "+s.Name, s.Requires, cmds)
			} else {
				logrus.Infof("windows service controller register service %s success(exist)", s.Name)
				return nil
			}
		}
		if err != nil {
			logrus.Errorf("windows service controller register service %s failure %s", s.Name, err.Error())
			return err
		}
	}
	logrus.Infof("windows service controller register service %s success", s.Name)
	return nil
}
func (w *windowsServiceController) RemoveConfig(name string) error {
	return windows.UnRegisterService(name)
}
func (w *windowsServiceController) EnableService(name string) error {
	return nil
}
func (w *windowsServiceController) DisableService(name string) error {
	//return windows.UnRegisterService(name)
	return nil
}
func (w *windowsServiceController) CheckBeforeStart() bool {
	return true
}
