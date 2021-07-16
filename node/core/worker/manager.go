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

package worker

import (
	"github.com/gridworkz/kato/node/api/model"
	"sync"

	"golang.org/x/net/context"
)

//Worker
type Worker interface {
	Start()
	Stop() error
	Result()
}

//Manager
type Manager struct {
	workers map[string]Worker
	lock    sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	closed  chan struct{}
}

//NewManager
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := Manager{
		ctx:    ctx,
		cancel: cancel,
	}
	return &m
}

//Start
func (m *Manager) Start() error {
	return nil
}

//Stop
func (m *Manager) Stop() error {
	return nil
}

//AddWorker
func (m *Manager) AddWorker(worker Worker) error {
	return nil
}

//RemoveWorker
func (m *Manager) RemoveWorker(worker Worker) error {
	return nil
}

//NewTaskWorker
func (m *Manager) NewTaskWorker(task *model.Task) Worker {
	return &taskWorker{task}
}

//NewTaskGroupWorker
func (m *Manager) NewTaskGroupWorker(taskgroup *model.TaskGroup) Worker {
	return &taskGroupWorker{taskgroup}
}
