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

	"github.com/pquerna/ffjson/ffjson"
)

//Shell
type Shell struct {
	Cmd []string `json:"cmd"`
}

//TaskTemp
type TaskTemp struct {
	Name       string            `json:"name" validate:"name|required"`
	ID         string            `json:"id" validate:"id|uuid"`
	Shell      Shell             `json:"shell"`
	Envs       map[string]string `json:"envs,omitempty"`
	Input      string            `json:"input,omitempty"`
	Args       []string          `json:"args,omitempty"`
	Depends    []DependStrategy  `json:"depends,omitempty"`
	Timeout    int               `json:"timeout" validate:"timeout|required|numeric"`
	CreateTime time.Time         `json:"create_time"`
	Labels     map[string]string `json:"labels,omitempty"`
}

//DependStrategy
type DependStrategy struct {
	DependTaskID      string `json:"depend_task_id"`
	DetermineStrategy string `json:"strategy"`
}

//AtLeastOnceStrategy
var AtLeastOnceStrategy = "AtLeastOnce"

//SameNodeStrategy
var SameNodeStrategy = "SameNode"

func (t TaskTemp) String() string {
	res, _ := ffjson.Marshal(&t)
	return string(res)
}

//Task
type Task struct {
	Name    string    `json:"name" validate:"name|required"`
	ID      string    `json:"id" validate:"id|uuid"`
	TempID  string    `json:"temp_id,omitempty" validate:"temp_id|uuid"`
	Temp    *TaskTemp `json:"temp,omitempty"`
	GroupID string    `json:"group_id,omitempty"`
	//Node of execution
	Nodes []string `json:"nodes"`
	//Execution time definition
	//For example, execute every 30 minutes: @every 30m
	Timer   string `json:"timer"`
	TimeOut int64  `json:"time_out"`
	//The number of failed retry attempts
	//The default is 0, do not try again
	Retry int `json:"retry"`
	//Failed task execution retry interval
	//Unit second, if not greater than 0, try again immediately
	Interval int `json:"interval"`
	//ExecCount
	ExecCount int `json:"exec_count"`
	//Execution status of each execution node
	Status       map[string]TaskStatus `json:"status,omitempty"`
	Scheduler    Scheduler             `json:"scheduler"`
	CreateTime   time.Time             `json:"create_time"`
	StartTime    time.Time             `json:"start_time"`
	CompleteTime time.Time             `json:"complete_time"`
	ResultPath   string                `json:"result_path"`
	EventID      string                `json:"event_id"`
	RunMode      string                `json:"run_mode"`
	OutPut       []*TaskOutPut         `json:"out_put"`
}

func (t Task) String() string {
	res, _ := ffjson.Marshal(&t)
	return string(res)
}

//Decode
func (t *Task) Decode(data []byte) error {
	return ffjson.Unmarshal(data, t)
}

//UpdataOutPut
func (t *Task) UpdataOutPut(output TaskOutPut) {
	updateIndex := -1
	for i, oldOut := range t.OutPut {
		if oldOut.NodeID == output.NodeID {
			updateIndex = i
			break
		}
	}
	if updateIndex != -1 {
		t.OutPut[updateIndex] = &output
		return
	}
	t.OutPut = append(t.OutPut, &output)
}

//CanBeDelete
func (t Task) CanBeDelete() bool {
	if t.Status == nil || len(t.Status) == 0 {
		return true
	}
	for _, v := range t.Status {
		if v.Status != "create" {
			return false
		}
	}
	return true
}

//Scheduler
type Scheduler struct {
	Mode   string                     `json:"mode"` //Immediate scheduling (Intime), trigger scheduling (Passive)
	Status map[string]SchedulerStatus `json:"status"`
}

//SchedulerStatus
type SchedulerStatus struct {
	Status          string    `json:"status"`
	Message         string    `json:"message"`
	SchedulerTime   time.Time `json:"scheduler_time"`   //Scheduling time
	SchedulerMaster string    `json:"scheduler_master"` //Scheduled management node
}

//TaskOutPut
type TaskOutPut struct {
	NodeID string            `json:"node_id"`
	JobID  string            `json:"job_id"`
	Global map[string]string `json:"global"`
	Inner  map[string]string `json:"inner"`
	//Return data type, test result type (check), execution installation type (install), common type (common)
	Type       string             `json:"type"`
	Status     []TaskOutPutStatus `json:"status"`
	ExecStatus string             `json:"exec_status"`
	Body       string             `json:"body"`
}

//ParseTaskOutPut json parse
func ParseTaskOutPut(body string) (t TaskOutPut, err error) {
	t.Body = body
	err = ffjson.Unmarshal([]byte(body), &t)
	return
}

//TaskOutPutStatus
type TaskOutPutStatus struct {
	Name string `json:"name"`
	//Node attributes
	ConditionType string `json:"condition_type"`
	//Node attribute value
	ConditionStatus string   `json:"condition_status"`
	NextTask        []string `json:"next_tasks,omitempty"`
	NextGroups      []string `json:"next_groups,omitempty"`
}

//TaskStatus
type TaskStatus struct {
	JobID        string    `json:"job_id"`
	Status       string    `json:"status"` //Execution status，create init exec complete timeout
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	TakeTime     int       `json:"take_time"`
	CompleStatus string    `json:"comple_status"`
	//Script exit code
	ShellCode int    `json:"shell_code"`
	Message   string `json:"message,omitempty"`
}

//TaskGroup
type TaskGroup struct {
	Name       string           `json:"name" validate:"name|required"`
	ID         string           `json:"id" validate:"id|uuid"`
	Tasks      []*Task          `json:"tasks"`
	CreateTime time.Time        `json:"create_time"`
	Status     *TaskGroupStatus `json:"status"`
}

func (t TaskGroup) String() string {
	res, _ := ffjson.Marshal(&t)
	return string(res)
}

//CanBeDelete
func (t TaskGroup) CanBeDelete() bool {
	if t.Status == nil || len(t.Status.TaskStatus) == 0 {
		return true
	}
	for _, v := range t.Status.TaskStatus {
		if v.Status != "create" {
			return false
		}
	}
	return true
}

//TaskGroupStatus
type TaskGroupStatus struct {
	TaskStatus map[string]TaskStatus `json:"task_status"`
	InitTime   time.Time             `json:"init_time"`
	StartTime  time.Time             `json:"start_time"`
	EndTime    time.Time             `json:"end_time"`
	Status     string                `json:"status"` //create init exec complete timeout
}
