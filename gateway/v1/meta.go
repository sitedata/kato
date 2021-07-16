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

//Meta Common meta
type Meta struct {
	Index      int64  `json:"index"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	ServiceID  string `json:"service_id"`
	PluginName string `json:"plugin_name"`
}

//Equals -
func (m *Meta) Equals(c *Meta) bool {
	if m == c {
		return true
	}
	if m == nil || c == nil {
		return false
	}
	if m.Name != c.Name {
		return false
	}
	if m.Namespace != c.Namespace {
		return false
	}
	if m.PluginName != c.PluginName {
		return false
	}
	if m.ServiceID != c.ServiceID {
		return false
	}
	return true
}
