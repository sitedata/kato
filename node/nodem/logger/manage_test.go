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

package logger

import (
	"bufio"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/cmd/node/option"
)

func TestWatchConatainer(t *testing.T) {
	dc, err := client.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}
	cm := CreatContainerLogManage(&option.Conf{
		DockerCli: dc,
	})
	cm.Start()
	select {}
}

func TestGetConatainerLogger(t *testing.T) {
	dc, err := client.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}
	cm := CreatContainerLogManage(&option.Conf{
		DockerCli: dc,
	})

	stdout, _, err := cm.getContainerLogReader(cm.ctx, "9874f23cbfc8201571bc654955aad94124f256ccb37e10f914f80865d734a4c5")
	if err != nil {
		t.Fatal(err)
	}
	buffer := bufio.NewReader(stdout)
	for {
		line, _, err := buffer.ReadLine()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(line))
	}
}

func TestHostConfig(t *testing.T) {
	cj := new(types.ContainerJSON)
	if cj.ContainerJSONBase == nil || cj.HostConfig == nil || cj.HostConfig.LogConfig.Type == "" {
		fmt.Println("jsonBase is nil")
		cj.ContainerJSONBase = new(types.ContainerJSONBase)
	}
	if cj.ContainerJSONBase == nil || cj.HostConfig == nil || cj.HostConfig.LogConfig.Type == "" {
		fmt.Println("hostConfig is nil")
		cj.HostConfig = &container.HostConfig{}
	}
	if cj.ContainerJSONBase == nil || cj.HostConfig == nil || cj.HostConfig.LogConfig.Type == "" {
		fmt.Println("logconfig is nil won't panic")
		return
	}

}
