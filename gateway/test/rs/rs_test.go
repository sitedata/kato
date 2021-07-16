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

package rs

import (
	"testing"

	"github.com/gridworkz/kato/gateway/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestReplicaSetTimestamp(t *testing.T) {
	clientset, err := controller.NewClientSet("/opt/kato/etc/kubernetes/kubecfg/admin.kubeconfig")
	if err != nil {
		t.Errorf("can't create Kubernetes's client: %v", err)
	}

	ns := "c1a29fe4d7b0413993dc859430cf743d"
	rs, err := clientset.ExtensionsV1beta1().ReplicaSets(ns).Get("88d8c4c55657217522f3bb86cfbded7e-deployment-7545b75dbd", metav1.GetOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
	t.Logf("%+v", rs)
}
