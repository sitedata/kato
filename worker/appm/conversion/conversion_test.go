// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

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

package conversion

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/dao"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/util"
	"github.com/gridworkz/kato/worker/appm/types/v1"
	"testing"
)

func TestTenantServiceBase(t *testing.T) {
	t.Run("third-party service", func(t *testing.T) {
		as := &v1.AppService{}
		as.ServiceID = util.NewUUID()
		as.TenantID = util.NewUUID()
		as.TenantName = "abcdefg"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbm := db.NewMockManager(ctrl)
		// TenantServiceDao
		tenantServiceDao := dao.NewMockTenantServiceDao(ctrl)
		tenantService := &model.TenantServices{
			TenantID:  as.TenantID,
			ServiceID: as.ServiceID,
			Kind:      model.ServiceKindThirdParty.String(),
		}
		tenantServiceDao.EXPECT().GetServiceByID(as.ServiceID).Return(tenantService, nil)
		dbm.EXPECT().TenantServiceDao().Return(tenantServiceDao)
		// TenantDao
		tenantDao := dao.NewMockTenantDao(ctrl)
		tenant := &model.Tenants{
			UUID: as.TenantID,
			Name: as.TenantName,
		}
		tenantDao.EXPECT().GetTenantByUUID(as.TenantID).Return(tenant, nil)
		dbm.EXPECT().TenantDao().Return(tenantDao)
		if err := TenantServiceBase(as, dbm); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
