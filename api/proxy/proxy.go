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

package proxy

import "net/http"

//Proxy proxy
type Proxy interface {
	Proxy(w http.ResponseWriter, r *http.Request)
	Do(r *http.Request) (*http.Response, error)
	UpdateEndpoints(endpoints ...string) // format: ["name=>ip:port", ...]
}

//CreateProxy create proxy
func CreateProxy(name string, mode string, endpoints []string) Proxy {
	switch mode {
	case "websocket":
		return createWebSocketProxy(name, endpoints)
	case "http":
		if name == "eventlog" {
			return createHTTPProxy(name, endpoints, NewSelectBalance())
		}
		return createHTTPProxy(name, endpoints, nil)
	default:
		return createHTTPProxy(name, endpoints, nil)
	}
}
