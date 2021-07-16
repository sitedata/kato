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

package l4

import (
	"fmt"
	"github.com/gridworkz/kato/gateway/annotations/parser"
	"github.com/gridworkz/kato/gateway/annotations/resolver"
	extensions "k8s.io/api/extensions/v1beta1"
)

type Config struct {
	L4Enable bool
	L4Host   string
	L4Port   int
}

type l4 struct {
	r resolver.Resolver
}

func NewParser(r resolver.Resolver) parser.IngressAnnotation {
	return l4{r}
}

func (l l4) Parse(ing *extensions.Ingress) (interface{}, error) {
	l4Enable, _ := parser.GetBoolAnnotation("l4-enable", ing)
	l4Host, _ := parser.GetStringAnnotation("l4-host", ing)
	if l4Host == "" {
		l4Host = "0.0.0.0"
	}

	l4Port, _ := parser.GetIntAnnotation("l4-port", ing)
	if l4Enable && (l4Port <= 0 || l4Port > 65535) {
		return nil, fmt.Errorf("error l4Port: %d", l4Port)
	}
	return &Config{
		L4Enable: l4Enable,
		L4Host:   l4Host,
		L4Port:   l4Port,
	}, nil
}
