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
	"fmt"
	"strings"
	"time"

	"github.com/gridworkz/kato/util"
	"github.com/gridworkz/kato/worker/server/pb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//GetServicePods get service pods list
func (a *AppRuntimeSyncClient) GetServicePods(serviceID string) (*pb.ServiceAppPodList, error) {
	ctx, cancel := context.WithTimeout(a.ctx, time.Second*5)
	defer cancel()
	return a.AppRuntimeSyncClient.GetAppPods(ctx, &pb.ServiceRequest{ServiceId: serviceID})
}

//GetMultiServicePods get multi service pods list
func (a *AppRuntimeSyncClient) GetMultiServicePods(serviceIDs []string) (*pb.MultiServiceAppPodList, error) {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		defer util.Elapsed(fmt.Sprintf("[AppRuntimeSyncClient] [GetMultiServicePods] component nums: %d", len(serviceIDs)))()
	}

	ctx, cancel := context.WithTimeout(a.ctx, time.Second*5)
	defer cancel()
	return a.AppRuntimeSyncClient.GetMultiAppPods(ctx, &pb.ServicesRequest{ServiceIds: strings.Join(serviceIDs, ",")})
}

// GetComponentPodNums -
func (a *AppRuntimeSyncClient) GetComponentPodNums(ctx context.Context, componentIDs []string) (map[string]int32, error) {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		defer util.Elapsed(fmt.Sprintf("[AppRuntimeSyncClient] get component pod nums: %d", len(componentIDs)))()
	}

	res, err := a.AppRuntimeSyncClient.GetComponentPodNums(ctx, &pb.ServicesRequest{ServiceIds: strings.Join(componentIDs, ",")})
	if err != nil {
		return nil, errors.Wrap(err, "get component pod nums")
	}

	return res.PodNums, nil
}

// GetPodDetail -
func (a *AppRuntimeSyncClient) GetPodDetail(sid, name string) (*pb.PodDetail, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return a.AppRuntimeSyncClient.GetPodDetail(ctx, &pb.GetPodDetailReq{
		Sid:     sid,
		PodName: name,
	})
}
