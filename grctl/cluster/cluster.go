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

package cluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/gridworkz/kato-operator/api/v1alpha1"
	"github.com/gridworkz/kato/grctl/clients"
	"github.com/oam-dev/kubevela/pkg/utils/apply"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"
)

// Cluster represents a kato cluster.
type Cluster struct {
	katoCluster *v1alpha1.KatoCluster
	namespace       string
	newVersion      string
}

// NewCluster creates new cluster.
func NewCluster(namespace, newVersion string) (*Cluster, error) {
	katoCluster, err := getKatoCluster(namespace)
	if err != nil {
		return nil, err
	}
	return &Cluster{
		katoCluster: katoCluster,
		namespace:       namespace,
		newVersion:      newVersion,
	}, nil
}

// Upgrade upgrade cluster.
func (c *Cluster) Upgrade() error {
	logrus.Infof("upgrade cluster from %s to %s", c.katoCluster.Spec.InstallVersion, c.newVersion)

	if errs := c.createCrds(); len(errs) > 0 {
		return errors.New(strings.Join(errs, ","))
	}

	if errs := c.updateRbdComponents(); len(errs) > 0 {
		return fmt.Errorf("update kato components: %s", strings.Join(errs, ","))
	}

	if err := c.updateCluster(); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) createCrds() []string {
	crds := c.getCrds()
	if crds == nil {
		return nil
	}
	logrus.Info("start creating crds")
	var errs []string
	for _, crd := range crds {
		if err := c.createCrd(crd); err != nil {
			errs = append(errs, err.Error())
		}
	}
	logrus.Info("crds applyed")
	return errs
}

func (c *Cluster) createCrd(crdStr string) error {
	var crd apiextensionsv1beta1.CustomResourceDefinition
	if err := yaml.Unmarshal([]byte(crdStr), &crd); err != nil {
		return fmt.Errorf("unmarshal crd: %v", err)
	}
	applyer := apply.NewAPIApplicator(clients.KatoKubeClient)
	if err := applyer.Apply(context.Background(), &crd); err != nil {
		return fmt.Errorf("apply crd: %v", err)
	}
	return nil
}

func (c *Cluster) getCrds() []string {
	for v, versionConfig := range versions {
		if strings.Contains(c.newVersion, v) {
			return versionConfig.CRDs
		}
	}
	return nil
}

func (c *Cluster) updateRbdComponents() []string {
	componentNames := []string{
		"rbd-api",
		"rbd-chaos",
		"rbd-mq",
		"rbd-eventlog",
		"rbd-gateway",
		"rbd-node",
		"rbd-resource-proxy",
		"rbd-webcli",
		"rbd-worker",
	}
	var errs []string
	for _, name := range componentNames {
		err := c.updateRbdComponent(name)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	return errs
}

func (c *Cluster) updateRbdComponent(name string) error {
	var cpt v1alpha1.RbdComponent
	err := clients.KatoKubeClient.Get(context.Background(),
		types.NamespacedName{Namespace: c.namespace, Name: name}, &cpt)
	if err != nil {
		return fmt.Errorf("get rbdcomponent %s: %v", name, err)
	}

	ref, err := reference.Parse(cpt.Spec.Image)
	if err != nil {
		return fmt.Errorf("parse image %s: %v", cpt.Spec.Image, err)
	}
	repo := ref.(reference.Named)
	newImage := repo.Name() + ":" + c.newVersion

	oldImageName := cpt.Spec.Image
	cpt.Spec.Image = newImage
	if err := clients.KatoKubeClient.Update(context.Background(), &cpt); err != nil {
		return fmt.Errorf("update rbdcomponent %s: %v", name, err)
	}

	logrus.Infof("update rbdcomponent %s \nfrom %s \nto   %s", name, oldImageName, newImage)
	return nil
}

func (c *Cluster) updateCluster() error {
	c.katoCluster.Spec.InstallVersion = c.newVersion
	if err := clients.KatoKubeClient.Update(context.Background(), c.katoCluster); err != nil {
		return fmt.Errorf("update kato cluster: %v", err)
	}
	logrus.Infof("update kato cluster to %s", c.newVersion)
	return nil
}

func getKatoCluster(namespace string) (*v1alpha1.KatoCluster, error) {
	var cluster v1alpha1.KatoCluster
	err := clients.KatoKubeClient.Get(context.Background(),
		types.NamespacedName{Namespace: namespace, Name: "katocluster"}, &cluster)
	if err != nil {
		return nil, fmt.Errorf("get kato cluster: %v", err)
	}
	return &cluster, nil
}
