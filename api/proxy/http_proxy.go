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

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

//HTTPProxy
type HTTPProxy struct {
	name      string
	endpoints EndpointList
	lb        LoadBalance
	client    *http.Client
}

//Proxy http proxy
func (h *HTTPProxy) Proxy(w http.ResponseWriter, r *http.Request) {
	endpoint := h.lb.Select(r, h.endpoints)
	endURL, err := url.Parse(endpoint.GetHTTPAddr())
	if err != nil {
		logrus.Errorf("parse endpoint url error,%s", err.Error())
		w.WriteHeader(502)
		return
	}
	if endURL.Scheme == "" {
		endURL.Scheme = "http"
	}
	proxy := httputil.NewSingleHostReverseProxy(endURL)
	proxy.ServeHTTP(w, r)
}

//UpdateEndpoints 
func (h *HTTPProxy) UpdateEndpoints(endpoints ...string) {
	ends := []string{}
	for _, end := range endpoints {
		if kv := strings.Split(end, "=>"); len(kv) > 1 {
			ends = append(ends, kv[1])
		} else {
			ends = append(ends, end)
		}
	}
	h.endpoints = CreateEndpoints(ends)
}

//Do proxy
func (h *HTTPProxy) Do(r *http.Request) (*http.Response, error) {
	endpoint := h.lb.Select(r, h.endpoints)
	if strings.HasPrefix(endpoint.String(), "http") {
		r.URL.Host = strings.Replace(endpoint.String(), "http://", "", 1)
	} else {
		r.URL.Host = endpoint.String()
	}
	//default is http
	r.URL.Scheme = "http"
	return h.client.Do(r)
}

func createHTTPProxy(name string, endpoints []string, lb LoadBalance) *HTTPProxy {
	ends := []string{}
	for _, end := range endpoints {
		if kv := strings.Split(end, "=>"); len(kv) > 1 {
			ends = append(ends, kv[1])
		} else {
			ends = append(ends, end)
		}
	}
	if lb == nil {
		lb = NewRoundRobin()
	}
	timeout, _ := strconv.Atoi(os.Getenv("PROXY_TIMEOUT"))
	if timeout == 0 {
		timeout = 10
	}
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Transport: netTransport,
		Timeout:   time.Second * time.Duration(timeout),
	}
	return &HTTPProxy{name, CreateEndpoints(ends), lb, client}
}
