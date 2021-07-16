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
	"path"

	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigFileVolume config file volume struct
type ConfigFileVolume struct {
	Base
	envs          []corev1.EnvVar
	envVarSecrets []*corev1.Secret
}

// CreateVolume config file volume create volume
func (v *ConfigFileVolume) CreateVolume(define *Define) error {
	// environment variables
	configs := make(map[string]string)
	for _, sec := range v.envVarSecrets {
		for k, v := range sec.Data {
			// The priority of component environment variable is higher than the one of the application.
			if val := configs[k]; val == string(v) {
				continue
			}
			configs[k] = string(v)
		}
	}
	for _, env := range v.envs {
		configs[env.Name] = env.Value
	}
	cf, err := v.dbmanager.TenantServiceConfigFileDao().GetByVolumeName(v.as.ServiceID, v.svm.VolumeName)
	if err != nil {
		logrus.Errorf("error getting config file by volume name(%s): %v", v.svm.VolumeName, err)
		return fmt.Errorf("error getting config file by volume name(%s): %v", v.svm.VolumeName, err)
	}
	cmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.NewUUID(),
			Namespace: v.as.TenantID,
			Labels:    v.as.GetCommonLabels(),
		},
		Data: make(map[string]string),
	}
	cmap.Data[path.Base(v.svm.VolumePath)] = util.ParseVariable(cf.FileContent, configs)
	v.as.SetConfigMap(cmap)
	define.SetVolumeCMap(cmap, path.Base(v.svm.VolumePath), v.svm.VolumePath, false)
	return nil
}

// CreateDependVolume config file volume create depend volume
func (v *ConfigFileVolume) CreateDependVolume(define *Define) error {
	configs := make(map[string]string)
	for _, env := range v.envs {
		configs[env.Name] = env.Value
	}
	_, err := v.dbmanager.TenantServiceVolumeDao().GetVolumeByServiceIDAndName(v.smr.DependServiceID, v.smr.VolumeName)
	if err != nil {
		return fmt.Errorf("error getting TenantServiceVolume according to serviceID(%s) and volumeName(%s): %v",
			v.smr.DependServiceID, v.smr.VolumeName, err)
	}
	cf, err := v.dbmanager.TenantServiceConfigFileDao().GetByVolumeName(v.smr.DependServiceID, v.smr.VolumeName)
	if err != nil {
		return fmt.Errorf("error getting TenantServiceConfigFile according to volumeName(%s): %v", v.smr.VolumeName, err)
	}

	cmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      util.NewUUID(),
			Namespace: v.as.TenantID,
			Labels:    v.as.GetCommonLabels(),
		},
		Data: make(map[string]string),
	}
	cmap.Data[path.Base(v.smr.VolumePath)] = util.ParseVariable(cf.FileContent, configs)
	v.as.SetConfigMap(cmap)

	define.SetVolumeCMap(cmap, path.Base(v.smr.VolumePath), v.smr.VolumePath, false)
	return nil
}
