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

package cluster

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/cmd/gateway/option"
	"github.com/sirupsen/logrus"
)

//NodeManager
type NodeManager struct {
	config    option.Config
	ipManager IPManager
}

//CreateNodeManager
func CreateNodeManager(ctx context.Context, config option.Config, etcdcli *clientv3.Client) (*NodeManager, error) {
	nm := &NodeManager{
		config: config,
	}
	ipManager, err := CreateIPManager(ctx, config, etcdcli)
	if err != nil {
		return nil, err
	}
	nm.ipManager = ipManager
	return nm, nil
}

// Start
func (n *NodeManager) Start() error {
	if err := n.ipManager.Start(); err != nil {
		return err
	}
	if ok := n.checkGatewayPort(); !ok {
		return fmt.Errorf("Check gateway node port failure")
	}
	return nil
}

// Stop
func (n *NodeManager) Stop() {
	n.ipManager.Stop()
}

func (n *NodeManager) checkGatewayPort() bool {
	ports := []uint32{
		uint32(n.config.ListenPorts.Health),
		uint32(n.config.ListenPorts.HTTP),
		uint32(n.config.ListenPorts.HTTPS),
		uint32(n.config.ListenPorts.Status),
		uint32(n.config.ListenPorts.Stream),
	}
	return n.CheckPortAvailable("tcp", ports...)
}

//CheckPortAvailable checks whether the specified port is available
func (n *NodeManager) CheckPortAvailable(protocol string, ports ...uint32) bool {
	if protocol == "" {
		protocol = "tcp"
	}
	timeout := time.Second * 3
	for _, port := range ports {
		c, _ := net.DialTimeout(protocol, fmt.Sprintf("0.0.0.0:%d", port), timeout)
		if c != nil {
			logrus.Errorf("Gateway must need listen port %d, but it has been uesd.", port)
			return false
		}
	}
	return true
}

//IPManager
func (n *NodeManager) IPManager() IPManager {
	return n.ipManager
}
