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
	"testing"
	"time"
)

var cli Interface

func init() {
	cli, _ = NewPrometheus(&Options{
		Endpoint: "9999.grc42f14.8wsfp0ji.a24839.grapps.ca",
	})
}
func TestGetMetric(t *testing.T) {
	metric := cli.GetMetric("up{job=\"rbdapi\"}", time.Now())
	if len(metric.MetricData.MetricValues) == 0 {
		t.Fatal("not found metric")
	}
	t.Log(metric.MetricData.MetricValues[0].Sample.Value())
}

func TestGetMetricOverTime(t *testing.T) {
	metric := cli.GetMetricOverTime("up{job=\"rbdapi\"}", time.Now().Add(-time.Second*60), time.Now(), time.Second*10)
	if len(metric.MetricData.MetricValues) == 0 {
		t.Fatal("not found metric")
	}
	if len(metric.MetricData.MetricValues[0].Series) < 6 {
		t.Fatalf("metric series length %d is less than 6", len(metric.MetricData.MetricValues[0].Series))
	}
	t.Log(metric.MetricData.MetricValues[0].Series)
}

func TestGetMetadata(t *testing.T) {
	metas := cli.GetMetadata("rbd-system")
	if len(metas) == 0 {
		t.Fatal("meta length is 0")
	}
	for _, meta := range metas {
		t.Log(meta.Metric)
	}
}

func TestGetAppMetadata(t *testing.T) {
	metas := cli.GetAppMetadata("rbd-system", "482")
	if len(metas) == 0 {
		t.Fatal("meta length is 0")
	}
	for _, meta := range metas {
		t.Log(meta.Metric)
	}
}

func TestGetComponentMetadata(t *testing.T) {
	metas := cli.GetComponentMetadata("3be96e95700a480c9b37c6ef5daf3566", "d89ffc075ca74476b6040c8e8bae9756")
	if len(metas) == 0 {
		t.Fatal("meta length is 0")
	}
	for _, meta := range metas {
		t.Log(meta.Metric)
	}
}
