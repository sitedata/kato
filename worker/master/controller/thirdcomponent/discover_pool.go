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

package thirdcomponent

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	dis "github.com/gridworkz/kato/worker/master/controller/thirdcomponent/discover"
	"github.com/gridworkz/kato/worker/master/controller/thirdcomponent/prober"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DiscoverPool -
type DiscoverPool struct {
	ctx            context.Context
	lock           sync.Mutex
	discoverWorker map[string]*Worker
	updateChan     chan *v1alpha1.ThirdComponent
	reconciler     *Reconciler

	recorder record.EventRecorder
}

// NewDiscoverPool -
func NewDiscoverPool(ctx context.Context,
	reconciler *Reconciler,
	recorder record.EventRecorder) *DiscoverPool {
	dp := &DiscoverPool{
		ctx:            ctx,
		discoverWorker: make(map[string]*Worker),
		updateChan:     make(chan *v1alpha1.ThirdComponent, 1024),
		reconciler:     reconciler,
		recorder:       recorder,
	}
	go dp.Start()
	return dp
}

// GetSize -
func (d *DiscoverPool) GetSize() float64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	return float64(len(d.discoverWorker))
}

// Start -
func (d *DiscoverPool) Start() {
	logrus.Infof("third component discover pool started")
	for {
		select {
		case <-d.ctx.Done():
			logrus.Infof("third component discover pool stoped")
			return
		case component := <-d.updateChan:
			func() {
				ctx, cancel := context.WithTimeout(d.ctx, time.Second*10)
				defer cancel()
				var old v1alpha1.ThirdComponent
				name := client.ObjectKey{Name: component.Name, Namespace: component.Namespace}
				d.reconciler.Client.Get(ctx, name, &old)
				if !reflect.DeepEqual(component.Status.Endpoints, old.Status.Endpoints) {
					if err := d.reconciler.updateStatus(ctx, component); err != nil {
						if apierrors.IsNotFound(err) {
							d.RemoveDiscover(component)
							return
						}
						logrus.Errorf("update component status failure: %s", err.Error())
					}
					logrus.Infof("update component %s status success by discover pool", name)
				}
			}()
		}
	}
}

func (d *DiscoverPool) newWorker(dis dis.Discover) *Worker {
	ctx, cancel := context.WithCancel(d.ctx)

	worker := &Worker{
		ctx:        ctx,
		discover:   dis,
		cancel:     cancel,
		updateChan: d.updateChan,
	}

	component := dis.GetComponent()
	if component.Spec.IsStaticEndpoints() {
		proberManager := prober.NewManager(d.recorder)
		dis.SetProberManager(proberManager)
		worker.proberManager = proberManager
	}

	return worker
}

// AddDiscover -
func (d *DiscoverPool) AddDiscover(dis dis.Discover) {
	d.lock.Lock()
	defer d.lock.Unlock()
	component := dis.GetComponent()
	if component == nil {
		return
	}
	key := component.Namespace + component.Name
	olddis, exist := d.discoverWorker[key]
	if exist {
		olddis.UpdateDiscover(dis)
		if olddis.IsStop() {
			go olddis.Start()
		}
		return
	}
	worker := d.newWorker(dis)
	if component.Spec.IsStaticEndpoints() {
		worker.proberManager.AddThirdComponent(dis.GetComponent())
	}
	go worker.Start()
	d.discoverWorker[key] = worker
}

// RemoveDiscover -
func (d *DiscoverPool) RemoveDiscover(component *v1alpha1.ThirdComponent) {
	d.lock.Lock()
	defer d.lock.Unlock()
	key := component.Namespace + component.Name
	olddis, exist := d.discoverWorker[key]
	if exist {
		olddis.Stop()
		delete(d.discoverWorker, key)
	}
}

// RemoveDiscoverByName -
func (d *DiscoverPool) RemoveDiscoverByName(req types.NamespacedName) {
	d.lock.Lock()
	defer d.lock.Unlock()
	key := req.Namespace + req.Name
	olddis, exist := d.discoverWorker[key]
	if exist {
		olddis.Stop()
		delete(d.discoverWorker, key)
	}
}
