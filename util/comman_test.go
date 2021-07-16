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

package util

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestOpenOrCreateFile(t *testing.T) {
	file, err := OpenOrCreateFile("./test.log")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
}

func TestDeweight(t *testing.T) {
	data := []string{"asd", "asd", "12", "12"}
	Deweight(&data)
	t.Log(data)
}

func TestGetDirSize(t *testing.T) {
	t.Log(GetDirSize("/go"))
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
}

func TestGetDirSizeByCmd(t *testing.T) {
	t.Log(GetDirSizeByCmd("/go"))
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
}

func TestZip(t *testing.T) {
	if err := Zip("/tmp/cache", "/tmp/cache.zip"); err != nil {
		t.Fatal(err)
	}
}

func TestUnzip(t *testing.T) {
	if err := Unzip("/tmp/cache.zip", "/tmp/cache0"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateVersionByTime(t *testing.T) {
	if re := CreateVersionByTime(); re != "" {
		t.Log(re)
	}
}

func TestGetDirList(t *testing.T) {
	list, err := GetDirList("/tmp", 2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}

func TestMergeDir(t *testing.T) {
	t.Log(filepath.Dir("/tmp/cache/asdasd"))
	if err := MergeDir("/tmp/ctr-944254844/", "/tmp/cache"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateHostID(t *testing.T) {
	uid, err := CreateHostID()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uid)
}

func TestGetCurrentDir(t *testing.T) {
	t.Log(GetCurrentDir())
}

func TestCopyFile(t *testing.T) {
	if err := CopyFile("/tmp/test2.zip", "/tmp/test4.zip"); err != nil {
		t.Fatal(err)
	}
}

func TestParseVariable(t *testing.T) {
	configs := make(map[string]string, 0)
	result := ParseVariable("sada${XXX:aaa}dasd${XXX:aaa} ${YYY:aaa} ASDASD ${ZZZ:aaa}", configs)
	t.Log(result)

	t.Log(ParseVariable("sada${XXX:aaa}dasd${XXX:aaa} ${YYY:aaa} ASDASD ${ZZZ:aaa}", map[string]string{
		"XXX": "123DDD",
		"ZZZ": ",.,.,.,.",
	}))
}

func TestTimeFormat(t *testing.T) {
	tt := "2019-08-24 11:11:30.165753932 +0800 CST m=+55557.682499470"
	timeF, err := time.Parse(time.RFC3339, strings.Replace(tt[0:19]+"+08:00", " ", "T", 1))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(timeF.Format(time.RFC3339))
}
