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
	"fmt"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/sirupsen/logrus"
)

//VersionInfo
type VersionInfo struct {
	Model
	BuildVersion string `gorm:"column:build_version;size:40" json:"build_version"` //unique
	EventID      string `gorm:"column:event_id;size:40" json:"event_id"`
	ServiceID    string `gorm:"column:service_id;size:40" json:"service_id"`
	Kind         string `gorm:"column:kind;size:40" json:"kind"` //kind
	//DeliveredType app version delivered type
	// image: this is a docker image
	// slug: this is a source code tar file
	DeliveredType string `gorm:"column:delivered_type;size:40" json:"delivered_type"`  //kind
	DeliveredPath string `gorm:"column:delivered_path;size:250" json:"delivered_path"` //deliverable path
	ImageName     string `gorm:"column:image_name;size:250" json:"image_name"`         //run image name
	Cmd           string `gorm:"column:cmd;size:2048" json:"cmd"`                      //start command
	RepoURL       string `gorm:"column:repo_url;size:2047" json:"repo_url"`
	CodeVersion   string `gorm:"column:code_version;size:40" json:"code_version"`
	CodeBranch    string `gorm:"column:code_branch;size:40" json:"code_branch"`
	CommitMsg     string `gorm:"column:code_commit_msg;size:1024" json:"code_commit_msg"`
	Author        string `gorm:"column:code_commit_author;size:40" json:"code_commit_author"`
	//FinalStatus app version status
	// success: version available
	// failure: build failure
	// lost: there is nothing delivered
	FinalStatus string    `gorm:"column:final_status;size:40" json:"final_status"`
	FinishTime  time.Time `gorm:"column:finish_time;" json:"finish_time"`
	PlanVersion string  `gorm:"column:plan_version;size:250" json:"plan_version"`
}

//TableName
func (t *VersionInfo) TableName() string {
	return "tenant_service_version"
}

//CreateShareImage
func (t *VersionInfo) CreateShareImage(hubURL, namespace, appVersion string) (string, error) {
	_, err := reference.ParseAnyReference(t.DeliveredPath)
	if err != nil {
		logrus.Errorf("reference image error: %s", err.Error())
		return "", err
	}
	image := ParseImage(t.DeliveredPath)
	if hubURL != "" {
		image.Host = hubURL
	}
	if namespace != "" {
		image.Namespace = namespace
	}
	image.Name = fmt.Sprintf("%s:%s", t.ServiceID, t.BuildVersion)
	return image.String(), nil
}
