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

package region

// //DefineCloudAuth DefineCloudAuth
// func (t *tenant) DefineCloudAuth(gt *api_model.GetUserToken) DefineCloudAuthInterface {
// 	return &DefineCloudAuth{
// 		GT: gt,
// 	}
// }

// //DefineCloudAuth DefineCloudAuth
// type DefineCloudAuth struct {
// 	GT *api_model.GetUserToken
// }

// //DefineCloudAuthInterface DefineCloudAuthInterface
// type DefineCloudAuthInterface interface {
// 	GetToken() ([]byte, error)
// 	PostToken() ([]byte, error)
// 	PutToken() error
// }

// //GetToken GetToken
// func (d *DefineCloudAuth) GetToken() ([]byte, error) {
// 	resp, code, err := request(
// 		fmt.Sprintf("/cloud/auth/%s", d.GT.Body.EID),
// 		"GET",
// 		nil,
// 	)
// 	if err != nil {
// 		return nil, util.CreateAPIHandleError(code, err)
// 	}
// 	if code > 400 {
// 		if code == 404 {
// 			return nil, util.CreateAPIHandleError(code, fmt.Errorf("eid %s is not exist", d.GT.Body.EID))
// 		}
// 		return nil, util.CreateAPIHandleError(code, fmt.Errorf("get eid infos %s failed", d.GT.Body.EID))
// 	}
// 	//valJ, err := simplejson.NewJson(resp)
// 	return resp, nil
// }

// //PostToken PostToken
// func (d *DefineCloudAuth) PostToken() ([]byte, error) {
// 	data, err := ffjson.Marshal(d.GT.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, code, err := request(
// 		"/cloud/auth",
// 		"POST",
// 		data,
// 	)
// 	if err != nil {
// 		logrus.Errorf("create auth token error, %v", err)
// 		return nil, util.CreateAPIHandleError(code, err)
// 	}
// 	if code > 400 {
// 		logrus.Errorf("create auth token error")
// 		return nil, util.CreateAPIHandleError(code, fmt.Errorf("cretae auth token failed"))
// 	}
// 	return resp, nil
// }

// //PutToken PutToken
// func (d *DefineCloudAuth) PutToken() error {
// 	data, err := ffjson.Marshal(d.GT.Body)
// 	if err != nil {
// 		return err
// 	}
// 	_, code, err := request(
// 		fmt.Sprintf("/cloud/auth/%s", d.GT.Body.EID),
// 		"PUT",
// 		data,
// 	)
// 	if err != nil {
// 		logrus.Errorf("create auth token error, %v", err)
// 		return util.CreateAPIHandleError(code, err)
// 	}
// 	if code > 400 {
// 		logrus.Errorf("create auth token error")
// 		return util.CreateAPIHandleError(code, fmt.Errorf("cretae auth token failed"))
// 	}
// 	return nil
// }
