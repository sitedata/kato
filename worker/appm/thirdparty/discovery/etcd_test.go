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
	"testing"
	"time"
)

func TestEtcd_Watch(t *testing.T) {
	cfg := &model.ThirdPartySvcDiscoveryCfg{
		Type:    model.DiscorveryTypeEtcd.String(),
		Servers: "http://127.0.0.1:2379",
		Key:     "/foobar/eps",
	}
	updateCh := channels.NewRingChannel(1024)
	stopCh := make(chan struct{})
	defer close(stopCh)

	go func() {
		for {
			select {
			case event := <-updateCh.Out():
				fmt.Printf("%+v", event)
			case <-stopCh:
				break
			}
		}
	}()

	etcd := NewEtcd(cfg, updateCh, stopCh)
	if err := etcd.Connect(); err != nil {
		t.Fatalf("error connecting etcd: %v", err)
	}
	defer etcd.Close()

	etcd.Watch()

	time.Sleep(10 * time.Second)
}
