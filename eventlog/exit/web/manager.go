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

package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gridworkz/kato/eventlog/cluster"
	"github.com/gridworkz/kato/eventlog/cluster/discover"
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/exit/monitor"
	"github.com/gridworkz/kato/eventlog/store"
	"github.com/gridworkz/kato/util"
	httputil "github.com/gridworkz/kato/util/http"

	"github.com/coreos/etcd/clientv3"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"golang.org/x/net/context"
)

//SocketServer
type SocketServer struct {
	conf                 conf.WebSocketConf
	discoverConf         conf.DiscoverConf
	log                  *logrus.Entry
	cancel               func()
	context              context.Context
	storemanager         store.Manager
	listenErr, errorStop chan error
	reStart              int
	timeout              time.Duration
	cluster              cluster.Cluster
	healthInfo           map[string]string
	etcdClient           *clientv3.Client
	pubsubCtx            map[string]*PubContext
}

//NewSocket - create zmq sub client
func NewSocket(conf conf.WebSocketConf, discoverConf conf.DiscoverConf, etcdClient *clientv3.Client, log *logrus.Entry, storeManager store.Manager, c cluster.Cluster, healthInfo map[string]string) *SocketServer {
	ctx, cancel := context.WithCancel(context.Background())
	d, err := time.ParseDuration(conf.TimeOut)
	if err != nil {
		d = time.Minute * 1
	}

	return &SocketServer{
		conf:         conf,
		discoverConf: discoverConf,
		log:          log,
		cancel:       cancel,
		context:      ctx,
		storemanager: storeManager,
		listenErr:    make(chan error),
		errorStop:    make(chan error),
		timeout:      d,
		cluster:      c,
		healthInfo:   healthInfo,
		etcdClient:   etcdClient,
		pubsubCtx:    make(map[string]*PubContext),
	}
}

func (s *SocketServer) pushEventMessage(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:    s.conf.ReadBufferSize,
		WriteBufferSize:   s.conf.WriteBufferSize,
		EnableCompression: s.conf.EnableCompression,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {

		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Create web socket conn error.", err.Error())
		return
	}
	defer conn.Close()
	_, me, err := conn.ReadMessage()
	if err != nil {
		s.log.Error("Read EventID from first message error.", err.Error())
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("ok"))
	info := strings.Split(string(me), "=")
	if len(info) != 2 {
		s.log.Error("Read EventID from first message error. The data format is not correct")
		return
	}
	EventID := info[1]
	if EventID == "" {
		s.log.Error("Event ID can not be empty when get socket message")
		return
	}
	s.log.Infof("Begin push event message of event (%s)", EventID)
	SubID := uuid.NewV4().String()
	ch := s.storemanager.WebSocketMessageChan("event", EventID, SubID)
	if ch == nil {
		// w.Write([]byte("Real-time message does not exist."))
		// w.Header().Set("Status Code", "200")
		s.log.Error("get web socket message chan from storemanager error.")
		return
	}
	defer func() {
		s.log.Debug("Push event message request closed")
		s.storemanager.RealseWebSocketMessageChan("event", EventID, SubID)
	}()
	stop := make(chan struct{})
	go s.reader(conn, stop)
	pingTicker := time.NewTicker(s.timeout * 8 / 10)
	defer pingTicker.Stop()
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				return
			}
			if message != nil {
				//s.log.Debugf("websocket push a message,%s", message.Message)
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err = conn.WriteMessage(websocket.TextMessage, message.Content)
				if err != nil {
					s.log.Warn("Push message to client error.", err.Error())
					return
				}
			}
		case <-stop:
			return
		case <-s.context.Done():
			return
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}

}

func (s *SocketServer) pushDockerLog(w http.ResponseWriter, r *http.Request) {
	// if r.FormValue("host") == "" || r.FormValue("host") != s.cluster.GetInstanceID() {
	// 	w.WriteHeader(404)
	// 	return
	// }
	upgrader := websocket.Upgrader{
		ReadBufferSize:    s.conf.ReadBufferSize,
		WriteBufferSize:   s.conf.WriteBufferSize,
		EnableCompression: s.conf.EnableCompression,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {

		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Create web socket conn error.", err.Error())
		return
	}
	defer conn.Close()
	_, me, err := conn.ReadMessage()
	if err != nil {
		s.log.Error("Read ServiceID from first message error.", err.Error())
		return
	}
	info := strings.Split(string(me), "=")
	if len(info) != 2 {
		s.log.Error("Read ServiceID from first message error. The data format is not correct")
		return
	}
	ServiceID := info[1]
	if ServiceID == "" {
		s.log.Error("ServiceID ID can not be empty when get socket message")
		return
	}
	s.log.Infof("Begin push docker message of service (%s)", ServiceID)
	SubID := uuid.NewV4().String()
	ch := s.storemanager.WebSocketMessageChan("docker", ServiceID, SubID)
	if ch == nil {
		// w.Write([]byte("Real-time message does not exist."))
		// w.Header().Set("Status Code", "200")
		s.log.Error("get web socket message chan from storemanager error.")
		return
	}
	defer func() {
		s.log.Debug("Push docker log message request closed")
		s.storemanager.RealseWebSocketMessageChan("docker", ServiceID, SubID)
	}()
	conn.WriteMessage(websocket.TextMessage, []byte("ok"))
	stop := make(chan struct{})
	go s.reader(conn, stop)
	pingTicker := time.NewTicker(s.timeout * 8 / 10)
	defer pingTicker.Stop()
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				return
			}
			if message != nil {
				s.log.Debugf("websocket push a message: %v", message)
				err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					s.log.Warningf("error setting write deadline: %v", err)
				}
				err = conn.WriteMessage(websocket.TextMessage, message.Content)
				if err != nil {
					s.log.Warn("Push message to client error.", err.Error())
					return
				}
			}
		case <-stop:
			return
		case <-s.context.Done():
			return
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}

}
func (s *SocketServer) pushMonitorMessage(w http.ResponseWriter, r *http.Request) {
	// if r.FormValue("host") == "" || r.FormValue("host") != s.cluster.GetInstanceID() {
	// 	w.WriteHeader(404)
	// 	return
	// }
	upgrader := websocket.Upgrader{
		ReadBufferSize:    s.conf.ReadBufferSize,
		WriteBufferSize:   s.conf.WriteBufferSize,
		EnableCompression: s.conf.EnableCompression,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {

		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Create web socket conn error.", err.Error())
		return
	}
	defer conn.Close()
	_, me, err := conn.ReadMessage()
	if err != nil {
		s.log.Error("Read tag key from first message error.", err.Error())
		return
	}
	info := strings.Split(string(me), "=")
	if len(info) != 2 {
		s.log.Error("Read tag key from first message error. The data format is not correct")
		return
	}
	ServiceID := info[1]
	if ServiceID == "" {
		s.log.Error("tag key can not be empty when get socket message")
		return
	}
	s.log.Infof("Begin push monitor message of service (%s)", ServiceID)
	SubID := uuid.NewV4().String()
	ch := s.storemanager.WebSocketMessageChan("monitor", ServiceID, SubID)
	if ch == nil {
		// w.Write([]byte("Real-time message does not exist."))
		// w.Header().Set("Status Code", "200")
		s.log.Error("get web socket message chan from storemanager error.")
		return
	}
	defer func() {
		s.log.Debug("Push docker log message request closed")
		s.storemanager.RealseWebSocketMessageChan("monitor", ServiceID, SubID)
	}()
	conn.WriteMessage(websocket.TextMessage, []byte("ok"))
	stop := make(chan struct{})
	go s.reader(conn, stop)
	pingTicker := time.NewTicker(s.timeout * 8 / 10)
	defer pingTicker.Stop()
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				return
			}
			if message != nil {
				s.log.Debugf("websocket push a monitor message,%s", string(message.MonitorData))
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err = conn.WriteMessage(websocket.TextMessage, message.MonitorData)
				if err != nil {
					s.log.Warn("Push message to client error.", err.Error())
					return
				}
			}
		case <-stop:
			return
		case <-s.context.Done():
			return
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}

}
func (s *SocketServer) pushNewMonitorMessage(w http.ResponseWriter, r *http.Request) {
	// if r.FormValue("host") == "" || r.FormValue("host") != s.cluster.GetInstanceID() {
	// 	w.WriteHeader(404)
	// 	return
	// }
	upgrader := websocket.Upgrader{
		ReadBufferSize:    s.conf.ReadBufferSize,
		WriteBufferSize:   s.conf.WriteBufferSize,
		EnableCompression: s.conf.EnableCompression,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {

		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Create web socket conn error.", err.Error())
		return
	}
	defer conn.Close()
	_, me, err := conn.ReadMessage()
	if err != nil {
		s.log.Error("Read tag key from first message error.", err.Error())
		return
	}
	info := strings.Split(string(me), "=")
	if len(info) != 2 {
		s.log.Error("Read tag key from first message error. The data format is not correct")
		return
	}
	ServiceID := info[1]
	if ServiceID == "" {
		s.log.Error("tag key can not be empty when get socket message")
		return
	}
	s.log.Infof("Begin push monitor message of service (%s)", ServiceID)
	SubID := uuid.NewV4().String()
	ch := s.storemanager.WebSocketMessageChan("newmonitor", ServiceID, SubID)
	if ch == nil {
		// w.Write([]byte("Real-time message does not exist."))
		// w.Header().Set("Status Code", "200")
		s.log.Error("get web socket message chan from storemanager error.")
		return
	}
	defer func() {
		s.log.Debug("Push new monitor message request closed")
		s.storemanager.RealseWebSocketMessageChan("newmonitor", ServiceID, SubID)
	}()
	conn.WriteMessage(websocket.TextMessage, []byte("ok"))
	stop := make(chan struct{})
	go s.reader(conn, stop)
	pingTicker := time.NewTicker(s.timeout * 8 / 10)
	defer pingTicker.Stop()
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				return
			}
			if message != nil {
				s.log.Debugf("websocket push a new monitor message")
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err = conn.WriteMessage(websocket.TextMessage, message.MonitorData)
				if err != nil {
					s.log.Warn("Push message to client error.", err.Error())
					return
				}
			}
		case <-stop:
			return
		case <-s.context.Done():
			return
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}

}
func (s *SocketServer) reader(ws *websocket.Conn, ch chan struct{}) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(s.timeout)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
	s.log.Debug("socket conn ping/pong time out ,will closed.")
	close(ch)
}

//Run
func (s *SocketServer) Run() error {
	s.log.Info("WebSocker Server start")
	go s.listen()
	go s.checkHealth()
	return nil
}
func (s *SocketServer) listen() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// deprecated
	r.Get("/event_log", s.pushEventMessage)
	// deprecated
	r.Get("/docker_log", s.pushDockerLog)
	// deprecated
	r.Get("/monitor_message", s.pushMonitorMessage)
	// deprecated
	r.Get("/new_monitor_message", s.pushNewMonitorMessage)

	r.Get("/monitor", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	r.Get("/docker-instance", func(w http.ResponseWriter, r *http.Request) {
		ServiceID := r.FormValue("service_id")
		if ServiceID == "" {
			w.WriteHeader(412)
			w.Write([]byte(`{"message":"service id can not be empty.","status":"failure"}`))
			return
		}
		s.log.Info("ServiceID:" + ServiceID)
		instance := s.cluster.GetSuitableInstance(ServiceID)
		err := discover.SaveDockerLogInInstance(s.etcdClient, s.discoverConf, ServiceID, instance.HostID)
		if err != nil {
			s.log.Error("Save docker service and instance id to etcd error.")
			w.WriteHeader(500)
			w.Write([]byte(`{"message":"Save docker service and instance id to etcd error.","status":"failure"}`))
			return
		}
		w.WriteHeader(200)
		url := fmt.Sprintf("tcp://%s:%d", instance.HostIP, instance.DockerLogPort)
		w.Write([]byte(`{"host":"` + url + `","status":"success"}`))
	})
	r.Get("/event_push", s.receiveEventMessage)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if s.healthInfo["status"] != "health" {
			httputil.ReturnError(r, w, 400, "eventlog service unusual")
		}
		httputil.ReturnSuccess(r, w, s.healthInfo)
	})
	// new websocket pubsub
	r.Get("/services/{serviceID}/pubsub", s.pubsub)
	r.Get("/tenants/{tenantName}/services/{serviceID}/logs", s.getDockerLogs)
	//monitor setting
	s.prometheus(r)
	//pprof debug
	util.ProfilerSetup(r)

	if s.conf.SSL {
		go func() {
			addr := fmt.Sprintf("%s:%d", s.conf.BindIP, s.conf.SSLBindPort)
			s.log.Infof("web socket ssl server listen %s", addr)
			err := http.ListenAndServeTLS(addr, s.conf.CertFile, s.conf.KeyFile, r)
			if err != nil {
				s.log.Error("websocket listen error.", err.Error())
				s.listenErr <- err
			}
		}()
	}
	addr := fmt.Sprintf("%s:%d", s.conf.BindIP, s.conf.BindPort)
	s.log.Infof("web socket server listen %s", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		s.log.Error("websocket listen error.", err.Error())
		s.listenErr <- err
	}
}
func (s *SocketServer) checkHealth() {
	tike := time.Tick(time.Minute * 10)
	for {
		select {
		case <-s.context.Done():
			return
		case <-tike:
			s.reStart = 0
		case err := <-s.listenErr:
			if s.reStart > s.conf.MaxRestartCount {
				s.log.Error("Web socket server listen error count more than max restart count.")
				s.errorStop <- err
			} else {
				go s.listen()
				s.reStart++
			}
		}
	}
}

//ListenError - return error channel
func (s *SocketServer) ListenError() chan error {
	return s.errorStop
}

//Stop
func (s *SocketServer) Stop() {
	s.log.Info("WebSocker Server stop")
	s.cancel()
}

//ReceiveEventMessage - receive operation log API
func (s *SocketServer) receiveEventMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var re ResponseType
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		re = NewResponseType(500, err.Error(), "Error reading event message content", nil, nil)
	} else {
		select {
		case s.storemanager.ReceiveMessageChan() <- message:
			re = NewSuccessResponse(nil, nil)
			w.WriteHeader(200)
		default:
			re = NewResponseType(500, "event message chan is block", "Event message channel blocked", nil, nil)
			w.WriteHeader(500)
		}
	}
	if r.Body != nil {
		r.Body.Close()
	}
	json.NewEncoder(w).Encode(re)
	return
}

func (s *SocketServer) prometheus(r *chi.Mux) {
	prometheus.MustRegister(version.NewCollector("event_log"))
	exporter := monitor.NewExporter(s.storemanager, s.cluster)
	prometheus.MustRegister(exporter)
	r.Handle(s.conf.PrometheusMetricPath, promhttp.Handler())
}

//ResponseType - Return content
type ResponseType struct {
	Code      int          `json:"code"`
	Message   string       `json:"msg"`
	MessageCN string       `json:"msgcn"`
	Body      ResponseBody `json:"body,omitempty"`
}

//ResponseBody - Back to the body
type ResponseBody struct {
	Bean     interface{}   `json:"bean,omitempty"`
	List     []interface{} `json:"list,omitempty"`
	PageNum  int           `json:"pageNumber,omitempty"`
	PageSize int           `json:"pageSize,omitempty"`
	Total    int           `json:"total,omitempty"`
}

//NewResponseType - Build the return structure
func NewResponseType(code int, message string, messageCN string, bean interface{}, list []interface{}) ResponseType {
	return ResponseType{
		Code:      code,
		Message:   message,
		MessageCN: messageCN,
		Body: ResponseBody{
			Bean: bean,
			List: list,
		},
	}
}

//NewSuccessResponse - Create successful return structure
func NewSuccessResponse(bean interface{}, list []interface{}) ResponseType {
	return NewResponseType(200, "", "", bean, list)
}
