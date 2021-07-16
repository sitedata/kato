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

package helmapp

import (
	"context"
	"github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	"github.com/gridworkz/kato/pkg/generated/listers/kato/v1alpha1"
	"github.com/sirupsen/logrus"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Controller -
type Controller struct {
	storer      Storer
	stopCh      chan struct{}
	controlLoop *ControlLoop
	finalizer   *Finalizer
}

// NewController creates a new helm app controller.
func NewController(ctx context.Context,
	stopCh chan struct{},
	kubeClient clientset.Interface,
	clientset versioned.Interface,
	informer cache.SharedIndexInformer,
	lister v1alpha1.HelmAppLister,
	repoFile, repoCache, chartCache string) *Controller {
	workQueue := workqueue.New()
	finalizerQueue := workqueue.New()
	storer := NewStorer(informer, lister, workQueue, finalizerQueue)

	controlLoop := NewControlLoop(ctx, kubeClient, clientset, storer, workQueue, repoFile, repoCache, chartCache)
	finalizer := NewFinalizer(ctx, kubeClient, clientset, finalizerQueue, repoFile, repoCache, chartCache)

	return &Controller{
		storer:      storer,
		stopCh:      stopCh,
		controlLoop: controlLoop,
		finalizer:   finalizer,
	}
}

// Start starts the controller.
func (c *Controller) Start() {
	logrus.Info("start helm app controller")
	c.storer.Run(c.stopCh)
	go c.controlLoop.Run()
	c.finalizer.Run()
}

// Stop stops the controller.
func (c *Controller) Stop() {
	c.controlLoop.Stop()
	c.finalizer.Stop()
}
