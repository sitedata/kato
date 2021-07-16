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

package v1

import "time"

const (
	// StatHealthy -
	StatHealthy string = "healthy"
	// StatUnhealthy -
	StatUnhealthy string = "unhealthy"
	// StatDeath -
	StatDeath string = "death"
)

// Service
type Service struct {
	Sid           string  `json:"service_id"`
	Name          string  `json:"name"`
	ServiceHealth *Health `json:"health"`
	Disable       bool    `json:"disable"`
}

// Equal check if the left service(l) is equal to the right service(r)
func (l *Service) Equal(r *Service) bool {
	if l == r {
		return true
	}
	if l.Sid != r.Sid {
		return false
	}
	if l.Name != r.Name {
		return false
	}
	if l.Disable != r.Disable {
		return false
	}
	if !l.ServiceHealth.Equal(r.ServiceHealth) {
		return false
	}
	return true
}

//Health ServiceHealth
type Health struct {
	Name             string `json:"name"`
	Model            string `json:"model"`
	IP               string `json:"ip"`
	Port             int    `json:"port"`
	Address          string `json:"address"`
	TimeInterval     int    `json:"time_interval"`
	MaxErrorsNum     int    `json:"max_errors_num"`
	MaxTimeoutSecond int    `json:"max_timeout"`
}

// Equal check if the left health(l) is equal to the right health(r)
func (l *Health) Equal(r *Health) bool {
	if l == r {
		return true
	}
	if l.Name != r.Name {
		return false
	}
	if l.Model != r.Model {
		return false
	}
	if l.IP != r.IP {
		return false
	}
	if l.Port != r.Port {
		return false
	}
	if l.Address != r.Address {
		return false
	}
	if l.TimeInterval != r.TimeInterval {
		return false
	}
	if l.MaxErrorsNum != r.MaxErrorsNum {
		return false
	}
	return true
}

//HealthStatus
type HealthStatus struct {
	Name           string        `json:"name"`
	Status         string        `json:"status"`
	ErrorNumber    int           `json:"error_number"`
	ErrorDuration  time.Duration `json:"error_duration"`
	StartErrorTime time.Time     `json:"start_error_time"`
	Info           string        `json:"info"`
	LastStatus     string        `json:"last_status"`
	StatusChange   bool          `json:"status_change"`
}
