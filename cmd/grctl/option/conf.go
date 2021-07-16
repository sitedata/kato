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
	"io/ioutil"
	"os"
	"path"

	"github.com/gridworkz/kato/api/region"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
	//"strings"
)

var config Config

//Config
type Config struct {
	Kubernets     Kubernets      `yaml:"kube"`
	RegionAPI     region.APIConf `yaml:"region_api"`
	DockerLogPath string         `yaml:"docker_log_path"`
}

//RegionMysql
type RegionMysql struct {
	URL      string `yaml:"url"`
	Pass     string `yaml:"pass"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
}

//Kubernetes
type Kubernets struct {
	KubeConf string `yaml:"kube-conf"`
}

//LoadConfig
func LoadConfig(ctx *cli.Context) (Config, error) {
	config = Config{
		RegionAPI: region.APIConf{
			Endpoints: []string{"http://127.0.0.1:8888"},
		},
	}
	configfile := ctx.GlobalString("config")
	if configfile == "" {
		home, _ := sources.Home()
		configfile = path.Join(home, ".rbd", "grctl.yaml")
	}
	_, err := os.Stat(configfile)
	if err != nil {
		return config, nil
	}
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		logrus.Warning("Read config file error ,will get config from region.", err.Error())
		return config, err
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		logrus.Warning("Read config file error ,will get config from region.", err.Error())
		return config, err
	}
	return config, nil
}

//GetConfig
func GetConfig() Config {
	return config
}

// Get tenantNamePath
func GetTenantNamePath() (tenantnamepath string, err error) {
	home, err := sources.Home()
	if err != nil {
		logrus.Warn("Get Home Dir error.", err.Error())
		return tenantnamepath, err
	}
	tenantnamepath = path.Join(home, ".rbd", "tenant.txt")
	return tenantnamepath, err
}
