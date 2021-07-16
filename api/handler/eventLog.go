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

package handler

import (
	"bytes"
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/api/model"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	eventdb "github.com/gridworkz/kato/eventlog/db"
	"github.com/gridworkz/kato/util/constants"
)

//LogAction  log action struct
type LogAction struct {
	EtcdCli *clientv3.Client
	eventdb *eventdb.EventFilePlugin
}

//CreateLogManager get log manager
func CreateLogManager(cli *clientv3.Client) *LogAction {
	return &LogAction{
		EtcdCli: cli,
		eventdb: &eventdb.EventFilePlugin{
			HomePath: "/grdata/logs/",
		},
	}
}

// GetEvents get target logs
func (l *LogAction) GetEvents(target, targetID string, page, size int) ([]*dbmodel.ServiceEvent, int, error) {
	if target == "tenant" {
		return db.GetManager().ServiceEventDao().GetEventsByTenantID(targetID, (page-1)*size, size)
	}
	return db.GetManager().ServiceEventDao().GetEventsByTarget(target, targetID, (page-1)*size, size)
}

//GetLogList get log list
func (l *LogAction) GetLogList(serviceAlias string) ([]*model.HistoryLogFile, error) {
	logDIR := path.Join(constants.GrdataLogPath, serviceAlias)
	_, err := os.Stat(logDIR)
	if os.IsNotExist(err) {
		return nil, err
	}
	fileList, err := ioutil.ReadDir(logDIR)
	if err != nil {
		return nil, err
	}

	var logFiles []*model.HistoryLogFile
	for _, file := range fileList {
		logfile := &model.HistoryLogFile{
			Filename:     file.Name(),
			RelativePath: path.Join("logs", serviceAlias, file.Name()),
		}
		logFiles = append(logFiles, logfile)
	}
	return logFiles, nil
}

//GetLogFile GetLogFile
func (l *LogAction) GetLogFile(serviceAlias, fileName string) (string, string, error) {
	logPath := path.Join(constants.GrdataLogPath, serviceAlias)
	fullPath := path.Join(logPath, fileName)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return "", "", err
	}
	return logPath, fullPath, err
}

//GetLogInstance get log web socket instance
func (l *LogAction) GetLogInstance(serviceID string) (string, error) {
	value, err := l.EtcdCli.Get(context.Background(), fmt.Sprintf("/event/dockerloginstacne/%s", serviceID))
	if err != nil {
		return "", err
	}
	if len(value.Kvs) > 0 {
		return string(value.Kvs[0].Value), nil
	}

	return "", nil
}

//GetLevelLog get event log
func (l *LogAction) GetLevelLog(eventID string, level string) (*api_model.DataLog, error) {
	re, err := l.eventdb.GetMessages(eventID, level, 0)
	if err != nil {
		return nil, err
	}
	if re != nil {
		messageList, ok := re.(eventdb.MessageDataList)
		if ok {
			return &api_model.DataLog{
				Status: "success",
				Data:   messageList,
			}, nil
		}
	}
	return &api_model.DataLog{
		Status: "success",
		Data:   nil,
	}, nil
}

//Decompress zlib decoding
func decompress(zb []byte) ([]byte, error) {
	b := bytes.NewReader(zb)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(&out, r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func checkLevel(level, info string) bool {
	switch level {
	case "error":
		if info == "error" {
			return true
		}
		return false
	case "info":
		if info == "info" || info == "error" {
			return true
		}
		return false
	case "debug":
		if info == "info" || info == "error" || info == "debug" {
			return true
		}
		return false
	default:
		if info == "info" || info == "error" {
			return true
		}
		return false
	}
}

func uncompress(source []byte) (re []byte, err error) {
	r, err := zlib.NewReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	io.Copy(&buffer, r)
	r.Close()
	return buffer.Bytes(), nil
}

func bubSort(d []api_model.MessageData) []api_model.MessageData {
	for i := 0; i < len(d); i++ {
		for j := i + 1; j < len(d); j++ {
			if d[i].Unixtime > d[j].Unixtime {
				temp := d[i]
				d[i] = d[j]
				d[j] = temp
			}
		}
	}
	return d
}
