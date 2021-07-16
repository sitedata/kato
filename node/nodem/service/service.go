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

package service

import (
	"time"

	"github.com/gridworkz/kato/util"

	yaml "gopkg.in/yaml.v2"
)

const (
	Stat_Unknow    string = "unknow"    //unknown health
	Stat_healthy   string = "healthy"   //healthy
	Stat_unhealthy string = "unhealthy" //unhealthy
	Stat_death     string = "death"     //dead
)

//Service
type Service struct {
	Name            string      `yaml:"name"`
	Endpoints       []*Endpoint `yaml:"endpoints,omitempty"`
	ServiceHealth   *Health     `yaml:"health"`
	OnlyHealthCheck bool        `yaml:"only_health_check"`
	IsInitStart     bool        `yaml:"is_init_start"`
	Disable         bool        `yaml:"disable"`
	After           []string    `yaml:"after"`
	Requires        []string    `yaml:"requires"`
	Type            string      `yaml:"type,omitempty"`
	PreStart        string      `yaml:"pre_start,omitempty"`
	Start           string      `yaml:"start"`
	Stop            string      `yaml:"stop,omitempty"`
	RestartPolicy   string      `yaml:"restart_policy,omitempty"`
	RestartSec      string      `yaml:"restart_sec,omitempty"`
}

//Equal
func (s *Service) Equal(e *Service) bool {
	sb, err := yaml.Marshal(s)
	if err != nil {
		return false
	}
	eb, err := yaml.Marshal(e)
	if err != nil {
		return false
	}
	if util.BytesSliceEqual(sb, eb) {
		return true
	}
	return false
}

//Services default config of all services
type Services struct {
	Version  string     `yaml:"version"`
	Services []*Service `yaml:"services"`
	FromFile string     `yaml:"-"`
}

//Endpoint
type Endpoint struct {
	Name     string `yaml:"name"`
	Protocol string `yaml:"protocol"`
	Port     string `yaml:"port"`
}

//Health ServiceHealth
type Health struct {
	Name         string `yaml:"name"`
	Model        string `yaml:"model"`
	Address      string `yaml:"address"`
	TimeInterval int    `yaml:"time_interval"`
	MaxErrorsNum int    `yaml:"max_errors_num"`
}

//HealthStatus
type HealthStatus struct {
	Name           string
	Status         string
	ErrorNumber    int
	ErrorDuration  time.Duration
	StartErrorTime time.Time
	Info           string
}
