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

package thirdparty

import (
	"github.com/eapache/channels"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/worker/appm/thirdparty/discovery"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
)

// Interacter is the interface that wraps the required methods to interact
// with DB or service registry that holds the endpoints information.
type Interacter interface {
	List() ([]*v1.RbdEndpoint, error)
	// if endpoints type is static, do nothing.
	// if endpoints type is dynamic, watch the changes in endpoints.
	Watch()
}

// NewInteracter creates a new Interacter.
func NewInteracter(sid string, updateCh *channels.RingChannel, stopCh chan struct{}) Interacter {
	cfg, err := db.GetManager().ThirdPartySvcDiscoveryCfgDao().GetByServiceID(sid)
	if err != nil {
		logrus.Warningf("ServiceID: %s;error getting third-party discovery configuration"+
			": %s", sid, err.Error())
	}
	if err == nil && cfg != nil {
		d := &dynamic{
			cfg:      cfg,
			updateCh: updateCh,
			stopCh:   stopCh,
		}
		return d
	}
	return &static{
		sid: sid,
	}
}

// NewStaticInteracter creates a new static interacter.
func NewStaticInteracter(sid string) Interacter {
	return &static{
		sid: sid,
	}
}

type static struct {
	sid string
}

func (s *static) List() ([]*v1.RbdEndpoint, error) {
	eps, err := db.GetManager().EndpointsDao().List(s.sid)
	if err != nil {
		return nil, err
	}
	var res []*v1.RbdEndpoint
	for _, ep := range eps {
		res = append(res, &v1.RbdEndpoint{
			UUID:     ep.UUID,
			Sid:      ep.ServiceID,
			IP:       ep.IP,
			Port:     ep.Port,
			IsOnline: *ep.IsOnline,
		})
	}
	return res, nil
}

func (s *static) Watch() {
	// do nothing
}

// NewDynamicInteracter creates a new static interacter.
func NewDynamicInteracter(sid string, updateCh *channels.RingChannel, stopCh chan struct{}) Interacter {
	cfg, err := db.GetManager().ThirdPartySvcDiscoveryCfgDao().GetByServiceID(sid)
	if err != nil {
		logrus.Warningf("ServiceID: %s;error getting third-party discovery configuration"+
			": %s", sid, err.Error())
		return nil
	}
	if cfg == nil {
		return nil
	}
	d := &dynamic{
		cfg:      cfg,
		updateCh: updateCh,
		stopCh:   stopCh,
	}
	return d

}

type dynamic struct {
	cfg *model.ThirdPartySvcDiscoveryCfg

	updateCh *channels.RingChannel
	stopCh   chan struct{}
}

func (d *dynamic) List() ([]*v1.RbdEndpoint, error) {
	discoverier, err := discovery.NewDiscoverier(d.cfg, d.updateCh, d.stopCh)
	if err != nil {
		return nil, err
	}
	if err := discoverier.Connect(); err != nil {
		return nil, err
	}
	defer discoverier.Close()
	return discoverier.Fetch()
}

func (d *dynamic) Watch() {
	discoverier, err := discovery.NewDiscoverier(d.cfg, d.updateCh, d.stopCh)
	if err != nil {
		logrus.Warningf("error creating discoverier: %s", err.Error())
		return
	}
	if err := discoverier.Connect(); err != nil {
		logrus.Warningf("error connecting service discovery center: %s", err.Error())
		return
	}
	defer discoverier.Close()
	discoverier.Watch()
}
