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

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//CreateKubeService create kube service
func CreateKubeService(client kubernetes.Interface, namespace string, services ...*corev1.Service) error {
	var retryService []*corev1.Service
	for i := range services {
		createService := services[i]
		if _, err := client.CoreV1().Services(namespace).Create(context.Background(), createService, metav1.CreateOptions{}); err != nil {
			// Ignore if the Service is invalid with this error message:
			// 	Service "kube-dns" is invalid: spec.clusterIP: Invalid value: "10.96.0.10": provided IP is already allocated
			if !errors.IsAlreadyExists(err) && !errors.IsInvalid(err) {
				retryService = append(retryService, createService)
				continue
			}
			if _, err := client.CoreV1().Services(namespace).Update(context.Background(), createService, metav1.UpdateOptions{}); err != nil {
				retryService = append(retryService, createService)
				continue
			}
		}
	}
	//second attempt
	for _, service := range retryService {
		_, err := client.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				continue
			}
			return err
		}
	}
	return nil
}
