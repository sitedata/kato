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

package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/worker/client"
)

// ThirdPartyServiceHanlder handles business logic for all third-party services
type ThirdPartyServiceHanlder struct {
	dbmanager db.Manager
	statusCli *client.AppRuntimeSyncClient
}

// Create3rdPartySvcHandler creates a new *ThirdPartyServiceHanlder.
func Create3rdPartySvcHandler(dbmanager db.Manager, statusCli *client.AppRuntimeSyncClient) *ThirdPartyServiceHanlder {
	return &ThirdPartyServiceHanlder{
		dbmanager: dbmanager,
		statusCli: statusCli,
	}
}

// AddEndpoints adds endpoints for third-party service.
func (t *ThirdPartyServiceHanlder) AddEndpoints(sid string, d *model.AddEndpiontsReq) error {
	address, port := convertAddressPort(d.Address)
	if port == 0 {
		//set default port by service port
		ports, _ := t.dbmanager.TenantServicesPortDao().GetPortsByServiceID(sid)
		if len(ports) > 0 {
			port = ports[0].ContainerPort
		}
	}
	ep := &dbmodel.Endpoint{
		UUID:      util.NewUUID(),
		ServiceID: sid,
		IP:        address,
		Port:      port,
		IsOnline:  &d.IsOnline,
	}
	if err := t.dbmanager.EndpointsDao().AddModel(ep); err != nil {
		return err
	}

	logrus.Debugf("add new endpoint[address: %s, port: %d]", address, port)
	t.statusCli.AddThirdPartyEndpoint(ep)
	return nil
}

// UpdEndpoints updates endpoints for third-party service.
func (t *ThirdPartyServiceHanlder) UpdEndpoints(d *model.UpdEndpiontsReq) error {
	ep, err := t.dbmanager.EndpointsDao().GetByUUID(d.EpID)
	if err != nil {
		logrus.Warningf("EpID: %s; error getting endpoints: %v", d.EpID, err)
		return err
	}
	if d.Address != "" {
		address, port := convertAddressPort(d.Address)
		ep.IP = address
		ep.Port = port
	}
	ep.IsOnline = &d.IsOnline
	if err := t.dbmanager.EndpointsDao().UpdateModel(ep); err != nil {
		return err
	}

	t.statusCli.UpdThirdPartyEndpoint(ep)

	return nil
}

func convertAddressPort(s string) (address string, port int) {
	prefix := ""
	if strings.HasPrefix(s, "https://") {
		s = strings.Split(s, "https://")[1]
		prefix = "https://"
	}
	if strings.HasPrefix(s, "http://") {
		s = strings.Split(s, "http://")[1]
		prefix = "http://"
	}

	if strings.Contains(s, ":") {
		sp := strings.Split(s, ":")
		address = prefix + sp[0]
		port, _ = strconv.Atoi(sp[1])
	} else {
		address = prefix + s
	}

	return address, port
}

// DelEndpoints deletes endpoints for third-party service.
func (t *ThirdPartyServiceHanlder) DelEndpoints(epid, sid string) error {
	ep, err := t.dbmanager.EndpointsDao().GetByUUID(epid)
	if err != nil {
		logrus.Warningf("EpID: %s; error getting endpoints: %v", epid, err)
		return err
	}
	if err := t.dbmanager.EndpointsDao().DelByUUID(epid); err != nil {
		return err
	}
	t.statusCli.DelThirdPartyEndpoint(ep)

	return nil
}

// ListEndpoints lists third-party service endpoints.
func (t *ThirdPartyServiceHanlder) ListEndpoints(sid string) ([]*model.EndpointResp, error) {
	endpoints, err := t.dbmanager.EndpointsDao().List(sid)
	if err != nil {
		logrus.Warningf("ServiceID: %s; error listing endpoints from db; %v", sid, err)
	}
	m := make(map[string]*model.EndpointResp)
	for _, item := range endpoints {
		ep := &model.EndpointResp{
			EpID: item.UUID,
			Address: func(ip string, p int) string {
				if p != 0 {
					return fmt.Sprintf("%s:%d", ip, p)
				}
				return ip
			}(item.IP, item.Port),
			Status:   "-",
			IsOnline: false,
			IsStatic: true,
		}
		m[ep.Address] = ep
	}
	thirdPartyEndpoints, err := t.statusCli.ListThirdPartyEndpoints(sid)
	if err != nil {
		logrus.Warningf("ServiceID: %s; grpc; error listing third-party endpoints: %v", sid, err)
		return nil, err
	}
	if thirdPartyEndpoints != nil && thirdPartyEndpoints.Obj != nil {
		for _, item := range thirdPartyEndpoints.Obj {
			ep := m[fmt.Sprintf("%s:%d", item.Ip, item.Port)]
			if ep != nil {
				ep.IsOnline = true
				ep.Status = item.Status
				continue
			}
			rep := &model.EndpointResp{
				EpID:     item.Uuid,
				Address:  item.Ip,
				Status:   item.Status,
				IsOnline: true,
				IsStatic: false,
			}
			m[rep.Address] = rep
		}
	}
	var res []*model.EndpointResp
	for _, item := range m {
		res = append(res, item)
	}
	return res, nil
}
