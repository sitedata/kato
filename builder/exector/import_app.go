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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato-oam/pkg/localimport"
	"github.com/gridworkz/kato-oam/pkg/ram/v1alpha1"
	"github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/event"
	"github.com/sirupsen/logrus"
)

func init() {
	RegisterWorker("import_app", NewImportApp)
}

//ImportApp Export app to specified format(kato-app or dockercompose)
type ImportApp struct {
	EventID       string             `json:"event_id"`
	Format        string             `json:"format"`
	SourceDir     string             `json:"source_dir"`
	Apps          []string           `json:"apps"`
	ServiceImage  model.ServiceImage `json:"service_image"`
	Logger        event.Logger
	DockerClient  *client.Client
	oldAPPPath    map[string]string
	oldPluginPath map[string]string
}

//NewImportApp create
func NewImportApp(in []byte, m *exectorManager) (TaskWorker, error) {
	var importApp ImportApp
	if err := json.Unmarshal(in, &importApp); err != nil {
		return nil, err
	}
	if importApp.ServiceImage.HubURL == "" || importApp.ServiceImage.HubURL == "gridworkz" {
		importApp.ServiceImage.HubURL = builder.REGISTRYDOMAIN
		importApp.ServiceImage.HubUser = builder.REGISTRYUSER
		importApp.ServiceImage.HubPassword = builder.REGISTRYPASS
	}
	logrus.Infof("load app image to hub %s", importApp.ServiceImage.HubURL)
	importApp.Logger = event.GetManager().GetLogger(importApp.EventID)
	importApp.DockerClient = m.DockerClient
	importApp.oldAPPPath = make(map[string]string)
	importApp.oldPluginPath = make(map[string]string)
	return &importApp, nil
}

//Stop
func (i *ImportApp) Stop() error {
	return nil
}

//Name return worker name
func (i *ImportApp) Name() string {
	return "import_app"
}

//GetLogger
func (i *ImportApp) GetLogger() event.Logger {
	return i.Logger
}

//ErrorCallBack if run error will callback
func (i *ImportApp) ErrorCallBack(err error) {
	i.updateStatus("failed")
}

//Run
func (i *ImportApp) Run(timeout time.Duration) error {
	if i.Format == "kato-app" {
		err := i.importApp()
		if err != nil {
			logrus.Errorf("load kato app failure %s", err.Error())
			return err
		}
		return nil
	}
	return errors.New("Unsupported the format: " + i.Format)
}

// importApp
// support batch import
func (i *ImportApp) importApp() error {
	oldSourceDir := i.SourceDir
	var datas []v1alpha1.KatoApplicationConfig
	var wait sync.WaitGroup
	for _, app := range i.Apps {
		wait.Add(1)
		go func(app string) {
			defer wait.Done()
			appFile := filepath.Join(oldSourceDir, app)
			tmpDir := path.Join(oldSourceDir, app+"-cache")
			li := localimport.New(logrus.StandardLogger(), i.DockerClient, tmpDir)
			if err := i.updateStatusForApp(app, "importing"); err != nil {
				logrus.Errorf("Failed to update status to importing for app %s: %v", app, err)
			}
			ram, err := li.Import(appFile, v1alpha1.ImageInfo{
				HubURL:      i.ServiceImage.HubURL,
				HubUser:     i.ServiceImage.HubUser,
				HubPassword: i.ServiceImage.HubPassword,
				Namespace:   i.ServiceImage.NameSpace,
			})
			if err != nil {
				logrus.Errorf("Failed to load app %s: %v", appFile, err)
				i.updateStatusForApp(app, "failed")
				return
			}
			os.Rename(appFile, appFile+".success")
			datas = append(datas, *ram)
			logrus.Infof("Successful import app: %s", appFile)
			os.Remove(tmpDir)
		}(app)
	}
	wait.Wait()
	metadatasFile := fmt.Sprintf("%s/metadatas.json", i.SourceDir)
	dataBytes, _ := json.Marshal(datas)
	if err := ioutil.WriteFile(metadatasFile, []byte(dataBytes), 0644); err != nil {
		logrus.Errorf("Failed to load apps %s: %v", i.SourceDir, err)
		return err
	}
	if err := i.updateStatus("success"); err != nil {
		logrus.Errorf("Failed to load apps %s: %v", i.SourceDir, err)
		return err
	}
	return nil
}

func (i *ImportApp) updateStatus(status string) error {
	logrus.Debug("Update app status in database to: ", status)
	// Get the status information of the application from the database
	res, err := db.GetManager().AppDao().GetByEventId(i.EventID)
	if err != nil {
		err = fmt.Errorf("failed to get app %s from db: %s", i.EventID, err.Error())
		logrus.Error(err)
		return err
	}

	// Update the status information of the application in the database
	res.Status = status

	if err := db.GetManager().AppDao().UpdateModel(res); err != nil {
		err = fmt.Errorf("failed to update app %s: %s", i.EventID, err.Error())
		logrus.Error(err)
		return err
	}

	return nil
}

func (i *ImportApp) updateStatusForApp(app, status string) error {
	logrus.Debugf("Update status in database for app %s to: %s", app, status)
	// Get the status information of the application from the database
	res, err := db.GetManager().AppDao().GetByEventId(i.EventID)
	if err != nil {
		err = fmt.Errorf("Failed to get app %s from db: %s", i.EventID, err.Error())
		logrus.Error(err)
		return err
	}

	// Update the status information of the application in the database
	appsMap := str2map(res.Apps)
	appsMap[app] = status
	res.Apps = map2str(appsMap)

	if err := db.GetManager().AppDao().UpdateModel(res); err != nil {
		err = fmt.Errorf("Failed to update app %s: %s", i.EventID, err.Error())
		logrus.Error(err)
		return err
	}

	return nil
}

func str2map(str string) map[string]string {
	result := make(map[string]string, 10)

	for _, app := range strings.Split(str, ",") {
		appMap := strings.Split(app, ":")
		result[appMap[0]] = appMap[1]
	}

	return result
}

func map2str(m map[string]string) string {
	var result string

	for k, v := range m {
		kv := k + ":" + v

		if result == "" {
			result += kv
		} else {
			result += "," + kv
		}
	}

	return result
}

// Only keep the part after "/" and remove illegal characters, generally used to restore the exported image file to the image name
func buildFromLinuxFileName(fileName string) string {
	if fileName == "" {
		return fileName
	}

	arr := strings.Split(fileName, "/")

	if str := arr[len(arr)-1]; str == "" {
		fileName = strings.Replace(fileName, "---", "/", -1)
	} else {
		fileName = str
	}

	fileName = strings.Replace(fileName, "--", ":", -1)
	fileName = strings.TrimSpace(fileName)

	return fileName
}
