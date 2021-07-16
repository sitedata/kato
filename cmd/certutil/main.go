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

package main

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"net"
	"os"
	"sort"

	version "github.com/gridworkz/kato/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

//Config
type Config struct {
	CrtName, KeyName  string
	Address           []string
	IsCa              bool
	CAName, CAKeyName string
	Domains           []string
}

func main() {
	App := cli.NewApp()
	App.Version = version.GetVersion()
	App.Commands = []cli.Command{
		cli.Command{
			Name: "create",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "crt-name",
					Value: "",
					Usage: "creat crt file name",
				},
				cli.StringFlag{
					Name:  "crt-key-name",
					Value: "",
					Usage: "creat crt key file name",
				},
				cli.StringSliceFlag{
					Name:  "address",
					Value: &cli.StringSlice{"127.0.0.1"},
					Usage: "address list",
				},
				cli.StringSliceFlag{
					Name:  "domains",
					Value: &cli.StringSlice{""},
					Usage: "domain list",
				},
				cli.StringFlag{
					Name:  "ca-name",
					Value: "./ca.pem",
					Usage: "creat or read ca file name",
				},
				cli.StringFlag{
					Name:  "ca-key-name",
					Value: "./ca.key.pem",
					Usage: "creat or read ca key file name",
				},
				cli.BoolFlag{
					Name:   "is-ca",
					Hidden: false,
					Usage:  "is create ca",
				},
			},
			Action: create,
		},
	}
	sort.Sort(cli.FlagsByName(App.Flags))
	sort.Sort(cli.CommandsByName(App.Commands))
	App.Run(os.Args)
}
func parseConfig(ctx *cli.Context) Config {
	var c Config
	c.Address = ctx.StringSlice("address")
	c.CAKeyName = ctx.String("ca-key-name")
	c.CAName = ctx.String("ca-name")
	c.CrtName = ctx.String("crt-name")
	c.KeyName = ctx.String("crt-key-name")
	c.Domains = ctx.StringSlice("domains")
	c.IsCa = ctx.Bool("is-ca")
	return c
}
func create(ctx *cli.Context) error {
	c := parseConfig(ctx)
	info := c.CreateCertInformation()
	if c.IsCa {
		err := CreateCRT(nil, nil, info)
		if err != nil {
			logrus.Fatal("Create crt error,Error info:", err)
		}
	} else {
		info.Names = []pkix.AttributeTypeAndValue{
			pkix.AttributeTypeAndValue{
				Type:  asn1.ObjectIdentifier{2, 1, 3},
				Value: "MAC_ADDR",
			},
		}
		crt, pri, err := Parse(c.CAName, c.CAKeyName)
		if err != nil {
			logrus.Fatal("Parse crt error,Error info:", err)
		}
		err = CreateCRT(crt, pri, info)
		if err != nil {
			logrus.Fatal("Create crt error,Error info:", err)
		}
	}
	fmt.Println("create success")
	return nil
}

//CreateCertInformation
func (c *Config) CreateCertInformation() CertInformation {
	baseinfo := CertInformation{
		Country:            []string{"CA"},
		Organization:       []string{"Gridworkz"},
		IsCA:               c.IsCa,
		OrganizationalUnit: []string{"gridworkz kato"},
		EmailAddress:       []string{"gdevs@gridworkz.com"},
		Locality:           []string{"Toronto"},
		Province:           []string{"Ontario"},
		CommonName:         "kato",
		CrtName:            c.CrtName,
		KeyName:            c.KeyName,
		Domains:            c.Domains,
	}
	if c.IsCa {
		baseinfo.CrtName = c.CAName
		baseinfo.KeyName = c.CAKeyName
	}
	var address []net.IP
	for _, a := range c.Address {
		address = append(address, net.ParseIP(a))
	}
	baseinfo.IPAddresses = address
	return baseinfo
}
