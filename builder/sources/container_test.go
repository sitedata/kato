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

package sources

import (
	"log"
	"testing"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func createService() *DockerService {
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
	service := CreateDockerService(context.Background(), client)
	return service
}
func TestCreateContainer(t *testing.T) {
	service := createService()
	containerID, err := service.CreateContainer(&ContainerConfig{
		Metadata: &ContainerMetadata{
			Name: "test",
		},
		Image: &ImageSpec{
			Image: "nginx",
		},
		Args: []string{""},
		Devices: []*Device{
			&Device{
				ContainerPath: "/data",
				HostPath:      "/tmp/test",
			},
		},
		Mounts: []*Mount{
			&Mount{
				ContainerPath: "/data2",
				HostPath:      "/tmp/test2",
				Readonly:      false,
			},
		},
		Envs: []*KeyValue{
			&KeyValue{Key: "ABC", Value: "ASDA ASDASD ASDASD \"ASDAASD"},
		},
		Stdin: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(containerID)
}

func TestStartContainer(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	service := createService()
	containerID, err := service.CreateContainer(&ContainerConfig{
		Metadata: &ContainerMetadata{
			Name: "test",
		},
		Image: &ImageSpec{
			Image: "containertest",
		},
		Mounts: []*Mount{
			&Mount{
				ContainerPath: "/data2",
				HostPath:      "/tmp/test2",
				Readonly:      false,
			},
		},
		Envs: []*KeyValue{
			&KeyValue{Key: "ABC", Value: "ASDA ASDASD ASDASD \"ASDAASD"},
		},
		Stdin: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	errchan := make(chan error, 1)
	//start the container
	if err := service.StartContainer(containerID); err != nil {
		<-errchan
		t.Fatal(err)
	}
	if errchan != nil {
		if err := <-errchan; err != nil {
			logrus.Debugf("Error hijack: %s", err)
			t.Fatal(err)
		}
	}
	t.Log(containerID)
}
