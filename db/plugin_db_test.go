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

package db

import (
	"context"
	"fmt"
	dbconfig "github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/db/model"
	"github.com/testcontainers/testcontainers-go"
	"testing"
	"time"
)

func TestManager_PluginBuildVersionDaoImpl_ListSuccessfulOnesByPluginIDs(t *testing.T) {
	dbname := "region"
	rootpw := "kato"

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mariadb",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": rootpw,
			"MYSQL_DATABASE":      dbname,
		},
		Cmd: []string{"character-set-server=utf8mb4", "collation-server=utf8mb4_unicode_ci"},
	}
	mariadb, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer mariadb.Terminate(ctx)

	host, err := mariadb.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	port, err := mariadb.MappedPort(ctx, "3306")
	if err != nil {
		t.Error(err)
	}

	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root",
		rootpw, host, port.Int(), dbname)
	tryTimes := 3
	for {
		if err := CreateManager(dbconfig.Config{
			DBType:              "mysql",
			MysqlConnectionInfo: connInfo,
		}); err != nil {
			if tryTimes == 0 {
				t.Fatalf("Connect info: %s; error creating db manager: %v", connInfo, err)
			} else {
				tryTimes = tryTimes - 1
				time.Sleep(10 * time.Second)
				continue
			}
		}
		break
	}

	// prepare test data
	oridata := []struct {
		pluginID, status string
	}{
		{pluginID: "ff6aad8a70324384a7578285799e50d9", status: "complete"},
		{pluginID: "ff6aad8a70324384a7578285799e50d9", status: "failure"},
		{pluginID: "ff6aad8a70324384a7578285799e50d9", status: "failure"},
		{pluginID: "4998ca78d41f45149e71c1f03ad0aa22", status: "complete"},
	}
	for _, od := range oridata {
		buildVersion := &model.TenantPluginBuildVersion{
			PluginID:      od.pluginID,
			Status:        od.status,
			DeployVersion: time.Now().Format("20060102150405.000000"),
		}
		if err := GetManager().TenantPluginBuildVersionDao().AddModel(buildVersion); err != nil {
			t.Fatalf("failed to create plugin build version: %v", err)
		}
	}

	pluginIDs := []string{"ff6aad8a70324384a7578285799e50d9", "4998ca78d41f45149e71c1f03ad0aa22"}
	verions, err := GetManager().TenantPluginBuildVersionDao().ListSuccessfulOnesByPluginIDs(pluginIDs)
	if err != nil {
		t.Errorf("received unexpected error: %v", err)
	}
	for _, p := range verions {
		t.Logf("version id: %d; deploy version: %s; status: %s", p.ID, p.DeployVersion, p.Status)
	}
}
