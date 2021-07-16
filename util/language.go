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

package util

import "os"

var translationMetadata = map[string]string{
	"write console level metadata success":           "write console level metadata success",
	"write region level metadata success":            "write region level metadata success",
	"Asynchronous tasks are sent successfully":       "Asynchronous tasks are sent successfully",
	"create ftp client error":                        "create ftp client error",
	"push slug file to local dir error":              "push slug file to local dir error",
	"push slug file to sftp server error":            "push slug file to sftp server error",
	"down slug file from sftp server error":          "down slug file from sftp server error",
	"save image to local dir error":                  "save image to local dir error",
	"save image to hub error":                        "save image to hub error",
	"Please try again or contact customer service":   "Please try again or contact customer service",
	"unzip metadata file error":                      "unzip metadata file error",
	"start service error":                            "start service error",
	"start service timeout":                          "start service timeout",
	"stop service error":                             "stop service error",
	"stop service timeout":                           "stop service timeout",
	"(restart)stop service error":                    "(restart)stop service error",
	"(restart)Application model init create failure": "(restart)Application model init create failure",
	"horizontal scaling service error":               "horizontal scaling service error",
	"horizontal scaling service timeout":             "horizontal scaling service timeout",
	"upgrade service error":                          "upgrade service error",
	"upgrade service timeout":                        "upgrade service timeout",
	"Check for log location code errors":             "Check for log location code errors",
	"Check for log location imgae source errors":     "Check for log location imgae source errors",
	"create share image task error":                  "create share image task error",
	"get rbd-repo ip failure":                        "get rbd-repo ip failure",
}

//Translation
func Translation(english string) string {
	if chinese, ok := translationMetadata[english]; ok {
		if os.Getenv("KATO_LANG") == "en" {
			return english
		}
		return chinese
	}
	return english
}
