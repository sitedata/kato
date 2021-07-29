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
	"path"

	"github.com/gridworkz/kato/util"
)

func init() {
	checkFuncList = append(checkFuncList, dockerfile)
	checkFuncList = append(checkFuncList, javaJar)
	checkFuncList = append(checkFuncList, javaWar)
	checkFuncList = append(checkFuncList, javaMaven)
	checkFuncList = append(checkFuncList, php)
	checkFuncList = append(checkFuncList, python)
	checkFuncList = append(checkFuncList, nodeJSStatic)
	checkFuncList = append(checkFuncList, nodejs)
	checkFuncList = append(checkFuncList, ruby)
	checkFuncList = append(checkFuncList, static)
	checkFuncList = append(checkFuncList, clojure)
	checkFuncList = append(checkFuncList, golang)
	checkFuncList = append(checkFuncList, gradle)
	checkFuncList = append(checkFuncList, grails)
	checkFuncList = append(checkFuncList, scala)
	checkFuncList = append(checkFuncList, netcore)
}

//ErrCodeNotExist Code is empty error
var ErrCodeNotExist = fmt.Errorf("code is not exist")

//ErrCodeDirNotExist Code directory does not exist
var ErrCodeDirNotExist = fmt.Errorf("code dir is not exist")

//ErrCodeUnableIdentify The code does not recognize the language
var ErrCodeUnableIdentify = fmt.Errorf("code lang unable to identify")

//ErrKatoFileNotFound kato file not found
var ErrKatoFileNotFound = fmt.Errorf("kato file not found")

//Lang Language type
type Lang string

//String return lang string
func (l Lang) String() string {
	return string(l)
}

//NO Empty language type
var NO Lang = "no"

//Dockerfile Lang
var Dockerfile Lang = "dockerfile"

//Docker Lang
var Docker Lang = "docker"

//Python Lang
var Python Lang = "Python"

//Ruby Lang
var Ruby Lang = "Ruby"

//PHP Lang
var PHP Lang = "PHP"

//JavaMaven Lang
var JavaMaven Lang = "Java-maven"

//JaveWar Lang
var JaveWar Lang = "Java-war"

//JavaJar Lang
var JavaJar Lang = "Java-jar"

//Nodejs Lang
var Nodejs Lang = "Node.js"

//NodeJSStatic static Lang
var NodeJSStatic Lang = "NodeJSStatic"

//Static Lang
var Static Lang = "static"

//Clojure Lang
var Clojure Lang = "Clojure"

//Golang Lang
var Golang Lang = "Go"

//Gradle Lang
var Gradle Lang = "Gradle"

//Grails Lang
var Grails Lang = "Grails"

//NetCore Lang
var NetCore Lang = ".NetCore"

//OSS Lang
var OSS Lang = "OSS"

//GetLangType check code lang
func GetLangType(homepath string) (Lang, error) {
	if ok, _ := util.FileExists(homepath); !ok {
		return NO, ErrCodeDirNotExist
	}
	//Determine whether there is a code
	if ok := util.IsHaveFile(homepath); !ok {
		return NO, ErrCodeNotExist
	}
	//Get a certain language
	for _, check := range checkFuncList {
		if lang := check(homepath); lang != NO {
			return lang, nil
		}
	}
	//Get possible languages
	//Unrecognized
	return NO, ErrCodeUnableIdentify
}

type langTypeFunc func(homepath string) Lang

var checkFuncList []langTypeFunc

func dockerfile(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "Dockerfile")); !ok {
		return NO
	}
	return Dockerfile
}
func python(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "requirements.txt")); ok {
		return Python
	}
	if ok, _ := util.FileExists(path.Join(homepath, "setup.py")); ok {
		return Python
	}
	if ok, _ := util.FileExists(path.Join(homepath, "Pipfile")); ok {
		return Python
	}
	return NO
}
func ruby(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "Gemfile")); ok {
		return Ruby
	}
	return NO
}
func php(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "composer.json")); ok {
		return PHP
	}
	if ok := util.SearchFile(homepath, "index.php", 2); ok {
		return PHP
	}
	return NO
}
func javaMaven(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "pom.xml")); ok {
		return JavaMaven
	}
	if ok, _ := util.FileExists(path.Join(homepath, "pom.atom")); ok {
		return JavaMaven
	}
	if ok, _ := util.FileExists(path.Join(homepath, "pom.clj")); ok {
		return JavaMaven
	}
	if ok, _ := util.FileExists(path.Join(homepath, "pom.groovy")); ok {
		return JavaMaven
	}
	if ok, _ := util.FileExists(path.Join(homepath, "pom.rb")); ok {
		return JavaMaven
	}
	if ok, _ := util.FileExists(path.Join(homepath, "pom.scala")); ok {
		return JavaMaven
	}
	if ok, _ := util.FileExists(path.Join(homepath, "pom.yaml")); ok {
		return JavaMaven
	}
	if ok, _ := util.FileExists(path.Join(homepath, "pom.yml")); ok {
		return JavaMaven
	}
	return NO
}
func javaWar(homepath string) Lang {
	if ok := util.FileExistsWithSuffix(homepath, ".war"); ok {
		return JaveWar
	}
	return NO
}

//javaJar Procfile must be defined
func javaJar(homepath string) Lang {
	if ok := util.FileExistsWithSuffix(homepath, ".jar"); ok {
		return JavaJar
	}
	return NO
}
func nodejs(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "package.json")); ok {
		return Nodejs
	}
	return NO
}
func nodeJSStatic(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "package.json")); ok {
		if ok, _ := util.FileExists(path.Join(homepath, "nodestatic.json")); ok {
			return NodeJSStatic
		}
	}
	return NO
}

func static(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "index.html")); ok {
		return Static
	}
	if ok, _ := util.FileExists(path.Join(homepath, "index.htm")); ok {
		return Static
	}
	if ok, _ := util.FileExists(path.Join(homepath, "static.json")); ok {
		return Static
	}
	return NO
}

func clojure(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "project.clj")); ok {
		return Clojure
	}
	return NO
}
func golang(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "go.mod")); ok {
		return Golang
	}
	if ok, _ := util.FileExists(path.Join(homepath, "Gopkg.lock")); ok {
		return Golang
	}
	if ok, _ := util.FileExists(path.Join(homepath, "Godeps", "Godeps.json")); ok {
		return Golang
	}
	if ok, _ := util.FileExists(path.Join(homepath, "vendor", "vendor.json")); ok {
		return Golang
	}
	if ok, _ := util.FileExists(path.Join(homepath, "glide.yaml")); ok {
		return Golang
	}
	if ok := util.FileExistsWithSuffix(path.Join(homepath, "src"), ".go"); ok {
		return Golang
	}
	return NO
}
func gradle(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "build.gradle")); ok {
		return Gradle
	}
	if ok, _ := util.FileExists(path.Join(homepath, "gradlew")); ok {
		return Gradle
	}
	if ok, _ := util.FileExists(path.Join(homepath, "settings.gradle")); ok {
		return Gradle
	}
	return NO
}
func grails(homepath string) Lang {
	if ok, _ := util.FileExists(path.Join(homepath, "grails-app")); ok {
		return Grails
	}
	return NO
}

//netcore
func netcore(homepath string) Lang {
	if ok := util.FileExistsWithSuffix(homepath, ".sln"); ok {
		return NetCore
	}
	if ok := util.FileExistsWithSuffix(homepath, ".csproj"); ok {
		return NetCore
	}
	return NO
}

//Not currently supported
func scala(homepath string) Lang {
	return NO
}
