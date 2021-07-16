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

package probe

import (
	"context"

	v1 "github.com/gridworkz/kato/util/prober/types/v1"
)

//Probe
type Probe interface {
	Check()
	Stop()
}

//CreateProbe
func CreateProbe(ctx context.Context, statusChan chan *v1.HealthStatus, v *v1.Service) Probe {
	timeoutSecond := v.ServiceHealth.MaxTimeoutSecond
	if timeoutSecond <= 2 {
		timeoutSecond = 5
	}
	interval := v.ServiceHealth.TimeInterval
	if interval <= 2 {
		interval = 5
	}
	ctx, cancel := context.WithCancel(ctx)
	if v.ServiceHealth.Model == "tcp" {
		t := &TCPProbe{
			Name:          v.ServiceHealth.Name,
			Address:       v.ServiceHealth.Address,
			Ctx:           ctx,
			Cancel:        cancel,
			ResultsChan:   statusChan,
			TimeInterval:  interval,
			MaxErrorsNum:  v.ServiceHealth.MaxErrorsNum,
			TimeoutSecond: timeoutSecond,
		}
		return t
	}
	if v.ServiceHealth.Model == "http" {
		t := &HTTPProbe{
			Name:          v.ServiceHealth.Name,
			Address:       v.ServiceHealth.Address,
			Ctx:           ctx,
			Cancel:        cancel,
			ResultsChan:   statusChan,
			TimeInterval:  v.ServiceHealth.TimeInterval,
			MaxErrorsNum:  v.ServiceHealth.MaxErrorsNum,
			TimeoutSecond: v.ServiceHealth.MaxTimeoutSecond,
		}
		return t
	}
	cancel()
	return nil
}
