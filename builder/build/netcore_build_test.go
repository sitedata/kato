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
	"testing"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder/parser/code"
	"github.com/gridworkz/kato/event"
)

func TestBuildNetCore(t *testing.T) {
	build, err := netcoreBuilder()
	if err != nil {
		t.Fatal(err)
	}
	dockerCli, _ := client.NewEnvClient()
	req := &Request{
		SourceDir:     "/Users/devs/gridworkz/dotnet-docker/samples/aspnetapp/test",
		CacheDir:      "/Users/devs/gridworkz/dotnet-docker/samples/aspnetapp/test/cache",
		RepositoryURL: "https://github.com/dotnet/dotnet-docker.git",
		ServiceAlias:  "gr123456",
		DeployVersion: "666666",
		Commit:        Commit{User: "gridworkz"},
		Lang:          code.NetCore,
		Logger:        event.GetTestLogger(),
		DockerClient:  dockerCli,
	}
	res, err := build.Build(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*res)
}
