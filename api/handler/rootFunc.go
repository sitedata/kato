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
	"strconv"
	"strings"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/cmd/api/option"

	"github.com/sirupsen/logrus"
)

//RootAction  root function action struct
type RootAction struct{}

//CreateRootFuncManager get root func manager
func CreateRootFuncManager(conf option.Config) *RootAction {
	return &RootAction{}
}

//VersionInfo VersionInfo
type VersionInfo struct {
	Version []*LangInfo `json:"version"`
}

//LangInfo LangInfo
type LangInfo struct {
	Lang  string `json:"lang"`
	Major []*MajorInfo
}

//MajorInfo MajorInfo
type MajorInfo struct {
	Major int `json:"major"`
	Minor []*MinorInfo
}

//MinorInfo MinorInfo
type MinorInfo struct {
	Minor int   `json:"minor"`
	Patch []int `json:"patch"`
}

//{"php":{"3":{"4":3, "5":2}, "4":5}}
//MinorVersion := make(map[string]int)
//MajorVersion := make(map[string]minorVersion)
//ListVersion := make(map[string](make(map[string](make(map[string]int)))))

//ResolvePHP php - application build
func (r *RootAction) ResolvePHP(cs *api_model.ComposerStruct) (string, error) {
	lang := cs.Body.Lang
	data := cs.Body.Data
	logrus.Debugf("Composer got default_runtime=%v, json body=%v", lang, data)
	jsonData := cs.Body.Data.JSON
	if cs.Body.Data.JSON.PlatForm.PHP == "" {
		jsonData = cs.Body.Data.Lock
	}
	listVersions, err := createListVersion(cs)
	if err != nil {
		return "", err
	}
	isphp := jsonData.PlatForm.PHP
	if isphp != "" {
		var maxVersion string
		var errV error
		if strings.HasPrefix(isphp, "~") {
			si := strings.Split(isphp, "~")
			mm := strings.Split(si[1], ".")
			major, minor, _, err := transAtoi(mm)
			if err != nil {
				return "", err
			}
			maxVersion, errV = getMaxVersion(lang, &listVersions, major, minor)
			if errV != nil {
				return "", errV
			}
		} else if strings.HasPrefix(isphp, ">=") {
			si := strings.Split(isphp, ">=")
			mm := strings.Split(si[1], ".")
			major, minor, _, err := transAtoi(mm)
			if err != nil {
				return "", err
			}
			maxVersion, errV = getMaxVersion(lang, &listVersions, major, minor)
			if errV != nil {
				return "", errV
			}
		} else {
			mm := strings.Split(isphp, ".")
			major, minor, patch, err := transAtoi(mm)
			if err != nil {
				return "", err
			}
			maxVersion, errV = getMaxVersion(lang, &listVersions, major, minor, patch)
			if errV != nil {
				return "", errV
			}
		}
		return fmt.Sprintf("{%s|composer.json|%s|%s", lang, isphp, maxVersion), nil
	}
	maxVersion, errM := getMaxVersion(lang, &listVersions)
	if errM != nil {
		return "", errM
	}
	return fmt.Sprintf("{%s|default|*|%s}", lang, maxVersion), nil
}

func createListVersion(cs *api_model.ComposerStruct) (map[string]*VersionInfo, error) {
	//listVersions := make(map[string]*VersionInfo)
	/*
		listVersions := make(map[string]interface{})
		var vi VersionInfo
		for _, p := range cs.Body.Data.Packages {
			mm := strings.Split(p, "-")
			name := mm[0]
			version := mm[1]
			var li LangInfo
			mp := strings.Split(version, ".")
			major, minor, patch, errT := transAtoi(mp)
			if errT != nil {
				return nil, errT
			}

				if _, ok := listVersions[name]; !ok {
					var nn VersionInfo
					listVersions[name] = &nn
				}
				mp := strings.Split(version, ".")
				major, minor, patch, errT := transAtoi(mp)
				if errT != nil {
					return nil, errT
				}
				listVersions[name].Majon = major
				listVersions[name].Minor = minor
				listVersions[name].Patch = patch
		}*/
	return nil, nil
}

func transAtoi(mm []string) (int, int, int, error) {
	major, minor, patch := 0, 0, 0
	major, errMa := strconv.Atoi(mm[0])
	if errMa != nil {
		return 0, 0, 0, errMa
	}
	minor, errMi := strconv.Atoi(mm[1])
	if errMi != nil {
		return 0, 0, 0, errMi
	}
	if len(mm) == 3 {
		var err error
		patch, err = strconv.Atoi(mm[2])
		if err != nil {
			return 0, 0, 0, err
		}
	} else if len(mm) == 2 {
		patch = 0
	} else {
		return 0, 0, 0, fmt.Errorf("version length error")
	}
	return major, minor, patch, nil
}

func getMaxVersion(l string, lv *map[string]*VersionInfo, opts ...int) (string, error) {

	return "", nil
}
