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

package license

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

//LicenseInfo license data
type LicenseInfo struct {
	Code      string    `json:"code"`
	Company   string    `json:"company"`
	Node      int64     `json:"node"`
	Memory    int64     `json:"memory"`
	EndTime   string    `json:"end_time"`
	StartTime string    `json:"start_time"`
	Features  []Feature `json:"features"`
}

func (l *LicenseInfo) HaveFeature(code string) bool {
	for _, f := range l.Features {
		if f.Code == strings.ToUpper(code) {
			return true
		}
	}
	return false
}

type Feature struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var licenseInfo *LicenseInfo

//ReadLicense -
func ReadLicense() *LicenseInfo {
	if licenseInfo != nil {
		return licenseInfo
	}
	licenseFile := os.Getenv("LICENSE_PATH")
	if licenseFile == "" {
		return nil
	}
	//step1 read license file
	_, err := os.Stat(licenseFile)
	if err != nil {
		logrus.Error("read LICENSE file failure：" + err.Error())
		return nil
	}
	infoBody, err := ioutil.ReadFile(licenseFile)
	if err != nil {
		logrus.Error("read LICENSE file failure：" + err.Error())
		return nil
	}

	//step2 decryption info
	key := os.Getenv("LICENSE_KEY")
	if key == "" {
		logrus.Error("not define license Key")
		return nil
	}
	infoData, err := Decrypt(getKey(key), string(infoBody))
	if err != nil {
		logrus.Error("decrypt LICENSE failure " + err.Error())
		return nil
	}
	info := LicenseInfo{}
	err = json.Unmarshal(infoData, &info)
	if err != nil {
		logrus.Error("decrypt LICENSE json failure " + err.Error())
		return nil
	}
	licenseInfo = &info
	return &info
}

func Decrypt(key []byte, encrypted string) ([]byte, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}
func getKey(source string) []byte {
	if len(source) > 32 {
		return []byte(source[:32])
	}
	return append(defaultKey[len(source):], []byte(source)...)
}

var defaultKey = []byte{113, 119, 101, 114, 116, 121, 117, 105, 111, 112, 97, 115, 100, 102, 103, 104, 106, 107, 108, 122, 120, 99, 118, 98, 110, 109, 49, 50, 51, 52, 53, 54}
