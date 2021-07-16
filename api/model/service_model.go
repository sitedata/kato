// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the MIT License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-MIT usage of Kato,
// one or multiple Commercial Licenses authorized by Gridworkz Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// MIT License for more details.

// You should have received a copy of the MIT License
// along with this program. If not, see <http://opensource.org/licenses/MIT>.

package model

import "time"

// BuildListRespVO is the response value object for build-list api.
type BuildListRespVO struct {
	DeployVersion string      `json:"deploy_version"`
	List          interface{} `json:"list"`
}

// BuildVersion -
type BuildVersion struct {
	BuildVersion string `json:"build_version"` //only
	EventID      string `json:"event_id"`
	ServiceID    string `json:"service_id"`
	Kind         string `json:"kind"` //kind
	//DeliveredType app version delivered type
	//image: this is a docker image
	//slug: this is a source code tar file
	DeliveredType string `json:"delivered_type"` //kind
	DeliveredPath string `json:"delivered_path"` //deliverable path
	Cmd           string `json:"cmd"`            //start command
	RepoURL       string `json:"repo_url"`       //source image name or source code url

	CodeBranch  string `json:"code_branch"`
	CodeVersion string `json:"code_version"`
	CommitMsg   string `json:"code_commit_msg"`
	Author      string `json:"code_commit_author"`

	ImageName   string `json:"image_name"` // runtime image name
	ImageRepo   string `json:"image_repo"`
	ImageDomain string `json:"image_domain"`
	ImageTag    string `json:"image_tag"`

	//FinalStatus app version status
	//success: version available
	//failure: build failure
	//lost: there is no delivered
	CreateTime  string    `json:"create_time"`
	FinalStatus string    `json:"final_status"`
	FinishTime  time.Time `json:"finish_time"`
	PlanVersion string    `json:"plan_version"`
}
