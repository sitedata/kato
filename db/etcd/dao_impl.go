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

package etcd

import (
	"github.com/gridworkz/kato/db/dao"
)

// TenantDao
func (m *Manager) TenantDao() dao.TenantDao {
	return nil
}

// TenantServiceDao
func (m *Manager) TenantServiceDao() dao.TenantServiceDao {
	return nil
}

// TenantServicesPortDao
func (m *Manager) TenantServicesPortDao() dao.TenantServicesPortDao {
	return nil
}

// TenantServiceRelationDao
func (m *Manager) TenantServiceRelationDao() dao.TenantServiceRelationDao {
	return nil
}

// TenantServiceEnvVarDao
func (m *Manager) TenantServiceEnvVarDao() dao.TenantServiceEnvVarDao {
	return nil
}

// TenantServiceMountRelationDao
func (m *Manager) TenantServiceMountRelationDao() dao.TenantServiceMountRelationDao {
	return nil
}

// TenantServiceVolumeDao
func (m *Manager) TenantServiceVolumeDao() dao.TenantServiceVolumeDao {
	return nil
}

// func (m *Manager) K8sServiceDao() dao.K8sServiceDao {
// 	return nil
// }
// func (m *Manager) K8sDeployReplicationDao() dao.K8sDeployReplicationDao {
// 	return nil
// }
// func (m *Manager) K8sPodDao() dao.K8sPodDao {
// 	return nil
// }

// ServiceProbeDao
func (m *Manager) ServiceProbeDao() dao.ServiceProbeDao {
	return nil
}
