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

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gridworkz/kato/gateway/controller/openresty/model"
	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

//NewCmdGateway
func NewCmdGateway() cli.Command {
	c := cli.Command{
		Name:  "gateway",
		Usage: "Gateway management related commands",
		Subcommands: []cli.Command{
			{
				Name:  "endpoints",
				Usage: "list gateway http endpoints",
				Subcommands: []cli.Command{
					{
						Name:  "http",
						Usage: "list gateway http endpoints",
						Flags: []cli.Flag{
							cli.IntFlag{
								Name:  "port",
								Usage: "gateway http endpoint query port",
								Value: 18080,
							},
						},
						Action: func(c *cli.Context) error {
							return listHTTPEndpoint(c)
						},
					},
					{
						Name:  "stream",
						Usage: "list gateway stream endpoints",
						Flags: []cli.Flag{
							cli.IntFlag{
								Name:  "port",
								Usage: "gateway stream endpoint query port",
								Value: 18081,
							},
						},
						Action: func(c *cli.Context) error {
							return listStreamEndpoint(c)
						},
					},
				},
			},
		},
	}
	return c
}

func listStreamEndpoint(c *cli.Context) error {
	return tcpGetAndPrint(fmt.Sprintf("127.0.0.1:%d", c.Int("port")))
}

func listHTTPEndpoint(c *cli.Context) error {
	return httpGetAndPrint(fmt.Sprintf("http://127.0.0.1:%d/config/backends", c.Int("port")))
}
func tcpGetAndPrint(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("GET\r\n"))
	if err != nil {
		return err
	}
	print(conn)
	return nil
}
func httpGetAndPrint(url string) error {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if res.Body != nil {
		defer res.Body.Close()
		print(res.Body)
	}
	return nil
}

func print(reader io.Reader) {
	decoder := json.NewDecoder(reader)
	var backends []*model.Backend
	if err := decoder.Decode(&backends); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	table := uitable.New()
	table.Wrap = true // wrap columns
	for _, b := range backends {
		table.AddRow(b.Name, strings.Join(func() []string {
			var re []string
			for _, e := range b.Endpoints {
				re = append(re, fmt.Sprintf("%s:%s %d", e.Address, e.Port, e.Weight))
			}
			return re
		}(), ";"))
	}
	fmt.Println(table)
}
