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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"context"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/dao"
	"github.com/gridworkz/kato/node/nodem/client"
	"github.com/gridworkz/kato/worker/appm/store"
	"github.com/gridworkz/kato/worker/master/volumes/provider/lib/controller"

	"k8s.io/client-go/kubernetes"

	"github.com/sirupsen/logrus"

	httputil "github.com/gridworkz/kato/util/http"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type katosslcProvisioner struct {
	// The directory to create PV-backing directories in
	name    string
	kubecli kubernetes.Interface
	store   store.Storer
}

// NewKatosslcProvisioner creates a new Kato statefulset share volume provisioner
func NewKatosslcProvisioner (kubecli kubernetes.Interface, store store.Storer) controller.Provisioner {
	return &katosslcProvisioner{
		name:    "kato.io/provisioner-sslc",
		kubecli: kubecli,
		store: store,
	}
}

var _ controller.Provisioner = & katosslcProvisioner {}

//selectNode select an appropriate node with the largest resource surplus
func (p *katosslcProvisioner) selectNode(ctx context.Context, nodeOS, ignore string) (*v1.Node, error) {
	allnode, err := p.kubecli.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var maxavailable int64
	var selectnode *v1.Node
	for _, node := range allnode.Items {
		nodeReady := false
		if node.Labels[client.LabelOS] != nodeOS {
			continue
		}

		// filter out ignore nodes
		if strings.Contains(ignore, node.Name) {
			logrus.Debugf("[katosslcProvisioner] [selectNode] ignore node %s based on %s", node.Name, ignore)
			continue
		}

		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady {
				nodeReady = true
				if condition.Status == v1.ConditionTrue {
					ip := ""
					for _, address := range node.Status.Addresses {
						if address.Type == v1.NodeInternalIP {
							ip = address.Address
							break
						}
					}
					if ip == "" {
						logrus.Warningf("Node: %s; node internal address not found", node.Name)
						break
					}
					//only contains kato pod
					//pods, err := p.store.GetPodLister().Pods(v1.NamespaceAll).List(labels.NewSelector())
					pods, err := p.kubecli.CoreV1().Pods(v1.NamespaceAll).List(ctx, metav1.ListOptions{
						FieldSelector: "spec.nodeName=" + node.Name,
					})
					if err != nil {
						logrus.Errorf("list pods list from node ip error %s", err.Error())
						break
					}
					var nodeUsedMemory int64
					for _, pod := range pods.Items {
						for _, con := range pod.Spec.Containers {
							memory := con.Resources.Requests.Memory()
							nodeUsedMemory += memory.Value()
						}
					}
					available := node.Status.Allocatable.Memory().Value() - nodeUsedMemory
					if available >= maxavailable {
						logrus.Infof("select node: %s", node.Name)
						maxavailable = available
						selectnode = node.DeepCopy()
					} else {
						logrus.Infof("Node: %s; node available memory(%d) is less than max available "+
							"memory(%d)", node.Name, available, maxavailable)
					}
				}
			}
		}
		if !nodeReady {
			logrus.Warningf("Node: %s; not ready", node.Name)
		}
	}
	return selectnode, nil
}
func (p *katosslcProvisioner) createPath(options controller.VolumeOptions) (string, error) {
	tenantID := options.PVC.Labels["tenant_id"]
	serviceID := options.PVC.Labels["service_id"]
	volumeID := getVolumeIDByPVCName(options.PVC.Name)
	if volumeID != 0 {
		volume, err := db.GetManager().TenantServiceVolumeDao().GetVolumeByID(volumeID)
		if err != nil {
			logrus.Warningf("get volume by id %d failure %s", volumeID, err.Error())
			return "", err
		}
		reqoptions := map[string]string{
			"tenant_id":   tenantID,
			"service_id":  serviceID,
			"pvcname":     options.PVC.Name,
			"volume_name": volume.VolumeName,
			"pod_name":    getPodNameByPVCName(options.PVC.Name),
		}
		var ip string
		for _, address := range options.SelectedNode.Status.Addresses {
			if address.Type == v1.NodeInternalIP {
				ip = address.Address
			}
		}
		if ip == "" {
			return "", fmt.Errorf("do not find node ip")
		}
		retry := 3
		var path string
		for retry > 0 {
			retry--
			body := bytes.NewBuffer(nil)
			if err := json.NewEncoder(body).Encode(reqoptions); err != nil {
				return "", fmt.Errorf("create volume body failure %s", err.Error())
			}
			res, err := http.Post(fmt.Sprintf("http://%s:6100/v2/localvolumes/create", ip), "application/json", body)
			if err != nil {
				logrus.Errorf("do request node api failure %s", err.Error())
			}
			if res != nil && res.StatusCode == 200 && res.Body != nil {
				if res, err := httputil.ParseResponseBody(res.Body, "application/json"); err == nil {
					if info, ok := res.Bean.(map[string]interface{}); ok {
						path = info["path"].(string)
						break
					} else {
						logrus.Errorf("request create local volume failure: parse body info failure  ")
					}
				} else {
					logrus.Errorf("request create local volume failure: parse body failure %s ", err.Error())
				}
			}
			if res != nil {
				logrus.Errorf("request create local volume failure code:%d", res.StatusCode)
			}
			time.Sleep(time.Second * 2)
		}
		return path, nil
	}
	return "", fmt.Errorf("can not parse volume id")
}

// Provision creates a storage asset and returns a PV object representing it.
func (p *katosslcProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	logrus.Debugf("[katosslcProvisioner] start creating PV object. paramters: %+v", options.Parameters)
	//runtime select an appropriate node with the largest resource surplus
	// storageclass VolumeBinding set WaitForFirstConsumer, SelectedNode should be assigned.
	if options.SelectedNode == nil {
		var err error
		var ignoreNodes string
		if options.Parameters != nil {
			ignoreNodes = options.Parameters["ignoreNodes"]
		}
		options.SelectedNode, err = p.selectNode(context.Background(), options.PVC.Annotations[client.LabelOS], ignoreNodes)
		if err != nil {
			return nil, fmt.Errorf("node OS: %s; error selecting node: %s", options.PVC.Annotations[client.LabelOS], err.Error())
		}
		if options.SelectedNode == nil {
			return nil, fmt.Errorf("do not select an appropriate node for local volume")
		}
		if _, ok := options.SelectedNode.Labels["kubernetes.io/hostname"]; !ok {
			return nil, fmt.Errorf("select node(%s) do not have label kubernetes.io/hostname ", options.SelectedNode.Name)
		}
	}
	path, err := p.createPath(options)
	if err != nil {
		if err == dao.ErrVolumeNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("create local volume from node %s failure %s", options.SelectedNode.Name, err.Error())
	}
	if path == "" {
		return nil, fmt.Errorf("create local volume failure,local path is not create")
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
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: path,
				},
			},
			MountOptions: options.MountOptions,
			NodeAffinity: &v1.VolumeNodeAffinity{
				Required: &v1.NodeSelector{
					NodeSelectorTerms: []v1.NodeSelectorTerm{
						{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      "kubernetes.io/hostname",
									Operator: v1.NodeSelectorOpIn,
									Values:   []string{options.SelectedNode.Labels["kubernetes.io/hostname"]},
								},
							},
						},
					},
				},
			},
		},
	}
	logrus.Infof("create katosslc pv %s for pvc %s", pv.Name, options.PVC.Name)
	return pv, nil
}

// Delete removes the storage asset that was created by Provision represented
// by the given PV.
func (p *katosslcProvisioner) Delete(volume *v1.PersistentVolume) error {

	return nil
}

func (p *katosslcProvisioner) Name() string {
	return p.name
}
