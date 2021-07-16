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

package sources

import (
	"fmt"
	"testing"

	"github.com/gridworkz/kato/event"
	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TestImageName(t *testing.T) {
	imageName := []string{
		"hub.gridworkz.com/nginx:v1",
		"hub.gridworkz.cn/nginx",
		"nginx:v2",
		"tomcat",
	}
	for _, i := range imageName {
		in := ImageNameHandle(i)
		fmt.Printf("host: %s, name: %s, tag: %s\n", in.Host, in.Name, in.Tag)
	}
}

func TestBuildImage(t *testing.T) {
	dc, _ := client.NewEnvClient()
	buildOptions := types.ImageBuildOptions{
		Tags:   []string{"java:test"},
		Remove: true,
	}
	if _, err := ImageBuild(dc, "/Users/gridworkz/coding/java/Demo-RestAPI-Servlet2", buildOptions, nil, 20); err != nil {
		t.Fatal(err)
	}
}

func TestPushImage(t *testing.T) {
	dc, _ := client.NewEnvClient()
	if err := ImagePush(dc, "hub.gridworkz.com/devs-test/etcd:v2.2.0", "devs-test", "devs-test", nil, 2); err != nil {
		t.Fatal(err)
	}
}

func TestTrustedImagePush(t *testing.T) {
	dc, _ := client.NewEnvClient()
	if err := TrustedImagePush(dc, "hub.gridworkz.com/devs-test/etcd:v2.2.0", "devs-test", "devs-test", nil, 2); err != nil {
		t.Fatal(err)
	}
}

func TestCheckTrustedRepositories(t *testing.T) {
	err := CheckTrustedRepositories("hub.gridworkz.com/devs-test/etcd2:v2.2.0", "devs-test", "devs-test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestImageSave(t *testing.T) {
	dc, _ := client.NewEnvClient()
	if err := ImageSave(dc, "hub.gridworkz.com/devs-test/etcd:v2.2.0", "/tmp/testsaveimage.tar", nil); err != nil {
		t.Fatal(err)
	}
}

func TestMultiImageSave(t *testing.T) {
	dc, _ := client.NewEnvClient()
	if err := MultiImageSave(context.Background(), dc, "/tmp/testsaveimage.tar", nil,
		"registry.gitlab.com/gridworkz/rbd-node:V5.3.0-cloud",
		"registry.gitlab.com/gridworkz/rbd-resource-proxy:V5.3.0-cloud"); err != nil {
		t.Fatal(err)
	}
}

func TestImageImport(t *testing.T) {
	dc, _ := client.NewEnvClient()
	if err := ImageImport(dc, "hub.gridworkz.com/devs-test/etcd:v2.2.0", "/tmp/testsaveimage.tar", nil); err != nil {
		t.Fatal(err)
	}
}

func TestImagePull(t *testing.T) {
	dc, _ := client.NewEnvClient()
	_, err := ImagePull(dc, "gridworkz/collabora:190422", "", "", event.GetTestLogger(), 60)
	if err != nil {
		t.Fatal(err)
	}
}
