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

package model

import (
	"strconv"

	"github.com/sirupsen/logrus"

	v1 "github.com/gridworkz/kato/gateway/v1"
	apiv1 "k8s.io/api/core/v1"
)

//Config update config
type Config struct {
	Backends []*Backend `json:"backends"`
}

// Backend describes one or more remote server/s (endpoints) associated with a service
type Backend struct {
	// Name represents an unique apiv1.Service name formatted as <namespace>-<name>-<port>
	Name string `json:"name"`

	Endpoints []Endpoint `json:"endpoints,omitempty"`
	// StickySessionAffinitySession contains the StickyConfig object with stickyness configuration
	SessionAffinity SessionAffinityConfig `json:"sessionAffinityConfig"`
	// Consistent hashing by NGINX variable
	UpstreamHashBy string `json:"upstream-hash-by,omitempty"`
	// LB algorithm configuration per ingress
	LoadBalancing string `json:"load-balance,omitempty"`
}

// SessionAffinityConfig describes different affinity configurations for new sessions.
// Once a session is mapped to a backend based on some affinity setting, it
// retains that mapping till the backend goes down, or the ingress controller
// restarts. Exactly one of these values will be set on the upstream, since multiple
// affinity values are incompatible. Once set, the backend makes no guarantees
// about honoring updates.
type SessionAffinityConfig struct {
	AffinityType          string                `json:"name"`
	CookieSessionAffinity CookieSessionAffinity `json:"cookieSessionAffinity"`
}

// CookieSessionAffinity defines the structure used in Affinity configured by Cookies.
// +k8s:deepcopy-gen=true
type CookieSessionAffinity struct {
	Name      string              `json:"name"`
	Hash      string              `json:"hash"`
	Expires   string              `json:"expires,omitempty"`
	MaxAge    string              `json:"maxage,omitempty"`
	Locations map[string][]string `json:"locations,omitempty"`
	Path      string              `json:"path,omitempty"`
}

// Endpoint describes a kubernetes endpoint in a backend
// +k8s:deepcopy-gen=true
type Endpoint struct {
	// Address IP address of the endpoint
	Address string `json:"address"`
	// Port number of the TCP port
	Port string `json:"port"`
	// Weight weight of the endpoint
	Weight int `json:"weight"`
	// Target returns a reference to the object providing the endpoint
	Target *apiv1.ObjectReference `json:"target,omitempty"`
}

//CreateBackendByPool create backend by pool
func CreateBackendByPool(pool *v1.Pool) *Backend {
	var backend = Backend{
		Name: pool.Name,
	}
	switch pool.LoadBalancingType {
	case v1.RoundRobin:
		backend.LoadBalancing = "round_robin"
	case v1.CookieSessionAffinity:
		logrus.Infof("pool %s use cookie-session-affinity load balance", pool.Name)
		backend.SessionAffinity = SessionAffinityConfig{
			AffinityType: "cookie",
			CookieSessionAffinity: CookieSessionAffinity{
				Name: "kato-route",
			},
		}
	}
	backend.UpstreamHashBy = pool.UpstreamHashBy
	var endpoints []Endpoint
	for _, node := range pool.Nodes {
		endpoints = append(endpoints, Endpoint{
			Address: node.Host,
			Port:    strconv.Itoa(int(node.Port)),
			Weight:  node.Weight,
		})
	}
	backend.Endpoints = endpoints
	return &backend
}
