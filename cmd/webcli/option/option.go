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

package option

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

//Config server
type Config struct {
	EtcdEndPoints        []string
	EtcdCaFile           string
	EtcdCertFile         string
	EtcdKeyFile          string
	Address              string
	HostIP               string
	HostName             string
	Port                 int
	SessionKey           string
	PrometheusMetricPath string
	K8SConfPath          string
}

//WebCliServer
type WebCliServer struct {
	Config
	LogLevel string
}

//NewWebCliServer
func NewWebCliServer() *WebCliServer {
	return &WebCliServer{}
}

//AddFlags config
func (a *WebCliServer) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&a.LogLevel, "log-level", "info", "the webcli log level")
	fs.StringSliceVar(&a.EtcdEndPoints, "etcd-endpoints", []string{"http://127.0.0.1:2379"}, "etcd v3 cluster endpoints.")
	fs.StringVar(&a.EtcdCaFile, "etcd-ca", "", "etcd tls ca file ")
	fs.StringVar(&a.EtcdCertFile, "etcd-cert", "", "etcd tls cert file")
	fs.StringVar(&a.EtcdKeyFile, "etcd-key", "", "etcd http tls cert key file")
	fs.StringVar(&a.Address, "address", "0.0.0.0", "server listen address")
	fs.StringVar(&a.HostIP, "hostIP", "", "Current node Intranet IP")
	fs.StringVar(&a.HostName, "hostName", "", "Current node host name")
	fs.StringVar(&a.K8SConfPath, "kube-conf", "", "absolute path to the kubeconfig file")
	fs.IntVar(&a.Port, "port", 7171, "server listen port")
	fs.StringVar(&a.PrometheusMetricPath, "metric", "/metrics", "prometheus metrics path")
}

//SetLog
func (a *WebCliServer) SetLog() {
	level, err := logrus.ParseLevel(a.LogLevel)
	if err != nil {
		fmt.Println("set log level error." + err.Error())
		return
	}
	logrus.SetLevel(level)
}
