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

import "github.com/gridworkz/kato/pkg/apis/kato/v1alpha1"

//ThirdComponentProperties third component properties
type ThirdComponentProperties struct {
	Kubernetes *ThirdComponentKubernetes          `json:"kubernetes,omitempty"`
	Endpoints  []*v1alpha1.ThirdComponentEndpoint `json:"endpoints,omitempty"`
	Port       []*ThirdComponentPort              `json:"port"`
	Probe      *v1alpha1.Probe                    `json:"probe,omitempty"`
}

// ThirdComponentPort -
type ThirdComponentPort struct {
	Name      string `json:"name"`
	Port      int    `json:"port"`
	OpenInner bool   `json:"openInner"`
	OpenOuter bool   `json:"openOuter"`
}

// ThirdComponentKubernetes -
type ThirdComponentKubernetes struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
