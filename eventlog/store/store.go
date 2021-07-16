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

package store

import (
	"sync"
	"time"

	"github.com/gridworkz/kato/eventlog/db"

	"github.com/prometheus/client_golang/prometheus"

	"golang.org/x/net/context"
)

//MessageStore store
type MessageStore interface {
	InsertMessage(*db.EventLogMessage)
	InsertGarbageMessage(...*db.EventLogMessage)
	GetHistoryMessage(eventID string, length int) []string
	SubChan(eventID, subID string) chan *db.EventLogMessage
	RealseSubChan(eventID, subID string)
	GetMonitorData() *db.MonitorData
	Run()
	Gc()
	stop()
	Scrape(ch chan<- prometheus.Metric, namespace, exporter, from string) error
}

//NewStore
func NewStore(storeType string, manager *storeManager) MessageStore {
	ctx, cancel := context.WithCancel(context.Background())
	if storeType == "handle" {
		handle := &handleMessageStore{
			barrels:   make(map[string]*EventBarrel, 100),
			conf:      manager.conf,
			log:       manager.log.WithField("module", "HandleMessageStore"),
			garbageGC: make(chan int),
			ctx:       ctx,
			cancel:    cancel,
			//TODO:
			//If this channel is too small, it will block the insertion of received messages and cause deadlock
			//Change the persistent event to non-blocking insert
			barrelEvent:         make(chan []string, 100),
			dbPlugin:            manager.dbPlugin,
			handleEventCoreSize: 2,
			stopGarbage:         make(chan struct{}),
			manager:             manager,
		}
		handle.pool = &sync.Pool{
			New: func() interface{} {
				barrel := &EventBarrel{
					barrel:            make([]*db.EventLogMessage, 0),
					persistenceBarrel: make([]*db.EventLogMessage, 0),
					barrelEvent:       handle.barrelEvent,
					cacheNumber:       manager.conf.PeerEventMaxCacheLogNumber,
					maxNumber:         manager.conf.PeerEventMaxLogNumber,
				}
				return barrel
			},
		}
		go handle.handleGarbageMessage()
		for i := 0; i < handle.handleEventCoreSize; i++ {
			go handle.handleBarrelEvent()
		}
		return handle
	}
	if storeType == "read" {
		read := &readMessageStore{
			barrels: make(map[string]*readEventBarrel, 100),
			conf:    manager.conf,
			log:     manager.log.WithField("module", "SubMessageStore"),
			ctx:     ctx,
			cancel:  cancel,
		}
		read.pool = &sync.Pool{
			New: func() interface{} {
				reb := &readEventBarrel{
					subSocketChan: make(map[string]chan *db.EventLogMessage, 0),
				}
				return reb
			},
		}
		return read
	}
	if storeType == "docker_log" {
		docker := &dockerLogStore{
			barrels:    make(map[string]*dockerLogEventBarrel, 100),
			conf:       manager.conf,
			log:        manager.log.WithField("module", "DockerLogStore"),
			ctx:        ctx,
			cancel:     cancel,
			filePlugin: manager.filePlugin,
			//TODO:
			//If this channel is too small, it will block the insertion of received messages and cause deadlock
			//Change the persistent event to non-blocking insert
			barrelEvent: make(chan []string, 100),
		}
		docker.pool = &sync.Pool{
			New: func() interface{} {
				reb := &dockerLogEventBarrel{
					subSocketChan:   make(map[string]chan *db.EventLogMessage, 0),
					cacheSize:       manager.conf.PeerDockerMaxCacheLogNumber,
					barrelEvent:     docker.barrelEvent,
					persistenceTime: time.Now(),
				}
				return reb
			},
		}
		return docker
	}
	if storeType == "newmonitor" {
		monitor := &newMonitorMessageStore{
			barrels: make(map[string]*CacheMonitorMessageList, 100),
			conf:    manager.conf,
			log:     manager.log.WithField("module", "NewMonitorMessageStore"),
			ctx:     ctx,
			cancel:  cancel,
		}
		return monitor
	}

	return nil
}
