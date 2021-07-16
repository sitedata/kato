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

//TableName
func (t *RegionProcotols) TableName() string {
	return "region_protocols"
}

//RegionProcotols
type RegionProcotols struct {
	Model
	ProtocolGroup string `gorm:"column:protocol_group;size:32;" json:"protocol_group"`
	ProtocolChild string `gorm:"column:protocol_child;size:32;" json:"protocol_child"`
	APIVersion    string `gorm:"column:api_version;size:8" json:"api_version"`
	IsSupport     bool   `gorm:"column:is_support;default:false" json:"is_support"`
}

//STREAMGROUP
var STREAMGROUP = "stream"

//HTTPGROUP
var HTTPGROUP = "http"

//MYSQLPROTOCOL
var MYSQLPROTOCOL = "mysql"

//UDPPROTOCOL
var UDPPROTOCOL = "udp"

//TCPPROTOCOL
var TCPPROTOCOL = "tcp"

//V2VERSION region version
var V2VERSION = "v2"
