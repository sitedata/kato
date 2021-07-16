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

package service

import (
	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/api/model"
	"github.com/gridworkz/kato/node/masterserver"
	"github.com/gridworkz/kato/node/utils"
)

//PrometheusService
type PrometheusService struct {
	prometheusAPI *model.PrometheusAPI
	conf          *option.Conf
	ms            *masterserver.MasterServer
}

var prometheusService *PrometheusService

//CreatePrometheusService
func CreatePrometheusService(c *option.Conf, ms *masterserver.MasterServer) *PrometheusService {
	if prometheusService == nil {
		prometheusService = &PrometheusService{
			prometheusAPI: &model.PrometheusAPI{API: c.PrometheusAPI},
			conf:          c,
			ms:            ms,
		}
	}
	return prometheusService
}

//Exec prometheus query
func (ts *PrometheusService) Exec(expr string) (*model.Prome, *utils.APIHandleError) {
	resp, err := ts.prometheusAPI.Query(expr)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//ExecRange - exec prometheus query range
func (ts *PrometheusService) ExecRange(expr, start, end, step string) (*model.Prome, *utils.APIHandleError) {
	resp, err := ts.prometheusAPI.QueryRange(expr, start, end, step)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
