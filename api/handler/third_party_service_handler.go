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
	"sort"
	"strconv"
	"strings"

	"github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	"github.com/gridworkz/kato/worker/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ThirdPartyServiceHanlder handles business logic for all third-party services
type ThirdPartyServiceHanlder struct {
	logger    *logrus.Entry
	dbmanager db.Manager
	statusCli *client.AppRuntimeSyncClient
}

// Create3rdPartySvcHandler creates a new *ThirdPartyServiceHanlder.
func Create3rdPartySvcHandler(dbmanager db.Manager, statusCli *client.AppRuntimeSyncClient) *ThirdPartyServiceHanlder {
	return &ThirdPartyServiceHanlder{
		logger:    logrus.WithField("WHO", "ThirdPartyServiceHanlder"),
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
func (t *ThirdPartyServiceHanlder) ListEndpoints(componentID string) ([]*model.ThirdEndpoint, error) {
	logger := t.logger.WithField("Method", "ListEndpoints").
		WithField("ComponentID", componentID)

	runtimeEndpoints, err := t.listRuntimeEndpoints(componentID)
	if err != nil {
		logger.Warning(err.Error())
	}

	staticEndpoints, err := t.listStaticEndpoints(componentID)
	if err != nil {
		staticEndpoints = map[string]*model.ThirdEndpoint{}
		logger.Warning(err.Error())
	}

	// Merge runtimeEndpoints with staticEndpoints
	for _, ep := range runtimeEndpoints {
		sep, ok := staticEndpoints[ep.EpID]
		if !ok {
			continue
		}
		ep.IsStatic = sep.IsStatic
		ep.Address = sep.Address
		delete(staticEndpoints, ep.EpID)
	}

	// Add offline static endpoints
	for _, ep := range staticEndpoints {
		runtimeEndpoints = append(runtimeEndpoints, ep)
	}

	sort.Sort(model.ThirdEndpoints(runtimeEndpoints))
	return runtimeEndpoints, nil
}

func (t *ThirdPartyServiceHanlder) listRuntimeEndpoints(componentID string) ([]*model.ThirdEndpoint, error) {
	runtimeEndpoints, err := t.statusCli.ListThirdPartyEndpoints(componentID)
	if err != nil {
		return nil, errors.Wrap(err, "list runtime third endpoints")
	}

	var endpoints []*model.ThirdEndpoint
	for _, item := range runtimeEndpoints.Items {
		endpoints = append(endpoints, &model.ThirdEndpoint{
			EpID:    item.Name,
			Address: item.Address,
			Status:  item.Status,
		})
	}
	return endpoints, nil
}

func (t *ThirdPartyServiceHanlder) listStaticEndpoints(componentID string) (map[string]*model.ThirdEndpoint, error) {
	staticEndpoints, err := t.dbmanager.EndpointsDao().List(componentID)
	if err != nil {
		return nil, errors.Wrap(err, "list static endpoints")
	}

	endpoints := make(map[string]*model.ThirdEndpoint)
	for _, item := range staticEndpoints {
		address := func(ip string, p int) string {
			if p != 0 {
				return fmt.Sprintf("%s:%d", ip, p)
			}
			return ip
		}(item.IP, item.Port)
		endpoints[item.UUID] = &model.ThirdEndpoint{
			EpID:     item.UUID,
			Address:  address,
			Status:   "-",
			IsStatic: true,
		}
	}
	return endpoints, nil
}
