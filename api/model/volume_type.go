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

// VolumeTypeStruct volume option struct
type VolumeTypeStruct struct {
	VolumeType         string                 `json:"volume_type" validate:"volume_type|required"`
	NameShow           string                 `json:"name_show"`
	CapacityValidation map[string]interface{} `json:"capacity_validation"`
	Description        string                 `json:"description"`
	AccessMode         []string               `json:"access_mode"`    // read and write mode（Important! A volume can only be mounted using one access mode at a time, even if it supports many. For example, a GCEPersistentDisk can be mounted as ReadWriteOnce by a single node or ReadOnlyMany by many nodes, but not at the same time. #https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes）
	SharePolicy        []string               `json:"share_policy"`   // sharing mode
	BackupPolicy       []string               `json:"backup_policy"`  // backup strategy
	ReclaimPolicy      string                 `json:"reclaim_policy"` // recycling strategy: delete, retain, recyle
	Provisioner        string                 `json:"provisioner"`    // storage provider
	StorageClassDetail map[string]interface{} `json:"storage_class_detail" validate:"storage_class_detail|required"`
	Sort               int                    `json:"sort"`   // sort
	Enable             bool                   `json:"enable"` // does it take effect
}

// VolumeTypePageStruct volume option struct with page
type VolumeTypePageStruct struct {
	list     *VolumeTypeStruct
	page     int
	pageSize int
	count    int
}
