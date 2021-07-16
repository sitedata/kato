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

package daemon

import (
	"github.com/gridworkz/kato/api/util"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

//StartRegionAPI start up
func StartRegionAPI(ch chan os.Signal) {
	logrus.Info("old region api begin start..")
	arg := []string{"region_api.wsgi", "-b=127.0.0.1:8887", "--max-requests=5000", "--reload", "--debug", "--workers=4", "--log-file", "-", "--access-logfile", "-", "--error-logfile", "-"}
	cmd := exec.Command("gunicorn", arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go func() {
		if err := cmd.Start(); err != nil {
			logrus.Error("start region old api error.", err.Error())
		}
		tick := time.NewTicker(time.Second * 5)
		select {
		case si := <-ch:
			cmd.Process.Signal(si)
			return
		case <-tick.C:
			monitor()
		}
	}()
	return
}
func monitor() {
	response, err := http.Get("http://127.0.0.1:8887/monitor")
	if err != nil {
		logrus.Error("monitor region old api error.", err.Error())
		return
	}
	defer util.CloseResponse(response)
	if response != nil && response.StatusCode/100 > 2 {
		logrus.Errorf("monitor region old api error. response code is %s", response.Status)
	}

}
