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
	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	"github.com/gridworkz/kato/pkg/helm"
	corev1 "k8s.io/api/core/v1"
)

// Detector is responsible for detecting the helm app.
type Detector struct {
	helmApp *v1alpha1.HelmApp
	repo    *helm.Repo
	app     *App
}

// NewDetector creates a new Detector.
func NewDetector(helmApp *v1alpha1.HelmApp, app *App, repo *helm.Repo) *Detector {
	return &Detector{
		helmApp: helmApp,
		repo:    repo,
		app:     app,
	}
}

// Detect detects the helm app.
func (d *Detector) Detect() error {
	// add repo
	if !d.helmApp.Status.IsConditionTrue(v1alpha1.HelmAppChartReady) {
		appStore := d.helmApp.Spec.AppStore
		if err := d.repo.Add(appStore.Name, appStore.URL, "", ""); err != nil {
			d.helmApp.Status.SetCondition(*v1alpha1.NewHelmAppCondition(
				v1alpha1.HelmAppChartReady, corev1.ConditionFalse, "RepoFailed", err.Error()))
			return err
		}
	}

	// load chart
	if !d.helmApp.Status.IsConditionTrue(v1alpha1.HelmAppChartReady) {
		err := d.app.LoadChart()
		if err != nil {
			d.helmApp.Status.UpdateCondition(v1alpha1.NewHelmAppCondition(
				v1alpha1.HelmAppChartReady, corev1.ConditionFalse, "ChartFailed", err.Error()))
			return err
		}
		d.helmApp.Status.UpdateConditionStatus(v1alpha1.HelmAppChartReady, corev1.ConditionTrue)
		return nil
	}

	// check if the chart is valid
	if !d.helmApp.Status.IsConditionTrue(v1alpha1.HelmAppPreInstalled) {
		if err := d.app.PreInstall(); err != nil {
			d.helmApp.Status.UpdateCondition(v1alpha1.NewHelmAppCondition(
				v1alpha1.HelmAppPreInstalled, corev1.ConditionFalse, "PreInstallFailed", err.Error()))
			return err
		}
		d.helmApp.Status.UpdateConditionStatus(v1alpha1.HelmAppPreInstalled, corev1.ConditionTrue)
		return nil
	}

	return nil
}
