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

package client

import (
	"context"
	"testing"
	"time"

	"github.com/gridworkz/kato/worker/server/pb"
)

func TestGetAppStatus(t *testing.T) {
	client, err := NewClient(context.Background(), AppRuntimeSyncClientConf{
		EtcdEndpoints: []string{"192.168.3.252:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	status, err := client.GetAppStatus(context.Background(), &pb.AppStatusReq{
		AppId: "43eaae441859eda35b02075d37d83589",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(status)
}

func TestGetPodDetail(t *testing.T) {
	client, err := NewClient(context.Background(), AppRuntimeSyncClientConf{
		EtcdEndpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		t.Fatalf("error creating grpc client: %v", err)
	}

	time.Sleep(3 * time.Second)

	sid := "bc2e153243e4917acaf0c6088789fa12"
	podname := "bc2e153243e4917acaf0c6088789fa12-deployment-6cc445949b-2dh25"
	detail, err := client.GetPodDetail(sid, podname)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("pod detail: %+v", detail)
}
