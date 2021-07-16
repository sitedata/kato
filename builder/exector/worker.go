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

package exector

/*
Copyright 2017 The Gridworkz Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gridworkz/kato/util"
	"os"
	"os/exec"
	"time"
)

//Worker
type Worker struct {
	cmd  *exec.Cmd
	user string
}

func (w *Worker) run(timeout time.Duration) ([]byte, error) {
	stdout := &bytes.Buffer{}
	isTimeout, err := util.CmdRunWithTimeout(w.cmd, timeout)
	if err != nil {
		return nil, workerErr(err, stdout.Bytes())
	}
	if isTimeout {
		return nil, fmt.Errorf("exec worker timeout")
	}
	return stdout.Bytes(), nil
}

//NewWorker
func NewWorker(cmdpath, user string, envs []string, in []byte) *Worker {

	stdout := &bytes.Buffer{}
	c := &exec.Cmd{
		Env:    envs,
		Path:   "/usr/bin/python",
		Args:   []string{"python", cmdpath},
		Stdin:  bytes.NewBuffer(in),
		Stdout: stdout,
		Stderr: os.Stderr,
	}
	return &Worker{cmd: c, user: user}
}

//Error
type Error struct {
	Code    uint   `json:"code"`
	Msg     string `json:"msg"`
	Details string `json:"details,omitempty"`
}

func workerErr(err error, output []byte) error {
	if _, ok := err.(*exec.ExitError); ok {
		emsg := Error{}
		if perr := json.Unmarshal(output, &emsg); perr != nil {
			return fmt.Errorf("netplugin failed but error parsing its diagnostic message %q: %v", string(output), perr)
		}
		details := ""
		if emsg.Details != "" {
			details = fmt.Sprintf("; %v", emsg.Details)
		}
		return fmt.Errorf("%v%v", emsg.Msg, details)
	}
	return err
}
