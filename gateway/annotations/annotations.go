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

package annotations

import (
	"github.com/gridworkz/kato/gateway/annotations/cookie"
	"github.com/gridworkz/kato/gateway/annotations/header"
	"github.com/gridworkz/kato/gateway/annotations/l4"
	"github.com/gridworkz/kato/gateway/annotations/lbtype"
	"github.com/gridworkz/kato/gateway/annotations/parser"
	"github.com/gridworkz/kato/gateway/annotations/proxy"
	"github.com/gridworkz/kato/gateway/annotations/resolver"
	"github.com/gridworkz/kato/gateway/annotations/rewrite"
	"github.com/gridworkz/kato/gateway/annotations/upstreamhashby"
	weight "github.com/gridworkz/kato/gateway/annotations/wight"
	"github.com/gridworkz/kato/util/ingress-nginx/ingress/errors"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeniedKeyName name of the key that contains the reason to deny a location
const DeniedKeyName = "Denied"

// Ingress defines the valid annotations present in one NGINX Ingress rule
type Ingress struct {
	metav1.ObjectMeta
	Header            header.Config
	Cookie            cookie.Config
	Weight            weight.Config
	Rewrite           rewrite.Config
	L4                l4.Config
	UpstreamHashBy    string
	LoadBalancingType string
	Proxy             proxy.Config
}

// Extractor defines the annotation parsers to be used in the extraction of annotations
type Extractor struct {
	annotations map[string]parser.IngressAnnotation
}

// NewAnnotationExtractor creates a new annotations extractor
func NewAnnotationExtractor(cfg resolver.Resolver) Extractor {
	return Extractor{
		map[string]parser.IngressAnnotation{
			"Header":            header.NewParser(cfg),
			"Cookie":            cookie.NewParser(cfg),
			"Weight":            weight.NewParser(cfg),
			"Rewrite":           rewrite.NewParser(cfg),
			"L4":                l4.NewParser(cfg),
			"UpstreamHashBy":    upstreamhashby.NewParser(cfg),
			"LoadBalancingType": lbtype.NewParser(cfg),
			"Proxy":             proxy.NewParser(cfg),
		},
	}
}

// Extract extracts the annotations from an Ingress
func (e Extractor) Extract(ing *extensions.Ingress) *Ingress {
	pia := &Ingress{
		ObjectMeta: ing.ObjectMeta,
	}

	data := make(map[string]interface{})
	for name, annotationParser := range e.annotations {
		val, err := annotationParser.Parse(ing)
		if err != nil {
			if errors.IsMissingAnnotations(err) {
				continue
			}

			if !errors.IsLocationDenied(err) {
				continue
			}

			_, alreadyDenied := data[DeniedKeyName]
			if !alreadyDenied {
				data[DeniedKeyName] = err
				logrus.Errorf("error reading %v annotation in Ingress %v/%v: %v", name, ing.GetNamespace(), ing.GetName(), err)
				continue
			}

			logrus.Infof("error reading %v annotation in Ingress %v/%v: %v", name, ing.GetNamespace(), ing.GetName(), err)
		}

		if val != nil {
			data[name] = val
		}
	}

	err := mergo.MapWithOverwrite(pia, data)
	if err != nil {
		logrus.Errorf("unexpected error merging extracted annotations: %v", err)
	}

	return pia
}
