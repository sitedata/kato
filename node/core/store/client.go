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

package store

import (
	"errors"
	"strings"
	"time"

	client "github.com/coreos/etcd/clientv3"

	conf "github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/utils"

	"context"

	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/sirupsen/logrus"
)

var (
	//DefalutClient etcd client
	DefalutClient *Client
)

//Client - etcd client
type Client struct {
	*client.Client
	reqTimeout time.Duration
}

//NewClient
func NewClient(ctx context.Context, cfg *conf.Conf, etcdClientArgs *etcdutil.ClientArgs) (err error) {
	cli, err := etcdutil.NewClient(ctx, etcdClientArgs)
	if err != nil {
		return
	}
	if cfg.ReqTimeout < 3 {
		cfg.ReqTimeout = 3
	}
	c := &Client{
		Client:     cli,
		reqTimeout: time.Duration(cfg.ReqTimeout) * time.Second,
	}
	logrus.Infof("init etcd client, endpoint is:%v", cfg.EtcdEndpoints)
	DefalutClient = c
	return
}

//ErrKeyExists
var ErrKeyExists = errors.New("key already exists")

// Post attempts to create the given key, only succeeding if the key did
// not yet exist.
func (c *Client) Post(key, val string, opts ...client.OpOption) (*client.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	cmp := client.Compare(client.Version(key), "=", 0)
	req := client.OpPut(key, val, opts...)
	txnresp, err := c.Client.Txn(ctx).If(cmp).Then(req).Commit()
	if err != nil {
		return nil, err
	}
	if !txnresp.Succeeded {
		return nil, ErrKeyExists
	}
	return txnresp.OpResponse().Put(), nil
}

//Put etcd v3
func (c *Client) Put(key, val string, opts ...client.OpOption) (*client.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.Put(ctx, key, val, opts...)
}

//NewRunnable
func (c *Client) NewRunnable(key, val string, opts ...client.OpOption) (*client.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.Put(ctx, key, val, opts...)
}

//DelRunnable
func (c *Client) DelRunnable(key string, opts ...client.OpOption) (*client.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.Delete(ctx, key, opts...)
}

//PutWithModRev
func (c *Client) PutWithModRev(key, val string, rev int64) (*client.PutResponse, error) {
	if rev == 0 {
		return c.Put(key, val)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	tresp, err := DefalutClient.Txn(ctx).
		If(client.Compare(client.ModRevision(key), "=", rev)).
		Then(client.OpPut(key, val)).
		Commit()
	cancel()
	if err != nil {
		return nil, err
	}

	if !tresp.Succeeded {
		return nil, utils.ErrValueMayChanged
	}

	resp := client.PutResponse(*tresp.Responses[0].GetResponsePut())
	return &resp, nil
}

//IsRunnable
func (c *Client) IsRunnable(key string, opts ...client.OpOption) bool {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	resp, err := c.Client.Get(ctx, key, opts...)
	if err != nil {
		logrus.Infof("get key %s from etcd failed ,details %s", key, err.Error())
		return false
	}
	if resp.Count <= 0 {
		logrus.Infof("get nothing from etcd by key %s", key)
		return false
	}
	return true
}

//Get
func (c *Client) Get(key string, opts ...client.OpOption) (*client.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.Get(ctx, key, opts...)
}

//Delete v3 etcd
func (c *Client) Delete(key string, opts ...client.OpOption) (*client.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.Delete(ctx, key, opts...)
}

//Watch etcd v3
func (c *Client) Watch(key string, opts ...client.OpOption) client.WatchChan {
	return c.Client.Watch(context.Background(), key, opts...)
}

//WatchByCtx - watch by ctx
func (c *Client) WatchByCtx(ctx context.Context, key string, opts ...client.OpOption) client.WatchChan {
	return c.Client.Watch(ctx, key, opts...)
}

//KeepAliveOnce etcd v3
func (c *Client) KeepAliveOnce(id client.LeaseID) (*client.LeaseKeepAliveResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.KeepAliveOnce(ctx, id)
}

//GetLock
func (c *Client) GetLock(key string, id client.LeaseID) (bool, error) {
	key = conf.Config.LockPath + key
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	resp, err := DefalutClient.Txn(ctx).
		If(client.Compare(client.CreateRevision(key), "=", 0)).
		Then(client.OpPut(key, "", client.WithLease(id))).
		Commit()
	cancel()

	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

//DelLock
func (c *Client) DelLock(key string) error {
	_, err := c.Delete(conf.Config.LockPath + key)
	return err
}

//Grant etcd v3
func (c *Client) Grant(ttl int64) (*client.LeaseGrantResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.Grant(ctx, ttl)
}

//IsValidAsKeyPath
func IsValidAsKeyPath(s string) bool {
	return strings.IndexByte(s, '/') == -1
}
