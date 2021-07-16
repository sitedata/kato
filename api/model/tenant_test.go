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

package model

import (
	"sort"
	"testing"
)

func TestTenantList(t *testing.T) {
	var tenants TenantList
	t1 := &TenantAndResource{
		MemoryRequest: 100,
	}
	t1.LimitMemory = 30
	tenants.Add(t1)

	t2 := &TenantAndResource{
		MemoryRequest: 80,
	}
	t2.LimitMemory = 40
	tenants.Add(t2)

	t3 := &TenantAndResource{
		MemoryRequest: 0,
	}
	t3.LimitMemory = 60
	t4 := &TenantAndResource{
		MemoryRequest: 0,
	}
	t4.LimitMemory = 70

	t5 := &TenantAndResource{
		RunningAppNum: 10,
	}
	t5.LimitMemory = 0

	tenants.Add(t3)
	tenants.Add(t4)
	tenants.Add(t5)
	sort.Sort(tenants)
	for _, ten := range tenants {
		t.Logf("%+v", ten)
	}
}
