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
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/nodem/service"
)

//NewController At the stage you want to load the configurations of all kato components
func NewController(conf *option.Conf, manager *ManagerService) Controller {
	return &testController{}
}

type testController struct {
}

func (t *testController) InitStart(services []*service.Service) error {
	return nil
}
func (t *testController) StartService(name string) error {
	return nil
}
func (t *testController) StopService(name string) error {
	return nil
}
func (t *testController) StartList(list []*service.Service) error {
	return nil
}
func (t *testController) StopList(list []*service.Service) error {
	return nil
}
func (t *testController) RestartService(s *service.Service) error {
	return nil
}
func (t *testController) WriteConfig(s *service.Service) (bool, error) {
	return false, nil
}
func (t *testController) RemoveConfig(name string) error {
	return nil
}
func (t *testController) EnableService(name string) error {
	return nil
}
func (t *testController) DisableService(name string) error {
	return nil
}
func (t *testController) CheckBeforeStart() bool {
	return false
}
