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

package option

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Config
type Config struct {
	EtcdEndpointsLine    string
	EtcdEndpoints        []string
	EtcdCaFile           string
	EtcdCertFile         string
	EtcdKeyFile          string
	LogLevel             string
	AdvertiseAddr        string
	BindIP               string
	Port                 int
	StartArgs            []string
	ConfigFile           string
	AlertingRulesFile    string
	AlertManagerURL      []string
	LocalStoragePath     string
	Web                  Web
	Tsdb                 Tsdb
	WebTimeout           string
	RemoteFlushDeadline  string
	AlertmanagerCapacity string
	AlertmanagerTimeout  string
	QueryLookbackDelta   string
	QueryTimeout         string
	QueryMaxConcurrency  string
	CadvisorListenPort   int
	MysqldExporter       string
	KSMExporter          string
	KubeConfig           string
}

// Web Options for the web Handler.
type Web struct {
	ListenAddress        string
	ReadTimeout          time.Duration
	MaxConnections       int
	ExternalURL          string
	RoutePrefix          string
	UseLocalAssets       bool
	UserAssetsPath       string
	ConsoleTemplatesPath string
	ConsoleLibrariesPath string
	EnableLifecycle      bool
	EnableAdminAPI       bool
}

// Tsdb Options of the DB storage.
type Tsdb struct {
	// The interval at which the write ahead log is flushed to disk.
	WALFlushInterval time.Duration

	// The timestamp range of head blocks after which they get persisted.
	// It's the minimum duration of any persisted block.
	MinBlockDuration string

	// The maximum timestamp range of compacted blocks.
	MaxBlockDuration string

	// Duration for how long to retain data.
	Retention string

	// Disable creation and consideration of lockfile.
	NoLockfile bool
}

// NewConfig
func NewConfig() *Config {
	host, _ := os.Hostname()

	config := &Config{
		EtcdEndpointsLine:    "http://127.0.0.1:2379",
		EtcdEndpoints:        []string{},
		AdvertiseAddr:        host + ":9999",
		BindIP:               host,
		Port:                 9999,
		LogLevel:             "info",
		KubeConfig:           "",
		ConfigFile:           "/etc/prometheus/prometheus.yml",
		AlertingRulesFile:    "/etc/prometheus/rules.yml",
		AlertManagerURL:      []string{},
		LocalStoragePath:     "/prometheusdata",
		WebTimeout:           "5m",
		RemoteFlushDeadline:  "1m",
		AlertmanagerCapacity: "10000",
		AlertmanagerTimeout:  "10s",
		QueryLookbackDelta:   "5m",
		QueryTimeout:         "2m",
		QueryMaxConcurrency:  "20",
		Web: Web{
			ListenAddress:        "0.0.0.0:9999",
			ReadTimeout:          time.Minute * 5,
			MaxConnections:       512,
			ConsoleTemplatesPath: "consoles",
			ConsoleLibrariesPath: "console_libraries",
		},
		Tsdb: Tsdb{
			MinBlockDuration: "2h",
			Retention:        "7d",
		},
		CadvisorListenPort: 10250,
	}

	return config
}

//AddFlag monitor flag
func (c *Config) AddFlag(cmd *pflag.FlagSet) {
	cmd.StringVar(&c.EtcdEndpointsLine, "etcd-endpoints", c.EtcdEndpointsLine, "etcd endpoints list.")
	cmd.StringVar(&c.EtcdCaFile, "etcd-ca", "", "etcd tls ca file ")
	cmd.StringVar(&c.EtcdCertFile, "etcd-cert", "", "etcd tls cert file")
	cmd.StringVar(&c.EtcdKeyFile, "etcd-key", "", "etcd http tls cert key file")
	cmd.StringVar(&c.AdvertiseAddr, "advertise-addr", c.AdvertiseAddr, "advertise address, and registry into etcd.")
	cmd.IntVar(&c.CadvisorListenPort, "cadvisor-listen-port", c.CadvisorListenPort, "kubelet cadvisor listen port in all node")
	cmd.StringSliceVar(&c.AlertManagerURL, "alertmanager-address", c.AlertManagerURL, "AlertManager url.")
	cmd.StringVar(&c.MysqldExporter, "mysqld-exporter", c.MysqldExporter, "mysqld exporter address. eg: 127.0.0.1:9104")
	cmd.StringVar(&c.KSMExporter, "kube-state-metrics", c.KSMExporter, "kube-state-metrics, current server's kube-state-metrics address")
	cmd.StringVar(&c.KubeConfig, "kube-config", "", "kubernetes api server config file")
}

//AddPrometheusFlag
func (c *Config) AddPrometheusFlag(cmd *pflag.FlagSet) {
	cmd.StringVar(&c.ConfigFile, "config.file", c.ConfigFile, "Prometheus configuration file path.")

	cmd.StringVar(&c.AlertingRulesFile, "rules-config.file", c.AlertingRulesFile, "Prometheus alerting rules config file path.")

	cmd.StringVar(&c.Web.ListenAddress, "web.listen-address", c.Web.ListenAddress, "Address to listen on for UI, API, and telemetry.")

	cmd.StringVar(&c.WebTimeout, "web.read-timeout", c.WebTimeout, "Maximum duration before timing out read of the request, and closing idle connections.")

	cmd.IntVar(&c.Web.MaxConnections, "web.max-connections", c.Web.MaxConnections, "Maximum number of simultaneous connections.")

	cmd.StringVar(&c.Web.ExternalURL, "web.external-url", c.Web.ExternalURL, "The URL under which Prometheus is externally reachable (for example, if Prometheus is served via a reverse proxy). Used for generating relative and absolute links back to Prometheus itself. If the URL has a path portion, it will be used to prefix all HTTP endpoints served by Prometheus. If omitted, relevant URL components will be derived automatically.")

	cmd.StringVar(&c.Web.RoutePrefix, "web.route-prefix", c.Web.RoutePrefix, "Prefix for the internal routes of Web endpoints. Defaults to path of --Web.external-url.")

	cmd.StringVar(&c.Web.UserAssetsPath, "web.user-assets", c.Web.UserAssetsPath, "Path to static asset directory, available at /user.")

	cmd.BoolVar(&c.Web.EnableLifecycle, "web.enable-lifecycle", c.Web.EnableLifecycle, "Enable shutdown and reload via HTTP request.")

	cmd.BoolVar(&c.Web.EnableAdminAPI, "web.enable-admin-api", c.Web.EnableAdminAPI, "Enable API endpoints for admin control actions.")

	cmd.StringVar(&c.Web.ConsoleTemplatesPath, "web.console.templates", c.Web.ConsoleTemplatesPath, "Path to the console template directory, available at /consoles.")

	cmd.StringVar(&c.Web.ConsoleLibrariesPath, "web.console.libraries", c.Web.ConsoleLibrariesPath, "Path to the console library directory.")

	cmd.StringVar(&c.LocalStoragePath, "storage.tsdb.path", c.LocalStoragePath, "Base path for metrics storage.")

	cmd.StringVar(&c.Tsdb.MinBlockDuration, "storage.tsdb.min-block-duration", c.Tsdb.MinBlockDuration, "Minimum duration of a data block before being persisted. For use in testing.")

	cmd.StringVar(&c.Tsdb.MaxBlockDuration, "storage.tsdb.max-block-duration", c.Tsdb.MaxBlockDuration,
		"Maximum duration compacted blocks may span. For use in testing. (Defaults to 10% of the retention period).")

	cmd.StringVar(&c.Tsdb.Retention, "storage.tsdb.retention", c.Tsdb.Retention, "How long to retain samples in storage.")

	cmd.BoolVar(&c.Tsdb.NoLockfile, "storage.tsdb.no-lockfile", c.Tsdb.NoLockfile, "Do not create lockfile in data directory.")

	cmd.StringVar(&c.RemoteFlushDeadline, "storage.remote.flush-deadline", c.RemoteFlushDeadline, "How long to wait flushing sample on shutdown or config reload.")

	cmd.StringVar(&c.AlertmanagerCapacity, "alertmanager.notification-queue-capacity", c.AlertmanagerCapacity, "The capacity of the queue for pending Alertmanager notifications.")

	cmd.StringVar(&c.AlertmanagerTimeout, "alertmanager.timeout", c.AlertmanagerTimeout, "Timeout for sending alerts to Alertmanager.")

	cmd.StringVar(&c.QueryLookbackDelta, "query.lookback-delta", c.QueryLookbackDelta, "The delta difference allowed for retrieving metrics during expression evaluations.")

	cmd.StringVar(&c.QueryTimeout, "query.timeout", c.QueryTimeout, "Maximum time a query may take before being aborted.")

	cmd.StringVar(&c.QueryMaxConcurrency, "query.max-concurrency", c.QueryMaxConcurrency, "Maximum number of queries executed concurrently.")

	cmd.StringVar(&c.LogLevel, "log.level", c.LogLevel, "log level.")
}

// CompleteConfig
func (c *Config) CompleteConfig() {
	// parse etcd urls line to array
	for _, url := range strings.Split(c.EtcdEndpointsLine, ",") {
		c.EtcdEndpoints = append(c.EtcdEndpoints, url)
	}

	if len(c.EtcdEndpoints) < 1 {
		logrus.Error("Must define the etcd endpoints by --etcd-endpoints")
		os.Exit(17)
	}

	// parse values from prometheus options to config
	ipPort := strings.TrimLeft(c.AdvertiseAddr, "shttp://")
	ipPortArr := strings.Split(ipPort, ":")
	c.BindIP = ipPortArr[0]
	port, err := strconv.Atoi(ipPortArr[1])
	if err == nil {
		c.Port = port
	}

	defaultOptions := "--log.level=%s --web.listen-address=%s --config.file=%s --storage.tsdb.path=%s --storage.tsdb.retention.time=%s"
	defaultOptions = fmt.Sprintf(defaultOptions, c.LogLevel, c.Web.ListenAddress, c.ConfigFile, c.LocalStoragePath, c.Tsdb.Retention)
	if c.Tsdb.NoLockfile {
		defaultOptions += " --storage.tsdb.no-lockfile"
	}
	if c.Web.EnableAdminAPI {
		defaultOptions += " --web.enable-admin-api"
	}
	if c.Web.EnableLifecycle {
		defaultOptions += " --web.enable-lifecycle"
	}
	if c.Web.ExternalURL != "" {
		defaultOptions += fmt.Sprintf(" --web.external-url=%s", c.Web.ExternalURL)
	}

	args := strings.Split(defaultOptions, " ")
	c.StartArgs = append(c.StartArgs, args...)

	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		fmt.Println("ERROR set log level:", err)
		return
	}
	logrus.SetLevel(level)

	logrus.Info("Start with options: ", c)
}
