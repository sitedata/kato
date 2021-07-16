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
	"context"
	"fmt"
	"strings"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoyv2 "github.com/gridworkz/kato/node/core/envoy/v2"
	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

//NewCmdEnvoy
func NewCmdEnvoy() cli.Command {
	c := cli.Command{
		Name:  "envoy",
		Usage: "envoy management related commands",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "address",
				Usage: "node envoy api address",
				Value: "127.0.0.1:6101",
			},
			cli.StringFlag{
				Name:  "node",
				Usage: "envoy node name",
			},
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "endpoints",
				Usage: "list envoy node endpoints",
				Action: func(c *cli.Context) error {
					return listEnvoyEndpoint(c)
				},
			},
		},
	}
	return c
}

func listEnvoyEndpoint(c *cli.Context) error {
	if c.GlobalString("node") == "" {
		showError("node name can not be empty,please define by --node")
	}
	cli, err := grpc.Dial(c.GlobalString("address"), grpc.WithInsecure())
	if err != nil {
		showError(err.Error())
	}
	endpointDiscover := v2.NewEndpointDiscoveryServiceClient(cli)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := endpointDiscover.FetchEndpoints(ctx, &v2.DiscoveryRequest{
		Node: &core.Node{
			Cluster: c.GlobalString("node"),
			Id:      c.GlobalString("node"),
		},
	})
	if err != nil {
		showError(err.Error())
	}
	if len(res.Resources) == 0 {
		showError("not find endpoints")
	}
	endpoints := envoyv2.ParseLocalityLbEndpointsResource(res.Resources)
	table := uitable.New()
	table.Wrap = true // wrap columns
	for _, end := range endpoints {
		table.AddRow(end.ClusterName, strings.Join(func() []string {
			var re []string
			for _, e := range end.Endpoints {
				for _, a := range e.LbEndpoints {
					if lbe, ok := a.HostIdentifier.(*endpoint.LbEndpoint_Endpoint); ok && lbe != nil {
						if address, ok := lbe.Endpoint.Address.Address.(*core.Address_SocketAddress); ok && address != nil {
							if port, ok := address.SocketAddress.PortSpecifier.(*core.SocketAddress_PortValue); ok && port != nil {
								re = append(re, fmt.Sprintf("%s:%d", address.SocketAddress.Address, port.PortValue))
							}
						}
					}
				}
			}
			return re
		}(), ";"))
	}
	fmt.Println(table)
	return nil
}
