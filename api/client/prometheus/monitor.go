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

package prometheus

//Level -
type Level int

const (
	LevelCluster = 1 << iota
	LevelNode
	LevelWorkspace
	LevelNamespace
	LevelWorkload
	LevelPod
	LevelContainer
	LevelPVC
	LevelComponent
)

type QueryOption interface {
	Apply(*QueryOptions)
}

type QueryOptions struct {
	Level Level

	ResourceFilter            string
	NodeName                  string
	WorkspaceName             string
	NamespaceName             string
	WorkloadKind              string
	WorkloadName              string
	PodName                   string
	ContainerName             string
	StorageClassName          string
	PersistentVolumeClaimName string
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

type ClusterOption struct{}

func (_ ClusterOption) Apply(o *QueryOptions) {
	o.Level = LevelCluster
}

type NodeOption struct {
	ResourceFilter string
	NodeName       string
}

func (no NodeOption) Apply(o *QueryOptions) {
	o.Level = LevelNode
	o.ResourceFilter = no.ResourceFilter
	o.NodeName = no.NodeName
}

type WorkspaceOption struct {
	ResourceFilter string
	WorkspaceName  string
}

func (wo WorkspaceOption) Apply(o *QueryOptions) {
	o.Level = LevelWorkspace
	o.ResourceFilter = wo.ResourceFilter
	o.WorkspaceName = wo.WorkspaceName
}

type NamespaceOption struct {
	ResourceFilter string
	WorkspaceName  string
	NamespaceName  string
}

func (no NamespaceOption) Apply(o *QueryOptions) {
	o.Level = LevelNamespace
	o.ResourceFilter = no.ResourceFilter
	o.WorkspaceName = no.WorkspaceName
	o.NamespaceName = no.NamespaceName
}

type WorkloadOption struct {
	ResourceFilter string
	NamespaceName  string
	WorkloadKind   string
}

func (wo WorkloadOption) Apply(o *QueryOptions) {
	o.Level = LevelWorkload
	o.ResourceFilter = wo.ResourceFilter
	o.NamespaceName = wo.NamespaceName
	o.WorkloadKind = wo.WorkloadKind
}

type PodOption struct {
	ResourceFilter string
	NodeName       string
	NamespaceName  string
	WorkloadKind   string
	WorkloadName   string
	PodName        string
}

func (po PodOption) Apply(o *QueryOptions) {
	o.Level = LevelPod
	o.ResourceFilter = po.ResourceFilter
	o.NodeName = po.NodeName
	o.NamespaceName = po.NamespaceName
	o.WorkloadKind = po.WorkloadKind
	o.WorkloadName = po.WorkloadName
	o.PodName = po.PodName
}

type ContainerOption struct {
	ResourceFilter string
	NamespaceName  string
	PodName        string
	ContainerName  string
}

func (co ContainerOption) Apply(o *QueryOptions) {
	o.Level = LevelContainer
	o.ResourceFilter = co.ResourceFilter
	o.NamespaceName = co.NamespaceName
	o.PodName = co.PodName
	o.ContainerName = co.ContainerName
}

type PVCOption struct {
	ResourceFilter            string
	NamespaceName             string
	StorageClassName          string
	PersistentVolumeClaimName string
}

func (po PVCOption) Apply(o *QueryOptions) {
	o.Level = LevelPVC
	o.ResourceFilter = po.ResourceFilter
	o.NamespaceName = po.NamespaceName
	o.StorageClassName = po.StorageClassName
	o.PersistentVolumeClaimName = po.PersistentVolumeClaimName
}

type ComponentOption struct{}

func (_ ComponentOption) Apply(o *QueryOptions) {
	o.Level = LevelComponent
}
