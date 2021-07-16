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
	"io/ioutil"
	"path"
	"strings"

	"github.com/gridworkz/kato/util"
)

//Specification
type Specification struct {
	//Whether it meets the specification
	Conform bool
	//Non-conforming project solution
	Noconform map[string]string
	//Proposal to standardize the project
	Advice map[string]string
}

//Various types of language specifications
var specification map[Lang]func(buildPath string) Specification

func init() {
	specification = make(map[Lang]func(buildPath string) Specification)
	specification[JavaJar] = javaJarCheck
	specification[JavaMaven] = javaMavenCheck
	specification[PHP] = phpCheck
	specification[NodeJSStatic] = nodeCheck
	specification[Nodejs] = nodeCheck
	specification[Golang] = golangCheck
}

//CheckCodeSpecification
func CheckCodeSpecification(buildPath string, lang Lang) Specification {
	if check, ok := specification[lang]; ok {
		return check(buildPath)
	}
	return common()
}

//Procfile must be defined
//The jar package defined in the Procfile must exist
func javaJarCheck(buildPath string) Specification {
	procfile, spec := checkProcfile(buildPath)
	if spec != nil {
		return *spec
	}
	if !procfile {
		return Specification{
			Conform:   false,
			Noconform: map[string]string{"Recognized as JavaJar language, Procfile file is not defined": "The main directory defines the Procfile file to specify the jar package startup mode, reference format:\n web: java $JAVA_OPTS -jar demo.jar"},
		}
	}
	return common()
}

// Find whether the pom.xml file contains org.springframework.boot
// Not included (mandatory)
// Must be introduced in pom.xml webapp-runner.jar
// The packaging type must be war
// It is recommended to write Procfile (if not written, the platform is set by default)
func javaMavenCheck(buildPath string) Specification {
	procfile, spec := checkProcfile(buildPath)
	if spec != nil {
		return *spec
	}
	if ok, _ := util.FileExists(path.Join(buildPath, "pom.xml")); !ok {
		return Specification{
			Conform:   false,
			Noconform: map[string]string{"Recognized as JavaMaven language, no pom.xml file found in the working directory": "define the pom.xml file"},
		}
	}
	ok := util.SearchFileBody(path.Join(buildPath, "pom.xml"), "<modules>")
	if ok {
		return common()
	}
	//Determine whether pom.xml contains org.springframework.boot definition
	ok = util.SearchFileBody(path.Join(buildPath, "pom.xml"), "org.springframework.boot")
	if !ok {
		//By default, it can only be packaged as a war package
		war := util.SearchFileBody(path.Join(buildPath, "pom.xml"), "<packaging>war</packaging>")
		if !war && !procfile {
			//If defined as a jar package, Procfile must be defined
			return Specification{
				Conform:   false,
				Noconform: map[string]string{"Recognized as JavaMaven language, non-SpringBoot projects only support War packaging by default": "see the official JaveMaven project code configuration"},
			}
		}
		//TODO: Check if the procfile definition is correct
	}
	return common()
}

//checkProcfile
func checkProcfile(buildPath string) (bool, *Specification) {
	if ok, _ := util.FileExists(path.Join(buildPath, "Procfile")); !ok {
		return false, nil
	}
	procfile, err := ioutil.ReadFile(path.Join(buildPath, "Procfile"))
	if err != nil {
		return false, nil
	}
	infos := strings.Split(strings.TrimRight(string(procfile), " "), " ")
	if len(infos) < 2 {
		return true, &Specification{
			Conform:   false,
			Noconform: map[string]string{"Procfile file does not meet specifications": "reference format\n web: start command"},
		}
	}
	if infos[0] != "web:" {
		return true, &Specification{
			Conform:   false,
			Noconform: map[string]string{"The Procfile file specification currently only supports web: beginning": "reference format\n web: start command"},
		}
	}
	return true, nil
}
func common() Specification {
	return Specification{
		Conform: true,
	}
}

func phpCheck(buildPath string) Specification {
	if ok, _ := util.FileExists(path.Join(buildPath, "composer.lock")); !ok {
		return Specification{
			Conform:   false,
			Noconform: map[string]string{"Recognized as PHP language, no composer.lock file found in the code directory": "Composer.lock file must be generated"},
		}
	}
	return common()
}
func nodeCheck(buildPath string) Specification {
	var yarn, npm bool
	if ok, _ := util.FileExists(path.Join(buildPath, "yarn.lock")); ok {
		yarn = true
	}
	if ok, _ := util.FileExists(path.Join(buildPath, "package-lock.json")); ok {
		npm = true
	}
	if !yarn && !npm {
		return Specification{
			Conform:   false,
			Noconform: map[string]string{"No yarn.lock or package-lock.json files found in the code directory": "must generate and submit the yarn.lock or package-lock.json file"},
		}
	}
	return common()
}

func golangCheck(buildPath string) Specification {
	return common()
}
