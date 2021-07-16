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
	"io"
	"testing"
	"time"

	"github.com/gridworkz/kato/event"
)

func TestGitClone(t *testing.T) {
	start := time.Now()
	csi := CodeSourceInfo{
		RepositoryURL: "git@gitee.com:zhoujunhaogridworkz/webhook_test.git",
		Branch:        "master",
	}
	res, err := GitClone(csi, "/tmp/katodoc3", event.GetTestLogger(), 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Take %d ms", time.Now().Unix()-start.Unix())
	commit, err := GetLastCommit(res)
	t.Logf("%+v %+v", commit, err)
}
func TestGitCloneByTag(t *testing.T) {
	start := time.Now()
	csi := CodeSourceInfo{
		RepositoryURL: "https://github.com/gridworkz/kato-ui.git",
		Branch:        "master",
	}
	res, err := GitClone(csi, "/tmp/katodoc4", event.GetTestLogger(), 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Take %d ms", time.Now().Unix()-start.Unix())
	commit, err := GetLastCommit(res)
	t.Logf("%+v %+v", commit, err)
}

func TestGitPull(t *testing.T) {
	csi := CodeSourceInfo{
		RepositoryURL: "git@gitee.com:zhoujunhaogridworkz/webhook_test.git",
		Branch:        "master2",
	}
	res, err := GitPull(csi, "/tmp/master2", event.GetTestLogger(), 1)
	if err != nil {
		t.Fatal(err)
	}
	commit, err := GetLastCommit(res)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", commit)
}

func TestGitPullOrClone(t *testing.T) {
	csi := CodeSourceInfo{
		RepositoryURL: "git@gitee.com:zhoujunhaogridworkz/webhook_test.git",
	}
	res, err := GitCloneOrPull(csi, "/tmp/gridworkzweb2", event.GetTestLogger(), 1)
	if err != nil {
		t.Fatal(err)
	}
	//get last commit
	commit, err := GetLastCommit(res)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	t.Logf("%+v", commit)
}

func TestGetCodeCacheDir(t *testing.T) {
	csi := CodeSourceInfo{
		RepositoryURL: "git@121.196.222.148:summersoft/yycx_push.git",
		Branch:        "test",
	}
	t.Log(csi.GetCodeSourceDir())
}

func TestGetShowURL(t *testing.T) {
	t.Log(getShowURL("https://zsl1526:79890ffc74014b34b49040d42b95d5af@github.com:9090/zsl1549/python-demo.git"))
}
