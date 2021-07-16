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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/cmd"
	envoyv2 "github.com/gridworkz/kato/node/core/envoy/v2"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		cmd.ShowVersion("sidecar")
	}
	if len(os.Args) > 1 && os.Args[1] == "wait" {
		var timeoutSeconds = 60
		var envoyReadyURL = "http://127.0.0.1:65533/ready"
		var envoyListennerReadyURL = "http://127.0.0.1:65533/listeners"
		var periodMillis = 500
		var requestTimeoutMillis = 500
		client := &http.Client{
			Timeout: time.Duration(requestTimeoutMillis) * time.Millisecond,
		}
		logrus.Infof("Waiting for Envoy proxy to be ready (timeout: %d seconds)...", timeoutSeconds)

		var err error
		timeoutAt := time.Now().Add(time.Duration(timeoutSeconds) * time.Second)
		for time.Now().Before(timeoutAt) {
			err = checkEnvoyIfReady(client, envoyReadyURL)
			if err == nil {
				logrus.Infof("Sidecar server is ready!")
				break
			}
			logrus.Infof("Not ready yet: %v", err)
			time.Sleep(time.Duration(periodMillis) * time.Millisecond)
		}
		if len(os.Args) > 2 && os.Args[2] != "0" {
			for time.Now().Before(timeoutAt) {
				err = checkEnvoyListenerIfReady(client, envoyListennerReadyURL, os.Args[2])
				if err == nil {
					logrus.Infof("Sidecar is ready!")
					os.Exit(0)
				}
				logrus.Infof("Not ready yet: %v", err)
				time.Sleep(time.Duration(periodMillis) * time.Millisecond)
			}
		} else {
			logrus.Infof("Sidecar is ready!")
			os.Exit(0)
		}
		logrus.Errorf("timeout waiting for Envoy proxy to become ready. Last error: %v", err)
		os.Exit(1)
	}
	if len(os.Args) > 1 && os.Args[1] == "run" {
		if err := run(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	loggerFile, _ := os.Create("/var/log/sidecar.log")
	if loggerFile != nil {
		defer loggerFile.Close()
		logrus.SetOutput(loggerFile)
	}
	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

//Run run
func Run() error {
	// start run first
	run()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	//step finally: listen Signal
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-term:
			logrus.Warn("Received SIGTERM, exiting gracefully...")
			return nil
		case <-ticker.C:
			run()
		}
	}
}

func run() error {
	configs := discoverConfig()
	if configs != nil {
		if hosts := getHosts(configs); hosts != nil {
			if err := writeHosts(hosts); err != nil {
				logrus.Errorf("write hosts failure %s", err.Error())
				return err
			}
			logrus.Debugf("rewrite hosts file success, %+v", hosts)
		}
	}
	return nil
}

func discoverConfig() *api_model.ResourceSpec {
	discoverURL := fmt.Sprintf("http://%s:6100/v1/resources/%s/%s/%s", os.Getenv("XDS_HOST_IP"), os.Getenv("TENANT_ID"), os.Getenv("SERVICE_NAME"), os.Getenv("PLUGIN_ID"))
	http.DefaultClient.Timeout = time.Second * 5
	res, err := http.Get(discoverURL)
	if err != nil {
		logrus.Errorf("get config failure %s", err.Error())
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close()
		var rs api_model.ResourceSpec
		if err := json.NewDecoder(res.Body).Decode(&rs); err != nil {
			logrus.Errorf("parse config body failure %s", err.Error())
		} else {
			return &rs
		}
	}
	return nil
}

func getHosts(configs *api_model.ResourceSpec) map[string]string {
	hosts := make(map[string]string)
	for _, service := range configs.BaseServices {
		options := envoyv2.GetOptionValues(service.Options)
		for _, domain := range options.Domains {
			if domain != "" && domain != "*" {
				if strings.Contains(domain, ":") {
					domain = strings.Split(domain, ":")[0]
				}
				hosts[domain] = "127.0.0.1"
			}
		}
	}
	if len(hosts) == 0 {
		return nil
	}
	return hosts
}

func writeHosts(ipnames map[string]string) error {
	hostFilePath := os.Getenv("HOST_FILE_PATH")
	if hostFilePath == "" {
		hostFilePath = "/etc/hosts"
	}
	hosts, err := util.NewHosts(hostFilePath)
	if err != nil {
		return err
	}
	for name, ip := range ipnames {
		hosts.Add(ip, name)
	}
	return hosts.Flush()
}

func checkEnvoyIfReady(client *http.Client, envoyReadyURL string) error {
	req, err := http.NewRequest(http.MethodGet, envoyReadyURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	reBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 || !strings.Contains(string(reBody), "LIVE") {
		return fmt.Errorf("HTTP status code %d, body: %s", resp.StatusCode, string(reBody))
	}
	return nil
}

func checkEnvoyListenerIfReady(client *http.Client, url string, port string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	reBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 || !strings.Contains(string(reBody), fmt.Sprintf(":%s", port)) {
		return fmt.Errorf("check Listeners HTTP status code %v, body is %s", resp.StatusCode, string(reBody))
	}
	return nil
}
