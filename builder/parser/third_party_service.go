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

package parser

import (
	"encoding/json"
	"strings"

	"github.com/gridworkz/kato/builder/parser/discovery"
	"github.com/gridworkz/kato/event"
	"github.com/sirupsen/logrus"
)

// ThirdPartyServiceParse is one of the implematation of parser.Parser
type ThirdPartyServiceParse struct {
	sourceBody string

	endpoints []*discovery.Endpoint

	errors []ParseError
	logger event.Logger
}

// CreateThirdPartyServiceParse creates a new ThirdPartyServiceParse.
func CreateThirdPartyServiceParse(sourceBody string, logger event.Logger) Parser {
	return &ThirdPartyServiceParse{
		sourceBody: sourceBody,
		logger:     logger,
	}
}

// Parse blablabla
func (t *ThirdPartyServiceParse) Parse() ParseErrorList {
	// empty t.sourceBody means the service has static endpoints
	// static endpoints is no need to do service check.
	if strings.Replace(t.sourceBody, " ", "", -1) == "" {
		return nil
	}

	var info discovery.Info
	if err := json.Unmarshal([]byte(t.sourceBody), &info); err != nil {
		logrus.Errorf("wrong source_body: %v, source_body: %s", err, t.sourceBody)
		t.logger.Error("third party check input parameter error", map[string]string{"step": "parse"})
		t.errors = append(t.errors, ParseError{FatalError, "wrong input data", ""})
		return t.errors
	}
	// TODO: validate data

	d := discovery.NewDiscoverier(&info)
	err := d.Connect()
	if err != nil {
		t.logger.Error("error connecting discovery center", map[string]string{"step": "parse"})
		t.errors = append(t.errors, ParseError{FatalError, "error connecting discovery center", "please make sure " +
			"the configuration is right and the discovery center is working."})
		return t.errors
	}
	defer d.Close()
	eps, err := d.Fetch()
	if err != nil {
		t.logger.Error("error fetching endpints", map[string]string{"step": "parse"})
		t.errors = append(t.errors, ParseError{FatalError, "error fetching endpints", "please check the given key."})
		return t.errors
	}
	t.endpoints = eps

	return nil
}

// GetServiceInfo returns information of third-party service from
// the receiver *ThirdPartyServiceParse.
func (t *ThirdPartyServiceParse) GetServiceInfo() []ServiceInfo {
	serviceInfo := ServiceInfo{
		Endpoints: t.endpoints,
	}
	return []ServiceInfo{serviceInfo}
}

// GetImage is a dummy method. there is no image for Third-party service.
func (t *ThirdPartyServiceParse) GetImage() Image {
	return Image{}
}
