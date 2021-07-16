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

package controller

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/gridworkz/kato/mq/api/mq"

	"github.com/sirupsen/logrus"

	discovermodel "github.com/gridworkz/kato/worker/discover/model"

	restful "github.com/emicklei/go-restful"
)

//Register
func Register(container *restful.Container, mq mq.ActionMQ) {
	MQSource{mq}.Register(container)
}

//MQSource
type MQSource struct {
	mq mq.ActionMQ
}

//Register
func (u MQSource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/mq").
		Doc("message queue interface").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.POST("/{topic}").To(u.enqueue).
		// docs
		Doc("send a task message to the topic queue").
		Operation("enqueue").
		Param(ws.PathParameter("topic", "queue topic name").DataType("string")).
		Reads(discovermodel.Task{}).
		Returns(201, "Message sent successfully", ResponseType{}).
		ReturnsError(400, "The message format is wrong", ResponseType{}))

	ws.Route(ws.GET("/{topic}").To(u.dequeue).
		// docs
		Doc("get a task message from the topic queue").
		Operation("dequeue").
		Param(ws.PathParameter("topic", "queue topic name").DataType("string")).
		Writes(ResponseType{
			Body: ResponseBody{
				Bean: discovermodel.Task{},
			},
		})) // from the request
	ws.Route(ws.GET("/topics").To(u.getAllTopics).
		// docs
		Doc("get all support topic").
		Operation("getAllTopics").
		Writes(ResponseType{
			Body: ResponseBody{
				Bean: discovermodel.Task{},
			},
		})) // from the request
	container.Add(ws)
}

//ResponseType
type ResponseType struct {
	Code      int          `json:"code"`
	Message   string       `json:"msg"`
	MessageCN string       `json:"msgcn"`
	Body      ResponseBody `json:"body,omitempty"`
}

//ResponseBody
type ResponseBody struct {
	Bean     interface{}   `json:"bean,omitempty"`
	List     []interface{} `json:"list,omitempty"`
	PageNum  int           `json:"pageNumber,omitempty"`
	PageSize int           `json:"pageSize,omitempty"`
	Total    int           `json:"total,omitempty"`
}

//NewResponseType
func NewResponseType(code int, message string, messageCN string, bean interface{}, list []interface{}) ResponseType {
	return ResponseType{
		Code:      code,
		Message:   message,
		MessageCN: messageCN,
		Body: ResponseBody{
			Bean: bean,
			List: list,
		},
	}
}

//NewPostSuccessResponse
func NewPostSuccessResponse(bean interface{}, list []interface{}, response *restful.Response) {
	response.WriteHeaderAndJson(201, NewResponseType(201, "", "", bean, list), restful.MIME_JSON)
	return
}

//NewSuccessResponse
func NewSuccessResponse(bean interface{}, list []interface{}, response *restful.Response) {
	response.WriteHeaderAndJson(200, NewResponseType(200, "", "", bean, list), restful.MIME_JSON)
	return
}

//NewSuccessMessageResponse
func NewSuccessMessageResponse(bean interface{}, list []interface{}, message, messageCN string, response *restful.Response) {
	response.WriteHeaderAndJson(200, NewResponseType(200, message, messageCN, bean, list), restful.MIME_JSON)
	return
}

//NewFaliResponse
func NewFaliResponse(code int, message string, messageCN string, response *restful.Response) {
	response.WriteHeaderAndJson(code, NewResponseType(code, message, messageCN, nil, nil), restful.MIME_JSON)
	return
}

func (u *MQSource) enqueue(request *restful.Request, response *restful.Response) {
	topic := request.PathParameter("topic")
	if topic == "" || !u.mq.TopicIsExist(topic) {
		NewFaliResponse(400, "topic can not be empty or topic is not define", "The subject cannot be empty or the current subject is not registered", response)
		return
	}
	body, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		NewFaliResponse(500, "request body error."+err.Error(), "Error reading data", response)
		return
	}
	request.Request.Body.Close()
	_, err = discovermodel.NewTask(body)
	if err != nil {
		NewFaliResponse(400, "request body error."+err.Error(), "Error in reading data, invalid data", response)
		return
	}

	ctx, cancel := context.WithCancel(request.Request.Context())
	defer cancel()
	err = u.mq.Enqueue(ctx, topic, string(body))
	if err != nil {
		NewFaliResponse(500, "enqueue error."+err.Error(), "Message queue error", response)
		return
	}
	logrus.Debugf("Add a task to queue :%s", strings.TrimSpace(string(body)))
	NewPostSuccessResponse(nil, nil, response)
}

func (u *MQSource) dequeue(request *restful.Request, response *restful.Response) {
	topic := request.PathParameter("topic")
	if topic == "" || !u.mq.TopicIsExist(topic) {
		NewFaliResponse(400, "topic can not be empty or topic is not define", "The subject cannot be empty or the current subject is not registered", response)
		return
	}
	ctx, cancel := context.WithCancel(request.Request.Context())
	defer cancel()
	value, err := u.mq.Dequeue(ctx, topic)
	if err != nil {
		NewFaliResponse(500, "dequeue error."+err.Error(), "Message out of queue error", response)
		return
	}
	task, err := discovermodel.NewTask([]byte(value))
	if err != nil {
		NewFaliResponse(500, "dequeue error."+err.Error(), "The queue read message format is illegal", response)
		return
	}
	NewSuccessResponse(task, nil, response)
	logrus.Debugf("Consume a task from queue :%s", strings.TrimSpace(value))
}

func (u *MQSource) getAllTopics(request *restful.Request, response *restful.Response) {
	topics := u.mq.GetAllTopics()
	var list []interface{}
	for _, t := range topics {
		list = append(list, t)
	}
	NewSuccessResponse(nil, list, response)
}
