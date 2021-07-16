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

package cert

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"testing"
)

func Test_crt(t *testing.T) {
	baseinfo := CertInformation{Country: []string{"CN"}, Organization: []string{"Gridworkz"}, IsCA: true,
		OrganizationalUnit: []string{"work-stacks"}, EmailAddress: []string{"gdevs@gridworkz.com"},
		Locality: []string{"Ontario"}, Province: []string{"Ontario"}, CommonName: "Work-Stacks",
		Domains: []string{"gridworkz"}, CrtName: "../../test/ssl/ca.pem", KeyName: "../../test/ssl/ca.key"}

	err := CreateCRT(nil, nil, baseinfo)
	if err != nil {
		t.Log("Create crt error,Error info:", err)
		return
	}
	crtinfo := baseinfo
	crtinfo.IsCA = false
	crtinfo.CrtName = "../../test/ssl/api_server.pem"
	crtinfo.KeyName = "../../test/ssl/api_server.key"
	crtinfo.Names = []pkix.AttributeTypeAndValue{
		pkix.AttributeTypeAndValue{
			Type:  asn1.ObjectIdentifier{2, 1, 3},
			Value: "MAC_ADDR",
		},
	}

	crt, pri, err := Parse(baseinfo.CrtName, baseinfo.KeyName)
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	err = CreateCRT(crt, pri, crtinfo)
	if err != nil {
		t.Log("Create crt error,Error info:", err)
	}
	//os.Remove(baseinfo.CrtName)
	//os.Remove(baseinfo.KeyName)
	//os.Remove(crtinfo.CrtName)
	//os.Remove(crtinfo.KeyName)
}
