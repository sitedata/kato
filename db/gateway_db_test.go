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

func TestGwRuleConfig(t *testing.T) {
	dbm, err := CreateTestManager()
	if err != nil {
		t.Fatalf("error creating test db manager: %v", err)
	}
	rid := util.NewUUID()
	cfg := &model.GwRuleConfig{
		RuleID: rid,
		Key:    "set-header-Host",
		Value:  "$http_host",
	}
	if err := dbm.GwRuleConfigDao().AddModel(cfg); err != nil {
		t.Fatalf("error create rule config: %v", err)
	}
	cfg2 := &model.GwRuleConfig{
		RuleID: rid,
		Key:    "set-header-foo",
		Value:  "bar",
	}
	if err := dbm.GwRuleConfigDao().AddModel(cfg2); err != nil {
		t.Fatalf("error create rule config: %v", err)
	}
	list, err := dbm.GwRuleConfigDao().ListByRuleID(rid)
	if err != nil {
		t.Fatalf("error listing configs: %v", err)
	}
	if list == nil && len(list) != 2 {
		t.Errorf("Expected 2 for the length fo list, but returned %d", len(list))
	}

	if err := dbm.GwRuleConfigDao().DeleteByRuleID(rid); err != nil {
		t.Fatalf("error deleting rule config: %v", err)
	}
	list, err = dbm.GwRuleConfigDao().ListByRuleID(rid)
	if err != nil {
		t.Fatalf("error listing configs: %v", err)
	}
	if list != nil && len(list) > 0 {
		t.Errorf("Expected empty for list, but returned %+v", list)
	}
}

func TestCertificateDaoImpl_AddOrUpdate(t *testing.T) {
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

	cert := &model.Certificate{
		UUID:            util.NewUUID(),
		CertificateName: "dummy-name",
		Certificate:     "dummy-certificate",
		PrivateKey:      "dummy-privateKey",
	}
	err = GetManager().CertificateDao().AddOrUpdate(cert)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := GetManager().CertificateDao().GetCertificateByID(cert.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.UUID != cert.UUID {
		t.Errorf("Expected %s for resp.UUID, but returned %s", cert.UUID, resp.UUID)
	}
	if resp.CertificateName != cert.CertificateName {
		t.Errorf("Expected %s for resp.CertificateName, but returned %s", cert.CertificateName, resp.CertificateName)
	}
	if resp.Certificate != cert.Certificate {
		t.Errorf("Expected %s for resp.Certificate, but returned %s", cert.Certificate, resp.Certificate)
	}
	if resp.PrivateKey != cert.PrivateKey {
		t.Errorf("Expected %s for resp.UUID, but returned %s", cert.PrivateKey, resp.PrivateKey)
	}

	cert.Certificate = "update-certificate"
	cert.PrivateKey = "update-privateKey"
	err = GetManager().CertificateDao().AddOrUpdate(cert)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = GetManager().CertificateDao().GetCertificateByID(cert.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if resp.UUID != cert.UUID {
		t.Errorf("Expected %s for resp.UUID, but returned %s", cert.UUID, resp.UUID)
	}
	if resp.CertificateName != cert.CertificateName {
		t.Errorf("Expected %s for resp.CertificateName, but returned %s", cert.CertificateName, resp.CertificateName)
	}
	if resp.Certificate != cert.Certificate {
		t.Errorf("Expected %s for resp.Certificate, but returned %s", cert.Certificate, resp.Certificate)
	}
	if resp.PrivateKey != cert.PrivateKey {
		t.Errorf("Expected %s for resp.UUID, but returned %s", cert.PrivateKey, resp.PrivateKey)
	}
}
