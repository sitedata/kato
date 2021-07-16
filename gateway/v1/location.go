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
	"github.com/gridworkz/kato/gateway/annotations/proxy"
	"github.com/gridworkz/kato/gateway/annotations/rewrite"
)

// ConditionType condition type
type ConditionType string

// HeaderType -
var HeaderType ConditionType = "header"

// CookieType -
var CookieType ConditionType = "cookie"

// DefaultType -
var DefaultType ConditionType = "default"

// Location -
type Location struct {
	Path          string
	NameCondition map[string]*Condition // papping between backend name and condition
	// Rewrite describes the redirection this location.
	// +optional
	Rewrite rewrite.Config `json:"rewrite,omitempty"`
	// Proxy contains information about timeouts and buffer sizes
	// to be used in connections against endpoints
	// +optional
	Proxy            proxy.Config `json:"proxy,omitempty"`
	DisableProxyPass bool
}

// Condition is the condition that the traffic can reach the specified backend
type Condition struct {
	Type  ConditionType
	Value map[string]string
}

// Equals determines if two locations are equal
func (l *Location) Equals(c *Location) bool {
	if l == c {
		return true
	}
	if l == nil || c == nil {
		return false
	}
	if l.Path != c.Path {
		return false
	}

	if len(l.NameCondition) != len(c.NameCondition) {
		return false
	}
	for name, lc := range l.NameCondition {
		if cc, exists := c.NameCondition[name]; !exists || !cc.Equals(lc) {
			return false
		}
	}

	if !l.Proxy.Equal(&c.Proxy) {
		return false
	}

	return true
}

// Equals determines if two conditions are equal
func (c *Condition) Equals(cc *Condition) bool {
	if c == cc {
		return true
	}
	if c == nil || cc == nil {
		return false
	}
	if c.Type != cc.Type {
		return false
	}

	if len(c.Value) != len(cc.Value) {
		return false
	}
	for k, v := range c.Value {
		if vv, ok := cc.Value[k]; !ok || v != vv {
			return false
		}
	}

	return true
}

func newFakeLocation() *Location {
	return &Location{
		Path: "foo-path",
	}
}
