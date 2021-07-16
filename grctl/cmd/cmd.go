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
	"fmt"
	"os"
	"strings"

	conf "github.com/gridworkz/kato/cmd/grctl/option"
	"github.com/gridworkz/kato/grctl/clients"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

//GetCmds
func GetCmds() []cli.Command {
	cmds := []cli.Command{}
	cmds = append(cmds, NewCmdInstall())
	cmds = append(cmds, NewCmdService())
	cmds = append(cmds, NewCmdTenant())
	cmds = append(cmds, NewCmdNode())
	cmds = append(cmds, NewCmdCluster())
	cmds = append(cmds, NewSourceBuildCmd())
	cmds = append(cmds, NewCmdAnsible())
	cmds = append(cmds, NewCmdLicense())
	cmds = append(cmds, NewCmdGateway())
	cmds = append(cmds, NewCmdEnvoy())
	cmds = append(cmds, NewCmdConfig())
	return cmds
}

//Common
func Common(c *cli.Context) {
	config, err := conf.LoadConfig(c)
	if err != nil {
		logrus.Warn("Load config file error.", err.Error())
	}
	kc := c.GlobalString("kubeconfig")
	if kc != "" {
		config.Kubernets.KubeConf = kc
	}
	if err := clients.InitClient(config.Kubernets.KubeConf); err != nil {
		logrus.Errorf("error config k8s,details %s", err.Error())
	}
	//clients.SetInfo(config.RegionAPI.URL, config.RegionAPI.Token)
	if err := clients.InitRegionClient(config.RegionAPI); err != nil {
		logrus.Fatal("error config region")
	}

}

//CommonWithoutRegion Common
func CommonWithoutRegion(c *cli.Context) {
	config, err := conf.LoadConfig(c)
	if err != nil {
		logrus.Warn("Load config file error.", err.Error())
	}
	kc := c.GlobalString("kubeconfig")
	if kc != "" {
		config.Kubernets.KubeConf = kc
	}
	if err := clients.InitClient(config.Kubernets.KubeConf); err != nil {
		logrus.Errorf("error config k8s,details %s", err.Error())
	}
}

// fatal prints the message (if provided) and then exits. If V(2) or greater,
// glog.Fatal is invoked for extended information.
func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

//GetTenantNamePath
func GetTenantNamePath() string {
	tenantnamepath, err := conf.GetTenantNamePath()
	if err != nil {
		logrus.Warn("Ger Home error", err.Error())
		return tenantnamepath
	}
	return tenantnamepath
}
