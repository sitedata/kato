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

package code

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

//KatoFileConfig - source code configuration file
type KatoFileConfig struct {
	Language  string                 `yaml:"language"`
	BuildPath string                 `yaml:"buildpath"`
	Ports     []Port                 `yaml:"ports"`
	Envs      map[string]interface{} `yaml:"envs"`
	Cmd       string                 `yaml:"cmd"`
	Services  []*Service             `yaml:"services"`
}

// Service contains
type Service struct {
	Name  string            `yaml:"name"`
	Ports []Port            `yaml:"ports"`
	Envs  map[string]string `yaml:"envs"`
}

//Port
type Port struct {
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"`
}

//ReadKatoFile - read cloud help code configuration
func ReadKatoFile(homepath string) (*KatoFileConfig, error) {
	if ok, _ := util.FileExists(path.Join(homepath, "katofile")); !ok {
		return nil, ErrKatoFileNotFound
	}
	body, err := ioutil.ReadFile(path.Join(homepath, "katofile"))
	if err != nil {
		logrus.Error("read kato file error,", err.Error())
		return nil, fmt.Errorf("read kato file error")
	}
	var rbdfile KatoFileConfig
	if err := yaml.Unmarshal(body, &rbdfile); err != nil {
		logrus.Error("marshal kato file error,", err.Error())
		return nil, fmt.Errorf("marshal kato file error")
	}
	return &rbdfile, nil
}
