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
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	client "github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/util"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/sirupsen/logrus"
)

//KeepAlive - Service registration
type KeepAlive struct {
	cancel        context.CancelFunc
	EtcdClentArgs *etcdutil.ClientArgs
	ServerName    string
	HostName      string
	Endpoint      string
	TTL           int64
	LID           clientv3.LeaseID
	Done          chan struct{}
	etcdClient    *client.Client
}

//CreateKeepAlive - create keepalive for server
func CreateKeepAlive(etcdClientArgs *etcdutil.ClientArgs, serverName string, hostName string, HostIP string, Port int) (*KeepAlive, error) {
	if serverName == "" || Port == 0 {
		return nil, fmt.Errorf("servername or serverport can not be empty")
	}
	if hostName == "" {
		var err error
		hostName, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	if HostIP == "" {
		ip, err := util.LocalIP()
		if err != nil {
			logrus.Errorf("get ip failed,details %s", err.Error())
			return nil, err
		}
		HostIP = ip.String()
	}
	return &KeepAlive{
		EtcdClentArgs: etcdClientArgs,
		ServerName:    serverName,
		HostName:      hostName,
		Endpoint:      fmt.Sprintf("%s:%d", HostIP, Port),
		TTL:           10,
		Done:          make(chan struct{}),
	}, nil
}

//Start
func (k *KeepAlive) Start() error {
	duration := time.Duration(k.TTL) * time.Second
	timer := time.NewTimer(duration)
	ctx, cancel := context.WithCancel(context.Background())
	etcdclient, err := etcdutil.NewClient(ctx, k.EtcdClentArgs)
	if err != nil {
		return err
	}
	k.etcdClient = etcdclient
	k.cancel = cancel
	go func() {
		for {
			select {
			case <-k.Done:
				return
			case <-timer.C:
				if k.LID > 0 {
					func() {
						ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
						defer cancel()
						_, err := k.etcdClient.KeepAliveOnce(ctx, k.LID)
						if err == nil {
							timer.Reset(duration)
							return
						}
						logrus.Warnf("%s lid[%x] keepAlive err: %s, try to reset...", k.Endpoint, k.LID, err.Error())
						k.LID = 0
						timer.Reset(duration)
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
	return fmt.Sprintf("/traefik/backends/%s/servers/%s/url", k.ServerName, k.HostName)
}

func (k *KeepAlive) reg() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := k.etcdClient.Grant(ctx, k.TTL+3)
	if err != nil {
		return err
	}
	if _, err := k.etcdClient.Put(ctx,
		k.etcdKey(),
		k.Endpoint,
		clientv3.WithLease(resp.ID)); err != nil {
		return err
	}
	logrus.Infof("Register a %s server endpoint %s to cluster", k.ServerName, k.Endpoint)
	k.LID = resp.ID
	return nil
}

//Stop
func (k *KeepAlive) Stop() error {
	close(k.Done)
	defer k.cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if _, err := k.etcdClient.Delete(ctx, k.etcdKey()); err != nil {
		return err
	}
	logrus.Infof("cancel %s server endpoint %s from etcd", k.ServerName, k.Endpoint)
	return nil
}
