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
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

//ContextKey
type ContextKey string

//RoundRobin loadBalance impl
type RoundRobin struct {
	ops *uint64
}

//LoadBalance
type LoadBalance interface {
	Select(r *http.Request, endpoints EndpointList) Endpoint
}

//Endpoint
type Endpoint string

func (e Endpoint) String() string {
	return string(e)
}

//GetName get endpoint name
func (e Endpoint) GetName() string {
	if kv := strings.Split(string(e), "=>"); len(kv) > 1 {
		return kv[0]
	}
	return string(e)
}

//GetAddr
func (e Endpoint) GetAddr() string {
	if kv := strings.Split(string(e), "=>"); len(kv) > 1 {
		return kv[1]
	}
	return string(e)
}

//GetHTTPAddr get http url
func (e Endpoint) GetHTTPAddr() string {
	if kv := strings.Split(string(e), "=>"); len(kv) > 1 {
		return withScheme(kv[1])
	}
	return withScheme(string(e))
}

func withScheme(s string) string {
	if strings.HasPrefix(s, "http") {
		return s
	}
	return "http://" + s
}

//EndpointList
type EndpointList []Endpoint

//Len
func (e *EndpointList) Len() int {
	return len(*e)
}

//Add
func (e *EndpointList) Add(endpoints ...string) {
	for _, end := range endpoints {
		*e = append(*e, Endpoint(end))
	}
}

//Delete
func (e *EndpointList) Delete(endpoints ...string) {
	var new EndpointList
	for _, endpoint := range endpoints {
		for _, old := range *e {
			if string(old) != endpoint {
				new = append(new, old)
			}
		}
	}
	*e = new
}

//Select
func (e *EndpointList) Selec(i int) Endpoint {
	return (*e)[i]
}

//HaveEndpoint whether or not there is an endpoint
func (e *EndpointList) HaveEndpoint(endpoint string) bool {
	for _, en := range *e {
		if en.String() == endpoint {
			return true
		}
	}
	return false
}

//CreateEndpoints
func CreateEndpoints(endpoints []string) EndpointList {
	var epl EndpointList
	for _, e := range endpoints {
		epl = append(epl, Endpoint(e))
	}
	return epl
}

// NewRoundRobin create a RoundRobin
func NewRoundRobin() LoadBalance {
	var ops uint64
	ops = 0
	return RoundRobin{
		ops: &ops,
	}
}

// Select select a server from servers using RoundRobin
func (rr RoundRobin) Select(r *http.Request, endpoints EndpointList) Endpoint {
	l := uint64(endpoints.Len())
	if 0 >= l {
		return ""
	}
	selec := int(atomic.AddUint64(rr.ops, 1) % l)
	return endpoints.Selec(selec)
}

//SelectBalance selective load balancing
type SelectBalance struct {
	hostIDMap map[string]string
}

//NewSelectBalance create selective load balancing
func NewSelectBalance() *SelectBalance {
	return &SelectBalance{
		hostIDMap: map[string]string{"local": "rbd-eventlog:6363"},
	}
}

//Select load
func (s *SelectBalance) Select(r *http.Request, endpoints EndpointList) Endpoint {
	if r.URL == nil {
		return Endpoint(s.hostIDMap["local"])
	}

	id2ip := map[string]string{"local": "rbd-eventlog:6363"}
	for _, end := range endpoints {
		if kv := strings.Split(string(end), "=>"); len(kv) > 1 {
			id2ip[kv[0]] = kv[1]
		}
	}

	if r.URL != nil {
		hostID := r.URL.Query().Get("host_id")
		if hostID == "" {
			hostIDFromContext := r.Context().Value(ContextKey("host_id"))
			if hostIDFromContext != nil {
				hostID = hostIDFromContext.(string)
			}
		}
		if e, ok := id2ip[hostID]; ok {
			logrus.Infof("[lb selelct] find host %s from name %s success", e, hostID)
			return Endpoint(e)
		}
	}

	if len(endpoints) > 0 {
		logrus.Infof("default endpoint is %s", endpoints[len(endpoints)-1])
		return endpoints[len(endpoints)-1]
	}

	return Endpoint(s.hostIDMap["local"])
}
