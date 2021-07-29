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
	"os"

	"github.com/gridworkz/kato/grctl/clients"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NewCmdConfig config command
func NewCmdConfig() cli.Command {
	c := cli.Command{
		Name: "config",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "output,o",
				Usage: "write region api config to file",
			},
			cli.StringFlag{
				Name:  "namespace,ns",
				Usage: "kato default namespace",
				Value: "rbd-system",
			},
		},
		Usage: "show region config file",
		Action: func(c *cli.Context) {
			Common(c)
			namespace := c.String("namespace")
			configMap, err := clients.K8SClient.CoreV1().ConfigMaps(namespace).Get(context.Background(), "region-config", metav1.GetOptions{})
			if err != nil {
				showError(err.Error())
			}
			regionConfig := map[string]string{
				"client.pem":          string(configMap.BinaryData["client.pem"]),
				"client.key.pem":      string(configMap.BinaryData["client.key.pem"]),
				"ca.pem":              string(configMap.BinaryData["ca.pem"]),
				"apiAddress":          configMap.Data["apiAddress"],
				"websocketAddress":    configMap.Data["websocketAddress"],
				"defaultDomainSuffix": configMap.Data["defaultDomainSuffix"],
				"defaultTCPHost":      configMap.Data["defaultTCPHost"],
			}
			body, err := yaml.Marshal(regionConfig)
			if err != nil {
				showError(err.Error())
			}
			if c.String("o") != "" {
				file, err := os.OpenFile(c.String("o"), os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					showError(err.Error())
				}
				defer file.Close()
				_, err = file.Write(body)
				if err != nil {
					showError(err.Error())
				}
			} else {
				fmt.Println(string(body))
			}
		},
	}
	return c
}
