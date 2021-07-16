// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

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

package discovery

import (
	"strings"
)

// Discoverier blablabla
type Discoverier interface {
	Connect() error
	Fetch() ([]*Endpoint, error)
	Close() error
}

// NewDiscoverier creates a new Discoverier
func NewDiscoverier(info *Info) Discoverier {
	switch strings.ToUpper(info.Type) {
	case "ETCD":
		return NewEtcd(info)
	}
	return nil
}

// Info holds service discovery center information.
type Info struct {
	Type     string   `json:"type"`
	Servers  []string `json:"servers"`
	Key      string   `json:"key"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

// Endpoint holds endpoint and endpoint status(online or offline).
type Endpoint struct {
	Ep       string `json:"endpoint"`
	IsOnline bool   `json:"is_online"`
}
