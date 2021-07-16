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

package event

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/discover/config"
	eventclient "github.com/gridworkz/kato/eventlog/entry/grpc/client"
	eventpb "github.com/gridworkz/kato/eventlog/entry/grpc/pb"
	"github.com/gridworkz/kato/util"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

//Manager Operation log, client service
//Client load balancing
type Manager interface {
	GetLogger(eventID string) Logger
	Start() error
	Close() error
	ReleaseLogger(Logger)
}

// EventConfig event config struct
type EventConfig struct {
	EventLogServers []string
	DiscoverArgs    *etcdutil.ClientArgs
}
type manager struct {
	ctx            context.Context
	cancel         context.CancelFunc
	config         EventConfig
	qos            int32
	loggers        map[string]Logger
	handles        map[string]handle
	lock           sync.Mutex
	eventServer    []string
	abnormalServer map[string]string
	dis            discover.Discover
}

var defaultManager Manager

const (
	//REQUESTTIMEOUT  time out
	REQUESTTIMEOUT = 1000 * time.Millisecond
	//MAXRETRIES retry
	MAXRETRIES = 3 //  Before we abandon
	buffersize = 1000
)

//NewManager
func NewManager(conf EventConfig) error {
	dis, err := discover.GetDiscover(config.DiscoverConfig{EtcdClientArgs: conf.DiscoverArgs})
	if err != nil {
		logrus.Error("create discover manager error.", err.Error())
		if len(conf.EventLogServers) < 1 {
			return err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defaultManager = &manager{
		ctx:            ctx,
		cancel:         cancel,
		config:         conf,
		loggers:        make(map[string]Logger, 1024),
		handles:        make(map[string]handle),
		eventServer:    conf.EventLogServers,
		dis:            dis,
		abnormalServer: make(map[string]string),
	}
	return defaultManager.Start()
}

//GetManager
func GetManager() Manager {
	return defaultManager
}

// NewTestManager
func NewTestManager(m Manager) {
	defaultManager = m
}

//CloseManager
func CloseManager() {
	if defaultManager != nil {
		defaultManager.Close()
	}
}

func (m *manager) Start() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for i := 0; i < len(m.eventServer); i++ {
		h := handle{
			cacheChan: make(chan []byte, buffersize),
			stop:      make(chan struct{}),
			server:    m.eventServer[i],
			manager:   m,
			ctx:       m.ctx,
		}
		m.handles[m.eventServer[i]] = h
		go h.HandleLog()
	}
	if m.dis != nil {
		m.dis.AddProject("event_log_event_grpc", m)
	}
	go m.GC()
	return nil
}

func (m *manager) UpdateEndpoints(endpoints ...*config.Endpoint) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if endpoints == nil || len(endpoints) < 1 {
		return
	}
	//clear unavailable node information, focusing on service discovery
	m.abnormalServer = make(map[string]string)
	//add new node
	var new = make(map[string]string)
	for _, end := range endpoints {
		new[end.URL] = end.URL
		if _, ok := m.handles[end.URL]; !ok {
			h := handle{
				cacheChan: make(chan []byte, buffersize),
				stop:      make(chan struct{}),
				server:    end.URL,
				manager:   m,
				ctx:       m.ctx,
			}
			m.handles[end.URL] = h
			logrus.Infof("Add event server endpoint,%s", end.URL)
			go h.HandleLog()
		}
	}
	//delete old node
	for k := range m.handles {
		if _, ok := new[k]; !ok {
			delete(m.handles, k)
			logrus.Infof("Remove event server endpoint,%s", k)
		}
	}
	var eventServer []string
	for k := range new {
		eventServer = append(eventServer, k)
	}
	m.eventServer = eventServer
	logrus.Debugf("update event handle core success,handle core count:%d, event server count:%d", len(m.handles), len(m.eventServer))
}

func (m *manager) Error(err error) {

}
func (m *manager) Close() error {
	m.cancel()
	if m.dis != nil {
		m.dis.Stop()
	}
	return nil
}

func (m *manager) GC() {
	util.IntermittentExec(m.ctx, func() {
		m.lock.Lock()
		defer m.lock.Unlock()
		var needRelease []string
		for k, l := range m.loggers {
			//1min unreleased, automatic gc
			if l.CreateTime().Add(time.Minute).Before(time.Now()) {
				needRelease = append(needRelease, k)
			}
		}
		if len(needRelease) > 0 {
			for _, event := range needRelease {
				logrus.Infof("start auto release event logger. %s", event)
				delete(m.loggers, event)
			}
		}
	}, time.Second*20)
}

//GetLogger
//Must call the ReleaseLogger method after use
func (m *manager) GetLogger(eventID string) Logger {
	m.lock.Lock()
	defer m.lock.Unlock()
	if eventID == " " || len(eventID) == 0 {
		eventID = "system"
	}
	if l, ok := m.loggers[eventID]; ok {
		return l
	}
	l := NewLogger(eventID, m.getLBChan())
	m.loggers[eventID] = l
	return l
}

func (m *manager) ReleaseLogger(l Logger) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if l, ok := m.loggers[l.Event()]; ok {
		delete(m.loggers, l.Event())
	}
}

type handle struct {
	server    string
	stop      chan struct{}
	cacheChan chan []byte
	ctx       context.Context
	manager   *manager
}

func (m *manager) DiscardedLoggerChan(cacheChan chan []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for k, v := range m.handles {
		if v.cacheChan == cacheChan {
			logrus.Warnf("event server %s can not link, will ignore it.", k)
			m.abnormalServer[k] = k
		}
	}
	for _, v := range m.loggers {
		if v.GetChan() == cacheChan {
			v.SetChan(m.getLBChan())
		}
	}
}

func (m *manager) getLBChan() chan []byte {
	for i := 0; i < len(m.eventServer); i++ {
		index := m.qos % int32(len(m.eventServer))
		m.qos = atomic.AddInt32(&(m.qos), 1)
		server := m.eventServer[index]
		if _, ok := m.abnormalServer[server]; ok {
			logrus.Warnf("server[%s] is abnormal, skip it", server)
			continue
		}
		if h, ok := m.handles[server]; ok {
			return h.cacheChan
		}
		h := handle{
			cacheChan: make(chan []byte, buffersize),
			stop:      make(chan struct{}),
			server:    server,
			manager:   m,
			ctx:       m.ctx,
		}
		m.handles[server] = h
		go h.HandleLog()
		return h.cacheChan
	}
	//not selected, return first handle chan
	for _, v := range m.handles {
		return v.cacheChan
	}
	return nil
}
func (m *manager) RemoveHandle(server string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.handles[server]; ok {
		delete(m.handles, server)
	}
}
func (m *handle) HandleLog() error {
	defer m.manager.RemoveHandle(m.server)
	return util.Exec(m.ctx, func() error {
		ctx, cancel := context.WithCancel(m.ctx)
		defer cancel()
		client, err := eventclient.NewEventClient(ctx, m.server)
		if err != nil {
			logrus.Error("create event client error.", err.Error())
			return err
		}
		logrus.Infof("start a event log handle core. connect server %s", m.server)
		logClient, err := client.Log(ctx)
		if err != nil {
			logrus.Error("create event log client error.", err.Error())
			//Switch the logger using this chan to other chan
			m.manager.DiscardedLoggerChan(m.cacheChan)
			return err
		}
		for {
			select {
			case <-m.ctx.Done():
				logClient.CloseSend()
				return nil
			case <-m.stop:
				logClient.CloseSend()
				return nil
			case me := <-m.cacheChan:
				err := logClient.Send(&eventpb.LogMessage{Log: me})
				if err != nil {
					logrus.Error("send event log error.", err.Error())
					logClient.CloseSend()
					//Switch the logger using this chan to other chan
					m.manager.DiscardedLoggerChan(m.cacheChan)
					return nil
				}
			}
		}
	}, time.Second*3)
}

func (m *handle) Stop() {
	close(m.stop)
}

//Logger
type Logger interface {
	Info(string, map[string]string)
	Error(string, map[string]string)
	Debug(string, map[string]string)
	Event() string
	CreateTime() time.Time
	GetChan() chan []byte
	SetChan(chan []byte)
	GetWriter(step, level string) LoggerWriter
}

// NewLogger
func NewLogger(eventID string, sendCh chan []byte) Logger {
	return &logger{
		event:      eventID,
		sendChan:   sendCh,
		createTime: time.Now(),
	}
}

type logger struct {
	event      string
	sendChan   chan []byte
	createTime time.Time
}

func (l *logger) GetChan() chan []byte {
	return l.sendChan
}
func (l *logger) SetChan(ch chan []byte) {
	l.sendChan = ch
}
func (l *logger) Event() string {
	return l.event
}
func (l *logger) CreateTime() time.Time {
	return l.createTime
}
func (l *logger) Info(message string, info map[string]string) {
	if info == nil {
		info = make(map[string]string)
	}
	info["level"] = "info"
	l.send(message, info)
}
func (l *logger) Error(message string, info map[string]string) {
	if info == nil {
		info = make(map[string]string)
	}
	info["level"] = "error"
	l.send(message, info)
}
func (l *logger) Debug(message string, info map[string]string) {
	if info == nil {
		info = make(map[string]string)
	}
	info["level"] = "debug"
	l.send(message, info)
}
func (l *logger) send(message string, info map[string]string) {
	info["event_id"] = l.event
	info["message"] = message
	info["time"] = time.Now().Format(time.RFC3339)
	log, err := ffjson.Marshal(info)
	if err == nil && l.sendChan != nil {
		util.SendNoBlocking(log, l.sendChan)
	}
}

//LoggerWriter
type LoggerWriter interface {
	io.Writer
	SetFormat(map[string]interface{})
}

func (l *logger) GetWriter(step, level string) LoggerWriter {
	return &loggerWriter{
		l:     l,
		step:  step,
		level: level,
	}
}

type loggerWriter struct {
	l           *logger
	step        string
	level       string
	fmt         map[string]interface{}
	tmp         []byte
	lastMessage string
}

func (l *loggerWriter) SetFormat(f map[string]interface{}) {
	l.fmt = f
}
func (l *loggerWriter) Write(b []byte) (n int, err error) {
	if b != nil && len(b) > 0 {
		if !strings.HasSuffix(string(b), "\n") {
			l.tmp = append(l.tmp, b...)
			return len(b), nil
		}
		var message string
		if len(l.tmp) > 0 {
			message = string(append(l.tmp, b...))
			l.tmp = l.tmp[:0]
		} else {
			message = string(b)
		}
		// if loggerWriter has format, and then use it format message
		if len(l.fmt) > 0 {
			newLineMap := make(map[string]interface{}, len(l.fmt))
			for k, v := range l.fmt {
				if v == "%s" {
					newLineMap[k] = fmt.Sprintf(v.(string), message)
				} else {
					newLineMap[k] = v
				}
			}
			messageb, _ := ffjson.Marshal(newLineMap)
			message = string(messageb)
		}
		if l.step == "build-progress" {
			if strings.HasPrefix(message, "Progress ") && strings.HasPrefix(l.lastMessage, "Progress ") {
				l.lastMessage = message
				return len(b), nil
			}
			// send last message
			if !strings.HasPrefix(message, "Progress ") && strings.HasPrefix(l.lastMessage, "Progress ") {
				l.l.send(message, map[string]string{"step": l.lastMessage, "level": l.level})
			}
		}
		l.l.send(message, map[string]string{"step": l.step, "level": l.level})
		l.lastMessage = message
	}
	return len(b), nil
}

//GetTestLogger
func GetTestLogger() Logger {
	return &testLogger{}
}

type testLogger struct {
}

func (l *testLogger) GetChan() chan []byte {
	return nil
}
func (l *testLogger) SetChan(ch chan []byte) {

}
func (l *testLogger) Event() string {
	return "test"
}
func (l *testLogger) CreateTime() time.Time {
	return time.Now()
}
func (l *testLogger) Info(message string, info map[string]string) {
	fmt.Println("info:", message)
}
func (l *testLogger) Error(message string, info map[string]string) {
	fmt.Println("error:", message)
}
func (l *testLogger) Debug(message string, info map[string]string) {
	fmt.Println("debug:", message)
}

type testLoggerWriter struct {
}

func (l *testLoggerWriter) SetFormat(f map[string]interface{}) {

}
func (l *testLoggerWriter) Write(b []byte) (n int, err error) {
	return os.Stdout.Write(b)
}

func (l *testLogger) GetWriter(step, level string) LoggerWriter {
	return &testLoggerWriter{}
}
