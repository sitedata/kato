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

package conversion

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/db/dao"
	"github.com/gridworkz/kato/db/model"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/gridworkz/kato/worker/appm/volume"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTenantServiceVersion(t *testing.T) {
	var as v1.AppService
	TenantServiceVersion(&as, nil)
}

func TestConvertRulesToEnvs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbmanager := db.NewMockManager(ctrl)

	as := &v1.AppService{}
	as.ServiceID = "dummy service id"
	as.TenantName = "dummy tenant name"
	as.ServiceAlias = "dummy service alias"

	httpRuleDao := dao.NewMockHTTPRuleDao(ctrl)
	httpRuleDao.EXPECT().GetHTTPRuleByServiceIDAndContainerPort(as.ServiceID, 0).Return(nil, nil)
	dbmanager.EXPECT().HTTPRuleDao().Return(httpRuleDao)

	port := &model.TenantServicesPort{
		TenantID:       "dummy tenant id",
		ServiceID:      as.ServiceID,
		ContainerPort:  0,
		Protocol:       "http",
		PortAlias:      "GRD835895000",
		IsInnerService: func() *bool { b := false; return &b }(),
		IsOuterService: func() *bool { b := true; return &b }(),
	}

	renvs := convertRulesToEnvs(as, dbmanager, []*model.TenantServicesPort{port})
	if len(renvs) > 0 {
		t.Errorf("Expected 0 for the length rule envs, but return %d", len(renvs))
	}
}

func TestCreateVolume(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db.CreateManager(config.Config{DBType: "mysql", MysqlConnectionInfo: "oc6Poh:noot6Mea@tcp(192.168.2.203:3306)/region"})
	dbmanager := db.GetManager()

	as := &v1.AppService{}
	as.ServiceID = "dummy service id"
	as.TenantName = "dummy tenant name"
	as.ServiceAlias = "dummy service alias"
	var replicas int32
	as.SetStatefulSet(&appv1.StatefulSet{Spec: appv1.StatefulSetSpec{Replicas: &replicas, Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"version": "version"}}}}})

	serviceVolume, err := db.GetManager().TenantServiceVolumeDao().GetVolumeByID(25)
	if err != nil {
		t.Log(err)
		return
	}
	version := &model.VersionInfo{}

	vol := volume.NewVolumeManager(as, serviceVolume, nil, version, nil, nil, dbmanager)
	var define = &volume.Define{}
	vol.CreateVolume(define)
}
func TestFoobar(t *testing.T) {
	memory := 64
	cpuRequest, cpuLimit := int64(memory)/128*30, int64(memory)/128*80
	t.Errorf("request: %d; limit: %d", cpuRequest, cpuLimit)
}
