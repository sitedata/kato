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

package multi

import (
	"os"
	"testing"
)

func TestMaven_ParsePom(t *testing.T) {
	pomPath := "./pom.xml"
	pom, err := parsePom(pomPath)
	if err != nil {
		t.Fatal(err)
	}
	if pom.Packaging != "pom" {
		t.Errorf("Expected pom for pom.Packaging, but returned %s", pom.Packaging)
	}
	if pom.Modules == nil || len(pom.Modules) != 3 {
		t.Error("Modules not found")
	} else {
		if pom.Modules[0] != "rbd-api" {
			t.Errorf("Expected 'rbd-api' for pom.Modules[0], but returned %s", pom.Modules[0])
		}
		if pom.Modules[1] != "rbd-worker" {
			t.Errorf("Expected 'rbd-worker' for pom.Modules[0], but returned %s", pom.Modules[0])
		}
		if pom.Modules[2] != "rbd-gateway" {
			t.Errorf("Expected 'rbd-gateway' for pom.Modules[0], but returned %s", pom.Modules[0])
		}
	}
}

func TestMaven_ListModules(t *testing.T) {
	path := os.Getenv("GOPATH") + "/src/github.com/gridworkz/kato/builder/parser/code/multisvc/"
	m := maven{}
	res, err := m.ListModules(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 3 {
		t.Errorf("Expected 3 for the length of mudules, but returned %d", len(res))
	}
	for _, svc := range res {
		for _, env := range svc.Envs {
			t.Logf("Name: %s; Value: %s", env.Name, env.Value)
		}
	}
	t.Error("test")
}
