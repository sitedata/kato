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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gridworkz/kato/util"
	corev1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	"github.com/gridworkz/kato/pkg/generated/informers/externalversions"
	k8sutil "github.com/gridworkz/kato/util/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var ctx = context.Background()
var kubeClient clientset.Interface
var katoClient versioned.Interface
var testEnv *envtest.Environment
var stopCh = make(chan struct{})

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"HelmApp Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	projectHome := os.Getenv("PROJECT_HOME")
	kubeconfig := os.Getenv("KUBE_CONFIG")

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join(projectHome, "config", "crd")},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    util.Bool(true),
	}

	_, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())

	restConfig, err := k8sutil.NewRestConfig(kubeconfig)
	Expect(err).NotTo(HaveOccurred())

	katoClient = versioned.NewForConfigOrDie(restConfig)
	kubeClient = clientset.NewForConfigOrDie(restConfig)
	katoInformer := externalversions.NewSharedInformerFactoryWithOptions(katoClient, 10*time.Second,
		externalversions.WithNamespace(corev1.NamespaceAll))

	ctrl := NewController(ctx, stopCh, kubeClient, katoClient, katoInformer.Kato().V1alpha1().HelmApps().Informer(), katoInformer.Kato().V1alpha1().HelmApps().Lister(), "/tmp/helm/repo/repositories.yaml", "/tmp/helm/cache", "/tmp/helm/chart")
	go ctrl.Start()

	// create namespace

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the helmCmd app controller")
	close(stopCh)

	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
