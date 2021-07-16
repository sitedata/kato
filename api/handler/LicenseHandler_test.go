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
	"fmt"
	"github.com/gridworkz/kato/api/db"
	"github.com/gridworkz/kato/cmd/api/option"
	"testing"
)

func TestLicenseInfo(t *testing.T) {
	conf := option.Config{
		DBType:           "mysql",
		DBConnectionInfo: "admin:admin@tcp(127.0.0.1:3306)/region",
	}
	//create db manager
	if err := db.CreateDBManager(conf); err != nil {
		fmt.Printf("create db manager error, %v", err)

	}
	//create a license verification manager
	if err := CreateLicensesInfoManager(); err != nil {
		fmt.Printf("create license check manager error, %v", err)
	}
	lists, err := GetLicensesInfosHandler().ShowInfos()
	if err != nil {
		fmt.Printf("get list error, %v", err)
	}
	for _, v := range lists {
		fmt.Printf("license value is %v", v)
	}
}
