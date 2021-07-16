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

package builder

import (
	"os"
	"path"
	"strings"

	"github.com/gridworkz/kato/util/constants"
)

func init() {
	if os.Getenv("BUILD_IMAGE_REPOSTORY_DOMAIN") != "" {
		REGISTRYDOMAIN = os.Getenv("BUILD_IMAGE_REPOSTORY_DOMAIN")
	}
	if os.Getenv("BUILD_IMAGE_REPOSTORY_USER") != "" {
		REGISTRYUSER = os.Getenv("BUILD_IMAGE_REPOSTORY_USER")
	}
	if os.Getenv("BUILD_IMAGE_REPOSTORY_PASS") != "" {
		REGISTRYPASS = os.Getenv("BUILD_IMAGE_REPOSTORY_PASS")
	}
	RUNNERIMAGENAME = "/runner"
	if os.Getenv("RUNNER_IMAGE_NAME") != "" {
		RUNNERIMAGENAME = os.Getenv("RUNNER_IMAGE_NAME")
	}
	RUNNERIMAGENAME = path.Join(REGISTRYDOMAIN, RUNNERIMAGENAME)
	BUILDERIMAGENAME = "builder"
	if os.Getenv("BUILDER_IMAGE_NAME") != "" {
		BUILDERIMAGENAME = os.Getenv("BUILDER_IMAGE_NAME")
	}

	BUILDERIMAGENAME = path.Join(REGISTRYDOMAIN, BUILDERIMAGENAME)
}

// GetImageUserInfoV2 -
func GetImageUserInfoV2(domain, user, pass string) (string, string) {
	if user != "" && pass != "" {
		return user, pass
	}
	if strings.HasPrefix(domain, REGISTRYDOMAIN) {
		return REGISTRYUSER, REGISTRYPASS
	}
	return "", ""
}

//GetImageRepo -
func GetImageRepo(imageRepo string) string {
	if imageRepo == "" {
		return REGISTRYDOMAIN
	}
	return imageRepo
}

//REGISTRYDOMAIN REGISTRY_DOMAIN
var REGISTRYDOMAIN = constants.DefImageRepository

//REGISTRYUSER REGISTRY USER NAME
var REGISTRYUSER = ""

//REGISTRYPASS REGISTRY PASSWORD
var REGISTRYPASS = ""

//RUNNERIMAGENAME runner image name
var RUNNERIMAGENAME string

//BUILDERIMAGENAME builder image name
var BUILDERIMAGENAME string
