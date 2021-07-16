// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

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

package node

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCluster_handleNodeStatus(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/gridworkz/Documents/company/gridworkz/remote/192.168.2.200/admin.kubeconfig")
	if err != nil {
		return
	}
	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}

	node, err := cli.CoreV1().Nodes().Get("192.168.2.200", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("node is :%+v", node)
	t.Logf("cpu:%v", node.Status.Allocatable.Cpu().Value())
	t.Logf("mem: %v", node.Status.Allocatable.Memory().Value())
}
