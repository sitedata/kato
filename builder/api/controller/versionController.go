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

package controller

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/db"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/sirupsen/logrus"
)

func GetVersionByEventID(w http.ResponseWriter, r *http.Request) {
	eventID := strings.TrimSpace(chi.URLParam(r, "eventID"))

	version, err := db.GetManager().VersionInfoDao().GetVersionByEventID(eventID)
	if err != nil {
		httputil.ReturnError(r, w, 404, err.Error())
	}
	httputil.ReturnSuccess(r, w, version)
}

func UpdateVersionByEventID(w http.ResponseWriter, r *http.Request) {
	eventID := strings.TrimSpace(chi.URLParam(r, "eventID"))

	version, err := db.GetManager().VersionInfoDao().GetVersionByEventID(eventID)
	if err != nil {
		httputil.ReturnError(r, w, 404, err.Error())
		return
	}
	in, err := ioutil.ReadAll(r.Body)
	json, err := simplejson.NewJson(in)
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}

	if author, err := json.Get("code_commit_author").String(); err != nil || author == "" {
		logrus.Debugf("error get code_commit_author from version body ,details %s", err.Error())
	} else {
		version.Author = author
	}

	if msg, err := json.Get("code_commit_msg").String(); err != nil || msg == "" {
		logrus.Debugf("error get code_commit_msg from version body ,details %s", err.Error())
	} else {
		version.CommitMsg = msg
	}
	if cVersion, err := json.Get("code_version").String(); err != nil || cVersion == "" {
		logrus.Debugf("error get code_version from version body ,details %s", err.Error())
	} else {
		version.CodeVersion = cVersion
	}

	if status, err := json.Get("final_status").String(); err != nil || status == "" {
		logrus.Debugf("error get final_status from version body ,details %s", err.Error())
	} else {
		version.FinalStatus = status
	}
	err = db.GetManager().VersionInfoDao().UpdateModel(version)
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}
func GetVersionByServiceID(w http.ResponseWriter, r *http.Request) {
	serviceID := strings.TrimSpace(chi.URLParam(r, "serviceID"))

	versions, err := db.GetManager().VersionInfoDao().GetVersionByServiceID(serviceID)
	if err != nil {
		httputil.ReturnError(r, w, 404, err.Error())
	}
	httputil.ReturnSuccess(r, w, versions)
}
func DeleteVersionByEventID(w http.ResponseWriter, r *http.Request) {
	eventID := strings.TrimSpace(chi.URLParam(r, "eventID"))

	versionInfo, err := db.GetManager().VersionInfoDao().GetVersionByEventID(eventID)
	if versionInfo.DeliveredType == "" || versionInfo.DeliveredPath == "" {
		httputil.ReturnError(r, w, 500, errors.New("deliverable type and delivery path are empty").Error())
		return
	}
	if versionInfo.DeliveredType == "code" {
		//todo It's dangerous.
	} else {
		cmd := exec.Command("docker", "rmi", versionInfo.DeliveredPath)
		err := cmd.Start()
		if err != nil {
			logrus.Errorf("error delete image :%s ,details %s", versionInfo.DeliveredPath, err.Error())
		}
	}
	err = db.GetManager().VersionInfoDao().DeleteVersionByEventID(eventID)
	if err != nil {
		httputil.ReturnError(r, w, 404, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}
func UpdateDeliveredPath(w http.ResponseWriter, r *http.Request) {
	in, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}
	logrus.Infof("update build info to %s", string(in))
	jsonc, err := simplejson.NewJson(in)
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}
	event, err := jsonc.Get("event_id").String()
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}
	dt, err := jsonc.Get("type").String()
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}
	dp, err := jsonc.Get("path").String()
	if err != nil {
		httputil.ReturnError(r, w, 400, err.Error())
		return
	}
	version, err := db.GetManager().VersionInfoDao().GetVersionByEventID(event)
	if err != nil {
		httputil.ReturnError(r, w, 404, err.Error())
		return
	}

	version.DeliveredType = dt
	version.DeliveredPath = dp
	if version.DeliveredType == "slug" {
		version.ImageName = builder.RUNNERIMAGENAME
	} else {
		version.ImageName = dp
	}
	err = db.GetManager().VersionInfoDao().UpdateModel(version)
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
	return
}
