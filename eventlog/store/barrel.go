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
	"errors"
	"sync"
	"time"

	"github.com/gridworkz/kato/eventlog/db"

	"github.com/sirupsen/logrus"
)

//EventBarrel
//Coroutine cannot be started on this structure
type EventBarrel struct {
	eventID           string
	barrel            []*db.EventLogMessage
	persistenceBarrel []*db.EventLogMessage
	needPersistence   bool
	persistencelock   sync.Mutex
	barrelEvent       chan []string
	maxNumber         int64
	cacheNumber       int
	size              int64
	isCallback        bool
	updateTime        time.Time
}

//insert message
func (e *EventBarrel) insert(m *db.EventLogMessage) error {
	if e.size > e.maxNumber {
		return errors.New("received message number more than peer event max message number")
	}
	if m.Step == "progress" { //Progress log is not stored
		return nil
	}
	e.barrel = append(e.barrel, m)
	e.size++      //Total number of message data
	e.analysis(m) //Synchronous analysis
	if len(e.barrel) >= e.cacheNumber {
		e.persistence()
	}
	e.updateTime = time.Now()
	return nil
}

//Return the value to zero and return to the object pool
func (e *EventBarrel) empty() {
	e.size = 0
	e.eventID = ""
	e.barrel = e.barrel[:0]
	e.persistenceBarrel = e.persistenceBarrel[:0]
	e.needPersistence = false
	e.isCallback = false
}

//Real-time analysis
func (e *EventBarrel) analysis(newMessage *db.EventLogMessage) {
	if newMessage.Step == "last" || newMessage.Step == "callback" {
		e.persistence()
		e.barrelEvent <- []string{"callback", e.eventID, newMessage.Status, newMessage.Message}
	}
	if newMessage.Step == "code-version" {
		e.barrelEvent <- []string{"code-version", e.eventID, newMessage.Message}
	}
}

//persistence
func (e *EventBarrel) persistence() {
	e.persistencelock.Lock()
	defer e.persistencelock.Unlock()
	e.needPersistence = true
	e.persistenceBarrel = append(e.persistenceBarrel, e.barrel...) //Data goes to the persistent waiting queue
	e.barrel = e.barrel[:0]                                        //Cache queue empty
	select {
	case e.barrelEvent <- []string{"persistence", e.eventID}: //Issue a persistent command
	default:
		logrus.Debug("event message log persistence delay")
	}
}

//Caller lock
func (e *EventBarrel) gcPersistence() {
	e.needPersistence = true
	e.persistenceBarrel = append(e.persistenceBarrel, e.barrel...) //Data goes to the persistent waiting queue
	e.barrel = nil
}

type readEventBarrel struct {
	barrel        []*db.EventLogMessage
	subSocketChan map[string]chan *db.EventLogMessage
	subLock       sync.Mutex
	updateTime    time.Time
}

func (r *readEventBarrel) empty() {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if r.barrel != nil {
		r.barrel = r.barrel[:0]
	}
	//Close subscription chan
	for _, ch := range r.subSocketChan {
		close(ch)
	}
	r.subSocketChan = make(map[string]chan *db.EventLogMessage)
}

func (r *readEventBarrel) insertMessage(message *db.EventLogMessage) {
	r.barrel = append(r.barrel, message)
	r.updateTime = time.Now()
	r.subLock.Lock()
	defer r.subLock.Unlock()
	for _, v := range r.subSocketChan { //Send messages to subscribed channels
		select {
		case v <- message:
		default:
		}
	}
}

func (r *readEventBarrel) pushCashMessage(ch chan *db.EventLogMessage, subID string) {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	//send cache message
	for _, m := range r.barrel {
		ch <- m
	}
	r.subSocketChan[subID] = ch
}

//Add socket subscription
func (r *readEventBarrel) addSubChan(subID string) chan *db.EventLogMessage {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if sub, ok := r.subSocketChan[subID]; ok {
		return sub
	}
	ch := make(chan *db.EventLogMessage, 10)
	go r.pushCashMessage(ch, subID)
	return ch
}

//Delete socket subscription
func (r *readEventBarrel) delSubChan(subID string) {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if ch, ok := r.subSocketChan[subID]; ok {
		close(ch)
		delete(r.subSocketChan, subID)
	}
}

type dockerLogEventBarrel struct {
	name              string
	barrel            []*db.EventLogMessage
	subSocketChan     map[string]chan *db.EventLogMessage
	subLock           sync.Mutex
	updateTime        time.Time
	size              int
	cacheSize         int64
	persistencelock   sync.Mutex
	persistenceTime   time.Time
	needPersistence   bool
	persistenceBarrel []*db.EventLogMessage
	barrelEvent       chan []string
}

func (r *dockerLogEventBarrel) empty() {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if r.barrel != nil {
		r.barrel = r.barrel[:0]
	}
	for _, ch := range r.subSocketChan {
		close(ch)
	}
	r.subSocketChan = make(map[string]chan *db.EventLogMessage)
	r.size = 0
	r.name = ""
	r.persistenceBarrel = r.persistenceBarrel[:0]
	r.needPersistence = false
}

func (r *dockerLogEventBarrel) insertMessage(message *db.EventLogMessage) {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	r.barrel = append(r.barrel, message)
	if r.name == "" {
		r.name = message.EventID
	}
	r.updateTime = time.Now()
	for _, v := range r.subSocketChan { //Send messages to subscribed channels
		select {
		case v <- message:
		default:
		}
	}
	r.size++
	if int64(len(r.barrel)) >= r.cacheSize {
		r.persistence()
	}
}

func (r *dockerLogEventBarrel) pushCashMessage(ch chan *db.EventLogMessage, subID string) {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	r.subSocketChan[subID] = ch
}

//Add socket subscription
func (r *dockerLogEventBarrel) addSubChan(subID string) chan *db.EventLogMessage {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if sub, ok := r.subSocketChan[subID]; ok {
		return sub
	}
	ch := make(chan *db.EventLogMessage, 100)
	go r.pushCashMessage(ch, subID)
	return ch
}

//Delete socket subscription
func (r *dockerLogEventBarrel) delSubChan(subID string) {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if ch, ok := r.subSocketChan[subID]; ok {
		close(ch)
		delete(r.subSocketChan, subID)
	}
}

//persistence
func (r *dockerLogEventBarrel) persistence() {
	r.persistencelock.Lock()
	defer r.persistencelock.Unlock()
	r.needPersistence = true
	r.persistenceBarrel = append(r.persistenceBarrel, r.barrel...) //Data goes to the persistent waiting queue
	r.barrel = r.barrel[:0]                                        //Cache queue empty
	select {
	case r.barrelEvent <- []string{"persistence", r.name}: //Issue a persistent command
		r.persistenceTime = time.Now()
	default:
		logrus.Errorln("docker log persistence delay")
	}
}

//Caller lock
func (r *dockerLogEventBarrel) gcPersistence() {
	r.needPersistence = true
	r.persistenceBarrel = append(r.persistenceBarrel, r.barrel...) //Data goes to the persistent waiting queue
	r.barrel = nil
}
func (r *dockerLogEventBarrel) GetSubChanLength() int {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	return len(r.subSocketChan)
}

type monitorMessageBarrel struct {
	barrel        []*db.EventLogMessage
	subSocketChan map[string]chan *db.EventLogMessage
	subLock       sync.Mutex
	updateTime    time.Time
}

func (r *monitorMessageBarrel) empty() {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if r.barrel != nil {
		r.barrel = r.barrel[:0]
	}
	for _, ch := range r.subSocketChan {
		close(ch)
	}
	r.subSocketChan = make(map[string]chan *db.EventLogMessage)
}

func (r *monitorMessageBarrel) insertMessage(message *db.EventLogMessage) {
	r.updateTime = time.Now()
	r.subLock.Lock()
	defer r.subLock.Unlock()
	for _, v := range r.subSocketChan { //Send messages to subscribed channels
		select {
		case v <- message:
		default:
		}
	}
}

func (r *monitorMessageBarrel) pushCashMessage(ch chan *db.EventLogMessage, subID string) {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	r.subSocketChan[subID] = ch
}

//Add socket subscription
func (r *monitorMessageBarrel) addSubChan(subID string) chan *db.EventLogMessage {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if sub, ok := r.subSocketChan[subID]; ok {
		return sub
	}
	ch := make(chan *db.EventLogMessage, 10)
	go r.pushCashMessage(ch, subID)
	return ch
}

//Delete socket subscription
func (r *monitorMessageBarrel) delSubChan(subID string) {
	r.subLock.Lock()
	defer r.subLock.Unlock()
	if ch, ok := r.subSocketChan[subID]; ok {
		close(ch)
		delete(r.subSocketChan, subID)
	}
}
