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
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
)

type filePlugin struct {
	homePath string
}

func (m *filePlugin) getStdFilePath(serviceID string) (string, error) {
	apath := path.Join(m.homePath, GetServiceAliasID(serviceID))
	_, err := os.Stat(apath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(apath, 0755)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return apath, nil
}

func (m *filePlugin) SaveMessage(events []*EventLogMessage) error {
	if len(events) == 0 {
		return nil
	}
	key := events[0].EventID
	var logfile *os.File
	filePathDir, err := m.getStdFilePath(key)
	if err != nil {
		return err
	}
	logFile, err := os.Stat(path.Join(filePathDir, "stdout.log"))
	if err != nil {
		if os.IsNotExist(err) {
			logfile, err = os.Create(path.Join(filePathDir, "stdout.log"))
			if err != nil {
				return err
			}
			defer logfile.Close()
		} else {
			return err
		}
	} else {
		if logFile.ModTime().Day() != time.Now().Day() {
			err := MvLogFile(fmt.Sprintf("%s/%d-%d-%d.log.gz", filePathDir, logFile.ModTime().Year(), logFile.ModTime().Month(), logFile.ModTime().Day()), path.Join(filePathDir, "stdout.log"))
			if err != nil {
				return err
			}
		}
	}
	if logfile == nil {
		logfile, err = os.OpenFile(path.Join(filePathDir, "stdout.log"), os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer logfile.Close()
	} else {
		defer logfile.Close()
	}
	var contant [][]byte
	for _, e := range events {
		contant = append(contant, e.Content)
	}
	body := bytes.Join(contant, []byte("\n"))
	body = append(body, []byte("\n")...)
	_, err = logfile.Write(body)
	return err
}
func (m *filePlugin) GetMessages(serviceID, level string, length int) (interface{}, error) {
	if length <= 0 {
		return nil, nil
	}
	filePathDir, err := m.getStdFilePath(serviceID)
	if err != nil {
		return nil, err
	}
	filePath := path.Join(filePathDir, "stdout.log")
	if ok, err := util.FileExists(filePath); !ok {
		if err != nil {
			logrus.Errorf("check file exist error %s", err.Error())
		}
		return nil, nil
	}
	f, err := exec.Command("tail", "-n", fmt.Sprintf("%d", length), filePath).Output()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(bytes.NewBuffer(f))
	var lines []string
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if len(line) == 0 {
			continue
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

func (m *filePlugin) Close() error {
	return nil
}

//GetServiceAliasID python:
//new_word = str(ord(string[10])) + string + str(ord(string[3])) + 'log' + str(ord(string[2]) / 7)
//new_id = hashlib.sha224(new_word).hexdigest()[0:16]
//
func GetServiceAliasID(ServiceID string) string {
	if len(ServiceID) > 11 {
		newWord := strconv.Itoa(int(ServiceID[10])) + ServiceID + strconv.Itoa(int(ServiceID[3])) + "log" + strconv.Itoa(int(ServiceID[2])/7)
		ha := sha256.New224()
		ha.Write([]byte(newWord))
		return fmt.Sprintf("%x", ha.Sum(nil))[0:16]
	}
	return ServiceID
}

//MvLogFile - change file name, compress
func MvLogFile(newName string, filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	reader, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	// Write the contents of the compressed document to the file
	f, err := os.OpenFile(newName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	defer zw.Close()
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	new, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer new.Close()
	return nil
}
