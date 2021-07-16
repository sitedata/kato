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
	"strings"
	"time"

	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	"github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	"github.com/gridworkz/kato/pkg/helm"
	"github.com/sirupsen/logrus"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/workqueue"
)

const (
	defaultTimeout = 3 * time.Second
)

var defaultConditionTypes = []v1alpha1.HelmAppConditionType{
	v1alpha1.HelmAppChartReady,
	v1alpha1.HelmAppPreInstalled,
	v1alpha1.HelmAppInstalled,
}

// ControlLoop is a control loop to get helm app and reconcile it.
type ControlLoop struct {
	ctx        context.Context
	log        *logrus.Entry
	kubeClient clientset.Interface
	clientset  versioned.Interface
	storer     Storer
	workQueue  workqueue.Interface
	repo       *helm.Repo
	repoFile   string
	repoCache  string
	chartCache string
}

// NewControlLoop -
func NewControlLoop(ctx context.Context,
	kubeClient clientset.Interface,
	clientset versioned.Interface,
	storer Storer,
	workQueue workqueue.Interface,
	repoFile string,
	repoCache string,
	chartCache string,
) *ControlLoop {
	repo := helm.NewRepo(repoFile, repoCache)
	return &ControlLoop{
		ctx:        ctx,
		log:        logrus.WithField("WHO", "Helm App ControlLoop"),
		kubeClient: kubeClient,
		clientset:  clientset,
		storer:     storer,
		workQueue:  workQueue,
		repo:       repo,
		repoFile:   repoFile,
		repoCache:  repoCache,
		chartCache: chartCache,
	}
}

// Run runs the control loop.
func (c *ControlLoop) Run() {
	for {
		obj, shutdown := c.workQueue.Get()
		if shutdown {
			return
		}

		c.run(obj)
	}
}

// Stop stops the control loop.
func (c *ControlLoop) Stop() {
	c.log.Info("stopping...")
	c.workQueue.ShutDown()
}

func (c *ControlLoop) run(obj interface{}) {
	key, ok := obj.(string)
	if !ok {
		return
	}
	defer c.workQueue.Done(obj)
	name, ns := nameNamespace(key)

	helmApp, err := c.storer.GetHelmApp(ns, name)
	if err != nil {
		logrus.Warningf("[HelmAppController] [ControlLoop] get helm app(%s): %v", key, err)
		return
	}

	if err := c.Reconcile(helmApp); err != nil {
		// ignore the error, informer will push the same time into queue later.
		logrus.Warningf("[HelmAppController] [ControlLoop] [Reconcile]: %v", err)
		return
	}
}

// nameNamespace -
func nameNamespace(key string) (string, string) {
	strs := strings.Split(key, "/")
	return strs[0], strs[1]
}

// Reconcile -
func (c *ControlLoop) Reconcile(helmApp *v1alpha1.HelmApp) error {
	app, err := NewApp(c.ctx, c.kubeClient, c.clientset, helmApp, c.repoFile, c.repoCache, c.chartCache)
	if err != nil {
		return err
	}

	app.log.Debug("start reconcile")

	// update running status
	defer app.UpdateRunningStatus()

	// setups the default values of the helm app.
	if app.NeedSetup() {
		return app.Setup()
	}

	// detect the helm app.
	if app.NeedDetect() {
		return app.Detect()
	}

	// install or update the helm app.
	if app.NeedUpdate() {
		return app.InstallOrUpdate()
	}

	return nil
}
