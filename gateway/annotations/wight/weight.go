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

package weight

import (
	"github.com/gridworkz/kato/gateway/annotations/parser"
	"github.com/gridworkz/kato/gateway/annotations/resolver"
	"github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	"strconv"
)

// Config contains weight or router
type Config struct {
	Weight int
}

type weight struct {
	r resolver.Resolver
}

// NewParser creates a new parser
func NewParser(r resolver.Resolver) parser.IngressAnnotation {
	return weight{r}
}

func (c weight) Parse(ing *networkingv1.Ingress) (interface{}, error) {
	wstr, err := parser.GetStringAnnotation("weight", ing)
	var w int
	if err != nil || wstr == "" {
		w = 1
	} else {
		w, err = strconv.Atoi(wstr)
		if err != nil {
			logrus.Warnf("Unexpected error occurred when convert string(%s) to int: %v", wstr, err)
			w = 1
		}
	}
	return &Config{
		Weight: w,
	}, nil
}
