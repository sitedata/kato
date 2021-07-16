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
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	eventutil "github.com/gridworkz/kato/eventlog/util"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
)

//EventFilePlugin
type EventFilePlugin struct {
	HomePath string
}

//SaveMessage
func (m *EventFilePlugin) SaveMessage(events []*EventLogMessage) error {
	if len(events) == 0 {
		return nil
	}
	filePath := eventutil.EventLogFilePath(m.HomePath)
	if err := util.CheckAndCreateDir(filePath); err != nil {
		return err
	}
	filename := eventutil.EventLogFileName(filePath, events[0].EventID)
	writeFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer writeFile.Close()
	var lastTime int64
	for _, e := range events {
		if e == nil {
			continue
		}
		writeFile.Write(GetLevelFlag(e.Level))
		logtime := GetTimeUnix(e.Time)
		if logtime != 0 {
			lastTime = logtime
		}
		writeFile.Write([]byte(fmt.Sprintf("%d ", lastTime)))
		writeFile.Write([]byte(e.Message))
		writeFile.Write([]byte("\n"))
	}
	return nil
}

//MessageData message data - obtain the operation log of the specified operation
type MessageData struct {
	Message  string `json:"message"`
	Time     string `json:"time"`
	Unixtime int64  `json:"utime"`
}

//MessageDataList
type MessageDataList []MessageData

func (a MessageDataList) Len() int           { return len(a) }
func (a MessageDataList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MessageDataList) Less(i, j int) bool { return a[i].Unixtime <= a[j].Unixtime }

//GetMessages
func (m *EventFilePlugin) GetMessages(eventID, level string, length int) (interface{}, error) {
	var message MessageDataList
	apath := path.Join(m.HomePath, "eventlog", eventID+".log")
	if ok, err := util.FileExists(apath); !ok {
		if err != nil {
			logrus.Errorf("check file exist error %s", err.Error())
		}
		return message, nil
	}
	eventFile, err := os.Open(apath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer eventFile.Close()
	reader := bufio.NewReader(eventFile)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				logrus.Error("read event log file error:", err.Error())
			}
			break
		}
		if len(line) > 2 {
			flag := line[0]
			if CheckLevel(string(flag), level) {
				info := strings.SplitN(string(line), " ", 3)
				if len(info) == 3 {
					timeunix := info[1]
					unix, _ := strconv.ParseInt(timeunix, 10, 64)
					tm := time.Unix(unix, 0)
					md := MessageData{
						Message:  info[2],
						Unixtime: unix,
						Time:     tm.Format(time.RFC3339),
					}
					message = append(message, md)
					if len(message) > length && length != 0 {
						break
					}
				}
			}
		}
	}
	return message, nil
}

//CheckLevel
func CheckLevel(flag, level string) bool {
	switch flag {
	case "0":
		return true
	case "1":
		if level != "error" {
			return true
		}
	case "2":
		if level == "debug" {
			return true
		}
	}
	return false
}

//GetTimeUnix
func GetTimeUnix(timeStr string) int64 {
	var timeLayout string
	if strings.Contains(timeStr, ".") {
		timeLayout = "2006-01-02T15:04:05"
	} else {
		timeLayout = "2006-01-02T15:04:05+08:00"
	}
	loc, _ := time.LoadLocation("Local")
	utime, err := time.ParseInLocation(timeLayout, timeStr, loc)
	if err != nil {
		logrus.Errorf("Parse log time error %s", err.Error())
		return 0
	}
	return utime.Unix()
}

//GetLevelFlag
func GetLevelFlag(level string) []byte {
	switch level {
	case "error":
		return []byte("0 ")
	case "info":
		return []byte("1 ")
	case "debug":
		return []byte("2 ")
	default:
		return []byte("0 ")
	}
}

//Close
func (m *EventFilePlugin) Close() error {
	return nil
}
