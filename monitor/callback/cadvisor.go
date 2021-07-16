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

package callback

import (
	"fmt"
	"time"

	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/discover/config"
	"github.com/gridworkz/kato/monitor/prometheus"
	"github.com/gridworkz/kato/monitor/utils"
	"github.com/gridworkz/kato/util"
	"github.com/gridworkz/kato/util/watch"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// Cadvisor refers to container monitoring data, which comes from kubelet on all child nodes
// 127.0.0.1:4194/metrics
type Cadvisor struct {
	discover.Callback
	Prometheus      *prometheus.Manager
	sortedEndpoints []string
	ListenPort      int

	endpoints []*config.Endpoint
}

//UpdateEndpoints
func (c *Cadvisor) UpdateEndpoints(endpoints ...*config.Endpoint) {
	newArr := utils.TrimAndSort(endpoints)

	if utils.ArrCompare(c.sortedEndpoints, newArr) {
		logrus.Debugf("The endpoints is not modify: %s", c.Name())
		return
	}

	c.sortedEndpoints = newArr

	scrape := c.toScrape()
	c.Prometheus.UpdateScrape(scrape)
}

func (c *Cadvisor) Error(err error) {
	logrus.Error(err)
}

//Name
func (c *Cadvisor) Name() string {
	return "cadvisor"
}

func (c *Cadvisor) toScrape() *prometheus.ScrapeConfig {
	apiServerHost := util.Getenv("KUBERNETES_SERVICE_HOST", "kubernetes.default.svc")
	apiServerPort := util.Getenv("KUBERNETES_SERVICE_PORT", "443")

	return &prometheus.ScrapeConfig{
		JobName:        c.Name(),
		ScrapeInterval: model.Duration(15 * time.Second),
		ScrapeTimeout:  model.Duration(10 * time.Second),
		Scheme:         "https",
		HTTPClientConfig: prometheus.HTTPClientConfig{
			TLSConfig: prometheus.TLSConfig{
				CAFile:             "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
				InsecureSkipVerify: true,
			},
			BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
		},
		ServiceDiscoveryConfig: prometheus.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*prometheus.SDConfig{
				{
					Role: "node",
				},
			},
		},
		RelabelConfigs: []*prometheus.RelabelConfig{
			{
				TargetLabel: "__address__",
				Replacement: apiServerHost + ":" + apiServerPort,
			},
			{
				SourceLabels: []model.LabelName{
					"__meta_kubernetes_node_name",
				},
				Regex:       prometheus.MustNewRegexp("(.+)"),
				TargetLabel: "__metrics_path__",
				Replacement: "/api/v1/nodes/${1}/proxy/metrics/cadvisor",
			},
			{
				Action: prometheus.RelabelAction("labelmap"),
				Regex:  prometheus.MustNewRegexp("__meta_kubernetes_node_label_(.+)"),
			},
		},
		MetricRelabelConfigs: []*prometheus.RelabelConfig{
			{
				SourceLabels: []model.LabelName{"name"},
				Regex:        prometheus.MustNewRegexp("k8s_(.*)_(.*)_(.*)_(.*)_(.*)"),
				TargetLabel:  "service_id",
				Replacement:  "${1}",
			},
			{
				SourceLabels: []model.LabelName{"name"},
				//k8s_POD_709dfaa8d9b9498a827fd5c503e0d1a1-deployment-8679ff667-j8fj8_5201d8a00fa743c18eb6553778f77c84_d6670db0-00a7-4d2c-a92e-18a19541268d_0
				Regex:       prometheus.MustNewRegexp("k8s_POD_(.*)-deployment-(.*)"),
				TargetLabel: "service_id",
				Replacement: "${1}",
			},
		},
	}
}

//AddEndpoint
func (c *Cadvisor) AddEndpoint(end *config.Endpoint) {
	c.endpoints = append(c.endpoints, end)
	c.UpdateEndpoints(c.endpoints...)
}

//Add
func (c *Cadvisor) Add(event *watch.Event) {
	url := fmt.Sprintf("%s:%d", gjson.Get(event.GetValueString(), "internal_ip").String(), c.ListenPort)
	end := &config.Endpoint{
		Name: event.GetKey(),
		URL:  url,
	}
	c.AddEndpoint(end)
}

//Modify
func (c *Cadvisor) Modify(event *watch.Event) {
	var update bool
	url := fmt.Sprintf("%s:%d", gjson.Get(event.GetValueString(), "internal_ip").String(), c.ListenPort)
	for i, end := range c.endpoints {
		if end.Name == event.GetKey() {
			c.endpoints[i].URL = url
			c.UpdateEndpoints(c.endpoints...)
			update = true
			break
		}
	}
	if !update {
		c.endpoints = append(c.endpoints, &config.Endpoint{
			Name: event.GetKey(),
			URL:  url,
		})
	}
}

//Delete
func (c *Cadvisor) Delete(event *watch.Event) {
	for i, end := range c.endpoints {
		if end.Name == event.GetKey() {
			c.endpoints = append(c.endpoints[:i], c.endpoints[i+1:]...)
			c.UpdateEndpoints(c.endpoints...)
			break
		}
	}
}
