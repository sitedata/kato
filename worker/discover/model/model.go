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
	"time"

	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

//TaskType task type
type TaskType string

//Task task
type Task struct {
	Type       TaskType  `json:"type"`
	Body       TaskBody  `json:"body"`
	CreateTime time.Time `json:"time,omitempty"`
	User       string    `json:"user"`
}

//NewTask 从json bytes data create task
func NewTask(data []byte) (*Task, error) {
	taskType := gjson.GetBytes(data, "type").String()
	body := CreateTaskBody(taskType)
	task := Task{
		Body: body,
	}
	err := ffjson.Unmarshal(data, &task)
	if err != nil {
		return nil, err
	}
	return &task, err
}

//TransTask transtask
func TransTask(task *pb.TaskMessage) (*Task, error) {
	timeT, _ := time.Parse(time.RFC3339, task.CreateTime)
	return &Task{
		Type:       TaskType(task.TaskType),
		Body:       NewTaskBody(task.TaskType, task.TaskBody),
		CreateTime: timeT,
		User:       task.User,
	}, nil
}

//NewTaskBody new task body
func NewTaskBody(taskType string, body []byte) TaskBody {
	switch taskType {
	case "start":
		b := StartTaskBody {}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "stop":
		b := StopTaskBody {}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "restart":
		b := RestartTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "rolling_upgrade":
		b := RollingUpgradeTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "rollback":
		b := RollBackTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "group_start":
		b := GroupStartTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "group_stop":
		b := GroupStopTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "horizontal_scaling":
		b := HorizontalScalingTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "vertical_scaling":
		b := VerticalScalingTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "apply_rule":
		b := ApplyRuleTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			logrus.Debugf("error unmarshal data: %v", err)
			return nil
		}
		return &b
	case "apply_plugin_config":
		b := &ApplyPluginConfigTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "service_gc":
		b := ServiceGCTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "delete_tenant":
		b := &DeleteTenantTaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	case "refreshhpa":
		b := &RefreshHPATaskBody{}
		err := ffjson.Unmarshal(body, &b)
		if err != nil {
			return nil
		}
		return b
	default:
		return DefaultTaskBody{}
	}
}

//CreateTaskBody creates an entity through the type string
func CreateTaskBody(taskType string) TaskBody {
	switch taskType {
	case "start":
		return StartTaskBody{}
	case "stop":
		return StopTaskBody{}
	case "restart":
		return RestartTaskBody{}
	case "rolling_upgrade":
		return RollingUpgradeTaskBody{}
	case "rollback":
		return RollBackTaskBody{}
	case "group_start":
		return GroupStartTaskBody{}
	case "group_stop":
		return GroupStopTaskBody{}
	case "horizontal_scaling":
		return HorizontalScalingTaskBody{}
	case "vertical_scaling":
		return VerticalScalingTaskBody{}
	case "apply_plugin_config":
		return ApplyPluginConfigTaskBody{}
	case "delete_tenant":
		return DeleteTenantTaskBody{}
	case "refreshhpa":
		return RefreshHPATaskBody{}
	default:
		return DefaultTaskBody{}
	}
}

//TaskBody task body
type TaskBody interface{}

//StartTaskBody starts the operation task body
type StartTaskBody struct {
	TenantID      string            `json:"tenant_id"`
	ServiceID     string            `json:"service_id"`
	DeployVersion string            `json:"deploy_version"`
	EventID       string            `json:"event_id"`
	Configs       map[string]string `json:"configs"`
	// When determining the startup sequence of services, you need to know the services they depend on
	DepServiceIDInBootSeq []string `json:"dep_service_ids_in_boot_seq"`
}

//StopTaskBody stops the operation task body
type StopTaskBody struct {
	TenantID      string            `json:"tenant_id"`
	ServiceID     string            `json:"service_id"`
	DeployVersion string            `json:"deploy_version"`
	EventID       string            `json:"event_id"`
	Configs       map[string]string `json:"configs"`
}

//HorizontalScalingTaskBody Horizontal scaling operation task body
type HorizontalScalingTaskBody struct {
	TenantID  string `json:"tenant_id"`
	ServiceID string `json:"service_id"`
	Replicas  int32  `json:"replicas"`
	EventID   string `json:"event_id"`
	Username  string `json:"username"`
}

//VerticalScalingTaskBody vertical scaling operation task body
type VerticalScalingTaskBody struct {
	TenantID        string `json:"tenant_id"`
	ServiceID       string `json:"service_id"`
	ContainerCPU    *int   `json:"container_cpu"`
	ContainerMemory *int   `json:"container_memory"`
	ContainerGPU    *int   `json:"container_gpu"`
	EventID         string `json:"event_id"`
}

//RestartTaskBody restart operation task body
type RestartTaskBody struct {
	TenantID      string `json:"tenant_id"`
	ServiceID     string `json:"service_id"`
	DeployVersion string `json:"deploy_version"`
	EventID       string `json:"event_id"`
	//Restart policy, this policy is not guaranteed to take effect
	//For example, if the application is a stateful service, if this policy is configured to be started and then closed, this policy will not take effect
	//Stateless services are used by default to start and then close to ensure that the service is not affected
	Strategy []string          `json:"strategy"`
	Configs  map[string]string `json:"configs"`
}

//StrategyIsValid verifies whether the strategy is valid
//The strategy includes the following values:
// prestart starts first and then closes
// prestop is closed first and then started
// rollingupdate rolling form
// grayupdate gray form
// bluegreenupdate blue-green form
//
func StrategyIsValid(strategy []string, serviceDeployType string) bool {
	return false
}

//RollingUpgradeTaskBody Upgrade operation task body
type RollingUpgradeTaskBody struct {
	TenantID         string            `json:"tenant_id"`
	ServiceID        string            `json:"service_id"`
	NewDeployVersion string            `json:"deploy_version"`
	EventID          string            `json:"event_id"`
	Strategy         []string          `json:"strategy"`
	Configs          map[string]string `json:"configs"`
}

//RollBackTaskBody rollback operation task body
type RollBackTaskBody struct {
	TenantID  string `json:"tenant_id"`
	ServiceID string `json:"service_id"`
	//current version
	CurrentDeployVersion string `json:"current_deploy_version"`
	//Roll back the target version
	OldDeployVersion string `json:"old_deploy_version"`
	EventID          string `json:"event_id"`
	//Restart policy, this policy is not guaranteed to take effect
	//For example, if the application is a stateful service, if this policy is configured to be started and then closed, this policy will not take effect
	//Stateless services are used by default to start and then close to ensure that the service is not affected
	//If you need to use rolling upgrades and other strategies, use multiple strategies
	Strategy []string `json:"strategy"`
}

//GroupStartTaskBody group application start operation task body
type GroupStartTaskBody struct {
	Services    []StartTaskBody `json:"services"`
	Dependences []Dependence    `json:"dependences"`
	//Group startup strategy
	//Sequential start, disorderly concurrent start
	Strategy []string `json:"strategy"`
}

// ApplyRuleTaskBody contains information for ApplyRuleTask
type ApplyRuleTaskBody struct {
	ServiceID     string            `json:"service_id"`
	EventID       string            `json:"event_id"`
	ServiceKind   string            `json:"service_kind"`
	Action        string            `json:"action"`
	Port          int               `json:"port"`
	IsInner       bool              `json:"is_inner"`
	Limit         map[string]string `json:"limit"`
}

// ApplyPluginConfigTaskBody apply plugin dynamic discover config
type ApplyPluginConfigTaskBody struct {
	ServiceID string `json:"service_id"`
	PluginID  string `json:"plugin_id"`
	EventID   string `json:"event_id"`
	//Action put delete
	Action string `json:"action"`
}

//Dependence dependency
type Dependence struct {
	CurrentServiceID string `json:"current_service_id"`
	DependServiceID  string `json:"depend_service_id"`
}

//GroupStopTaskBody group application stops operating the task body
type GroupStopTaskBody struct {
	Services    []StartTaskBody `json:"services"`
	Dependences []Dependence    `json:"dependences"`
	//Group shutdown strategy
	//Sequence relationship, unordered concurrent close
	Strategy []string `json:"strategy"`
}

// ServiceGCTaskBody holds the request body to execute service gc task.
type ServiceGCTaskBody struct {
	TenantID  string   `json:"tenant_id"`
	ServiceID string   `json:"service_id"`
	EventIDs  []string `json:"event_ids"`
}

// DeleteTenantTaskBody -
type DeleteTenantTaskBody struct {
	TenantID string `json:"tenant_id"`
}

// RefreshHPATaskBody -
type RefreshHPATaskBody struct {
	ServiceID string `json:"service_id"`
	RuleID    string `json:"rule_id"`
	EventID   string `json:"eventID"`
}

//DefaultTaskBody default operation task body
type DefaultTaskBody map[string]interface{}
