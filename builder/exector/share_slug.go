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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/event"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
)

//SlugShareItem
type SlugShareItem struct {
	Namespace     string `json:"namespace"`
	TenantName    string `json:"tenant_name"`
	ServiceID     string `json:"service_id"`
	ServiceAlias  string `json:"service_alias"`
	SlugPath      string `json:"slug_path"`
	LocalSlugPath string `json:"local_slug_path"`
	ShareID       string `json:"share_id"`
	Logger        event.Logger
	ShareInfo     struct {
		ServiceKey string `json:"service_key" `
		AppVersion string `json:"app_version" `
		EventID    string `json:"event_id"`
		ShareUser  string `json:"share_user"`
		ShareScope string `json:"share_scope"`
		SlugInfo   struct {
			Namespace   string `json:"namespace"`
			FTPHost     string `json:"ftp_host"`
			FTPPort     string `json:"ftp_port"`
			FTPUser     string `json:"ftp_username"`
			FTPPassword string `json:"ftp_password"`
		} `json:"slug_info,omitempty"`
	} `json:"share_info"`
	EtcdCli     *clientv3.Client
	PackageName string
}

//NewSlugShareItem
func NewSlugShareItem(in []byte, etcdCli *clientv3.Client) (*SlugShareItem, error) {
	var ssi SlugShareItem
	if err := ffjson.Unmarshal(in, &ssi); err != nil {
		return nil, err
	}
	eventID := ssi.ShareInfo.EventID
	ssi.Logger = event.GetManager().GetLogger(eventID)
	ssi.EtcdCli = etcdCli
	return &ssi, nil
}

//ShareService
func (i *SlugShareItem) ShareService() error {

	logrus.Debugf("share app local slug path: %s ,target path: %s", i.LocalSlugPath, i.SlugPath)
	if _, err := os.Stat(i.LocalSlugPath); err != nil {
		i.Logger.Error(fmt.Sprintf("The data center application code package does not exist, please build the application first"), map[string]string{"step": "slug-share", "status": "failure"})
		return err
	}
	if i.ShareInfo.SlugInfo.FTPHost != "" && i.ShareInfo.SlugInfo.FTPPort != "" {
		if err := i.ShareToFTP(); err != nil {
			return err
		}
	} else {
		if err := i.ShareToLocal(); err != nil {
			return err
		}
	}
	return nil
}

func createMD5(packageName string) (string, error) {
	md5Path := packageName + ".md5"
	_, err := os.Stat(md5Path)
	if err == nil {
		//md5 file exist
		return md5Path, nil
	}
	f, err := exec.Command("md5sum", packageName).Output()
	if err != nil {
		return "", err
	}
	md5In := strings.Split(string(f), "")
	if err := ioutil.WriteFile(md5Path, []byte(md5In[0]), 0644); err != nil {
		return "", err
	}
	return md5Path, nil
}

//ShareToFTP
func (i *SlugShareItem) ShareToFTP() error {
	i.Logger.Info("Start uploading application media to FTP server", map[string]string{"step": "slug-share"})
	sFTPClient, err := sources.NewSFTPClient(i.ShareInfo.SlugInfo.FTPUser, i.ShareInfo.SlugInfo.FTPPassword, i.ShareInfo.SlugInfo.FTPHost, i.ShareInfo.SlugInfo.FTPPort)
	if err != nil {
		i.Logger.Error("Failed to create FTP client", map[string]string{"step": "slug-share", "status": "failure"})
		return err
	}
	defer sFTPClient.Close()
	if err := sFTPClient.PushFile(i.LocalSlugPath, i.SlugPath, i.Logger); err != nil {
		i.Logger.Error("Failed to upload source package file", map[string]string{"step": "slug-share", "status": "failure"})
		return err
	}
	i.Logger.Info("Share gridworkz cloud remote server completed", map[string]string{"step": "slug-share", "status": "success"})
	return nil
}

//ShareToLocal
func (i *SlugShareItem) ShareToLocal() error {
	file := i.LocalSlugPath
	i.Logger.Info("Start sharing the application to the local directory", map[string]string{"step": "slug-share"})
	md5, err := createMD5(file)
	if err != nil {
		i.Logger.Error("Failed to generate md5", map[string]string{"step": "slug-share", "status": "success"})
		return err
	}
	if err := sources.CopyFileWithProgress(i.LocalSlugPath, i.SlugPath, i.Logger); err != nil {
		os.Remove(i.SlugPath)
		logrus.Errorf("copy file to share path error: %s", err.Error())
		i.Logger.Error("Failed to copy file", map[string]string{"step": "slug-share", "status": "failure"})
		return err
	}
	if err := sources.CopyFileWithProgress(md5, i.SlugPath+".md5", i.Logger); err != nil {
		os.Remove(i.SlugPath)
		os.Remove(i.SlugPath + ".md5")
		logrus.Errorf("copy file to share path error: %s", err.Error())
		i.Logger.Error("Failed to copy md5 file", map[string]string{"step": "slug-share", "status": "failure"})
		return err
	}
	i.Logger.Info("Shared data center completed locally", map[string]string{"step": "slug-share", "status": "success"})
	return nil
}

//UpdateShareStatus
func (i *SlugShareItem) UpdateShareStatus(status string) error {
	var ss = ShareStatus{
		ShareID: i.ShareID,
		Status:  status,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := i.EtcdCli.Put(ctx, fmt.Sprintf("/kato/shareresult/%s", i.ShareID), ss.String())
	if err != nil {
		logrus.Errorf("put shareresult  %s into etcd error, %v", i.ShareID, err)
		i.Logger.Error("Failed to save sharing results.", map[string]string{"step": "callback", "status": "failure"})
	}
	if status == "success" {
		i.Logger.Info("Create and share results are successful, shared successfully", map[string]string{"step": "last", "status": "success"})
	} else {
		i.Logger.Info("Create and share results succeeded, share failed", map[string]string{"step": "callback", "status": "failure"})
	}
	return nil
}

//CheckMD5FileExist
func (i *SlugShareItem) CheckMD5FileExist(md5path, packageName string) bool {
	return false
}
