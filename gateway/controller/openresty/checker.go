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

package openresty

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Check returns if the nginx healthz endpoint is returning ok (status code 200)
func (o *OrService) Check() error {
	url := fmt.Sprintf("http://127.0.0.1:%v/%v", o.ocfg.ListenPorts.Status, o.ocfg.HealthPath)
	timeout := o.ocfg.HealthCheckTimeout
	statusCode, err := simpleGet(url, timeout)
	if err != nil {
		logrus.Errorf("error checking %s healthz: %v", url, err)
		return err
	}
	if statusCode != 200 {
		return fmt.Errorf("ingress controller is not healthy")
	}

	return nil
}

func simpleGet(url string, timeout time.Duration) (int, error) {
	client := &http.Client{
		Timeout:   timeout * time.Second,
		Transport: &http.Transport{DisableKeepAlives: true},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	res, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}
