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
	"strings"
	"time"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/cmd/api/option"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
)

//TokenIdenAction
type TokenIdenAction struct{}

//CreateTokenIdenManager token identification
func CreateTokenIdenManager(conf option.Config) (*TokenIdenAction, error) {
	return &TokenIdenAction{}, nil
}

//AddTokenIntoMap
func (t *TokenIdenAction) AddTokenIntoMap(rui *dbmodel.RegionUserInfo) {
	m := GetDefaultTokenMap()
	m[rui.Token] = rui
}

//DeleteTokenFromMap
func (t *TokenIdenAction) DeleteTokenFromMap(oldtoken string, rui *dbmodel.RegionUserInfo) {
	m := GetDefaultTokenMap()
	t.AddTokenIntoMap(rui)
	delete(m, oldtoken)
}

//GetAPIManager
func (t *TokenIdenAction) GetAPIManager() map[string][]*dbmodel.RegionAPIClass {
	return GetDefaultSourceURI()
}

//AddAPIManager
func (t *TokenIdenAction) AddAPIManager(am *api_model.APIManager) *util.APIHandleError {
	m := GetDefaultSourceURI()
	ra := &dbmodel.RegionAPIClass{
		ClassLevel: am.Body.ClassLevel,
		Prefix:     am.Body.Prefix,
	}
	if sourceList, ok := m[am.Body.ClassLevel]; ok {
		sourceList = append(sourceList, ra)
	} else {
		//support for new types
		newL := []*dbmodel.RegionAPIClass{ra}
		m[am.Body.ClassLevel] = newL
	}
	ra.URI = am.Body.URI
	ra.Alias = am.Body.Alias
	ra.Remark = am.Body.Remark
	if err := db.GetManager().RegionAPIClassDao().AddModel(ra); err != nil {
		return util.CreateAPIHandleErrorFromDBError("add api manager", err)
	}
	return nil
}

//DeleteAPIManager
func (t *TokenIdenAction) DeleteAPIManager(am *api_model.APIManager) *util.APIHandleError {
	m := GetDefaultSourceURI()
	if sourceList, ok := m[am.Body.ClassLevel]; ok {
		var newL []*dbmodel.RegionAPIClass
		for _, s := range sourceList {
			if s.Prefix == am.Body.Prefix {
				continue
			}
			newL = append(newL, s)
		}
		if len(newL) == 0 {
			//when the level group is empty, delete the resource group
			delete(m, am.Body.ClassLevel)
		} else {
			m[am.Body.ClassLevel] = newL
		}
	} else {
		return util.CreateAPIHandleError(400, fmt.Errorf("have no api class level about %v", am.Body.ClassLevel))
	}
	if err := db.GetManager().RegionAPIClassDao().DeletePrefixInClass(am.Body.ClassLevel, am.Body.Prefix); err != nil {
		return util.CreateAPIHandleErrorFromDBError("delete api prefix", err)
	}
	return nil
}

//CheckToken
func (t *TokenIdenAction) CheckToken(token, uri string) bool {
	m := GetDefaultTokenMap()
	//logrus.Debugf("default token map is %v", m)
	regionInfo, ok := m[token]
	if !ok {
		var err error
		regionInfo, err = db.GetManager().RegionUserInfoDao().GetTokenByTokenID(token)
		if err != nil {
			return false
		}
		SetTokenCache(regionInfo)
	}
	if regionInfo.ValidityPeriod < int(time.Now().Unix()) {
		return false
	}
	switch regionInfo.APIRange {
	case dbmodel.ALLPOWER:
		return true
	case dbmodel.SERVERSOURCE:
		sm := GetDefaultSourceURI()
		smL, ok := sm[dbmodel.SERVERSOURCE]
		if !ok {
			return false
		}
		rc := false
		for _, urinfo := range smL {
			if strings.HasPrefix(uri, urinfo.Prefix) {
				rc = true
			}
		}
		return rc
	case dbmodel.NODEMANAGER:
		sm := GetDefaultSourceURI()
		smL, ok := sm[dbmodel.NODEMANAGER]
		if !ok {
			return false
		}
		rc := false
		for _, urinfo := range smL {
			if strings.HasPrefix(uri, urinfo.Prefix) {
				rc = true
			}
		}
		return rc
	}
	return false
}

//InitTokenMap
func (t *TokenIdenAction) InitTokenMap() error {
	ruis, err := db.GetManager().RegionUserInfoDao().GetALLTokenInValidityPeriod()
	if err != nil {
		return err
	}
	m := GetDefaultTokenMap()
	for _, rui := range ruis {
		m[rui.Token] = rui
	}
	return nil
}
