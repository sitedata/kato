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

package distribution

import (
	"github.com/gridworkz/kato/eventlog/cluster/discover"
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/db"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func TestGetSuitableInstance(t *testing.T) {
	dis := discover.New(conf.DiscoverConf{}, logrus.WithField("Module", "Test"))
	ctx, cancel := context.WithCancel(context.Background())
	d := &Distribution{
		cancel:       cancel,
		context:      ctx,
		discover:     dis,
		updateTime:   make(map[string]time.Time),
		abnormalNode: make(map[string]int),
		log:          logrus.WithField("Module", "Test"),
	}
	d.monitorDatas = map[string]*db.MonitorData{
		"a": &db.MonitorData{
			InstanceID: "a", LogSizePeerM: 200,
		},
		"b": &db.MonitorData{
			InstanceID: "b", LogSizePeerM: 150,
		},
	}
	d.GetSuitableInstance("todo service id")
}
