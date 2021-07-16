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

package provider

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/util"
	"github.com/gridworkz/kato/worker/master/volumes/provider/lib/controller"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type katossscProvisioner struct {
	// The directory to create PV-backing directories in
	pvDir string
	name  string
}

// NewKatossscProvisioner creates a new Kato statefulset share volume provisioner
func NewKatossscProvisioner() controller.Provisioner {
	sharePath := os.Getenv("SHARE_DATA_PATH")
	if sharePath == "" {
		sharePath = "/grdata"
	}
	return &katossscProvisioner{
		pvDir: sharePath,
		name:  "kato.io/provisioner-sssc",
	}
}

var _ controller.Provisioner = &katossscProvisioner{}

// Provision creates a storage asset and returns a PV object representing it.
func (p *katossscProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	logrus.Debugf("[katossscProvisioner] start creating PV object. paramters: %+v", options.Parameters)

	tenantID := options.PVC.Labels["tenant_id"]
	serviceID := options.PVC.Labels["service_id"]
	_, stateless := options.PVC.Labels["stateless"]
	// v5.0.4 Previous versions
	hostpath := path.Join(p.pvDir, "tenant", tenantID, "service", serviceID, options.PVC.Name)

	// after v5.0.4,change host path
	// Directory path has nothing to do with volume ID
	// Directory path bound to volume mount path
	if util.DirIsEmpty(hostpath) {
		podName := getPodNameByPVCName(options.PVC.Name)
		volumeID := getVolumeIDByPVCName(options.PVC.Name)
		if volumeID != 0 {
			volume, err := db.GetManager().TenantServiceVolumeDao().GetVolumeByID(volumeID)
			if err != nil {
				logrus.Errorf("get volume by id %d failure %s", volumeID, err.Error())
				return nil, err
			}
			hostpath = volume.HostPath
			if !stateless {
				hostpath = path.Join(volume.HostPath, podName)
			}
		} else {
			return nil, fmt.Errorf("can not parse volume id")
		}
	}
	if err := util.CheckAndCreateDirByMode(hostpath, 0777); err != nil {
		return nil, err
	}
	// new volume path
	persistentVolumeSource, err := updatePathForPersistentVolumeSource(&options.PersistentVolumeSource, hostpath)
	if err != nil {
		return nil, err
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   options.PVName,
			Labels: options.PVC.Labels,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceStorage: options.PVC.Spec.Resources.Requests[v1.ResourceStorage],
			},
			MountOptions:           options.MountOptions,
			PersistentVolumeSource: *persistentVolumeSource,
		},
	}
	logrus.Infof("create katosssc pv %s for pvc %s", pv.Name, options.PVC.Name)
	return pv, nil
}

// Delete removes the storage asset that was created by Provision represented
// by the given PV.
func (p *katossscProvisioner) Delete(volume *v1.PersistentVolume) error {

	return nil
}

func (p *katossscProvisioner) Name() string {
	return p.name
}

func getPodNameByPVCName(pvcName string) string {
	pvcNames := strings.SplitN(pvcName, "-", 2)
	if len(pvcNames) == 2 {
		return pvcNames[1]
	}
	return pvcName
}

func getVolumeIDByPVCName(pvcName string) int {
	logrus.Debug("parse volume id from pvc name", pvcName)
	pvcNames := strings.SplitN(pvcName, "-", 2)
	if len(pvcNames) == 2 {
		fmt.Println(idStr)
		id, _ := strconv.Atoi(idStr)
		return id
	}
	if strings.HasPrefix(pvcName, "manual") {
		idStr := strings.TrimPrefix(pvcName, "manual")
		id, _ := strconv.Atoi(idStr)
		return id
	}
	return 0
}

func updatePathForPersistentVolumeSource(persistentVolumeSource *v1.PersistentVolumeSource, hostpath string) (*v1.PersistentVolumeSource, error) {
	newPath := func(new string) string {
		p := strings.Replace(hostpath, "/grdata", "", 1)
		return path.Join(new, p)
	}
	source := &v1.PersistentVolumeSource{}
	switch {
	case persistentVolumeSource.NFS != nil:
		source.NFS = persistentVolumeSource.NFS
		source.NFS.Path = newPath(persistentVolumeSource.NFS.Path)
	case persistentVolumeSource.CSI != nil && persistentVolumeSource.CSI.Driver == "nasplugin.csi.alibabacloud.com":
		// convert aliyun nas to nfs
		if persistentVolumeSource.CSI.VolumeAttributes != nil {
			source.NFS = &v1.NFSVolumeSource{
				Server: persistentVolumeSource.CSI.VolumeAttributes["server"],
				Path:   newPath(persistentVolumeSource.CSI.VolumeAttributes["path"]),
			}
		}
	case persistentVolumeSource.Glusterfs != nil:
		//glusterfs:
		//	endpoints: glusterfs-cluster
		//	path: myVol1
		glusterfs := &v1.GlusterfsPersistentVolumeSource{
			EndpointsName:      persistentVolumeSource.Glusterfs.EndpointsName,
			EndpointsNamespace: persistentVolumeSource.Glusterfs.EndpointsNamespace,
			Path:               newPath(persistentVolumeSource.Glusterfs.Path),
		}
		source.Glusterfs = glusterfs
	default:
		return nil, fmt.Errorf("unsupported persistence volume source")
	}
	return source, nil
}
