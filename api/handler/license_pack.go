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

package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

//Info license information
type Info struct {
	Code       string   `json:"code"`
	Company    string   `json:"company"`
	Node       int64    `json:"node"`
	CPU        int64    `json:"cpu"`
	Memory     int64    `json:"memory"`
	Tenant     int64    `json:"tenant"`
	EndTime    string   `json:"end_time"`
	StartTime  string   `json:"start_time"`
	DataCenter int64    `json:"data_center"`
	ModuleList []string `json:"module_list"`
}

var key = []byte("qa123zxswe3532crfvtg123bnhymjuki")

//decryption algorithm
func decrypt(key []byte, encrypted string) ([]byte, error) {
	return []byte{}, nil
}

//ReadLicenseFromFile
func ReadLicenseFromFile(licenseFile string) (Info, error) {

	info := Info{}
	//step1 read license file
	_, err := os.Stat(licenseFile)
	if err != nil {
		return info, err
	}
	infoBody, err := ioutil.ReadFile(licenseFile)
	if err != nil {
		return info, errors.New("LICENSE file is not readable")
	}

	//step2 decryption info
	infoData, err := decrypt(key, string(infoBody))
	if err != nil {
		return info, errors.New("An error occurred during LICENSE decryption")
	}
	err = json.Unmarshal(infoData, &info)
	if err != nil {
		return info, errors.New("An error occurred while decoding the LICENSE file")
	}
	return info, nil
}

//BasePack base pack
func BasePack(text []byte) (string, error) {
	token := ""
	encodeStr := base64.StdEncoding.EncodeToString(text)
	begin := 0
	if len([]byte(encodeStr)) > 40 {
		begin = randInt(0, (len([]byte(encodeStr)) - 40))
	} else {
		return token, fmt.Errorf("error license")
	}
	token = string([]byte(encodeStr)[begin:(begin + 40)])
	return token, nil
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
