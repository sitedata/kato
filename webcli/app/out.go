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
	"io"
	"os"

	"github.com/gridworkz/kato/webcli/term"
	"github.com/sirupsen/logrus"
)

//Out
type Out struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

//CreateOut
func CreateOut(tty *os.File) *Out {
	return &Out{
		Stdin:  tty,
		Stdout: tty,
		Stderr: tty,
	}
}

//SetTTY
func (o *Out) SetTTY() term.TTY {
	t := term.TTY{
		Out: o.Stdout,
		In:  o.Stdin,
	}
	if !t.IsTerminalIn() {
		logrus.Errorf("stdin is not tty")
		return t
	}
	// if we get to here, the user wants to attach stdin, wants a TTY, and o.In is a terminal, so we
	// can safely set t.Raw to true
	t.Raw = true
	return t
}
