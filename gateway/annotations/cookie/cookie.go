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

package cookie

import (
	"github.com/gridworkz/kato/gateway/annotations/parser"
	"github.com/gridworkz/kato/gateway/annotations/resolver"
	extensions "k8s.io/api/extensions/v1beta1"
	"strings"
)

type Config struct {
	Cookie map[string]string `json:"cookie"`
}

type cookie struct {
	r resolver.Resolver
}

func NewParser(r resolver.Resolver) parser.IngressAnnotation {
	return cookie{r}
}

func (c cookie) Parse(ing *extensions.Ingress) (interface{}, error) {
	co, err := parser.GetStringAnnotation("cookie", ing)
	if err != nil {
		return nil, err
	}
	hmap := transform(co)

	return &Config{
		Cookie: hmap,
	}, nil
}

// transform transfers string to map
func transform(target string) map[string]string {
	target = strings.Replace(target, " ", "", -1)
	result := make(map[string]string)
	for _, item := range strings.Split(target, ";") {
		split := strings.Split(item, "=")
		if len(split) < 2 || split[0] == "" || split[1] == "" {
			continue
		}
		result[split[0]] = split[1]
	}
	if len(result) == 0 {
		return nil
	}

	return result
}
