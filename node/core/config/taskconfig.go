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

package config

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/core/store"
	"github.com/sirupsen/logrus"
)

//GroupContext
type GroupContext struct {
	ctx     context.Context
	groupID string
}

//NewGroupContext
func NewGroupContext(groupID string) *GroupContext {
	return &GroupContext{
		ctx:     context.Background(),
		groupID: groupID,
	}
}

//Add
func (g *GroupContext) Add(k, v interface{}) {
	g.ctx = context.WithValue(g.ctx, k, v)
	store.DefalutClient.Put(fmt.Sprintf("%s/group/%s/%s", option.Config.ConfigStoragePath, g.groupID, k), v.(string))
}

//Get
func (g *GroupContext) Get(k interface{}) interface{} {
	if v := g.ctx.Value(k); v != nil {
		return v
	}
	res, _ := store.DefalutClient.Get(fmt.Sprintf("%s/group/%s/%s", option.Config.ConfigStoragePath, g.groupID, k))
	if res.Count > 0 {
		return string(res.Kvs[0].Value)
	}
	return ""
}

//GetString
func (g *GroupContext) GetString(k interface{}) string {
	if v := g.ctx.Value(k); v != nil {
		return v.(string)
	}
	res, _ := store.DefalutClient.Get(fmt.Sprintf("%s/group/%s/%s", option.Config.ConfigStoragePath, g.groupID, k))
	if res.Count > 0 {
		return string(res.Kvs[0].Value)
	}
	return ""
}

var reg = regexp.MustCompile(`(?U)\$\{.*\}`)

//GetConfigKey
func GetConfigKey(rk string) string {
	if len(rk) < 4 {
		return ""
	}
	left := strings.Index(rk, "{")
	right := strings.Index(rk, "}")
	return rk[left+1 : right]
}

//ResettingArray
func ResettingArray(groupCtx *GroupContext, source []string) ([]string, error) {
	sourcecopy := make([]string, len(source))
	// Use copy
	for i, s := range source {
		sourcecopy[i] = s
	}
	for i, s := range sourcecopy {
		resultKey := reg.FindAllString(s, -1)
		for _, rk := range resultKey {
			key := strings.ToUpper(GetConfigKey(rk))
			// if len(key) < 1 {
			// 	return nil, fmt.Errorf("%s Parameter configuration error.please make sure `${XXX}`", s)
			// }
			value := GetConfig(groupCtx, key)
			sourcecopy[i] = strings.Replace(s, rk, value, -1)
		}
	}
	return sourcecopy, nil
}

//GetConfig
func GetConfig(groupCtx *GroupContext, key string) string {
	if groupCtx != nil {
		value := groupCtx.Get(key)
		if value != nil {
			switch value.(type) {
			case string:
				if value.(string) != "" {
					return value.(string)
				}
			case int:
				if value.(int) != 0 {
					return strconv.Itoa(value.(int))
				}
			case []string:
				if value.([]string) != nil {
					result := strings.Join(value.([]string), ",")
					if strings.HasSuffix(result, ",") {
						return result
					}
					return result + ","
				}
			case []interface{}:
				if value.([]interface{}) != nil && len(value.([]interface{})) > 0 {
					result := ""
					for _, v := range value.([]interface{}) {
						switch v.(type) {
						case string:
							result += v.(string) + ","
						case int:
							result += strconv.Itoa(v.(int)) + ","
						}
					}
					return result
				}
			}
		}
	}
	if dataCenterConfig == nil {
		return ""
	}
	cn := dataCenterConfig.GetConfig(key)
	if cn != nil && cn.Value != nil {
		if cn.ValueType == "string" || cn.ValueType == "" {
			return cn.Value.(string)
		}
		if cn.ValueType == "array" {
			switch cn.Value.(type) {
			case []string:
				return strings.Join(cn.Value.([]string), ",")
			case []interface{}:
				vas := cn.Value.([]interface{})
				result := ""
				for _, va := range vas {
					switch va.(type) {
					case string:
						result += va.(string) + ","
					case int:
						result += strconv.Itoa(va.(int)) + ","
					}
				}
				return result
			}
		}
		if cn.ValueType == "int" {
			return strconv.Itoa(cn.Value.(int))
		}
	}
	logrus.Warnf("can not find config for key %s", key)
	return ""
}

//ResettingString
func ResettingString(groupCtx *GroupContext, source string) (string, error) {
	resultKey := reg.FindAllString(source, -1)
	for _, rk := range resultKey {
		key := strings.ToUpper(GetConfigKey(rk))
		// if len(key) < 1 {
		// 	return nil, fmt.Errorf("%s Parameter configuration error.please make sure `${XXX}`", s)
		// }
		value := GetConfig(groupCtx, key)
		source = strings.Replace(source, rk, value, -1)
	}
	return source, nil
}

//ResettingMap
func ResettingMap(groupCtx *GroupContext, source map[string]string) (map[string]string, error) {
	sourcecopy := make(map[string]string, len(source))
	for k, v := range source {
		sourcecopy[k] = v
	}
	for k, s := range sourcecopy {
		resultKey := reg.FindAllString(s, -1)
		for _, rk := range resultKey {
			key := strings.ToUpper(GetConfigKey(rk))
			// if len(key) < 1 {
			// 	return nil, fmt.Errorf("%s Parameter configuration error.please make sure `${XXX}`", s)
			// }
			value := GetConfig(groupCtx, key)
			sourcecopy[k] = strings.Replace(s, rk, value, -1)
		}
	}
	return sourcecopy, nil
}
