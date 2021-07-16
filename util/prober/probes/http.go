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

package probe

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gridworkz/kato/node/nodem/service"
	v1 "github.com/gridworkz/kato/util/prober/types/v1"
	"github.com/sirupsen/logrus"
)

// HTTPProbe probes through the http protocol
type HTTPProbe struct {
	Name          string
	Address       string
	ResultsChan   chan *v1.HealthStatus
	Ctx           context.Context
	Cancel        context.CancelFunc
	TimeInterval  int
	MaxErrorsNum  int
	TimeoutSecond int
}

//Check starts http probe.
func (h *HTTPProbe) Check() {
	go h.HTTPCheck()
}

//Stop http probe.
func (h *HTTPProbe) Stop() {
	h.Cancel()
}

//HTTPCheck
func (h *HTTPProbe) HTTPCheck() {
	if h.TimeInterval == 0 {
		h.TimeInterval = 5
	}
	timer := time.NewTimer(time.Second * time.Duration(h.TimeInterval))
	defer timer.Stop()
	for {
		HealthMap := h.GetHTTPHealth()
		result := &v1.HealthStatus{
			Name:   h.Name,
			Status: HealthMap["status"],
			Info:   HealthMap["info"],
		}
		h.ResultsChan <- result
		timer.Reset(time.Second * time.Duration(h.TimeInterval))
		select {
		case <-h.Ctx.Done():
			return
		case <-timer.C:
		}
	}
}

// Return true if the underlying error indicates a http.Client timeout.
//
// Use for errors returned from http.Client methods (Get, Post).
func isClientTimeout(err error) bool {
	if uerr, ok := err.(*url.Error); ok {
		if nerr, ok := uerr.Err.(net.Error); ok && nerr.Timeout() {
			return true
		}
	}
	return false
}

//GetHTTPHealth
func (h *HTTPProbe) GetHTTPHealth() map[string]string {
	address := h.Address
	c := &http.Client{
		Timeout: time.Duration(h.TimeoutSecond) * time.Second,
	}
	if strings.HasPrefix(address, "https://") {
		c.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		logrus.Warnf("address %s do not has scheme, auto add http scheme", address)
		address = "http://" + address
	}
	addr, err := url.Parse(address)
	if err != nil {
		logrus.Errorf("%s is invalid %s", address, err.Error())
		return map[string]string{"status": service.Stat_healthy, "info": "check url is invalid"}
	}
	if addr.Scheme == "" {
		addr.Scheme = "http"
	}
	logrus.Debugf("http probe check address; %s", address)
	resp, err := c.Get(addr.String())
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if isClientTimeout(err) {
			return map[string]string{"status": service.Stat_death, "info": "Request service timeout"}
		}
		logrus.Debugf("http probe request error %s", err.Error())
		return map[string]string{"status": service.Stat_unhealthy, "info": err.Error()}
	}
	if resp.StatusCode >= 400 {
		logrus.Debugf("http probe check address %s return code %d", address, resp.StatusCode)
		return map[string]string{"status": service.Stat_unhealthy, "info": "Service unhealthy"}
	}
	return map[string]string{"status": service.Stat_healthy, "info": "service health"}
}
