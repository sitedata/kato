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

package entry

import (
	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/store"

	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

//Entry - data entry
type Entry struct {
	supervisor   *suture.Supervisor
	log          *logrus.Entry
	conf         conf.EntryConf
	storeManager store.Manager
}

//NewEntry
func NewEntry(conf conf.EntryConf, log *logrus.Entry, storeManager store.Manager) *Entry {
	return &Entry{
		log:          log,
		conf:         conf,
		storeManager: storeManager,
	}
}

//Start
func (e *Entry) Start() error {
	supervisor := suture.New("Entry Server", suture.Spec{
		Log: func(m string) {
			e.log.Info(m)
		},
	})
	eventServer, err := NewEventLogServer(e.conf.EventLogServer, e.log.WithField("server", "EventLog"), e.storeManager)
	if err != nil {
		return err
	}
	dockerServer, err := NewDockerLogServer(e.conf.DockerLogServer, e.log.WithField("server", "DockerLog"), e.storeManager)
	if err != nil {
		return err
	}
	monitorServer, err := NewMonitorMessageServer(e.conf.MonitorMessageServer, e.log.WithField("server", "MonitorMessage"), e.storeManager)
	if err != nil {
		return err
	}
	newmonitorServer, err := NewNMonitorMessageServer(e.conf.NewMonitorMessageServerConf, e.log.WithField("server", "NewMonitorMessage"), e.storeManager)
	if err != nil {
		return err
	}

	supervisor.Add(eventServer)
	supervisor.Add(dockerServer)
	supervisor.Add(monitorServer)
	supervisor.Add(newmonitorServer)
	supervisor.ServeBackground()
	e.supervisor = supervisor
	return nil
}

//Stop
func (e *Entry) Stop() {
	if e.supervisor != nil {
		e.supervisor.Stop()
	}
}
