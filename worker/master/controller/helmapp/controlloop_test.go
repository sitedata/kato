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

	katov1alpha1 "github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	"github.com/gridworkz/kato/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ControlLoop", func() {
	var namespace string
	var helmApp *katov1alpha1.HelmApp
	BeforeEach(func() {
		// create namespace
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: util.NewUUID(),
			},
		}
		namespace = ns.Name
		By("create namespace: " + namespace)
		_, err := kubeClient.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		helmApp = &katov1alpha1.HelmApp{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "phpmyadmin",
				Namespace: namespace,
				Labels: map[string]string{
					"app": "phpmyadmin",
				},
			},
			Spec: katov1alpha1.HelmAppSpec{
				EID:          "5bfba91b0ead72f612732535ef802217",
				TemplateName: "phpmyadmin",
				Version:      "8.2.0",
				AppStore: &katov1alpha1.HelmAppStore{
					Name: "bitnami",
					URL:  "https://charts.bitnami.com/bitnami",
				},
			},
		}
		By("create helm app: " + helmApp.Name)
		_, err = katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Create(context.Background(), helmApp, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		By("delete namespace: " + namespace)
		err := kubeClient.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Reconcile", func() {
		Context("HelmApp created", func() {
			It("should fulfill default values", func() {
				watch, err := katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Watch(context.Background(), metav1.ListOptions{
					LabelSelector: "app=phpmyadmin",
					Watch:         true,
				})
				Expect(err).NotTo(HaveOccurred())

				By("wait until the default values of the helm app were setup")
				for event := range watch.ResultChan() {
					newHelmApp := event.Object.(*katov1alpha1.HelmApp)
					// wait status
					for _, conditionType := range defaultConditionTypes {
						_, condition := newHelmApp.Status.GetCondition(conditionType)
						if condition == nil {
							break
						}
					}
					if newHelmApp.Status.Phase == "" {
						continue
					}

					// wait spec
					if newHelmApp.Spec.PreStatus == "" {
						continue
					}

					break
				}
			})

			It("should start detecting", func() {
				newHelmApp, err := katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Get(context.Background(), helmApp.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				Expect(newHelmApp.Status.Phase).NotTo(Equal(katov1alpha1.HelmAppStatusPhaseDetecting))

				By("wait until condition detecting conditions become true")
				watch, err := katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Watch(context.Background(), metav1.ListOptions{
					LabelSelector: "app=phpmyadmin",
					Watch:         true,
				})
				Expect(err).NotTo(HaveOccurred())

				conditionTypes := []katov1alpha1.HelmAppConditionType{
					katov1alpha1.HelmAppChartReady,
					katov1alpha1.HelmAppPreInstalled,
				}

				for event := range watch.ResultChan() {
					newHelmApp = event.Object.(*katov1alpha1.HelmApp)
					isFinished := true
					for _, conditionType := range conditionTypes {
						_, condition := newHelmApp.Status.GetCondition(conditionType)
						if condition == nil || condition.Status == corev1.ConditionFalse {
							isFinished = false
							break
						}
					}
					if isFinished {
						break
					}
				}
			})

			It("should start configuring", func() {
				By("wait until phase become configuring")
				err := waitUntilConfiguring(helmApp)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Install HelmApp", func() {
			It("should ok", func() {
				err := waitUntilConfiguring(helmApp)
				Expect(err).NotTo(HaveOccurred())

				newHelmApp, err := katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Get(context.Background(), helmApp.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				By("install helm app: " + helmApp.Name)
				newHelmApp.Spec.PreStatus = katov1alpha1.HelmAppPreStatusConfigured
				_, err = katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Update(context.Background(), newHelmApp, metav1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())

				err = waitUntilInstalled(helmApp)
				Expect(err).NotTo(HaveOccurred())

				err = waitUntilDeployed(helmApp)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

func waitUntilConfiguring(helmApp *katov1alpha1.HelmApp) error {
	_, err := waitPhaseUntil(helmApp, katov1alpha1.HelmAppStatusPhaseConfiguring)
	return err
}

func waitUntilInstalled(helmApp *katov1alpha1.HelmApp) error {
	_, err := waitPhaseUntil(helmApp, katov1alpha1.HelmAppStatusPhaseInstalled)
	return err
}

func waitPhaseUntil(helmApp *katov1alpha1.HelmApp, phase katov1alpha1.HelmAppStatusPhase) (*katov1alpha1.HelmApp, error) {
	watch, err := katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: "app=phpmyadmin",
		Watch:         true,
	})
	if err != nil {
		return nil, err
	}

	// TODO: timeout
	for event := range watch.ResultChan() {
		newHelmApp := event.Object.(*katov1alpha1.HelmApp)
		if newHelmApp.Status.Phase == phase {
			return newHelmApp, nil
		}
	}

	return nil, nil
}

func waitUntilDeployed(helmApp *katov1alpha1.HelmApp) error {
	return waitStatusUntil(helmApp, katov1alpha1.HelmAppStatusDeployed)
}

func waitStatusUntil(helmApp *katov1alpha1.HelmApp, status katov1alpha1.HelmAppStatusStatus) error {
	watch, err := katoClient.KatoV1alpha1().HelmApps(helmApp.Namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: "app=phpmyadmin",
		Watch:         true,
	})
	if err != nil {
		return err
	}

	// TODO: timeout
	for event := range watch.ResultChan() {
		newHelmApp := event.Object.(*katov1alpha1.HelmApp)
		if newHelmApp.Status.Status == status {
			return nil
		}
	}

	return nil
}
