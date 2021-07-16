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

package discovery

import (
	"fmt"
	"github.com/eapache/channels"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/worker/appm/types/v1"
	"strings"
)

// EventType
type EventType string

const (
	// CreateEvent event associated with new objects in a service discovery center
	CreateEvent EventType = "CREATE"
	// UpdateEvent event associated with an object update in a service discovery center
	UpdateEvent EventType = "UPDATE"
	// DeleteEvent event associated when an object is removed from a service discovery center
	DeleteEvent EventType = "DELETE"
	// UnhealthyEvent -
	UnhealthyEvent EventType = "UNHEALTHY"
	// HealthEvent -
	HealthEvent EventType = "HEALTH"
)

// Event holds the context of an event.
type Event struct {
	Type EventType
	Obj  interface{}
}

// Discoverier is the interface that wraps the required methods to gather
// information about third-party service endpoints.
type Discoverier interface {
	Connect() error
	Close() error
	Fetch() ([]*v1.RbdEndpoint, error)
	Watch()
}

// NewDiscoverier creates a new Discoverier.
func NewDiscoverier(cfg *model.ThirdPartySvcDiscoveryCfg,
	updateCh *channels.RingChannel,
	stopCh chan struct{}) (Discoverier, error) {
	switch strings.ToLower(cfg.Type) {
	case strings.ToLower(string(model.DiscorveryTypeEtcd)):
		return NewEtcd(cfg, updateCh, stopCh), nil
	default:
		return nil, fmt.Errorf("Unsupported discovery type: %s", cfg.Type)
	}
}
