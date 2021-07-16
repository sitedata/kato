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

package nginxcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	nginxBinary      = "nginx"
	defaultNginxConf = "/run/nginx/conf/nginx.conf"
	//ErrorCheck check config file failure
	ErrorCheck  = fmt.Errorf("error check config")
	updateCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nginx",
		Subsystem: "",
		Name:      "update",
		Help:      "Number of nginx updates inside the gateway",
	})
	errUpdateCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "nginx",
		Subsystem: "",
		Name:      "update_err",
		Help:      "Number of nginx error updates inside the gateway",
	})
)

func init() {
	nginxBinary = path.Join(os.Getenv("OPENRESTY_HOME"), "/nginx/sbin/nginx")
	ngx := os.Getenv("NGINX_BINARY")
	if ngx != "" {
		nginxBinary = ngx
	}
}

//SetDefaultNginxConf set
func SetDefaultNginxConf(path string) {
	defaultNginxConf = path
}

//PromethesuScrape prometheus scrape
func PromethesuScrape(ch chan<- *prometheus.Desc) {
	updateCount.Describe(ch)
	errUpdateCount.Describe(ch)
}

//PrometheusCollect prometheus collect
func PrometheusCollect(ch chan<- prometheus.Metric) {
	updateCount.Collect(ch)
	errUpdateCount.Collect(ch)
}

//CreateNginxCommand create nginx command
func CreateNginxCommand(args ...string) *exec.Cmd {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "-c", defaultNginxConf)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(nginxBinary, cmdArgs...)
	return cmd
}

//ExecNginxCommand exec nginx command
func ExecNginxCommand(args ...string) error {
	cmd := CreateNginxCommand(args...)
	if body, err := cmd.Output(); err != nil {
		if eerr, ok := err.(*exec.ExitError); ok {
			logrus.Errorf("nginx exec failure:%s", string(eerr.Stderr))
		}
		if len(body) > 0 {
			logrus.Errorf("nginx exec failure:%s", string(body))
		}
		return err
	}
	return nil
}

//CheckConfig check nginx config file
func CheckConfig() error {
	if err := ExecNginxCommand("-t"); err != nil {
		return ErrorCheck
	}
	return nil
}

//Reload reload nginx config
func Reload() error {
	updateCount.Inc()
	if err := ExecNginxCommand("-s", "reload"); err != nil {
		errUpdateCount.Inc()
		return err
	}
	logrus.Infof("nginx config reload success")
	return nil
}
