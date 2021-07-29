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
	"os"
	"path"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/build"
	"github.com/gridworkz/kato/builder/parser"
	"github.com/gridworkz/kato/builder/parser/code"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"k8s.io/client-go/kubernetes"
)

//SourceCodeBuildItem SouceCodeBuildItem
type SourceCodeBuildItem struct {
	Namespace     string       `json:"namespace"`
	TenantName    string       `json:"tenant_name"`
	GRDataPVCName string       `json:"gr_data_pvc_name"`
	CachePVCName  string       `json:"cache_pvc_name"`
	CacheMode     string       `json:"cache_mode"`
	CachePath     string       `json:"cache_path"`
	ServiceAlias  string       `json:"service_alias"`
	Action        string       `json:"action"`
	DestImage     string       `json:"dest_image"`
	Logger        event.Logger `json:"logger"`
	EventID       string       `json:"event_id"`
	CacheDir      string       `json:"cache_dir"`
	TGZDir        string       `json:"tgz_dir"`
	DockerClient  *client.Client
	KubeClient    kubernetes.Interface
	RbdNamespace  string
	RbdRepoName   string
	TenantID      string
	ServiceID     string
	DeployVersion string
	Lang          string
	Runtime       string
	BuildEnvs     map[string]string
	CodeSouceInfo sources.CodeSourceInfo
	RepoInfo      *sources.RepostoryBuildInfo
	commit        Commit
	Configs       map[string]gjson.Result `json:"configs"`
	Ctx           context.Context
}

//Commit code Commit
type Commit struct {
	Hash    string
	Author  string
	Message string
}

//NewSouceCodeBuildItem create
func NewSouceCodeBuildItem(in []byte) *SourceCodeBuildItem {
	eventID := gjson.GetBytes(in, "event_id").String()
	logger := event.GetManager().GetLogger(eventID)
	csi := sources.CodeSourceInfo{
		ServerType:    strings.Replace(gjson.GetBytes(in, "server_type").String(), " ", "", -1),
		RepositoryURL: gjson.GetBytes(in, "repo_url").String(),
		Branch:        gjson.GetBytes(in, "branch").String(),
		User:          gjson.GetBytes(in, "user").String(),
		Password:      gjson.GetBytes(in, "password").String(),
		TenantID:      gjson.GetBytes(in, "tenant_id").String(),
		ServiceID:     gjson.GetBytes(in, "service_id").String(),
	}
	envs := gjson.GetBytes(in, "envs").String()
	be := make(map[string]string)
	if err := ffjson.Unmarshal([]byte(envs), &be); err != nil {
		logrus.Errorf("unmarshal build envs error: %s", err.Error())
	}
	scb := &SourceCodeBuildItem{
		Namespace:     gjson.GetBytes(in, "tenant_id").String(),
		TenantName:    gjson.GetBytes(in, "tenant_name").String(),
		ServiceAlias:  gjson.GetBytes(in, "service_alias").String(),
		TenantID:      gjson.GetBytes(in, "tenant_id").String(),
		ServiceID:     gjson.GetBytes(in, "service_id").String(),
		Action:        gjson.GetBytes(in, "action").String(),
		DeployVersion: gjson.GetBytes(in, "deploy_version").String(),
		Logger:        logger,
		EventID:       eventID,
		CodeSouceInfo: csi,
		Lang:          gjson.GetBytes(in, "lang").String(),
		Runtime:       gjson.GetBytes(in, "runtime").String(),
		Configs:       gjson.GetBytes(in, "configs").Map(),
		BuildEnvs:     be,
	}
	scb.CacheDir = fmt.Sprintf("/cache/build/%s/cache/%s", scb.TenantID, scb.ServiceID)
	//scb.SourceDir = scb.CodeSouceInfo.GetCodeSourceDir()
	scb.TGZDir = fmt.Sprintf("/grdata/build/tenant/%s/slug/%s", scb.TenantID, scb.ServiceID)
	return scb
}

//Run Run
func (i *SourceCodeBuildItem) Run(timeout time.Duration) error {
	// 1.clone
	// 2.check dockerfile/ source_code
	// 3.build
	// 4.upload image /upload slug
	rbi, err := sources.CreateRepostoryBuildInfo(i.CodeSouceInfo.RepositoryURL, i.CodeSouceInfo.ServerType, i.CodeSouceInfo.Branch, i.TenantID, i.ServiceID)
	if err != nil {
		i.Logger.Error("Git project repository address format error", map[string]string{"step": "parse"})
		return err
	}
	i.RepoInfo = rbi
	if err := i.prepare(); err != nil {
		logrus.Errorf("prepare build code error: %s", err.Error())
		i.Logger.Error("Failed to prepare source code to build", map[string]string{"step": "builder-exector", "status": "failure"})
		return err
	}
	i.CodeSouceInfo.RepositoryURL = rbi.RepostoryURL
	switch i.CodeSouceInfo.ServerType {
	case "svn":
		csi := i.CodeSouceInfo
		svnclient := sources.NewClient(csi, rbi.GetCodeHome(), i.Logger)
		rs, err := svnclient.UpdateOrCheckout(rbi.BuildPath)
		if err != nil {
			logrus.Errorf("checkout svn code error: %s", err.Error())
			i.Logger.Error(fmt.Sprintf("Checkout svn code failed, please make sure the code can be downloaded properly"), map[string]string{"step": "builder-exector", "status": "failure"})
			return err
		}
		if rs.Logs == nil || len(rs.Logs.CommitEntrys) < 1 {
			logrus.Errorf("get code commit info error: %s", err.Error())
			i.Logger.Error(fmt.Sprintf("Failed to read code version information"), map[string]string{"step": "builder-exector", "status": "failure"})
			return err
		}
		i.commit = Commit{
			Hash:    rs.Logs.CommitEntrys[0].Revision,
			Message: rs.Logs.CommitEntrys[0].Msg,
			Author:  rs.Logs.CommitEntrys[0].Author,
		}
	case "oss":
		i.commit = Commit{}
	default:
		//default git
		rs, err := sources.GitCloneOrPull(i.CodeSouceInfo, rbi.GetCodeHome(), i.Logger, 5)
		if err != nil {
			logrus.Errorf("pull git code error: %s", err.Error())
			i.Logger.Error("Failed to pull the code, please make sure the code can be downloaded normally", map[string]string{"step": "builder-exector", "status": "failure"})
			return err
		}
		//get last commit
		commit, err := sources.GetLastCommit(rs)
		if err != nil || commit == nil {
			logrus.Errorf("get code commit info error: %s", err.Error())
			i.Logger.Error("Failed to read code version information", map[string]string{"step": "builder-exector", "status": "failure"})
			return err
		}
		i.commit = Commit{
			Hash:    commit.Hash.String(),
			Author:  commit.Author.Name,
			Message: commit.Message,
		}
	}
	// clean cache code
	defer func() {
		if err := os.RemoveAll(rbi.GetCodeHome()); err != nil {
			logrus.Warningf("remove source code: %v", err)
		}
	}()

	hash := i.commit.Hash
	if len(hash) >= 8 {
		hash = i.commit.Hash[0:7]
	}
	info := fmt.Sprintf("CodeVersion:%s Author:%s Commit:%s ", hash, i.commit.Author, i.commit.Message)
	i.Logger.Info(info, map[string]string{"step": "code-version"})
	if _, ok := i.BuildEnvs["REPARSE"]; ok {
		_, lang, err := parser.ReadRbdConfigAndLang(rbi)
		if err != nil {
			logrus.Errorf("reparse code lange error %s", err.Error())
			i.Logger.Error(fmt.Sprintf("Re-parse code language error"), map[string]string{"step": "builder-exector", "status": "failure"})
			return err
		}
		i.Lang = string(lang)
	}

	i.Logger.Info("pull or clone code successfully, start code build", map[string]string{"step": "codee-version"})
	res, err := i.codeBuild()
	if err != nil {
		if err.Error() == context.DeadlineExceeded.Error() {
			i.Logger.Error("Build app version from source code timeout, the maximum time is 60 minutes", map[string]string{"step": "builder-exector", "status": "failure"})
		} else {
			i.Logger.Error("Build app version from source code failure,"+err.Error(), map[string]string{"step": "builder-exector", "status": "failure"})
		}
		return err
	}
	if err := i.UpdateBuildVersionInfo(res); err != nil {
		return err
	}
	return nil
}

func (i *SourceCodeBuildItem) codeBuild() (*build.Response, error) {
	codeBuild, err := build.GetBuild(code.Lang(i.Lang))
	if err != nil {
		logrus.Errorf("get code build error: %s lang %s", err.Error(), i.Lang)
		i.Logger.Error(util.Translation("No way of compiling to support this source type was found"), map[string]string{"step": "builder-exector", "status": "failure"})
		return nil, err
	}
	hostAlias, err := i.getHostAlias()
	if err != nil {
		i.Logger.Error(util.Translation("get rbd-repo ip failure"), map[string]string{"step": "builder-exector", "status": "failure"})
		return nil, err
	}
	buildReq := &build.Request{
		RbdNamespace:  i.RbdNamespace,
		SourceDir:     i.RepoInfo.GetCodeBuildAbsPath(),
		CacheDir:      i.CacheDir,
		TGZDir:        i.TGZDir,
		RepositoryURL: i.RepoInfo.RepostoryURL,
		CodeSouceInfo: i.CodeSouceInfo,
		ServiceAlias:  i.ServiceAlias,
		ServiceID:     i.ServiceID,
		TenantID:      i.TenantID,
		ServerType:    i.CodeSouceInfo.ServerType,
		Runtime:       i.Runtime,
		Branch:        i.CodeSouceInfo.Branch,
		DeployVersion: i.DeployVersion,
		Commit:        build.Commit{User: i.commit.Author, Message: i.commit.Message, Hash: i.commit.Hash},
		Lang:          code.Lang(i.Lang),
		BuildEnvs:     i.BuildEnvs,
		Logger:        i.Logger,
		DockerClient:  i.DockerClient,
		KubeClient:    i.KubeClient,
		HostAlias:     hostAlias,
		Ctx:           i.Ctx,
		GRDataPVCName: i.GRDataPVCName,
		CachePVCName:  i.CachePVCName,
		CacheMode:     i.CacheMode,
		CachePath:     i.CachePath,
	}
	res, err := codeBuild.Build(buildReq)
	return res, err
}

func (i *SourceCodeBuildItem) getHostAlias() (hostAliasList []build.HostAlias, err error) {
	endpoints, err := i.KubeClient.CoreV1().Endpoints(i.RbdNamespace).Get(context.Background(), i.RbdRepoName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("do not found ep by name: %s in namespace: %s", i.RbdRepoName, i.Namespace)
		return nil, err
	}
	hostNames := []string{"maven.gridworkz.me", "lang.gridworkz.me"}
	for _, subset := range endpoints.Subsets {
		for _, addr := range subset.Addresses {
			hostAliasList = append(hostAliasList, build.HostAlias{IP: addr.IP, Hostnames: hostNames})
		}
	}
	return
}

//IsDockerfile CheckDockerfile
func (i *SourceCodeBuildItem) IsDockerfile() bool {
	filepath := path.Join(i.RepoInfo.GetCodeBuildAbsPath(), "Dockerfile")
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return true
}

func (i *SourceCodeBuildItem) prepare() error {
	if err := util.CheckAndCreateDir(i.CacheDir); err != nil {
		return err
	}
	if err := util.CheckAndCreateDir(i.TGZDir); err != nil {
		return err
	}
	if err := util.CheckAndCreateDir(i.RepoInfo.GetCodeHome()); err != nil {
		return err
	}
	if _, ok := i.BuildEnvs["NO_CACHE"]; ok {
		if !util.DirIsEmpty(i.RepoInfo.GetCodeHome()) {
			os.RemoveAll(i.RepoInfo.GetCodeHome())
		}
		if err := os.RemoveAll(i.CacheDir); err != nil {
			logrus.Error("remove cache dir error", err.Error())
		}
		if err := os.MkdirAll(i.CacheDir, 0755); err != nil {
			logrus.Error("make cache dir error", err.Error())
		}
	}
	os.Chown(i.CacheDir, 200, 200)
	os.Chown(i.TGZDir, 200, 200)
	return nil
}

//UpdateVersionInfo Update build application service version info
func (i *SourceCodeBuildItem) UpdateVersionInfo(vi *dbmodel.VersionInfo) error {
	version, err := db.GetManager().VersionInfoDao().GetVersionByDeployVersion(i.DeployVersion, i.ServiceID)
	if err != nil {
		return err
	}
	if vi.DeliveredType != "" {
		version.DeliveredType = vi.DeliveredType
	}
	if vi.DeliveredPath != "" {
		version.DeliveredPath = vi.DeliveredPath
		if vi.DeliveredType == "image" {
			version.ImageName = vi.DeliveredPath
		}
	}
	if vi.FinalStatus != "" {
		version.FinalStatus = vi.FinalStatus
	}
	version.CommitMsg = vi.CommitMsg
	version.Author = vi.Author
	version.CodeVersion = vi.CodeVersion
	version.CodeBranch = vi.CodeBranch
	version.FinishTime = time.Now()
	if err := db.GetManager().VersionInfoDao().UpdateModel(version); err != nil {
		return err
	}
	return nil
}

//UpdateBuildVersionInfo update service build version info to db
func (i *SourceCodeBuildItem) UpdateBuildVersionInfo(res *build.Response) error {
	vi := &dbmodel.VersionInfo{
		DeliveredType: string(res.MediumType),
		DeliveredPath: res.MediumPath,
		EventID:       i.EventID,
		ImageName:     builder.RUNNERIMAGENAME,
		FinalStatus:   "success",
		CodeBranch:    i.CodeSouceInfo.Branch,
		CodeVersion:   i.commit.Hash,
		CommitMsg:     i.commit.Message,
		Author:        i.commit.Author,
		FinishTime:    time.Now(),
	}
	if err := i.UpdateVersionInfo(vi); err != nil {
		logrus.Errorf("update version info error: %s", err.Error())
		i.Logger.Error("Update application service version information failed", map[string]string{"step": "build-code", "status": "failure"})
		return err
	}
	return nil
}

//UpdateCheckResult UpdateCheckResult
func (i *SourceCodeBuildItem) UpdateCheckResult(result *dbmodel.CodeCheckResult) error {
	return nil
}
