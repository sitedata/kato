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

package helmapp

import (
	"context"

	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	"github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// Status represents the status of helm app.
type Status struct {
	ctx            context.Context
	katoClient versioned.Interface
	helmApp        *v1alpha1.HelmApp
}

// NewStatus creates a new helm app status.
func NewStatus(ctx context.Context, app *v1alpha1.HelmApp, katoClient versioned.Interface) *Status {
	return &Status{
		ctx:            ctx,
		helmApp:        app,
		katoClient: katoClient,
	}
}

// Update updates helm app status.
func (s *Status) Update() error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		ctx, cancel := context.WithTimeout(s.ctx, defaultTimeout)
		defer cancel()

		helmApp, err := s.katoClient.KatoV1alpha1().HelmApps(s.helmApp.Namespace).Get(ctx, s.helmApp.Name, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "get helm app before update")
		}

		s.helmApp.Status.Phase = s.getPhase()
		s.helmApp.ResourceVersion = helmApp.ResourceVersion
		_, err = s.katoClient.KatoV1alpha1().HelmApps(s.helmApp.Namespace).UpdateStatus(ctx, s.helmApp, metav1.UpdateOptions{})
		return err
	})
}

func (s *Status) getPhase() v1alpha1.HelmAppStatusPhase {
	phase := v1alpha1.HelmAppStatusPhaseDetecting
	if s.isDetected() {
		phase = v1alpha1.HelmAppStatusPhaseConfiguring
	}
	if s.helmApp.Spec.PreStatus == v1alpha1.HelmAppPreStatusConfigured {
		phase = v1alpha1.HelmAppStatusPhaseInstalling
	}
	idx, condition := s.helmApp.Status.GetCondition(v1alpha1.HelmAppInstalled)
	if idx != -1 && condition.Status == corev1.ConditionTrue {
		phase = v1alpha1.HelmAppStatusPhaseInstalled
	}
	return phase
}

func (s *Status) isDetected() bool {
	types := []v1alpha1.HelmAppConditionType{
		v1alpha1.HelmAppChartReady,
		v1alpha1.HelmAppPreInstalled,
	}
	for _, t := range types {
		if !s.helmApp.Status.IsConditionTrue(t) {
			return false
		}
	}
	return true
}
