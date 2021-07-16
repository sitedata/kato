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

package logger

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	bufSize  = 16 * 1024
	readSize = 2 * 1024
)

// Copier can copy logs from specified sources to Logger and attach Timestamp.
// Writes are concurrent, so you need implement some sync in your logger.
type Copier struct {
	logfile *LogFile
	dst     []Logger
	closed  chan struct{}
	reader  *LogWatcher
	since   time.Time
	once    sync.Once
}

// NewCopier
func NewCopier(logfile *LogFile, dst []Logger, since time.Time) *Copier {
	return &Copier{
		logfile: logfile,
		reader:  NewLogWatcher(),
		dst:     dst,
		since:   since,
	}
}

// Run starts logs copying
func (c *Copier) Run() {
	c.closed = make(chan struct{})
	go c.logfile.ReadLogs(ReadConfig{Follow: true, Since: c.since, Tail: 0}, c.reader)
	go c.copySrc()
}

func (c *Copier) copySrc() {
	defer c.reader.ConsumerGone()
lool:
	for {
		select {
		case <-c.closed:
			return
		case err := <-c.reader.Err:
			logrus.Errorf("read container log file error %s, will retry after 5 seconds", err.Error())
			//If there is an error in the collection log process,
			//the collection should be restarted and not stopped
			time.Sleep(time.Second * 5)
			go c.logfile.ReadLogs(ReadConfig{Follow: true, Since: c.since, Tail: 0}, c.reader)
			continue
		case msg, ok := <-c.reader.Msg:
			if !ok {
				break lool
			}
			for _, d := range c.dst {
				if err := d.Log(msg); err != nil {
					logrus.Debugf("copy container log failure %s", err.Error())
				}
			}
		}
	}
}

// Close closes the copier
func (c *Copier) Close() {
	c.once.Do(func() {
		if c.dst != nil {
			for _, d := range c.dst {
				if err := d.Close(); err != nil {
					logrus.Errorf("close log driver failure %s", err.Error())
				}
			}
		}
		close(c.closed)
	})
}
