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

//AddVolumeStruct
//swagger:parameters addVolumes
type AddVolumeStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// Types of "application;app_publish"
		// in: body
		// required: true
		Category string `json:"category"`
		// container mount directory
		// in: body
		// required: true
		VolumePath string `json:"volume_path" validate:"volume_path|required|regex:^/"`
		// storage type (share,local,tmpfs)
		// in: body
		// required: true
		VolumeType string `json:"volume_type" validate:"volume_type|required"`
		// storage name (unique to the same application)
		// in: body
		// required: true
		VolumeName  string `json:"volume_name" validate:"volume_name|required|max:50"`
		FileContent string `json:"file_content"`
		// Storage driver alias (storageClass alias)
		VolumeProviderName string `json:"volume_provider_name"`
		IsReadOnly         bool   `json:"is_read_only"`
		// VolumeCapacity - storage size
		VolumeCapacity int64 `json:"volume_capacity"` // Unit: Mi
		// AccessMode read and write mode (important! A volume can only be mounted using one access mode at a time, even if it supports many. For example, a GCEPersistentDisk can be mounted as ReadWriteOnce by a single node or ReadOnlyMany by many nodes, but not at the same time. #https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes）
		AccessMode string `json:"access_mode"`
		// SharePolicy - sharing mode
		SharePolicy string `json:"share_policy"`
		// BackupPolicy - backup strategy
		BackupPolicy string `json:"backup_policy"`
		// ReclaimPolicy - recycling strategy
		ReclaimPolicy string `json:"reclaim_policy"`
		// AllowExpansion - whether to support expansion
		AllowExpansion bool `json:"allow_expansion"`
	}
}

//DeleteVolumeStruct
//swagger:parameters deleteVolumes
type DeleteVolumeStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// storage name
	// in: path
	// required: true
	VolumeName string `json:"volume_name"`
}

//AddVolumeDependencyStruct
//swagger:parameters addDepVolume
type AddVolumeDependencyStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// dependent service id
		// in: body
		// required: true
		DependServiceID string `json:"depend_service_id"  validate:"depend_service_id|required"`
		// container mount directory
		// in: body
		// required: true
		VolumePath string `json:"volume_path" validate:"volume_path|required|regex:^/"`
		// dependent storage name
		// in: body
		// required: true
		VolumeName string `json:"volume_name" validate:"volume_name|required|max:50"`

		VolumeType string `json:"volume_type" validate:"volume_type|required|in:share-file,config-file"`
	}
}

//DeleteVolumeDependencyStruct
//swagger:parameters  delDepVolume
type DeleteVolumeDependencyStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// dependent service id
		// in: body
		// required: true
		DependServiceID string `json:"depend_service_id" validate:"depend_service_id|required|max:32"`
		// dependent storage name
		// in: body
		// required: true
		VolumeName string `json:"volume_name" validate:"volume_name|required|max:50"`
	}
}

//The following is the v2 old version API parameter definition

//V2AddVolumeStruct AddVolumeStruct
//swagger:parameters addVolume
type V2AddVolumeStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// Types of "application;app_publish"
		// in: body
		// required: true
		Category string `json:"category"`
		// container mount directory
		// in: body
		// required: true
		VolumePath string `json:"volume_path" validate:"volume_path|required|regex:^/"`
		// host mount directory
		// in: body
		// required: true
		HostPath string `json:"host_path" validate:"volume_path|required|regex:^/"`
		// storage driver name
		VolumeProviderName string `json:"volume_provider_name"`
		// storage size
		VolumeCapacity int64 `json:"volume_capacity" validate:"volume_capacity|required|min:1"` // Unit Mi
		// AccessMode read and write mode（Important! A volume can only be mounted using one access mode at a time, even if it supports many. For example, a GCEPersistentDisk can be mounted as ReadWriteOnce by a single node or ReadOnlyMany by many nodes, but not at the same time. #https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes）)
		AccessMode string `gorm:"column:access_mode" json:"access_mode"`
		// SharePolicy - sharing mode
		SharePolicy string `gorm:"column:share_policy" json:"share_policy"`
		// BackupPolicy - backup strategy
		BackupPolicy string `gorm:"column:backup_policy" json:"backup_policy"`
		// ReclaimPolicy - recycling strategy
		ReclaimPolicy string `json:"reclaim_policy"`
		// AllowExpansion - whether to support expansion
		AllowExpansion bool `gorm:"column:allow_expansion" json:"allow_expansion"`
	}
}

//V2DelVolumeStruct AddVolumeStruct
//swagger:parameters deleteVolume
type V2DelVolumeStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// Types of "application;app_publish"
		// in: body
		// required: true
		Category string `json:"category"`
		// container mount directory
		// in: body
		// required: true
		VolumePath string `json:"volume_path" validate:"volume_path|required|regex:^/"`
	}
}

//V2AddVolumeDependencyStruct AddVolumeDependencyStruct
//swagger:parameters addVolumeDependency
type V2AddVolumeDependencyStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// dependent service id
		// in: body
		// required: true
		DependServiceID string `json:"depend_service_id"  validate:"depend_service_id|required"`
		// mount directory
		// in: body
		// required: true
		MntDir string `json:"mnt_dir" validate:"mnt_dir|required"`
		// the name of the directory in the mount container
		// in: body
		// required: true
		MntName string `json:"mnt_name" validate:"mnt_name|required"`
	}
}

//V2DelVolumeDependencyStruct
//swagger:parameters deleteVolumeDependency
type V2DelVolumeDependencyStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// dependent service id
		// in: body
		// required: true
		DependServiceID string `json:"depend_service_id"  validate:"depend_service_id|required"`
	}
}

// UpdVolumeReq is a value struct holding request for updating volume.
type UpdVolumeReq struct {
	VolumeName  string `json:"volume_name" validate:"required"`
	VolumeType  string `json:"volume_type" validate:"volume_type|required"`
	FileContent string `json:"file_content"`
	VolumePath  string `json:"volume_path" validate:"volume_path|required"`
}

// VolumeWithStatusResp volume status
type VolumeWithStatusResp struct {
	ServiceID string `json:"service_id"`
	//storage name
	Status map[string]string `json:"status"`
}

// VolumeWithStatusStruct volume with status struct
type VolumeWithStatusStruct struct {
	ServiceID string `json:"service_id"`
	//service type
	Category string `json:"category"`
	//storage type（share,local,tmpfs）
	VolumeType string `json:"volume_type"`
	//storage name
	VolumeName string `json:"volume_name"`
	//host address
	HostPath string `json:"host_path"`
	//mount address
	VolumePath string `json:"volume_path"`
	//read-only
	IsReadOnly bool `json:"is_read_only"`
	// VolumeCapacity -storage size
	VolumeCapacity int64 `json:"volume_capacity"`
	// AccessMode Read and write mode（Important! A volume can only be mounted using one access mode at a time, even if it supports many. For example, a GCEPersistentDisk can be mounted as ReadWriteOnce by a single node or ReadOnlyMany by many nodes, but not at the same time. #https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes）
	AccessMode string `json:"access_mode"`
	// SharePolicy - sharing mode
	SharePolicy string `json:"share_policy"`
	// BackupPolicy - backup strategy
	BackupPolicy string `json:"backup_policy"`
	// ReclaimPolicy - recycling strategy
	ReclaimPolicy string `json:"reclaim_policy"`
	// AllowExpansion - whether to support expansion
	AllowExpansion bool `json:"allow_expansion"`
	// VolumeProviderName - storage driver alias used
	VolumeProviderName string `json:"volume_provider_name"`
	Status             string `json:"status"`
}
