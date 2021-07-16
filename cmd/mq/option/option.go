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

import "github.com/spf13/pflag"
import "github.com/sirupsen/logrus"
import "fmt"

//Config server
type Config struct {
	EtcdEndPoints        []string
	EtcdCaFile           string
	EtcdCertFile         string
	EtcdKeyFile          string
	EtcdTimeout          int
	EtcdPrefix           string
	ClusterName          string
	APIPort              int
	PrometheusMetricPath string
	RunMode              string //http grpc
	HostIP               string
	HostName             string
}

//MQServer lb worker server
type MQServer struct {
	Config
	LogLevel string
}

//NewMQServer
func NewMQServer() *MQServer {
	return &MQServer{}
}

//AddFlags config
func (a *MQServer) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&a.LogLevel, "log-level", "info", "the mq log level")
	fs.StringSliceVar(&a.EtcdEndPoints, "etcd-endpoints", []string{"http://127.0.0.1:2379"}, "etcd v3 cluster endpoints.")
	fs.IntVar(&a.EtcdTimeout, "etcd-timeout", 10, "etcd http timeout seconds")
	fs.StringVar(&a.EtcdCaFile, "etcd-ca", "", "etcd tls ca file ")
	fs.StringVar(&a.EtcdCertFile, "etcd-cert", "", "etcd tls cert file")
	fs.StringVar(&a.EtcdKeyFile, "etcd-key", "", "etcd http tls cert key file")
	fs.StringVar(&a.EtcdPrefix, "etcd-prefix", "/mq", "the etcd data save key prefix ")
	fs.IntVar(&a.APIPort, "api-port", 6300, "the api server listen port")
	fs.StringVar(&a.RunMode, "mode", "grpc", "the api server run mode grpc or http")
	fs.StringVar(&a.PrometheusMetricPath, "metric", "/metrics", "prometheus metrics path")
	fs.StringVar(&a.HostIP, "hostIP", "", "Current node Intranet IP")
	fs.StringVar(&a.HostName, "hostName", "", "Current node host name")
}

//SetLog
func (a *MQServer) SetLog() {
	level, err := logrus.ParseLevel(a.LogLevel)
	if err != nil {
		fmt.Println("set log level error." + err.Error())
		return
	}
	logrus.SetLevel(level)
}
