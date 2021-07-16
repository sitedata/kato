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
	"runtime"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder/parser/compose"
	"github.com/gridworkz/kato/builder/parser/types"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	"github.com/sirupsen/logrus"
)

//DockerComposeParse docker compose - file analysis
type DockerComposeParse struct {
	services     map[string]*ServiceInfoFromDC
	errors       []ParseError
	dockerclient *client.Client
	logger       event.Logger
	source       string
	user         string
	password     string
}

//ServiceInfoFromDC service info from dockercompose
type ServiceInfoFromDC struct {
	ports       map[int]*types.Port
	volumes     map[string]*types.Volume
	envs        map[string]*types.Env
	source      string
	memory      int
	image       Image
	args        []string
	depends     []string
	imageAlias  string
	serviceType string
	name        string
}

//GetPorts 
func (d *ServiceInfoFromDC) GetPorts() (ports []types.Port) {
	for _, cv := range d.ports {
		ports = append(ports, *cv)
	}
	return ports
}

//GetVolumes 
func (d *ServiceInfoFromDC) GetVolumes() (volumes []types.Volume) {
	for _, cv := range d.volumes {
		volumes = append(volumes, *cv)
	}
	return
}

//GetEnvs
func (d *ServiceInfoFromDC) GetEnvs() (envs []types.Env) {
	for _, cv := range d.envs {
		envs = append(envs, *cv)
	}
	return
}

//CreateDockerComposeParse
func CreateDockerComposeParse(source string, dockerclient *client.Client, user, pass string, logger event.Logger) Parser {
	return &DockerComposeParse{
		source:       source,
		dockerclient: dockerclient,
		logger:       logger,
		services:     make(map[string]*ServiceInfoFromDC),
		user:         user,
		password:     pass,
	}
}

//Parse
func (d *DockerComposeParse) Parse() ParseErrorList {
	if d.source == "" {
		d.errappend(Errorf(FatalError, "source can not be empty"))
		return d.errors
	}
	comp := compose.Compose{}
	co, err := comp.LoadBytes([][]byte{[]byte(d.source)})
	if err != nil {
		logrus.Warning("parse compose file error,", err.Error())
		d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("ComposeFile parsing error"), SolveAdvice("modify_compose", "Please confirm whether the input of ComposeFile is grammatically correct")))
		return d.errors
	}
	for kev, sc := range co.ServiceConfigs {
		logrus.Debugf("service config is %v, container name is %s", sc, sc.ContainerName)
		ports := make(map[int]*types.Port)

		if sc.Image == "" {
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("ComposeFile parsing error"), SolveAdvice(fmt.Sprintf("Service %s has no image specified", kev), fmt.Sprintf("Please specify a mirror for %s", kev))))
			continue
		}

		for _, p := range sc.Port {
			pro := string(p.Protocol)
			if pro != "udp" {
				pro = GetPortProtocol(int(p.ContainerPort))
			}
			ports[int(p.ContainerPort)] = &types.Port{
				ContainerPort: int(p.ContainerPort),
				Protocol:      pro,
			}
		}
		volumes := make(map[string]*types.Volume)
		for _, v := range sc.Volumes {
			if strings.Contains(v.MountPath, ":") {
				infos := strings.Split(v.MountPath, ":")
				if len(infos) > 1 {
					volumes[v.MountPath] = &types.Volume{
						VolumePath: infos[1],
						VolumeType: model.ShareFileVolumeType.String(),
					}
				}
			} else {
				volumes[v.MountPath] = &types.Volume{
					VolumePath: v.MountPath,
					VolumeType: model.ShareFileVolumeType.String(),
				}
			}
		}
		envs := make(map[string]*types.Env)
		for _, e := range sc.Environment {
			envs[e.Name] = &types.Env{
				Name:  e.Name,
				Value: e.Value,
			}
		}
		service := ServiceInfoFromDC{
			ports:      ports,
			volumes:    volumes,
			envs:       envs,
			memory:     int(sc.MemLimit / 1024 / 1024),
			image:      ParseImageName(sc.Image),
			args:       sc.Args,
			depends:    sc.Links,
			imageAlias: sc.ContainerName,
			name:       kev,
		}
		if sc.DependsON != nil {
			service.depends = sc.DependsON
		}
		service.serviceType = DetermineDeployType(service.image)
		d.services[kev] = &service
	}
	for serviceName, service := range d.services {
		//Verify that depends is complete
		existDepends := []string{}
		for i, depend := range service.depends {
			if strings.Contains(depend, ":") {
				service.depends[i] = strings.Split(depend, ":")[0]
			}
			if _, ok := d.services[service.depends[i]]; !ok {
				d.errappend(ErrorAndSolve(NegligibleError, fmt.Sprintf("Service %s dependency definition error", serviceName), SolveAdvice("modify_compose", fmt.Sprintf("Please confirm whether the dependent service of %s service is correct", serviceName))))
			} else {
				existDepends = append(existDepends, service.depends[i])
			}
		}
		service.depends = existDepends
		var hubUser = d.user
		var hubPass = d.password
		for _, env := range service.GetEnvs() {
			if env.Name == "HUB_USER" {
				hubUser = env.Value
			}
			if env.Name == "HUB_PASSWORD" {
				hubPass = env.Value
			}
		}
		//do not pull image, but check image exist
		d.logger.Debug(fmt.Sprintf("start check service %s ", service.name), map[string]string{"step": "service_check", "status": "running"})
		exist, err := sources.ImageExist(service.image.String(), hubUser, hubPass)
		if err != nil {
			logrus.Errorf("check image(%s) exist failure %s", service.image.String(), err.Error())
		}
		if !exist {
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Service %s mirror %s detection failed", serviceName, service.image.String()), SolveAdvice("modify_compose", fmt.Sprintf("Please confirm whether the image name of the %s service is correct or whether the mirror warehouse access is normal", serviceName))))
		}
	}
	return d.errors
}

func (d *DockerComposeParse) errappend(pe ParseError) {
	d.errors = append(d.errors, pe)
}

//GetServiceInfo
func (d *DockerComposeParse) GetServiceInfo() []ServiceInfo {
	var sis []ServiceInfo
	for _, service := range d.services {
		si := ServiceInfo{
			Ports:          service.GetPorts(),
			Envs:           service.GetEnvs(),
			Volumes:        service.GetVolumes(),
			Image:          service.image,
			Args:           service.args,
			DependServices: service.depends,
			ImageAlias:     service.imageAlias,
			ServiceType:    service.serviceType,
			Name:           service.name,
			Cname:          service.name,
			OS:             runtime.GOOS,
		}
		if service.memory != 0 {
			si.Memory = service.memory
		} else {
			si.Memory = 512
		}
		sis = append(sis, si)
	}
	return sis
}

//GetImage
func (d *DockerComposeParse) GetImage() Image {
	return Image{}
}
