package cmd

import (
	"fmt"

	licutil "github.com/gridworkz/kato/util/license"
	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

//NewCmdLicense -
func NewCmdLicense() cli.Command {
	c := cli.Command{
		Name:  "license",
		Usage: "kato license manage cmd",
		Subcommands: []cli.Command{
			{
				Name:  "show",
				Usage: "show license information",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "lic-path, lp",
						Usage: "license file path",
						Value: "/opt/kato/etc/license/license.yb",
					},
					cli.StringFlag{
						Name:  "lic-so-path, lsp",
						Usage: "license.so file path",
						Value: "/opt/kato/etc/license/license.so",
					},
				},
				Action: func(c *cli.Context) error {
					Common(c)
					licPath := c.String("lic-path")
					licSoPath := c.String("lic-so-path")
					licInfo, err := licutil.GetLicInfo(licPath, licSoPath)
					if err != nil {
						showError(err.Error())
					}

					if licInfo == nil {
						fmt.Println("non-enterprise version, no license information")
						return nil
					}

					table := uitable.New()
					table.AddRow("Authorized company name:", licInfo.Company)
					table.AddRow("Authorized company code:", licInfo.Code)
					table.AddRow("Number of authorized single data center nodes:", licInfo.Node)
					table.AddRow("Authorization start time:", licInfo.StartTime)
					table.AddRow("Authorization expiration time:", licInfo.EndTime)
					table.AddRow("Authorization key:", licInfo.LicKey)
					fmt.Println(table)
					return nil
				},
			},
			{
				Name:  "genkey",
				Usage: "generate a license key for the machine",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "lic-so-path, lsp",
						Usage: "license.so file path",
						Value: "/opt/kato/etc/license/license.so",
					},
				},
				Action: func(c *cli.Context) error {
					Common(c)
					licSoPath := c.String("lic-so-path")
					licKey, err := licutil.GenLicKey(licSoPath)
					if err != nil {
						showError(err.Error())
					}

					if licKey == "" {
						fmt.Println("non-enterprise version, no license key")
						return nil
					}
					fmt.Println(licKey)
					return nil
				},
			},
		},
	}
	return c
}
