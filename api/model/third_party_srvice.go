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

package model

// AddEndpiontsReq is one of the Endpoints in the request to add the endpints.
type AddEndpiontsReq struct {
	Address string `json:"address" validate:"address|required"`
}

// UpdEndpiontsReq is one of the Endpoints in the request to update the endpints.
type UpdEndpiontsReq struct {
	EpID    string `json:"ep_id" validate:"required|len:32"`
	Address string `json:"address"`
}

// DelEndpiontsReq is one of the Endpoints in the request to update the endpints.
type DelEndpiontsReq struct {
	EpID string `json:"ep_id" validate:"required|len:32"`
}

// ThirdEndpoint is one of the Endpoints list in the response to list, add,
// update or delete the endpints.
type ThirdEndpoint struct {
	EpID     string `json:"ep_id"`
	Address  string `json:"address"`
	Status   string `json:"status"`
	IsStatic bool   `json:"is_static"`
}

// ThirdEndpoints -
type ThirdEndpoints []*ThirdEndpoint

// Len is part of sort.Interface.
func (e ThirdEndpoints) Len() int {
	return len(e)
}

// Swap is part of sort.Interface.
func (e ThirdEndpoints) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (e ThirdEndpoints) Less(i, j int) bool {
	return e[i].Address < e[j].Address
}

// ThridPartyServiceProbe is the json obejct in the request
// to update or fetch the ThridPartyServiceProbe.
type ThridPartyServiceProbe struct {
	Scheme       string `json:"scheme"`
	Path         string `json:"path"`
	Port         int    `json:"port"`
	TimeInterval int    `json:"time_interval"`
	MaxErrorNum  int    `json:"max_error_num"`
	Action       string `json:"action"`
}
