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

package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gridworkz/kato/api/handler"
	api_model "github.com/gridworkz/kato/api/model"
	ctxutil "github.com/gridworkz/kato/api/util/ctx"
	"github.com/gridworkz/kato/cmd/api/option"
	"github.com/gridworkz/kato/mq/client"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	k8svalidation "k8s.io/apimachinery/pkg/util/validation"
)

// GatewayStruct -
type GatewayStruct struct {
	MQClient client.MQClient
	cfg      *option.Config
}

// HTTPRule is used to add, update or delete http rule which enables
// external traffic to access applications through the gateway
func (g *GatewayStruct) HTTPRule(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		g.addHTTPRule(w, r)
	case "PUT":
		g.updateHTTPRule(w, r)
	case "DELETE":
		g.deleteHTTPRule(w, r)
	}
}

func validateDomain(domain string) []string {
	if strings.TrimSpace(domain) == "" {
		return nil
	}
	var errs []string
	if strings.Contains(domain, "*") {
		errs = k8svalidation.IsWildcardDNS1123Subdomain(domain)
	} else {
		errs = k8svalidation.IsDNS1123Subdomain(domain)
	}
	return errs
}

func (g *GatewayStruct) addHTTPRule(w http.ResponseWriter, r *http.Request) {
	var req api_model.AddHTTPRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	// verify request
	values := url.Values{}
	if req.ContainerPort == 0 {
		values["container_port"] = []string{"The container_port field is required"}
	}
	if strings.Replace(req.CertificateID, " ", "", -1) != "" {
		if req.Certificate == "" {
			values["certificate"] = []string{"The certificate field is required"}
		}
		if req.PrivateKey == "" {
			values["private_key"] = []string{"The private_key field is required"}
		}
	}
	errs := validateDomain(req.Domain)
	if errs != nil && len(errs) > 0 {
		logrus.Debugf("Invalid domain: %s", strings.Join(errs, ";"))
		values["domain"] = []string{"The domain field is invalid"}
	}
	if len(values) != 0 {
		httputil.ReturnValidationError(r, w, values)
		return
	}

	h := handler.GetGatewayHandler()
	err := h.AddHTTPRule(&req)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Unexpected error occorred while adding http rule: %v", err))
		return
	}

	httputil.ReturnSuccess(r, w, req)
}

func (g *GatewayStruct) updateHTTPRule(w http.ResponseWriter, r *http.Request) {
	var req api_model.UpdateHTTPRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}
	// verify request
	values := url.Values{}
	if strings.Replace(req.CertificateID, " ", "", -1) != "" {
		if req.Certificate == "" {
			values["certificate"] = []string{"The certificate field is required"}
		}
		if req.PrivateKey == "" {
			values["private_key"] = []string{"The private_key field is required"}
		}
	}
	if len(req.RuleExtensions) > 0 {
		for _, re := range req.RuleExtensions {
			if re.Key == "" {
				values["key"] = []string{"The key field is required"}
				break
			}
			if re.Value == "" {
				values["value"] = []string{"The value field is required"}
				break
			}
		}
	}
	errs := validateDomain(req.Domain)
	if errs != nil && len(errs) > 0 {
		logrus.Debugf("Invalid domain: %s", strings.Join(errs, ";"))
		values["domain"] = []string{"The domain field is invalid"}
	}
	if len(values) != 0 {
		httputil.ReturnValidationError(r, w, values)
		return
	}

	h := handler.GetGatewayHandler()
	err := h.UpdateHTTPRule(&req)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Unexpected error occorred while "+
			"updating http rule: %v", err))
		return
	}

	httputil.ReturnSuccess(r, w, "success")
}

func (g *GatewayStruct) deleteHTTPRule(w http.ResponseWriter, r *http.Request) {
	var req api_model.DeleteHTTPRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	h := handler.GetGatewayHandler()
	err := h.DeleteHTTPRule(&req)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Unexpected error occorred while delete http rule: %v", err))
		return
	}

	httputil.ReturnSuccess(r, w, "success")
}

// TCPRule is used to add, update or delete tcp rule which enables
// external traffic to access applications through the gateway
func (g *GatewayStruct) TCPRule(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		g.AddTCPRule(w, r)
	case "PUT":
		g.updateTCPRule(w, r)
	case "DELETE":
		g.deleteTCPRule(w, r)
	}
}

// AddTCPRule adds a tcp rule
func (g *GatewayStruct) AddTCPRule(w http.ResponseWriter, r *http.Request) {
	var req api_model.AddTCPRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	h := handler.GetGatewayHandler()
	// verify request
	values := url.Values{}
	if req.ContainerPort == 0 {
		values["container_port"] = []string{"The container_port field is required"}
	}
	if req.Port == 0 {
		values["port"] = []string{"The port field is required"}
	} else if req.Port <= g.cfg.MinExtPort {
		values["port"] = []string{fmt.Sprintf("The port field should be greater than %d", g.cfg.MinExtPort)}
	} else {
		// check if the port exists
		if h.TCPIPPortExists(req.IP, req.Port) {
			values["port"] = []string{fmt.Sprintf("The ip %s port(%v) already exists", req.IP, req.Port)}
		}
	}
	if len(req.RuleExtensions) > 0 {
		for _, re := range req.RuleExtensions {
			if re.Key == "" {
				values["key"] = []string{"The key field is required"}
				break
			}
			if re.Value == "" {
				values["value"] = []string{"The value field is required"}
				break
			}
		}
	}
	if len(values) != 0 {
		httputil.ReturnValidationError(r, w, values)
		return
	}
	err := h.AddTCPRule(&req)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Unexpected error occorred while "+
			"adding tcp rule: %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, "success")
}

func (g *GatewayStruct) updateTCPRule(w http.ResponseWriter, r *http.Request) {
	var req api_model.UpdateTCPRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	h := handler.GetGatewayHandler()
	// verify reqeust
	values := url.Values{}
	if req.Port != 0 && req.Port <= g.cfg.MinExtPort {
		values["port"] = []string{fmt.Sprintf("The port field should be greater than %d", g.cfg.MinExtPort)}
	}
	if len(req.RuleExtensions) > 0 {
		for _, re := range req.RuleExtensions {
			if re.Key == "" {
				values["key"] = []string{"The key field is required"}
				break
			}
			if re.Value == "" {
				values["value"] = []string{"The value field is required"}
				break
			}
		}
	}
	if len(values) != 0 {
		httputil.ReturnValidationError(r, w, values)
		return
	}

	err := h.UpdateTCPRule(&req, g.cfg.MinExtPort)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Unexpected error occorred while "+
			"updating tcp rule: %v", err))
		return
	}

	httputil.ReturnSuccess(r, w, "success")
}

func (g *GatewayStruct) deleteTCPRule(w http.ResponseWriter, r *http.Request) {
	var req api_model.DeleteTCPRuleStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	h := handler.GetGatewayHandler()
	err := h.DeleteTCPRule(&req)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Unexpected error occorred while "+
			"deleting tcp rule: %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, "success")
}

// GetAvailablePort returns a available port
func (g *GatewayStruct) GetAvailablePort(w http.ResponseWriter, r *http.Request) {
	h := handler.GetGatewayHandler()
	lock, _ := strconv.ParseBool(r.FormValue("lock"))
	res, err := h.GetAvailablePort("0.0.0.0", lock)
	if err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Unexpected error occorred while "+
			"getting available port: %v", err))
		return
	}
	httputil.ReturnSuccess(r, w, res)
}

// RuleConfig is used to add, update or delete rule config.
func (g *GatewayStruct) RuleConfig(w http.ResponseWriter, r *http.Request) {
	var req api_model.RuleConfigReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	sid := r.Context().Value(ctxutil.ContextKey("service_id")).(string)
	eventID := r.Context().Value(ctxutil.ContextKey("event_id")).(string)
	req.ServiceID = sid
	req.EventID = eventID
	if err := handler.GetGatewayHandler().RuleConfig(&req); err != nil {
		httputil.ReturnError(r, w, 500, fmt.Sprintf("Rule id: %s; error update rule config: %v", req.RuleID, err))
		return
	}
	httputil.ReturnSuccess(r, w, "success")
}

// Certificate -
func (g *GatewayStruct) Certificate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		g.updCertificate(w, r)
	}
}

//updCertificate updates certificate and refresh http rules based on certificate id
func (g *GatewayStruct) updCertificate(w http.ResponseWriter, r *http.Request) {
	var req api_model.UpdCertificateReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &req, nil)
	if !ok {
		return
	}

	if err := handler.GetGatewayHandler().UpdCertificate(&req); err != nil {
		logrus.Errorf("update certificate: %v", err)
		if err == gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 404, err.Error())
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}

	httputil.ReturnSuccess(r, w, nil)
}

//GetGatewayIPs get gateway ips
func GetGatewayIPs(w http.ResponseWriter, r *http.Request) {
	ips := handler.GetGatewayHandler().GetGatewayIPs()
	httputil.ReturnSuccess(r, w, ips)
}
