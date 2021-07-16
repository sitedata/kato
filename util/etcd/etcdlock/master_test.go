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

package etcdlock

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestMasterLock(t *testing.T) {
	master1, err := CreateMasterLock(nil, "/kato/appruntimesyncmaster", "127.0.0.1:1", 10)
	if err != nil {
		t.Fatal(err)
	}
	master1.Start()
	defer master1.Stop()
	master2, err := CreateMasterLock(nil, "/kato/appruntimesyncmaster", "127.0.0.1:2", 10)
	if err != nil {
		t.Fatal(err)
	}
	master2.Start()
	defer master2.Stop()
	logrus.Info("start receive event")
	for {
		select {
		case event := <-master1.EventsChan():
			logrus.Info("master1", event, time.Now())
			//master1.Stop()
			if event.Type == MasterDeleted {
				logrus.Info("delete")
				return
			}
		case event := <-master2.EventsChan():
			logrus.Info("master2", event, time.Now())
			//master2.Stop()
			if event.Type == MasterDeleted {
				logrus.Info("delete")
				return
			}
		}
	}
}
