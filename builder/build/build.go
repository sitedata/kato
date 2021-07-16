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

package build

import (
	"context"
	"fmt"
	"strings"

	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/parser/code"
	"github.com/gridworkz/kato/event"
	"k8s.io/client-go/kubernetes"

	"github.com/docker/docker/client"
)

func init() {
	buildcreaters = make(map[code.Lang]CreaterBuild)
	buildcreaters[code.Dockerfile] = dockerfileBuilder
	buildcreaters[code.Docker] = dockerfileBuilder
	buildcreaters[code.NetCore] = netcoreBuilder
	buildcreaters[code.JavaJar] = slugBuilder
	buildcreaters[code.JavaMaven] = slugBuilder
	buildcreaters[code.JaveWar] = slugBuilder
	buildcreaters[code.PHP] = slugBuilder
	buildcreaters[code.Python] = slugBuilder
	buildcreaters[code.Nodejs] = slugBuilder
	buildcreaters[code.Golang] = slugBuilder
}

var buildcreaters map[code.Lang]CreaterBuild

//Build app build pack
type Build interface {
	Build(*Request) (*Response, error)
}

//CreaterBuild
type CreaterBuild func() (Build, error)

//MediumType Build output medium type
type MediumType string

//ImageMediumType image type
var ImageMediumType MediumType = "image"

//SlugMediumType slug type
var SlugMediumType MediumType = "slug"

//Response build result
type Response struct {
	MediumPath string
	MediumType MediumType
}

//Request build input
type Request struct {
	RbdNamespace  string
	GRDataPVCName string
	CachePVCName  string
	CacheMode     string
	CachePath     string
	TenantID      string
	SourceDir     string
	CacheDir      string
	TGZDir        string
	RepositoryURL string
	Branch        string
	ServiceAlias  string
	ServiceID     string
	DeployVersion string
	Runtime       string
	ServerType    string
	Commit        Commit
	Lang          code.Lang
	BuildEnvs     map[string]string
	Logger        event.Logger
	DockerClient  *client.Client
	KubeClient    kubernetes.Interface
	ExtraHosts    []string
	HostAlias     []HostAlias
	Ctx           context.Context
}

// HostAlias holds the mapping between IP and hostnames that will be injected as an entry in the
// pod's hosts file.
type HostAlias struct {
	// IP address of the host file entry.
	IP string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`
	// Hostnames for the above IP address.
	Hostnames []string `json:"hostnames,omitempty" protobuf:"bytes,2,rep,name=hostnames"`
}

//Commit
type Commit struct {
	User    string
	Message string
	Hash    string
}

//GetBuild
func GetBuild(lang code.Lang) (Build, error) {
	if fun, ok := buildcreaters[lang]; ok {
		return fun()
	}
	return slugBuilder()
}

//CreateImageName
func CreateImageName(serviceID, deployversion string) string {
	return strings.ToLower(fmt.Sprintf("%s/%s:%s", builder.REGISTRYDOMAIN, serviceID, deployversion))
}
