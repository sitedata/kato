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
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/builder/model"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	httputil "github.com/gridworkz/kato/util/http"
	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/gridworkz/kato/builder/discover"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func AddCodeCheck(w http.ResponseWriter, r *http.Request) {
	//b,_:=ioutil.ReadAll(r.Body)
	//{\"url_repos\": \"https://github.com/bay1ts/zk_cluster_mini.git\", \"check_type\": \"first_check\", \"code_from\": \"gitlab_manual\", \"service_id\": \"c24dea8300b9401b1461dd975768881a\", \"code_version\": \"master\", \"git_project_id\": 0, \"condition\": \"{\\\"language\\\":\\\"docker\\\",\\\"runtimes\\\":\\\"false\\\", \\\"dependencies\\\":\\\"false\\\",\\\"procfile\\\":\\\"false\\\"}\", \"git_url\": \"--branch master --depth 1 https://github.com/bay1ts/zk_cluster_mini.git\"}
	//logrus.Infof("request recive %s",string(b))
	result := new(model.CodeCheckResult)

	b, _ := ioutil.ReadAll(r.Body)
	j, err := simplejson.NewJson(b)
	if err != nil {
		logrus.Errorf("error decode json,details %s", err.Error())
		httputil.ReturnError(r, w, 400, "bad request")
		return
	}
	result.URLRepos, _ = j.Get("url_repos").String()
	result.CheckType, _ = j.Get("check_type").String()
	result.CodeFrom, _ = j.Get("code_from").String()
	result.ServiceID, _ = j.Get("service_id").String()
	result.CodeVersion, _ = j.Get("code_version").String()
	result.GitProjectId, _ = j.Get("git_project_id").String()
	result.Condition, _ = j.Get("condition").String()
	result.GitURL, _ = j.Get("git_url").String()

	defer r.Body.Close()

	dbmodel := convertModelToDB(result)
	//checkAndGet
	db.GetManager().CodeCheckResultDao().AddModel(dbmodel)
	httputil.ReturnSuccess(r, w, nil)
}
func Update(w http.ResponseWriter, r *http.Request) {
	serviceID := strings.TrimSpace(chi.URLParam(r, "serviceID"))
	result := new(model.CodeCheckResult)

	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	logrus.Infof("update receive %s", string(b))
	j, err := simplejson.NewJson(b)
	if err != nil {
		logrus.Errorf("error decode json,details %s", err.Error())
		httputil.ReturnError(r, w, 400, "bad request")
		return
	}
	result.BuildImageName, _ = j.Get("image").String()
	portList, err := j.Get("port_list").Map()
	if err != nil {
		portList = make(map[string]interface{})
	}
	volumeList, err := j.Get("volume_list").StringArray()
	if err != nil {
		volumeList = nil
	}
	strMap := make(map[string]string)
	for k, v := range portList {
		strMap[k] = v.(string)
	}
	result.VolumeList = volumeList
	result.PortList = strMap
	result.ServiceID = serviceID
	dbmodel := convertModelToDB(result)
	dbmodel.DockerFileReady = true
	db.GetManager().CodeCheckResultDao().UpdateModel(dbmodel)
	httputil.ReturnSuccess(r, w, nil)
}
func convertModelToDB(result *model.CodeCheckResult) *dbmodel.CodeCheckResult {
	r := dbmodel.CodeCheckResult{}
	r.ServiceID = result.ServiceID
	r.CheckType = result.CheckType
	r.CodeFrom = result.CodeFrom
	r.CodeVersion = result.CodeVersion
	r.Condition = result.Condition
	r.GitProjectId = result.GitProjectId
	r.GitURL = result.GitURL
	r.URLRepos = result.URLRepos

	if result.Condition != "" {
		bs := []byte(result.Condition)
		l, err := simplejson.NewJson(bs)
		if err != nil {
			logrus.Errorf("error get condition,details %s", err.Error())
		}
		language, err := l.Get("language").String()
		if err != nil {
			logrus.Errorf("error get language,details %s", err.Error())
		}
		r.Language = language
	}
	r.BuildImageName = result.BuildImageName
	r.InnerPort = result.InnerPort
	pl, _ := json.Marshal(result.PortList)
	r.PortList = string(pl)
	vl, _ := json.Marshal(result.VolumeList)
	r.VolumeList = string(vl)
	r.VolumeMountPath = result.VolumeMountPath
	return &r
}
func GetCodeCheck(w http.ResponseWriter, r *http.Request) {
	serviceID := strings.TrimSpace(chi.URLParam(r, "serviceID"))
	//findResultByServiceID
	cr, err := db.GetManager().CodeCheckResultDao().GetCodeCheckResult(serviceID)
	if err != nil {
		logrus.Errorf("error get check result,details %s", err.Error())
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, cr)
}

func CheckHalth(w http.ResponseWriter, r *http.Request) {
	healthInfo := discover.HealthCheck()
	if healthInfo["status"] != "health" {
		httputil.ReturnError(r, w, 400, "builder service unusual")
	}
	httputil.ReturnSuccess(r, w, healthInfo)
}
