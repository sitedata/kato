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

package v1

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/apps/v1"
)

func TestGetStatefulsetModifiedConfiguration(t *testing.T) {
	var replicas int32 = 1
	var replicasnew int32 = 2
	bytes, err := getStatefulsetModifiedConfiguration(&v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "teststatefulset",
			Labels: map[string]string{
				"version": "1",
			},
		},
		Spec: v1.StatefulSetSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeName: "v1",
					NodeSelector: map[string]string{
						"test": "1111",
					},
					Containers: []corev1.Container{
						corev1.Container{
							Image: "nginx",
							Name:  "nginx1",
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "version",
									Value: "V1",
								},
								corev1.EnvVar{
									Name:  "delete",
									Value: "true",
								},
							},
						},
						corev1.Container{
							Image: "nginx",
							Name:  "nginx2",
						},
					},
				},
			},
		},
	}, &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "teststatefulset",
			Labels: map[string]string{
				"version": "2",
			},
		},
		Spec: v1.StatefulSetSpec{
			Replicas: &replicasnew,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeName: "v2",
					NodeSelector: map[string]string{
						"test": "1111",
					},
					Containers: []corev1.Container{
						corev1.Container{
							Image: "nginx",
							Name:  "nginx1",
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "version",
									Value: "V2",
								},
							},
						},
						corev1.Container{
							Image: "nginx",
							Name:  "nginx3",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes))
}
