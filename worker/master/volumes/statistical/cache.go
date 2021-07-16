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

package statistical

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//DiskCache disk asynchronous statistics
type DiskCache struct {
	cache []struct {
		Key   string
		Value float64
	}
	dbmanager db.Manager
	ctx       context.Context
	cancel    context.CancelFunc
}

//CreatDiskCache creation
func CreatDiskCache(ctx context.Context) *DiskCache {
	cctx, cancel := context.WithCancel(ctx)
	return &DiskCache{
		dbmanager: db.GetManager(),
		ctx:       cctx,
		cancel:    cancel,
	}
}

//Start start statistics
func (d *DiskCache) Start() {
	d.setcache ()
	timer := time.NewTimer(time.Minute * 5)
	defer timer.Stop()
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-timer.C:
			d.setcache ()
			timer.Reset(time.Minute * 5)
		}
	}
}

//Stop stop
func (d *DiskCache) Stop() {
	logrus.Info("stop disk cache statistics")
	d.cancel()
}
func (d *DiskCache) setcache() {
	logrus.Info("start get all service disk size")
	start := time.Now()
	var diskcache []struct {
		Key   string
		Value float64
	}
	services, err := d.dbmanager.TenantServiceDao().GetAllServicesID()
	if err != nil {
		logrus.Errorln("Error get tenant service when select db :", err)
		return
	}
	_, err = d.dbmanager.TenantServiceVolumeDao().GetAllVolumes()
	if err != nil {
		logrus.Errorln("Error get tenant service volume when select db :", err)
		return
	}
	sharePath := os.Getenv("SHARE_DATA_PATH")
	if sharePath == "" {
		sharePath = "/grdata"
	}
	var cache = make(map[string]*model.TenantServices)
	for _, service := range services {
		//service nfs volume
		size := util.GetDirSize(fmt.Sprintf("%s/tenant/%s/service/%s", sharePath, service.TenantID, service.ServiceID))
		if size != 0 {
			diskcache = append(diskcache, struct {
				Key   string
				Value float64
			}{
				Key:   service.ServiceID + "_" + service.AppID + "_" + service.TenantID,
				Value: size,
			})
		}
		cache[service.ServiceID] = service
	}
	d.cache = diskcache
	logrus.Infof("end get all service disk size,time consum %2.f s", time.Since(start).Seconds())
}

//Get to obtain disk statistics
func (d *DiskCache) Get() map[string]float64 {
	newcache := make(map[string]float64)
	for _, v := range d.cache {
		newcache[v.Key] += v.Value
	}
	return newcache
}

// GetTenantDisk GetTenantDisk
func (d *DiskCache) GetTenantDisk(tenantID string) float64 {
	var value float64
	for _, v := range d.cache {
		if strings.HasSuffix(v.Key, "_"+tenantID) {
			value += v.Value
		}
	}
	return value
}

//GetServiceDisk GetServiceDisk
func (d *DiskCache) GetServiceDisk(serviceID string) float64 {
	var value float64
	for _, v := range d.cache {
		if strings.HasPrefix(v.Key, serviceID+"_") {
			value += v.Value
		}
	}
	return value
}
