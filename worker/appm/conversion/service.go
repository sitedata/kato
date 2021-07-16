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
	"fmt"
	"strings"

	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/jinzhu/gorm"
	yaml "gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//ServiceSource conv ServiceSource
func ServiceSource(as *v1.AppService, dbmanager db.Manager) error {
	sscs, err := dbmanager.ServiceSourceDao().GetServiceSource(as.ServiceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("conv service source failure %s", err.Error())
	}
	for _, ssc := range sscs {
		switch ssc.SourceType {
		case "deployment":
			var dm appsv1.Deployment
			if err := decoding(ssc.SourceBody, &dm); err != nil {
				return decodeError(err)
			}
			as.SetDeployment(&dm)
		case "statefulset":
			var ss appsv1.StatefulSet
			if err := decoding(ssc.SourceBody, &ss); err != nil {
				return decodeError(err)
			}
			as.SetStatefulSet(&ss)
		case "configmap":
			var cm corev1.ConfigMap
			if err := decoding(ssc.SourceBody, &cm); err != nil {
				return decodeError(err)
			}
			as.SetConfigMap(&cm)
		}
	}
	return nil
}
func decodeError(err error) error {
	return fmt.Errorf("decode service source failure %s", err.Error())
}
func decoding(source string, target interface{}) error {
	return yaml.Unmarshal([]byte(source), target)
}
func int32Ptr(i int) *int32 {
	j := int32(i)
	return &j
}

//TenantServiceBase conv tenant service base info
func TenantServiceBase(as *v1.AppService, dbmanager db.Manager) error {
	tenantService, err := dbmanager.TenantServiceDao().GetServiceByID(as.ServiceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrServiceNotFound
		}
		return fmt.Errorf("error getting service base info by serviceID(%s) %s", as.ServiceID, err.Error())
	}
	as.ServiceKind = dbmodel.ServiceKind(tenantService.Kind)
	tenant, err := dbmanager.TenantDao().GetTenantByUUID(tenantService.TenantID)
	if err != nil {
		return fmt.Errorf("get tenant info failure %s", err.Error())
	}
	as.TenantID = tenantService.TenantID
	if as.DeployVersion == "" {
		as.DeployVersion = tenantService.DeployVersion
	}
	as.ContainerCPU = tenantService.ContainerCPU
	as.AppID = tenantService.AppID
	as.ContainerMemory = tenantService.ContainerMemory
	as.Replicas = tenantService.Replicas
	as.ServiceAlias = tenantService.ServiceAlias
	as.UpgradeMethod = v1.TypeUpgradeMethod(tenantService.UpgradeMethod)
	if as.CreaterID == "" {
		as.CreaterID = string(util.NewTimeVersion())
	}
	as.TenantName = tenant.Name
	if err := initTenant(as, tenant); err != nil {
		return fmt.Errorf("conversion tenant info failure %s", err.Error())
	}
	if tenantService.Kind == dbmodel.ServiceKindThirdParty.String() {
		return nil
	}
	label, err := dbmanager.TenantServiceLabelDao().GetLabelByNodeSelectorKey(as.ServiceID, "windows")
	if label != nil {
		as.IsWindowsService = true
	}
	if !tenantService.IsState() {
		initBaseDeployment(as, tenantService)
		return nil
	}
	if tenantService.IsState() {
		initBaseStatefulSet(as, tenantService)
		return nil
	}
	return fmt.Errorf("Kind: %s; do not decision build type for service %s",
		tenantService.Kind, as.ServiceAlias)
}

func initTenant(as *v1.AppService, tenant *dbmodel.Tenants) error {
	if tenant == nil || tenant.UUID == "" {
		return fmt.Errorf("tenant is invalid")
	}
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   tenant.UUID,
			Labels: map[string]string{"creator": "Kato"},
		},
	}
	as.SetTenant(namespace)
	return nil
}
func initSelector(selector *metav1.LabelSelector, service *dbmodel.TenantServices) {
	if selector.MatchLabels == nil {
		selector.MatchLabels = make(map[string]string)
	}
	selector.MatchLabels["name"] = service.ServiceAlias
	selector.MatchLabels["tenant_id"] = service.TenantID
	selector.MatchLabels["service_id"] = service.ServiceID
	//selector.MatchLabels["version"] = service.DeployVersion
}
func initBaseStatefulSet(as *v1.AppService, service *dbmodel.TenantServices) {
	as.ServiceType = v1.TypeStatefulSet
	stateful := as.GetStatefulSet()
	if stateful == nil {
		stateful = &appsv1.StatefulSet{}
	}
	stateful.Namespace = as.TenantID
	stateful.Spec.Replicas = int32Ptr(service.Replicas)
	if stateful.Spec.Selector == nil {
		stateful.Spec.Selector = &metav1.LabelSelector{}
	}
	initSelector(stateful.Spec.Selector, service)
	stateful.Spec.ServiceName = service.ServiceName
	stateful.Name = service.ServiceName
	if stateful.Spec.ServiceName == "" {
		stateful.Spec.ServiceName = service.ServiceAlias
		stateful.Name = service.ServiceAlias
	}
	stateful.Namespace = service.TenantID
	stateful.GenerateName = service.ServiceAlias
	stateful.Labels = as.GetCommonLabels(stateful.Labels, map[string]string{
		"name":    service.ServiceAlias,
		"version": service.DeployVersion,
	})
	stateful.Spec.UpdateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType
	if as.UpgradeMethod == v1.OnDelete {
		stateful.Spec.UpdateStrategy.Type = appsv1.OnDeleteStatefulSetStrategyType
	}
	as.SetStatefulSet(stateful)
}

func initBaseDeployment(as *v1.AppService, service *dbmodel.TenantServices) {
	as.ServiceType = v1.TypeDeployment
	deployment := as.GetDeployment()
	if deployment == nil {
		deployment = &appsv1.Deployment{}
	}
	deployment.Namespace = as.TenantID
	deployment.Spec.Replicas = int32Ptr(service.Replicas)
	if deployment.Spec.Selector == nil {
		deployment.Spec.Selector = &metav1.LabelSelector{}
	}
	initSelector(deployment.Spec.Selector, service)
	deployment.Namespace = service.TenantID
	deployment.Name = service.ServiceID + "-deployment"
	deployment.GenerateName = strings.Replace(service.ServiceAlias, "_", "-", -1)
	deployment.Labels = as.GetCommonLabels(deployment.Labels, map[string]string{
		"name":    service.ServiceAlias,
		"version": service.DeployVersion,
	})
	deployment.Spec.Strategy.Type = appsv1.RollingUpdateDeploymentStrategyType
	if as.UpgradeMethod == v1.OnDelete {
		deployment.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType
	}
	as.SetDeployment(deployment)
}
