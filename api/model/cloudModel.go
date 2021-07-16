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

//GetUserToken
//swagger:parameters createToken
type GetUserToken struct {
	// in: body
	Body struct {
		// eid
		// in: body
		// required: true
		EID string `json:"eid" validate:"eid|required"`
		// Controllable range:all_power|node_manager|server_source
		// in: body
		// required: false
		Range string `json:"range" validate:"range"`
		// valid period
		// in: body
		// required: true
		ValidityPeriod int `json:"validity_period" validate:"validity_period|required"` //1549812345
		// data center identification
		// in: body
		// required: false
		RegionTag  string `json:"region_tag" validate:"region_tag"`
		BeforeTime int    `json:"before_time"`
	}
}

//GetTokenInfo
//swagger:parameters getTokenInfo
type GetTokenInfo struct {
	// in: path
	// required: true
	EID string `json:"eid" validate:"eid|required"`
}

//UpdateToken
//swagger:parameters updateToken
type UpdateToken struct {
	// in: path
	// required: true
	EID string `json:"eid" validate:"eid|required"`
	//in: body
	Body struct {
		// valid period
		// in: body
		// required: true
		ValidityPeriod int `json:"validity_period" validate:"validity_period|required"` //1549812345
	}
}

//TokenInfo
type TokenInfo struct {
	EID   string `json:"eid"`
	Token string `json:"token"`
	CA    string `json:"ca"`
	Key   string `json:"key"`
}

//APIManager
//swagger:parameters addAPIManager deleteAPIManager
type APIManager struct {
	//in: body
	Body struct {
		//api level
		//in: body
		//required: true
		ClassLevel string `json:"class_level" validate:"class_level|reqired"`
		//uri front
		//in: body
		//required: true
		Prefix string `json:"prefix" validate:"prefix|required"`
		//full uri
		//in: body
		//required: false
		URI string `json:"uri" validate:"uri"`
		//nickname
		//in: body
		//required: false
		Alias string `json:"alias" validate:"alias"`
		//additional information
		//in:body
		//required: false
		Remark string `json:"remark" validate:"remark"`
	}
}
