// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

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
	"bytes"
	"fmt"
	"os/exec"

	"github.com/urfave/cli"
)

//NewCmdDomain domain cmd
//v5.2 need refactoring
func NewCmdDomain() cli.Command {
	c := cli.Command{
		Name: "domain",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "ip",
				Usage: "ip address",
			},
			cli.StringFlag{
				Name:  "domain",
				Usage: "domain",
			},
		},
		Usage: "Default *.grapps.ca domain resolution",
		Action: func(c *cli.Context) error {
			ip := c.String("ip")
			if len(ip) == 0 {
				fmt.Println("ip must not null")
				return nil
			}
			domain := c.String("domain")
			cmd := exec.Command("bash", "/opt/kato/bin/.domain.sh", ip, domain)
			outbuf := bytes.NewBuffer(nil)
			cmd.Stdout = outbuf
			cmd.Run()
			out := outbuf.String()
			fmt.Println(out)
			return nil
		},
	}
	return c
}
