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
	"strings"

	simplejson "github.com/bitly/go-simplejson"

	"github.com/gridworkz/kato/util"
)

//ErrRuntimeNotSupport runtime not supported
var ErrRuntimeNotSupport = fmt.Errorf("runtime version not supported")

//CheckRuntime
func CheckRuntime(buildPath string, lang Lang) (map[string]string, error) {
	switch lang {
	case PHP:
		return readPHPRuntimeInfo(buildPath)
	case Python:
		return readPythonRuntimeInfo(buildPath)
	case JavaMaven, JaveWar, JavaJar:
		return readJavaRuntimeInfo(buildPath)
	case Nodejs:
		return readNodeRuntimeInfo(buildPath)
	case NodeJSStatic:
		runtime, err := readNodeRuntimeInfo(buildPath)
		if err != nil {
			return nil, err
		}
		runtime["RUNTIMES_SERVER"] = "nginx"
		return runtime, nil
	case Static:
		return map[string]string{"RUNTIMES_SERVER": "nginx"}, nil
	default:
		return nil, nil
	}
}

func readPHPRuntimeInfo(buildPath string) (map[string]string, error) {
	var phpRuntimeInfo = make(map[string]string, 1)
	if ok, _ := util.FileExists(path.Join(buildPath, "composer.json")); !ok {
		return phpRuntimeInfo, nil
	}
	body, err := ioutil.ReadFile(path.Join(buildPath, "composer.json"))
	if err != nil {
		return phpRuntimeInfo, nil
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return phpRuntimeInfo, nil
	}
	getPhpNewVersion := func(v string) string {
		version := v
		switch v {
		case "5.5":
			version = "5.5.38"
		case "5.6":
			version = "5.6.35"
		case "7.0":
			version = "7.0.29"
		case "7.1":
			version = "7.1.33"
		case "7.2":
			version = "7.2.26"
		case "7.3":
			version = "7.3.13"
		}
		return version
	}
	if json.Get("require") != nil {
		if phpVersion := json.Get("require").Get("php"); phpVersion != nil {
			version, _ := phpVersion.String()
			if version != "" {
				if len(version) < 4 || (version[0:2] == ">=" && len(version) < 5) {
					return nil, ErrRuntimeNotSupport
				}
				if version[0:2] == ">=" {
					if !util.StringArrayContains([]string{"5.5", "5.6", "7.0", "7.1", "7.3"}, version[2:3]) {
						return nil, ErrRuntimeNotSupport
					}
					version = getPhpNewVersion(version[2:3])
				}
				if version[0] == '~' {
					if !util.StringArrayContains([]string{"5.5", "5.6", "7.0", "7.1", "7.3"}, version[1:3]) {
						return nil, ErrRuntimeNotSupport
					}
					version = getPhpNewVersion(version[1:3])
				} else {
					if !util.StringArrayContains([]string{"5.5", "5.6", "7.0", "7.1", "7.3"}, version[0:3]) {
						return nil, ErrRuntimeNotSupport
					}
					version = getPhpNewVersion(version[0:3])
				}
				phpRuntimeInfo["RUNTIMES"] = version
			}
		}
		if hhvmVersion := json.Get("require").Get("hhvm"); hhvmVersion != nil {
			phpRuntimeInfo["RUNTIMES_HHVM"], _ = hhvmVersion.String()
		}
	}
	return phpRuntimeInfo, nil
}

func readPythonRuntimeInfo(buildPath string) (map[string]string, error) {
	var runtimeInfo = make(map[string]string, 1)
	if ok, _ := util.FileExists(path.Join(buildPath, "runtime.txt")); !ok {
		return runtimeInfo, nil
	}
	body, err := ioutil.ReadFile(path.Join(buildPath, "runtime.txt"))
	if err != nil {
		return runtimeInfo, nil
	}
	runtimeInfo["RUNTIMES"] = string(body)
	return runtimeInfo, nil
}

func readJavaRuntimeInfo(buildPath string) (map[string]string, error) {
	var runtimeInfo = make(map[string]string, 1)
	ok, err := util.FileExists(path.Join(buildPath, "system.properties"))
	if !ok || err != nil {
		return runtimeInfo, nil
	}
	cmd := fmt.Sprintf(`grep -i "java.runtime.version" %s | grep  -E -o "[0-9]+(.[0-9]+)?(.[0-9]+)?"`, path.Join(buildPath, "system.properties"))
	runtime, err := util.CmdExec(cmd)
	if err != nil {
		return runtimeInfo, nil
	}
	if runtime != "" {
		runtimeInfo["RUNTIMES"] = runtime
	}
	return runtimeInfo, nil
}

func readNodeRuntimeInfo(buildPath string) (map[string]string, error) {
	var runtimeInfo = make(map[string]string, 1)
	if ok, _ := util.FileExists(path.Join(buildPath, "package.json")); !ok {
		return runtimeInfo, nil
	}
	if ok, _ := util.FileExists(path.Join(buildPath, "yarn.lock")); ok {
		runtimeInfo["PACKAGE_TOOL"] = "yarn"
	}
	if ok, _ := util.FileExists(path.Join(buildPath, "package-lock.json")); ok {
		runtimeInfo["PACKAGE_TOOL"] = "npm"
	}
	body, err := ioutil.ReadFile(path.Join(buildPath, "package.json"))
	if err != nil {
		return runtimeInfo, nil
	}
	json, err := simplejson.NewJson(body)
	if err != nil {
		return runtimeInfo, nil
	}
	if json.Get("engines") != nil {
		if v := json.Get("engines").Get("node"); v != nil {
			nodeVersion, _ := v.String()
			// The latest version is used by default. (11.1.0 is latest version in ui)
			if strings.HasPrefix(nodeVersion, ">") || strings.HasPrefix(nodeVersion, "*") || strings.HasPrefix(nodeVersion, "^") {
				nodeVersion = "11.1.0"
			}
			runtimeInfo["RUNTIMES"] = nodeVersion
		}
	}
	return runtimeInfo, nil
}
