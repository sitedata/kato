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

package conversion

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

//Allocate the CPU at the ratio of 4g memory to 1 core CPU
func createResourcesByDefaultCPU(memory int, setCPURequest, setCPULimit int64) corev1.ResourceRequirements {
	var cpuRequest, cpuLimit int64
	base := int64(memory) / 128
	if base <= 0 {
		base = 1
	}
	if memory < 512 {
		cpuRequest, cpuLimit = base*30, base*80
	} else if memory <= 1024 {
		cpuRequest, cpuLimit = base*30, base*160
	} else {
		cpuRequest, cpuLimit = base*30, ((int64(memory)-1024)/1024*500 + 1280)
	}
	if setCPULimit > 0 {
		cpuLimit = setCPULimit
	}
	if setCPURequest > 0 {
		cpuRequest = setCPURequest
	}

	limits := corev1.ResourceList{}
	limits[corev1.ResourceCPU] = *resource.NewMilliQuantity(cpuLimit, resource.DecimalSI)
	limits[corev1.ResourceMemory] = *resource.NewQuantity(int64(memory*1024*1024), resource.BinarySI)

	request := corev1.ResourceList{}
	request[corev1.ResourceCPU] = *resource.NewMilliQuantity(cpuRequest, resource.DecimalSI)
	request[corev1.ResourceMemory] = *resource.NewQuantity(int64(memory*1024*1024), resource.BinarySI)

	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: request,
	}
}
