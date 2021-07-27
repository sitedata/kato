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
	"errors"
	"fmt"
	"time"

	"github.com/gridworkz/kato/db/config"
	"github.com/gridworkz/kato/db/dao"
	"github.com/gridworkz/kato/db/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//Manager db manager
type Manager interface {
	CloseManager() error
	Begin() *gorm.DB
	DB() *gorm.DB
	EnsureEndTransactionFunc() func(tx *gorm.DB)
	VolumeTypeDao() dao.VolumeTypeDao
	LicenseDao() dao.LicenseDao
	AppDao() dao.AppDao
	ApplicationDao() dao.ApplicationDao
	ApplicationDaoTransactions(db *gorm.DB) dao.ApplicationDao
	AppConfigGroupDao() dao.AppConfigGroupDao
	AppConfigGroupDaoTransactions(db *gorm.DB) dao.AppConfigGroupDao
	AppConfigGroupServiceDao() dao.AppConfigGroupServiceDao
	AppConfigGroupServiceDaoTransactions(db *gorm.DB) dao.AppConfigGroupServiceDao
	AppConfigGroupItemDao() dao.AppConfigGroupItemDao
	AppConfigGroupItemDaoTransactions(db *gorm.DB) dao.AppConfigGroupItemDao
	EnterpriseDao() dao.EnterpriseDao
	TenantDao() dao.TenantDao
	TenantDaoTransactions(db *gorm.DB) dao.TenantDao
	TenantServiceDao() dao.TenantServiceDao
	TenantServiceDeleteDao() dao.TenantServiceDeleteDao
	TenantServiceDaoTransactions(db *gorm.DB) dao.TenantServiceDao
	TenantServiceDeleteDaoTransactions(db *gorm.DB) dao.TenantServiceDeleteDao
	TenantServicesPortDao() dao.TenantServicesPortDao
	TenantServicesPortDaoTransactions(*gorm.DB) dao.TenantServicesPortDao
	TenantServiceRelationDao() dao.TenantServiceRelationDao
	TenantServiceRelationDaoTransactions(*gorm.DB) dao.TenantServiceRelationDao
	TenantServiceEnvVarDao() dao.TenantServiceEnvVarDao
	TenantServiceEnvVarDaoTransactions(*gorm.DB) dao.TenantServiceEnvVarDao
	TenantServiceMountRelationDao() dao.TenantServiceMountRelationDao
	TenantServiceMountRelationDaoTransactions(db *gorm.DB) dao.TenantServiceMountRelationDao
	TenantServiceVolumeDao() dao.TenantServiceVolumeDao
	TenantServiceVolumeDaoTransactions(*gorm.DB) dao.TenantServiceVolumeDao
	TenantServiceConfigFileDao() dao.TenantServiceConfigFileDao
	TenantServiceConfigFileDaoTransactions(*gorm.DB) dao.TenantServiceConfigFileDao
	ServiceProbeDao() dao.ServiceProbeDao
	ServiceProbeDaoTransactions(*gorm.DB) dao.ServiceProbeDao
	TenantServiceLBMappingPortDao() dao.TenantServiceLBMappingPortDao
	TenantServiceLBMappingPortDaoTransactions(*gorm.DB) dao.TenantServiceLBMappingPortDao
	TenantServiceLabelDao() dao.TenantServiceLabelDao
	TenantServiceLabelDaoTransactions(db *gorm.DB) dao.TenantServiceLabelDao
	LocalSchedulerDao() dao.LocalSchedulerDao
	TenantPluginDaoTransactions(db *gorm.DB) dao.TenantPluginDao
	TenantPluginDao() dao.TenantPluginDao
	TenantPluginDefaultENVDaoTransactions(db *gorm.DB) dao.TenantPluginDefaultENVDao
	TenantPluginDefaultENVDao() dao.TenantPluginDefaultENVDao
	TenantPluginBuildVersionDao() dao.TenantPluginBuildVersionDao
	TenantPluginBuildVersionDaoTransactions(db *gorm.DB) dao.TenantPluginBuildVersionDao
	TenantPluginVersionENVDao() dao.TenantPluginVersionEnvDao
	TenantPluginVersionENVDaoTransactions(db *gorm.DB) dao.TenantPluginVersionEnvDao
	TenantPluginVersionConfigDao() dao.TenantPluginVersionConfigDao
	TenantPluginVersionConfigDaoTransactions(db *gorm.DB) dao.TenantPluginVersionConfigDao
	TenantServicePluginRelationDao() dao.TenantServicePluginRelationDao
	TenantServicePluginRelationDaoTransactions(db *gorm.DB) dao.TenantServicePluginRelationDao
	TenantServicesStreamPluginPortDao() dao.TenantServicesStreamPluginPortDao
	TenantServicesStreamPluginPortDaoTransactions(db *gorm.DB) dao.TenantServicesStreamPluginPortDao

	CodeCheckResultDao() dao.CodeCheckResultDao
	CodeCheckResultDaoTransactions(db *gorm.DB) dao.CodeCheckResultDao

	ServiceEventDao() dao.EventDao
	ServiceEventDaoTransactions(db *gorm.DB) dao.EventDao

	VersionInfoDao() dao.VersionInfoDao
	VersionInfoDaoTransactions(db *gorm.DB) dao.VersionInfoDao

	RegionUserInfoDao() dao.RegionUserInfoDao
	RegionUserInfoDaoTransactions(db *gorm.DB) dao.RegionUserInfoDao

	RegionAPIClassDao() dao.RegionAPIClassDao
	RegionAPIClassDaoTransactions(db *gorm.DB) dao.RegionAPIClassDao

	NotificationEventDao() dao.NotificationEventDao
	AppBackupDao() dao.AppBackupDao
	AppBackupDaoTransactions(db *gorm.DB) dao.AppBackupDao
	ServiceSourceDao() dao.ServiceSourceDao

	// gateway
	CertificateDao() dao.CertificateDao
	CertificateDaoTransactions(db *gorm.DB) dao.CertificateDao
	RuleExtensionDao() dao.RuleExtensionDao
	RuleExtensionDaoTransactions(db *gorm.DB) dao.RuleExtensionDao
	HTTPRuleDao() dao.HTTPRuleDao
	HTTPRuleDaoTransactions(db *gorm.DB) dao.HTTPRuleDao
	TCPRuleDao() dao.TCPRuleDao
	TCPRuleDaoTransactions(db *gorm.DB) dao.TCPRuleDao
	GwRuleConfigDao() dao.GwRuleConfigDao
	GwRuleConfigDaoTransactions(db *gorm.DB) dao.GwRuleConfigDao

	// third-party service
	EndpointsDao() dao.EndpointsDao
	EndpointsDaoTransactions(db *gorm.DB) dao.EndpointsDao
	ThirdPartySvcDiscoveryCfgDao() dao.ThirdPartySvcDiscoveryCfgDao
	ThirdPartySvcDiscoveryCfgDaoTransactions(db *gorm.DB) dao.ThirdPartySvcDiscoveryCfgDao

	TenantServceAutoscalerRulesDao() dao.TenantServceAutoscalerRulesDao
	TenantServceAutoscalerRulesDaoTransactions(db *gorm.DB) dao.TenantServceAutoscalerRulesDao
	TenantServceAutoscalerRuleMetricsDao() dao.TenantServceAutoscalerRuleMetricsDao
	TenantServceAutoscalerRuleMetricsDaoTransactions(db *gorm.DB) dao.TenantServceAutoscalerRuleMetricsDao
	TenantServiceScalingRecordsDao() dao.TenantServiceScalingRecordsDao
	TenantServiceScalingRecordsDaoTransactions(db *gorm.DB) dao.TenantServiceScalingRecordsDao

	TenantServiceMonitorDao() dao.TenantServiceMonitorDao
	TenantServiceMonitorDaoTransactions(db *gorm.DB) dao.TenantServiceMonitorDao
}

var defaultManager Manager

var supportDrivers map[string]struct{}

func init() {
	supportDrivers = map[string]struct{}{
		"mysql":       {},
		"cockroachdb": {},
	}
}

//CreateManager Create Manager
func CreateManager(config config.Config) (err error) {
	if _, ok := supportDrivers[config.DBType]; !ok {
		return fmt.Errorf("DB drivers: %s not supported", config.DBType)
	}

	for {
		defaultManager, err = mysql.CreateManager(config)
		if err == nil {
			logrus.Infof("db manager is ready")
			break
		}
		logrus.Errorf("get db manager failed, try time is %d,%s", 10, err.Error())
		time.Sleep(10 * time.Second)
	}
	//TODO:etcd db plugin
	//defaultManager, err = etcd.CreateManager(config)
	return
}

//CloseManager close db manager
func CloseManager() error {
	if defaultManager == nil {
		return errors.New("default db manager not init")
	}
	return defaultManager.CloseManager()
}

//GetManager get db manager
func GetManager() Manager {
	return defaultManager
}

// SetTestManager sets the default manager for unit test
func SetTestManager(m Manager) {
	defaultManager = m
}
