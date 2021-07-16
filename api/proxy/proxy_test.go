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

package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestHttpProxy(t *testing.T) {
	proxy := CreateProxy("prometheus", "http", []string{"http://106.14.145.76:9999"})

	query := fmt.Sprintf(`sum(app_resource_appfs{tenant_id=~"%s"}) by(tenant_id)`, strings.Join([]string{"824b2e9dcc4d461a852ddea20369d377"}, "|"))
	query = strings.Replace(query, " ", "%20", -1)
	fmt.Printf("http://127.0.0.1:9999/api/v1/query?query=%s", query)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:9999/api/v1/query?query=%s", query), nil)
	if err != nil {
		logrus.Error("create request prometheus api error ", err.Error())
		return
	}
	result, err := proxy.Do(req)
	if err != nil {
		logrus.Error("do proxy request prometheus api error ", err.Error())
		return
	}
	if result.Body != nil {
		defer result.Body.Close()
		if result.StatusCode != 200 {
			fmt.Println(result.StatusCode)
		}
		// var qres queryResult
		// err = json.NewDecoder(result.Body).Decode(&qres)
		// fmt.Println(qres)
		B, _ := ioutil.ReadAll(result.Body)
		fmt.Println(string(B))
	}
}
