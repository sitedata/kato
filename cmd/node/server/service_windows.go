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

package server

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gridworkz/kato/cmd/node/option"
	utilwindows "github.com/gridworkz/kato/util/windows"
	"github.com/spf13/pflag"
	"golang.org/x/sys/windows"
)

var (
	flRegisterService   *bool
	flUnregisterService *bool
	flServiceName       *string
	flRunService        *bool

	setStdHandle = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetStdHandle")
	oldStderr    windows.Handle
	panicFile    *os.File
)

//InstallServiceFlags install service flag set
func InstallServiceFlags(flags *pflag.FlagSet) {
	flServiceName = flags.String("service-name", "kato-node", "Set the Windows service name")
	flRegisterService = flags.Bool("register-service", false, "Register the service and exit")
	flUnregisterService = flags.Bool("unregister-service", false, "Unregister the service and exit")
	flRunService = flags.Bool("run-service", false, "")
	flags.MarkHidden("run-service")
}
func getServicePath() (string, error) {
	p, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(p)
}

// initService is the entry point for running the daemon as a Windows
// service. It returns an indication to stop (if registering/un-registering);
// an indication of whether it is running as a service; and an error.
func initService(conf *option.Conf, startfunc, stopfunc func() error) error {
	if *flUnregisterService {
		if *flRegisterService {
			return errors.New("--register-service and --unregister-service cannot be used together")
		}
		return unregisterService()
	}

	if *flRegisterService {
		return registerService()
	}
	if !*flRunService {
		return startfunc()
	}
	return utilwindows.RunAsService(*flServiceName, startfunc, stopfunc, false)
}

func unregisterService() error {
	if err := utilwindows.StopService(*flServiceName); err != nil && !strings.Contains(err.Error(), "service has not been started") {
		return err
	}
	return utilwindows.UnRegisterService(*flServiceName)
}

func registerService() error {
	p, err := getServicePath()
	if err != nil {
		return err
	}
	// Configure the service to launch with the arguments that were just passed.
	args := []string{"--run-service"}
	for _, a := range os.Args[1:] {
		if a != "--register-service" && a != "--unregister-service" {
			args = append(args, a)
		}
	}
	return utilwindows.RegisterService(*flServiceName, p, "Kato NodeManager", []string{}, args)
}
