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
	"fmt"
	"os/exec"
	"testing"
	"time"

	api_db "github.com/gridworkz/kato/api/db"
	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/cmd/api/option"

	"github.com/sirupsen/logrus"
)

func TestEmessage(t *testing.T) {
	fmt.Printf("begin.\n")
	conf := option.Config{
		DBType:           "mysql",
		DBConnectionInfo: "admin:admin@tcp(127.0.0.1:3306)/region",
	}
	//Create db manager
	if err := api_db.CreateDBManager(conf); err != nil {
		fmt.Printf("create db manager error, %v", err)
	}
	getLevelLog("dd09a25eb9744afa9b3ad5f5541013e7", "info")
	fmt.Printf("end.\n")
}

func getLevelLog(eventID string, level string) (*api_model.DataLog, error) {
	//messages, err := db.GetManager().EventLogDao().GetEventLogMessages(eventID)
	//if err != nil {
	//	return nil, err
	//}
	////var d []api_model.MessageData
	//for _, v := range messages {
	//	log, err := uncompress(v.Message)
	//	if err != nil {
	//		return nil, err
	//	}
	//	logrus.Debugf("log is %v", log)
	//	fmt.Printf("log is %v", string(log))
	//
	//	var mLogs []msgStruct
	//	if err := ffjson.Unmarshal(log, &mLogs); err != nil {
	//		return nil, err
	//	}
	//	fmt.Printf("jlog %v", mLogs)
	//	break
	//}
	return nil, nil
}

type msgStruct struct {
	EventID string `json:"event_id"`
	Step    string `json:"step"`
	Message string `json:"message"`
	Level   string `json:"level"`
	Time    string `json:"time"`
}

func TestLines(t *testing.T) {
	filePath := "/Users/pujielan/Downloads/log"
	logrus.Debugf("file path is %s", filePath)
	n := 1000
	f, err := exec.Command("tail", "-n", fmt.Sprintf("%d", n), filePath).Output()
	if err != nil {
		fmt.Printf("err if %v", err)
	}
	fmt.Printf("f is %v", string(f))
}

func TestTimes(t *testing.T) {
	//toBeCharge := "2015-01-01 00:00:00"     //the string to be converted into a timestamp. Note that the hours and minutes must also be written in seconds, because they follow the template. You don’t need to write it if you modify the template.
	toBeCharge := "2017-09-29T10:02:44+08:00" //the string to be converted into a timestamp. Note that the hours and minutes must also be written in seconds, because they follow the template. You don’t need to write it if you modify the template.
	timeLayout := "2006-01-02T15:04:05"       //template required for transformation
	loc, _ := time.LoadLocation("Local")      //important: Get the time zone
	//toBeCharge = strings.Split(toBeCharge, ".")[0]
	fmt.Println(toBeCharge)
	theTime, err := time.ParseInLocation(timeLayout, toBeCharge, loc) //use the template to convert to the time.time type in the corresponding time zone
	fmt.Println(err)
	sr := theTime.Unix() //converted to timestamp type is int64
	fmt.Println(theTime) //print out the time 2015-01-01 15:15:00 +0800 CST
	fmt.Println(sr)
}

func TestSort(t *testing.T) {
	arr := [...]int{3, 41, 24, 76, 11, 45, 3, 3, 64, 21, 69, 19, 36}
	fmt.Println(arr)
	num := len(arr)

	//Circular sort
	for i := 0; i < num; i++ {
		for j := i + 1; j < num; j++ {
			if arr[i] > arr[j] {
				temp := arr[i]
				arr[i] = arr[j]
				arr[j] = temp
			}
		}
	}
	fmt.Println(arr)
}

func quickSort(array []int, left int, right int) {
	if left < right {
		key := array[left]
		low := left
		high := right
		for low < high {
			for low < high && array[high] > key {
				high--
			}
			array[low] = array[high]
			for low < high && array[low] < key {
				low++
			}
			array[high] = array[low]
		}
		array[low] = key
		quickSort(array, left, low-1)
		quickSort(array, low+1, right)
	}
}

func qsort(array []int, low, high int) {
	if low < high {
		m := partition(array, low, high)
		// fmt.Println(m)
		qsort(array, low, m-1)
		qsort(array, m+1, high)
	}
}

func partition(array []int, low, high int) int {
	key := array[low]
	tmpLow := low
	tmpHigh := high
	for {
		//Find an element less than or equal to key. The position of the element must be between tmpLow and high, because array[tmpLow] and the left element are less than or equal to key, and will not cross the boundary
		for array[tmpHigh] > key {
			tmpHigh--
		}
		//Find an element greater than key. The position of the element must be between low and tmpHigh+1. Because array[tmpHigh+1] must be greater than key
		for array[tmpLow] <= key && tmpLow < tmpHigh {
			tmpLow++
		}

		if tmpLow >= tmpHigh {
			break
		}
		// swap(array[tmpLow], array[tmpHigh])
		array[tmpLow], array[tmpHigh] = array[tmpHigh], array[tmpLow]
		fmt.Println(array)
	}
	array[tmpLow], array[low] = array[low], array[tmpLow]
	return tmpLow
}

func TestFastSort(t *testing.T) {
	var sortArray = []int{3, 41, 24, 76, 11, 45, 3, 3, 64, 21, 69, 19, 36}
	fmt.Println(sortArray)
	qsort(sortArray, 0, len(sortArray)-1)
	fmt.Println(sortArray)
}
