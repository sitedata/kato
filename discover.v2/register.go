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

package discover

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	client "github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/util"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	grpcutil "github.com/gridworkz/kato/util/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/naming"
)

//KeepAlive
type KeepAlive struct {
	cancel         context.CancelFunc
	EtcdClientArgs *etcdutil.ClientArgs
	ServerName     string
	HostName       string
	Endpoint       string
	TTL            int64
	LID            clientv3.LeaseID
	Done           chan struct{}
	etcdClient     *client.Client
	gRPCResolver   *grpcutil.GRPCResolver
	once           sync.Once
}

//CreateKeepAlive create keepalive for server
func CreateKeepAlive(etcdClientArgs *etcdutil.ClientArgs, ServerName string, Protocol string, HostIP string, Port int) (*KeepAlive, error) {
	if ServerName == "" || Port == 0 {
		return nil, fmt.Errorf("servername or serverport can not be empty")
	}
	if HostIP == "" {
		ip, err := util.LocalIP()
		if err != nil {
			logrus.Errorf("get ip failed,details %s", err.Error())
			return nil, err
		}
		HostIP = ip.String()
	}

	ctx, cancel := context.WithCancel(context.Background())
	etcdclient, err := etcdutil.NewClient(ctx, etcdClientArgs)
	if err != nil {
		cancel()
		return nil, err
	}

	k := &KeepAlive{
		EtcdClientArgs: etcdClientArgs,
		ServerName:     ServerName,
		Endpoint:       fmt.Sprintf("%s:%d", HostIP, Port),
		TTL:            5,
		Done:           make(chan struct{}),
		etcdClient:     etcdclient,
		cancel:         cancel,
	}
	if Protocol == "" {
		k.Endpoint = fmt.Sprintf("%s:%d", HostIP, Port)
	} else {
		k.Endpoint = fmt.Sprintf("%s://%s:%d", Protocol, HostIP, Port)
	}
	return k, nil
}

//Start
func (k *KeepAlive) Start() error {
	duration := time.Duration(k.TTL) * time.Second
	timer := time.NewTimer(duration)

	go func() {
		for {
			select {
			case <-k.Done:
				return
			case <-timer.C:
				if k.LID > 0 {
					func() {
						ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
						defer cancel()
						defer timer.Reset(duration)
						_, err := k.etcdClient.KeepAliveOnce(ctx, k.LID)
						if err == nil {
							return
						}
						logrus.Warnf("%s lid[%x] keepAlive err: %s, try to reset...", k.Endpoint, k.LID, err.Error())
						k.LID = 0
					}()
				} else {
					if err := k.reg(); err != nil {
						logrus.Warnf("%s set lid err: %s, try to reset after %d seconds...", k.Endpoint, err.Error(), k.TTL)
					} else {
						logrus.Infof("%s set lid[%x] success", k.Endpoint, k.LID)
					}
					timer.Reset(duration)
				}
			}
		}
	}()
	return nil
}

func (k *KeepAlive) etcdKey() string {
	return fmt.Sprintf("/kato/discover/%s", k.ServerName)
}

func (k *KeepAlive) reg() error {
	k.gRPCResolver = &grpcutil.GRPCResolver{Client: k.etcdClient}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := k.etcdClient.Grant(ctx, k.TTL+3)
	if err != nil {
		return err
	}
	if err := k.gRPCResolver.Update(ctx, k.etcdKey(), naming.Update{Op: naming.Add, Addr: k.Endpoint}, clientv3.WithLease(resp.ID)); err != nil {
		return err
	}
	logrus.Infof("Register a %s server endpoint %s to cluster", k.ServerName, k.Endpoint)
	k.LID = resp.ID
	return nil
}

//Stop
func (k *KeepAlive) Stop() {
	k.once.Do(func() {
		close(k.Done)
		k.cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if k.gRPCResolver != nil {
			if err := k.gRPCResolver.Update(ctx, k.etcdKey(), naming.Update{Op: naming.Delete, Addr: k.Endpoint}); err != nil {
				logrus.Errorf("cancel %s server endpoint %s from etcd error %s", k.ServerName, k.Endpoint, err.Error())
			} else {
				logrus.Infof("cancel %s server endpoint %s from etcd", k.ServerName, k.Endpoint)
			}
		}
	})
}
