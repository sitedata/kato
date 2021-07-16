package util

import (
	"github.com/sirupsen/logrus"
	"strings"

	"encoding/json"
	api_model "github.com/gridworkz/kato/api/model"
	dbmodel "github.com/gridworkz/kato/db/model"
	v1 "github.com/gridworkz/kato/worker/appm/types/v1"
	"github.com/gridworkz/kato/worker/server/pb"
	"github.com/twinj/uuid"
)

// SetVolumeDefaultValue
func SetVolumeDefaultValue(info *dbmodel.TenantServiceVolume) {
	if info.VolumeName == "" {
		info.VolumeName = uuid.NewV4().String()
	}

	if info.AccessMode != "" {
		info.AccessMode = strings.ToUpper(info.AccessMode)
	} else {
		info.AccessMode = "RWO"
	}

	if info.SharePolicy == "" {
		info.SharePolicy = "exclusive"
	}

	if info.BackupPolicy == "" {
		info.BackupPolicy = "exclusive"
	}

	if info.ReclaimPolicy == "" {
		info.ReclaimPolicy = "retain"
	}
}

// ParseVolumeTypeOption parse volume Option name show and volume type
func ParseVolumeTypeOption(detail *pb.StorageClassDetail) string {
	volumeType := transferVolumeTypeOption(detail.Name, detail.Provisioner, detail.Parameters)
	if volumeType != nil {
		return volumeType.String()
	}

	return "unknown"
}

func transferVolumeTypeOption(name string, opts ...interface{}) *dbmodel.VolumeType {
	if name == v1.KatoStatefuleShareStorageClass {
		return &dbmodel.ShareFileVolumeType
	}
	if name == v1.KatoStatefuleLocalStorageClass {
		return &dbmodel.LocalVolumeType
	}
	vt := dbmodel.MakeNewVolume(name)
	return &vt
}

// opts[0]: kind is storageClass's provisioner
// opts[1]: parameters is storageClass's parameter
func transferCustomVolumeOptionName2Kind(opts ...interface{}) *dbmodel.VolumeType {
	if len(opts) != 2 {
		return nil
	}
	kind := opts[0].(string)
	if strings.HasSuffix(kind, "rbd") {
		if parameters, ok := opts[1].(map[string]string); ok {
			if parameters["adminId"] != "" && parameters["monitors"] != "" && parameters["pool"] != "" && parameters["userId"] != "" {
				return &dbmodel.CephRBDVolumeType
			}
		}
	}
	if kind == "alicloud/disk" {
		return &dbmodel.AliCloudVolumeType
	}
	return nil
}

// HackVolumeOptionDetailFromDB hack volumeOptionDetail from db
func HackVolumeOptionDetailFromDB(detail *api_model.VolumeTypeStruct, data *dbmodel.TenantServiceVolumeType) {
	if data != nil {
		detail.Description = data.Description
		detail.NameShow = data.NameShow
		if err := json.Unmarshal([]byte(data.CapacityValidation), &detail.CapacityValidation); err != nil {
			logrus.Warnf("unmarshal volumetype's capacityValidation error: %s, set capacityValidation to default", err.Error())
			detail.CapacityValidation = defaultcapacityValidation
		}
		detail.AccessMode = strings.Split(data.AccessMode, ",")
		detail.SharePolicy = strings.Split(data.SharePolicy, ",")
		detail.BackupPolicy = strings.Split(data.BackupPolicy, ",")
		detail.ReclaimPolicy = data.ReclaimPolicy
		detail.Sort = data.Sort
	}
}

var defaultcapacityValidation map[string]interface{}

func init() {
	capacityValidation := make(map[string]interface{})
	capacityValidation["min"] = 1
	capacityValidation["required"] = false
	capacityValidation["max"] = 999999999
}

// HackVolumeOptionDetail hack volume Option detail, like accessMode, sharePolicy, backupPolicy
func HackVolumeOptionDetail(volumeType string, detail *api_model.VolumeTypeStruct, more ...interface{}) {
	/*
		RWO - ReadWriteOnce
		ROX - ReadOnlyMany
		RWX - ReadWriteMany
	*/
	detail.AccessMode = append(detail.AccessMode, hackVolumeOptionAccessMode(volumeType)...)
	detail.SharePolicy = append(detail.SharePolicy, hackVolumeOptionSharePolicy(volumeType)...)
	detail.BackupPolicy = append(detail.BackupPolicy, hackVolumeOptionBackupPolicy(volumeType)...)
	detail.CapacityValidation = hackVolumeOptionCapacityValidation(volumeType)
	detail.Description = hackVolumeOptionDesc(volumeType)
	detail.NameShow = hackVolumeOptionNameShow(volumeType)
	if len(more) == 4 {
		detail.ReclaimPolicy = more[1].(string)
	}
}

func hackVolumeOptionNameShow(volumeType string) string {
	nameShow := volumeType
	if volumeType == "alicloud-disk-available" {
		nameShow = "Alibaba Cloud Disk (Smart Choice)"
	} else if volumeType == "alicloud-disk-common" {
		nameShow = "Alibaba Cloud Disk (Basic)"
	} else if volumeType == "alicloud-disk-efficiency" {
		nameShow = "Alibaba Cloud Disk (efficient)"
	} else if volumeType == "alicloud-disk-ssd" {
		nameShow = "Alibaba Cloud Disk (SSD)"
	}
	return nameShow
}

func hackVolumeOptionDesc(vt string) string {
	volumeType := dbmodel.VolumeType(vt)
	switch volumeType {
	case dbmodel.ShareFileVolumeType:
		return "default distributed file storage, which can be shared and mounted within tenants, suitable for all types of applications"
	case dbmodel.LocalVolumeType:
		return "default local storage device, suitable for stateful database services"
	case dbmodel.MemoryFSVolumeType:
		return "default memory-based storage device, the capacity is limited by the amount of memory. Data is lost when the application restarts, suitable for high-speed temporary storage of data"
	default:
		return ""
	}
}

func hackVolumeOptionCapacityValidation(volumeType string) map[string]interface{} {
	data := make(map[string]interface{})
	data["required"] = true
	data["default"] = 1
	if strings.HasPrefix(volumeType, "alicloud-disk") {
		data["min"] = 20
		data["default"] = 20
		data["max"] = 32768 // [ali-cloud-disk usage limit](https://help.aliyun.com/document_detail/25412.html?spm=5176.2020520101.0.0.41d84df5faliP4)
	} else {
		data["min"] = 1
		data["max"] = 999999999
	}
	return data
}

/*

## volume accessMode
---

Volume Plugin 		| ReadWriteOnce        |    ReadOnlyMany          | ReadWriteMany
--------------------|----------------------|--------------------------|-----------------------
AWSElasticBlockStore| 	✓		           |   	-	      	          |  -
AzureFile			|    ✓			       |   	✓		  	          |  ✓
AzureDisk			|    ✓			       |   	-		  	          |  -
CephFS			    |    ✓		           |      ✓		          	  |  ✓
Cinder			    |    ✓		           |      -		  	          |  -
CSI					| depends on the driver|	depends on the driver |	depends on the driver
FC					|  ✓				   |   ✓					  | -
FlexVolume			| ✓					   |	✓					  | depends on the driver
Flocker				|	✓				   |  -						  | -
GCEPersistentDisk	|	✓				   | ✓						  | -
Glusterfs			|  ✓				   | ✓	                      | ✓
HostPath	        |  ✓				   | -						  | -
iSCSI				| ✓					   | ✓						  | -
Quobyte				| ✓					   | ✓						  | ✓
NFS					| ✓					   | ✓						  | ✓
RBD					| ✓					   | ✓						  | -
VsphereVolume		| ✓					   | -						  | - (works when Pods are collocated)
PortworxVolume		| ✓					   | -						  | ✓
ScaleIO				| ✓					   | ✓						  | -
StorageOS			| ✓					   | -						  | -

*/
func hackVolumeOptionAccessMode(vt string) []string {
	volumeType := dbmodel.VolumeType(vt)
	switch volumeType {
	case dbmodel.ShareFileVolumeType:
		return []string{"RWO", "ROX", "RWX"}
	case dbmodel.LocalVolumeType:
		return []string{"RWO", "ROX", "RWX"}
	case dbmodel.ConfigFileVolumeType:
		return []string{"ROX"}
	case dbmodel.MemoryFSVolumeType:
		return []string{"ROX"}
	default:
		return []string{"RWO"}
	}
}

// TODO finish volume share policy
func hackVolumeOptionSharePolicy(volumeType string) []string {
	return []string{"exclusive"}
}

// TODO finish vollume backup policy
func hackVolumeOptionBackupPolicy(volumeType string) []string {
	return []string{"exclusive"}
}
