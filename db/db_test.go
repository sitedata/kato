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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	dbconfig "github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
)

func CreateTestManager() (Manager, error) {
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
		return nil, err
	}
	defer mariadb.Terminate(ctx)

	host, err := mariadb.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := mariadb.MappedPort(ctx, "3306")
	if err != nil {
		return nil, err
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
				return nil, fmt.Errorf("Connect info: %s; error creating db manager: %v", connInfo, err)
			}
			tryTimes = tryTimes - 1
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	return GetManager(), nil
}

func TestTenantDao(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "root:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	err := GetManager().TenantDao().AddModel(&model.Tenants{
		Name: "gridworkz4",
		UUID: util.NewUUID(),
	})
	if err != nil {
		t.Fatal(err)
	}
	tenant, err := GetManager().TenantDao().GetTenantByUUID("27bbdd119b24444696dc51aa2f41eef8")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tenant)
}

func TestTenantServiceDao(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "root:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	service, err := GetManager().TenantServiceDao().GetServiceByTenantIDAndServiceAlias("27bbdd119b24444696dc51aa2f41eef8", "grb58f90")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
	service, err = GetManager().TenantServiceDao().GetServiceByID("2f29882148c19f5f84e3a7cedf6097c7")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
	services, err := GetManager().TenantServiceDao().GetServiceAliasByIDs([]string{"2f29882148c19f5f84e3a7cedf6097c7"})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range services {
		t.Log(s)
	}
}

func TestGetServiceEnvs(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "root:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	envs, err := GetManager().TenantServiceEnvVarDao().GetServiceEnvs("2f29882148c19f5f84e3a7cedf6097c7", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range envs {
		t.Log(e)
	}
}

func TestSetServiceLabel(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "root:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	label := model.TenantServiceLable{
		LabelKey:   "labelkey",
		LabelValue: "labelvalue",
		ServiceID:  "889bb1f028f655bebd545f24aa184a0b",
	}
	label.CreatedAt = time.Now()
	label.ID = 1
	err := GetManager().TenantServiceLabelDao().UpdateModel(&label)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTenantServiceLBMappingPort(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "root:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	mapPort, err := GetManager().TenantServiceLBMappingPortDao().CreateTenantServiceLBMappingPort("889bb1f028f655bebd545f24aa184a0b", 8080)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mapPort)
}

func TestCreateTenantServiceLBMappingPortTran(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "root:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	tx := GetManager().Begin()
	mapPort, err := GetManager().TenantServiceLBMappingPortDaoTransactions(tx).CreateTenantServiceLBMappingPort("889bb1f028f655bebd545f24aa184a0b", 8082)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
		return
	}
	tx.Commit()
	t.Log(mapPort)
}

func TestGetMem(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "admin:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	// err := GetManager().TenantDao().AddModel(&model.Tenants{
	// 	Name: "gridworkz3",
	// 	UUID: util.NewUUID(),
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func TestYugabyteDBCreateTable(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "postgresql://root@localhost:5432/region?sslmode=disable",
		DBType:              "yugabytedb",
	}); err != nil {
		t.Fatal(err)
	}
}
func TestYugabyteDBCreateService(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "postgresql://root@localhost:5432/region?sslmode=disable",
		DBType:              "yugabytedb",
	}); err != nil {
		t.Fatal(err)
	}
	err := GetManager().TenantServiceDao().AddModel(&model.TenantServices{
		TenantID:     "asdasd",
		ServiceID:    "asdasdasdasd",
		ServiceAlias: "grasdasdasdads",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestYugabyteDBDeleteService(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "postgresql://root@localhost:5432/region?sslmode=disable",
		DBType:              "yugabytedb",
	}); err != nil {
		t.Fatal(err)
	}
	err := GetManager().TenantServiceDao().DeleteServiceByServiceID("asdasdasdasd")
	if err != nil {
		t.Fatal(err)
	}
}

//func TestYugabyteDBSaveDeployInfo(t *testing.T) {
//	if err := CreateManager(dbconfig.Config{
//		MysqlConnectionInfo: "postgresql://root@localhost:5432/region?sslmode=disable",
//		DBType:              "yugabytedb",
//	}); err != nil {
//		t.Fatal(err)
//	}
//	err := GetManager().K8sDeployReplicationDao().AddModel(&model.K8sDeployReplication{
//		TenantID:        "asdasd",
//		ServiceID:       "asdasdasdasd",
//		ReplicationID:   "asdasdadsasdasdasd",
//		ReplicationType: model.TypeReplicationController,
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//}

//func TestYugabyteDBDeleteDeployInfo(t *testing.T) {
//	if err := CreateManager(dbconfig.Config{
//		MysqlConnectionInfo: "postgresql://root@localhost:5432/region?sslmode=disable",
//		DBType:              "yugabytedb",
//	}); err != nil {
//		t.Fatal(err)
//	}
//err := GetManager().K8sDeployReplicationDao().DeleteK8sDeployReplication("asdasdadsasdasdasd")
//if err != nil {
//	t.Fatal(err)
//}
//}

func TestGetCertificateByID(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "admin:admin@tcp(127.0.0.1:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}

	cert, err := GetManager().TenantServiceLabelDao().GetTenantNodeAffinityLabel("105bb7d4b94774f922edb3051bdf8ce1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cert)
}

var shareFileVolumeType = &model.TenantServiceVolumeType{
	VolumeType:  "share-file",
	NameShow:    "Shared storage (file)",
	Description: "Distributed file storage, shared and mounted within tenants, suitable for all types of applications",
	// CapacityValidation:   string(bs),
	ReclaimPolicy: "Retain",
	BackupPolicy:  "exclusive",
	SharePolicy:   "exclusive",
	AccessMode:    "RWO,ROX,RWX",
}

var localVolumeType = &model.TenantServiceVolumeType{
	VolumeType:  "local",
	NameShow:    "Local storage",
	Description: "Local storage device, suitable for stateful database services",
	// CapacityValidation:   string(bs),
	ReclaimPolicy: "Retain",
	BackupPolicy:  "exclusive",
	SharePolicy:   "exclusive",
	AccessMode:    "RWO,ROX,RWX",
}

var memoryFSVolumeType = &model.TenantServiceVolumeType{
	VolumeType:  "memoryfs",
	NameShow:    "Memory file storage",
	Description: "Memory-based storage devices whose capacity is limited by the amount of memory. Data is lost when the application restarts, suitable for high-speed temporary storage of data",
	// CapacityValidation:   string(bs),
	ReclaimPolicy: "Retain",
	BackupPolicy:  "exclusive",
	SharePolicy:   "exclusive",
	AccessMode:    "RWO,ROX,RWX",
}

var alicloudDiskAvailableVolumeType = &model.TenantServiceVolumeType{
	VolumeType:  "alicloud-disk-available",
	NameShow:    "Alibaba Cloud Disk (Smart Choice)",
	Description: "Alibaba Cloud chooses cloud disks intelligently. It will try to create cloud disk types supported by the current Alibaba Cloud region in the order of efficient cloud disk, SSD, and basic cloud disk",
	// CapacityValidation:   string(bs),
	ReclaimPolicy: "Delete",
	BackupPolicy:  "exclusive",
	SharePolicy:   "exclusive",
	AccessMode:    "RWO",
	Sort:          10,
}

var alicloudDiskcommonVolumeType = &model.TenantServiceVolumeType{
	VolumeType:  "alicloud-disk-common",
	NameShow:    "Alibaba Cloud Disk (Basic)",
	Description: "Alibaba Cloud ordinary basic cloud disk. Minimum 5G",
	// CapacityValidation:   string(bs),
	ReclaimPolicy: "Delete",
	BackupPolicy:  "exclusive",
	SharePolicy:   "exclusive",
	AccessMode:    "RWO",
	Sort:          13,
}
var alicloudDiskEfficiencyVolumeType = &model.TenantServiceVolumeType{
	VolumeType:  "alicloud-disk-efficiency",
	NameShow:    "Alibaba Cloud Disk (efficient)",
	Description: "Alibaba Cloud Efficient Cloud Disk. Minimum limit 20G",
	// CapacityValidation:   string(bs),
	ReclaimPolicy: "Delete",
	BackupPolicy:  "exclusive",
	SharePolicy:   "exclusive",
	AccessMode:    "RWO",
	Sort:          11,
}
var alicloudDiskeSSDVolumeType = &model.TenantServiceVolumeType{
	VolumeType:  "alicloud-disk-ssd",
	NameShow:    "Alibaba Cloud Disk (SSD)",
	Description: "Alibaba Cloud SSD type cloud disk. Minimum limit 20G",
	// CapacityValidation:   string(bs),
	ReclaimPolicy: "Delete",
	BackupPolicy:  "exclusive",
	SharePolicy:   "exclusive",
	AccessMode:    "RWO",
	Sort:          12,
}

func TestVolumeType(t *testing.T) {

	capacityValidation := make(map[string]interface{})
	capacityValidation["required"] = true
	capacityValidation["default"] = 20
	capacityValidation["min"] = 20
	capacityValidation["max"] = 32768 // [ali-cloud-disk usage limit](https://help.aliyun.com/document_detail/25412.html?spm=5176.2020520101.0.0.41d84df5faliP4)
	bs, _ := json.Marshal(capacityValidation)
	shareFileVolumeType.CapacityValidation = string(bs)
	localVolumeType.CapacityValidation = string(bs)
	memoryFSVolumeType.CapacityValidation = string(bs)
	alicloudDiskeSSDVolumeType.CapacityValidation = string(bs)
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "ieZoo9:Maigoed0@tcp(192.168.2.108:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
	if err := GetManager().VolumeTypeDao().AddModel(alicloudDiskeSSDVolumeType); err != nil {
		t.Error(err)
	} else {
		t.Log("yes")
	}
}

func initDBManager(t *testing.T) {
	if err := CreateManager(dbconfig.Config{
		MysqlConnectionInfo: "ieZoo9:Maigoed0@tcp(192.168.2.108:3306)/region",
		DBType:              "mysql",
	}); err != nil {
		t.Fatal(err)
	}
}
func TestGetVolumeType(t *testing.T) {
	initDBManager(t)
	vts, err := GetManager().VolumeTypeDao().GetAllVolumeTypes()
	if err != nil {
		t.Fatal(err)
	}
	for _, vt := range vts {
		t.Logf("%+v", vt)
	}
	t.Logf("volume type len is : %v", len(vts))
}

func TestGetVolumeTypeByType(t *testing.T) {
	initDBManager(t)
	vt, err := GetManager().VolumeTypeDao().GetVolumeTypeByType("ceph-rbd")
	if err != nil {
		t.Fatal("get volumeType by type error: ", err.Error())
	}
	t.Logf("%+v", vt)
}
