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

	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	"github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	"github.com/sirupsen/logrus"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/workqueue"
)

// Finalizer does some cleanup work when helmApp is deleted
type Finalizer struct {
	ctx        context.Context
	log        *logrus.Entry
	kubeClient clientset.Interface
	clientset  versioned.Interface
	queue      workqueue.Interface
	repoFile   string
	repoCache  string
	chartCache string
}

// NewFinalizer creates a new finalizer.
func NewFinalizer(ctx context.Context,
	kubeClient clientset.Interface,
	clientset versioned.Interface,
	workQueue workqueue.Interface,
	repoFile string,
	repoCache string,
	chartCache string,
) *Finalizer {
	return &Finalizer{
		ctx:        ctx,
		log:        logrus.WithField("WHO", "Helm App Finalizer"),
		kubeClient: kubeClient,
		clientset:  clientset,
		queue:      workQueue,
		repoFile:   repoFile,
		repoCache:  repoCache,
		chartCache: chartCache,
	}
}

// Run runs the finalizer.
func (c *Finalizer) Run() {
	for {
		obj, shutdown := c.queue.Get()
		if shutdown {
			return
		}

		err := c.run(obj)
		if err != nil {
			c.log.Warningf("run: %v", err)
			continue
		}
		c.queue.Done(obj)
	}
}

// Stop stops the finalizer.
func (c *Finalizer) Stop() {
	c.log.Info("stopping...")
	c.queue.ShutDown()
}

func (c *Finalizer) run(obj interface{}) error {
	helmApp, ok := obj.(*v1alpha1.HelmApp)
	if !ok {
		return nil
	}

	logrus.Infof("start uninstall helm app: %s/%s", helmApp.Name, helmApp.Namespace)

	app, err := NewApp(c.ctx, c.kubeClient, c.clientset, helmApp, c.repoFile, c.repoCache, c.chartCache)
	if err != nil {
		return err
	}

	return app.Uninstall()
}
