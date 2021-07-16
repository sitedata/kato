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

// +build linux

package info

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/gridworkz/kato/node/nodem/client"
)

//GetSystemInfo
func GetSystemInfo() (info client.NodeSystemInfo) {
	info.Architecture = runtime.GOARCH
	b, _ := ioutil.ReadFile("/etc/machine-id")
	info.MachineID = string(b)
	output, _ := exec.Command("uname", "-r").Output()
	info.KernelVersion = string(output)
	osInfo := readOS()
	if name, ok := osInfo["NAME"]; ok {
		info.OSImage = name
	}
	info.OperatingSystem = runtime.GOOS
	info.MemorySize, _ = getMemory()
	info.NumCPU = int64(runtime.NumCPU())
	return info
}

func readOS() map[string]string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return nil
	}
	defer f.Close()
	var info = make(map[string]string)
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return info
		}
		lines := strings.Split(string(line), "=")
		if len(lines) >= 2 {
			info[lines[0]] = lines[1]
		}
	}
}

func getMemory() (total uint64, free uint64) {
	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err == nil {
		return uint64(sysInfo.Totalram) * uint64(sysInfo.Unit), sysInfo.Freeram * uint64(syscall.Getpagesize())
	}
	return 0, 0
}
