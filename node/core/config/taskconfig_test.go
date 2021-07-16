// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

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

package config

import "testing"
import "fmt"

func TestResettingArray(t *testing.T) {
	c := CreateDataCenterConfig()
	c.Start()
	defer c.Stop()
	groupCtx := NewGroupContext("")
	groupCtx.Add("SADAS", "Test")
	result, err := ResettingArray(groupCtx, []string{"Sdd${sadas}asd", "${MYSQL_HOST}", "12_${MYSQL_PASS}_sd"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestResettingString(t *testing.T) {
	c := CreateDataCenterConfig()
	c.Start()
	defer c.Stop()
	groupCtx := NewGroupContext("")
	groupCtx.Add("SADAS", "Test")
	result, err := ResettingString(nil, "${MYSQL_HOST}Sdd${sadas}asd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestGroupConfig(t *testing.T) {
	groupCtx := NewGroupContext("")
	v := groupCtx.Get("API")
	fmt.Println("asdasd:", v)
}
