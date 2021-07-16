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
	"net/http"
	"strings"

	httputil "github.com/gridworkz/kato/util/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

//APPRegister - Service registration
func APPRegister(w http.ResponseWriter, r *http.Request) {
	appName := strings.TrimSpace(chi.URLParam(r, "app_name"))
	logrus.Infof(appName)
}

//APPDiscover - Service discovery
//Used in scenarios where real-time requirements are not high, for example, docker finds the event_log address
//Request API to return available addresses
func APPDiscover(w http.ResponseWriter, r *http.Request) {
	appName := strings.TrimSpace(chi.URLParam(r, "app_name"))
	endpoints := appService.FindAppEndpoints(appName)
	if endpoints == nil || len(endpoints) == 0 {
		httputil.ReturnError(r, w, 404, "app endpoints not found")
		return
	}
	httputil.ReturnSuccess(r, w, endpoints)
}

//APPList - List registered apps
func APPList(w http.ResponseWriter, r *http.Request) {

}
