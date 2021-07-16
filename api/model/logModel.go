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
	eventdb "github.com/gridworkz/kato/eventlog/db"
)

//LogData
type LogData struct {
	num int
	msg string
}

//MessageData message data - obtain the operation log of the specified operation
type MessageData struct {
	Message  string `json:"message"`
	Time     string `json:"time"`
	Unixtime int64  `json:"utime"`
}

//DataLog - obtain the operation log of the specified operation
type DataLog struct {
	Status string
	Data   eventdb.MessageDataList
}

//LogByLevelStruct GetLogByLevelStruct
//swagger:parameters logByAction
type LogByLevelStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
	// in: body
	Body struct {
		// log level - info/debug/error
		// in: body
		// required: true
		Level string `json:"level" validate:"level|required"`
		// eventID
		// in: body
		// required: true
		EventID string `json:"event_id" validate:"event_id|required"`
	}
}

//TenantLogByLevelStruct GetTenantLogByLevelStruct
//swagger:parameters tenantLogByAction
type TenantLogByLevelStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: body
	Body struct {
		// log level : info/debug/error
		// in: body
		// required: true
		Level string `json:"level" validate:"level|required"`
		// eventID
		// in: body
		// required: true
		EventID string `json:"event_id" validate:"event_id|required"`
	}
}

//LogSocketStruct LogSocketStruct
//swagger:parameters logSocket logList
type LogSocketStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
}

//LogFileStruct LogFileStruct
//swagger:parameters logFile
type LogFileStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
	// in: path
	// required: true
	FileName string `json:"file_name" validate:"file_name"`
}

//MsgStruct msg struct in eventlog_message
type MsgStruct struct {
	EventID string `json:"event_id"`
	Step    string `json:"step"`
	Message string `json:"message"`
	Level   string `json:"level"`
	Time    string `json:"time"`
}

//LastLinesStruct LastLinesStruct
//swagger:parameters lastLinesLogs
type LastLinesStruct struct {
	// in: path
	// required: true
	TenantName string `json:"tenant_name" validate:"tenant_name"`
	// in: path
	// required: true
	ServiceAlias string `json:"service_alias" validate:"service_alias"`
	// in: body
	Body struct {
		// rows
		// in: body
		// required: true
		Lines int `json:"lines" validate:"lines"`
	}
}
