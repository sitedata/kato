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

package testlog

import (
	"bytes"
	"fmt"

	"github.com/gridworkz/kato/node/nodem/logger"
	"github.com/sirupsen/logrus"
)

// Name is the name of the file that the jsonlogger logs to.
const Name = "test"

// TestLogger is Logger implementation for test
type TestLogger struct {
	buf *bytes.Buffer // json-encoded extra attributes
}

func init() {
	if err := logger.RegisterLogDriver(Name, New); err != nil {
		logrus.Fatal(err)
	}
	if err := logger.RegisterLogOptValidator(Name, ValidateLogOpt); err != nil {
		logrus.Fatal(err)
	}
}

// New creates new JSONFileLogger which writes to filename passed in
// on given context.
func New(info logger.Info) (logger.Logger, error) {
	logrus.Debugf("create logger driver for %s", info.ContainerName)
	return &TestLogger{
		buf: bytes.NewBuffer(nil),
	}, nil
}

// Log converts logger.Message to jsonlog.JSONLog and serializes it to file.
func (l *TestLogger) Log(msg *logger.Message) error {
	fmt.Println(string(msg.Line))
	return nil
}

// ValidateLogOpt looks for json specific log options max-file & max-size.
func ValidateLogOpt(cfg map[string]string) error {
	for key := range cfg {
		switch key {
		case "max-file":
		case "max-size":
		case "labels":
		case "env":
		default:
			return fmt.Errorf("unknown log opt '%s' for json-file log driver", key)
		}
	}
	return nil
}

// Close closes underlying file and signals all readers to stop.
func (l *TestLogger) Close() error {
	return nil
}

// Name returns name of this logger.
func (l *TestLogger) Name() string {
	return Name
}
