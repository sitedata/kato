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

package clean

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/util"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder/sources"
)

//Manager CleanManager
type Manager struct {
	dclient *client.Client
	ctx     context.Context
	cancel  context.CancelFunc
}

//CreateCleanManager
func CreateCleanManager() (*Manager, error) {
	dclient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &Manager{
		dclient: dclient,
		ctx:     ctx,
		cancel:  cancel,
	}
	return c, nil
}

//Start clean
func (t *Manager) Start(errchan chan error) error {
	logrus.Info("CleanManager is starting.")
	run := func() {
		err := util.Exec(t.ctx, func() error {
			now := time.Now()
			datetime := now.AddDate(0, -1, 0)
			// Find more than five versions
			results, err := db.GetManager().VersionInfoDao().SearchVersionInfo()
			if err != nil {
				logrus.Error(err)
			}
			var serviceIDList []string
			for _, v := range results {
				serviceIDList = append(serviceIDList, v.ServiceID)
			}
			versions, err := db.GetManager().VersionInfoDao().GetVersionInfo(datetime, serviceIDList)
			if err != nil {
				logrus.Error(err)
			}

			for _, v := range versions {
				versions, err := db.GetManager().VersionInfoDao().GetVersionByServiceID(v.ServiceID)
				if err != nil {
					logrus.Error("GetVersionByServiceID error: ", err.Error())
					continue
				}
				if len(versions) <= 5 {
					continue
				}
				if v.DeliveredType == "image" {
					imagePath := v.DeliveredPath
					//remove local image, However, it is important to note that the version image is stored in the image repository
					err := sources.ImageRemove(t.dclient, imagePath)
					if err != nil {
						logrus.Error(err)
					}
					if err := db.GetManager().VersionInfoDao().DeleteVersionInfo(v); err != nil {
						logrus.Error(err)
						continue
					}
					logrus.Info("Image deletion successful:", imagePath)
				}
				if v.DeliveredType == "slug" {
					filePath := v.DeliveredPath
					if err := os.Remove(filePath); err != nil {
						logrus.Error(err)
					}
					if err := db.GetManager().VersionInfoDao().DeleteVersionInfo(v); err != nil {
						logrus.Error(err)
						continue
					}
					logrus.Info("file deletion successful:", filePath)

				}

			}
			return nil
		}, 24*time.Hour)
		if err != nil {
			errchan <- err
		}
	}
	go run()
	return nil
}

//Stop
func (t *Manager) Stop() error {
	logrus.Info("CleanManager is stoping.")
	t.cancel()
	return nil
}
