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
	"github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/common"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cueTemplate = `
output: {
	apiVersion: "kato.io/v1alpha1"
	kind:       "ThirdComponent"
	metadata: {
		name: context.componentID
		namespace: context.namespace
	}
	spec: {
		endpointSource: {
			if parameter["kubernetes"] != _|_ {
				kubernetesService: {
					namespace: parameter["kubernetes"]["namespace"],
					name: parameter["kubernetes"]["name"]
				}
			}
		}
		if parameter["port"] != _|_ {
			ports: parameter["port"]
		}
	}
}

parameter: {
	kubernetes?: {
		namespace?: string
		name: string
	}
	port?: [...{
		name:   string
		port:   >0 & <=65533
		openInner: bool
		openOuter: bool
	}]
}
`
var thirdComponetDefineName = "core-thirdcomponent"
var thirdComponetDefine = v1alpha1.ComponentDefinition{
	TypeMeta: v1.TypeMeta{
		Kind:       "ComponentDefinition",
		APIVersion: "kato.io/v1alpha1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name: thirdComponetDefineName,
		Annotations: map[string]string{
			"definition.oam.dev/description": "Kato built-in component type that defines third-party service components.",
		},
	},
	Spec: v1alpha1.ComponentDefinitionSpec{
		Workload: common.WorkloadTypeDescriptor{
			Type: "ThirdComponent",
			Definition: common.WorkloadGVK{
				APIVersion: "kato.io/v1alpha1",
				Kind:       "ThirdComponent",
			},
		},
		Schematic: &v1alpha1.Schematic{
			CUE: &common.CUE{
				Template: cueTemplate,
			},
		},
	},
}
