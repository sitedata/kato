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
	"fmt"
	"time"

	katov1alpha1 "github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	"github.com/gridworkz/kato/pkg/generated/listers/kato/v1alpha1"
	k8sutil "github.com/gridworkz/kato/util/k8s"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Storer -
type Storer interface {
	Run(stopCh <-chan struct{})
	GetHelmApp(ns, name string) (*katov1alpha1.HelmApp, error)
}

type store struct {
	informer cache.SharedIndexInformer
	lister   v1alpha1.HelmAppLister
}

// NewStorer creates a new storer.
func NewStorer(informer cache.SharedIndexInformer,
	lister v1alpha1.HelmAppLister,
	workqueue workqueue.Interface,
	finalizerQueue workqueue.Interface) Storer {
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			helmApp := obj.(*katov1alpha1.HelmApp)
			workqueue.Add(k8sutil.ObjKey(helmApp))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			helmApp := newObj.(*katov1alpha1.HelmApp)
			workqueue.Add(k8sutil.ObjKey(helmApp))
		},
		DeleteFunc: func(obj interface{}) {
			// Two purposes of using finalizerQueue
			// 1. non-block DeleteFunc
			// 2. retry if the finalizer is failed
			finalizerQueue.Add(obj)
		},
	})
	return &store{
		informer: informer,
		lister:   lister,
	}
}

func (i *store) Run(stopCh <-chan struct{}) {
	go i.informer.Run(stopCh)

	// wait for all involved caches to be synced before processing items
	// from the queue
	if !cache.WaitForCacheSync(stopCh,
		i.informer.HasSynced,
	) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
	}

	// in big clusters, deltas can keep arriving even after HasSynced
	// functions have returned 'true'
	time.Sleep(1 * time.Second)
}

func (i *store) GetHelmApp(ns, name string) (*katov1alpha1.HelmApp, error) {
	return i.lister.HelmApps(ns).Get(name)
}
