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

package exector

import (
	"fmt"
	"os"

	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/event"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"

	//"github.com/docker/docker/api/types"

	//"github.com/docker/docker/client"

	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
)

//MarketSlugItem MarketSlugItem
type MarketSlugItem struct {
	TenantName    string       `json:"tenant_name"`
	ServiceAlias  string       `json:"service_alias"`
	Logger        event.Logger `json:"logger"`
	EventID       string       `json:"event_id"`
	Operator      string       `json:"operator"`
	DeployVersion string       `json:"deploy_version"`
	TenantID      string       `json:"tenant_id"`
	ServiceID     string       `json:"service_id"`
	Action        string       `json:"action"`
	TGZPath       string
	Configs       map[string]string `json:"configs"`
	SlugInfo      struct {
		SlugPath    string `json:"slug_path"`
		FTPHost     string `json:"ftp_host"`
		FTPPort     string `json:"ftp_port"`
		FTPUser     string `json:"ftp_username"`
		FTPPassword string `json:"ftp_password"`
	} `json:"slug_info"`
}

//NewMarketSlugItem
func NewMarketSlugItem(in []byte) (*MarketSlugItem, error) {
	var msi MarketSlugItem
	if err := ffjson.Unmarshal(in, &msi); err != nil {
		return nil, err
	}
	msi.Logger = event.GetManager().GetLogger(msi.EventID)
	msi.TGZPath = fmt.Sprintf("/grdata/build/tenant/%s/slug/%s/%s.tgz", msi.TenantID, msi.ServiceID, msi.DeployVersion)
	return &msi, nil
}

//Run
func (i *MarketSlugItem) Run() error {
	if i.SlugInfo.FTPHost != "" && i.SlugInfo.FTPPort != "" {
		sFTPClient, err := sources.NewSFTPClient(i.SlugInfo.FTPUser, i.SlugInfo.FTPPassword, i.SlugInfo.FTPHost, i.SlugInfo.FTPPort)
		if err != nil {
			i.Logger.Error("failed to create FTP client", map[string]string{"step": "slug-share", "status": "failure"})
			return err
		}
		defer sFTPClient.Close()
		if err := sFTPClient.DownloadFile(i.SlugInfo.SlugPath, i.TGZPath, i.Logger); err != nil {
			i.Logger.Error("remote FTP acquisition of source code package failed, installation failed", map[string]string{"step": "slug-share", "status": "failure"})
			logrus.Errorf("copy slug file error when build service, %s", err.Error())
			return nil
		}
	} else {
		if err := sources.CopyFileWithProgress(i.SlugInfo.SlugPath, i.TGZPath, i.Logger); err != nil {
			i.Logger.Error("failed to obtain the source code package locally, and failed to install", map[string]string{"step": "slug-share", "status": "failure"})
			logrus.Errorf("copy slug file error when build service, %s", err.Error())
			return nil
		}
	}
	if err := os.Chown(i.TGZPath, 200, 200); err != nil {
		os.Remove(i.TGZPath)
		i.Logger.Error("failed to obtain the source code package locally, and failed to install", map[string]string{"step": "slug-share", "status": "failure"})
		logrus.Errorf("chown slug file error when build service, %s", err.Error())
		return nil
	}
	i.Logger.Info("application built", map[string]string{"step": "build-code", "status": "success"})
	vi := &dbmodel.VersionInfo{
		DeliveredType: "slug",
		DeliveredPath: i.TGZPath,
		EventID:       i.EventID,
		FinalStatus:   "success",
	}
	if err := i.UpdateVersionInfo(vi); err != nil {
		logrus.Errorf("update version info error: %s", err.Error())
		i.Logger.Error("failed to update app version information", map[string]string{"step": "slug-share", "status": "failure"})
		return err
	}
	return nil
}

//UpdateVersionInfo
func (i *MarketSlugItem) UpdateVersionInfo(vi *dbmodel.VersionInfo) error {
	version, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(i.DeployVersion, i.ServiceID)
	if err != nil {
		return err
	}
	if vi.DeliveredType != "" {
		version.DeliveredType = vi.DeliveredType
	}
	if vi.DeliveredPath != "" {
		version.DeliveredPath = vi.DeliveredPath
	}
	if vi.EventID != "" {
		version.EventID = vi.EventID
	}
	if vi.FinalStatus != "" {
		version.FinalStatus = vi.FinalStatus
	}
	if vi.DeliveredType == "slug" {
		version.ImageName = builder.RUNNERIMAGENAME
	}
	if err := db.GetManager().VersionInfoDao().UpdateModel(version); err != nil {
		return err
	}
	return nil
}
