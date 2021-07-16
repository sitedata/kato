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
	"crypto/x509"
	"time"
)

// SSLCert describes a SSL certificate
type SSLCert struct {
	*Meta
	CertificateStr string            `json:"certificate_str"`
	Certificate    *x509.Certificate `json:"certificate,omitempty"`
	PrivateKey     string            `json:"private_key"`
	CertificatePem string            `json:"certificate_pem"`
	// CN contains all the common names defined in the SSL certificate
	CN []string `json:"cn"`
	// ExpiresTime contains the expiration of this SSL certificate in timestamp format
	ExpireTime time.Time `json:"expires"`
}

//Equals -
func (s *SSLCert) Equals(c *SSLCert) bool {
	if s == c {
		return true
	}
	if s == nil || c == nil {
		return false
	}
	if !s.Meta.Equals(c.Meta) {
		return false
	}
	if (s.Certificate == nil) != (c.Certificate == nil) {
		return false
	}
	if s.Certificate != nil && c.Certificate != nil {
		if !s.Certificate.Equal(c.Certificate) {
			return false
		}
	}
	if s.CertificateStr != c.CertificateStr {
		return false
	}
	if s.PrivateKey != c.PrivateKey {
		return false
	}

	if len(s.CN) != len(c.CN) {
		return false
	}
	for _, scn := range s.CN {
		flag := false
		for _, ccn := range c.CN {
			if scn != ccn {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}

	if !s.ExpireTime.Equal(c.ExpireTime) {
		return false
	}
	return true
}
