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

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/pquerna/ffjson/ffjson"

	"github.com/gridworkz/kato/cmd/node/option"
	"github.com/gridworkz/kato/node/api/model"
	"github.com/gridworkz/kato/node/core/store"
	"github.com/gridworkz/kato/util"

	client "github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

//DataCenterConfig
type DataCenterConfig struct {
	config  *model.GlobalConfig
	options *option.Conf
	ctx     context.Context
	cancel  context.CancelFunc
	//group config
	groupConfigs map[string]*GroupContext
}

var dataCenterConfig *DataCenterConfig

//GetDataCenterConfig
func GetDataCenterConfig() *DataCenterConfig {
	if dataCenterConfig == nil {
		return CreateDataCenterConfig()
	}
	return dataCenterConfig
}

//CreateDataCenterConfig
func CreateDataCenterConfig() *DataCenterConfig {
	ctx, cancel := context.WithCancel(context.Background())
	dataCenterConfig = &DataCenterConfig{
		options: option.Config,
		ctx:     ctx,
		cancel:  cancel,
		config: &model.GlobalConfig{
			Configs: make(map[string]*model.ConfigUnit),
		},
		groupConfigs: make(map[string]*GroupContext),
	}
	res, err := store.DefalutClient.Get(dataCenterConfig.options.ConfigStoragePath+"/global", client.WithPrefix())
	if err != nil {
		logrus.Error("load datacenter config error.", err.Error())
	}
	if res != nil {
		if len(res.Kvs) < 1 {
			dgc := &model.GlobalConfig{
				Configs: make(map[string]*model.ConfigUnit),
			}
			dataCenterConfig.config = dgc
		} else {
			for _, kv := range res.Kvs {
				dataCenterConfig.PutConfigKV(kv)
			}
		}
	}
	return dataCenterConfig
}

//Start, monitor configuration changes
func (d *DataCenterConfig) Start() {
	go util.Exec(d.ctx, func() error {
		ctx, cancel := context.WithCancel(d.ctx)
		defer cancel()
		logrus.Info("datacenter config listener start")
		ch := store.DefalutClient.WatchByCtx(ctx, d.options.ConfigStoragePath+"/global", client.WithPrefix())
		for event := range ch {
			for _, e := range event.Events {
				switch {
				case e.IsCreate(), e.IsModify():
					d.PutConfigKV(e.Kv)
				case e.Type == client.EventTypeDelete:
					d.DeleteConfig(util.GetIDFromKey(string(e.Kv.Key)))
				}
			}
		}
		return nil
	}, 1)
}

//Stop
func (d *DataCenterConfig) Stop() {
	d.cancel()
	logrus.Info("datacenter config listener stop")
}

//GetDataCenterConfig
func (d *DataCenterConfig) GetDataCenterConfig() (*model.GlobalConfig, error) {
	return d.config, nil
}

//PutDataCenterConfig
func (d *DataCenterConfig) PutDataCenterConfig(c *model.GlobalConfig) (err error) {
	if c == nil {
		return
	}
	for k, v := range c.Configs {
		d.config.Add(*v)
		_, err = store.DefalutClient.Put(d.options.ConfigStoragePath+"/global/"+k, v.String())
	}
	return err
}

//GetConfig
func (d *DataCenterConfig) GetConfig(name string) *model.ConfigUnit {
	return d.config.Get(name)
}

//CacheConfig
func (d *DataCenterConfig) CacheConfig(c *model.ConfigUnit) error {
	if c.Name == "" {
		return fmt.Errorf("config name can not be empty")
	}
	logrus.Debugf("add config %v", c)
	//Convert the value type from []interface{} to []string
	if c.ValueType == "array" {
		switch c.Value.(type) {
		case []interface{}:
			var data []string
			for _, v := range c.Value.([]interface{}) {
				data = append(data, v.(string))
			}
			c.Value = data
		}
		oldC := d.config.Get(c.Name)
		if oldC != nil {

			switch oldC.Value.(type) {
			case string:
				value := append(c.Value.([]string), oldC.Value.(string))
				util.Deweight(&value)
				c.Value = value
			case []string:
				value := append(c.Value.([]string), oldC.Value.([]string)...)
				util.Deweight(&value)
				c.Value = value
			default:
			}
		}
	}
	d.config.Add(*c)
	return nil
}

//PutConfig - Add or update configuration
func (d *DataCenterConfig) PutConfig(c *model.ConfigUnit) error {
	if c.Name == "" {
		return fmt.Errorf("config name can not be empty")
	}
	logrus.Debugf("add config %v", c)
	//Convert the value type from []interface{} to []string
	if c.ValueType == "array" {
		switch c.Value.(type) {
		case []interface{}:
			var data []string
			for _, v := range c.Value.([]interface{}) {
				data = append(data, v.(string))
			}
			c.Value = data
		}
		oldC := d.config.Get(c.Name)
		if oldC != nil {

			switch oldC.Value.(type) {
			case string:
				value := append(c.Value.([]string), oldC.Value.(string))
				util.Deweight(&value)
				c.Value = value
			case []string:
				value := append(c.Value.([]string), oldC.Value.([]string)...)
				util.Deweight(&value)
				c.Value = value
			default:
			}
		}
	}
	d.config.Add(*c)
	//Persistence
	_, err := store.DefalutClient.Put(d.options.ConfigStoragePath+"/global/"+c.Name, c.String())
	if err != nil {
		logrus.Error("put datacenter config to etcd error.", err.Error())
		return err
	}
	return nil
}

//PutConfigKV
func (d *DataCenterConfig) PutConfigKV(kv *mvccpb.KeyValue) {
	var cn model.ConfigUnit
	if err := ffjson.Unmarshal(kv.Value, &cn); err == nil {
		d.CacheConfig(&cn)
	} else {
		logrus.Errorf("parse config error,%s", err.Error())
	}
}

//DeleteConfig
func (d *DataCenterConfig) DeleteConfig(name string) {
	d.config.Delete(name)
}

//GetGroupConfig
func (d *DataCenterConfig) GetGroupConfig(groupID string) *GroupContext {
	if c, ok := d.groupConfigs[groupID]; ok {
		return c
	}
	c := NewGroupContext(groupID)
	d.groupConfigs[groupID] = c
	return c
}
