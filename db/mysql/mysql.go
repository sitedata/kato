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

package mysql

import (
	"sync"

	"github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/db/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	// import sql driver manually
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Manager - db manager
type Manager struct {
	db      *gorm.DB
	config  config.Config
	initOne sync.Once
	models  []model.Interface
}

//CreateManager
func CreateManager(config config.Config) (*Manager, error) {
	var db *gorm.DB
	if config.DBType == "mysql" {
		var err error
		db, err = gorm.Open("mysql", config.MysqlConnectionInfo+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			return nil, err
		}
	}
	if config.DBType == "yugabytedb" {
		var err error
		addr := config.MysqlConnectionInfo
		db, err = gorm.Open("postgres", addr)
		if err != nil {
			return nil, err
		}
	}
	if config.ShowSQL {
		db = db.Debug()
	}
	manager := &Manager{
		db:      db,
		config:  config,
		initOne: sync.Once{},
	}
	db.SetLogger(manager)
	manager.RegisterTableModel()
	manager.CheckTable()
	logrus.Debug("mysql db driver create")
	return manager, nil
}

//CloseManager
func (m *Manager) CloseManager() error {
	return m.db.Close()
}

//Begin a transaction
func (m *Manager) Begin() *gorm.DB {
	return m.db.Begin()
}

// DB returns the db.
func (m *Manager) DB() *gorm.DB {
	return m.db
}

// EnsureEndTransactionFunc
func (m *Manager) EnsureEndTransactionFunc() func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		if r := recover(); r != nil {
			logrus.Errorf("Unexpected panic occurred, rollback transaction: %v", r)
			tx.Rollback()
		}
	}
}

//Print
func (m *Manager) Print(v ...interface{}) {
	logrus.Info(v...)
}

//RegisterTableModel
func (m *Manager) RegisterTableModel() {
	m.models = append(m.models, &model.Tenants{})
	m.models = append(m.models, &model.TenantServices{})
	m.models = append(m.models, &model.TenantServicesPort{})
	m.models = append(m.models, &model.TenantServiceRelation{})
	m.models = append(m.models, &model.TenantServiceEnvVar{})
	m.models = append(m.models, &model.TenantServiceMountRelation{})
	m.models = append(m.models, &model.TenantServiceVolume{})
	m.models = append(m.models, &model.TenantServiceLable{})
	m.models = append(m.models, &model.TenantServiceProbe{})
	m.models = append(m.models, &model.LicenseInfo{})
	m.models = append(m.models, &model.TenantServicesDelete{})
	m.models = append(m.models, &model.TenantServiceLBMappingPort{})
	m.models = append(m.models, &model.TenantPlugin{})
	m.models = append(m.models, &model.TenantPluginBuildVersion{})
	m.models = append(m.models, &model.TenantServicePluginRelation{})
	m.models = append(m.models, &model.TenantPluginVersionEnv{})
	m.models = append(m.models, &model.TenantPluginVersionDiscoverConfig{})
	m.models = append(m.models, &model.CodeCheckResult{})
	m.models = append(m.models, &model.ServiceEvent{})
	m.models = append(m.models, &model.VersionInfo{})
	m.models = append(m.models, &model.RegionUserInfo{})
	m.models = append(m.models, &model.TenantServicesStreamPluginPort{})
	m.models = append(m.models, &model.RegionAPIClass{})
	m.models = append(m.models, &model.RegionProcotols{})
	m.models = append(m.models, &model.LocalScheduler{})
	m.models = append(m.models, &model.NotificationEvent{})
	m.models = append(m.models, &model.AppStatus{})
	m.models = append(m.models, &model.AppBackup{})
	m.models = append(m.models, &model.ServiceSourceConfig{})
	m.models = append(m.models, &model.Application{})
	m.models = append(m.models, &model.ApplicationConfigGroup{})
	m.models = append(m.models, &model.ConfigGroupService{})
	m.models = append(m.models, &model.ConfigGroupItem{})
	// gateway
	m.models = append(m.models, &model.Certificate{})
	m.models = append(m.models, &model.RuleExtension{})
	m.models = append(m.models, &model.HTTPRule{})
	m.models = append(m.models, &model.TCPRule{})
	m.models = append(m.models, &model.TenantServiceConfigFile{})
	m.models = append(m.models, &model.Endpoint{})
	m.models = append(m.models, &model.ThirdPartySvcDiscoveryCfg{})
	m.models = append(m.models, &model.GwRuleConfig{})

	// volumeType
	m.models = append(m.models, &model.TenantServiceVolumeType{})
	// pod autoscaler
	m.models = append(m.models, &model.TenantServiceAutoscalerRules{})
	m.models = append(m.models, &model.TenantServiceAutoscalerRuleMetrics{})
	m.models = append(m.models, &model.TenantServiceScalingRecords{})
	m.models = append(m.models, &model.TenantServiceMonitor{})
}

//CheckTable
func (m *Manager) CheckTable() {
	m.initOne.Do(func() {
		for _, md := range m.models {
			if !m.db.HasTable(md) {
				if m.config.DBType == "mysql" {
					err := m.db.Set("gorm:table_options", "ENGINE=InnoDB charset=utf8").CreateTable(md).Error
					if err != nil {
						logrus.Errorf("auto create table %s to db error."+err.Error(), md.TableName())
					} else {
						logrus.Infof("auto create table %s to db success", md.TableName())
					}
				} else { //yugabytedb
					err := m.db.CreateTable(md).Error
					if err != nil {
						logrus.Errorf("auto create yugabytedb table %s to db error."+err.Error(), md.TableName())
					} else {
						logrus.Infof("auto create yugabytedb table %s to db success", md.TableName())
					}
				}
			} else {
				if err := m.db.AutoMigrate(md).Error; err != nil {
					logrus.Errorf("auto Migrate table %s to db error."+err.Error(), md.TableName())
				}
			}
		}
		m.patchTable()
	})
}

func (m *Manager) patchTable() {
	//modify tenant service env max size to 1024
	if err := m.db.Exec("alter table tenant_services_envs modify column attr_value varchar(1024);").Error; err != nil {
		logrus.Errorf("alter table tenant_services_envs error %s", err.Error())
	}

	if err := m.db.Exec("alter table tenant_services_event modify column request_body varchar(1024);").Error; err != nil {
		logrus.Errorf("alter table tenant_services_envent error %s", err.Error())
	}

	if err := m.db.Exec("update gateway_tcp_rule set ip=? where ip=?", "0.0.0.0", "").Error; err != nil {
		logrus.Errorf("update gateway_tcp_rule data error %s", err.Error())
	}
	if err := m.db.Exec("alter table tenant_services_volume modify column volume_type varchar(64);").Error; err != nil {
		logrus.Errorf("alter table tenant_services_volume error: %s", err.Error())
	}
}
