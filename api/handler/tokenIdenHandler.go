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
	"os"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/cmd/api/option"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
)

//TokenMapHandler DefaultTokenMapHandler
type TokenMapHandler interface {
	AddTokenIntoMap(rui *dbmodel.RegionUserInfo)
	DeleteTokenFromMap(oldtoken string, rui *dbmodel.RegionUserInfo)
	CheckToken(token, uri string) bool
	GetAPIManager() map[string][]*dbmodel.RegionAPIClass
	AddAPIManager(am *api_model.APIManager) *util.APIHandleError
	DeleteAPIManager(am *api_model.APIManager) *util.APIHandleError
	InitTokenMap() error
}

var defaultTokenIdenHandler TokenMapHandler

//TokenMap
type TokenMap map[string]*dbmodel.RegionUserInfo

var defaultTokenMap map[string]*dbmodel.RegionUserInfo

var defaultSourceURI map[string][]*dbmodel.RegionAPIClass

//CreateTokenIdenHandler create token identification handler
func CreateTokenIdenHandler(conf option.Config) error {
	CreateDefaultTokenMap(conf)
	var err error
	if defaultTokenIdenHandler != nil {
		return nil
	}
	defaultTokenIdenHandler, err = CreateTokenIdenManager(conf)
	if err != nil {
		return err
	}
	return defaultTokenIdenHandler.InitTokenMap()
}

func createDefaultSourceURI() error {
	if defaultSourceURI != nil {
		return nil
	}
	var err error
	defaultSourceURI, err = resourceURI()
	if err != nil {
		return err
	}
	return nil
}

func resourceURI() (map[string][]*dbmodel.RegionAPIClass, error) {
	sourceMap := make(map[string][]*dbmodel.RegionAPIClass)
	nodeSource, err := db.GetManager().RegionAPIClassDao().GetPrefixesByClass(dbmodel.NODEMANAGER)
	if err != nil {
		return nil, err
	}
	sourceMap[dbmodel.NODEMANAGER] = nodeSource

	serverSource, err := db.GetManager().RegionAPIClassDao().GetPrefixesByClass(dbmodel.SERVERSOURCE)
	if err != nil {
		return nil, err
	}
	sourceMap[dbmodel.SERVERSOURCE] = serverSource
	return sourceMap, nil
}

//CreateDefaultTokenMap
func CreateDefaultTokenMap(conf option.Config) {
	createDefaultSourceURI()
	if defaultTokenMap != nil {
		return
	}
	consoleToken := "defaulttokentoken"
	if os.Getenv("TOKEN") != "" {
		consoleToken = os.Getenv("TOKEN")
	}
	rui := &dbmodel.RegionUserInfo{
		Token:          consoleToken,
		APIRange:       dbmodel.ALLPOWER,
		ValidityPeriod: 3257894000,
	}
	tokenMap := make(map[string]*dbmodel.RegionUserInfo)
	tokenMap[consoleToken] = rui
	defaultTokenMap = tokenMap
	return
}

//GetTokenIdenHandler
func GetTokenIdenHandler() TokenMapHandler {
	return defaultTokenIdenHandler
}

//GetDefaultTokenMap
func GetDefaultTokenMap() map[string]*dbmodel.RegionUserInfo {
	return defaultTokenMap
}

//SetTokenCache
func SetTokenCache(info *dbmodel.RegionUserInfo) {
	defaultTokenMap[info.Token] = info
}

//GetDefaultSourceURI
func GetDefaultSourceURI() map[string][]*dbmodel.RegionAPIClass {
	return defaultSourceURI
}
