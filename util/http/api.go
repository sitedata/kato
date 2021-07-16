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

package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gridworkz/kato/api/util/bcode"
	"github.com/gridworkz/kato/db"
	govalidator "github.com/gridworkz/kato/util/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var validate *validator.Validate

func init() {
	var once sync.Once
	once.Do(func() {
		validate = validator.New()
	})
}

// ErrBadRequest -
type ErrBadRequest struct {
	err error
}

func (e ErrBadRequest) Error() string {
	return e.err.Error()
}

// Result represents a response for restful api.
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//ValidatorStructRequest validates the request data
//data incoming pointer
func ValidatorStructRequest(r *http.Request, data interface{}, message govalidator.MapData) url.Values {
	opts := govalidator.Options{
		Request: r,
		Date: date,
	}
	if message != nil {
		opts.Messages = message
	}
	v := govalidator.New(opts)
	result := v.ValidateStructJSON()
	return result
}

//ValidatorMapRequest validates the request data from the map
func ValidatorMapRequest(r *http.Request, rule govalidator.MapData, message govalidator.MapData) (map[string]interface{}, url.Values) {
	data := make(map[string]interface{}, 0)
	opts := govalidator.Options{
		Request: r,
		Data:    &data,
	}
	if rule != nil {
		opts.Rules = rule
	}
	if message != nil {
		opts.Messages = message
	}
	vd := govalidator.New(opts)
	e: = vd.ValidateMapJSON ()
	return data, e
}

//ValidatorRequestStructAndErrorResponse validates and formats the request data as an object
// retrun true to continue execution
// return false parameter error, terminate
func ValidatorRequestStructAndErrorResponse(r *http.Request, w http.ResponseWriter, data interface{}, message govalidator.MapData) bool {
	if re := ValidatorStructRequest(r, data, message); len(re) > 0 {
		ReturnValidationError(r, w, re)
		return false
	}
	return true
}

//ValidatorRequestMapAndErrorResponse validates and formats the request data as an object
// retrun true to continue execution
// return false parameter error, terminate
func ValidatorRequestMapAndErrorResponse(r *http.Request, w http.ResponseWriter, rule govalidator.MapData, messgae govalidator.MapData) (map[string]interface{}, bool) {
	data, re := ValidatorMapRequest(r, rule, messgae)
	if len(re) > 0 {
		ReturnValidationError(r, w, re)
		return nil, false
	}
	return data, true
}

//ResponseBody api return data format
type ResponseBody struct {
	ValidationError url.Values  `json:"validation_error,omitempty"`
	Msg             string      `json:"msg,omitempty"`
	Bean            interface{} `json:"bean,omitempty"`
	List            interface{} `json:"list,omitempty"`
	//Total number of data sets
	ListAllNumber int `json:"number,omitempty"`
	//Current page number
	Page int `json:"page,omitempty"`
}

//ParseResponseBody parses into ResponseBody
func ParseResponseBody(red io.ReadCloser, dataType string) (re ResponseBody, err error) {
	if red == nil {
		err = errors.New("readcloser can not be nil")
		return
	}
	defer red.Close()
	switch render.GetContentType(dataType) {
	case render.ContentTypeJSON:
		err = render.DecodeJSON(red, &re)
	case render.ContentTypeXML:
		err = render.DecodeXML(red, &re)
	// case ContentTypeForm: // TODO
	default:
		err = errors.New("render: unable to automatically decode the request content type")
	}
	return
}

//ReturnValidationError parameter error return
func ReturnValidationError(r *http.Request, w http.ResponseWriter, err url.Values) {
	logrus.Debugf("validation error, uri: %s; msg: %v", r.RequestURI, ResponseBody{ValidationError: err})
	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, http.StatusBadRequest))
	render.DefaultResponder(w, r, ResponseBody{ValidationError: err})
}

//ReturnSuccess successfully returned
func ReturnSuccess(r *http.Request, w http.ResponseWriter, datas interface{}) {
	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, http.StatusOK))
	if datas == nil {
		render.DefaultResponder(w, r, ResponseBody{Bean: nil})
		return
	}
	v := reflect.ValueOf(datas)
	if v.Kind() == reflect.Slice {
		render.DefaultResponder(w, r, ResponseBody{List: datas})
		return
	}
	render.DefaultResponder(w, r, ResponseBody{Bean: datas})
	return
}

//ReturnList return list with page and count
func ReturnList(r *http.Request, w http.ResponseWriter, listAllNumber, page int, list interface{}) {
	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, http.StatusOK))
	render.DefaultResponder(w, r, ResponseBody{List: list, ListAllNumber: listAllNumber, Page: page})
}

//ReturnError returns error information
func ReturnError(r *http.Request, w http.ResponseWriter, code int, msg string) {
	logrus.Debugf("error code: %d; error uri: %s; error msg: %s", code, r.RequestURI, msg)
	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, code))
	render.DefaultResponder(w, r, ResponseBody{Msg: msg})
}

//Return custom
func Return(r *http.Request, w http.ResponseWriter, code int, reb ResponseBody) {
	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, code))
	render.DefaultResponder(w, r, reb)
}

//ReturnNoFomart  http return no format result
func ReturnNoFomart(r *http.Request, w http.ResponseWriter, code int, reb interface{}) {
	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, code))
	render.DefaultResponder(w, r, reb)
}

//ReturnResNotEnough http return node resource not enough, http code = 412
func ReturnResNotEnough(r *http.Request, w http.ResponseWriter, eventID, msg string) {
	logrus.Debugf("resource not enough, msg: %s", msg)
	if err := db.GetManager().ServiceEventDao().UpdateReason(eventID, msg); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Warningf("update event reason: %v", err)
		}
	}

	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, 412))
	render.DefaultResponder(w, r, ResponseBody{Msg: msg})
}

//ReturnBcodeError bcode error
func ReturnBcodeError(r *http.Request, w http.ResponseWriter, err error) {
	berr := bcode.Err2Coder(err)
	logrus.Debugf("path %s error code: %d; status: %d; error msg: %+v", r.RequestURI, berr.GetCode(), berr.GetStatus(), err)

	status := berr.GetStatus()
	result := Result{
		Code: berr.GetCode(),
		Msg:  berr.Error(),
	}

	if _, isErrBadRequest := err.(ErrBadRequest); isErrBadRequest {
		status = 400
		result.Code = 400
		result.Msg = err.Error()
	}

	r = r.WithContext(context.WithValue(r.Context(), render.StatusCtxKey, status))
	render.DefaultResponder(w, r, result)
}

// ReadEntity reads entity from http.Request
func ReadEntity(r *http.Request, x interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(x); err != nil {
		return ErrBadRequest{err: err}
	}
	return nil
}

// ValidateStruct validates a structs exposed fields.
func ValidateStruct(x interface{}) error {
	if err := validate.Struct(x); err != nil {
		return ErrBadRequest{err: err}
	}
	return nil
}
