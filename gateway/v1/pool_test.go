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

package v1

import "testing"

func TestPool_Equals(t *testing.T) {
	node1 := newFakeNode()
	node1.Name = "node-a"
	node2 := newFakeNode()
	node2.Name = "node-b"
	p := NewFakePoolWithoutNodes()
	p.Nodes = []*Node{
		node1,
		node2,
	}

	node3 := newFakeNode()
	node3.Name = "node-a"
	node4 := newFakeNode()
	node4.Name = "node-b"
	c := NewFakePoolWithoutNodes()
	c.Nodes = []*Node{
		node3,
		node4,
	}

	if !p.Equals(c) {
		t.Errorf("Pool p shoul equal Pool c")
	}
}

func NewFakePoolWithoutNodes() *Pool {
	return &Pool{
		Meta: Meta{
			Index:      888,
			Name:       "foo-pool",
			Namespace:  "gateway",
			PluginName: "Nginx",
		},
		ServiceID:         "foo-service-id",
		ServiceVersion:    "1.0.0",
		ServicePort:       80,
		Note:              "foo",
		NodeNumber:        8,
		LoadBalancingType: RoundRobin,
		Monitors: []Monitor{
			"monitor-a",
			"monitor-b",
			"monitor-c",
		},
	}
}
