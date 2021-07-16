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
	"sync"

	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	corev1 "k8s.io/api/core/v1"
)

//ResourceCache resource cache
type ResourceCache struct {
	lock      sync.Mutex
	resources map[string]*NamespaceResource
}

//NewResourceCache new resource cache
func NewResourceCache() *ResourceCache {
	return &ResourceCache{
		resources: make(map[string]*NamespaceResource),
	}
}

//NamespaceResource namespace resource
type NamespaceResource map[string]*v1.PodResource

//SetPodResource set pod resource
func (r *NamespaceResource) SetPodResource(podName string, pr *v1.PodResource) {
	(*r)[podName] = pr
}

//RemovePod remove pod resource
func (r *NamespaceResource) RemovePod(podName string) {
	delete(*r, podName)
}

//TenantResource tenant resource
type TenantResource struct {
	Namespace     string
	MemoryRequest int64
	MemoryLimit   int64
	CPURequest    int64
	CPULimit      int64
}

//SetPodResource set pod resource
func (r *ResourceCache) SetPodResource(pod *corev1.Pod) {
	r.lock.Lock()
	defer r.lock.Unlock()
	namespace := pod.Namespace
	re := v1.CalculatePodResource(pod)
	if nr, ok := r.resources[namespace]; ok && nr != nil {
		nr.SetPodResource(pod.Name, re)
	} else {
		nameR := make(NamespaceResource)
		nameR.SetPodResource(pod.Name, re)
		r.resources[namespace] = &nameR
	}
}

//RemovePod remove pod resource
func (r *ResourceCache) RemovePod(pod *corev1.Pod) {
	r.lock.Lock()
	defer r.lock.Unlock()
	namespace := pod.Namespace
	if nr, ok := r.resources[namespace]; ok && nr != nil {
		nr.RemovePod(pod.Name)
	}
}

//GetTenantResource get tenant resource
func (r *ResourceCache) GetTenantResource(namespace string) (tr TenantResource) {
	r.lock.Lock()
	defer r.lock.Unlock()
	tr = r.getTenantResource(r.resources[namespace])
	tr.Namespace = namespace
	return tr
}

func (r *ResourceCache) getTenantResource(namespaceRe *NamespaceResource) (tr TenantResource) {
	if namespaceRe == nil {
		return
	}
	for _, v := range *namespaceRe {
		tr.CPULimit += v.CPULimit
		tr.MemoryLimit += v.MemoryLimit
		tr.CPURequest += v.CPURequest
		tr.MemoryRequest += v.MemoryRequest
	}
	return
}

//GetAllTenantResource get all tenant resources
func (r *ResourceCache) GetAllTenantResource() (trs []TenantResource) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for k := range r.resources {
		tr := r.getTenantResource(r.resources[k])
		tr.Namespace = k
		trs = append(trs, tr)
	}
	return
}
