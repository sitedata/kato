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

package prometheus

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	mv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	yaml "gopkg.in/yaml.v2"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/util/workqueue"
)

var smYaml = `
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: region-tokenchecker
  namespace: default
spec:
  jobLabel: service_alias
  endpoints:
  - interval: 10s
    path: /metrics
    port: tcp-9090
    relabelings:
      - sourceLabels: [__address__]
        targetLabel: app_name
        replacement: "region-tokenchecker"
  namespaceSelector:
    any: true
  selector:
    matchLabels:
      service_port: 9090
      port_protocol: http
      name: gr0a581fService
`

func TestCreateScrapeBySM(t *testing.T) {
	var smc ServiceMonitorController
	var sm mv1.ServiceMonitor
	k8syaml.NewYAMLOrJSONDecoder(bytes.NewBuffer([]byte(smYaml)), 1024).Decode(&sm)
	var scrapes []*ScrapeConfig
	t.Logf("%+v", sm)
	for i, ep := range sm.Spec.Endpoints {
		scrapes = append(scrapes, smc.createScrapeBySM(&sm, ep, i))
	}
	out, err := yaml.Marshal(scrapes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
}

func TestQueue(t *testing.T) {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "sm-monitor")
	defer queue.ShutDown()

	go func() {
		for i := 0; i < 10; i++ {
			queue.Add("abc")
			time.Sleep(time.Second * 1)
		}
	}()
	for {
		item, close := queue.Get()
		if close {
			t.Fatal("queue closed")
		}
		time.Sleep(time.Second * 2)
		fmt.Println(item)
		queue.Forget(item)
		queue.Done(item)
	}
}
