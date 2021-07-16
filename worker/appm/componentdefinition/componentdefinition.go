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

package componentdefinition

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	katoversioned "github.com/gridworkz/kato/pkg/generated/clientset/versioned"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ErrNotSupport = fmt.Errorf("not supported component definition")
var ErrOnlyCUESupport = fmt.Errorf("component definition only supports cue template")

type ComponentDefinitionBuilder struct {
	definitions map[string]*v1alpha1.ComponentDefinition
	namespace   string
	lock        sync.Mutex
}

var componentDefinitionBuilder *ComponentDefinitionBuilder

func NewComponentDefinitionBuilder(namespace string) *ComponentDefinitionBuilder {
	componentDefinitionBuilder = &ComponentDefinitionBuilder{
		definitions: make(map[string]*v1alpha1.ComponentDefinition),
		namespace:   namespace,
	}
	return componentDefinitionBuilder
}

func GetComponentDefinitionBuilder() *ComponentDefinitionBuilder {
	return componentDefinitionBuilder
}

func (c *ComponentDefinitionBuilder) OnAdd(obj interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	cd, ok := obj.(*v1alpha1.ComponentDefinition)
	if ok {
		logrus.Infof("load componentdefinition %s", cd.Name)
		c.definitions[cd.Name] = cd
	}
}
func (c *ComponentDefinitionBuilder) OnUpdate(oldObj, newObj interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	cd, ok := newObj.(*v1alpha1.ComponentDefinition)
	if ok {
		logrus.Infof("update componentdefinition %s", cd.Name)
		c.definitions[cd.Name] = cd
	}
}
func (c *ComponentDefinitionBuilder) OnDelete(obj interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	cd, ok := obj.(*v1alpha1.ComponentDefinition)
	if ok {
		logrus.Infof("delete componentdefinition %s", cd.Name)
		delete(c.definitions, cd.Name)
	}
}

func (c *ComponentDefinitionBuilder) GetComponentDefinition(name string) *v1alpha1.ComponentDefinition {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.definitions[name]
}

func (c *ComponentDefinitionBuilder) GetComponentProperties(as *v1.AppService, dbm db.Manager, cd *v1alpha1.ComponentDefinition) interface{} {
	//TODO: support custom component properties
	switch cd.Name {
	case thirdComponetDefineName:
		properties := &ThirdComponentProperties{}
		tpsd, err := dbm.ThirdPartySvcDiscoveryCfgDao().GetByServiceID(as.ServiceID)
		if err != nil {
			logrus.Errorf("query component %s third source config failure %s", as.ServiceID, err.Error())
		}
		if tpsd != nil {
			// support other source type
			if tpsd.Type == dbmodel.DiscorveryTypeKubernetes.String() {
				properties.Kubernetes = ThirdComponentKubernetes{
					Name:      tpsd.ServiceName,
					Namespace: tpsd.Namespace,
				}
			}
		}
		ports, err := dbm.TenantServicesPortDao().GetPortsByServiceID(as.ServiceID)
		if err != nil {
			logrus.Errorf("query component %s ports failure %s", as.ServiceID, err.Error())
		}

		for _, port := range ports {
			properties.Port = append(properties.Port, &ThirdComponentPort{
				Port:      port.ContainerPort,
				Name:      strings.ToLower(port.PortAlias),
				OpenInner: *port.IsInnerService,
				OpenOuter: *port.IsOuterService,
			})
		}
		if properties.Port == nil {
			properties.Port = []*ThirdComponentPort{}
		}
		return properties
	default:
		return nil
	}
}

func (c *ComponentDefinitionBuilder) BuildWorkloadResource(as *v1.AppService, dbm db.Manager) error {
	cd := c.GetComponentDefinition(as.GetComponentDefinitionName())
	if cd == nil {
		return ErrNotSupport
	}
	if cd.Spec.Schematic == nil || cd.Spec.Schematic.CUE == nil {
		return ErrOnlyCUESupport
	}
	ctx := NewTemplateContext(as, cd.Spec.Schematic.CUE.Template, c.GetComponentProperties(as, dbm, cd))
	manifests, err := ctx.GenerateComponentManifests()
	if err != nil {
		return err
	}
	ctx.SetContextValue(manifests)
	as.SetManifests(manifests)
	if len(manifests) > 0 {
		as.SetWorkload(manifests[0])
	}
	return nil
}

//InitCoreComponentDefinition init the built-in component type definition.
//Should be called after the store is initialized.
func (c *ComponentDefinitionBuilder) InitCoreComponentDefinition(katoClient katoversioned.Interface) {
	coreComponentDefinition := []*v1alpha1.ComponentDefinition{&thirdComponetDefine}
	for _, ccd := range coreComponentDefinition {
		if c.GetComponentDefinition(ccd.Name) == nil {
			logrus.Infof("create core componentdefinition %s", ccd.Name)
			if _, err := katoClient.KatoV1alpha1().ComponentDefinitions(c.namespace).Create(context.Background(), ccd, metav1.CreateOptions{}); err != nil {
				logrus.Errorf("create core componentdefinition %s failure %s", ccd.Name, err.Error())
			}
		}
	}
	logrus.Infof("success check core componentdefinition from cluster")
}
