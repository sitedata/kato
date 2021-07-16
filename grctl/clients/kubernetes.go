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

package clients

import (
	"fmt"
	"os"
	"path"

	"github.com/gridworkz/kato-operator/pkg/generated/clientset/versioned"
	"github.com/gridworkz/kato/builder/sources"
	k8sutil "github.com/gridworkz/kato/util/k8s"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

//K8SClient
var K8SClient kubernetes.Interface

//KatoKubeClient kato custom resource client
var KatoKubeClient versioned.Interface

//InitClient init k8s client
func InitClient(kubeconfig string) error {
	if kubeconfig == "" {
		homePath, _ := sources.Home()
		kubeconfig = path.Join(homePath, ".kube/config")
	}
	_, err := os.Stat(kubeconfig)
	if err != nil {
		fmt.Printf("Please make sure the kube-config file(%s) exists\n", kubeconfig)
		os.Exit(1)
	}
	// use the current context in kubeconfig
	config, err := k8sutil.NewRestConfig(kubeconfig)
	if err != nil {
		return err
	}
	config.QPS = 50
	config.Burst = 100

	K8SClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Error("Create kubernetes client error.", err.Error())
		return err
	}
	KatoKubeClient = versioned.NewForConfigOrDie(config)
	return nil
}
