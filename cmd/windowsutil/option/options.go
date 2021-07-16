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

package option

import (
	"fmt"
	"os"
	"path"

	"github.com/gridworkz/kato/util"

	"github.com/gridworkz/kato/util/windows"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

//Config
type Config struct {
	Debug        bool
	RunShell     string
	ServiceName  string
	RunAsService bool
	LogFile      string
}

var removeService bool

//AddFlags config
func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.RunShell, "run", "", "Specify startup command")
	fs.StringVar(&c.ServiceName, "service-name", "", "Specify windows service name")
	fs.StringVar(&c.LogFile, "log-file", "c:\\windwosutil.log", "service log outputfile")
	fs.BoolVar(&c.RunAsService, "run-as-service", true, "run as windows service")
	fs.BoolVar(&c.Debug, "debug", false, "debug mode run ")
	fs.BoolVar(&removeService, "remove-service", false, "remove windows service")
}

//Check config
func (c *Config) Check() bool {
	if c.ServiceName == "" {
		logrus.Errorf("service name can not be empty")
		return false
	}
	if c.RunShell == "" && !removeService {
		logrus.Errorf("run shell can not be empty")
		return false
	}
	if err := util.CheckAndCreateDir(path.Dir(c.LogFile)); err != nil {
		logrus.Errorf("create node log file dir failure %s", err.Error())
		os.Exit(1)
	}
	logfile, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		logrus.Fatalf("open log file %s failure %s", c.LogFile, err.Error())
	}
	logrus.SetOutput(logfile)
	if removeService {
		if err := windows.UnRegisterService(c.ServiceName); err != nil {
			fmt.Printf("remove service %s failure %s", c.ServiceName, err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	return true
}
