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

package model

import (
	dbmodel "github.com/gridworkz/kato/db/model"
)

//TenantResList TenantResList
type TenantResList []*TenantResource

//PagedTenantResList PagedTenantResList
type PagedTenantResList struct {
	List   []*TenantResource `json:"list"`
	Length int               `json:"length"`
}

//TenantResource abandoned
type TenantResource struct {
	//without plugin
	AllocatedCPU int `json:"alloc_cpu"`
	//without plugin
	AllocatedMEM int `json:"alloc_memory"`
	//with plugin
	UsedCPU int `json:"used_cpu"`
	//with plugin
	UsedMEM  int     `json:"used_memory"`
	UsedDisk float64 `json:"used_disk"`
	Name     string  `json:"name"`
	UUID     string  `json:"uuid"`
	EID      string  `json:"eid"`
}

func (list TenantResList) Len() int {
	return len(list)
}

func (list TenantResList) Less(i, j int) bool {
	if list[i].UsedMEM > list[j].UsedMEM {
		return true
	} else if list[i].UsedMEM < list[j].UsedMEM {
		return false
	} else {
		return list[i].UsedCPU > list[j].UsedCPU
	}
}

func (list TenantResList) Swap(i, j int) {
	temp := list[i]
	list[i] = list[j]
	list[j] = temp
}

//TenantAndResource tenant and resource strcut
type TenantAndResource struct {
	dbmodel.Tenants
	CPURequest            int64 `json:"cpu_request"`
	CPULimit              int64 `json:"cpu_limit"`
	MemoryRequest         int64 `json:"memory_request"`
	MemoryLimit           int64 `json:"memory_limit"`
	RunningAppNum         int64 `json:"running_app_num"`
	RunningAppInternalNum int64 `json:"running_app_internal_num"`
	RunningAppThirdNum    int64 `json:"running_app_third_num"`
}

//TenantList Tenant list struct
type TenantList []*TenantAndResource

//Add add
func (list *TenantList) Add(tr *TenantAndResource) {
	*list = append(*list, tr)
}
func (list TenantList) Len() int {
	return len(list)
}

func (list TenantList) Less(i, j int) bool {
	// Highest priority
	if list[i].MemoryRequest > list[j].MemoryRequest {
		return true
	}
	if list[i].MemoryRequest == list[j].MemoryRequest {
		if list[i].CPURequest > list[j].CPURequest {
			return true
		}
		if list[i].CPURequest == list[j].CPURequest {
			if list[i].RunningAppNum > list[j].RunningAppNum {
				return true
			}
			if list[i].RunningAppNum == list[j].RunningAppNum {
				// Minimum priority
				if list[i].Tenants.LimitMemory > list[j].Tenants.LimitMemory {
					return true
				}
			}
		}
	}
	return false
}

func (list TenantList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

//Paging
func (list TenantList) Paging(page, pageSize int) map[string]interface{} {
	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize
	var relist TenantList
	if startIndex < list.Len() && endIndex < list.Len() {
		relist = list[startIndex:endIndex]
	}
	if startIndex < list.Len() && endIndex >= list.Len() {
		relist = list[startIndex:]
	}
	return map[string]interface{}{
		"list":     relist,
		"page":     page,
		"pageSize": pageSize,
		"total":    list.Len(),
	}
}
