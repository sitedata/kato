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

package parser

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/builder"
	"github.com/gridworkz/kato/builder/parser/code"
	multi "github.com/gridworkz/kato/builder/parser/code/multisvc"
	"github.com/gridworkz/kato/builder/parser/types"
	"github.com/gridworkz/kato/builder/sources"
	"github.com/gridworkz/kato/db/model"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport" //"github.com/docker/docker/client"
)

//SourceCodeParse docker run - command resolution or direct image name resolution
type SourceCodeParse struct {
	ports        map[int]*types.Port
	volumes      map[string]*types.Volume
	envs         map[string]*types.Env
	source       string
	memory       int
	image        Image
	args         []string
	branchs      []string
	errors       []ParseError
	dockerclient *client.Client
	logger       event.Logger
	Lang         code.Lang

	Runtime      bool `json:"runtime"`
	Dependencies bool `json:"dependencies"`
	Procfile     bool `json:"procfile"`

	isMulti  bool
	services []*types.Service
}

//CreateSourceCodeParse
func CreateSourceCodeParse(source string, logger event.Logger) Parser {
	return &SourceCodeParse{
		source:  source,
		ports:   make(map[int]*types.Port),
		volumes: make(map[string]*types.Volume),
		envs:    make(map[string]*types.Env),
		logger:  logger,
		image:   ParseImageName(builder.RUNNERIMAGENAME),
		args:    []string{"start", "web"},
	}
}

//Parse - get the code Analyze the code Verify the code
func (d *SourceCodeParse) Parse() ParseErrorList {
	if d.source == "" {
		d.logger.Error("Source code check input parameter error", map[string]string{"step": "parse"})
		d.errappend(Errorf(FatalError, "source can not be empty"))
		return d.errors
	}
	var csi sources.CodeSourceInfo
	err := ffjson.Unmarshal([]byte(d.source), &csi)
	if err != nil {
		d.logger.Error("Source code check input parameter error", map[string]string{"step": "parse"})
		d.errappend(Errorf(FatalError, "source data can not be read"))
		return d.errors
	}
	if csi.Branch == "" {
		csi.Branch = "master"
	}
	if csi.RepositoryURL == "" {
		d.logger.Error("Git project warehouse address cannot be empty", map[string]string{"step": "parse"})
		d.errappend(ErrorAndSolve(FatalError, "Git project warehouse address format error", SolveAdvice("modify_url", "Please confirm and modify the warehouse address")))
		return d.errors
	}
	//Verify warehouse address
	buildInfo, err := sources.CreateRepostoryBuildInfo(csi.RepositoryURL, csi.ServerType, csi.Branch, csi.TenantID, csi.ServiceID)
	if err != nil {
		d.logger.Error("Git project warehouse address format error", map[string]string{"step": "parse"})
		d.errappend(ErrorAndSolve(FatalError, "Git project warehouse address format error", SolveAdvice("modify_url", "Please confirm and modify the warehouse address")))
		return d.errors
	}
	gitFunc := func() ParseErrorList {
		//get code
		if !util.DirIsEmpty(buildInfo.GetCodeHome()) {
			if err := sources.RemoveDir(buildInfo.GetCodeHome()); err != nil {
				logrus.Errorf("remove code dir failure %s", err.Error())
				return d.errors
			}
		}
		csi.RepositoryURL = buildInfo.RepostoryURL
		rs, err := sources.GitClone(csi, buildInfo.GetCodeHome(), d.logger, 5)
		if err != nil {
			if err == transport.ErrAuthenticationRequired || err == transport.ErrAuthorizationFailed {
				if buildInfo.GetProtocol() == "ssh" {
					d.errappend(ErrorAndSolve(FatalError, "Git project warehouse requires security verification", SolveAdvice("get_publickey", "Please obtain the authorization key to configure it in your warehouse project")))
				} else {
					d.errappend(ErrorAndSolve(FatalError, "Git project warehouse requires security verification", SolveAdvice("modify_userpass", "Please provide the correct account password")))
				}
				return d.errors
			}
			if err == plumbing.ErrReferenceNotFound {
				solve := "Please go to the code warehouse to check the correct branch situation"
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Git project warehouse specified branch %s does not exist", csi.Branch), solve))
				return d.errors
			}
			if err == transport.ErrRepositoryNotFound {
				solve := SolveAdvice("modify_repo", "Please confirm whether the warehouse address is correct")
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Git project repository does not exist"), solve))
				return d.errors
			}
			if err == transport.ErrEmptyRemoteRepository {
				solve := SolveAdvice("open_repo", "Please confirm that the code has been submitted")
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("No valid file in Git project warehouse"), solve))
				return d.errors
			}
			if strings.Contains(err.Error(), "ssh: unable to authenticate") {
				solve := SolveAdvice("get_publickey", "Please get the authorized key configuration to your warehouse project to try")
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Remote warehouse SSH authentication error"), solve))
				return d.errors
			}
			if strings.Contains(err.Error(), "context deadline exceeded") {
				solve := "Please confirm whether the source code repository can be accessed normally"
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Get code timed out"), solve))
				return d.errors
			}
			logrus.Errorf("git clone error,%s", err.Error())
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Failed to get code"), "Please confirm whether the warehouse can be accessed normally, or contact customer service for consultation"))
			return d.errors
		}
		//Get branch
		branch, err := rs.Branches()
		if err == nil {
			branch.ForEach(func(re *plumbing.Reference) error {
				name := re.Name()
				if name.IsBranch() {
					d.branchs = append(d.branchs, name.Short())
				}
				return nil
			})
		} else {
			d.branchs = append(d.branchs, csi.Branch)
		}
		return nil
	}

	svnFunc := func() ParseErrorList {
		if sources.CheckFileExist(buildInfo.GetCodeHome()) {
			if err := sources.RemoveDir(buildInfo.GetCodeHome()); err != nil {
				//d.errappend(ErrorAndSolve(err, "Clean up the cache dir error", "please submit the code to the warehouse"))
				return d.errors
			}
		}
		csi.RepositoryURL = buildInfo.RepostoryURL
		svnclient := sources.NewClient(csi, buildInfo.GetCodeHome(), d.logger)
		rs, err := svnclient.UpdateOrCheckout(buildInfo.BuildPath)
		if err != nil {
			if strings.Contains(err.Error(), "svn:E170000") {
				solve := "Please go to the code warehouse to check the correct branch situation"
				d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("svn project warehouse designated branch %s does not exist", csi.Branch), solve))
				return d.errors
			}
			logrus.Errorf("svn checkout or update error,%s", err.Error())
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("Failed to get code"), "Please confirm whether the warehouse can be accessed normally, or check the community documentation"))
			return d.errors
		}
		//get branches
		d.branchs = rs.Branchs
		return nil
	}
	logrus.Debugf("start get service code by %s server type", csi.ServerType)
	
	//Get the code repository
	switch csi.ServerType {
	case "git":
		if err := gitFunc(); err != nil && err.IsFatalError() {
			return err
		}
	case "svn":
		if err := svnFunc(); err != nil && err.IsFatalError() {
			return err
		}
	default:
		//default git
		logrus.Warningf("do not get void server type,default use git")
		if err := gitFunc(); err != nil && err.IsFatalError() {
			return err
		}
	}
	// The source code is useless after the test is completed, and needs to be deleted.
	defer func() {
		if sources.CheckFileExist(buildInfo.GetCodeHome()) {
			if err := sources.RemoveDir(buildInfo.GetCodeHome()); err != nil {
				logrus.Warningf("remove source code: %v", err)
			}
		}
	}()

	//read katofile
	rbdfileConfig, err := code.ReadKatoFile(buildInfo.GetCodeBuildAbsPath())
	if err != nil {
		if err != code.ErrKatoFileNotFound {
			d.errappend(ErrorAndSolve(NegligibleError, "The katofile definition format is incorrect", "You can refer to the document description to configure this file to define application attributes"))
		}
	}
	//Judgment target directory
	var buildPath = buildInfo.GetCodeBuildAbsPath()
	//Parse code type
	var lang code.Lang
	if rbdfileConfig != nil && rbdfileConfig.Language != "" {
		lang = code.Lang(rbdfileConfig.Language)
	} else {
		lang, err = code.GetLangType(buildPath)
		if err != nil {
			if err == code.ErrCodeDirNotExist {
				d.errappend(ErrorAndSolve(FatalError, "The source code directory does not exist", "Failed to get the code task, please contact customer service"))
			} else if err == code.ErrCodeNotExist {
				d.errappend(ErrorAndSolve(FatalError, "The code in the warehouse does not exist", "Please submit the code to the warehouse"))
			} else {
				d.errappend(ErrorAndSolve(FatalError, "The code cannot recognize the language type", "Please refer to the document to view the platform language support specification"))
			}
			return d.errors
		}
	}
	d.Lang = lang
	if lang == code.NO {
		d.errappend(ErrorAndSolve(FatalError, "The code cannot recognize the language type", "Please refer to the document to view the platform language support specification"))
		return d.errors
	}
	//check code Specification
	spec := code.CheckCodeSpecification(buildPath, lang)
	if spec.Advice != nil {
		for k, v := range spec.Advice {
			d.errappend(ErrorAndSolve(NegligibleError, k, v))
		}
	}
	if spec.Noconform != nil {
		for k, v := range spec.Noconform {
			d.errappend(ErrorAndSolve(FatalError, k, v))
		}
	}
	if !spec.Conform {
		return d.errors
	}
	//If it is a dockerfile, parse the dockerfile file
	if lang == code.Dockerfile {
		if ok := d.parseDockerfileInfo(path.Join(buildPath, "Dockerfile")); !ok {
			return d.errors
		}
	}
	runtimeInfo, err := code.CheckRuntime(buildPath, lang)
	if err != nil && err == code.ErrRuntimeNotSupport {
		d.errappend(ErrorAndSolve(FatalError, "The runtime version of the code selection is not supported", "Please refer to the document to view the runtime version supported by each language of the platform"))
		return d.errors
	}
	for k, v := range runtimeInfo {
		d.envs["BUILD_"+k] = &types.Env{
			Name:  "BUILD_" + k,
			Value: v,
		}
	}
	d.memory = getRecommendedMemory(lang)
	var ProcfileLine string
	d.Procfile, ProcfileLine = code.CheckProcfile(buildPath, lang)
	if ProcfileLine != "" {
		d.envs["BUILD_PROCFILE"] = &types.Env{
			Name:  "BUILD_PROCFILE",
			Value: ProcfileLine,
		}
	}

	// multi services
	m := multi.NewMultiServiceI(lang.String())
	if m != nil {
		logrus.Infof("Lang: %s; start listing multi modules", lang.String())
		services, err := m.ListModules(buildInfo.GetCodeBuildAbsPath())
		if err != nil {
			d.logger.Error("Failed to parse multi-module project", map[string]string{"step": "parse"})
			d.errappend(ErrorAndSolve(FatalError, fmt.Sprintf("error listing modules: %v", err), "check source code for multi-modules"))
			return d.errors
		}
		if services != nil && len(services) > 1 {
			d.isMulti = true
			d.services = services
		}

		if rbdfileConfig != nil && rbdfileConfig.Services != nil && len(rbdfileConfig.Services) > 0 {
			mm := make(map[string]*types.Service)
			for i := range services {
				mm[services[i].Name] = services[i]
			}
			for _, svc := range rbdfileConfig.Services {
				if item := mm[svc.Name]; item != nil {
					for k, v := range svc.Envs {
						if item.Envs == nil {
							item.Envs = make(map[string]*types.Env, len(rbdfileConfig.Envs))
						}
						item.Envs[k] = &types.Env{Name: k, Value: v}
					}
					for i := range svc.Ports {
						if item.Ports == nil {
							item.Ports = make(map[int]*types.Port, len(rbdfileConfig.Ports))
						}
						item.Ports[i] = &types.Port{
							ContainerPort: svc.Ports[i].Port, Protocol: svc.Ports[i].Protocol,
						}
					}
					for k, v := range rbdfileConfig.Envs {
						if item.Envs == nil {
							item.Envs = make(map[string]*types.Env, len(rbdfileConfig.Envs))
						}
						if item.Envs[k] == nil {
							item.Envs[k] = &types.Env{Name: k, Value: fmt.Sprintf("%v", v)}
						}
					}
					for _, port := range rbdfileConfig.Ports {
						if item.Ports == nil {
							item.Ports = make(map[int]*types.Port, len(rbdfileConfig.Ports))
						}
						if item.Ports[port.Port] == nil {
							item.Ports[port.Port] = &types.Port{
								ContainerPort: port.Port,
								Protocol:      port.Protocol,
							}
						}
					}
				}
			}
		}

		if rbdfileConfig != nil && d.isMulti {
			rbdfileConfig.Envs = nil
			rbdfileConfig.Ports = nil
		}
	}

	if rbdfileConfig != nil {
		//handle profile env
		for k, v := range rbdfileConfig.Envs {
			d.envs[k] = &types.Env{Name: k, Value: fmt.Sprintf("%v", v)}
		}
		//handle profile port
		for _, port := range rbdfileConfig.Ports {
			if port.Port == 0 {
				continue
			}
			if port.Protocol == "" {
				port.Protocol = GetPortProtocol(port.Port)
			}
			d.ports[port.Port] = &types.Port{ContainerPort: port.Port, Protocol: port.Protocol}
		}
		if rbdfileConfig.Cmd != "" {
			d.args = strings.Split(rbdfileConfig.Cmd, " ")
		}
	}
	return d.errors
}

//ReadRbdConfigAndLang read katofile and lang
func ReadRbdConfigAndLang(buildInfo *sources.RepostoryBuildInfo) (*code.KatoFileConfig, code.Lang, error) {
	rbdfileConfig, err := code.ReadKatoFile(buildInfo.GetCodeBuildAbsPath())
	if err != nil {
		return nil, code.NO, err
	}
	var lang code.Lang
	if rbdfileConfig != nil && rbdfileConfig.Language != "" {
		lang = code.Lang(rbdfileConfig.Language)
	} else {
		lang, err = code.GetLangType(buildInfo.GetCodeBuildAbsPath())
		if err != nil {
			return rbdfileConfig, code.NO, err
		}
	}
	return rbdfileConfig, lang, nil
}

func getRecommendedMemory(lang code.Lang) int {
	//java recommended 1024
	if lang == code.JavaJar || lang == code.JavaMaven || lang == code.JaveWar || lang == code.Gradle {
		return 1024
	}
	if lang == code.Python {
		return 512
	}
	if lang == code.Nodejs {
		return 512
	}
	if lang == code.PHP {
		return 512
	}
	return 512
}

func (d *SourceCodeParse) errappend(pe ParseError) {
	d.errors = append(d.errors, pe)
}

//GetBranches
func (d *SourceCodeParse) GetBranchs() []string {
	return d.branchs
}

//GetPorts
func (d *SourceCodeParse) GetPorts() (ports []types.Port) {
	for _, cv := range d.ports {
		ports = append(ports, *cv)
	}
	return ports
}

//GetVolumes
func (d *SourceCodeParse) GetVolumes() (volumes []types.Volume) {
	for _, cv := range d.volumes {
		volumes = append(volumes, *cv)
	}
	return
}

//GetValid
func (d *SourceCodeParse) GetValid() bool {
	return false
}

//GetEnvs
func (d *SourceCodeParse) GetEnvs() (envs []types.Env) {
	for _, cv := range d.envs {
		envs = append(envs, *cv)
	}
	return
}

//GetImage
func (d *SourceCodeParse) GetImage() Image {
	return d.image
}

//GetArgs
func (d *SourceCodeParse) GetArgs() []string {
	return d.args
}

//GetMemory
func (d *SourceCodeParse) GetMemory() int {
	return d.memory
}

//GetLang
func (d *SourceCodeParse) GetLang() code.Lang {
	return d.Lang
}

//GetServiceInfo
func (d *SourceCodeParse) GetServiceInfo() []ServiceInfo {
	serviceInfo := ServiceInfo{
		Ports:       d.GetPorts(),
		Envs:        d.GetEnvs(),
		Volumes:     d.GetVolumes(),
		Image:       d.GetImage(),
		Args:        d.GetArgs(),
		Branchs:     d.GetBranchs(),
		Memory:      d.memory,
		Lang:        d.GetLang(),
		ServiceType: model.ServiceTypeStatelessMultiple.String(),
		OS:          runtime.GOOS,
	}
	var res []ServiceInfo
	if d.isMulti && d.services != nil && len(d.services) > 0 {
		for idx := range d.services {
			svc := d.services[idx]
			info := serviceInfo
			info.ID = util.NewUUID()
			info.Name = svc.Name
			info.Cname = svc.Cname
			info.Packaging = svc.Packaging
			for i := range svc.Envs {
				info.Envs = append(info.Envs, *svc.Envs[i])
			}
			for i := range svc.Ports {
				info.Ports = append(info.Ports, *svc.Ports[i])
			}
			res = append(res, info)
		}
	} else {
		serviceInfo.Envs = d.GetEnvs()
		serviceInfo.Ports = d.GetPorts()
		res = []ServiceInfo{serviceInfo}
	}

	return res
}

func removeQuotes(value string) string {
	if len(value) > 0 && (value[0] == '"' || value[0] == '\'') {
		value = value[1:]
	}
	if len(value) > 0 && (value[len(value)-1] == '"' || value[0] == '\'') {
		value = value[:len(value)-1]
	}
	return value
}

func (d *SourceCodeParse) parseDockerfileInfo(dockerfile string) bool {
	commands, err := sources.ParseFile(dockerfile)
	if err != nil {
		d.errappend(ErrorAndSolve(FatalError, err.Error(), "Please confirm whether the Dockerfile format complies with the specification"))
		return false
	}

	for _, cm := range commands {
		switch cm.Cmd {
		case "arg":
			length := len(cm.Value)
			for i := 0; i < length; i++ {
				if kv := strings.Split(cm.Value[i], "="); len(kv) > 1 {
					key := "BUILD_ARG_" + kv[0]
					d.envs[key] = &types.Env{Name: key, Value: removeQuotes(kv[1])}
				} else {
					if i+1 >= length {
						logrus.Error("Parse ARG format error at ", cm.Value[i])
						continue
					}
					key := "BUILD_ARG_" + cm.Value[i]
					d.envs[key] = &types.Env{Name: key, Value: removeQuotes(cm.Value[i+1])}
					i++
				}
			}
		case "env":
			length := len(cm.Value)
			for i := 0; i < len(cm.Value); i++ {
				if kv := strings.Split(cm.Value[i], "="); len(kv) > 1 {
					d.envs[kv[0]] = &types.Env{Name: kv[0], Value: kv[1]}
				} else {
					if i+1 >= length {
						logrus.Error("Parse ENV format error at ", cm.Value[1])
						continue
					}
					d.envs[cm.Value[i]] = &types.Env{Name: cm.Value[i], Value: cm.Value[i+1]}
					i++
				}
			}
		case "expose":
			for _, v := range cm.Value {
				port, _ := strconv.Atoi(v)
				if port != 0 {
					d.ports[port] = &types.Port{ContainerPort: port, Protocol: GetPortProtocol(port)}
				}
			}
		case "volume":
			for _, v := range cm.Value {
				d.volumes[v] = &types.Volume{VolumePath: v, VolumeType: model.ShareFileVolumeType.String()}
			}
		}
	}
	// dockerfile empty args
	d.args = []string{}
	return true
}
