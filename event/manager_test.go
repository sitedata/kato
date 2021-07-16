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

package event

import (
	"testing"
	"time"

	etcdutil "github.com/gridworkz/kato/util/etcd"
)

func TestLogger(t *testing.T) {
	err := NewManager(EventConfig{
		EventLogServers: []string{"192.168.195.1:6366"},
		DiscoverArgs:    &etcdutil.ClientArgs{Endpoints: []string{"192.168.195.1:2379"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer GetManager().Close()
	time.Sleep(time.Second * 3)
	for i := 0; i < 500; i++ {
		GetManager().GetLogger("qwdawdasdasasfafa").Info("hello word", nil)
		GetManager().GetLogger("asdasdasdasdads").Debug("hello word", nil)
		GetManager().GetLogger("1234124124124").Error("hello word", nil)
		time.Sleep(time.Millisecond * 1)
	}
	select {}
}
