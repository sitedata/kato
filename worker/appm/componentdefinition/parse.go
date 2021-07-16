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
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/oam-dev/kubevela/pkg/cue/model"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// ParameterTag is the keyword in CUE template to define users' input
	ParameterTag = "parameter"
	// OutputFieldName is the reference of context base object
	OutputFieldName = "output"
	// OutputsFieldName is the reference of context Auxiliaries
	OutputsFieldName = "outputs"
	// ConfigFieldName is the reference of context config
	ConfigFieldName = "config"
	// ContextName is the name of context
	ContextName = "name"
	// ContextAppName is the appName of context
	ContextAppName = "appName"
	// ContextID is the componentID of context
	ContextID = "componentID"
	// ContextAppID is the appID of context
	ContextAppID = "appID"
	// ContextNamespace is the namespace of the app
	ContextNamespace = "namespace"
)

type TemplateContext struct {
	as            *v1.AppService
	componentName string
	appName       string
	componentID   string
	appID         string
	namespace     string
	template      string
	params        interface{}
}

func NewTemplateContext(as *v1.AppService, template string, params interface{}) *TemplateContext {
	return &TemplateContext{
		as:            as,
		componentName: as.ServiceAlias,
		appName:       as.AppID,
		componentID:   as.ServiceID,
		appID:         as.AppID,
		namespace:     as.TenantID,
		template:      template,
		params:        params,
	}
}

func (c *TemplateContext) GenerateComponentManifests() ([]*unstructured.Unstructured, error) {
	bi := build.NewContext().NewInstance("", nil)
	if err := bi.AddFile("-", c.template); err != nil {
		return nil, errors.WithMessagef(err, "invalid cue template of component %s", c.componentID)
	}
	var paramFile = "parameter: {}"
	if c.params != nil {
		bt, err := json.Marshal(c.params)
		if err != nil {
			return nil, errors.WithMessagef(err, "marshal parameter of component %s", c.componentID)
		}
		if string(bt) != "null" {
			paramFile = fmt.Sprintf("%s: %s", ParameterTag, string(bt))
		}
	}
	if err := bi.AddFile("parameter", paramFile); err != nil {
		return nil, errors.WithMessagef(err, "invalid parameter of component %s", c.componentID)
	}

	if err := bi.AddFile("-", c.ExtendedContextFile()); err != nil {
		return nil, err
	}
	var r cue.Runtime
	inst, err := r.Build(bi)
	if err != nil {
		return nil, err
	}
	if err := inst.Value().Validate(); err != nil {
		return nil, errors.WithMessagef(err, "invalid cue template of component %s after merge parameter and context", c.componentID)
	}

	output := inst.Lookup(OutputFieldName)

	base, err := model.NewBase(output)
	if err != nil {
		return nil, errors.WithMessagef(err, "invalid output of component %s", c.componentID)
	}
	workload, err := base.Unstructured()
	if err != nil {
		return nil, errors.WithMessagef(err, "invalid output of component %s", c.componentID)
	}

	manifests := []*unstructured.Unstructured{workload}

	outputs := inst.Lookup(OutputsFieldName)
	if !outputs.Exists() {
		return manifests, nil
	}
	st, err := outputs.Struct()
	if err != nil {
		return nil, errors.WithMessagef(err, "invalid outputs of workload %s", c.componentID)
	}
	for i := 0; i < st.Len(); i++ {
		fieldInfo := st.Field(i)
		if fieldInfo.IsDefinition || fieldInfo.IsHidden || fieldInfo.IsOptional {
			continue
		}
		other, err := model.NewOther(fieldInfo.Value)
		if err != nil {
			return nil, errors.WithMessagef(err, "invalid outputs(%s) of workload %s", fieldInfo.Name, c.componentID)
		}
		othermanifest, err := other.Unstructured()
		if err != nil {
			return nil, errors.WithMessagef(err, "invalid outputs(%s) of workload %s", fieldInfo.Name, c.componentID)
		}
		manifests = append(manifests, othermanifest)
	}
	return manifests, nil
}

func (c *TemplateContext) SetContextValue(manifests []*unstructured.Unstructured) {
	for i := range manifests {
		manifests[i].SetNamespace(c.namespace)
		manifests[i].SetLabels(c.as.GetCommonLabels(manifests[i].GetLabels()))
	}
}
func (c *TemplateContext) ExtendedContextFile() string {
	var buff string
	buff += fmt.Sprintf(ContextName+": \"%s\"\n", c.componentName)
	buff += fmt.Sprintf(ContextAppName+": \"%s\"\n", c.appName)
	buff += fmt.Sprintf(ContextNamespace+": \"%s\"\n", c.namespace)
	buff += fmt.Sprintf(ContextAppID+": \"%s\"\n", c.appID)
	buff += fmt.Sprintf(ContextID+": \"%s\"\n", c.componentID)
	return fmt.Sprintf("context: %s", structMarshal(buff))
}

func structMarshal(v string) string {
	skip := false
	v = strings.TrimFunc(v, func(r rune) bool {
		if !skip {
			if unicode.IsSpace(r) {
				return true
			}
			skip = true

		}
		return false
	})

	if strings.HasPrefix(v, "{") {
		return v
	}
	return fmt.Sprintf("{%s}", v)
}
