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
	"context"
	"encoding/json"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/gridworkz/kato/eventlog/db"
)

type newMonitorMessageStore struct {
	conf        conf.EventStoreConf
	log         *logrus.Entry
	barrels     map[string]*CacheMonitorMessageList
	lock        sync.RWMutex
	cancel      func()
	ctx         context.Context
	size        int64
	allLogCount float64
}

func (h *newMonitorMessageStore) Scrape(ch chan<- prometheus.Metric, namespace, exporter, from string) error {
	chanDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "new_monitor_store_barrel_count"),
		"the handle container log count size.",
		[]string{"from"}, nil,
	)
	ch <- prometheus.MustNewConstMetric(chanDesc, prometheus.GaugeValue, float64(len(h.barrels)), from)
	logDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, exporter, "new_monitor_store_log_count"),
		"the handle monitor log count size.",
		[]string{"from"}, nil,
	)
	ch <- prometheus.MustNewConstMetric(logDesc, prometheus.GaugeValue, h.allLogCount, from)

	return nil
}
func (h *newMonitorMessageStore) insertMessage(message *db.EventLogMessage) ([]MonitorMessage, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	mm := fromByte(message.MonitorData)
	if len(mm) < 1 {
		return mm, true
	}
	if mm[0].ServiceID == "" {
		return mm, true
	}
	if ba, ok := h.barrels[mm[0].ServiceID]; ok {
		ba.Insert(mm...)
		return mm, true
	}
	return mm, false
}

func (h *newMonitorMessageStore) InsertMessage(message *db.EventLogMessage) {
	if message == nil {
		return
	}
	h.size++
	h.allLogCount++
	mm, ok := h.insertMessage(message)
	if ok {
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	ba := CreateCacheMonitorMessageList(mm[0].ServiceID)
	ba.Insert(mm...)
	h.barrels[mm[0].ServiceID] = ba
}
func (h *newMonitorMessageStore) GetMonitorData() *db.MonitorData {
	data := &db.MonitorData{
		ServiceSize:  len(h.barrels),
		LogSizePeerM: h.size,
	}
	return data
}

func (h *newMonitorMessageStore) SubChan(eventID, subID string) chan *db.EventLogMessage {
	h.lock.Lock()
	defer h.lock.Unlock()
	if ba, ok := h.barrels[eventID]; ok {
		return ba.addSubChan(subID)
	}
	ba := CreateCacheMonitorMessageList(eventID)
	h.barrels[eventID] = ba
	return ba.addSubChan(subID)
}
func (h *newMonitorMessageStore) RealseSubChan(eventID, subID string) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if ba, ok := h.barrels[eventID]; ok {
		ba.delSubChan(subID)
	}
}
func (h *newMonitorMessageStore) Run() {
	go h.Gc()
}
func (h *newMonitorMessageStore) Gc() {
	tiker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-tiker.C:
		case <-h.ctx.Done():
			h.log.Debug("read message store gc stop.")
			tiker.Stop()
			return
		}
		h.size = 0
		if len(h.barrels) == 0 {
			continue
		}
		var gcEvent []string
		for k, v := range h.barrels {
			if len(v.subSocketChan) == 0 {
				if v.UpdateTime.Add(time.Minute * 3).Before(time.Now()) { // barrel - Timeout did not receive message
					gcEvent = append(gcEvent, k)
				}
			}
		}
		if gcEvent != nil && len(gcEvent) > 0 {
			for _, id := range gcEvent {
				h.log.Infof("monitor message barrel %s will be gc", id)
				barrel := h.barrels[id]
				barrel.empty()
				delete(h.barrels, id)
			}
		}
	}
}
func (h *newMonitorMessageStore) stop() {
	h.cancel()
}
func (h *newMonitorMessageStore) InsertGarbageMessage(message ...*db.EventLogMessage) {}
func (h *newMonitorMessageStore) GetHistoryMessage(eventID string, length int) (re []string) {
	return nil
}

//MonitorMessage - performance monitoring message system model
type MonitorMessage struct {
	ServiceID   string
	Port        string
	HostName    string
	MessageType string //mysql，http ...
	Key         string
	//total time
	CumulativeTime float64
	AverageTime    float64
	MaxTime        float64
	Count          uint64
	//Number of abnormal requests
	AbnormalCount uint64
}

//cacheMonitorMessage - Data cache for each instance
type cacheMonitorMessage struct {
	updateTime time.Time
	hostName   string
	mms        MonitorMessageList
}

//CacheMonitorMessageList - An application performance analysis data
type CacheMonitorMessageList struct {
	list          []*cacheMonitorMessage
	subSocketChan map[string]chan *db.EventLogMessage
	subLock       sync.Mutex
	message       db.EventLogMessage
	UpdateTime    time.Time
}

//CreateCacheMonitorMessageList - Create application monitoring information buffer
func CreateCacheMonitorMessageList(eventID string) *CacheMonitorMessageList {
	return &CacheMonitorMessageList{
		subSocketChan: make(map[string]chan *db.EventLogMessage),
		message: db.EventLogMessage{
			EventID: eventID,
		},
	}
}

//Insert - Think that the hostname of mms is consistent
//Gc every time a message is received
func (c *CacheMonitorMessageList) Insert(mms ...MonitorMessage) {
	if mms == nil || len(mms) < 1 {
		return
	}
	c.UpdateTime = time.Now()
	hostname := mms[0].HostName
	if len(c.list) == 0 {
		c.list = []*cacheMonitorMessage{
			&cacheMonitorMessage{
				updateTime: time.Now(),
				hostName:   hostname,
				mms:        mms,
			}}
	}
	var update bool
	for i := range c.list {
		cm := c.list[i]
		if cm.hostName == hostname {
			cm.updateTime = time.Now()
			cm.mms = mms
			update = true
			break
		}
	}
	if !update {
		c.list = append(c.list, &cacheMonitorMessage{
			updateTime: time.Now(),
			hostName:   hostname,
			mms:        mms,
		})
	}
	c.Gc()
	c.pushMessage()
}

//Gc - Clean data
func (c *CacheMonitorMessageList) Gc() {
	var list []*cacheMonitorMessage
	for i := range c.list {
		cmm := c.list[i]
		if !cmm.updateTime.Add(time.Second * 30).Before(time.Now()) {
			list = append(list, cmm)
		}
	}
	c.list = list
}

func (c *CacheMonitorMessageList) pushMessage() {
	if len(c.list) == 0 {
		return
	}
	var mdata []byte
	if len(c.list) == 1 {
		mdata = getByte(c.list[0].mms)
	}
	source := c.list[0].mms
	for i := 1; i < len(c.list); i++ {
		addSource := c.list[i].mms
		source = merge(source, addSource)
	}
	//Sort descending
	sort.Sort(sort.Reverse(&source))
	mdata = getByte(*source.Pop(20))
	c.message.MonitorData = mdata
	for _, ch := range c.subSocketChan {
		select {
		case ch <- &c.message:
		default:
		}
	}
}

// Add socket subscription
func (c *CacheMonitorMessageList) addSubChan(subID string) chan *db.EventLogMessage {
	c.subLock.Lock()
	defer c.subLock.Unlock()
	if sub, ok := c.subSocketChan[subID]; ok {
		return sub
	}
	ch := make(chan *db.EventLogMessage, 10)
	c.subSocketChan[subID] = ch
	c.pushMessage()
	return ch
}

//delSubChan delete socket sub chan
func (c *CacheMonitorMessageList) delSubChan(subID string) {
	c.subLock.Lock()
	defer c.subLock.Unlock()
	if ch, ok := c.subSocketChan[subID]; ok {
		close(ch)
		delete(c.subSocketChan, subID)
	}
}
func (c *CacheMonitorMessageList) empty() {
	c.subLock.Lock()
	defer c.subLock.Unlock()
	for _, v := range c.subSocketChan {
		close(v)
	}
}
func getByte(source []MonitorMessage) []byte {
	b, _ := json.Marshal(source)
	return b
}
func fromByte(source []byte) []MonitorMessage {
	var mm []MonitorMessage
	json.Unmarshal(source, &mm)
	return mm
}

func merge(source, addsource MonitorMessageList) (result MonitorMessageList) {
	var cache = make(map[string]MonitorMessage)
	for _, mm := range source {
		cache[mm.Key] = mm
	}
	for _, mm := range addsource {
		if oldmm, ok := cache[mm.Key]; ok {
			oldmm.Count += mm.Count
			oldmm.AbnormalCount += mm.AbnormalCount
			oldmm.AverageTime = Round((oldmm.AverageTime+mm.AverageTime)/2, 2)
			oldmm.CumulativeTime = Round(oldmm.CumulativeTime+mm.CumulativeTime, 2)
			if mm.MaxTime > oldmm.MaxTime {
				oldmm.MaxTime = mm.MaxTime
			}
			cache[mm.Key] = oldmm
			continue
		}
		cache[mm.Key] = mm
	}
	for _, c := range cache {
		result.Add(&c)
	}
	return
}

//Round
func Round(f float64, n int) float64 {
	pow10n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10n)*pow10n) / pow10n
}

//MonitorMessageList
type MonitorMessageList []MonitorMessage

//Add
func (m *MonitorMessageList) Add(mm *MonitorMessage) {
	*m = append(*m, *mm)
}

//Len - Is the total number of elements in the collection
func (m *MonitorMessageList) Len() int {
	return len(*m)
}

//Less - If the element with index i is smaller than the element with index j, return true, otherwise return false
func (m *MonitorMessageList) Less(i, j int) bool {
	return (*m)[i].CumulativeTime < (*m)[j].CumulativeTime
}

//Swap - Swap elements with index i and j
func (m *MonitorMessageList) Swap(i, j int) {
	tmp := (*m)[i]
	(*m)[i] = (*m)[j]
	(*m)[j] = tmp
}

//Pop
func (m *MonitorMessageList) Pop(i int) *MonitorMessageList {
	if len(*m) <= i {
		return m
	}
	cache := (*m)[:i]
	return &cache
}

//String json string
func (m *MonitorMessageList) String() string {
	body, _ := json.Marshal(m)
	return string(body)
}
