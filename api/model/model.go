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
	"net/url"
	"time"

	"github.com/gridworkz/kato/util"

	dbmodel "github.com/gridworkz/kato/db/model"
	dmodel "github.com/gridworkz/kato/worker/discover/model"
)

// AppType
const (
	AppTypeKato = "kato"
	AppTypeHelm     = "helm"
)

//ServiceGetCommon path parameter
//swagger:parameters getVolumes getDepVolumes
type ServiceGetCommon struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
}

//ComposerStruct ComposerStruct
// swagger:parameters resolve
type ComposerStruct struct {
	// in : body
	Body struct {
		Lang string `json:"default_runtime" validate:"default_runtime"`
		Data struct {
			JSON struct {
				PlatForm struct {
					PHP string `json:"php" validate:"php"`
				}
			}
			Packages []string `json:"packages" validate:"packages"`
			Lock     struct {
				PlatForm struct {
					PHP string `json:"php" validate:"php"`
				}
			}
		}
	}
}

//CreateServiceStruct serviceCreate struct
// swagger:parameters createService
type CreateServiceStruct struct {
	// in: path
	// required: true
	TenantName string `gorm:"column:tenant_name;size:32" json:"tenant_name" validate:"tenant_name"`
	// in:body
	Body struct {
		// tenant id
		// in: body
		// required: false
		TenantID string `gorm:"column:tenant_id;size:32" json:"tenant_id" validate:"tenant_id"`
		// application id
		// in: body
		// required: false
		ServiceID string `gorm:"column:service_id;size:32" json:"service_id" validate:"service_id"`
		// Operator
		// in: body
		// required: false
		Operator string `json:"operator" validate:"operator"`
		// Application label, value
		// in: body
		// required: false
		ServiceLabel string `json:"service_label" validate:"service_label"`
		// node label, format: v1, v2
		// in: body
		// required: false
		NodeLabel string `json:"node_label" validate:"node_label"`
		// Depend on id, format: []struct TenantServiceRelation
		// in: body
		// required: false
		DependIDs []dbmodel.TenantServiceRelation `json:"depend_ids" validate:"depend_ids"`
		// Persistent directory information, format: []struct TenantServiceVolume
		// in: body
		// required: false
		VolumesInfo []dbmodel.TenantServiceVolume `json:"volumes_info" validate:"volumes_info"`
		// Environment variable information, format: []struct TenantServiceEnvVar
		// in: body
		// required: false
		EnvsInfo []dbmodel.TenantServiceEnvVar `json:"envs_info" validate:"envs_info"`
		// Port information, format: []struct TenantServicesPort
		// in: body
		// required: false
		PortsInfo []dbmodel.TenantServicesPort `json:"ports_info" validate:"ports_info"`
		// service key
		// in: body
		// required: false
		ServiceKey string `gorm:"column:service_key;size:32" json:"service_key" validate:"service_key"`
		// service alias
		// in: body
		// required: true
		ServiceAlias string `gorm:"column:service_alias;size:30" json:"service_alias" validate:"service_alias"`
		// Service description
		// in: body
		// required: false
		Comment string `gorm:"column:comment" json:"comment" validate:"comment"`
		// Service version
		// in: body
		// required: false
		ServiceVersion string `gorm:"column:service_version;size:32" json:"service_version" validate:"service_version"`
		// Mirror name
		// in: body
		// required: false
		ImageName string `gorm:"column:image_name;size:100" json:"image_name" validate:"image_name"`
		// Container CPU weight
		// in: body
		// required: false
		ContainerCPU int `gorm:"column:container_cpu;default:500" json:"container_cpu" validate:"container_cpu"`
		// Maximum memory of the container
		// in: body
		// required: false
		ContainerMemory int `gorm:"column:container_memory;default:128" json:"container_memory" validate:"container_memory"`
		// Container start command
		// in: body
		// required: false
		ContainerCMD string `gorm:"column:container_cmd;size:2048" json:"container_cmd" validate:"container_cmd"`
		// Container environment variables
		// in: body
		// required: false
		ContainerEnv string `gorm:"column:container_env;size:255" json:"container_env" validate:"container_env"`
		// Volume name
		// in: body
		// required: false
		VolumePath string `gorm:"column:volume_path" json:"volume_path" validate:"volume_path"`
		// Container mount directory
		// in: body
		// required: false
		VolumeMountPath string `gorm:"column:volume_mount_path" json:"volume_mount_path" validate:"volume_mount_path"`
		// host directory
		// in: body
		// required: false
		HostPath string `gorm:"column:host_path" json:"host_path" validate:"host_path"`
		// Expansion method; 0: stateless; 1: stateful; 2: partition
		// in: body
		// required: false
		ExtendMethod string `gorm:"column:extend_method;default:'stateless';" json:"extend_method" validate:"extend_method"`
		// number of nodes
		// in: body
		// required: false
		Replicas int `gorm:"column:replicas;default:1" json:"replicas" validate:"replicas"`
		// Deployment version
		// in: body
		// required: false
		DeployVersion string `gorm:"column:deploy_version" json:"deploy_version" validate:"deploy_version"`
		// Service category: application, cache, store
		// in: body
		// required: false
		Category string `gorm:"column:category" json:"category" validate:"category"`
		// Latest operation ID
		// in: body
		// required: false
		EventID string `gorm:"column:event_id" json:"event_id" validate:"event_id"`
		// Service type
		// in: body
		// required: false
		ServiceType string `gorm:"column:service_type" json:"service_type" validate:"service_type"`
		// Mirror source
		// in: body
		// required: false
		Namespace string `gorm:"column:namespace" json:"namespace" validate:"namespace"`
		// Shared type shared, exclusive
		// in: body
		// required: false
		VolumeType string `gorm:"column:volume_type;default:'shared'" json:"volume_type" validate:"volume_type"`
		// Port type, one_outer; dif_protocol; multi_outer
		// in: body
		// required: false
		PortType string `gorm:"column:port_type;default:'multi_outer'" json:"port_type" validate:"port_type"`
		// update time
		// in: body
		// required: false
		UpdateTime time.Time `gorm:"column:update_time" json:"update_time" validate:"update_time"`
		// Service creation type cloud cloud city service, assistant cloud help service
		// in: body
		// required: false
		ServiceOrigin string `gorm:"column:service_origin;default:'assistant'" json:"service_origin" validate:"service_origin"`
		// Code source: gitlab,github
		// in: body
		// required: false
		CodeFrom string `gorm:"column:code_from" json:"code_from" validate:"code_from"`
	}
}

// UpdateServiceStruct service update
// swagger:parameters updateService
type UpdateServiceStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		// Container start command
		// in: body
		// required: false
		ContainerCMD string `gorm:"column:container_cmd;size:2048" json:"container_cmd" validate:"container_cmd"`
		// Mirror name
		// in: body
		// required: false
		ImageName string `gorm:"column:image_name;size:100" json:"image_name" validate:"image_name"`
		// Maximum memory of the container
		// in: body
		// required: false
		ContainerMemory int `gorm:"column:container_memory;default:128" json:"container_memory" validate:"container_memory"`
	}
}

//StartStopStruct start struct
type StartStopStruct struct {
	ServiceID     string
	TenantID      string
	DeployVersion string
	EventID       string
	TaskType      string
}

//LanguageSet set language
type LanguageSet struct {
	ServiceID string `json:"service_id"`
	Language  string `json:"language"`
}

//ServiceStruct service struct
type ServiceStruct struct {
	TenantID string `json:"tenant_id" validate:"tenant_id"`
	// in: path
	// required: true
	ServiceID string `json:"service_id" validate:"service_id"`
	// Service name, used for stateful service DNS
	// in: body
	// required: false
	ServiceName string `json:"service_name" validate:"service_name"`
	// service alias
	// in: body
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
	// Component type
	// in: body
	// required: true
	ServiceType string `json:"service_type" validate:"service_type"`
	// Service description
	// in: body
	// required: false
	Comment string `json:"comment" validate:"comment"`
	// Service version
	// in: body
	// required: false
	ServiceVersion string `json:"service_version" validate:"service_version"`
	// Mirror name
	// in: body
	// required: false
	ImageName string `json:"image_name" validate:"image_name"`
	// Container CPU weight
	// in: body
	// required: false
	ContainerCPU int `json:"container_cpu" validate:"container_cpu"`
	// Maximum memory of the container
	// in: body
	// required: false
	ContainerMemory int `json:"container_memory" validate:"container_memory"`
	// component gpu video memory
	// in: body
	// required: false
	ContainerGPU int `json:"container_gpu" validate:"container_gpu"`
	// Container start command
	// in: body
	// required: false
	ContainerCMD string `json:"container_cmd" validate:"container_cmd"`
	// Container environment variables
	// in: body
	// required: false
	ContainerEnv string `json:"container_env" validate:"container_env"`
	// Expansion method; 0: stateless; 1: stateful; 2: partition (v5.2 is used for the type of receiving component)
	// in: body
	// required: false
	ExtendMethod string `json:"extend_method" validate:"extend_method"`
	// number of nodes
	// in: body
	// required: false
	Replicas int `json:"replicas" validate:"replicas"`
	// Deployment version
	// in: body
	// required: false
	DeployVersion string `json:"deploy_version" validate:"deploy_version"`
	// Service category: application, cache, store
	// in: body
	// required: false
	Category string `json:"category" validate:"category"`
	// Current status of the service: undeploy,running,closed,unusual,starting,checking,stoping
	// in: body
	// required: false
	CurStatus string `json:"cur_status" validate:"cur_status"`
	// Latest operation ID
	// in: body
	// required: false
	EventID string `json:"event_id" validate:"event_id"`
	// Mirror source
	// in: body
	// required: false
	Namespace string `json:"namespace" validate:"namespace"`
	// update time
	// in: body
	// required: false
	UpdateTime time.Time `json:"update_time" validate:"update_time"`
	// Service creation type cloud cloud city service, assistant cloud help service
	// in: body
	// required: false
	ServiceOrigin string `json:"service_origin" validate:"service_origin"`
	Kind          string `json:"kind" validate:"kind|in:internal,third_party"`
	EtcdKey       string `json:"etcd_key" validate:"etcd_key"`
	//OSType runtime os type
	// in: body
	// required: false
	OSType            string                               `json:"os_type" validate:"os_type|in:windows,linux"`
	ServiceLabel      string                               `json:"service_label"  validate:"service_label|in:StatelessServiceType,StatefulServiceType"`
	NodeLabel         string                               `json:"node_label"  validate:"node_label"`
	Operator          string                               `json:"operator"  validate:"operator"`
	RepoURL           string                               `json:"repo_url" validate:"repo_url"`
	DependIDs         []dbmodel.TenantServiceRelation      `json:"depend_ids" validate:"depend_ids"`
	VolumesInfo       []TenantServiceVolumeStruct          `json:"volumes_info" validate:"volumes_info"`
	DepVolumesInfo    []dbmodel.TenantServiceMountRelation `json:"dep_volumes_info" validate:"dep_volumes_info"`
	EnvsInfo          []dbmodel.TenantServiceEnvVar        `json:"envs_info" validate:"envs_info"`
	PortsInfo         []dbmodel.TenantServicesPort         `json:"ports_info" validate:"ports_info"`
	Endpoints         *Endpoints                           `json:"endpoints" validate:"endpoints"`
	AppID             string                               `json:"app_id" validate:"required"`
	ComponentProbes   []ServiceProbe                       `json:"component_probes" validate:"component_probes"`
	ComponentMonitors []AddServiceMonitorRequestStruct     `json:"component_monitors" validate:"component_monitors"`
	HTTPRules         []AddHTTPRuleStruct                  `json:"http_rules" validate:"http_rules"`
	TCPRules          []AddTCPRuleStruct                   `json:"tcp_rules" validate:"tcp_rules"`
}

// Endpoints holds third-party service endpoints or configuraion to get endpoints.
type Endpoints struct {
	Static     []string            `json:"static" validate:"static"`
	Kubernetes * EndpointKubernetes `json:" kubernetes "validate:" kubernetes "`
}

// DbModel -
func (e *Endpoints) DbModel(componentID string) *dbmodel.ThirdPartySvcDiscoveryCfg {
	return &dbmodel.ThirdPartySvcDiscoveryCfg{
		ServiceID:   componentID,
		Type:        string(dbmodel.DiscorveryTypeKubernetes),
		Namespace:   e.Kubernetes.Namespace,
		ServiceName: e.Kubernetes.ServiceName,
	}
}

// EndpointKubernetes -
type EndpointKubernetes struct {
	Namespace   string `json:"namespace"`
	ServiceName string `json:"serviceName"`
}

//TenantServiceVolumeStruct -
type TenantServiceVolumeStruct struct {
	ServiceID string ` json:"service_id"`
	//Service type
	Category string `json:"category"`
	//Storage type (share, local, tmpfs)
	VolumeType string `json:"volume_type"`
	//Storage name
	VolumeName string `json:"volume_name"`
	//Host address
	HostPath string `json:"host_path"`
	//Mount address
	VolumePath string `json:"volume_path"`
	//Whether it is read-only
	IsReadOnly bool `json:"is_read_only"`

	FileContent string `json:"file_content"`
	// VolumeCapacity storage size
	VolumeCapacity int64 `json:"volume_capacity"`
	// AccessMode 读写模式（Important! A volume can only be mounted using one access mode at a time, even if it supports many. For example, a GCEPersistentDisk can be mounted as ReadWriteOnce by a single node or ReadOnlyMany by many nodes, but not at the same time. #https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes）
	AccessMode string `json:"access_mode"`
	// SharePolicy sharing mode
	SharePolicy string `json:"share_policy"`
	// BackupPolicy backup policy
	BackupPolicy string `json:"backup_policy"`
	// ReclaimPolicy recycling strategy
	ReclaimPolicy string `json:"reclaim_policy"`
	// Whether AllowExpansion supports expansion
	AllowExpansion bool `json:"allow_expansion"`
	// The storage driver alias used by VolumeProviderName
	VolumeProviderName string `json:"volume_provider_name"`
}

//DependService struct for depend service
type DependService struct {
	TenantID       string `json:"tenant_id"`
	ServiceID      string `json:"service_id"`
	DepServiceID   string `json:"dep_service_id"`
	DepServiceType string `json:"dep_service_type"`
	Action         string `json:"action"`
}

//Attr attr
type Attr struct {
	Action    string `json:"action"`
	TenantID  string `json:"tenant_id"`
	ServiceID string `json:"service_id"`
	AttrName  string `json:"env_name"`
	AttrValue string `json:"env_value"`
}

// DeleteServicePort service port
// swagger:parameters deletePort
type DeleteServicePort struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// container port
	// in: path
	// required: true
	Port int `json:"port"`
}

// TenantResources TenantResources
// swagger:parameters tenantResources
type TenantResources struct {
	// in: body
	Body struct {
		// in: body
		// required: true
		TenantNames []string `json:"tenant_name" validate:"tenant_name"`
	}
}

// ServicesResources ServicesResources
// swagger:parameters serviceResources
type ServicesResources struct {
	// in: body
	Body struct {
		// in: body
		// required: true
		ServiceIDs []string `json:"service_ids" validate:"service_ids"`
	}
}

// CommandResponse api unified return structure
// swagger:response commandResponse
type CommandResponse struct {
	// in: body
	Body struct {
		//Parameter verification error message
		ValidationError url.Values `json:"validation_error,omitempty"`
		//API error message
		Msg string `json:"msg,omitempty"`
		//Single resource entity
		Bean interface{} `json:"bean,omitempty"`
		//Resource list
		List interface{} `json:"list,omitempty"`
		//Total number of data sets
		ListAllNumber int `json:"number,omitempty"`
		//Current page number
		Page int `json:"page,omitempty"`
	}
}

// ServicePortInnerOrOuter service port
// swagger:parameters PortInnerController PortOuterController
type ServicePortInnerOrOuter struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: path
	// required: true
	Port int `json:"port"`
	//in: body
	Body struct {
		// Operation value `close` or `open`
		// in: body
		// required: true
		Operation      string `json:"operation"  validate:"operation|required|in:open,close"`
		IfCreateExPort bool   `json:"if_create_ex_port"`
	}
}

// ServiceLBPortChange change lb port
// swagger:parameters changelbport
type ServiceLBPortChange struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: path
	// required: true
	Port int `json:"port"`
	//in: body
	Body struct {
		// in: body
		// required: true
		ChangePort int `json:"change_port"  validate:"change_port|required"`
	}
}

//RollbackStruct struct
type RollbackStruct struct {
	TenantID      string `json:"tenant_id"`
	ServiceID     string `json:"service_id"`
	EventID       string `json:"event_id;default:system"`
	Operator      string `json:"operator"`
	DeployVersion string `json:"deploy_version"`
}

//StatusList status list
type StatusList struct {
	TenantID      string     `json:"tenant_id"`
	ServiceID     string     `json:"service_id"`
	ServiceAlias  string     `json:"service_alias"`
	DeployVersion string     `json:"deploy_version"`
	Replicas      int        `json:"replicas"`
	ContainerMem  int        `json:"container_memory"`
	CurStatus     string     `json:"cur_status"`
	ContainerCPU  int        `json:"container_cpu"`
	StatusCN      string     `json:"status_cn"`
	StartTime     string     `json:"start_time"`
	PodList       []PodsList `json:"pod_list"`
}

//PodsList pod list
type PodsList struct {
	PodIP    string `json:"pod_ip"`
	Phase    string `json:"phase"`
	PodName  string `json:"pod_name"`
	NodeName string `json:"node_name"`
}

//StatsInfo stats info
type StatsInfo struct {
	UUID string `json:"uuid"`
	CPU  int    `json:"cpu"`
	MEM  int    `json:"memory"`
}

//TotalStatsInfo total stats info
type TotalStatsInfo struct {
	Data []*StatsInfo `json:"data"`
}

//LicenseInfo license info
type LicenseInfo struct {
	Code       string   `json:"code"`
	Company    string   `json:"company"`
	Node       int      `json:"node"`
	CPU        int      `json:"cpu"`
	MEM        int      `json:"memory"`
	Tenant     int      `json:"tenant"`
	EndTime    string   `json:"end_time"`
	StartTime  string   `json:"start_time"`
	DataCenter int      `json:"data_center"`
	ModuleList []string `json:"module_list"`
}

// AddTenantStruct AddTenantStruct
// swagger:parameters addTenant
type AddTenantStruct struct {
	//in: body
	Body struct {
		// the tenant id
		// in: body
		// required: false
		TenantID string `json:"tenant_id" validate:"tenant_id"`
		// the tenant name
		// in: body
		// required: false
		TenantName string `json:"tenant_name" validate:"tenant_name"`
		// the eid
		// in : body
		// required: false
		Eid         string `json:"eid" validata:"eid"`
		Token       string `json:"token" validate:"token"`
		LimitMemory int    `json:"limit_memory" validate:"limit_memory"`
	}
}

// UpdateTenantStruct UpdateTenantStruct
// swagger:parameters updateTenant
type UpdateTenantStruct struct {
	//in: body
	Body struct {
		// the eid
		// in : body
		// required: false
		LimitMemory int `json:"limit_memory" validate:"limit_memory"`
	}
}

// ServicesInfoStruct ServicesInfoStruct
// swagger:parameters getServiceInfo
type ServicesInfoStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
}

// SetLanguageStruct SetLanguageStruct
// swagger:parameters setLanguage
type SetLanguageStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		// the tenant id
		// in: body
		// required: true
		EventID string `json:"event_id"`
		// the language
		// in: body
		// required: true
		Language string `json:"language"`
	}
}

//StartServiceStruct StartServiceStruct
//swagger:parameters startService stopService restartService
type StartServiceStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		// the tenant id
		// in: body
		// required: false
		EventID string `json:"event_id"`
	}
}

//VerticalServiceStruct VerticalServiceStruct
//swagger:parameters verticalService
type VerticalServiceStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		// the event id
		// in: body
		// required: false
		EventID string `json:"event_id"`
		// number of cpu
		// in: body
		// required: false
		ContainerCPU int `json:"container_cpu"`
		// memory size
		// in: body
		// required: false
		ContainerMemory int `json:"container_memory"`
	}
}

//HorizontalServiceStruct HorizontalServiceStruct
//swagger:parameters horizontalService
type HorizontalServiceStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		// the event id
		// in: body
		// required: false
		EventID string `json:"event_id"`
		// Number of scaling
		// in: body
		// required: false
		NodeNUM int `json:"node_num"`
	}
}

//BuildServiceStruct BuildServiceStruct
//swagger:parameters serviceBuild
type BuildServiceStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
	//in: body
	Body struct {
		// the event id
		// in: body
		// required: false
		EventID string `json:"event_id" validate:"event_id"`
		// variable
		// in: body
		// required: false
		ENVS map[string]string `json:"envs" validate:"envs"`
		// Application build type
		// in: body
		// required: true
		Kind string `json:"kind" validate:"kind|required"`
		// Follow-up action, one-click deployment based on the value, if the value is not passed, only the build is performed by default
		// in: body
		// required: false
		Action string `json:"action" validate:"action"`
		// Mirror address
		// in: body
		// required: false
		ImageURL string `json:"image_url" validate:"image_url"`
		// deployed version number
		// in: body
		// required: true
		DeployVersion string `json:"deploy_version" validate:"deploy_version|required"`
		// git address
		// in: body
		// required: false
		RepoURL string `json:"repo_url" validate:"repo_url"`
		// branch branch information
		// in: body
		// required: false
		Branch string `json:"branch" validate:"branch"`
		// operator
		// in: body
		// required: false
		Lang string `json:"lang" validate:"lang"`
		// Code server type
		// in: body
		// required: false
		ServerType   string `json:"server_type" validate:"server_type"`
		Runtime      string `json:"runtime" validate:"runtime"`
		ServiceType  string `json:"service_type" validate:"service_type"`
		User         string `json:"user" validate:"user"`
		Password     string `json:"password" validate:"password"`
		Operator     string `json:"operator" validate:"operator"`
		TenantName   string `json:"tenant_name"`
		ServiceAlias string `json:"service_alias"`
		Cmd          string `json:"cmd"`
		//Used for cloud city code package creation
		SlugInfo struct {
			SlugPath    string `json:"slug_path"`
			FTPHost     string `json:"ftp_host"`
			FTPPort     string `json:"ftp_port"`
			FTPUser     string `json:"ftp_username"`
			FTPPassword string `json:"ftp_password"`
		} `json:"slug_info"`
	}
}

//V1BuildServiceStruct V1BuildServiceStruct
type V1BuildServiceStruct struct {
	// in: path
	// required: true
	ServiceID string `json:"service_id" validate:"service_id"`
	Body      struct {
		ServiceID     string `json:"service_id" validate:"service_id"`
		EventID       string `json:"event_id" validate:"event_id"`
		ENVS          string `json:"envs" validate:"envs"`
		Kind          string `json:"kind" validate:"kind"`
		Action        string `json:"action" validate:"action"`
		ImageURL      string `json:"image_url" validate:"image_url"`
		DeployVersion string `json:"deploy_version" validate:"deploy_version|required"`
		RepoURL       string `json:"repo_url" validate:"repo_url"`
		GitURL        string `json:"gitUrl" validate:"gitUrl"`
		Operator      string `json:"operator" validate:"operator"`
	}
}

//UpgradeServiceStruct UpgradeServiceStruct
//swagger:parameters upgradeService
type UpgradeServiceStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		// the event id
		// in: body
		// required: false
		EventID string `json:"event_id"`
		// version number
		// in: body
		// required: true
		DeployVersion int `json:"deploy_version"`
		// operator
		// in: body
		// required: false
		Operator int `json:"operator"`
	}
}

//StatusServiceStruct StatusServiceStruct
//swagger:parameters serviceStatus
type StatusServiceStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
}

//StatusServiceListStruct StatusServiceListStruct
//swagger:parameters serviceStatuslist
type StatusServiceListStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: body
	// required: true
	Body struct {
		// The list of service IDs that need to get the status, if not specified, return the status of all applications of the tenant
		// in: body
		// required: true
		ServiceIDs []string `json:"service_ids" validate:"service_ids|required"`
	}
}

//AddServiceLabelStruct AddServiceLabelStruct
//swagger:parameters addServiceLabel updateServiceLabel
type AddServiceLabelStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// Tag value, the format is "v1"
		// in: bid
		// required: true
		LabelValues string `json:"label_values"`
	}
}

//AddNodeLabelStruct AddNodeLabelStruct
//swagger:parameters addNodeLabel deleteNodeLabel
type AddNodeLabelStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// Tag value, the format is "[v1, v2, v3]"
		// in: body
		// required: true
		LabelValues []string `json:"label_values" validate:"label_values|required"`
	}
}

// LabelsStruct blabla
type LabelsStruct struct {
	Labels []LabelStruct `json:"labels"`
}

// LabelStruct holds info for adding, updating or deleting label
type LabelStruct struct {
	LabelKey   string `json:"label_key" validate:"label_key|required"`
	LabelValue string `json:"label_value" validate:"label_value|required"`
}

// GetSingleServiceInfoStruct GetSingleServiceInfoStruct
//swagger:parameters getService deleteService
type GetSingleServiceInfoStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
}

//CheckCodeStruct CheckCodeStruct
//swagger:parameters checkCode
type CheckCodeStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: body
	Body struct {
		// git branch details
		// in: body
		// required: true
		GitURL string `json:"git_url" validate:"git_url|required"`
		// git address
		// in: body
		// required: true
		URLRepos string `json:"url_repos" validate:"url_repos|required"`
		// Check type, "first_check"
		// in: body
		// required: true
		CheckType string `json:"check_type" validate:"check_type|required"`
		// code branch
		// in: body
		// required: true
		CodeVersion string `json:"code_version" validate:"code_version|required"`
		// git project id, 0
		// in: body
		// required: true
		GitProjectID int `json:"git_project_id" validate:"git_project_id|required"`
		// git source, "gitlab_manual"
		// in: body
		// required: true
		CodeFrom string `json:"code_from" validate:"code_from|required"`
		// tenant id
		// in: body
		// required: false
		TenantID string `json:"tenant_id" validate:"tenant_id"`
		Action   string `json:"action"`
		// application id
		// in: body
		// required: true
		ServiceID string `json:"service_id"`
	}
}

//ServiceCheckStruct application detection, support source code detection, mirror detection, dockerrun detection
//swagger:parameters serviceCheck
type ServiceCheckStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: body
	Body struct {
		//uuid
		// in: body
		CheckUUID string `json:"uuid"`
		//Detection source type
		// in: body
		// required: true
		SourceType string `json:"source_type" validate:"source_type|required|in:docker-run,docker-compose,sourcecode,third-party-service"`

		CheckOS string `json:"check_os"`
		// Definition of detection source,
		// Code: https://github.com/gridworkz/kato.git master
		// docker-run: docker run --name xxx nginx:latest nginx
		// docker-compose: compose full text
		// in: body
		// required: true
		SourceBody string `json:"source_body"`
		TenantID   string
		Username   string `json:"username"`
		Password   string `json:"password"`
		EventID    string `json:"event_id"`
	}
}

//GetServiceCheckInfoStruct to obtain application detection information
//swagger:parameters getServiceCheckInfo
type GetServiceCheckInfoStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	UUID string `json:"uuid"`
}

//PublicShare share shared structure
type PublicShare struct {
	ServiceKey string         `json:"service_key" validate:"service_key"`
	APPVersion string         `json:"app_version" validate:"app_version"`
	IsOuter    bool           `json:"is_outer" validate:"is_outer"`
	Action     string         `json:"action" validate:"action"`
	ShareID    string         `json:"share_id" validate:"share_id"`
	EventID    string         `json:"event_id" validate:"event_id"`
	Dest       string         `json:"dest" validate:"dest|in:yb,ys"`
	ServiceID  string         `json:"service_id" validate:"service_id"`
	ShareConf  ShareConfItems `json:"share_conf" validate:"share_conf"`
}

//SlugShare Slug type
type SlugShare struct {
	PublicShare
	ServiceKey    string `json:"service_key" validate:"service_key"`
	APPVersion    string `json:"app_version" validate:"app_version"`
	DeployVersion string `json:"deploy_version" validate:"deploy_version"`
	TenantID      string `json:"tenant_id" validate:"tenant_id"`
	Dest          string `json:"dest" validate:"dest|in:yb,ys"`
}

// ImageShare image type
type ImageShare struct {
	PublicShare
	Image string `json:"image" validate:"image"`
}

//ShareConfItems share related configuration
type ShareConfItems struct {
	FTPHost       string `json:"ftp_host" validate:"ftp_host"`
	FTPPort       int    `json:"ftp_port" validate:"ftp_port"`
	FTPUserName   string `json:"ftp_username" valiate:"ftp_username"`
	FTPPassWord   string `json:"ftp_password" validate:"ftp_password"`
	FTPNamespace  string `json:"ftp_namespace" validate:"ftp_namespace"`
	OuterRegistry string `json:"outer_registry" validate:"outer_registry"`
}

//AddDependencyStruct AddDependencyStruct
//swagger:parameters addDependency deleteDependency
type AddDependencyStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// application id that is dependent
		// in: body
		// required: true
		DepServiceID string `json:"dep_service_id"`
		// The type of application to be relied on, you need to pass a value when adding, and you don't need to pass a value when deleting
		// in: body
		// required: false
		DepServiceType string `json:"dep_service_type"`
		// Unknown, default pass 1, but not pass
		// in: body
		// required: false
		DepOrder string `json:"dep_order"`
	}
}

//AddEnvStruct AddEnvStruct
//swagger:parameters addEnv deleteEnv
type AddEnvStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// port
		// in: body
		// required: false
		ContainerPort int `json:"container_port"`
		// name
		// in: body
		// required: false
		Name string `json:"name"`
		// variable name
		// in: body
		// required: true
		AttrName string `json:"env_name"`
		// Variable value, you need to pass the value when adding, you can not pass when deleting
		// in: body
		// required: false
		AttrValue string `json:"env_value"`
		// Can it be modified
		// in: body
		// required: false
		IsChange bool `json:"is_change"`
		// Application range: inner or outer or both
		// in: body
		// required: false
		Scope string `json:"scope"`
	}
}

//RollBackStruct RollBackStruct
//swagger:parameters rollback
type RollBackStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// event_id
		// in: body
		// required: false
		EventID string `json:"event_id"`
		// The version number to roll back to
		// in: body
		// required: true
		DeployVersion string `json:"deploy_version"`
		// Operator
		// in: body
		// required: false
		Operator string `json:"operator"`
	}
}

//AddProbeStruct AddProbeStruct
//swagger:parameters addProbe updateProbe
type AddProbeStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// Probe id
		// in: body
		// required: true
		ProbeID string `json:"probe_id"`
		// mode
		// in: body
		// required: false
		Mode string `json:"mode"`
		// mode
		// in: body
		// required: false
		Scheme string `json:"scheme"`
		// path
		// in: body
		// required: false
		Path string `json:"path"`
		// Port, default is 80
		// in: body
		// required: false
		Port int `json:"port"`
		// Run the command
		// in: body
		// required: false
		Cmd string `json:"cmd"`
		// http request header, key=value, key2=value2
		// in: body
		// required: false
		HTTPHeader string `json:"http_header"`
		// Initialize waiting time, default is 1
		// in: body
		// required: false
		InitialDelaySecond int `json:"initial_delay_second"`
		// Detection interval time, default is 3
		// in: body
		// required: false
		PeriodSecond int `json:"period_second"`
		// Detection timeout, default is 30
		// in: body
		// required: false
		TimeoutSecond int `json:"timeout_second"`
		// Whether to enable
		// in: body
		// required: false
		IsUsed int `json:"is_used"`
		// The number of failed detections marked
		// in: body
		// required: false
		FailureThreshold int `json:"failure_threshold"`
		// The number of successful detections marked
		// in: body
		// required: false
		SuccessThreshold int `json:"success_threshold"`
	}
}

//DeleteProbeStruct DeleteProbeStruct
//swagger:parameters deleteProbe
type DeleteProbeStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	// in: body
	Body struct {
		// Probe id
		// in: body
		// required: true
		ProbeID string `json:"probe_id"`
	}
}

//PodsStructStruct PodsStructStruct
//swagger:parameters getPodsInfo
type PodsStructStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
}

//Login SSHLoginStruct
//swagger:parameters login
type Login struct {
	// in: body
	Body struct {
		// ip: port
		// in: body
		// required: true
		HostPort string `json:"hostport"`
		// Login type
		// in: body
		// required: true
		LoginType bool `json:"type"`
		// node type
		// in: body
		// required: true
		HostType string `json:"hosttype"`
		// root password
		// in: body
		// required: false
		RootPwd string `json:"pwd,omitempty"`
	}
}

//Labels LabelsStruct
//swagger:parameters labels
type Labels struct {
	// in: path
	// required: true
	NodeID string `json:"node"`
	// in: body
	Body struct {
		// label value list
		// in: body
		// required: true
		Labels []string `json:"labels"`
	}
}

//Model default field
type Model struct {
	ID uint
	//CreatedAt time.Time
}

//AddTenantServiceEnvVar application environment variables
type AddTenantServiceEnvVar struct {
	Model
	TenantID      string `validate:"tenant_id|between:30,33" json:"tenant_id"`
	ServiceID     string `validate:"service_id|between:30,33" json:"service_id"`
	ContainerPort int    `validate:"container_port|numeric_between:1,65535" json:"container_port"`
	Name          string `validate:"name" json:"name"`
	AttrName      string `validate:"env_name|required" json:"env_name"`
	AttrValue     string `validate:"env_value" json:"env_value"`
	IsChange      bool   `validate:"is_change|bool" json:"is_change"`
	Scope         string `validate:"scope|in:outer,inner,both,build" json:"scope"`
}

// DbModel return database model
func (a *AddTenantServiceEnvVar) DbModel(tenantID, componentID string) *dbmodel.TenantServiceEnvVar {
	return &dbmodel.TenantServiceEnvVar{
		TenantID:      tenantID,
		ServiceID:     componentID,
		Name:          a.Name,
		AttrName:      a.AttrName,
		AttrValue:     a.AttrValue,
		ContainerPort: a.ContainerPort,
		IsChange:      true,
		Scope:         a.Scope,
	}
}

//DelTenantServiceEnvVar application environment variables
type DelTenantServiceEnvVar struct {
	Model
	TenantID      string `validate:"tenant_id|between:30,33" json:"tenant_id"`
	ServiceID     string `validate:"service_id|between:30,33" json:"service_id"`
	ContainerPort int    `validate:"container_port|numeric_between:1,65535" json:"container_port"`
	Name          string `validate:"name" json:"name"`
	AttrName      string `validate:"env_name|required" json:"env_name"`
	AttrValue     string `validate:"env_value" json:"env_value"`
	IsChange      bool   `validate:"is_change|bool" json:"is_change"`
	Scope         string `validate:"scope|in:outer,inner,both,build" json:"scope"`
}

//ServicePorts service ports
type ServicePorts struct {
	Port []*TenantServicesPort
}

//TenantServicesPort application port information
type TenantServicesPort struct {
	Model
	TenantID       string `gorm:"column:tenant_id;size:32" validate:"tenant_id|between:30,33" json:"tenant_id"`
	ServiceID      string `gorm:"column:service_id;size:32" validate:"service_id|between:30,33" json:"service_id"`
	ContainerPort  int    `gorm:"column:container_port" validate:"container_port|required|numeric_between:1,65535" json:"container_port"`
	MappingPort    int    `gorm:"column:mapping_port" validate:"mapping_port|required|numeric_between:1,65535" json:"mapping_port"`
	Protocol       string `gorm:"column:protocol" validate:"protocol|required|in:http,https,stream,grpc" json:"protocol"`
	PortAlias      string `gorm:"column:port_alias" validate:"port_alias|required|alpha_dash" json:"port_alias"`
	K8sServiceName string `gorm:"column:k8s_service_name" json:"k8s_service_name"`
	IsInnerService bool   `gorm:"column:is_inner_service" validate:"is_inner_service|bool" json:"is_inner_service"`
	IsOuterService bool   `gorm:"column:is_outer_service" validate:"is_outer_service|bool" json:"is_outer_service"`
}

// DbModel return database model
func (p *TenantServicesPort) DbModel(tenantID, componentID string) *dbmodel.TenantServicesPort {
	isInnerService := p.IsInnerService
	isOuterService := p.IsOuterService
	return &dbmodel.TenantServicesPort{
		TenantID:       tenantID,
		ServiceID:      componentID,
		ContainerPort:  p.ContainerPort,
		MappingPort:    p.MappingPort,
		Protocol:       p.Protocol,
		PortAlias: p.PortAlias,
		IsInnerService: &isInnerService,
		IsOuterService: &isOuterService,
		K8sServiceName: p.K8sServiceName,
	}
}

// AddServicePort service port
// swagger:parameters addPort updatePort
type AddServicePort struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		//in: body
		ServicePorts
	}
}

type plugin struct {
	// the container port for this serviceport
	// in: body
	// required: true
	ContainerPort int32 `json:"container_port"`
	// the mapping port for this serviceport
	// in: body
	// required: true
	MappingPort int32 `json:"mapping_port"`
	// the protocol for this serviceport
	// in: body
	// required: true
	Protocol string `json:"protocol"`
	// the port alias for this serviceport
	// in: body
	// required: true
	PortAlias string `json:"port_alias"`
	// Whether to enable internal services
	// in: body
	Inner bool `json:"is_inner_service"`
	// Whether to open external services
	// in: body
	Outer bool `json:"is_outer_service"`
}

//ServiceProbe application probe information
type ServiceProbe struct {
	Model
	ServiceID string `gorm:"column:service_id;size:32" json:"service_id" validate:"service_id|between:30,33"`
	ProbeID   string `gorm:"column:probe_id;size:32" json:"probe_id" validate:"probe_id|required|between:30,33"`
	Mode      string `gorm:"column:mode;default:'liveness'" json:"mode" validate:"mode"`
	Scheme    string `gorm:"column:scheme;default:'scheme'" json:"scheme" validate:"scheme"`
	Path      string `gorm:"column:path" json:"path" validate:"path"`
	Port      int    `gorm:"column:port;size:5;default:80" json:"port" validate:"port|numeric_between:1,65535"`
	Cmd       string `gorm:"column:cmd;size:150" json:"cmd" validate:"cmd"`
	//http request header, key=value,key2=value2
	HTTPHeader string `gorm:"column:http_header;size:300" json:"http_header" validate:"http_header"`
	//Initialize waiting time
	InitialDelaySecond int `gorm:"column:initial_delay_second;size:2;default:1" json:"initial_delay_second" validate:"initial_delay_second"`
	//Detection interval time
	PeriodSecond int `gorm:"column:period_second;size:2;default:3" json:"period_second" validate:"period_second"`
	//Detection timeout
	TimeoutSecond int `gorm:"column:timeout_second;size:3;default:30" json:"timeout_second" validate:"timeout_second"`
	//Whether to enable
	IsUsed int `gorm:"column:is_used;size:1;default:0" json:"is_used" validate:"is_used|in:0,1"`
	//The number of detections marked as failed
	FailureThreshold int `gorm:"column:failure_threshold;size:2;default:3" json:"failure_threshold" validate:"failure_threshold"`
	//The number of successful detections marked
	SuccessThreshold int    `gorm:"column:success_threshold;size:2;default:1" json:"success_threshold" validate:"success_threshold"`
	FailureAction    string `json:"failure_action" validate:"failure_action"`
}

// DbModel return database model
func (p *ServiceProbe) DbModel(componentID string) *dbmodel.TenantServiceProbe {
	return &dbmodel.TenantServiceProbe{
		ServiceID:          componentID,
		Cmd:                p.Cmd,
		FailureThreshold:   p.FailureThreshold,
		HTTPHeader:         p.HTTPHeader,
		InitialDelaySecond: p.InitialDelaySecond,
		IsUsed:             &p.IsUsed,
		Mode: p.Mode,
		Path:               p.Path,
		PeriodSecond:       p.PeriodSecond,
		Port: p.Port,
		ProbeID: p.ProbeID,
		Scheme:             p.Scheme,
		SuccessThreshold:   p.SuccessThreshold,
		TimeoutSecond:      p.TimeoutSecond,
		FailureAction:      p.FailureAction,
	}
}

//TenantServiceVolume application persistent record
type TenantServiceVolume struct {
	Model
	ServiceID string `gorm:"column:service_id;size:32" json:"service_id" validate:"service_id"`
	//Service type
	Category   string `gorm:"column:category;size:50" json:"category" validate:"category|required"`
	HostPath   string `gorm:"column:host_path" json:"host_path" validate:"host_path|required"`
	VolumePath string `gorm:"column:volume_path" json:"volume_path" validate:"volume_path|required"`
	IsReadOnly bool   `gorm:"column:is_read_only;default:false" json:"is_read_only" validate:"is_read_only|bool"`
}

// GetSupportProtocols GetSupportProtocols
// swagger:parameters getSupportProtocols
type GetSupportProtocols struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
}

//ServiceShare service share
// swagger:parameters shareService
type ServiceShare struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias"`
	//in: body
	Body struct {
		//in: body
		//App share Key
		ServiceKey string `json:"service_key" validate:"service_key|required"`
		AppVersion string `json:"app_version" validate:"app_version|required"`
		EventID    string `json:"event_id"`
		ShareUser  string `json:"share_user"`
		ShareScope string `json:"share_scope"`
		ImageInfo  struct {
			HubURL      string `json:"hub_url"`
			HubUser     string `json:"hub_user"`
			HubPassword string `json:"hub_password"`
			Namespace   string `json:"namespace"`
			IsTrust     bool   `json:"is_trust,omitempty" validate:"is_trust"`
		} `json:"image_info,omitempty"`
		SlugInfo struct {
			Namespace   string `json:"namespace"`
			FTPHost     string `json:"ftp_host"`
			FTPPort     string `json:"ftp_port"`
			FTPUser     string `json:"ftp_username"`
			FTPPassword string `json:"ftp_password"`
		} `json:"slug_info,omitempty"`
	}
}

//ExportAppStruct -
type ExportAppStruct struct {
	SourceDir string `json:"source_dir"`
	Body      struct {
		EventID       string `json:"event_id"`
		GroupKey string `json:"group_key"` // TODO consider removing
		Version string `json:"version"` // TODO consider removing
		Format        string `json:"format"`    // only kato-app/docker-compose
		GroupMetadata string `json:"group_metadata"`
	}
}

// BatchOperationReq beatch operation request body
type BatchOperationReq struct {
	Operator   string `json:"operator"`
	TenantName string `json:"tenant_name"`
	Body       struct {
		Operation string                 `json:"operation" validate:"operation|required|in:start,stop,build,upgrade"`
		Builds    []*ComponentBuildReq   `json:"build_infos,omitempty"`
		Starts    []*ComponentStartReq   `json:"start_infos,omitempty"`
		Stops     []*ComponentStopReq    `json:"stop_infos,omitempty"`
		Upgrades  []*ComponentUpgradeReq `json:"upgrade_infos,omitempty"`
	}
}

//BuildImageInfo -
type BuildImageInfo struct {
	// Mirror address
	// in: body
	// required: false
	ImageURL string `json:"image_url" validate:"image_url"`
	User     string `json:"user" validate:"user"`
	Password string `json:"password" validate:"password"`
	Cmd      string `json:"cmd"`
}

//BuildCodeInfo -
type BuildCodeInfo struct {
	// git address
	// in: body
	// required: false
	RepoURL string `json:"repo_url" validate:"repo_url"`
	// branch branch information
	// in: body
	// required: false
	Branch string `json:"branch" validate:"branch"`
	// operator
	// in: body
	// required: false
	Lang string `json:"lang" validate:"lang"`
	// Code server type
	// in: body
	// required: false
	ServerType string `json:"server_type" validate:"server_type"`
	Runtime    string `json:"runtime"`
	User       string `json:"user" validate:"user"`
	Password   string `json:"password" validate:"password"`
	//for .netcore source type, need cmd
	Cmd string `json:"cmd"`
}

//BuildSlugInfo -
type BuildSlugInfo struct {
	SlugPath    string `json:"slug_path"`
	FTPHost     string `json:"ftp_host"`
	FTPPort     string `json:"ftp_port"`
	FTPUser     string `json:"ftp_username"`
	FTPPassword string `json:"ftp_password"`
}

//FromImageBuildKing build from image
var FromImageBuildKing = "build_from_image"

//FromCodeBuildKing build from code
var FromCodeBuildKing = "build_from_source_code"

//FromMarketImageBuildKing build from market image
var FromMarketImageBuildKing = "build_from_market_image"

//FromMarketSlugBuildKing build from market slug
var FromMarketSlugBuildKing = "build_from_market_slug"

// ComponentBuildReq -
type ComponentBuildReq struct {
	ComponentOpGeneralReq
	// variable
	// in: body
	// required: false
	BuildENVs map[string]string `json:"envs" validate:"envs"`
	// Application build type
	// in: body
	// required: true
	Kind string `json:"kind" validate:"kind|required"`
	// Follow-up action, one-click deployment based on the value, if the value is not passed, only the build is performed by default
	// in: body
	// required: false
	Action string `json:"action" validate:"action"`
	// Plan Version
	PlanVersion string `json:"plan_version"`
	// Deployed version number, The version is generated by the API
	// in: body
	DeployVersion string `json:"deploy_version" validate:"deploy_version"`
	// Build task initiator
	//in: body
	Operator string `json:"operator" validate:"operator"`
	//build form image
	ImageInfo BuildImageInfo `json:"image_info,omitempty"`
	//build from code
	CodeInfo BuildCodeInfo `json:"code_info,omitempty"`
	//Used for cloud city code package creation
	SlugInfo BuildSlugInfo `json:"slug_info,omitempty"`
	//tenantName
	TenantName string `json:"-"`
}

// GetEventID -
func (b *ComponentBuildReq) GetEventID() string {
	if b.EventID == "" {
		b.EventID = util.NewUUID()
	}
	return b.EventID
}

// BatchOpFailureItem -
func (b *ComponentBuildReq) BatchOpFailureItem () *ComponentOpResult {
	return &ComponentOpResult{
		ServiceID: b.ServiceID,
		EventID:   b.EventID,
		Operation: "build",
		Status: BatchOpResultItemStatusFailure,
	}
}

// GetVersion -
func (b *ComponentBuildReq) GetVersion() string {
	return b.DeployVersion
}

//SetVersion -
func (b *ComponentBuildReq) SetVersion(string) {
	// no need
	return
}

// OpType -
func (b *ComponentBuildReq) OpType() string {
	return "build-service"
}

// GetComponentID -
func (b *ComponentBuildReq) GetComponentID() string {
	return b.ServiceID
}

// TaskBody returns a task body.
func (b *ComponentBuildReq) TaskBody(cpt *dbmodel.TenantServices) interface{} {
	return nil
}

// UpdateBuildVersionReq -
type UpdateBuildVersionReq struct {
	PlanVersion string `json:"plan_version" validate:"required"`
}

//ComponentUpgradeReq -
type ComponentUpgradeReq struct {
	ComponentOpGeneralReq
	//UpgradeVersion The target version of the upgrade
	//If empty, the same version is upgraded
	UpgradeVersion string `json:"upgrade_version"`
}

// GetEventID -
func (u *ComponentUpgradeReq) GetEventID() string {
	if u.EventID == "" {
		u.EventID = util.NewUUID()
	}
	return u.EventID
}

// BatchOpFailureItem -
func (u *ComponentUpgradeReq) BatchOpFailureItem () *ComponentOpResult {
	return &ComponentOpResult{
		ServiceID: u.ServiceID,
		EventID:   u.GetEventID(),
		Operation: "upgrade",
		Status: BatchOpResultItemStatusFailure,
	}
}

// GetVersion -
func (u *ComponentUpgradeReq) GetVersion() string {
	return u.UpgradeVersion
}

//SetVersion -
func (u *ComponentUpgradeReq) SetVersion(version string) {
	if u.UpgradeVersion == "" {
		u.UpgradeVersion = version
	}
}

// GetComponentID -
func (u *ComponentUpgradeReq) GetComponentID() string {
	return u.ServiceID
}

// TaskBody returns the task body.
func (u *ComponentUpgradeReq) TaskBody(cpt *dbmodel.TenantServices) interface{} {
	return &dmodel.RollingUpgradeTaskBody{
		TenantID:         cpt.TenantID,
		ServiceID: cpt.ServiceID,
		NewDeployVersion: u.UpgradeVersion,
		EventID:          u.GetEventID(),
		Configs:          u.Configs,
	}
}

// OpType -
func (u *ComponentUpgradeReq) OpType() string {
	return "upgrade-service"
}

//RollbackInfoRequestStruct -
type RollbackInfoRequestStruct struct {
	//RollBackVersion The target version of the rollback
	RollBackVersion string `json:"upgrade_version"`
	//Event trace ID
	EventID   string            `json:"event_id"`
	ServiceID string            `json:"service_id"`
	Configs   map[string]string `json:"configs"`
}

//BuildMQBodyFrom -
func BuildMQBodyFrom(app *ExportAppStruct) *MQBody {
	return &MQBody{
		EventID:   app.Body.EventID,
		GroupKey:  app.Body.GroupKey,
		Version:   app.Body.Version,
		Format:    app.Body.Format,
		SourceDir: app.SourceDir,
	}
}

//MQBody -
type MQBody struct {
	EventID   string `json:"event_id"`
	GroupKey  string `json:"group_key"`
	Version   string `json:"version"`
	Format    string `json:"format"` // only kato-app/docker-compose
	SourceDir string `json:"source_dir"`
}

//NewAppStatusFromExport -
func NewAppStatusFromExport(app *ExportAppStruct) *dbmodel.AppStatus {
	return &dbmodel.AppStatus{
		Format:    app.Body.Format,
		EventID:   app.Body.EventID,
		SourceDir: app.SourceDir,
		Status:    "exporting",
	}
}

//ImportAppStruct -
type ImportAppStruct struct {
	EventID      string       `json:"event_id"`
	SourceDir    string       `json:"source_dir"`
	Apps         []string     `json:"apps"`
	Format       string       `json:"format"`
	ServiceImage ServiceImage `json:"service_image"`
	ServiceSlug  ServiceSlug  `json:"service_slug"`
}

//ServiceImage -
type ServiceImage struct {
	HubURL      string `json:"hub_url"`
	HubUser     string `json:"hub_user"`
	HubPassword string `json:"hub_password"`
	NameSpace   string `json:"namespace"`
}

//ServiceSlug -
type ServiceSlug struct {
	FtpHost     string `json:"ftp_host"`
	FtpPort     string `json:"ftp_port"`
	FtpUsername string `json:"ftp_username"`
	FtpPassword string `json:"ftp_password"`
	NameSpace   string `json:"namespace"`
}

//NewAppStatusFromImport -
func NewAppStatusFromImport(app *ImportAppStruct) *dbmodel.AppStatus {
	var apps string
	for _, app := range app.Apps {
		app += ":pending"
		if apps == "" {
			apps += app
		} else {
			apps += "," + app
		}
	}

	return &dbmodel.AppStatus{
		EventID:   app.EventID,
		Format:    app.Format,
		SourceDir: app.SourceDir,
		Apps:      apps,
		Status:    "importing",
	}
}

// Application -
type Application struct {
	EID             string   `json:"eid" validate:"required"`
	AppName         string   `json:"app_name" validate:"required"`
	AppType         string   `json:"app_type" validate:"required,oneof=kato helm"`
	ConsoleAppID    int64    `json:"console_app_id"`
	AppID           string   `json:"app_id"`
	TenantID        string   `json:"tenant_id"`
	ServiceIDs      []string `json:"service_ids"`
	AppStoreName    string   `json:"app_store_name"`
	AppStoreURL     string   `json:"app_store_url"`
	AppTemplateName string   `json:"app_template_name"`
	Version         string   `json:"version"`
}

// CreateAppRequest -
type CreateAppRequest struct {
	AppsInfo []Application `json:"apps_info"`
}

// CreateAppResponse -
type CreateAppResponse struct {
	AppID       int64  `json:"app_id"`
	RegionAppID string `json:"region_app_id"`
}

// ListAppResponse -
type ListAppResponse struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
	Total    int64                  `json:"total"`
	Apps     []*dbmodel.Application `json:"apps"`
}

// ListServiceResponse -
type ListServiceResponse struct {
	Page     int                       `json:"page"`
	PageSize int                       `json:"pageSize"`
	Total    int64                     `json:"total"`
	Services []*dbmodel.TenantServices `json:"services"`
}

// UpdateAppRequest -
type UpdateAppRequest struct {
	AppName        string   `json:"app_name"`
	GovernanceMode string   `json:"governance_mode"`
	Overrides      []string `json:"overrides"`
	Version        string   `json:"version"`
	Revision       int      `json:"revision"`
}

// NeedUpdateHelmApp check if necessary to update the helm app.
func (u *UpdateAppRequest) NeedUpdateHelmApp() bool {
	return len(u.Overrides) > 0 || u.Version != "" || u.Revision != 0
}

// BindServiceRequest -
type BindServiceRequest struct {
	ServiceIDs []string `json:"service_ids"`
}

// InstallAppReq -
type InstallAppReq struct {
	Overrides []string `json:"overrides"`
}

// ParseAppServicesReq -
type ParseAppServicesReq struct {
	Values string `json:"values"`
}

// ConfigGroupService -
type ConfigGroupService struct {
	ServiceID    string `json:"service_id"`
	ServiceAlias string `json:"service_alias"`
}

// DbModel return database model
func (c ConfigGroupService) DbModel(appID, configGroupName string) *dbmodel.ConfigGroupService {
	return &dbmodel.ConfigGroupService{
		AppID: appID,
		ConfigGroupName: configGroupName,
		ServiceID:       c.ServiceID,
		ServiceAlias: c.ServiceAlias,
	}
}

// ConfigItem -
type ConfigItem struct {
	AppID           string `json:"-"`
	ConfigGroupName string `json:"-"`
	ItemKey         string `json:"item_key" validate:"required,max=255"`
	ItemValue       string `json:"item_value" validate:"required,max=65535"`
}

// DbModel return database model
func (c ConfigItem) DbModel(appID, configGroupName string) *dbmodel.ConfigGroupItem {
	return &dbmodel.ConfigGroupItem{
		AppID: appID,
		ConfigGroupName: configGroupName,
		ItemKey:         c.ItemKey,
		ItemValue:       c.ItemValue,
	}
}

// ApplicationConfigGroup -
type ApplicationConfigGroup struct {
	AppID           string       `json:"app_id"`
	ConfigGroupName string       `json:"config_group_name" validate:"required,alphanum,min=2,max=64"`
	DeployType      string       `json:"deploy_type" validate:"required,oneof=env configfile"`
	ServiceIDs      []string     `json:"service_ids"`
	ConfigItems     []ConfigItem `json:"config_items"`
	Enable          bool         `json:"enable"`
}

// AppConfigGroup Interface for synchronizing application configuration groups
type AppConfigGroup struct {
	ConfigGroupName     string               `json:"config_group_name" validate:"required,alphanum,min=2,max=64"`
	DeployType          string               `json:"deploy_type" validate:"required,oneof=env configfile"`
	ConfigItems         []ConfigItem         `json:"config_items"`
	ConfigGroupServices []ConfigGroupService `json:"config_group_services"`
	Enable              bool                 `json:"enable"`
}

// DbModel return database model
func (a AppConfigGroup) DbModel(appID string) *dbmodel.ApplicationConfigGroup {
	return &dbmodel.ApplicationConfigGroup{
		AppID: appID,
		ConfigGroupName: a.ConfigGroupName,
		DeployType:      a.DeployType,
		Enable:          a.Enable,
	}
}

// ApplicationConfigGroupResp -
type ApplicationConfigGroupResp struct {
	CreateTime      time.Time                     `json:"create_time"`
	AppID           string                        `json:"app_id"`
	ConfigGroupName string                        `json:"config_group_name"`
	DeployType      string                        `json:"deploy_type"`
	Services        []*dbmodel.ConfigGroupService `json:"services"`
	ConfigItems     []*dbmodel.ConfigGroupItem    `json:"config_items"`
	Enable          bool                          `json:"enable"`
}

// UpdateAppConfigGroupReq -
type UpdateAppConfigGroupReq struct {
	ServiceIDs  []string     `json:"service_ids"`
	ConfigItems []ConfigItem `json:"config_items" validate:"required"`
	Enable      bool         `json:"enable"`
}

// ListApplicationConfigGroupResp -
type ListApplicationConfigGroupResp struct {
	ConfigGroup []ApplicationConfigGroupResp `json:"config_group"`
	Total       int64                        `json:"total"`
	Page        int                          `json:"page"`
	PageSize    int                          `json:"pageSize"`
}

// CheckResourceNameReq -
type CheckResourceNameReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// CheckResourceNameResp -
type CheckResourceNameResp struct {
	Name string `json:"name"`
}

// HelmAppRelease -
type HelmAppRelease struct {
	Revision    int    `json:"revision"`
	Updated     string `json:"updated"`
	Status      string `json:"status"`
	Chart       string `json:"chart"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}

// AppConfigGroupRelations -
type AppConfigGroupRelations struct {
	ConfigGroupName string `json:"config_group_name"`
}

// DbModel return database model
func (a *AppConfigGroupRelations) DbModel(appID, serviceID, serviceAlias string) *dbmodel.ConfigGroupService {
	return &dbmodel.ConfigGroupService{
		AppID: appID,
		ConfigGroupName: a.ConfigGroupName,
		ServiceID:       serviceID,
		ServiceAlias: serviceAlias,
	}
}

// SyncAppConfigGroup -
type SyncAppConfigGroup struct {
	AppConfigGroups []AppConfigGroup `json:"app_config_groups"`
}

// AppStatusesReq -
type AppStatusesReq struct {
	AppIDs []string `json:"app_ids"`
}
