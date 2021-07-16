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

import "testing"

func TestVirtualService_Equals(t *testing.T) {
	v := newFakeVirtualService()
	vlocA := newFakeLocation()
	vlocB := newFakeLocation()
	v.Locations = append(v.Locations, vlocA)
	v.Locations = append(v.Locations, vlocB)
	v.SSLCert = newFakeSSLCert()

	c := newFakeVirtualService()
	clocA := newFakeLocation()
	clocB := newFakeLocation()
	c.Locations = append(c.Locations, clocA)
	c.Locations = append(c.Locations, clocB)
	c.SSLCert = newFakeSSLCert()

	if !v.Equals(c) {
		t.Errorf("v should equal c")
	}
}

func newFakeVirtualService() *VirtualService {
	return &VirtualService{
		Meta:                   newFakeMeta(),
		Enabled:                true,
		Protocol:               "Http",
		BackendProtocol:        "Http",
		Port:                   80,
		Listening:              []string{"a", "b", "c"},
		Note:                   "foo-node",
		DefaultPoolName:        "default-pool-name",
		RuleNames:              []string{"a", "b", "c"},
		SSLdecrypt:             true,
		DefaultCertificateName: "default-certificate-name",
		RequestLogEnable:       true,
		RequestLogFileName:     "/var/log/gateway/request.log",
		RequestLogFormat:       "request-log-format",
		ConnectTimeout:         70,
		Timeout:                70,
		ServerName:             "foo-server_name",
		PoolName:               "foo-pool-name",
		ForceSSLRedirect:       true,
	}
}
