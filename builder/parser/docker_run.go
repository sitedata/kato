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

package parser

import (
	"fmt"
	"github.com/docker/distribution/reference" //"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder/parser/types"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
	"runtime"
	"strconv"
	"strings" //"github.com/docker/docker/client"
)

//DockerRunOrImageParse docker run - command resolution or direct image name resolution
type DockerRunOrImageParse struct {
	user, pass   string
	ports        map[int]*types.Port
	volumes      map[string]*types.Volume
	envs         map[string]*types.Env
	source       string
	serviceType  string
	memory       int
	image        Image
	args         []string
	errors       []ParseError
	dockerclient *client.Client
	logger       event.Logger
}

//CreateDockerRunOrImageParse create parser
func CreateDockerRunOrImageParse(user, pass, source string, dockerclient *client.Client, logger event.Logger) *DockerRunOrImageParse {
	source = strings.TrimLeft(source, " ")
	source = strings.Replace(source, "\n", "", -1)
	source = strings.Replace(source, "\\", "", -1)
	source = strings.Replace(source, "  ", " ", -1)
	return &DockerRunOrImageParse{
		user:         user,
		pass:         pass,
		source:       source,
		dockerclient: dockerclient,
		ports:        make(map[int]*types.Port),
		volumes:      make(map[string]*types.Volume),
		envs:         make(map[string]*types.Env),
		logger:       logger,
	}
}

//Parse - decode, get the image, resolve the image
//eg. docker run -it -p 80:80 nginx
func (d *DockerRunOrImageParse) Parse() ParseErrorList {
	if d.source == "" {
		d.errappend(Errorf(FatalError, "source can not be empty"))
		return d.errors
	}
	//docker run
	if strings.HasPrefix(d.source, "docker") {
		d.ParseDockerun(d.source)
		if d.image.String() == "" || d.image.String() == ":" {
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Mirror name recognition failed"), SolveAdvice("modify_image", "Please confirm whether the input DockerRun command is correct")))
			return d.errors
		}
		if _, err := reference.ParseAnyReference(d.image.String()); err != nil {
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Image name (%s) is illegal", d.image.String()), SolveAdvice("modify_image", "Please confirm whether the input DockerRun command is correct")))
			return d.errors
		}
	} else {
		//else image
		_, err := reference.ParseAnyReference(d.source)
		if err != nil {
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("The image name (%s) is illegal", d.image.String()), SolveAdvice("modify_image", "Please confirm whether the input mirror name is correct")))
			return d.errors
		}
		d.image = ParseImageName(d.source)
	}
	//Get the image and verify if it exists
	imageInspect, err := sources.ImagePull(d.dockerclient, d.image.Source(), d.user, d.pass, d.logger, 10)
	if err != nil {
		if strings.Contains(err.Error(), "No such image") {
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Mirror (%s) does not exist", d.image.String()), SolveAdvice("modify_image", "Please confirm whether the input mirror name is correct")))
		} else {
			if d.image.IsOfficial() {
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Mirror (%s) acquisition failed, and domestic access to the official Docker warehouse is often unstable", d.image.String()), SolveAdvice("modify_image", "Please confirm that the input image can be obtained normally")))
			} else {
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Mirror (%s) acquisition failed", d.image.String()), SolveAdvice("modify_image", "Please confirm that the input image can be obtained normally")))
			}
		}
		return d.errors
	}
	if imageInspect != nil && imageInspect.ContainerConfig != nil {
		for _, env := range imageInspect.ContainerConfig.Env {
			envinfo := strings.Split(env, "=")
			if len(envinfo) == 2 {
				if _, ok := d.envs[envinfo[0]]; !ok {
					d.envs[envinfo[0]] = &types.Env{Name: envinfo[0], Value: envinfo[1]}
				}
			}
		}
		for k := range imageInspect.ContainerConfig.Volumes {
			if _, ok := d.volumes[k]; !ok {
				d.volumes[k] = &types.Volume{VolumePath: k, VolumeType: model.ShareFileVolumeType.String()}
			}
		}
		for k := range imageInspect.ContainerConfig.ExposedPorts {
			proto := k.Proto()
			port := k.Int()
			if proto != "udp" {
				proto = GetPortProtocol(port)
			}
			if _, ok := d.ports[port]; ok {
				d.ports[port].Protocol = proto
			} else {
				d.ports[port] = &types.Port{Protocol: proto, ContainerPort: port}
			}
		}
	}
	d.serviceType = DetermineDeployType(d.image)
	return d.errors
}

//ParseDockerun
func (d *DockerRunOrImageParse) ParseDockerun(cmd string) {
	var name string
	cmd = strings.TrimLeft(cmd, " ")
	cmd = strings.Replace(cmd, "\n", "", -1)
	cmd = strings.Replace(cmd, "\r", "", -1)
	cmd = strings.Replace(cmd, "\t", "", -1)
	cmd = strings.Replace(cmd, "\\", "", -1)
	cmd = strings.Replace(cmd, "  ", " ", -1)
	source := util.RemoveSpaces(strings.Split(cmd, " "))
	for i, s := range source {
		if s == "docker" || s == "run" {
			continue
		}
		if strings.HasPrefix(s, "-") {
			name = strings.TrimLeft(s, "-")
			index := strings.Index(name, "=")
			if index > 0 && index < len(name)-1 {
				s = name[index+1:]
				name = name[0:index]
				switch name {
				case "e", "env":
					info := strings.Split(s, "=")
					if len(info) == 2 {
						d.envs[info[0]] = &types.Env{Name: info[0], Value: info[1]}
					}
				case "p", "public":
					info := strings.Split(s, ":")
					if len(info) == 2 {
						port, _ := strconv.Atoi(info[0])
						if port != 0 {
							d.ports[port] = &types.Port{ContainerPort: port, Protocol: GetPortProtocol(port)}
						}
					}
				case "v", "volume":
					info := strings.Split(s, ":")
					if len(info) >= 2 {
						d.volumes[info[1]] = &types.Volume{VolumePath: info[1], VolumeType: model.ShareFileVolumeType.String()}
					}
				case "memory", "m":
					d.memory = readmemory(s)
				}
				name = ""
			}
		} else {
			switch name {
			case "e", "env":
				info := strings.Split(removeQuotes(s), "=")
				if len(info) == 2 {
					d.envs[info[0]] = &types.Env{Name: info[0], Value: info[1]}
				}
			case "p", "public":
				info := strings.Split(removeQuotes(s), ":")
				if len(info) == 2 {
					port, _ := strconv.Atoi(info[1])
					if port != 0 {
						d.ports[port] = &types.Port{ContainerPort: port, Protocol: GetPortProtocol(port)}
					}
				}
			case "v", "volume":
				info := strings.Split(removeQuotes(s), ":")
				if len(info) >= 2 {
					d.volumes[info[1]] = &types.Volume{VolumePath: info[1], VolumeType: model.ShareFileVolumeType.String()}
				}
			case "memory", "m":
				d.memory = readmemory(s)
			case "", "d", "i", "t", "it", "P", "rm", "init", "interactive", "no-healthcheck", "oom-kill-disable", "privileged", "read-only", "tty", "sig-proxy":
				d.image = ParseImageName(s)
				if len(source) > i+1 {
					d.args = source[i+1:]
				}
				return
			}
			name = ""
			continue
		}
	}

}

func (d *DockerRunOrImageParse) errappend(pe ParseError) {
	d.errors = append(d.errors, pe)
}

//GetBranchs
func (d *DockerRunOrImageParse) GetBranchs() []string {
	return nil
}

//GetPorts
func (d *DockerRunOrImageParse) GetPorts() (ports []types.Port) {
	for _, cv := range d.ports {
		ports = append(ports, *cv)
	}
	return ports
}

//GetVolumes
func (d *DockerRunOrImageParse) GetVolumes() (volumes []types.Volume) {
	for _, cv := range d.volumes {
		volumes = append(volumes, *cv)
	}
	return
}

//GetValid
func (d *DockerRunOrImageParse) GetValid() bool {
	return false
}

//GetEnvs
func (d *DockerRunOrImageParse) GetEnvs() (envs []types.Env) {
	for _, cv := range d.envs {
		envs = append(envs, *cv)
	}
	return
}

//GetImage
func (d *DockerRunOrImageParse) GetImage() Image {
	return d.image
}

//GetArgs
func (d *DockerRunOrImageParse) GetArgs() []string {
	return d.args
}

//GetMemory
func (d *DockerRunOrImageParse) GetMemory() int {
	return d.memory
}

//GetServiceInfo
func (d *DockerRunOrImageParse) GetServiceInfo() []ServiceInfo {
	serviceInfo := ServiceInfo{
		Ports:       d.GetPorts(),
		Envs:        d.GetEnvs(),
		Volumes:     d.GetVolumes(),
		Image:       d.GetImage(),
		Args:        d.GetArgs(),
		Branchs:     d.GetBranchs(),
		Memory:      d.memory,
		ServiceType: d.serviceType,
		OS:          runtime.GOOS,
	}
	if serviceInfo.Memory == 0 {
		serviceInfo.Memory = 512
	}
	return []ServiceInfo{serviceInfo}
}
