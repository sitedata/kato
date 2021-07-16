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

package envoy

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

//Event
type Event int

const (
	// EventAdd is sent when an object is added
	EventAdd Event = iota

	// EventUpdate is sent when an object is modified
	// Captures the modified object
	EventUpdate

	// EventDelete is sent when an object is deleted
	// Captures the object at the last known state
	EventDelete
)

// Queue of work tickets processed using a rate-limiting loop
type Queue interface {
	// Push a ticket
	Push(Task)
	// Run the loop until a signal on the channel
	Run(<-chan struct{})
}

// Handler specifies a function to apply on an object for a given event type
type Handler func(obj interface{}, event Event) error

// Task object for the event watchers; processes until handler succeeds
type Task struct {
	handler Handler
	obj     interface{}
	event   Event
}

// NewTask creates a task from a work item
func NewTask(handler Handler, obj interface{}, event Event) Task {
	return Task{handler: handler, obj: obj, event: event}
}

type queueImpl struct {
	delay   time.Duration
	queue   []Task
	cond    *sync.Cond
	closing bool
}

// NewQueue instantiates a queue with a processing function
func NewQueue(errorDelay time.Duration) Queue {
	return &queueImpl{
		delay:   errorDelay,
		queue:   make([]Task, 0),
		closing: false,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (q *queueImpl) Push(item Task) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if !q.closing {
		q.queue = append(q.queue, item)
	}
	q.cond.Signal()
}

func (q *queueImpl) Run(stop <-chan struct{}) {
	go func() {
		<-stop
		q.cond.L.Lock()
		q.closing = true
		q.cond.L.Unlock()
	}()

	for {
		q.cond.L.Lock()
		for !q.closing && len(q.queue) == 0 {
			q.cond.Wait()
		}

		if len(q.queue) == 0 {
			q.cond.L.Unlock()
			// We must be shutting down.
			return
		}

		var item Task
		item, q.queue = q.queue[0], q.queue[1:]
		q.cond.L.Unlock()

		if err := item.handler(item.obj, item.event); err != nil {
			logrus.Infof("Work item handle failed (%v), retry after delay %v", err, q.delay)
			time.AfterFunc(q.delay, func() {
				q.Push(item)
			})
		}

	}
}

// ChainHandler applies handlers in a sequence
type ChainHandler struct {
	funcs []Handler
}

// Apply is the handler function
func (ch *ChainHandler) Apply(obj interface{}, event Event) error {
	for _, f := range ch.funcs {
		if err := f(obj, event); err != nil {
			return err
		}
	}
	return nil
}

// Append a handler as the last handler in the chain
func (ch *ChainHandler) Append(h Handler) {
	ch.funcs = append(ch.funcs, h)
}
