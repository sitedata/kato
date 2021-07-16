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
	"testing"

	"github.com/gridworkz/kato/event"
)

func TestSvnCheckout(t *testing.T) {
	client := NewClient(CodeSourceInfo{
		Branch:        "docker",
		RepositoryURL: "svn://139.196.72.60:21495/runoob01/src",
	}, "/tmp/svn/docker", event.GetTestLogger())

	info, err := client.Checkout()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*info)
}
func TestSvnCheckoutBranch(t *testing.T) {
	client := NewClient(CodeSourceInfo{
		Branch:        "0.1/app",
		RepositoryURL: "https://github.com/gridworkz-apps/Cachet.git",
	}, "/tmp/svn/Cachet/branches/0.1/app", event.GetTestLogger())

	info, err := client.Checkout()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*info)
}

func TestSvnCheckoutTag(t *testing.T) {
	client := NewClient(CodeSourceInfo{
		Branch:        "tag:v2.3.6",
		RepositoryURL: "https://github.com/gridworkz-apps/Cachet.git",
	}, "/tmp/svn/Cachet/tags/v2.3.6", event.GetTestLogger())

	info, err := client.Checkout()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*info)
}

func TestUpdateOrCheckout(t *testing.T) {
	client := NewClient(CodeSourceInfo{
		Branch:        "trunk",
		RepositoryURL: "svn://ali-sh-s1.gridworkz.net:21097/testrepo",
	}, "/tmp/svn/testrepo", event.GetTestLogger())
	info, err := client.UpdateOrCheckout("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*info)
}
