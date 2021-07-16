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

package appm

import (
	"github.com/eapache/channels"
	opt "github.com/gridworkz/kato/cmd/worker/option"
	"github.com/gridworkz/kato/worker/appm/prober"
	"github.com/gridworkz/kato/worker/appm/store"
	"github.com/gridworkz/kato/worker/appm/thirdparty"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

// NewAPPMController creates a new appm controller.
func NewAPPMController(clientset kubernetes.Interface,
	store store.Storer,
	startCh *channels.RingChannel,
	updateCh *channels.RingChannel,
	probeCh *channels.RingChannel) *Controller {
	c := &Controller{
		store:    store,
		updateCh: updateCh,
		startCh:  startCh,
		probeCh:  probeCh,
		stopCh:   make(chan struct{}),
	}
	// create prober first, then thirdparty
	c.prober = prober.NewProber(c.store, c.probeCh, c.updateCh)
	c.thirdparty = thirdparty.NewThirdPartier(clientset, c.store, c.startCh, c.updateCh, c.stopCh, c.prober)
	return c
}

// Controller describes a new appm controller.
type Controller struct {
	cfg opt.Config

	store      store.Storer
	thirdparty thirdparty.ThirdPartier
	prober     prober.Prober

	startCh  *channels.RingChannel
	updateCh *channels.RingChannel
	probeCh  *channels.RingChannel
	stopCh   chan struct{}
}

// Start starts appm controller
func (c *Controller) Start() error {
	c.thirdparty.Start()
	c.prober.Start()
	logrus.Debugf("start thirdparty appm manager success")
	return nil
}

// Stop stops appm controller.
func (c *Controller) Stop() {
	close(c.stopCh)
	c.prober.Stop()
}
