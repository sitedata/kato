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
	"os"
	"time"

	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/discover/config"
	"github.com/gridworkz/kato/monitor/prometheus"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

//RbdAPI rbd api metrics
type RbdAPI struct {
	discover.Callback
	Prometheus      *prometheus.Manager
	sortedEndpoints []string
}

//UpdateEndpoints update endpoint
func (b *RbdAPI) UpdateEndpoints(endpoints ...*config.Endpoint) {
	scrape := b.toScrape()
	b.Prometheus.UpdateScrape(scrape)
}

//Error handle error
func (b *RbdAPI) Error(err error) {
	logrus.Error(err)
}

//Name
func (b *RbdAPI) Name() string {
	return "rbdapi"
}

func (b *RbdAPI) toScrape() *prometheus.ScrapeConfig {
	ts := make([]string, 0, len(b.sortedEndpoints))
	for _, end := range b.sortedEndpoints {
		ts = append(ts, end)
	}
	namespace := os.Getenv("NAMESPACE")

	return &prometheus.ScrapeConfig{
		JobName:        b.Name(),
		ScrapeInterval: model.Duration(time.Minute),
		ScrapeTimeout:  model.Duration(30 * time.Second),
		MetricsPath:    "/metrics",
		HonorLabels:    true,
		ServiceDiscoveryConfig: prometheus.ServiceDiscoveryConfig{
			KubernetesSDConfigs: []*prometheus.SDConfig{
				&prometheus.SDConfig{
					Role: prometheus.RoleEndpoint,
					NamespaceDiscovery: prometheus.NamespaceDiscovery{
						Names: []string{namespace},
					},
					Selectors: []prometheus.SelectorConfig{
						prometheus.SelectorConfig{
							Role:  prometheus.RoleEndpoint,
							Field: "metadata.name=rbd-api-api-inner",
						},
					},
				},
			},
		},
	}
}
