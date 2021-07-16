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
	"github.com/gridworkz/kato/grctl/clients"
	"github.com/gridworkz/kato/node/nodem/client"
	"github.com/urfave/cli"

	ansibleUtil "github.com/gridworkz/kato/util/ansible"
)

//NewCmdAnsible ansible config cmd
func NewCmdAnsible() cli.Command {
	c := cli.Command{
		Name:   "ansible",
		Usage:  "Manage the ansible environment",
		Hidden: true,
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "hosts",
				Usage: "Manage the ansible hosts config environment",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "hosts-file-path",
						Usage: "hosts file path",
						Value: "/opt/kato/kato-ansible/inventory/hosts",
					},
					cli.StringFlag{
						Name:  "config-file-path",
						Usage: "install config path",
						Value: "/opt/kato/kato-ansible/scripts/installer/global.sh",
					},
				},
				Action: func(c *cli.Context) error {
					Common(c)
					hosts, err := clients.RegionClient.Nodes().List()
					handleErr(err)
					return WriteHostsFile(c.String("hosts-file-path"), c.String("config-file-path"), hosts)
				},
			},
		},
	}
	return c
}

//WriteHostsFile write hosts file
func WriteHostsFile(filePath, installConfPath string, hosts []*client.HostNode) error {
	//get node list from api without condition list.
	//so will get condition
	for i := range hosts {
		nodeWithCondition, _ := clients.RegionClient.Nodes().Get(hosts[i].ID)
		if nodeWithCondition != nil {
			hosts[i] = nodeWithCondition
		}
	}
	return ansibleUtil.WriteHostsFile(filePath, installConfPath, hosts)
}
