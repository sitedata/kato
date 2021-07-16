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

package watch

import (
	"path"

	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type watchAPI struct {
	client *clientv3.Client
	// getOpts contains additional options that should be passed
	// to all Get() calls.
	getOps     []clientv3.OpOption
	pathPrefix string
	watcher    *watcher
}

// New returns an etcd3 implementation of Watch.
func New(c *clientv3.Client, prefix string) Watch {
	return newWatchAPI(c, true, prefix)
}

// NewWithNoQuorumRead returns etcd3 implementation of storage.Interface
// where Get operations don't require quorum read.
func NewWithNoQuorumRead(c *clientv3.Client, prefix string) Watch {
	return newWatchAPI(c, false, prefix)
}

func newWatchAPI(c *clientv3.Client, quorumRead bool, prefix string) *watchAPI {

	result := &watchAPI{
		client: c,
		// for compatibility with etcd2 impl.
		// no-op for default prefix of '/registry'.
		// keeps compatibility with etcd2 impl for custom prefixes that don't start with '/'
		pathPrefix: prefix,
		watcher:    newWatcher(c),
	}
	if !quorumRead {
		// In case of non-quorum reads, we can set WithSerializable()
		// options for all Get operations.
		result.getOps = append(result.getOps, clientv3.WithSerializable())
	}
	return result
}

// Watch implements storage.Interface.Watch.
func (s *watchAPI) Watch(ctx context.Context, key string, resourceVersion string) (Interface, error) {
	return s.watch(ctx, key, resourceVersion, false)
}

// WatchList implements storage.Interface.WatchList.
func (s *watchAPI) WatchList(ctx context.Context, key string, resourceVersion string) (Interface, error) {
	return s.watch(ctx, key, resourceVersion, true)
}

func (s *watchAPI) watch(ctx context.Context, key string, rv string, recursive bool) (Interface, error) {
	rev, err := ParseWatchResourceVersion(rv)
	if err != nil {
		return nil, err
	}
	key = path.Join(s.pathPrefix, key)
	return s.watcher.Watch(ctx, key, int64(rev), recursive)
}

// ttlOpts returns client options based on given ttl.
// ttl: if ttl is non-zero, it will attach the key to a lease with ttl of roughly the same length
func (s *watchAPI) ttlOpts(ctx context.Context, ttl int64) ([]clientv3.OpOption, error) {
	if ttl == 0 {
		return nil, nil
	}
	// TODO: one lease per ttl key is expensive. Based on current use case, we can have a long window to
	// put keys within into same lease. We shall benchmark this and optimize the performance.
	lcr, err := s.client.Lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	return []clientv3.OpOption{clientv3.WithLease(clientv3.LeaseID(lcr.ID))}, nil
}
