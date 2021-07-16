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

//Pool Application service endpoints pool
type Pool struct {
	Meta
	//application service id
	ServiceID string `json:"service_id"`
	//application service version
	ServiceVersion string `json:"service_version"`
	//application service port
	ServicePort int `json:"service_port"`
	//pool instructions
	Note              string            `json:"note"`
	NodeNumber        int               `json:"node_number"`
	LoadBalancingType LoadBalancingType `json:"load_balancing_type"`
	UpstreamHashBy    string            `json:"upstream_hash_by"`
	LeastConn         bool              `json:"least_conn"`
	Monitors          []Monitor         `json:"monitors"`
	Nodes             []*Node           `json:"nodes"`
}

//Equals -
func (p *Pool) Equals(c *Pool) bool {
	if p == c {
		return true
	}
	if p == nil || c == nil {
		return false
	}
	if !p.Meta.Equals(&c.Meta) {
		return false
	}
	if p.ServiceID != c.ServiceID {
		return false
	}
	if p.ServiceVersion != c.ServiceVersion {
		return false
	}
	if p.ServicePort != c.ServicePort {
		return false
	}
	if p.Note != c.Note {
		return false
	}
	if p.NodeNumber != c.NodeNumber {
		return false
	}
	if p.LoadBalancingType != c.LoadBalancingType {
		return false
	}

	if len(p.Monitors) != len(c.Monitors) {
		return false
	}
	for _, a := range p.Monitors {
		flag := false
		for _, b := range c.Monitors {
			if a == b {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}

	if len(p.Nodes) != len(c.Nodes) {
		return false
	}
	for _, a := range p.Nodes {
		flag := false
		for _, b := range c.Nodes {
			if a.Equals(b) {
				flag = true
			}
		}
		if !flag {
			return false
		}
	}

	return true
}
