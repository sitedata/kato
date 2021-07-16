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

// 

package initiate

import (
	"context"
	"errors"

	"github.com/gridworkz/kato/cmd/node/option"
	discover "github.com/gridworkz/kato/discover.v2"
	"github.com/gridworkz/kato/discover/config"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
)

var (
	// ErrRegistryAddressNotFound record not found error, happens when haven't find any matched registry address.
	ErrRegistryAddressNotFound = errors.New("registry address not found")
)

// HostManager is responsible for writing the resolution of the private image repository domain name to /etc/hosts.
type HostManager interface {
	Start()
}

// NewHostManager
func NewHostManager(cfg *option.Conf, discover discover.Discover) (HostManager, error) {
	hosts, err := util.NewHosts(cfg.HostsFile)
	if err != nil {
		return nil, err
	}
	callback := &hostCallback{
		cfg:   cfg,
		hosts: hosts,
	}
	return &hostManager{
		cfg:          cfg,
		discover:     discover,
		hostCallback: callback,
	}, nil
}

type hostManager struct {
	ctx          context.Context
	cfg          *option.Conf
	discover     discover.Discover
	hostCallback *hostCallback
}

func (h *hostManager) Start() {
	if h.cfg.ImageRepositoryHost == "" {
		// no need to write hosts file
		return
	}
	if h.cfg.GatewayVIP != "" {
		if err := h.hostCallback.hosts.Cleanup(); err != nil {
			logrus.Warningf("cleanup hosts file: %v", err)
			return
		}
		logrus.Infof("set hosts %s to %s", h.cfg.ImageRepositoryHost, h.cfg.GatewayVIP)
		lines := []string{
			util.StartOfSection,
			h.cfg.GatewayVIP + " " + h.cfg.ImageRepositoryHost,
			h.cfg.GatewayVIP + " " + "region.gridworkz",
			util.EndOfSection,
		}
		h.hostCallback.hosts.AddLines(lines...)
		if err := h.hostCallback.hosts.Flush(); err != nil {
			logrus.Warningf("flush hosts file: %v", err)
		}
		return
	}
	h.discover.AddProject("rbd-gateway", h.hostCallback)
}

type hostCallback struct {
	cfg   *option.Conf
	hosts util.Hosts
}

// TODO HA error
func (h *hostCallback) UpdateEndpoints(endpoints ...*config.Endpoint) {
	logrus.Info("hostCallback; update endpoints")
	if err := h.hosts.Cleanup(); err != nil {
		logrus.Warningf("cleanup hosts file: %v", err)
		return
	}

	if len(endpoints) > 0 {
		logrus.Infof("found endpints: %d; endpoint selected: %#v", len(endpoints), *endpoints[0])
		lines := []string{
			util.StartOfSection,
			endpoints[0].URL + " " + h.cfg.ImageRepositoryHost,
			endpoints[0].URL + " " + "region.gridworkz",
			util.EndOfSection,
		}
		h.hosts.AddLines(lines...)
	}

	if err := h.hosts.Flush(); err != nil {
		logrus.Warningf("flush hosts file: %v", err)
	}
}

func (h *hostCallback) Error(err error) {
	logrus.Warningf("unexpected error from host callback: %v", err)
}
