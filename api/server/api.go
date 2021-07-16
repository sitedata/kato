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

package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gridworkz/kato/api/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

	"github.com/gridworkz/kato/util"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/cmd/api/option"

	"github.com/gridworkz/kato/api/api_routers/doc"
	"github.com/gridworkz/kato/api/api_routers/license"
	"github.com/gridworkz/kato/api/metric"
	"github.com/gridworkz/kato/api/proxy"

	"github.com/gridworkz/kato/api/api_routers/cloud"
	"github.com/gridworkz/kato/api/api_routers/version2"
	"github.com/gridworkz/kato/api/api_routers/websocket"

	apimiddleware "github.com/gridworkz/kato/api/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

//Manager apiserver
type Manager struct {
	ctx             context.Context
	cancel          context.CancelFunc
	conf            option.Config
	stopChan        chan struct{}
	r               *chi.Mux
	prometheusProxy proxy.Proxy
	etcdcli         *clientv3.Client
	exporter        *metric.Exporter
}

//NewManager
func NewManager(c option.Config, etcdcli *clientv3.Client) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	manager := &Manager{
		ctx:      ctx,
		cancel:   cancel,
		conf:     c,
		stopChan: make(chan struct{}),
		etcdcli:  etcdcli,
	}
	r := chi.NewRouter()
	manager.r = r
	manager.SetMiddleware()
	return manager
}

//SetMiddleware set api meddleware
func (m *Manager) SetMiddleware() {
	c := m.conf
	r := m.r
	r.Use(m.RequestMetric)
	r.Use(middleware.RequestID)
	//Sets an http.Request's RemoteAddr to either X-Forwarded-For or X-Real-IP
	r.Use(middleware.RealIP)
	//Logs the start and end of each request with the elapsed processing time
	if c.LoggerFile != "" {
		logerFile, err := os.OpenFile(c.LoggerFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			logrus.Errorf("open logger file %s error %s", c.LoggerFile, err.Error())
			r.Use(middleware.DefaultLogger)
		} else {
			requestLog := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(logerFile, "", log.LstdFlags)})
			r.Use(requestLog)
		}
	} else {
		r.Use(middleware.DefaultLogger)
	}
	//Gracefully absorb panics and prints the stack trace
	r.Use(middleware.Recoverer)
	//request time out
	r.Use(middleware.Timeout(time.Second * 5))
	//simple authz
	if os.Getenv("TOKEN") != "" {
		r.Use(apimiddleware.FullToken)
	}
	//simple api version
	r.Use(apimiddleware.APIVersion)
	r.Use(apimiddleware.Proxy)
}

//Start manager
func (m *Manager) Start() error {
	go m.Do()
	logrus.Info("start api router success.")
	return nil
}

//Do it
func (m *Manager) Do() {
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			m.Run()
		}
	}
}

//Stop manager
func (m *Manager) Stop() error {
	logrus.Info("api router is stopped.")
	m.cancel()
	return nil
}

//Run
func (m *Manager) Run() {
	v2R := &version2.V2{
		Cfg: &m.conf,
	}
	m.Metric()
	if m.conf.Debug {
		util.ProfilerSetup(m.r)
	}
	m.r.Get("/monitor", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("ok"))
	})
	m.r.Mount("/v2", v2R.Routes())
	m.r.Mount("/cloud", cloud.Routes())
	m.r.Mount("/", doc.Routes())
	m.r.Mount("/license", license.Routes())
	//compatible with the old version of docker
	m.r.Get("/v1/etcd/event-log/instances", m.EventLogInstance)

	m.r.Get("/kubernetes/dashboard", m.KuberntesDashboardAPI)
	//prometheus single node agent
	m.r.Get("/api/v1/query", m.PrometheusAPI)
	m.r.Get("/api/v1/query_range", m.PrometheusAPI)
	//enable websocket service and file service to the browser
	go func() {
		websocketRouter := chi.NewRouter()
		websocketRouter.Mount("/", websocket.Routes())
		websocketRouter.Mount("/logs", websocket.LogRoutes())
		websocketRouter.Mount("/app", websocket.AppRoutes())
		if m.conf.WebsocketSSL {
			logrus.Infof("websocket listen on (HTTPs) %s", m.conf.WebsocketAddr)
			logrus.Fatal(http.ListenAndServeTLS(m.conf.WebsocketAddr, m.conf.WebsocketCertFile, m.conf.WebsocketKeyFile, websocketRouter))
		} else {
			logrus.Infof("websocket listen on (HTTP) %s", m.conf.WebsocketAddr)
			logrus.Fatal(http.ListenAndServe(m.conf.WebsocketAddr, websocketRouter))
		}
	}()
	if m.conf.APISSL {
		go func() {
			pool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(m.conf.APICaFile)
			if err != nil {
				logrus.Fatal("ReadFile ca err:", err)
				return
			}
			pool.AppendCertsFromPEM(caCrt)
			s := &http.Server{
				Addr:    m.conf.APIAddrSSL,
				Handler: m.r,
				TLSConfig: &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.RequireAndVerifyClientCert,
				},
			}
			logrus.Infof("api listen on (HTTPs) %s", m.conf.APIAddrSSL)
			logrus.Fatal(s.ListenAndServeTLS(m.conf.APICertFile, m.conf.APIKeyFile))
		}()
	}
	logrus.Infof("api listen on (HTTP) %s", m.conf.APIAddr)
	logrus.Fatal(http.ListenAndServe(m.conf.APIAddr, m.r))
}

//EventLogInstance - query event server instance
func (m *Manager) EventLogInstance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(m.ctx)
	defer cancel()

	res, err := m.etcdcli.Get(ctx, "/event/instance", clientv3.WithPrefix())
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if res.Kvs != nil && len(res.Kvs) > 0 {
		result := `{"data":{"instance":[`
		for _, kv := range res.Kvs {
			result += string(kv.Value) + ","
		}
		result = result[:len(result)-1] + `]},"ok":true}`
		w.Write([]byte(result))
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(404)
	return
}

//PrometheusAPI prometheus api proxy
func (m *Manager) PrometheusAPI(w http.ResponseWriter, r *http.Request) {
	handler.GetPrometheusProxy().Proxy(w, r)
}

//KuberntesDashboardAPI proxy traffic to kubernetes dashboard
func (m *Manager) KuberntesDashboardAPI(w http.ResponseWriter, r *http.Request) {
	handler.GetKubernetesDashboardProxy().Proxy(w, r)
}

//Metric prometheus metric
func (m *Manager) Metric() {
	prometheus.MustRegister(version.NewCollector("rbd_api"))
	exporter := metric.NewExporter()
	m.exporter = exporter
	prometheus.MustRegister(exporter)
	m.r.Handle("/metrics", promhttp.Handler())
}

//RequestMetric
func (m *Manager) RequestMetric(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			path := r.RequestURI
			if strings.Index(r.RequestURI, "?") > -1 {
				path = r.RequestURI[:strings.Index(r.RequestURI, "?")]
			}
			m.exporter.RequestInc(ww.Status(), path)
		}()
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}
