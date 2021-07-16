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
	"testing"

	"github.com/gridworkz/kato/event"
)

func TestPushFile(t *testing.T) {
	sftpClient, err := NewSFTPClient("admin", "9bc067dc", "47.92.168.60", "20012")
	if err != nil {
		t.Fatal(err)
	}
	if err := sftpClient.PushFile("/tmp/src.tgz", "/upload/team/servicekey/gridworkz.tgz", event.GetTestLogger()); err != nil {
		t.Fatal(err)
	}
}

func TestDownloadFile(t *testing.T) {
	sftpClient, err := NewSFTPClient("foo", "pass", "22.gr6ac909.0mi9zp2q.lfsdo.gridworkz.org", "20004")
	if err != nil {
		t.Fatal(err)
	}
	if err := sftpClient.DownloadFile("/upload/team/servicekey/gridworkz.tgz", "./team/servicekey/gridworkz.tgz", event.GetTestLogger()); err != nil {
		t.Fatal(err)
	}
}
