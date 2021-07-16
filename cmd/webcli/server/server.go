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
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gridworkz/kato/cmd/webcli/option"
	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/webcli/app"

	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/sirupsen/logrus"
)

//Run start
func Run(s *option.WebCliServer) error {
	errChan := make(chan error)
	option := app.DefaultOptions
	option.Address = s.Address
	option.Port = strconv.Itoa(s.Port)
	option.SessionKey = s.SessionKey
	option.K8SConfPath = s.K8SConfPath
	ap, err := app.New(&option)
	if err != nil {
		return err
	}
	err = ap.Run()
	if err != nil {
		return err
	}
	defer ap.Exit()
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: s.EtcdEndPoints,
		CaFile:    s.EtcdCaFile,
		CertFile:  s.EtcdCertFile,
		KeyFile:   s.EtcdKeyFile,
	}
	keepalive, err := discover.CreateKeepAlive(etcdClientArgs, "acp_webcli", s.HostName, s.HostIP, s.Port)
	if err != nil {
		return err
	}
	if err := keepalive.Start(); err != nil {
		return err
	}
	defer keepalive.Stop()
	//step finally: listen Signal
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		logrus.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-errChan:
		logrus.Errorf("Received a error %s, exiting gracefully...", err.Error())
	}
	logrus.Info("See you next time!")
	return nil
}
