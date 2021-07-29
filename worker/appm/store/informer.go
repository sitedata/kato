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

package store

import (
	"k8s.io/client-go/tools/cache"
)

//Informer kube-api client cache
type Informer struct {
	Namespace               cache.SharedIndexInformer
	Ingress                 cache.SharedIndexInformer
	Service                 cache.SharedIndexInformer
	Secret                  cache.SharedIndexInformer
	StatefulSet             cache.SharedIndexInformer
	Deployment              cache.SharedIndexInformer
	Pod                     cache.SharedIndexInformer
	ConfigMap               cache.SharedIndexInformer
	ReplicaSet              cache.SharedIndexInformer
	Endpoints               cache.SharedIndexInformer
	Nodes                   cache.SharedIndexInformer
	StorageClass            cache.SharedIndexInformer
	Claims                  cache.SharedIndexInformer
	Events                  cache.SharedIndexInformer
	HorizontalPodAutoscaler cache.SharedIndexInformer
	CRD                     cache.SharedIndexInformer
	HelmApp                 cache.SharedIndexInformer
	ComponentDefinition     cache.SharedIndexInformer
	ThirdComponent          cache.SharedIndexInformer
	CRS                     map[string]cache.SharedIndexInformer
}

//StartCRS -
func (i *Informer) StartCRS(stop chan struct{}) {
	for k := range i.CRS {
		go i.CRS[k].Run(stop)
	}
}

//Start start
func (i *Informer) Start(stop chan struct{}) {
	go i.Namespace.Run(stop)
	go i.Ingress.Run(stop)
	go i.Service.Run(stop)
	go i.Secret.Run(stop)
	go i.StatefulSet.Run(stop)
	go i.Deployment.Run(stop)
	go i.Pod.Run(stop)
	go i.ConfigMap.Run(stop)
	go i.ReplicaSet.Run(stop)
	go i.Endpoints.Run(stop)
	go i.Nodes.Run(stop)
	go i.StorageClass.Run(stop)
	go i.Events.Run(stop)
	go i.HorizontalPodAutoscaler.Run(stop)
	go i.Claims.Run(stop)
	go i.CRD.Run(stop)
	go i.HelmApp.Run(stop)
	go i.ComponentDefinition.Run(stop)
	go i.ThirdComponent.Run(stop)
}

//Ready if all kube informers is syncd, store is ready
func (i *Informer) Ready() bool {
	if i.Namespace.HasSynced() && i.Ingress.HasSynced() && i.Service.HasSynced() && i.Secret.HasSynced() &&
		i.StatefulSet.HasSynced() && i.Deployment.HasSynced() && i.Pod.HasSynced() &&
		i.ConfigMap.HasSynced() && i.Nodes.HasSynced() && i.Events.HasSynced() &&
		i.HorizontalPodAutoscaler.HasSynced() && i.StorageClass.HasSynced() && i.Claims.HasSynced() && i.CRD.HasSynced() {
		return true
	}
	return false
}
