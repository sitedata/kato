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

package db

import (
	"encoding/json"
	"time"

	tsdbClient "github.com/bluebreezecf/opentsdb-goclient/client"
	tsdbConfig "github.com/bluebreezecf/opentsdb-goclient/config"
	"github.com/gridworkz/kato/cmd/api/option"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/config"
	dbModel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/mq/api/grpc/pb"
	"github.com/gridworkz/kato/mq/client"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/gridworkz/kato/worker/discover/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ConDB struct
type ConDB struct {
	ConnectionInfo string
	DBType         string
}

//CreateDBManager get db manager
//TODO: need to retry when error happens, try 4 times
func CreateDBManager(conf option.Config) error {
	dbCfg := config.Config{
		MysqlConnectionInfo: conf.DBConnectionInfo,
		DBType:              conf.DBType,
		ShowSQL:             conf.ShowSQL,
	}
	if err := db.CreateManager(dbCfg); err != nil {
		logrus.Errorf("get db manager failed,%s", err.Error())
		return err
	}
	// api database initialization
	go dataInitialization()

	return nil
}

//CreateEventManager create event manager
func CreateEventManager(conf option.Config) error {
	var tryTime time.Duration
	var err error
	etcdClientArgs := &etcdutil.ClientArgs{
		Endpoints: conf.EtcdEndpoint,
		CaFile:    conf.EtcdCaFile,
		CertFile:  conf.EtcdCertFile,
		KeyFile:   conf.EtcdKeyFile,
	}
	for tryTime < 4 {
		tryTime++
		if err = event.NewManager(event.EventConfig{
			EventLogServers: conf.EventLogServers,
			DiscoverArgs:    etcdClientArgs,
		}); err != nil {
			logrus.Errorf("get event manager failed, try time is %v,%s", tryTime, err.Error())
			time.Sleep((5 + tryTime*10) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Errorf("get event manager failed. %v", err.Error())
		return err
	}
	logrus.Debugf("init event manager success")
	return nil
}

//MQManager mq manager
type MQManager struct {
	EtcdClientArgs *etcdutil.ClientArgs
	DefaultServer  string
}

//NewMQManager new mq manager
func (m *MQManager) NewMQManager() (client.MQClient, error) {
	client, err := client.NewMqClient(m.EtcdClientArgs, m.DefaultServer)
	if err != nil {
		logrus.Errorf("new mq manager error, %v", err)
		return client, err
	}
	return client, nil
}

//TaskStruct task struct
type TaskStruct struct {
	TaskType string
	TaskBody model.TaskBody
	User     string
}

//OpentsdbManager OpentsdbManager
type OpentsdbManager struct {
	Endpoint string
}

//NewOpentsdbManager NewOpentsdbManager
func (o *OpentsdbManager) NewOpentsdbManager() (tsdbClient.Client, error) {
	opentsdbCfg := tsdbConfig.OpenTSDBConfig{
		OpentsdbHost: o.Endpoint,
	}
	tc, err := tsdbClient.NewClient(opentsdbCfg)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

//BuildTask build task
func BuildTask(t *TaskStruct) (*pb.EnqueueRequest, error) {
	var er pb.EnqueueRequest
	taskJSON, err := json.Marshal(t.TaskBody)
	if err != nil {
		logrus.Errorf("tran task json error")
		return &er, err
	}
	er.Topic = "worker"
	er.Message = &pb.TaskMessage{
		TaskType:   t.TaskType,
		CreateTime: time.Now().Format(time.RFC3339),
		TaskBody:   taskJSON,
		User:       t.User,
	}
	return &er, nil
}

//GetBegin get db transaction
func GetBegin() *gorm.DB {
	return db.GetManager().Begin()
}

func dbInit() error {
	logrus.Info("api database initialization starting...")
	begin := GetBegin()
	// Permissions set
	var rac dbModel.RegionAPIClass
	if err := begin.Where("class_level=? and prefix=?", "server_source", "/v2/show").Find(&rac).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			data := map[string]string{
				"/v2/show":           "server_source",
				"/v2/cluster":        "server_source",
				"/v2/resources":      "server_source",
				"/v2/builder":        "server_source",
				"/v2/tenants":        "server_source",
				"/v2/app":            "server_source",
				"/v2/port":           "server_source",
				"/v2/volume-options": "server_source",
				"/api/v1":            "server_source",
				"/v2/events":         "server_source",
				"/v2/gateway/ips":    "server_source",
				"/v2/gateway/ports":  "server_source",
				"/v2/nodes":          "node_manager",
				"/v2/job":            "node_manager",
				"/v2/configs":        "node_manager",
			}
			tx := begin
			var rollback bool
			for k, v := range data {
				if err := db.GetManager().RegionAPIClassDaoTransactions(tx).AddModel(&dbModel.RegionAPIClass{
					ClassLevel: v,
					Prefix:     k,
				}); err != nil {
					tx.Rollback()
					rollback = true
					break
				}
			}
			if !rollback {
				tx.Commit()
			}
		} else {
			return err
		}
	}

	//Port Protocol support
	var rps dbModel.RegionProcotols
	if err := begin.Where("protocol_group=? and protocol_child=?", "http", "http").Find(&rps).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			data := map[string][]string{
				"http":   []string{"http"},
				"stream": []string{"mysql", "tcp", "udp"},
			}
			tx := begin
			var rollback bool
			for k, v := range data {
				for _, v1 := range v {
					if err := db.GetManager().RegionProcotolsDaoTransactions(tx).AddModel(&dbModel.RegionProcotols{
						ProtocolGroup: k,
						ProtocolChild: v1,
						APIVersion:    "v2",
						IsSupport:     true,
					}); err != nil {
						tx.Rollback()
						rollback = true
						break
					}
				}
			}
			if !rollback {
				tx.Commit()
			}
		} else {
			return err
		}
	}

	return nil
}

func dataInitialization() {
	timer := time.NewTimer(time.Second * 2)
	defer timer.Stop()
	for {
		err := dbInit()
		if err != nil {
			logrus.Error("Initializing database failed, ", err)
		} else {
			logrus.Info("api database initialization success!")
			return
		}
		select {
		case <-timer.C:
			timer.Reset(time.Second * 2)
		}
	}
}
