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

	"github.com/sirupsen/logrus"

	autoscalingv2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
)

var str2ResourceName = map[string]corev1.ResourceName{
	"cpu":    corev1.ResourceCPU,
	"memory": corev1.ResourceMemory,
}

// TenantServiceAutoscaler -
func TenantServiceAutoscaler(as *v1.AppService, dbmanager db.Manager) error {
	hpas, err := newHPAs(as, dbmanager)
	if err != nil {
		return fmt.Errorf("create HPAs: %v", err)
	}
	logrus.Debugf("the numbers of HPAs: %d", len(hpas))

	as.SetHPAs(hpas)

	return nil
}

func newHPAs(as *v1.AppService, dbmanager db.Manager) ([]*autoscalingv2.HorizontalPodAutoscaler, error) {
	xpaRules, err := dbmanager.TenantServceAutoscalerRulesDao().ListEnableOnesByServiceID(as.ServiceID)
	if err != nil {
		return nil, err
	}

	var hpas []*autoscalingv2.HorizontalPodAutoscaler
	for _, rule := range xpaRules {
		metrics, err := dbmanager.TenantServceAutoscalerRuleMetricsDao().ListByRuleID(rule.RuleID)
		if err != nil {
			return nil, err
		}

		var kind, name string
		if as.GetStatefulSet() != nil {
			kind, name = "StatefulSet", as.GetStatefulSet().GetName()
		} else {
			kind, name = "Deployment", as.GetDeployment().GetName()
		}

		labels := as.GetCommonLabels(map[string]string{
			"rule_id": rule.RuleID,
			"version": as.DeployVersion,
		})

		hpa := newHPA(as.TenantID, kind, name, labels, rule, metrics)

		hpas = append(hpas, hpa)
	}

	return hpas, nil
}

func createResourceMetrics(metric *model.TenantServiceAutoscalerRuleMetrics) autoscalingv2.MetricSpec {
	ms := autoscalingv2.MetricSpec{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: str2ResourceName[metric.MetricsName],
		},
	}

	if metric.MetricTargetType == "utilization" {
		value := int32(metric.MetricTargetValue)
		ms.Resource.Target = autoscalingv2.MetricTarget{
			Type:               autoscalingv2.UtilizationMetricType,
			AverageUtilization: &value,
		}
	}
	if metric.MetricTargetType == "average_value" {
		ms.Resource.Target.Type = autoscalingv2.AverageValueMetricType
		if metric.MetricsName == "cpu" {
			ms.Resource.Target.AverageValue = resource.NewMilliQuantity(int64(metric.MetricTargetValue), resource.DecimalSI)
		}
		if metric.MetricsName == "memory" {
			ms.Resource.Target.AverageValue = resource.NewQuantity(int64(metric.MetricTargetValue*1024*1024), resource.BinarySI)
		}
	}

	return ms
}

func newHPA(namespace, kind, name string, labels map[string]string, rule *model.TenantServiceAutoscalerRules, metrics []*model.TenantServiceAutoscalerRuleMetrics) *autoscalingv2.HorizontalPodAutoscaler {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rule.RuleID,
			Namespace: namespace,
			Labels:    labels,
		},
	}

	spec := autoscalingv2.HorizontalPodAutoscalerSpec{
		MinReplicas: util.Int32(int32(rule.MinReplicas)),
		MaxReplicas: int32(rule.MaxReplicas),
		ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
			Kind:       kind,
			Name:       name,
			APIVersion: "apps/v1",
		},
	}

	for _, metric := range metrics {
		if metric.MetricsType != "resource_metrics" {
			logrus.Warningf("rule id:  %s; unsupported metric type: %s", rule.RuleID, metric.MetricsType)
			continue
		}
		if metric.MetricTargetValue <= 0 {
			// TODO: If the target value of cpu and memory is 0, it will not take effect.
			// TODO: The target value of the custom indicator can be 0.
			continue
		}

		ms := createResourceMetrics(metric)
		spec.Metrics = append(spec.Metrics, ms)
	}
	if len(spec.Metrics) == 0 {
		return nil
	}
	hpa.Spec = spec

	return hpa
}
