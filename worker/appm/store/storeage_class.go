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

package store

import (
	"context"

	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//InitStorageclass init storage class
func (a *appRuntimeStore) initStorageclass() error {
	for _, storageclass := range v1.GetInitStorageClass() {
		old, err := a.conf.KubeClient.StorageV1().StorageClasses().Get(context.Background(), storageclass.Name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				_, err = a.conf.KubeClient.StorageV1().StorageClasses().Create(context.Background(), storageclass, metav1.CreateOptions{})
			}
			if err != nil {
				return err
			}
			logrus.Info("create storageclass %s", storageclass.Name)
		} else {
			update := false
			if old.VolumeBindingMode == nil {
				update = true
			}
			if !update && old.ReclaimPolicy == nil {
				update = true
			}
			if !update && string(*old.VolumeBindingMode) != string(*storageclass.VolumeBindingMode) {
				update = true
			}
			if !update && string(*old.ReclaimPolicy) != string(*storageclass.ReclaimPolicy) {
				update = true
			}
			if update {
				err := a.conf.KubeClient.StorageV1().StorageClasses().Delete(context.Background(), storageclass.Name, metav1.DeleteOptions{})
				if err == nil {
					_, err := a.conf.KubeClient.StorageV1().StorageClasses().Create(context.Background(), storageclass, metav1.CreateOptions{})
					if err != nil {
						logrus.Errorf("recreate strageclass %s failure %s", storageclass.Name, err.Error())
					}
					logrus.Infof("update storageclass %s success", storageclass.Name)
				} else {
					logrus.Errorf("recreate strageclass %s failure %s", err.Error())
				}
			}
		}
	}
	return nil
}

