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
	"strings"
	"time"

	"github.com/gridworkz/kato/db/model"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	grpcutil "github.com/gridworkz/kato/util/grpc"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/gridworkz/kato/worker/server/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//AppRuntimeSyncClient grpc client
type AppRuntimeSyncClient struct {
	pb.AppRuntimeSyncClient
	AppRuntimeSyncClientConf
	cc  *grpc.ClientConn
	ctx context.Context
}

//AppRuntimeSyncClientConf client conf
type AppRuntimeSyncClientConf struct {
	EtcdEndpoints        []string
	EtcdCaFile           string
	EtcdCertFile         string
	EtcdKeyFile          string
	DefaultServerAddress []string
}

//NewClient new client
//ctx must be cancel where client not used
func NewClient(ctx context.Context, conf AppRuntimeSyncClientConf) (*AppRuntimeSyncClient, error) {
	var arsc AppRuntimeSyncClient
	arsc.AppRuntimeSyncClientConf = conf
	arsc.ctx = ctx
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: conf.EtcdEndpoints,
		CaFile:    conf.EtcdCaFile,
		CertFile:  conf.EtcdCertFile,
		KeyFile:   conf.EtcdKeyFile,
	}
	c, err := etcdutil.NewClient(ctx, etcdClientArgs)
	if err != nil {
		return nil, err
	}
	r := &grpcutil.GRPCResolver{Client: c}
	b := grpc.RoundRobin(r)
	arsc.cc, err = grpc.DialContext(ctx, "/kato/discover/app_sync_runtime_server", grpc.WithBalancer(b), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	arsc.AppRuntimeSyncClient = pb.NewAppRuntimeSyncClient(arsc.cc)
	return &arsc, nil
}

//when watch occurred error,will exec this method
func (a *AppRuntimeSyncClient) Error(err error) {
	logrus.Errorf("discover app runtime sync server address occurred err:%s", err.Error())
}

//GetStatus get status
func (a *AppRuntimeSyncClient) GetStatus(serviceID string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	status, err := a.AppRuntimeSyncClient.GetAppStatusDeprecated(ctx, &pb.ServicesRequest{
		ServiceIds: serviceID,
	})
	if err != nil {
		return v1.UNKNOW
	}
	return status.Status[serviceID]
}

//GetStatuss get multiple app status
func (a *AppRuntimeSyncClient) GetStatuss(serviceIDs string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	status, err := a.AppRuntimeSyncClient.GetAppStatusDeprecated(ctx, &pb.ServicesRequest{
		ServiceIds: serviceIDs,
	})
	if err != nil {
		logrus.Errorf("get service status failure %s", err.Error())
		re := make(map[string]string, len(serviceIDs))
		for _, id := range strings.Split(serviceIDs, ",") {
			re[id] = v1.UNKNOW
		}
		return re
	}
	return status.Status
}

//GetAllStatus get all status
func (a *AppRuntimeSyncClient) GetAllStatus() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	status, err := a.AppRuntimeSyncClient.GetAppStatusDeprecated(ctx, &pb.ServicesRequest{
		ServiceIds: "",
	})
	if err != nil {
		return nil
	}
	return status.Status
}

//GetNeedBillingStatus get need billing status
func (a *AppRuntimeSyncClient) GetNeedBillingStatus() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	re, err := a.AppRuntimeSyncClient.GetAppStatusDeprecated(ctx, &pb.ServicesRequest{})
	if err != nil {
		return nil, err
	}
	var res = make(map[string]string)
	for k, v := range re.Status {
		if !a.IsClosedStatus(v) {
			res[k] = v
		}
	}
	return res, nil
}

//GetServiceDeployInfo get service deploy info
func (a *AppRuntimeSyncClient) GetServiceDeployInfo(serviceID string) (*pb.DeployInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	re, err := a.AppRuntimeSyncClient.GetDeployInfo(ctx, &pb.ServiceRequest{
		ServiceId: serviceID,
	})
	if err != nil {
		return nil, err
	}
	return re, nil
}

//IsClosedStatus  check status
func (a *AppRuntimeSyncClient) IsClosedStatus(curStatus string) bool {
	return curStatus == "" || curStatus == v1.BUILDEFAILURE || curStatus == v1.CLOSED || curStatus == v1.UNDEPLOY || curStatus == v1.BUILDING || curStatus == v1.UNKNOW
}

//GetTenantResource get tenant resource
func (a *AppRuntimeSyncClient) GetTenantResource(tenantID string) (*pb.TenantResource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return a.AppRuntimeSyncClient.GetTenantResource(ctx, &pb.TenantRequest{TenantId: tenantID})
}

//GetAllTenantResource get all tenant resource
func (a *AppRuntimeSyncClient) GetAllTenantResource() (*pb.TenantResourceList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return a.AppRuntimeSyncClient.GetTenantResources(ctx, &pb.Empty{})
}

// ListThirdPartyEndpoints -
func (a *AppRuntimeSyncClient) ListThirdPartyEndpoints(sid string) (*pb.ThirdPartyEndpoints, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := a.AppRuntimeSyncClient.ListThirdPartyEndpoints(ctx, &pb.ServiceRequest{
		ServiceId: sid,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AddThirdPartyEndpoint -
func (a *AppRuntimeSyncClient) AddThirdPartyEndpoint(req *model.Endpoint) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _ = a.AppRuntimeSyncClient.AddThirdPartyEndpoint(ctx, &pb.AddThirdPartyEndpointsReq{
		Uuid:     req.UUID,
		Sid:      req.ServiceID,
		Ip:       req.IP,
		Port:     int32(req.Port),
		IsOnline: *req.IsOnline,
	})
}

// UpdThirdPartyEndpoint -
func (a *AppRuntimeSyncClient) UpdThirdPartyEndpoint(req *model.Endpoint) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _ = a.AppRuntimeSyncClient.UpdThirdPartyEndpoint(ctx, &pb.UpdThirdPartyEndpointsReq{
		Uuid:     req.UUID,
		Sid:      req.ServiceID,
		Ip:       req.IP,
		Port:     int32(req.Port),
		IsOnline: *req.IsOnline,
	})
}

// DelThirdPartyEndpoint -
func (a *AppRuntimeSyncClient) DelThirdPartyEndpoint(req *model.Endpoint) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _ = a.AppRuntimeSyncClient.DelThirdPartyEndpoint(ctx, &pb.DelThirdPartyEndpointsReq{
		Uuid: req.UUID,
		Sid:  req.ServiceID,
		Ip:   req.IP,
		Port: int32(req.Port),
	})
}

// GetStorageClasses client GetStorageClasses
func (a *AppRuntimeSyncClient) GetStorageClasses() (storageclasses *pb.StorageClasses, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return a.AppRuntimeSyncClient.GetStorageClasses(ctx, &pb.Empty{})
}

// GetAppVolumeStatus get app volume status
func (a *AppRuntimeSyncClient) GetAppVolumeStatus(serviceID string) (*pb.ServiceVolumeStatusMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return a.AppRuntimeSyncClient.GetAppVolumeStatus(ctx, &pb.ServiceRequest{ServiceId: serviceID})
}

// GetAppResources -
func (a *AppRuntimeSyncClient) GetAppResources(appID string) (*pb.AppStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return a.AppRuntimeSyncClient.GetAppStatus(ctx, &pb.AppStatusReq{AppId: appID})
}
