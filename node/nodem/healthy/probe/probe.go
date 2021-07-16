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
	"fmt"
	"strings"

	"github.com/gridworkz/kato/node/nodem/client"
	"github.com/gridworkz/kato/node/nodem/service"
)

//Probe
type Probe interface {
	Check()
	Stop()
}

//CreateProbe
func CreateProbe(ctx context.Context, hostNode *client.HostNode, statusChan chan *service.HealthStatus, v *service.Service) (Probe, error) {
	ctx, cancel := context.WithCancel(ctx)
	model := strings.ToLower(strings.TrimSpace(v.ServiceHealth.Model))
	switch model {
	case "http":
		return &HttpProbe{
			Name:         v.ServiceHealth.Name,
			Address:      v.ServiceHealth.Address,
			Ctx:          ctx,
			Cancel:       cancel,
			ResultsChan:  statusChan,
			TimeInterval: v.ServiceHealth.TimeInterval,
			HostNode:     hostNode,
			MaxErrorsNum: v.ServiceHealth.MaxErrorsNum,
		}, nil
	case "tcp":
		return &TcpProbe{
			Name:         v.ServiceHealth.Name,
			Address:      v.ServiceHealth.Address,
			Ctx:          ctx,
			Cancel:       cancel,
			ResultsChan:  statusChan,
			TimeInterval: v.ServiceHealth.TimeInterval,
			HostNode:     hostNode,
			MaxErrorsNum: v.ServiceHealth.MaxErrorsNum,
		}, nil
	case "cmd":
		return &ShellProbe{
			Name:         v.ServiceHealth.Name,
			Address:      v.ServiceHealth.Address,
			Ctx:          ctx,
			Cancel:       cancel,
			ResultsChan:  statusChan,
			TimeInterval: v.ServiceHealth.TimeInterval,
			HostNode:     hostNode,
			MaxErrorsNum: v.ServiceHealth.MaxErrorsNum,
		}, nil
	default:
		cancel()
		return nil, fmt.Errorf("service %s probe mode %s not support ", v.Name, model)
	}
}
