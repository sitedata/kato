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

package test

import (
	"testing"
	"time"

	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/kubecache"
)

func TestGetCluster(t *testing.T) {
	c := &option.Conf{
		K8SConfPath:     "/opt/kato/etc/kubernetes/kubecfg/admin.kubeconfig",
		MinResyncPeriod: 10 * time.Second,
	}
	kubecli, err := kubecache.NewKubeClient(c)
	if err != nil {
		t.Fatalf("error creating kube client: %v", err)
	}
	defer kubecli.Stop()

	nodes, err := kubecli.GetNodes()
	if err != nil {
		t.Errorf("error getting nodes: %v", err)
	}
	t.Log(nodes)
	t.Error("")
}
