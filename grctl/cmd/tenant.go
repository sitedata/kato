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

	"github.com/gridworkz/kato/grctl/clients"
	"github.com/gridworkz/kato/util/termtables"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	//"github.com/gridworkz/kato/eventlog/conf"
	"errors"

	config "github.com/gridworkz/kato/cmd/grctl/option"
)

//NewCmdTenant tenant cmd
func NewCmdTenant() cli.Command {
	c := cli.Command{
		Name:  "tenant",
		Usage: "grctl tenant -h",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "list",
				Usage: "list all tenant info",
				Action: func(c *cli.Context) error {
					Common(c)
					return getAllTenant(c)
				},
			},
			cli.Command{
				Name:  "get",
				Usage: "get all app details by specified tenant name",
				Action: func(c *cli.Context) error {
					Common(c)
					return getTenantInfo(c)
				},
			},
			cli.Command{
				Name:  "res",
				Usage: "get tenant resource details by specified tenant name",
				Action: func(c *cli.Context) error {
					Common(c)
					return findTenantResourceUsage(c)
				},
			},
			cli.Command{
				Name:  "batchstop",
				Usage: "batch stop app by specified tenant name",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "f",
						Usage: "Continuous log output",
					},
					cli.StringFlag{
						Name:  "event_log_server",
						Usage: "event log server address",
					},
				},
				Action: func(c *cli.Context) error {
					Common(c)
					return stopTenantService(c)
				},
			},
			cli.Command{
				Name:  "setdefname",
				Usage: "set default tenant name",
				Action: func(c *cli.Context) error {
					err := CreateTenantFile(c.Args().First())
					if err != nil {
						logrus.Error("set default tenantname fail", err.Error())
					}
					return nil
				},
			},
		},
	}
	return c
}

// grctrl tenant TENANT_NAME
func getTenantInfo(c *cli.Context) error {
	tenantID := c.Args().First()
	if tenantID == "" {
		fmt.Println("Please provide tenant name")
		os.Exit(1)
	}
	services, err := clients.RegionClient.Tenants(tenantID).Services("").List()
	handleErr(err)
	if services != nil {
		runtable := termtables.CreateTable()
		closedtable := termtables.CreateTable()
		runtable.AddHeaders("Service alias", "Application status", "Deploy version", "Number of instances", "Memory footprint")
		closedtable.AddHeaders("Tenant ID", "Service ID", "Service alias", "Application status", "Deploy version")
		for _, service := range services {
			if service.CurStatus != "closed" && service.CurStatus != "closing" && service.CurStatus != "undeploy" && service.CurStatus != "deploying" {
				runtable.AddRow(service.ServiceAlias, service.CurStatus, service.DeployVersion, service.Replicas, fmt.Sprintf("%d Mb", service.ContainerMemory*service.Replicas))
			} else {
				closedtable.AddRow(service.TenantID, service.ServiceID, service.ServiceAlias, service.CurStatus, service.DeployVersion)
			}
		}
		fmt.Println("Running application:")
		fmt.Println(runtable.Render())
		fmt.Println("Applications that are not running:")
		fmt.Println(closedtable.Render())
		return nil
	}
	return nil
}
func findTenantResourceUsage(c *cli.Context) error {
	tenantName := c.Args().First()
	if tenantName == "" {
		fmt.Println("Please provide tenant name")
		os.Exit(1)
	}
	resources, err := clients.RegionClient.Resources().Tenants(tenantName).Get()
	handleErr(err)
	table := uitable.New()
	table.Wrap = true // wrap columns
	table.AddRow("Tenant name:", resources.Name)
	table.AddRow("Tenant ID:", resources.UUID)
	table.AddRow("Enterprise ID:", resources.EID)
	table.AddRow("Using CPU resources:", fmt.Sprintf("%.2f Core", float64(resources.UsedCPU)/1000))
	table.AddRow("Using memory resources:", fmt.Sprintf("%d %s", resources.UsedMEM, "Mb"))
	table.AddRow("Disk resources are being used:", fmt.Sprintf("%.2f Mb", resources.UsedDisk/1024))
	table.AddRow("Total allocated CPU resources:", fmt.Sprintf("%.2f Core", float64(resources.AllocatedCPU)/1000))
	table.AddRow("Total allocated memory resources:", fmt.Sprintf("%d %s", resources.AllocatedMEM, "Mb"))
	fmt.Println(table)
	return nil
}

func getAllTenant(c *cli.Context) error {
	tenants, err := clients.RegionClient.Tenants("").List()
	handleErr(err)
	tenantsTable := termtables.CreateTable()
	tenantsTable.AddHeaders("TenantAlias", "TenantID", "TenantLimit")
	for _, t := range tenants {
		tenantsTable.AddRow(t.Name, t.UUID, fmt.Sprintf("%d GB", t.LimitMemory))
	}
	fmt.Print(tenantsTable.Render())
	return nil
}

//CreateTenantFile Create Tenant File
func CreateTenantFile(tname string) error {
	filename, err := config.GetTenantNamePath()
	if err != nil {
		logrus.Warn("Load config file error.")
		return errors.New("Load config file error")
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logrus.Warn("load teantnamefile file", err.Error())
		f.Close()
		return err
	}
	_, err = f.WriteString(tname)
	if err != nil {
		logrus.Warn("write teantnamefile file", err.Error())
	}
	f.Close()
	return err
}
