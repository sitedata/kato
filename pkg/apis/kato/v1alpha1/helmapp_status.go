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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHelmAppCondition creates a new HelmApp condition.
func NewHelmAppCondition(condType HelmAppConditionType, status corev1.ConditionStatus, reason, message string) *HelmAppCondition {
	return &HelmAppCondition{
		Type:               condType,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

// GetCondition returns a HelmApp condition based on the given type.
func (in *HelmAppStatus) GetCondition(t HelmAppConditionType) (int, *HelmAppCondition) {
	for i, c := range in.Conditions {
		if t == c.Type {
			return i, &c
		}
	}
	return -1, nil
}

// IsConditionTrue checks if the condition is ready or not based on the given condition type.
func (in *HelmAppStatus) IsConditionTrue(t HelmAppConditionType) bool {
	idx, condition := in.GetCondition(t)
	return idx != -1 && condition.Status == corev1.ConditionTrue
}

// SetCondition setups the given HelmApp condition.
func (in *HelmAppStatus) SetCondition(c HelmAppCondition) {
	pos, cp := in.GetCondition(c.Type)
	if cp != nil &&
		cp.Status == c.Status && cp.Reason == c.Reason && cp.Message == c.Message {
		return
	}

	if cp != nil {
		in.Conditions[pos] = c
	} else {
		in.Conditions = append(in.Conditions, c)
	}
}

// UpdateCondition updates existing HelmApp condition or creates a new
// one. Sets LastTransitionTime to now if the status has changed.
// Returns true if HelmApp condition has changed or has been added.
func (in *HelmAppStatus) UpdateCondition(condition *HelmAppCondition) bool {
	condition.LastTransitionTime = metav1.Now()
	// Try to find this HelmApp condition.
	conditionIndex, oldCondition := in.GetCondition(condition.Type)

	if oldCondition == nil {
		// We are adding new HelmApp condition.
		in.Conditions = append(in.Conditions, *condition)
		return true
	}

	// We are updating an existing condition, so we need to check if it has changed.
	if condition.Status == oldCondition.Status {
		condition.LastTransitionTime = oldCondition.LastTransitionTime
	}

	isEqual := condition.Status == oldCondition.Status &&
		condition.Reason == oldCondition.Reason &&
		condition.Message == oldCondition.Message &&
		condition.LastTransitionTime.Equal(&oldCondition.LastTransitionTime)

	in.Conditions[conditionIndex] = *condition
	// Return true if one of the fields have changed.
	return !isEqual
}

func (in *HelmAppStatus) UpdateConditionStatus(conditionType HelmAppConditionType, conditionStatus corev1.ConditionStatus) {
	_, condition := in.GetCondition(conditionType)
	if condition != nil {
		condition.Status = conditionStatus
		if conditionStatus == corev1.ConditionTrue {
			condition.Reason = ""
			condition.Message = ""
		}
		in.UpdateCondition(condition)
		return
	}

	condition = NewHelmAppCondition(conditionType, conditionStatus, "", "")
	in.UpdateCondition(condition)
}
