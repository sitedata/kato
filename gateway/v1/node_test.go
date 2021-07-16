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

import (
	"testing"
)

func TestNode_Equals(t *testing.T) {
	n := newFakeNode()
	c := newFakeNode()
	if !n.Equals(c) {
		t.Errorf("n should equal c")
	}
	f := newFakeNode()
	f.MaxFails = 5
	if n.Equals(f) {
		t.Errorf("n should not equal c")
	}
}

func newFakeNode() *Node {
	return &Node{
		Meta: Meta{
			Index:      888,
			Name:       "foo-node",
			Namespace:  "ns",
			PluginName: "Nginx",
		},
		Host:        "www.gridworkz.com",
		Port:        80,
		Protocol:    "Http",
		State:       "ok",
		PoolName:    "foo-poolName",
		Ready:       true,
		Weight:      5,
		MaxFails:    3,
		FailTimeout: "5",
	}
}
