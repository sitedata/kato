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

package db

import (
	"testing"
	"time"
)

func TestEventFileSaveMessage(t *testing.T) {
	eventFilePlugin := EventFilePlugin{
		HomePath: "./test",
	}
	if err := eventFilePlugin.SaveMessage([]*EventLogMessage{
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "info",
			Message: "1ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk",
			Time:    time.Now().Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "debug",
			Message: "2ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk",
			Time:    time.Now().Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "error",
			Message: "3ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk",
			Time:    time.Now().Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "debug",
			Message: "4ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk",
			Time:    time.Now().Add(time.Hour).Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "info",
			Message: "5ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk",
			Time:    time.Now().Format(time.RFC3339),
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := eventFilePlugin.SaveMessage([]*EventLogMessage{
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "info",
			Message: "1ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk2",
			Time:    time.Now().Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "debug",
			Message: "2ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk2",
			Time:    time.Now().Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "error",
			Message: "3ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk2",
			Time:    time.Now().Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "debug",
			Message: "4ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk2",
			Time:    time.Now().Add(time.Hour).Format(time.RFC3339),
		},
		&EventLogMessage{
			EventID: "eventidsadasd",
			Level:   "info",
			Message: "5ajdsnadskfasndjn afnasdfnln asdfjnajksndfjk2",
			Time:    time.Now().Format(time.RFC3339),
		},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGetMessage(t *testing.T) {
	eventFilePlugin := EventFilePlugin{
		HomePath: "/tmp",
	}
	list, err := eventFilePlugin.GetMessages("eventidsadasd", "debug", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
