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
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestReadFile(t *testing.T) {
	reader, err := NewLogFile("../../../test/dockerlog/tes.log", 3, false, decodeFunc, 0640, getTailReader)
	if err != nil {
		t.Fatal(err)
	}
	watch := NewLogWatcher()
	reader.ReadLogs(ReadConfig{Follow: true, Tail: 0, Since: time.Now()}, watch)
	defer watch.ConsumerGone()
LogLoop:
	for {
		select {
		case msg, ok := <-watch.Msg:
			if !ok {
				break LogLoop
			}
			fmt.Println(string(msg.Line))
		case err := <-watch.Err:
			logrus.Errorf("Error streaming logs: %v", err)
			break LogLoop
		}
	}
}
