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

package compose

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/docker/libcompose/project"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

// Compose is docker compose file loader, implements Loader interface
type Compose struct {
}

// checkUnsupportedKey checks if libcompose project contains
// keys that are not supported by this loader.
// list of all unsupported keys are stored in unsupportedKey variable
// returns list of unsupported YAML keys from docker-compose
func checkUnsupportedKey(composeProject *project.Project) []string {

	// list of all unsupported keys for this loader
	// this is map to make searching for keys easier
	// to make sure that unsupported key is not going to be reported twice
	// by keeping record if already saw this key in another service
	var unsupportedKey = map[string]bool{
		"CgroupParent":  false,
		"CPUSet":        false,
		"CPUShares":     false,
		"Devices":       false,
		"DependsOn":     true,
		"DNS":           false,
		"DNSSearch":     false,
		"DomainName":    false,
		"EnvFile":       false,
		"ExternalLinks": false,
		"ExtraHosts":    false,
		"Hostname":      false,
		"Ipc":           false,
		"Logging":       false,
		"MacAddress":    false,
		"MemSwapLimit":  false,
		"NetworkMode":   false,
		"SecurityOpt":   false,
		"ShmSize":       false,
		"StopSignal":    false,
		"VolumeDriver":  false,
		"Uts":           false,
		"ReadOnly":      false,
		"Ulimits":       false,
		"Net":           false,
		"Sysctls":       false,
		"Networks":      false, // there are special checks for Network in checkUnsupportedKey function
		"Links":         true,
	}

	// collect all keys found in project
	var keysFound []string

	// Root level keys are not yet supported
	// Check to see if the default network is available and length is only equal to one.
	// Else, warn the user that root level networks are not supported (yet)
	if _, ok := composeProject.NetworkConfigs["default"]; ok && len(composeProject.NetworkConfigs) == 1 {
		logrus.Debug("Default network found")
	} else if len(composeProject.NetworkConfigs) > 0 {
		keysFound = append(keysFound, "root level networks")
	}

	// Root level volumes are not yet supported
	if len(composeProject.VolumeConfigs) > 0 {
		keysFound = append(keysFound, "root level volumes")
	}

	for _, serviceConfig := range composeProject.ServiceConfigs.All() {
		// this reflection is used in check for empty arrays
		val := reflect.ValueOf(serviceConfig).Elem()
		s := structs.New(serviceConfig)

		for _, f := range s.Fields() {
			// Check if given key is among unsupported keys, and skip it if we already saw this key
			if alreadySaw, ok := unsupportedKey[f.Name()]; ok && !alreadySaw {
				if f.IsExported() && !f.IsZero() {
					// IsZero returns false for empty array/slice ([])
					// this check if field is Slice, and then it checks its size
					if field := val.FieldByName(f.Name()); field.Kind() == reflect.Slice {
						if field.Len() == 0 {
							// array is empty it doesn't matter if it is in unsupportedKey or not
							continue
						}
					}
					//get yaml tag name instad of variable name
					yamlTagName := strings.Split(f.Tag("yaml"), ",")[0]
					if f.Name() == "Networks" {
						// networks always contains one default element, even it isn't declared in compose v2.
						if len(serviceConfig.Networks.Networks) == 1 && serviceConfig.Networks.Networks[0].Name == "default" {
							// this is empty Network definition, skip it
							continue
						} else {
							yamlTagName = "networks"
						}
					}

					if linksArray := val.FieldByName(f.Name()); f.Name() == "Links" && linksArray.Kind() == reflect.Slice {
						//Links has "SERVICE:ALIAS" style, we don't support SERVICE != ALIAS
						findUnsupportedLinksFlag := false
						for i := 0; i < linksArray.Len(); i++ {
							if tmpLink := linksArray.Index(i); tmpLink.Kind() == reflect.String {
								tmpLinkStr := tmpLink.String()
								tmpLinkStrSplit := strings.Split(tmpLinkStr, ":")
								if len(tmpLinkStrSplit) == 2 && tmpLinkStrSplit[0] != tmpLinkStrSplit[1] {
									findUnsupportedLinksFlag = true
									break
								}
							}
						}
						if !findUnsupportedLinksFlag {
							continue
						}

					}

					keysFound = append(keysFound, yamlTagName)
					unsupportedKey[f.Name()] = true
				}
			}
		}
	}
	return keysFound
}

// LoadBytes loads a compose file byte into KomposeObject
func (c *Compose) LoadBytes(bodys [][]byte) (ComposeObject, error) {

	// Load the json / yaml file in order to get the version value
	var version string

	for _, body := range bodys {
		composeVersion, err := getVersionFromByte(body)
		if err != nil {
			return ComposeObject{}, fmt.Errorf("Unable to load yaml/json file for version parsing,%s", err.Error())
		}

		// Check that the previous file loaded matches.
		if len(bodys) > 0 && version != "" && version != composeVersion {
			return ComposeObject{}, fmt.Errorf("All Docker Compose files must be of the same version")
		}
		version = composeVersion
	}

	logrus.Debugf("Docker Compose version: %s", version)

	// Convert based on version
	switch version {
	// Use libcompose for 1 or 2
	// If blank, it's assumed it's 1 or 2
	case "", "1", "1.0", "2", "2.0", "2.1", "2.2", "2.3", "2.4":
		co, err := parseV1V2(bodys)
		if err != nil {
			return ComposeObject{}, err
		}
		return co, nil
	// Use docker/cli for 3
	case "3", "3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "3.7":
		co, err := parseV3(bodys)
		if err != nil {
			return ComposeObject{}, err
		}
		return co, nil
	default:
		return ComposeObject{}, fmt.Errorf("Version %s of Docker Compose is not supported. Please use version 1, 2 or 3", version)
	}

}

func getVersionFromFile(file string) (string, error) {
	type ComposeVersion struct {
		Version string `json:"version"` // This affects YAML as well
	}
	var version ComposeVersion

	loadedFile, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(loadedFile, &version)
	if err != nil {
		return "", err
	}

	return version.Version, nil
}

func getVersionFromByte(body []byte) (string, error) {
	type ComposeVersion struct {
		Version string `json:"version"` // This affects YAML as well
	}
	var version ComposeVersion

	err := yaml.Unmarshal(body, &version)
	if err != nil {
		return "", err
	}
	return version.Version, nil
}
