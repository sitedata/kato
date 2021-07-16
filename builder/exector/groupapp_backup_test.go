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

package exector

import (
	"fmt"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	"strings"
	"testing"
)

func TestUploadPkg(t *testing.T) {
	b := &BackupAPPNew{
		SourceDir: "/tmp/groupbackup/0d65c6608729438aad0a94f6317c80d0_20191024180024.zip",
		Mode:      "full-online",
	}
	b.S3Config.Provider = "AliyunOSS"
	b.S3Config.Endpoint = "dummy"
	b.S3Config.AccessKey = "dummy"
	b.S3Config.SecretKey = "dummy"
	b.S3Config.BucketName = "dummy"

	if err := b.uploadPkg(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUploadPkg2(t *testing.T) {
	b := &BackupAPPRestore{}
	b.S3Config.Provider = "alioss"
	b.S3Config.Endpoint = "dummy"
	b.S3Config.AccessKey = "dummy"
	b.S3Config.SecretKey = "dummy"
	b.S3Config.BucketName = "hrhtest"

	cacheDir := fmt.Sprintf("/tmp/cache/tmp/%s/%s", "c6b05a2a6d664fda83dab8d3bcf1a941", util.NewUUID())
	if err := util.CheckAndCreateDir(cacheDir); err != nil {
		t.Errorf("create cache dir error %s", err.Error())
	}
	b.cacheDir = cacheDir

	sourceDir := "/tmp/groupbackup/c6b05a2a6d664fda83dab8d3bcf1a941_20191024185643.zip"
	if err := b.downloadFromS3(sourceDir); err != nil {
		t.Error(err)
	}
}

func TestBackupServiceVolume(t *testing.T) {
	volume := dbmodel.TenantServiceVolume{}
	sourceDir := ""
	serviceID := ""
	dstDir := fmt.Sprintf("%s/data_%s/%s.zip", sourceDir, serviceID, strings.Replace(volume.VolumeName, "/", "", -1))
	hostPath := volume.HostPath
	if hostPath != "" && !util.DirIsEmpty(hostPath) {
		if err := util.Zip(hostPath, dstDir); err != nil {
			t.Fatalf("backup service(%s) volume(%s) data error.%s", serviceID, volume.VolumeName, err.Error())
		}
	}
}
