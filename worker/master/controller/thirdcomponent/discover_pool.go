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
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DiscoverPool struct {
	ctx            context.Context
	lock           sync.Mutex
	discoverWorker map[string]*Worker
	updateChan     chan *v1alpha1.ThirdComponent
	reconciler     *Reconciler
}

func NewDiscoverPool(ctx context.Context, reconciler *Reconciler) *DiscoverPool {
	dp := &DiscoverPool{
		ctx:            ctx,
		discoverWorker: make(map[string]*Worker),
		updateChan:     make(chan *v1alpha1.ThirdComponent, 1024),
		reconciler:     reconciler,
	}
	go dp.Start()
	return dp
}

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
						logrus.Errorf("update component status failure", err.Error())
					}
					logrus.Infof("update component %s status success by discover pool", name)
				} else {
					logrus.Debugf("component %s status endpoints not change", name)
				}
			}()
		}
	}
}

type Worker struct {
	discover   Discover
	cancel     context.CancelFunc
	ctx        context.Context
	updateChan chan *v1alpha1.ThirdComponent
	stoped     bool
}

func (w *Worker) Start() {
	defer func() {
		logrus.Infof("discover endpoint list worker %s/%s stoed", w.discover.GetComponent().Namespace, w.discover.GetComponent().Name)
		w.stoped = true
	}()
	w.stoped = false
	logrus.Infof("discover endpoint list worker %s/%s  started", w.discover.GetComponent().Namespace, w.discover.GetComponent().Name)
	for {
		w.discover.Discover(w.ctx, w.updateChan)
		select {
		case <-w.ctx.Done():
			return
		default:
		}
	}
}

func (w *Worker) UpdateDiscover(discover Discover) {
	w.discover = discover
}

func (w *Worker) Stop() {
	w.cancel()
}

func (w *Worker) IsStop() bool {
	return w.stoped
}

func (d *DiscoverPool) newWorker(dis Discover) *Worker {
	ctx, cancel := context.WithCancel(d.ctx)
	return &Worker{
		ctx:        ctx,
		discover:   dis,
		cancel:     cancel,
		updateChan: d.updateChan,
	}
}

func (d *DiscoverPool) AddDiscover(dis Discover) {
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
	go worker.Start()
	d.discoverWorker[key] = worker
}

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
