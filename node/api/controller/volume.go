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

package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"

	httputil "github.com/gridworkz/kato/util/http"
)

//CreateLocalVolume
func CreateLocalVolume(w http.ResponseWriter, r *http.Request) {
	var requestopt = make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&requestopt); err != nil {
		w.WriteHeader(400)
		return
	}
	tenantID := requestopt["tenant_id"]
	serviceID := requestopt["service_id"]
	pvcName := requestopt["pvcname"]
	var volumeHostPath = ""
	localPath := os.Getenv("LOCAL_DATA_PATH")
	if runtime.GOOS == "windows" {
		if localPath == "" {
			localPath = `c:\`
		}
	} else {
		if localPath == "" {
			localPath = "/grlocaldata"
		}
	}
	volumeHostPath = path.Join(localPath, "tenant", tenantID, "service", serviceID, pvcName)
	volumePath, volumeok := requestopt["volume_path"]
	podName, podok := requestopt["pod_name"]
	if volumeok && podok {
		volumeHostPath = path.Join(localPath, "tenant", tenantID, "service", serviceID, volumePath, podName)
	}
	if err := util.CheckAndCreateDirByMode(volumeHostPath, 0777); err != nil {
		logrus.Errorf("check and create dir %s error %s", volumeHostPath, err.Error())
		w.WriteHeader(500)
		return
	}
	httputil.ReturnSuccess(r, w, map[string]string{"path": volumeHostPath})
}

// DeleteLocalVolume
func DeleteLocalVolume(w http.ResponseWriter, r *http.Request) {
	var requestopt = make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&requestopt); err != nil {
		w.WriteHeader(400)
		return
	}
	path := requestopt["path"]

	if err := os.RemoveAll(path); err != nil {
		logrus.Errorf("path: %s; remove pv path: %v", path, err)
		w.WriteHeader(500)
	}

	httputil.ReturnSuccess(r, w, nil)
}
