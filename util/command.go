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

package util

import (
	"bytes"
	"io"
	"os/exec"
)

//PipeCommand
type PipeCommand struct {
	stack                    []*exec.Cmd
	finalStdout, finalStderr io.Reader
	pipestack                []*io.PipeWriter
}

//NewPipeCommand
func NewPipeCommand(stack ...*exec.Cmd) (*PipeCommand, error) {
	var errorbuffer bytes.Buffer
	pipestack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdinpipe, stdoutpipe := io.Pipe()
		stack[i].Stdout = stdoutpipe
		stack[i].Stderr = &errorbuffer
		stack[i+1].Stdin = stdinpipe
		pipestack[i] = stdoutpipe
	}
	finalStdout, err := stack[i].StdoutPipe()
	if err != nil {
		return nil, err
	}
	finalStderr, err := stack[i].StderrPipe()
	if err != nil {
		return nil, err
	}
	pipeCommand := &PipeCommand{
		stack:       stack,
		pipestack:   pipestack,
		finalStdout: finalStdout,
		finalStderr: finalStderr,
	}
	return pipeCommand, nil
}

//Run
func (p *PipeCommand) Run() error {
	return call(p.stack, p.pipestack)
}

//GetFinalStdout get final command stdout reader
func (p *PipeCommand) GetFinalStdout() io.Reader {
	return p.finalStdout
}

//GetFinalStderr get final command stderr reader
func (p *PipeCommand) GetFinalStderr() io.Reader {
	return p.finalStderr
}

func call(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = call(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}
