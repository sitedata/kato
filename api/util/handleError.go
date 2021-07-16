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

package util

import (
	"fmt"
	"net/http"
	"strings"

	httputil "github.com/gridworkz/kato/util/http"

	"github.com/jinzhu/gorm"

	"github.com/sirupsen/logrus"
)

//APIHandleError handle create err for api
type APIHandleError struct {
	Code int
	Err  error
}

//CreateAPIHandleError create APIHandleError
func CreateAPIHandleError(code int, err error) *APIHandleError {
	return &APIHandleError{
		Code: code,
		Err:  err,
	}
}

//CreateAPIHandleErrorf create handle error
func CreateAPIHandleErrorf(code int, format string, args ...interface{}) *APIHandleError {
	return &APIHandleError{
		Code: code,
		Err:  fmt.Errorf(format, args...),
	}
}

//CreateAPIHandleErrorFromDBError from db error create APIHandleError
func CreateAPIHandleErrorFromDBError(msg string, err error) *APIHandleError {
	if err.Error() == gorm.ErrRecordNotFound.Error() {
		return &APIHandleError{
			Code: 404,
			Err:  fmt.Errorf("%s:%s", msg, err.Error()),
		}
	}
	if strings.HasSuffix(strings.TrimRight(err.Error(), " "), "is exist") {
		return &APIHandleError{
			Code: 400,
			Err:  fmt.Errorf("%s:%s", msg, err.Error()),
		}
	}
	return &APIHandleError{
		Code: 500,
		Err:  fmt.Errorf("%s:%s", msg, err.Error()),
	}
}
func (a *APIHandleError) Error() string {
	return a.Err.Error()
}

func (a *APIHandleError) String() string {
	return fmt.Sprintf("(Code:%d) %s", a.Code, a.Err.Error())
}

//Handle
func (a *APIHandleError) Handle(r *http.Request, w http.ResponseWriter) {
	if a.Code >= 500 {
		logrus.Error(a.String())
		httputil.ReturnError(r, w, a.Code, a.Error())
		return
	}
	httputil.ReturnError(r, w, a.Code, a.Error())
	return
}
