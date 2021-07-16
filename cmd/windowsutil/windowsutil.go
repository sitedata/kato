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

package main

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gridworkz/kato/util/windows"

	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/cmd/windowsutil/option"
	"github.com/spf13/pflag"
)

func main() {
	conf := option.Config{}
	conf.AddFlags(pflag.CommandLine)
	pflag.Parse()
	if !conf.Check() {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	shell := strings.Split(conf.RunShell, "&nbsp;")
	logrus.Infof("run shell: %s", shell)
	cmd := exec.CommandContext(ctx, shell[0], shell[1:]...)
	startFunc := func() error {
		cmd.Stdin = os.Stdin
		reader, err := cmd.StdoutPipe()
		if err != nil {
			logrus.Errorf("open command stdout error %s", err.Error())
		}
		errReader, err := cmd.StderrPipe()
		if err != nil {
			logrus.Errorf("open command stderr error %s", err.Error())
		}
		go readBuffer(reader, logrus.Info)
		go readBuffer(errReader, logrus.Error)
		go func() {
			logrus.Info("start run progress")
			err := cmd.Start()
			if err != nil {
				logrus.Errorf("start cmd failure %s", err.Error())
				cancel()
			}
		}()
		var s os.Signal = syscall.SIGTERM
		defer func() {
			if cmd.Process != nil {
				if err := cmd.Process.Signal(s); err != nil {
					logrus.Errorf("send SIGTERM signal to progress failure %s", err.Error())
				}
				time.Sleep(time.Second * 2)
			}
		}()
		//step finally: listen Signal
		term := make(chan os.Signal)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		select {
		case ls := <-term:
			s = ls
			logrus.Warn("Received SIGTERM, exiting gracefully...")
		case <-ctx.Done():
		}
		logrus.Info("See you next time!")
		return nil
	}
	stopFunc := func() error {
		cancel()
		return nil
	}
	if conf.RunAsService {
		if err := windows.RunAsService(conf.ServiceName, startFunc, stopFunc, conf.Debug); err != nil {
			logrus.Fatalf("run command failure %s", err.Error())
		}
	} else {
		startFunc()
	}
}

func readBuffer(reader io.ReadCloser, print func(args ...interface{})) {
	defer reader.Close()
	bufreader := bufio.NewReader(reader)
	for {
		line, _, err := bufreader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return
			}
			logrus.Errorf("read log buffer failure %s", err.Error())
			return
		}
		print(string(line))
	}
}
