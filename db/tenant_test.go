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
	"github.com/gridworkz/kato/util"
	"github.com/testcontainers/testcontainers-go"
	"testing"
	"time"
)

func TestTenantServicesDao_ListThirdPartyServices(t *testing.T) {
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

	svcs, err := GetManager().TenantServiceDao().ListThirdPartyServices()
	if err != nil {
		t.Fatalf("error listing third-party service: %v", err)
	}
	if len(svcs) != 0 {
		t.Errorf("Expected 0 for the length of third-party services, but returned %d", len(svcs))
	}

	for i := 0; i < 3; i++ {
		item1 := &model.TenantServices{
			TenantID:  util.NewUUID(),
			ServiceID: util.NewUUID(),
			Kind:      model.ServiceKindThirdParty.String(),
		}
		if err = GetManager().TenantServiceDao().AddModel(item1); err != nil {
			t.Fatalf("error create third-party service: %v", err)
		}
	}
	svcs, err = GetManager().TenantServiceDao().ListThirdPartyServices()
	if err != nil {
		t.Fatalf("error listing third-party service: %v", err)
	}
	if len(svcs) != 3 {
		t.Errorf("Expected 3 for the length of third-party services, but returned %d", len(svcs))
	}
}

func TestTenantServicesPortDao_HasOpenPort(t *testing.T) {
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

	t.Run("service doesn't exist", func(t *testing.T) {
		hasOpenPort := GetManager().TenantServicesPortDao().HasOpenPort("foobar")
		if hasOpenPort {
			t.Error("Expected false for hasOpenPort, but returned true")
		}
	})
	trueVal := true
	falseVal := true
	t.Run("outer service", func(t *testing.T) {
		port := &model.TenantServicesPort{
			ServiceID:      util.NewUUID(),
			IsOuterService: &trueVal,
		}
		if err := GetManager().TenantServicesPortDao().AddModel(port); err != nil {
			t.Fatalf("error creating TenantServicesPort: %v", err)
		}
		hasOpenPort := GetManager().TenantServicesPortDao().HasOpenPort(port.ServiceID)
		if !hasOpenPort {
			t.Errorf("Expected true for hasOpenPort, but returned %v", hasOpenPort)
		}
	})
	t.Run("inner service", func(t *testing.T) {
		port := &model.TenantServicesPort{
			ServiceID:      util.NewUUID(),
			IsInnerService: &trueVal,
		}
		if err := GetManager().TenantServicesPortDao().AddModel(port); err != nil {
			t.Fatalf("error creating TenantServicesPort: %v", err)
		}
		hasOpenPort := GetManager().TenantServicesPortDao().HasOpenPort(port.ServiceID)
		if !hasOpenPort {
			t.Errorf("Expected true for hasOpenPort, but returned %v", hasOpenPort)
		}
	})
	t.Run("not inner or outer service", func(t *testing.T) {
		port := &model.TenantServicesPort{
			ServiceID:      util.NewUUID(),
			IsInnerService: &falseVal,
			IsOuterService: &falseVal,
		}
		if err := GetManager().TenantServicesPortDao().AddModel(port); err != nil {
			t.Fatalf("error creating TenantServicesPort: %v", err)
		}
		hasOpenPort := GetManager().TenantServicesPortDao().HasOpenPort(port.ServiceID)
		if hasOpenPort {
			t.Errorf("Expected false for hasOpenPort, but returned %v", hasOpenPort)
		}
	})
}
