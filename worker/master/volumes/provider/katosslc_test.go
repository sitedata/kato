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

package provider

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes"
)

func TestSelectNode(t *testing.T) {
	c, err := clientcmd.BuildConfigFromFlags("", "../../../../test/admin.kubeconfig")
	if err != nil {
		t.Fatal(err)
	}
	client, _ := kubernetes.NewForConfig(c)
	pr := &katosslcProvisioner{
		name:    "kato.io/provisioner-sslc",
		kubecli: client,
	}
	node, err := pr.selectNode("linux", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(node)
}

func TestGetVolumeIDByPVCName(t *testing.T) {
	t.Log(getVolumeIDByPVCName("manual17-gra02c40-0"))
	t.Log(getVolumeIDByPVCName("manual17"))
}

func TestGetPodNameByPVCName(t *testing.T) {
	t.Log(getPodNameByPVCName("manual17-gra02c40-0"))
}
