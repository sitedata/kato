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

package prober

import (
	"context"
	"fmt"
	"github.com/gridworkz/kato/util/prober/types/v1"
	"testing"
)

func TestProbeManager_Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	m := NewProber(ctx, cancel)

	serviceList := make([]*v1.Service, 0, 10)

	h := &v1.Service{
		Name: "etcd",
		ServiceHealth: &v1.Health{
			Name:         "etcd",
			Model:        "tcp",
			Address:      "192.168.1.107:23790",
			TimeInterval: 3,
		},
	}
	serviceList = append(serviceList, h)
	m.SetServices(serviceList)
	watcher := m.WatchServiceHealthy("etcd")
	m.EnableWatcher(watcher.GetServiceName(), watcher.GetID())

	m.Start()

	for {
		v := <-watcher.Watch()
		if v != nil {
			fmt.Println("----", v.Name, v.Status, v.Info, v.ErrorNumber, v.ErrorNumber)
		} else {
			t.Log("nil nil nil")
		}
	}
}
