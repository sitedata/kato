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
func (t *CodeCheckResult) TableName() string {
	return "tenant_services_codecheck"
}

//CodeCheckResult - codecheck result struct
type CodeCheckResult struct {
	Model
	ServiceID       string `gorm:"column:service_id;size:70"`
	Condition       string `gorm:"column:condition"`
	Language        string `gorm:"column:language"`
	CheckType       string `gorm:"column:check_type"`
	GitURL          string `gorm:"column:git_url"`
	CodeVersion     string `gorm:"column:code_version"`
	GitProjectId    string `gorm:"column:git_project_id"`
	CodeFrom        string `gorm:"column:code_from"`
	URLRepos        string `gorm:"column:url_repos"`
	DockerFileReady bool   `gorm:"column:docker_file_ready"`
	InnerPort       string `gorm:"column:inner_port"`
	VolumeMountPath string `gorm:"column:volume_mount_path"`
	BuildImageName  string `gorm:"column:image"`
	PortList        string `gorm:"column:port_list"`
	VolumeList      string `gorm:"column:volume_list"`
}
