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
	"testing"
	"time"

	dbconfig "github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	"github.com/testcontainers/testcontainers-go"
)

func TestTenantServicesDao_GetOpenedPort(t *testing.T) {
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

	sid := util.NewUUID()
	trueVal := true
	falseVal := true
	err = GetManager().TenantServicesPortDao().AddModel(&model.TenantServicesPort{
		ServiceID:      sid,
		ContainerPort:  1111,
		MappingPort:    1111,
		IsInnerService: &falseVal,
		IsOuterService: &trueVal,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = GetManager().TenantServicesPortDao().AddModel(&model.TenantServicesPort{
		ServiceID:      sid,
		ContainerPort:  2222,
		MappingPort:    2222,
		IsInnerService: &trueVal,
		IsOuterService: &falseVal,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = GetManager().TenantServicesPortDao().AddModel(&model.TenantServicesPort{
		ServiceID:      sid,
		ContainerPort:  3333,
		MappingPort:    3333,
		IsInnerService: &falseVal,
		IsOuterService: &falseVal,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = GetManager().TenantServicesPortDao().AddModel(&model.TenantServicesPort{
		ServiceID:      sid,
		ContainerPort:  5555,
		MappingPort:    5555,
		IsInnerService: &trueVal,
		IsOuterService: &trueVal,
	})
	if err != nil {
		t.Fatal(err)
	}
	ports, err := GetManager().TenantServicesPortDao().GetOpenedPorts(sid)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 3 {
		t.Errorf("Expected 3 for the length of ports, but return %d", len(ports))
	}
}

func TestListInnerPorts(t *testing.T) {
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

	sid := util.NewUUID()
	trueVal := true
	falseVal := false
	err = GetManager().TenantServicesPortDao().AddModel(&model.TenantServicesPort{
		ServiceID:      sid,
		ContainerPort:  1111,
		MappingPort:    1111,
		IsInnerService: &trueVal,
		IsOuterService: &trueVal,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = GetManager().TenantServicesPortDao().AddModel(&model.TenantServicesPort{
		ServiceID:      sid,
		ContainerPort:  2222,
		MappingPort:    2222,
		IsInnerService: &trueVal,
		IsOuterService: &falseVal,
	})
	if err != nil {
		t.Fatal(err)
	}

	ports, err := GetManager().TenantServicesPortDao().ListInnerPortsByServiceIDs([]string{sid})
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 2 {
		t.Errorf("Expocted %d for ports, but got %d", 2, len(ports))
	}
}
