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

package parser

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/parser/types"
	"github.com/gridworkz/kato/builder/sources"
)

func TestParseDockerfileInfo(t *testing.T) {
	parse := &SourceCodeParse{
		source:  "source",
		ports:   make(map[int]*types.Port),
		volumes: make(map[string]*types.Volume),
		envs:    make(map[string]*types.Env),
		logger:  nil,
		image:   ParseImageName(builder.RUNNERIMAGENAME),
		args:    []string{"start", "web"},
	}
	parse.parseDockerfileInfo("./Dockerfile")
	fmt.Println(parse.GetServiceInfo())
}

//ServiceCheckResult
type ServiceCheckResult struct {
	//Detection status: Success/Failure
	CheckStatus string         `json:"check_status"`
	ErrorInfos  ParseErrorList `json:"error_infos"`
	ServiceInfo []ServiceInfo  `json:"service_info"`
}

func TestSourceCode(t *testing.T) {
	sc := sources.CodeSourceInfo{
		ServerType:    "",
		RepositoryURL: "https://github.com/gridworkz/fserver.git",
		Branch:        "master",
	}
	b, _ := json.Marshal(sc)
	p := CreateSourceCodeParse(string(b), nil)
	err := p.Parse()
	if err != nil && err.IsFatalError() {
		t.Fatal(err)
	}
	re := ServiceCheckResult{
		CheckStatus: "Failure",
		ErrorInfos:  err,
		ServiceInfo: p.GetServiceInfo(),
	}
	body, _ := json.Marshal(re)
	fmt.Printf("%s \n", string(body))
}
