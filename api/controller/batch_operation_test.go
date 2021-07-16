// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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

package controller

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gridworkz/kato/api/handler"

	"github.com/gridworkz/kato/api/model"
	httputil "github.com/gridworkz/kato/util/http"
)

func TestBatchOperation(t *testing.T) {
	var build = model.BeatchOperationRequestStruct{}
	buildInfo := []model.BuildInfoRequestStruct{
		model.BuildInfoRequestStruct{
			BuildENVs: map[string]string{
				"MAVEN_SETTING": "java",
			},
			Action: "upgrade",
			Kind:   model.FromCodeBuildKing,
			CodeInfo: model.BuildCodeInfo{
				RepoURL:    "https://github.com/gridworkz/java-maven-demo.git",
				Branch:     "master",
				Lang:       "Java-maven",
				ServerType: "git",
			},
			ServiceID: "qwertyuiopasdfghjklzxcvbn",
		},
		model.BuildInfoRequestStruct{
			Action: "upgrade",
			Kind:   model.FromImageBuildKing,
			ImageInfo: model.BuildImageInfo{
				ImageURL: "hub.gridworkz.com/xxx/xxx:latest",
				Cmd:      "start web",
			},
			ServiceID: "qwertyuiopasdfghjklzxcvbn",
		},
	}
	startInfo := []model.StartOrStopInfoRequestStruct{
		model.StartOrStopInfoRequestStruct{
			ServiceID: "qwertyuiopasdfghjkzxcvb",
		},
	}
	upgrade := []model.UpgradeInfoRequestStruct{
		model.UpgradeInfoRequestStruct{
			ServiceID:      "qwertyuiopasdfghjkzxcvb",
			UpgradeVersion: "2345678",
		},
	}
	build.Body.BuildInfos = buildInfo
	build.Body.StartInfos = startInfo
	build.Body.StopInfos = startInfo
	build.Body.UpgradeInfos = upgrade

	build.Body.Operation = "stop"
	out, _ := json.MarshalIndent(build.Body, "", "\t")
	fmt.Print(string(out))

	result := handler.BatchOperationResult{
		BatchResult: []handler.OperationResult{
			handler.OperationResult{
				ServiceID:     "qwertyuiopasdfghjkzxcvb",
				Operation:     "build",
				EventID:       "wertyuiodfghjcvbnm",
				Status:        "success",
				ErrMsg:        "",
				DeployVersion: "1234567890",
			},
		},
	}
	rebody := httputil.ResponseBody{
		Bean: result,
	}
	outre, _ := json.MarshalIndent(rebody, "", "\t")
	fmt.Print(string(outre))
}
