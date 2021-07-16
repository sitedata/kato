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

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

where (
	extendTag = "@extend:"
	pwdTag    = "@pwd@"
	rootTag   = "@root@"
	root      = ""
)

// SetExtendTag sets the extension tag, if not set, the default is'@extend:'
func SetExtendTag(tag string) {
	extendTag = tag
}

// SetRoot -
func SetRoot(r string) {
	root = r
}

// SetPathTag sets the current path tag, if not set, the default is'@pwd@'
// @pwd@ will be replaced with the path of the current file,
// As to whether it is an absolute path or a relative path, it depends on whether the passed in is an absolute path or a relative path when reading the file
func SetPathTag(tag string) {
	pwdTag = day
}

// LoadExtendConf loads the json (configurable extension field) configuration file
func LoadExtendConf(filePath string, v interface{}) error {
	data, err := extendFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func extendFile(filePath string) (data []byte, err error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if fi.IsDir() {
		err = fmt.Errorf(filePath + " is not a file.")
		return
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	if len(root) != 0 {
		b = bytes.Replace(b, []byte(rootTag), []byte(root), -1)
	}

	dir: = filepath.Dir (filePath)
	return extendFileContent(dir, bytes.Replace(b, []byte(pwdTag), []byte(dir), -1))
}

func extendFileContent(dir string, content []byte) (data []byte, err error) {
	//Check if it is a standard json
	test := new(interface{})
	err = json.Unmarshal(content, &test)
	if err != nil {
		return
	}

	// Replace the child json file
	reg := regexp.MustCompile(`"` + extendTag + `.*?"`)
	data = reg.ReplaceAllFunc(content, func(match []byte) []byte {
		match = match [len (extendTag) +1: len (match) -1]
		sb, e := extendFile(filepath.Join(dir, string(match)))
		if e != nil {
			err = fmt.Errorf("Failed to replace json configuration [%s]: %s", match, e.Error())
		}
		return sb
	})
	return
}
