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

package ansible

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gridworkz/kato/util"
)

//NodeInstallOption
type NodeInstallOption struct {
	HostRole   string
	HostName   string
	InternalIP string
	RootPass   string // ssh login password
	KeyPath    string // ssh login key path
	NodeID     string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	loginValue string
	linkModel  string
}

//RunNodeInstallCmd - install node
func RunNodeInstallCmd(option NodeInstallOption) (err error) {
	installNodeShellPath := os.Getenv("INSTALL_NODE_SHELL_PATH")
	if installNodeShellPath == "" {
		installNodeShellPath = "/opt/kato/kato-ansible/scripts/node.sh"
	}
	// ansible file must exists
	if ok, _ := util.FileExists(installNodeShellPath); !ok {
		return fmt.Errorf("install node scripts is not found")
	}
	// ansible's param can't send nil nor empty string
	if err := preCheckNodeInstall(&option); err != nil {
		return err
	}
	line := fmt.Sprintf("'%s' -r '%s' -i '%s' -t '%s' -k '%s' -u '%s'",
		installNodeShellPath, option.HostRole, option.InternalIP, option.linkModel, option.loginValue, option.NodeID)
	cmd := exec.Command("bash", "-c", line)
	cmd.Stdin = option.Stdin
	cmd.Stdout = option.Stdout
	cmd.Stderr = option.Stderr
	return cmd.Run()
}

// check param
func preCheckNodeInstall(option *NodeInstallOption) error {
	if strings.TrimSpace(option.HostRole) == "" {
		return fmt.Errorf("install node failed, install scripts needs param hostRole")
	}
	if strings.TrimSpace(option.InternalIP) == "" {
		return fmt.Errorf("install node failed, install scripts needs param internalIP")
	}
	//login key path first, and then rootPass, so keyPath and RootPass can't all be empty
	if strings.TrimSpace(option.KeyPath) == "" {
		if strings.TrimSpace(option.RootPass) == "" {
			return fmt.Errorf("install node failed, install scripts needs login key path or login password")
		}
		option.loginValue = strings.TrimSpace(option.RootPass)
		option.linkModel = "pass"
	} else {
		option.loginValue = strings.TrimSpace(option.KeyPath)
		option.linkModel = "key"
	}
	if strings.TrimSpace(option.NodeID) == "" {
		return fmt.Errorf("install node failed, install scripts needs param nodeID")
	}
	return nil
}
