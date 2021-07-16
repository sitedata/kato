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

package exector

import (
	"testing"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
)

func TestModify(t *testing.T) {
	var b = BackupAPPRestore{
		BackupID:      "test",
		EventID:       "test",
		serviceChange: make(map[string]*Info, 0),
		Logger:        event.GetTestLogger(),
	}
	var appSnapshots = []*RegionServiceSnapshot{
		&RegionServiceSnapshot{
			ServiceID: "1234",
			Service: &dbmodel.TenantServices{
				ServiceID: "1234",
			},
			ServiceMntRelation: []*dbmodel.TenantServiceMountRelation{
				&dbmodel.TenantServiceMountRelation{
					ServiceID:       "1234",
					DependServiceID: "123456",
				},
			},
		},
		&RegionServiceSnapshot{
			ServiceID: "123456",
			Service: &dbmodel.TenantServices{
				ServiceID: "1234",
			},
			ServiceEnv: []*dbmodel.TenantServiceEnvVar{
				&dbmodel.TenantServiceEnvVar{
					ServiceID: "123456",
					Name:      "testenv",
				},
				&dbmodel.TenantServiceEnvVar{
					ServiceID: "123456",
					Name:      "testenv2",
				},
			},
		},
	}
	appSnapshot := AppSnapshot{
		Services: appSnapshots,
	}
	b.modify(&appSnapshot)
	re, _ := ffjson.Marshal(appSnapshot)
	t.Log(string(re))
}

func TestUnzipAllDataFile(t *testing.T) {
	allDataFilePath := "/tmp/__all_data.zip"
	allTmpDir := "/tmp/4f25c53e864744ec95d037528acaa708"
	if err := util.Unzip(allDataFilePath, allTmpDir); err != nil {
		logrus.Errorf("unzip all data file failure %s", err.Error())
	}
}
