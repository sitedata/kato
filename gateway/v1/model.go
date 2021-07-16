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

import "github.com/sirupsen/logrus"

//LoadBalancingType Load Balancing type
type LoadBalancingType string

//RoundRobin Assign requests in turn to each node.
var RoundRobin LoadBalancingType = "round-robin"

//CookieSessionAffinity session affinity by cookie
var CookieSessionAffinity LoadBalancingType = "cookie-session-affinity"

//GetLoadBalancingType get load balancing
func GetLoadBalancingType(s string) LoadBalancingType {
	switch s {
	case "round-robin":
		return RoundRobin
	case "cookie-session-affinity":
		return CookieSessionAffinity
	default:
		return RoundRobin
	}
}

//Monitor monitor type
type Monitor string

//ConnectMonitor tcp connect monitor
var ConnectMonitor Monitor = "connect"

//PingMonitor ping monitor
var PingMonitor Monitor = "ping"

//SimpleHTTP http monitor
var SimpleHTTP Monitor = "simple http"

//SimpleHTTPS http monitor
var SimpleHTTPS Monitor = "simple https"

//HTTPRule Application service access rule for http
type HTTPRule struct {
	Meta
	Domain       string            `json:"domain"`
	Path         string            `json:"path"`
	Headers      map[string]string `json:"headers"`
	Redirect     RedirectConfig    `json:"redirect,omitempty"`
	HTTPSEnabale bool              `json:"https_enable"`
	SSLCertName  string            `json:"ssl_cert_name"`
	PoolName     string            `json:"pool_name"`
}

//RedirectConfig Config returns the redirect configuration for an  rule
type RedirectConfig struct {
	URL       string `json:"url"`
	Code      int    `json:"code"`
	FromToWWW bool   `json:"fromToWWW"`
}

// Config contains all the configuration of the gateway
type Config struct {
	HTTPPools []*Pool
	TCPPools  []*Pool
	L7VS      []*VirtualService
	L4VS      []*VirtualService
}

// Equals determines if cfg is equal to c
func (cfg *Config) Equals(c *Config) bool {
	if cfg == c {
		return true
	}

	if cfg == nil || c == nil {
		return false
	}

	if len(cfg.L7VS) != len(c.L7VS) {
		return false
	}
	for _, cfgv := range cfg.L7VS {
		flag := false
		for _, cv := range c.L7VS {
			if cfgv.Equals(cv) {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	logrus.Debugf("len if cnf.L4VS = %d, l4vs = %d", len(cfg.L4VS), len(c.L4VS))
	if len(cfg.L4VS) != len(c.L4VS) {
		return false
	}
	for _, cfgv := range cfg.L4VS {
		flag := false
		for _, cv := range c.L4VS {
			if cfgv.Equals(cv) {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}

	return true
}
