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

package db

import (
	"fmt"

	"github.com/gridworkz/kato/eventlog/conf"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	SaveMessage([]*EventLogMessage) error
	Close() error
	GetMessages(id, level string, length int) (interface{}, error)
}

//NewManager - create a storage manager
func NewManager(conf conf.DBConf, log *logrus.Entry) (Manager, error) {
	switch conf.Type {
	case "file":
		return &filePlugin{
			homePath: conf.HomePath,
		}, nil
	case "eventfile":
		return &EventFilePlugin{
			HomePath: conf.HomePath,
		}, nil
	default:
		return nil, fmt.Errorf("plugin not supported")
	}
}
