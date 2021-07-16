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
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/gridworkz/kato/discover/config"
	"github.com/sirupsen/logrus"
)

//TrimAndSort
func TrimAndSort(endpoints []*config.Endpoint) []string {
	arr := make([]string, 0, len(endpoints))
	for _, end := range endpoints {
		if strings.HasPrefix(end.URL, "https://") {
			url := strings.TrimLeft(end.URL, "https://")
			arr = append(arr, url)
			continue
		}
		url := strings.TrimLeft(end.URL, "http://")
		arr = append(arr, url)
	}
	sort.Strings(arr)
	return arr
}

//ArrCompare
func ArrCompare(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i, item := range arr1 {
		if item != arr2[i] {
			return false
		}
	}

	return true
}

//ListenStop
func ListenStop() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	signal.Ignore(syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	logrus.Warn("monitor manager received signal: ", sig.String())
	close(sigs)
}
