package dao

import (
	"github.com/gridworkz/kato/db/model"
	"github.com/jinzhu/gorm"
)

//EnterpriseDaoImpl
type EnterpriseDaoImpl struct {
	DB *gorm.DB
}

// GetEnterpriseTenants
func (e *EnterpriseDaoImpl) GetEnterpriseTenants(enterpriseID string) ([]*model.Tenants, error) {
	var tenants []*model.Tenants
	if enterpriseID == "" {
		return []*model.Tenants{}, nil
	}
	if err := e.DB.Where("eid= ?", enterpriseID).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}
