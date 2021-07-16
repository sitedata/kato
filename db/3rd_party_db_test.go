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

func TestEndpointDaoImpl_UpdateModel(t *testing.T) {
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

	trueVal := true
	falseVal := false
	ep := &model.Endpoint{
		UUID:      util.NewUUID(),
		ServiceID: util.NewUUID(),
		IP:        "10.10.10.10",
		IsOnline:  &trueVal,
	}
	err = GetManager().EndpointsDao().AddModel(ep)
	if err != nil {
		t.Fatalf("error adding endpoint: %v", err)
	}
	ep.IsOnline = &falseVal
	err = GetManager().EndpointsDao().UpdateModel(ep)
	if err != nil {
		t.Fatalf("error updating endpoint: %v", err)
	}
	e, err := GetManager().EndpointsDao().GetByUUID(ep.UUID)
	if err != nil {
		t.Fatalf("error getting endpoint: %v", err)
	}
	if *e.IsOnline != false {
		t.Errorf("Expected %v for e.IsOnline, but returned %v", false, e.IsOnline)
	}
}

func TestEndpointDaoImpl_AddModel(t *testing.T) {
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
		t.Fatal(err)
	}
	port, err := mariadb.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatal(err)
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

	falseVal := false
	ep := &model.Endpoint{
		UUID:      util.NewUUID(),
		ServiceID: util.NewUUID(),
		IP:        "10.10.10.10",
		IsOnline:  &falseVal,
	}
	err = GetManager().EndpointsDao().AddModel(ep)
	if err != nil {
		t.Fatalf("error adding endpoint: %v", err)
	}
	e, err := GetManager().EndpointsDao().GetByUUID(ep.UUID)
	if err != nil {
		t.Fatalf("error getting endpoint: %v", err)
	}
	if *e.IsOnline != false {
		t.Errorf("Expected %v for e.IsOnline, but returned %v", false, e.IsOnline)
	}
}
