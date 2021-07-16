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

package app

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/kr/pty"
)

func TestTTY(t *testing.T) {
	pty, tty, err := pty.Open()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		buf := make([]byte, 1024)
		pty.WriteString("hello word")
		var i int
		for {
			i++
			size, err := tty.Read(buf)
			if err != nil {
				log.Printf("Command exited for: %s", err.Error())
				return
			}
			pty.WriteString(string(buf[:size]) + strconv.Itoa(i))
			fmt.Println("tty write:", string(buf[:size]))
			time.Sleep(time.Second * 1)
		}
	}()
	// go func() {
	// 	buf := make([]byte, 1024)
	// 	for {
	// 		size, err := tty.Read(buf)
	// 		if err != nil {
	// 			log.Printf("Command exited for: %s", err.Error())
	// 			return
	// 		}
	// 		pty.Write(buf[:size])
	// 		fmt.Println("pty write:", buf[:size])
	// 	}
	// }()
	time.Sleep(time.Minute * 1)
}
