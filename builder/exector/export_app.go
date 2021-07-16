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

package exector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"time"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato-oam/pkg/export"
	"github.com/gridworkz/kato-oam/pkg/ram/v1alpha1"
	ramv1alpha1 "github.com/gridworkz/kato-oam/pkg/ram/v1alpha1"
	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/event"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var re = regexp.MustCompile(`\s`)

//ExportApp - export app to specified format (kato-app or dockercompose)
type ExportApp struct {
	EventID      string `json:"event_id"`
	Format       string `json:"format"`
	SourceDir    string `json:"source_dir"`
	Logger       event.Logger
	DockerClient *client.Client
}

func init() {
	RegisterWorker("export_app", NewExportApp)
}

//NewExportApp create
func NewExportApp(in []byte, m *exectorManager) (TaskWorker, error) {
	eventID := gjson.GetBytes(in, "event_id").String()
	logger := event.GetManager().GetLogger(eventID)
	return &ExportApp{
		Format:       gjson.GetBytes(in, "format").String(),
		SourceDir:    gjson.GetBytes(in, "source_dir").String(),
		Logger:       logger,
		EventID:      eventID,
		DockerClient: m.DockerClient,
	}, nil
}

//Run
func (i *ExportApp) Run(timeout time.Duration) error {
	// disable Md5 checksum
	// if ok := i.isLatest(); ok {
	// 	i.updateStatus("success")
	// 	return nil
	// }
	// Delete the old application group directory and then regenerate the application package
	if err := i.CleanSourceDir(); err != nil {
		return err
	}

	ram, err := i.parseRAM()
	if err != nil {
		return err
	}
	i.handleDefaultRepo(ram)
	var re *export.Result
	if i.Format == "kato-app" {
		re, err = i.exportKatoAPP(*ram)
		if err != nil {
			logrus.Errorf("export kato app package failure %s", err.Error())
			i.updateStatus("failed", "")
			return err
		}
	} else if i.Format == "docker-compose" {
		re, err = i.exportDockerCompose(*ram)
		if err != nil {
			logrus.Errorf("export docker compose app package failure %s", err.Error())
			i.updateStatus("failed", "")
			return err
		}
	} else {
		return errors.New("Unsupported the format: " + i.Format)
	}
	if re != nil {
		// move package file to download dir
		downloadPath := path.Dir(i.SourceDir)
		os.Rename(re.PackagePath, path.Join(downloadPath, re.PackageName))
		packageDownloadPath := path.Join("/v2/app/download/", i.Format, re.PackageName)
		// update export event status
		if err := i.updateStatus("success", packageDownloadPath); err != nil {
			return err
		}
		logrus.Infof("move export package %s to download dir success", re.PackageName)
		i.cacheMd5()
	}

	return nil
}

func (i *ExportApp) handleDefaultRepo(ram *v1alpha1.KatoApplicationConfig) {
	for i := range ram.Components {
		com := ram.Components[i]
		com.AppImage.HubUser, com.AppImage.HubPassword = builder.GetImageUserInfoV2(
			com.ShareImage, com.AppImage.HubUser, com.AppImage.HubPassword)
	}
	for i := range ram.Plugins {
		plugin := ram.Plugins[i]
		plugin.PluginImage.HubUser, plugin.PluginImage.HubPassword = builder.GetImageUserInfoV2(
			plugin.ShareImage, plugin.PluginImage.HubUser, plugin.PluginImage.HubPassword)
	}
}

// create md5 file
func (i *ExportApp) cacheMd5() {
	metadataFile := fmt.Sprintf("%s/metadata.json", i.SourceDir)
	if err := exec.Command("sh", "-c", fmt.Sprintf("md5sum %s > %s.md5", metadataFile, metadataFile)).Run(); err != nil {
		err = errors.New(fmt.Sprintf("Failed to create md5 file: %v", err))
		logrus.Error(err)
        }
	logrus.Infof("create md5 file success")
}

// exportKatoAPP export offline kato app
func (i *ExportApp) exportKatoAPP(ram v1alpha1.KatoApplicationConfig) (*export.Result, error) {
	ramExporter := export.New(export.RAM, i.SourceDir, ram, i.DockerClient, logrus.StandardLogger())
	return ramExporter.Export()
}

// exportDockerCompose export app to docker compose app
func (i *ExportApp) exportDockerCompose(ram v1alpha1.RainbondApplicationConfig) (*export.Result, error) {
	ramExporter := export.New(export.DC, i.SourceDir, ram, i.DockerClient, logrus.StandardLogger())
	return ramExporter.Export()
}

//Stop
func (i *ExportApp) Stop() error {
	return nil
}

//Name return worker name
func (i *ExportApp) Name() string {
	return "export_app"
}

//GetLogger
func (i *ExportApp) GetLogger() event.Logger {
	return i.Logger
}

// isLatest Returns true if the application is packaged and up to date
func (i *ExportApp) isLatest() bool {
	md5File := fmt.Sprintf("%s/metadata.json.md5", i.SourceDir)
	if _, err := os.Stat(md5File); os.IsNotExist(err) {
		logrus.Debug("The export app md5 file not found: ", md5File)
		return false
	}
	err := exec.Command("md5sum", "-c", md5File).Run()
	if err != nil {
		tarFile := i.SourceDir + ".tar"
		if _, err := os.Stat(tarFile); os.IsNotExist(err) {
			logrus.Debug("The export app tar file not found. ")
			return false
		}
		logrus.Info("The export app tar file is not the latest.")
		return false
	}
	logrus.Info("The export app tar file is the latest.")
	return true
}

//CleanSourceDir clean export dir
func (i *ExportApp) CleanSourceDir() error {
	logrus.Debug("Ready clean the source directory.")
	metaFile := fmt.Sprintf("%s/metadata.json", i.SourceDir)

	data, err := ioutil.ReadFile(metaFile)
	if err != nil {
		logrus.Error("Failed to read metadata file: ", err)
		return err
	}

	os.RemoveAll(i.SourceDir)
	os.MkdirAll(i.SourceDir, 0755)

	if err := ioutil.WriteFile(metaFile, data, 0644); err != nil {
		logrus.Error("Failed to write metadata file: ", err)
		return err
	}

	return nil
}
func (i *ExportApp) parseRAM() (*ramv1alpha1.KatoApplicationConfig, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/metadata.json", i.SourceDir))
	if err != nil {
		i.Logger.Error("Failed to export application, application information was not found", map[string]string{"step": "read-metadata", "status": "failure"})
		logrus.Error("Failed to read metadata file: ", err)
		return nil, err
	}
	var ram ramv1alpha1.KatoApplicationConfig
	if err := json.Unmarshal(data, &ram); err != nil {
		return nil, err
	}
	return &ram, nil
}

func (i *ExportApp) updateStatus(status, filePath string) error {
	logrus.Debug("Update app status in database to: ", status)
	res, err := db.GetManager().AppDao().GetByEventId(i.EventID)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to get app %s from db: %v", i.EventID, err))
		logrus.Error(err)
		return err
	}
	res.Status = status
	if filePath != "" {
		res.TarFileHref = filePath
	}
	if err := db.GetManager().AppDao().UpdateModel(res); err != nil {
		err = errors.New(fmt.Sprintf("Failed to update app %s: %v", i.EventID, err))
		logrus.Error(err)
		return err
	}
	return nil
}

//ErrorCallBack if run error will callback
func (i *ExportApp) ErrorCallBack(err error) {
	i.updateStatus("failed", "")
}
