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

package sources

import (
	"fmt"
	"path"
	"strings"

	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

//RepostoryBuildInfo
type RepostoryBuildInfo struct {
	RepostoryURL     string
	RepostoryURLType string
	BuildBranch      string
	BuildPath        string
	CodeHome         string
	ep               *transport.Endpoint
}

//GetCodeHome
func (r *RepostoryBuildInfo) GetCodeHome() string {
	if r.RepostoryURLType == "svn" {
		if ok, _ := util.FileExists(path.Join(r.CodeHome, "trunk")); ok && r.BuildBranch == "trunk" {
			return path.Join(r.CodeHome, "trunk")
		}
		if r.BuildBranch != "" && r.BuildBranch != "trunk" {
			if strings.HasPrefix(r.BuildBranch, "tag:") {
				codepath := path.Join(r.CodeHome, "tags", r.BuildBranch[4:])
				if ok, _ := util.FileExists(codepath); ok {
					return codepath
				}
				codepath = path.Join(r.CodeHome, "Tags", r.BuildBranch[4:])
				if ok, _ := util.FileExists(codepath); ok {
					return codepath
				}
			}
			codepath := path.Join(r.CodeHome, "branches", r.BuildBranch)
			if ok, _ := util.FileExists(codepath); ok {
				return codepath
			}
			codepath = path.Join(r.CodeHome, "Branches", r.BuildBranch)
			if ok, _ := util.FileExists(codepath); ok {
				return codepath
			}
		}
	}
	return r.CodeHome
}

//GetCodeBuildAbsPath
func (r *RepostoryBuildInfo) GetCodeBuildAbsPath() string {
	return path.Join(r.GetCodeHome(), r.BuildPath)
}

//GetCodeBuildPath
func (r *RepostoryBuildInfo) GetCodeBuildPath() string {
	return r.BuildPath
}

//GetProtocol
func (r *RepostoryBuildInfo) GetProtocol() string {
	if r.ep != nil {
		if r.ep.Protocol == "" {
			return "ssh"
		}
		return r.ep.Protocol
	}
	return ""
}

//CreateRepostoryBuildInfo
//repoType git or svn
func CreateRepostoryBuildInfo(repoURL, repoType, branch, tenantID string, ServiceID string) (*RepostoryBuildInfo, error) {
	// repoURL= github.com/gridworkz/xxx.git?dir=home
	ep, err := transport.NewEndpoint(repoURL)
	if err != nil {
		return nil, err
	}
	rbi := &RepostoryBuildInfo{
		ep:               ep,
		RepostoryURL:     repoURL,
		RepostoryURLType: repoType,
		BuildBranch:      branch,
	}
	index := strings.Index(repoURL, "?dir=")
	if index > -1 && len(repoURL) > index+5 {
		fmt.Println(repoURL[index+5:], repoURL[:index])
		rbi.BuildPath = repoURL[index+5:]
		rbi.CodeHome = GetCodeSourceDir(repoURL[:index], branch, tenantID, ServiceID)
		rbi.RepostoryURL = repoURL[:index]
	}
	rbi.CodeHome = GetCodeSourceDir(repoURL, branch, tenantID, ServiceID)
	logrus.Infof("cache code dir is %s for service %s", rbi.CodeHome, ServiceID)
	return rbi, nil
}
