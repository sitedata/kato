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

	v3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
)

// Queue implements a multi-reader, multi-writer distributed queue.
type Queue struct {
	client *v3.Client
	ctx    context.Context

	keyPrefix string
}

// NewQueue
func NewQueue(ctx context.Context, client *v3.Client, keyPrefix string) *Queue {
	return &Queue{client, ctx, keyPrefix}
}

// Enqueue
func (q *Queue) Enqueue(val string) error {
	_, err := newUniqueKV(q.ctx, q.client, q.keyPrefix, val)
	return err
}

// Dequeue returns Enqueue()'d elements in FIFO order. If the
// queue is empty, Dequeue blocks until elements are available.
func (q *Queue) Dequeue() (string, error) {
	for {
		// TODO: fewer round trips by fetching more than one key
		resp, err := q.client.Get(q.ctx, q.keyPrefix, v3.WithFirstRev()...)
		if err != nil {
			return "", err
		}

		kv, err := claimFirstKey(q.ctx, q.client, resp.Kvs)
		if err != nil {
			return "", err
		} else if kv != nil {
			return string(kv.Value), nil
		} else if resp.More {
			// missed some items, retry to read in more
			return q.Dequeue()
		}

		// nothing yet; wait on elements
		ev, err := WaitPrefixEvents(
			q.client,
			q.keyPrefix,
			resp.Header.Revision,
			[]mvccpb.Event_EventType{mvccpb.PUT})
		if err != nil {
			if err == ErrNoUpdateForLongTime {
				continue
			}
			return "", err
		}
		if ev == nil {
			return "", fmt.Errorf("event is nil")
		}
		if ev.Kv == nil {
			return "", fmt.Errorf("event key value is nil")
		}
		ok, err := deleteRevKey(q.ctx, q.client, string(ev.Kv.Key), ev.Kv.ModRevision)
		if err != nil {
			return "", err
		} else if !ok {
			return q.Dequeue()
		}
		return string(ev.Kv.Value), err
	}
}
