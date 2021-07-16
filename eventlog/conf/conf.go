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

package conf

import "time"

// Conf
type Conf struct {
	Entry       EntryConf
	EventStore  EventStoreConf
	Log         LogConf
	WebSocket   WebSocketConf
	WebHook     WebHookConf
	ClusterMode bool
	Cluster     ClusterConf
	Kubernetes  KubernetsConf
}

// WebHookConf
type WebHookConf struct {
	ConsoleURL   string
	ConsoleToken string
}

// DBConf
type DBConf struct {
	Type        string
	URL         string
	PoolSize    int
	PoolMaxSize int
	HomePath    string
}

// WebSocketConf
type WebSocketConf struct {
	BindIP               string
	BindPort             int
	SSLBindPort          int
	EnableCompression    bool
	ReadBufferSize       int
	WriteBufferSize      int
	MaxRestartCount      int
	TimeOut              string
	SSL                  bool
	CertFile             string
	KeyFile              string
	PrometheusMetricPath string
}

// LogConf
type LogConf struct {
	LogLevel   string
	LogOutType string
	LogPath    string
}

// EntryConf
type EntryConf struct {
	EventLogServer              EventLogServerConf
	DockerLogServer             DockerLogServerConf
	MonitorMessageServer        MonitorMessageServerConf
	NewMonitorMessageServerConf NewMonitorMessageServerConf
}

// EventLogServerConf
type EventLogServerConf struct {
	BindIP           string
	BindPort         int
	CacheMessageSize int
}

// DockerLogServerConf
type DockerLogServerConf struct {
	BindIP           string
	BindPort         int
	CacheMessageSize int
	Mode             string
}

// DiscoverConf
type DiscoverConf struct {
	Type          string
	EtcdAddr      []string
	EtcdCaFile    string
	EtcdCertFile  string
	EtcdKeyFile   string
	EtcdUser      string
	EtcdPass      string
	ClusterMode   bool
	InstanceIP    string
	HomePath      string
	DockerLogPort int
	WebPort       int
	NodeID        string
}

// PubSubConf
type PubSubConf struct {
	PubBindIP      string
	PubBindPort    int
	ClusterMode    bool
	PollingTimeout time.Duration
}

// EventStoreConf
type EventStoreConf struct {
	EventLogPersistenceLength   int64
	MessageType                 string
	GarbageMessageSaveType      string
	GarbageMessageFile          string
	PeerEventMaxLogNumber       int64 //Maximum number of logs per event
	PeerEventMaxCacheLogNumber  int
	PeerDockerMaxCacheLogNumber int64
	ClusterMode                 bool
	HandleMessageCoreNumber     int
	HandleSubMessageCoreNumber  int
	HandleDockerLogCoreNumber   int
	DB                          DBConf
}

// KubernetsConf
type KubernetsConf struct {
	Master string
}

// ClusterConf
type ClusterConf struct {
	PubSub   PubSubConf
	Discover DiscoverConf
}

// MonitorMessageServerConf
type MonitorMessageServerConf struct {
	SubAddress       []string
	SubSubscribe     string
	CacheMessageSize int
}

// NewMonitorMessageServerConf
type NewMonitorMessageServerConf struct {
	ListenerHost string
	ListenerPort int
}
