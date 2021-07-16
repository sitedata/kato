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

// TenantServiceVolumeType
type TenantServiceVolumeType struct {
	Model
	VolumeType         string `gorm:"column:volume_type; size:64" json:"volume_type"`
	NameShow           string `gorm:"column:name_show; size:64" json:"name_show"`
	CapacityValidation string `gorm:"column:capacity_validation; size:1024" json:"capacity_validation"`
	Description        string `gorm:"column:description; size:1024" json:"description"`
	AccessMode         string `gorm:"column:access_mode; size:128" json:"access_mode"`
	BackupPolicy       string `gorm:"column:backup_policy; size:128" json:"backup_policy"`
	ReclaimPolicy      string `gorm:"column:reclaim_policy; size:20" json:"reclaim_policy"`
	SharePolicy        string `gorm:"share_policy; size:128" json:"share_policy"`
	Provisioner        string `gorm:"provisioner; size:128" json:"provisioner"`
	StorageClassDetail string `gorm:"storage_class_detail; size:2048" json:"storage_class_detail"`
	Sort               int    `gorm:"sort; default:9999" json:"sort"`
	Enable             bool   `gorm:"enable" json:"enable"`
}

// TableName
func (t *TenantServiceVolumeType) TableName() string {
	return "tenant_services_volume_type"
}
