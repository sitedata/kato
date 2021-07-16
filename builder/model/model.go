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

//BuildPluginTaskBody
type BuildPluginTaskBody struct {
	VersionID     string `json:"version_id"`
	TenantID      string `json:"tenant_id"`
	PluginID      string `json:"plugin_id"`
	Operator      string `json:"operator"`
	Repo          string `json:"repo"`
	GitURL        string `json:"git_url"`
	GitUsername   string `json:"git_username"`
	GitPassword   string `json:"git_password"`
	ImageURL      string `json:"image_url"`
	EventID       string `json:"event_id"`
	DeployVersion string `json:"deploy_version"`
	Kind          string `json:"kind"`
	PluginCMD     string `json:"plugin_cmd"`
	PluginCPU     int    `json:"plugin_cpu"`
	PluginMemory  int    `json:"plugin_memory"`
	ImageInfo     struct {
		HubURL      string `json:"hub_url"`
		HubUser     string `json:"hub_user"`
		HubPassword string `json:"hub_password"`
		Namespace   string `json:"namespace"`
		IsTrust     bool   `json:"is_trust,omitempty"`
	} `json:"image_info,omitempty"`
}

//BuildPluginVersion
type BuildPluginVersion struct {
	SourceImage string `json:"source_image"`
	InnerImage  string `json:"inner_image"`
	CreateTime  string `json:"create_time"`
	Repo        string `json:"repo"`
}

//CodeCheckResult
type CodeCheckResult struct {
	ServiceID    string `json:"service_id"`
	Condition    string `json:"condition"`
	CheckType    string `json:"check_type"`
	GitURL       string `json:"git_url"`
	CodeVersion  string `json:"code_version"`
	GitProjectId string `json:"git_project_id"`
	CodeFrom     string `json:"code_from"`
	URLRepos     string `json:"url_repos"`

	DockerFileReady bool              `json:"docker_file_ready,omitempty"`
	InnerPort       string            `json:"inner_port,omitempty"`
	VolumeMountPath string            `json:"volume_mount_path,omitempty"`
	BuildImageName  string            `json:"image,omitempty"`
	PortList        map[string]string `json:"port_list,omitempty"`
	VolumeList      []string          `json:"volume_list,omitempty"`

	//DFR          *DockerFileResult `json:"dockerfile,omitempty"`
}

//ImageName
type ImageName struct {
	Host      string `json:"host"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Tag       string `json:"tag"`
}
