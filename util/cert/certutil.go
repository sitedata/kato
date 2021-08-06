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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"net"
	"os"
	"time"
)

func init() {
	rd.Seed(time.Now().UnixNano())
}

type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	Locality           []string
	CommonName         string
	CrtName, KeyName   string
	IsCA               bool
	Names              []pkix.AttributeTypeAndValue
	IPAddresses        []net.IP
	Domains            []string
}

//CreateCRT
func CreateCRT(RootCa *x509.Certificate, RootKey *rsa.PrivateKey, info CertInformation) error {
	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		//create ca cert
		buf, err = x509.CreateCertificate(rand.Reader, Crt, Crt, &Key.PublicKey, Key)
		if err != nil {
			return err
		}
		keybuf := x509.MarshalPKCS1PrivateKey(Key)
		err = write(info.KeyName, "PRIVATE KEY", keybuf)
	} else {
		//create cert by ca
		buf, err = x509.CreateCertificate(rand.Reader, Crt, RootCa, &Key.PublicKey, RootKey)
		if err != nil {
			return err
		}
		keybuf := x509.MarshalPKCS1PrivateKey(Key)
		err = write(info.KeyName, "RSA PRIVATE KEY", keybuf)
	}
	if err != nil {
		return err
	}
	err = write(info.CrtName, "CERTIFICATE", buf)
	if err != nil {
		return err
	}
	return nil
}

//Write encoding to file
func write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	var b = &pem.Block{Bytes: p, Type: Type}
	return pem.Encode(File, b)
}

//Parse
func Parse(crtPath, keyPath string) (rootcertificate *x509.Certificate, rootPrivateKey *rsa.PrivateKey, err error) {
	rootcertificate, err = ParseCrt(crtPath)
	if err != nil {
		return
	}
	rootPrivateKey, err = ParseKey(keyPath)
	return
}

//ParseCrt
func ParseCrt(path string) (*x509.Certificate, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	return x509.ParseCertificate(p.Bytes)
}

//ParseKey
func ParseKey(path string) (*rsa.PrivateKey, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, buf := pem.Decode(buf)
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func newCertificate(info CertInformation) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			CommonName:         info.CommonName,
			Locality:           info.Locality,
			ExtraNames:         info.Names,
		},
		NotBefore:             time.Now(),                                                                 //start time
		NotAfter:              time.Now().AddDate(20, 0, 0),                                               //end time
		BasicConstraintsValid: true,                                                                       //basic
		IsCA:                  info.IsCA,                                                                  //is it a root certificate?
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //certificate purpose
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		EmailAddresses:        info.EmailAddress,
		IPAddresses:           info.IPAddresses,
		DNSNames:              info.Domains,
	}
}

//CreateCertInformation
func CreateCertInformation() CertInformation {
	baseinfo := CertInformation{
		Country:            []string{"CN"},
		Organization:       []string{"Gridworkz"},
		OrganizationalUnit: []string{"gridworkz kato"},
		EmailAddress:       []string{"gdevs@gridworkz.com"},
		Locality:           []string{"Toronto"},
		Province:           []string{"Ontario"},
		CommonName:         "kato",
		CrtName:            "",
		KeyName:            "",
		Domains:            []string{"gridworkz"},
	}
	baseinfo.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	return baseinfo
}
