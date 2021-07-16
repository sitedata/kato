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

// RbdEndpoints is a collection of RbdEndpoint.
type RbdEndpoints struct {
	Port        int      `json:"port"`
	IPs         []string `json:"ips"`
	NotReadyIPs []string `json:"not_ready_ips"`
}

// RbdEndpoint hold information to create k8s endpoints.
type RbdEndpoint struct {
	UUID     string `json:"uuid"`
	Sid      string `json:"sid"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Status   string `json:"status"`
	IsOnline bool   `json:"is_online"`
	Action   string `json:"action"`
	IsDomain bool   `json:"is_domain"`
}

// Equal tests for equality between two RbdEndpoint types
func (l1 *RbdEndpoint) Equal(l2 *RbdEndpoint) bool {
	return false
}
