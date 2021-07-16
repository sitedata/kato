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

package volume

import (
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// MemoryFSVolume memory fs volume struct
type MemoryFSVolume struct {
	Base
}

// CreateVolume memory fs volume create volume
func (v *MemoryFSVolume) CreateVolume(define *Define) error {
	logrus.Debugf("create emptyDir volume type for: %s", v.svm.VolumePath)
	volumeMountName := fmt.Sprintf("manual%d", v.svm.ID)
	volumeMountPath := v.svm.VolumePath
	volumeReadOnly := false
	if volumeMountPath == "" {
		logrus.Warningf("service[%s]'s mount path is empty, skip create memoryfs", v.version.ServiceID)
		return nil
	}
	for _, m := range define.volumeMounts {
		if m.MountPath == volumeMountPath {
			logrus.Warningf("service[%s]'s found the same mount path: %s, skip create memoryfs", v.version.ServiceID, volumeMountPath)
			return nil
		}
	}
	vo := corev1.Volume{Name: volumeMountName} // !!!: volumeMount name of k8s model must equal to volume name of k8s model

	// V5.2  emptyDir's medium use default "" which means to use the node's default medium
	vo.EmptyDir = &corev1.EmptyDirVolumeSource{}

	// get service custom env
	es, err := v.dbmanager.TenantServiceEnvVarDao().GetServiceEnvs(v.as.ServiceID, []string{"inner"})
	if err != nil {
		logrus.Errorf("get service[%s] env failed: %s", v.as.ServiceID, err.Error())
		return err
	}
	for _, env := range es {
		// still support for memory medium
		if env.AttrName == "ES_ENABLE_EMPTYDIR_MEDIUM_MEMORY" && env.AttrValue == "true" {
			logrus.Debugf("use memory as medium of emptyDir for volume[name: %s; path: %s]", volumeMountName, volumeMountPath)
			vo.EmptyDir.Medium = corev1.StorageMediumMemory
		}
	}
	define.volumes = append(define.volumes, vo)
	vm := corev1.VolumeMount{
		MountPath: volumeMountPath,
		Name:      volumeMountName,
		ReadOnly:  volumeReadOnly,
		SubPath:   "",
	}
	define.volumeMounts = append(define.volumeMounts, vm)
	return nil
}

// CreateDependVolume empty func
func (v *MemoryFSVolume) CreateDependVolume(define *Define) error {
	return nil
}
