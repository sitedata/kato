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

package sources

import (
	"fmt"

	"github.com/gridworkz/kato/builder/sources/registry"

	"github.com/docker/distribution/reference"
	"github.com/sirupsen/logrus"
)

//GetTagFromNamedRef get image tag by name
func GetTagFromNamedRef(ref reference.Named) string {
	if digested, ok := ref.(reference.Digested); ok {
		return digested.Digest().String()
	}
	ref = reference.TagNameOnly(ref)
	if tagged, ok := ref.(reference.Tagged); ok {
		return tagged.Tag()
	}
	return ""
}

//ImageExist check image exist
func ImageExist(imageName, user, password string) (bool, error) {
	ref, err := reference.ParseAnyReference(imageName)
	if err != nil {
		logrus.Errorf("reference image error: %s", err.Error())
		return false, err
	}
	name, err := reference.ParseNamed(ref.String())
	if err != nil {
		logrus.Errorf("reference parse image name error: %s", err.Error())
		return false, err
	}
	domain := reference.Domain(name)
	if domain == "docker.io" {
		domain = "registry-1.docker.io"
	}
	retry := 2
	var rerr error
	for retry > 0 {
		retry--
		reg, err := registry.New(domain, user, password)
		if err != nil {
			logrus.Debugf("new registry client failure %s", err.Error())
			reg, err = registry.NewInsecure(domain, user, password)
			if err != nil {
				logrus.Debugf("new insecure registry client failure %s", err.Error())
				reg, err = registry.NewInsecure("http://"+domain, user, password)
				if err != nil {
					logrus.Errorf("new insecure registry http or https client all failure %s", err.Error())
					rerr = err
					continue
				}
			}
		}
		tag := GetTagFromNamedRef(name)
		if err := reg.CheckManifest(reference.Path(name), tag); err != nil {
			rerr = fmt.Errorf("[ImageExist] check manifest v2: %v", err)
			continue
		}
		return true, nil
	}
	return false, rerr
}
