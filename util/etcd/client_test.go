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
	"path"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"golang.org/x/net/context"
)

func TestNewETCDClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	basePath := "/Users/gridworkz/Downloads/etcdcas"
	caPath := path.Join(basePath, "etcdca.crt")
	certPath := path.Join(basePath, "apiserver-etcd-client.crt")
	keyPath := path.Join(basePath, "apiserver-etcd-client.key")
	// client, err := NewETCDClient(ctx, 10*time.Second, []string{"192.168.2.203:2379"}, "", "", "") // connection no tls success
	/**
	curl --cacert etcdca.crt  --cert apiserver-etcd-client.crt --key apiserver-etcd-client.key -L https://192.168.2.63:2379/v2/keys/foo -XGET
	cacert Specify the root certificate of the issuing authority used by the server, so you need to use the root certificate of the etcd issuing authority，It is not the root certificate of the issuing authority of kubernetes. The file path is /etc/kubernetes/pki/etcd/ca.crt
	cert Specify the client certificate, which is used here certificate of kube-apiserver，certificate, the file path is: /etc/kubernetes/pki/apiserver-etcd-client.crt，you can also use etcd's node certificate /etc/kubernetes/pki/etcd/peer.crt
	cert Specify the client certificate key, which is used here kube-apiserver key of the certificate，the file path is: /etc/kubernetes/pki/apiserver-etcd-client.key
	*/
	clientArgs := ClientArgs{
		Endpoints: []string{"192.168.2.63:2379"},
		CaFile:    caPath,
		CertFile:  certPath,
		KeyFile:   keyPath,
	}
	client, err := NewClient(ctx, &clientArgs)
	if err != nil {
		t.Fatal("create client error: ", err)
	}
	resp, err := client.Get(ctx, "/foo")
	if err != nil {
		t.Fatal("get key error", err)
	}
	t.Logf("resp is : %+v", resp)
	time.Sleep(30)
}

func TestEtcd(t *testing.T) {
	// test etcd retry connection
	fmt.Println("yes")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	etcdClientArgs := &ClientArgs{
		Endpoints: []string{"http://127.0.0.1:2359"},
	}
	etcdcli, err := NewClient(ctx, etcdClientArgs)
	if err != nil {
		logrus.Errorf("create etcd client v3 error, %v", err)
		t.Fatal(err)
	}
	memberList, err := etcdcli.MemberList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(memberList.Members) == 0 {
		fmt.Println("no members")
		return
	}
	t.Logf("members is: %s", memberList.Members[0].Name)
}
