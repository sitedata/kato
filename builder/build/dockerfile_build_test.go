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

package build

import (
	"testing"
)

func TestGetARGs(t *testing.T) {
	buildEnvs := make(map[string]string)
	buildEnvs["ARG_TEST"] = "abcdefg"
	buildEnvs["PROC_ENV"] = "{\"procfile\": \"\", \"dependencies\": {}, \"language\": \"dockerfile\", \"runtimes\": \"\"}"

	args := GetARGs(buildEnvs)
	if v := buildEnvs["ARG_TEST"]; *args["TEST"] != v {
		t.Errorf("Expected %s for arg[\"%s\"], but returned %s", buildEnvs["ARG_TEST"], "ARG_TEST", *args["TEST"])
	}
	if procEnv := args["PROC_ENV"]; procEnv != nil {
		t.Errorf("Expected nil for  args[\"PROC_ENV\"], but returned %v", procEnv)
	}
}
