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

//Node Application service endpoint
type Node struct {
	Meta
	Host        string `json:"host"`
	Port        int32  `json:"port"`
	Protocol    string `json:"protocol"`
	State       string `json:"state"`     //Active Draining Disabled
	PoolName    string `json:"pool_name"` //Belong to the pool
	Ready       bool   `json:"ready"`     //Whether ready
	Weight      int    `json:"weight"`
	MaxFails    int    `json:"max_fails"`
	FailTimeout string `json:"fail_timeout"`
}

//Equals -
func (n *Node) Equals(c *Node) bool { //
	if n == c {
		return true
	}
	if n == nil || c == nil {
		return false
	}
	if n.Meta != c.Meta {
		return false
	}
	if n.Host != c.Host {
		return false
	}
	if n.Port != c.Port {
		return false
	}
	if n.Protocol != c.Protocol {
		return false
	}
	if n.State != c.State {
		return false
	}
	if n.PoolName != c.PoolName {
		return false
	}
	if n.Ready != c.Ready {
		return false
	}
	if n.Weight != c.Weight {
		return false
	}
	if n.MaxFails != c.MaxFails {
		return false
	}
	if n.FailTimeout != c.FailTimeout {
		return false
	}
	return true
}
