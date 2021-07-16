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

package client_test

import (
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	conf "github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/core/store"
	"github.com/gridworkz/kato/node/nodem/client"
)

func TestHostNodeMergeLabels(t *testing.T) {
	t.Parallel() // TODO: parallel
	hostNode := client.HostNode{
		Labels: map[string]string{
			"label 1": "value 1",
			"label 2": "value 2",
		},
		CustomLabels: map[string]string{
			"label a": "value a",
			"label b": "value b",
		},
	}
	sysLabelsLen := len(hostNode.Labels)
	exp := map[string]string{
		"label 1": "value 1",
		"label 2": "value 2",
		"label a": "value a",
		"label b": "value b",
	}
	labels := hostNode.MergeLabels()
	if len(exp) != len(labels) {
		t.Errorf("Expected %d for lables, but returned %d.", len(exp), len(labels))
	}
	equal := true
	for k, v := range exp {
		if labels[k] != v {
			equal = false
		}
	}
	if !equal {
		t.Errorf("Expected %+v for labels, but returned %+v", exp, labels)
	}
	if sysLabelsLen != len(hostNode.Labels) {
		t.Errorf("Expected %d for the length of system labels, but returned %+v", sysLabelsLen, len(hostNode.Labels))
	}
}

func TestHostNode_DelEndpoints(t *testing.T) {
	cfg := &conf.Conf{
		Etcd: clientv3.Config{
			Endpoints:   []string{"http://192.168.3.252:2379"},
			DialTimeout: 3 * time.Second,
		},
	}
	err := store.NewClient(cfg)
	if err != nil {
		t.Fatalf("error creating etcd client: %v", err)
	}
	n := &client.HostNode{
		InternalIP: "192.168.2.54",
	}
	n.DelEndpoints()
}
