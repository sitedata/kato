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

	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/db"

	"golang.org/x/net/context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type readMessageStore struct {
	conf    conf.EventStoreConf
	log     *logrus.Entry
	barrels map[string]*readEventBarrel
	lock    sync.Mutex
	cancel  func()
	ctx     context.Context
	pool    *sync.Pool
}

func (h *readMessageStore) Scrape(ch chan<- prometheus.Metric, namespace, exporter, from string) error {
	return nil
}
func (h *readMessageStore) InsertMessage(message *db.EventLogMessage) {
	if message == nil || message.EventID == "" {
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	if ba, ok := h.barrels[message.EventID]; ok {
		ba.insertMessage(message)
	} else {
		ba := h.pool.Get().(*readEventBarrel)
		ba.insertMessage(message)
		h.barrels[message.EventID] = ba
	}
}
func (h *readMessageStore) GetMonitorData() *db.MonitorData {
	return nil
}

func (h *readMessageStore) SubChan(eventID, subID string) chan *db.EventLogMessage {
	h.lock.Lock()
	defer h.lock.Unlock()
	if ba, ok := h.barrels[eventID]; ok {
		return ba.addSubChan(subID)
	}
	ba := h.pool.Get().(*readEventBarrel)
	ba.updateTime = time.Now()
	h.barrels[eventID] = ba
	return ba.addSubChan(subID)
}
func (h *readMessageStore) RealseSubChan(eventID, subID string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if ba, ok := h.barrels[eventID]; ok {
		ba.delSubChan(subID)
	}
}
func (h *readMessageStore) Run() {
	go h.Gc()
}
func (h *readMessageStore) Gc() {
	tiker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-tiker.C:
		case <-h.ctx.Done():
			h.log.Debug("read message store gc stop.")
			tiker.Stop()
			return
		}
		if len(h.barrels) == 0 {
			continue
		}
		var gcEvent []string
		for k, v := range h.barrels {
			if v.updateTime.Add(time.Minute * 2).Before(time.Now()) { // barrel - Timeout did not receive message
				gcEvent = append(gcEvent, k)
			}
		}
		if gcEvent != nil && len(gcEvent) > 0 {
			for _, id := range gcEvent {
				barrel := h.barrels[id]
				barrel.empty()
				h.pool.Put(barrel) // Return to the object pool
				delete(h.barrels, id)
			}
		}
	}
}
func (h *readMessageStore) stop() {
	h.cancel()
}
func (h *readMessageStore) InsertGarbageMessage(message ...*db.EventLogMessage) {}

func (h *readMessageStore) GetHistoryMessage(eventID string, length int) (re []string) {
	return nil
}
