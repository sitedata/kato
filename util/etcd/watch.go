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

package etcd

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//ErrNoUpdateForLongTime no update for a long time , can re-observe synchronous data
var ErrNoUpdateForLongTime = fmt.Errorf("not updated for a long time")

//WaitPrefixEvents
func WaitPrefixEvents(c *clientv3.Client, prefix string, rev int64, evs []mvccpb.Event_EventType) (*clientv3.Event, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logrus.Debug("start watch message from etcd queue")
	wc := clientv3.NewWatcher(c).Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(rev))
	if wc == nil {
		return nil, ErrNoWatcher
	}
	event := waitEvents(wc, evs)
	if event != nil {
		return event, nil
	}
	logrus.Debug("queue watcher sync, because of not updated for a long time")
	return nil, ErrNoUpdateForLongTime
}

//waitEvents this will return nil
func waitEvents(wc clientv3.WatchChan, evs []mvccpb.Event_EventType) *clientv3.Event {
	i := 0
	timer := time.NewTimer(time.Second * 30)
	defer timer.Stop()
	for {
		select {
		case wresp := <-wc:
			if wresp.Err() != nil {
				logrus.Errorf("watch event failure %s", wresp.Err().Error())
				return nil
			}
			if len(wresp.Events) == 0 {
				return nil
			}
			for _, ev := range wresp.Events {
				if ev.Type == evs[i] {
					i++
					if i == len(evs) {
						return ev
					}
				}
			}
		case <-timer.C:
			return nil
		}
	}
}
