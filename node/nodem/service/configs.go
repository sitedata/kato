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

package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gridworkz/kato/util"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

//LoadServicesFromLocal load all service config from config file
func LoadServicesFromLocal(serviceListFile string) []*Service {
	var serviceList []*Service
	ok, err := util.IsDir(serviceListFile)
	if err != nil {
		logrus.Errorf("read service config file error,%s", err.Error())
		return nil
	}
	if !ok {
		services, err := loadServicesFromFile(serviceListFile)
		if err != nil {
			logrus.Errorf("read service config file %s error,%s", serviceListFile, err.Error())
			return nil
		}
		return services.Services
	}
	filepath.Walk(serviceListFile, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "yaml") && !info.IsDir() {
			services, err := loadServicesFromFile(path)
			if err != nil {
				logrus.Errorf("read service config file %s error,%s", path, err.Error())
				return nil
			}
			serviceList = append(serviceList, services.Services...)
		}
		return nil
	})
	result := removeRepByLoop(serviceList)
	logrus.Infof("load service config file success. load %d service", len(result))
	return result
}

//LoadServicesWithFileFromLocal
func LoadServicesWithFileFromLocal(serviceListFile string) []*Services {
	var serviceList []*Services
	ok, err := util.IsDir(serviceListFile)
	if err != nil {
		logrus.Errorf("read service config file error,%s", err.Error())
		return nil
	}
	if !ok {
		services, err := loadServicesFromFile(serviceListFile)
		if err != nil {
			logrus.Errorf("read service config file %s error,%s", serviceListFile, err.Error())
			return nil
		}
		serviceList = append(serviceList, services)
		return serviceList
	}
	filepath.Walk(serviceListFile, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "yaml") && !info.IsDir() {
			services, err := loadServicesFromFile(path)
			if err != nil {
				logrus.Errorf("read service config file %s error,%s", path, err.Error())
				return nil
			}
			serviceList = append(serviceList, services)
		}
		return nil
	})
	return serviceList
}

//WriteServicesWithFile write services config file
func WriteServicesWithFile(serviceLists ...*Services) error {
	for _, serviceList := range serviceLists {
		if serviceList.FromFile != "" {
			write := func() error {
				success := false
				if ok, _ := util.FileExists(serviceList.FromFile); ok {
					if err := os.Rename(serviceList.FromFile, serviceList.FromFile+".bak"); err == nil {
						defer func() {
							if success {
								os.Remove(serviceList.FromFile + ".bak")
							} else {
								os.Rename(serviceList.FromFile+".bak", serviceList.FromFile)
							}
						}()
					}
				}
				file, err := os.OpenFile(serviceList.FromFile, os.O_CREATE|os.O_RDWR, 0755)
				if err != nil {
					return err
				}
				defer file.Close()
				filebody, err := yaml.Marshal(serviceList)
				if err != nil {
					return err
				}
				if _, err := file.Write(filebody); err != nil {
					return err
				}
				success = true
				return nil
			}
			if err := write(); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeRepByLoop(source []*Service) (target []*Service) {
	for i, s := range source {
		flag := true
		for _, t := range target {
			if s.Name == t.Name {
				flag = false
				break
			}
		}
		if flag {
			target = append(target, source[i])
		}
	}
	return
}

func loadServicesFromFile(serviceListFile string) (*Services, error) {
	// load default-configs.yaml
	content, err := ioutil.ReadFile(serviceListFile)
	if err != nil {
		err = fmt.Errorf("Failed to read service list file: %s", err.Error())
		return nil, err
	}
	var defaultConfigs Services
	err = yaml.Unmarshal(content, &defaultConfigs)
	if err != nil {
		logrus.Error("Failed to parse default configs yaml file: ", err)
		return nil, err
	}
	defaultConfigs.FromFile = serviceListFile
	return &defaultConfigs, nil
}
