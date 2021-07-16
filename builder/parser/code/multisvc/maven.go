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

package multi

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/gridworkz/kato/builder/parser/types"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
)

// Build represents build in pom.xml
type Build struct {
	FinalName string   `xml:"finalName"`
	Plugins   *Plugins `xml:"plugins"`
}

// Plugins represents plugins in pom.xml
type Plugins struct {
	Plugin []*Plugin `xml:"plugin"`
}

// Plugin represents plugin in pom.xml
type Plugin struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	FinalName  string `xml:"configuration>finalName"`
}

// maven is an implementation of MutilModuler.
type maven struct {
}

// NewMaven creates a new MultiModuler
func NewMaven() ServiceInterface {
	return &maven{}
}

// pom represents a maven pom.xml file
type pom struct {
	XMLName    xml.Name `xml:"project"`
	Name       string   `xml:"name"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Packaging  string   `xml:"packaging"`
	Modules    []string `xml:"modules>module"`
	Build      *Build   `xml:"build"`
}

// module represents a maven module
type module struct {
	ID               string
	Name             string
	MavenCustomOpts  string
	MavenCustomGoals string
	MavenJavaOpts    string
	Packaging        string
	Procfile         string
}

// ListModules lists all maven modules from pom.xml
func (m *maven) ListModules(path string) ([]*types.Service, error) {
	modules, err := listModules(path, strings.TrimRight(path, "/")+"/", "")
	if err != nil {
		return nil, err
	}
	var res []*types.Service
	for _, item := range modules {
		envs := []*types.Env{
			{
				Name:  "BUILD_MAVEN_CUSTOM_OPTS",
				Value: item.MavenCustomOpts,
			},
			{
				Name:  "BUILD_MAVEN_CUSTOM_GOALS",
				Value: item.MavenCustomGoals,
			},
			{
				Name:  "BUILD_PROCFILE",
				Value: item.Procfile,
			},
		}
		mo := &types.Service{
			ID:   item.ID,
			Name: item.Name,
			Cname: func(name string) string {
				cnames := strings.Split(name, "/")
				return cnames[len(cnames)-1]
			}(item.Name),
			Packaging: item.Packaging,
			Envs:      make(map[string]*types.Env),
		}
		for _, env := range envs {
			mo.Envs[env.Name] = &types.Env{Name: env.Name, Value: env.Value}
		}
		res = append(res, mo)
	}
	return res, nil
}

func listModules(prefix, topPref, finalName string) ([]*module, error) {
	pomPath := path.Join(prefix, "pom.xml")
	pom, err := parsePom(pomPath)
	if err != nil {
		return nil, err
	}

	if pom.Build != nil && pom.Build.FinalName != "" {
		finalName = pom.Build.FinalName
	}

	var modules []*module // module names
	// recursive end condition
	if pom.isValidModule() {
		filename := pom.getExecuteFilename(finalName)
		// full module name. eg: foobar/rbd-worker
		name := strings.Replace(prefix, topPref, "", 1)
		mo := &module{
			ID:   util.NewUUID(),
			Name: name,
			Packaging: func() string {
				if pom.Packaging == "war" {
					return "war"
				}
				return "jar"
			}(),
			MavenCustomOpts:  fmt.Sprintf("-DskipTests"),
			MavenCustomGoals: fmt.Sprintf("clean dependency:list install -pl %s -am", name),
			Procfile: func() string {
				if pom.Packaging == "war" {
					return fmt.Sprintf("web: java $JAVA_OPTS -jar /opt/webapp-runner.jar "+
						"--port $PORT %s/target/%s", name, filename)
				}
				return fmt.Sprintf("web: java $JAVA_OPTS -jar %s/target/%s", name, filename)
			}(),
		}
		return []*module{mo}, nil
	}

	for _, name := range pom.Modules {
		// submodule names
		submodules, err := listModules(path.Join(prefix, name), topPref, finalName)
		if err != nil {
			logrus.Warningf("Prefix: %s; error getting module names: %v",
				path.Join(prefix, name), err)
			continue
		}
		if submodules != nil && len(submodules) > 0 {
			modules = append(modules, submodules...)
		}
	}
	return modules, nil
}

// parsePom parses the pom.xml file into a pom struct
func parsePom(pomPath string) (*pom, error) {
	bytes, err := ioutil.ReadFile(pomPath)
	if err != nil {
		return nil, err
	}
	var pom pom
	if err = xml.Unmarshal(bytes, &pom); err != nil {
		return nil, err
	}
	return &pom, nil
}

// checks if the pom has submodules.
func (p *pom) hasSubmodules() bool {
	return len(p.Modules) > 0
}

// TODO: read maven source code, learn how does maven get the final name
func (p *pom) getExecuteFilename(finalName string) string {
	// default finalName
	name := p.ArtifactID + "-*"
	if p.Build != nil {
		// the finalName in the plugin has a higher priority than in the build.
		if p.Build.FinalName != "" {
			finalName = p.Build.FinalName
		}
		if p.Build.Plugins != nil && p.Build.Plugins.Plugin != nil && len(p.Build.Plugins.Plugin) > 0 {
			for _, plugin := range p.Build.Plugins.Plugin {
				if plugin.ArtifactID == "spring-boot-maven-plugin" && plugin.GroupID == "org.springframework.boot" &&
					plugin.FinalName != "" {
					finalName = plugin.FinalName
					break
				}
			}
		}
	}
	if finalName == "${project.name}" {
		name = p.ArtifactID
		if p.Name != "" {
			name = p.Name
		}
	} else if finalName == "${project.artifactId}" {
		name = p.ArtifactID
	} else if strings.HasPrefix(finalName, "") && strings.HasSuffix(finalName, "}") {
		name = "*"
	} else if finalName != "" {
		name = finalName
	}
	suffix := func() string {
		if p.Packaging == "" {
			return "jar"
		}
		return p.Packaging
	}
	return name + "." + suffix()
}

func (p *pom) isValidModule() bool {
	if p.Packaging != "jar" && p.Packaging != "war" && p.Packaging != "" {
		return false
	}
	if p.Modules != nil || len(p.Modules) > 0 {
		return false
	}
	return true
}
