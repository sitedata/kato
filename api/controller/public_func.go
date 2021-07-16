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
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strconv"
	"strings"

	"github.com/gridworkz/kato/db"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

//public function part

//TransStatus - trans service status
func TransStatus(eStatus string) string {
	switch eStatus {
	case "starting":
		return "starting"
	case "abnormal":
		return "abnormal operation"
	case "upgrade":
		return "upgrade"
	case "closed":
		return "closed"
	case "stopping":
		return "stopping"
	case "checking":
		return "checking"
	case "unusual":
		return "abnormal operation"
	case "running":
		return "running"
	case "failure":
		return "failure"
	case "undeploy":
		return "undeploy"
	case "deployed":
		return "deployed"
	case "deploying":
		return "deploying"
	}
	return ""
}

//HTTPRequest public http request
func HTTPRequest(w http.ResponseWriter, r *http.Request, endpoint string) {
	DirectRequest(w, r, endpoint)
}

//DirectRequest direct request
func DirectRequest(w http.ResponseWriter, r *http.Request, endpoint string) {
	director := func(req *http.Request) {
		req = r
		req.URL.Scheme = "http"
		req.URL.Host = endpoint
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

//TransBody transfer body to map
func TransBody(r *http.Request) (map[string]interface{}, error) {
	var mmp map[string]interface{}
	err := render.Decode(r, &mmp)
	if err != nil {
		return nil, err
	}
	return mmp, nil
}

//RestInfo rest info
func RestInfo(code int, others ...string) []byte {
	switch code {
	case 200:
		msg := qMsg(others)
		if qMsg(others) != "" {
			return []byte(fmt.Sprintf(`{"ok": true, %v}`, msg))
		}
		return []byte(`{"ok": true}`)
	case 500:
		msg := qMsg(others)
		if qMsg(others) != "" {
			return []byte(fmt.Sprintf(`{"ok": false, %v}`, msg))
		}
		return []byte(`{"ok": false}`)
	case 404:
		msg := qMsg(others)
		if qMsg(others) != "" {
			return []byte(fmt.Sprintf(`{"ok": false, %v}`, msg))
		}
		return []byte(`{"ok": false}`)
	}
	return []byte(`{"ok": true}`)
}

func qMsg(others []string) string {
	var msg string
	if len(others) != 0 {
		for _, value := range others {
			if msg == "" {
				msg = value
			} else {
				msg += "," + value
			}
		}
	}
	return msg
}

//CheckLabel check label
func CheckLabel(serviceID string) bool {
	//true for v2, false for v1
	serviceLabel, err := db.GetManager().TenantServiceLabelDao().GetTenantServiceLabel(serviceID)
	if err != nil {
		return false
	}
	if serviceLabel != nil && len(serviceLabel) > 0 {
		return true
	}
	return false
}

//TestFunc test func
func TestFunc(w http.ResponseWriter) {
	w.Write([]byte("here"))
	return
}

//CheckMapKey CheckMapKey
func CheckMapKey(rebody map[string]interface{}, key string, defaultValue interface{}) map[string]interface{} {
	if _, ok := rebody[key]; ok {
		return rebody
	}
	rebody[key] = defaultValue
	return rebody
}

//GetServiceAliasID get service alias id
//python:
//new_word = str(ord(string[10])) + string + str(ord(string[3])) + 'log' + str(ord(string[2]) / 7)
//new_id = hashlib.sha224(new_word).hexdigest()[0:16]
func GetServiceAliasID(ServiceID string) string {
	if len(ServiceID) > 11 {
		newWord := strconv.Itoa(int(ServiceID[10])) + ServiceID + strconv.Itoa(int(ServiceID[3])) + "log" + strconv.Itoa(int(ServiceID[2])/7)
		ha := sha256.New224()
		ha.Write([]byte(newWord))
		return fmt.Sprintf("%x", ha.Sum(nil))[0:16]
	}
	return ServiceID
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

//Display - Traversal structure
func Display(i interface{}, j interface{}) {
	fmt.Printf("Display %T\n", i)
	display("", reflect.ValueOf(i), j)

}

func display(path string, v reflect.Value, j interface{}) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool,
		reflect.String:
		fmt.Printf("%s: %v (%s)\n", path, v, v.Type().Name())
		return
	case reflect.Ptr:
		v = v.Elem()
		display(path, v, j)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(" "+path+"["+strconv.Itoa(i)+"]", v.Index(i), j)
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			display(" "+path+"."+t.Field(i).Name, v.Field(i), j)
		}
	}
}
