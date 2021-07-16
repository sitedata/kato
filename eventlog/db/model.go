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

package db

//EventLogMessage - event log entity
type EventLogMessage struct {
	EventID string `json:"event_id"`
	Step    string `json:"step"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Level   string `json:"level"`
	Time    string `json:"time"`
	Content []byte `json:"-"`
	//monitor message usage
	MonitorData []byte `json:"monitorData,omitempty"`
}

type MonitorData struct {
	InstanceID   string
	ServiceSize  int
	LogSizePeerM int64
}
type ClusterMessageType string

const (
	//EventMessage - operation log sharing
	EventMessage ClusterMessageType = "event_log"
	//ServiceMonitorMessage - business monitoring data message
	ServiceMonitorMessage ClusterMessageType = "monitor_message"
	//ServiceNewMonitorMessage - new business monitoring data message
	ServiceNewMonitorMessage ClusterMessageType = "new_monitor_message"
	//MonitorMessage - node monitoring data
	MonitorMessage ClusterMessageType = "monitor"
)

type ClusterMessage struct {
	Data []byte
	Mode ClusterMessageType
}
