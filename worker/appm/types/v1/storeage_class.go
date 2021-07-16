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
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// kind: StorageClass
// apiVersion: storage.k8s.io/v1
// metadata:
//   name: local-storage
// provisioner: kubernetes.io/no-provisioner
// volumeBindingMode: WaitForFirstConsumer

var initStorageClass []*storagev1.StorageClass

//KatoStatefuleShareStorageClass kato support statefulset app share volume
var KatoStatefuleShareStorageClass = "katosssc"

//KatoStatefuleLocalStorageClass kato support statefulset app local volume
var KatoStatefuleLocalStorageClass = "katoslsc"

func init() {
	var volumeBindingImmediate = storagev1.VolumeBindingImmediate
	var columeWaitForFirstConsumer = storagev1.VolumeBindingWaitForFirstConsumer
	var Retain = v1.PersistentVolumeReclaimRetain
	initStorageClass = append(initStorageClass, &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: KatoStatefuleShareStorageClass,
		},
		Provisioner:       "kato.io/provisioner-sssc",
		VolumeBindingMode: &volumeBindingImmediate,
		ReclaimPolicy:     &Retain,
	})
	initStorageClass = append(initStorageClass, &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: KatoStatefuleLocalStorageClass,
		},
		Provisioner:       "kato.io/provisioner-sslc",
		VolumeBindingMode: &columeWaitForFirstConsumer,
		ReclaimPolicy:     &Retain,
	})
}

//GetInitStorageClass get init storageclass list
func GetInitStorageClass() []*storagev1.StorageClass {
	return initStorageClass
}
