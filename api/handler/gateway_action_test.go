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

ppackage handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	apimodel "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
)

var gm *GatewayAction

func init() {
	//conf := option.Config{
	//	DBType:           "mysql",
	//	DBConnectionInfo: "pohF4b:EiW6Eipu@tcp(192.168.56.101:3306)/region",
	//}
	////Create db manager
	//if err := db.CreateDBManager(conf); err != nil {
	//	fmt.Printf("create db manager error, %v", err)
	//}
	//cli, err := client.NewMqClient([]string{"http://192.168.56.101:2379"}, "192.168.56.101:6300")
	//if err != nil {
	//	fmt.Printf("create mq client error, %v", err)
	//}
	//gm = CreateGatewayManager(cdb.GetManager(), cli, nil)
}
func TestSelectAvailablePort(t *testing.T) {
	t.Log(selectAvailablePort([]int{9000}))                // less than minport
	t.Log(selectAvailablePort([]int{10000}))               // equal to minport
	t.Log(selectAvailablePort([]int{10003, 10001, 10001})) // more than minport and less than maxport
	t.Log(selectAvailablePort([]int{65535}))               // equal to maxport
	t.Log(selectAvailablePort([]int{10000, 65536}))        // more than maxport
}

func TestAddHTTPRule(t *testing.T) {
	for i := 200; i < 500; i++ {
		domain := fmt.Sprintf("5000-%d.gr5d6478.aq0g1f8i.4f3597.grapps.ca", i)
		err := gm.AddHTTPRule(&apimodel.AddHTTPRuleStruct{
			HTTPRuleID:    util.NewUUID(),
			ServiceID:     "68f1b4f28d49baeb68a06e1c5f5d6478",
			ContainerPort: 5000,
			Domain:        domain,
		})
		if err != nil {
			t.Fatal(err)
		}
		logrus.Infof("add domain %s", domain)
		waitReady(domain)
		time.Sleep(time.Second * 2)
	}
}
func TestWaitReady(t *testing.T) {
	waitReady("5000-1.gr5d6478.aq0g1f8i.4f3597.grapps.ca")
}
func waitReady(domain string) bool {
	start := time.Now()
	for {
		reqAddres := "http://192.168.56.101"
		if strings.Contains(domain, "192.168.56.101") {
			reqAddres = "http://" + domain
		}
		req, _ := http.NewRequest("GET", reqAddres, nil)
		req.Host = domain
		res, _ := http.DefaultClient.Do(req)
		if res != nil && res.StatusCode == 200 {
			if res.Body != nil {
				body, _ := ioutil.ReadAll(res.Body)
				res.Body.Close()
				if strings.Contains(string(body), "2048") {
					logrus.Infof("%s is ready take %s", domain, time.Now().Sub(start))
					return true
				}
			}
		}
		time.Sleep(time.Millisecond * 500)
		continue
	}
}
func TestAddTCPRule(t *testing.T) {
	for i := 1; i < 200; i++ {
		address := fmt.Sprintf("192.168.56.101:%d", 10000+i)
		gm.AddTCPRule(&apimodel.AddTCPRuleStruct{
			TCPRuleID:     util.NewUUID(),
			ServiceID:     "68f1b4f28d49baeb68a06e1c5f5d6478",
			ContainerPort: 5000,
			IP:            "192.168.56.101",
			Port:          10000 + i,
		})
		logrus.Infof("add tcp listen %s", address)
		waitReady(address)
		time.Sleep(time.Second * 2)
	}
}
func TestDeleteHTTPRule(t *testing.T) {
	gm.DeleteHTTPRule(&apimodel.DeleteHTTPRuleStruct{HTTPRuleID: ""})
}
