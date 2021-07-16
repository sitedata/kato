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

package sources

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/sirupsen/logrus"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"

	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
	netssh "golang.org/x/crypto/ssh"
	sshkey "golang.org/x/crypto/ssh"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

//CodeSourceInfo
type CodeSourceInfo struct {
	ServerType    string `json:"server_type"`
	RepositoryURL string `json:"repository_url"`
	Branch        string `json:"branch"`
	User          string `json:"user"`
	Password      string `json:"password"`
	//To avoid conflicts between projects, the code cache directory is increased to tenants
	TenantID  string `json:"tenant_id"`
	ServiceID string `json:"service_id"`
}

//GetCodeSourceDir get source storage directory
func (c CodeSourceInfo) GetCodeSourceDir() string {
	return GetCodeSourceDir(c.RepositoryURL, c.Branch, c.TenantID, c.ServiceID)
}

//GetCodeSourceDir get source storage directory
// it changes as gitrepostory address, branch, and service id change
func GetCodeSourceDir(RepositoryURL, branch, tenantID string, ServiceID string) string {
	sourceDir := os.Getenv("SOURCE_DIR")
	if sourceDir == "" {
		sourceDir = "/grdata/source"
	}
	h := sha1.New()
	h.Write([]byte(RepositoryURL + branch + ServiceID))
	bs := h.Sum(nil)
	bsStr := fmt.Sprintf("%x", bs)
	return path.Join(sourceDir, "build", tenantID, bsStr)
}

//CheckFileExist CheckFileExist
func CheckFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//RemoveDir
func RemoveDir(path string) error {
	if path == "/" {
		return fmt.Errorf("remove wrong dir")
	}
	return os.RemoveAll(path)
}
func getShowURL(rurl string) string {
	urlpath, _ := url.Parse(rurl)
	if urlpath != nil {
		showURL := fmt.Sprintf("%s://%s%s", urlpath.Scheme, urlpath.Host, urlpath.Path)
		return showURL
	}
	return ""
}

//GitClone
func GitClone(csi CodeSourceInfo, sourceDir string, logger event.Logger, timeout int) (*git.Repository, error) {
	GetPrivateFileParam := csi.TenantID
	if !strings.HasSuffix(csi.RepositoryURL, ".git") {
		csi.RepositoryURL = csi.RepositoryURL + ".git"
	}
	flag := true
Loop:
	if logger != nil {
		//Hide possible account key information
		logger.Info(fmt.Sprintf("Start clone source code from %s", getShowURL(csi.RepositoryURL)), map[string]string{"step": "clone_code"})
	}
	ep, err := transport.NewEndpoint(csi.RepositoryURL)
	if err != nil {
		return nil, err
	}
	if timeout < 1 {
		timeout = 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()
	writer := logger.GetWriter("progress", "debug")
	writer.SetFormat(map[string]interface{}{"progress": "%s", "id": "Clone:"})
	opts := &git.CloneOptions{
		URL:               csi.RepositoryURL,
		Progress:          writer,
		SingleBranch:      true,
		Tags:              git.NoTags,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Depth:             1,
	}
	if csi.Branch != "" {
		opts.ReferenceName = getBranch(csi.Branch)
	}
	var rs *git.Repository
	if ep.Protocol == "ssh" {
		publichFile := GetPrivateFile(GetPrivateFileParam)
		sshAuth, auerr := ssh.NewPublicKeysFromFile("git", publichFile, "")
		if auerr != nil {
			if logger != nil {
				logger.Error(fmt.Sprintf("Create PublicKeys failure"), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return nil, auerr
		}
		sshAuth.HostKeyCallbackHelper.HostKeyCallback = netssh.InsecureIgnoreHostKey()
		opts.Auth = sshAuth
		rs, err = git.PlainCloneContext(ctx, sourceDir, false, opts)
	} else {
		// only proxy github
		// but when setting, other request will be proxyed
		if strings.Contains(csi.RepositoryURL, "github.com") && os.Getenv("GITHUB_PROXY") != "" {
			proxyURL, err := url.Parse(os.Getenv("GITHUB_PROXY"))
			if err == nil {
				customClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
				customClient.Timeout = time.Minute * time.Duration(timeout)
				client.InstallProtocol("https", githttp.NewClient(customClient))
				defer func() {
					client.InstallProtocol("https", githttp.DefaultClient)
				}()
			} else {
				logrus.Error(err)
			}
		}
		if csi.User != "" && csi.Password != "" {
			httpAuth := &githttp.BasicAuth{
				Username: csi.User,
				Password: csi.Password,
			}
			opts.Auth = httpAuth
		}
		rs, err = git.PlainCloneContext(ctx, sourceDir, false, opts)
	}
	if err != nil {
		if reerr := os.RemoveAll(sourceDir); reerr != nil {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred while pulling the code. Failed to delete the code directory."), map[string]string{"step": "clone-code", "status": "failure"})
			}
		}
		if err == transport.ErrAuthenticationRequired {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred while pulling the code, and the code source needs to be authorized to access."), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return rs, err
		}
		if err == transport.ErrAuthorizationFailed {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred while pulling the code, and the code source authentication failed."), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return rs, err
		}
		if err == transport.ErrRepositoryNotFound {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred in the pull code, and the warehouse does not exist."), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return rs, err
		}
		if err == transport.ErrEmptyRemoteRepository {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred in the pull code, and the remote warehouse is empty."), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return rs, err
		}
		if err == plumbing.ErrReferenceNotFound {
			if logger != nil {
				logger.Error(fmt.Sprintf("The code branch (%s) does not exist", csi.Branch), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return rs, fmt.Errorf("branch %s is not exist", csi.Branch)
		}
		if strings.Contains(err.Error(), "ssh: unable to authenticate") {

			if flag {
				GetPrivateFileParam = "builder_rsa"
				flag = false
				goto Loop
			}
			if logger != nil {
				logger.Error(fmt.Sprintf("The remote code base needs to be configured with SSH Key."), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return rs, err
		}
		if strings.Contains(err.Error(), "context deadline exceeded") {
			if logger != nil {
				logger.Error(fmt.Sprintf("Get code timed out"), map[string]string{"step": "clone-code", "status": "failure"})
			}
			return rs, err
		}
	}
	return rs, err
}
func retryAuth(ep *transport.Endpoint, csi CodeSourceInfo) (transport.AuthMethod, error) {
	switch ep.Protocol {
	case "ssh":
		home, _ := Home()
		sshAuth, err := ssh.NewPublicKeysFromFile("git", path.Join(home, "/.ssh/id_rsa"), "")
		if err != nil {
			return nil, err
		}
		return sshAuth, nil
	case "http", "https":
		//return http.NewBasicAuth(csi.User, csi.Password), nil
	}
	return nil, nil
}

//GitPull
func GitPull(csi CodeSourceInfo, sourceDir string, logger event.Logger, timeout int) (*git.Repository, error) {
	GetPrivateFileParam := csi.TenantID
	flag := true
Loop:
	if logger != nil {
		logger.Info(fmt.Sprintf("Start pull source code from %s", csi.RepositoryURL), map[string]string{"step": "clone_code"})
	}
	if timeout < 1 {
		timeout = 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()
	writer := logger.GetWriter("progress", "debug")
	writer.SetFormat(map[string]interface{}{"progress": "%s", "id": "Pull:"})
	opts := &git.PullOptions{
		Progress:     writer,
		SingleBranch: true,
		Depth:        1,
	}
	if csi.Branch != "" {
		opts.ReferenceName = getBranch(csi.Branch)
	}
	ep, err := transport.NewEndpoint(csi.RepositoryURL)
	if err != nil {
		return nil, err
	}
	if ep.Protocol == "ssh" {
		publichFile := GetPrivateFile(GetPrivateFileParam)
		sshAuth, auerr := ssh.NewPublicKeysFromFile("git", publichFile, "")
		if auerr != nil {
			if logger != nil {
				logger.Error(fmt.Sprintf("Error creating PublicKeys"), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return nil, auerr
		}
		sshAuth.HostKeyCallbackHelper.HostKeyCallback = netssh.InsecureIgnoreHostKey()
		opts.Auth = sshAuth
	} else {
		// only proxy github
		// but when setting, other request will be proxyed
		if strings.Contains(csi.RepositoryURL, "github.com") && os.Getenv("GITHUB_PROXY") != "" {
			proxyURL, _ := url.Parse(os.Getenv("GITHUB_PROXY"))
			customClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
			customClient.Timeout = time.Minute * time.Duration(timeout)
			client.InstallProtocol("https", githttp.NewClient(customClient))
			defer func() {
				client.InstallProtocol("https", githttp.DefaultClient)
			}()
		}
		if csi.User != "" && csi.Password != "" {
			httpAuth := &githttp.BasicAuth{
				Username: csi.User,
				Password: csi.Password,
			}
			opts.Auth = httpAuth
		}
	}
	rs, err := git.PlainOpen(sourceDir)
	if err != nil {
		return nil, err
	}
	tree, err := rs.Worktree()
	if err != nil {
		return nil, err
	}
	err = tree.PullContext(ctx, opts)
	if err != nil {
		if err == transport.ErrAuthenticationRequired {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred in the update code, and the code source needs to be authorized to access."), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return rs, err
		}
		if err == transport.ErrAuthorizationFailed {

			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred when updating the code, and the code source authentication failed."), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return rs, err
		}
		if err == transport.ErrRepositoryNotFound {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred in the update code and the warehouse does not exist."), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return rs, err
		}
		if err == transport.ErrEmptyRemoteRepository {
			if logger != nil {
				logger.Error(fmt.Sprintf("An error occurred in the update code and the remote warehouse is empty."), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return rs, err
		}
		if err == plumbing.ErrReferenceNotFound {
			if logger != nil {
				logger.Error(fmt.Sprintf("The code branch (%s) does not exist.", csi.Branch), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return rs, fmt.Errorf("branch %s is not exist", csi.Branch)
		}
		if strings.Contains(err.Error(), "ssh: unable to authenticate") {
			if flag {
				GetPrivateFileParam = "builder_rsa"
				flag = false
				goto Loop
			}
			if logger != nil {
				logger.Error(fmt.Sprintf("The remote code base needs to be configured with SSH Key."), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return rs, err
		}
		if strings.Contains(err.Error(), "context deadline exceeded") {
			if logger != nil {
				logger.Error(fmt.Sprintf("Update code timed out"), map[string]string{"step": "pull-code", "status": "failure"})
			}
			return rs, err
		}
		if err == git.NoErrAlreadyUpToDate {
			return rs, nil
		}
	}
	return rs, err
}

//GitCloneOrPull if code exist in local,use git pull.
func GitCloneOrPull(csi CodeSourceInfo, sourceDir string, logger event.Logger, timeout int) (*git.Repository, error) {
	if ok, err := util.FileExists(path.Join(sourceDir, ".git")); err == nil && ok && !strings.HasPrefix(csi.Branch, "tag:") {
		re, err := GitPull(csi, sourceDir, logger, timeout)
		if err == nil && re != nil {
			return re, nil
		}
		logrus.Error("git pull source code error,", err.Error())
	}
	// empty the sourceDir
	if reerr := os.RemoveAll(sourceDir); reerr != nil {
		logrus.Error("empty the source code dir error,", reerr.Error())
		if logger != nil {
			logger.Error(fmt.Sprintf("Failed to clear the code directory."), map[string]string{"step": "clone-code", "status": "failure"})
		}
	}
	return GitClone(csi, sourceDir, logger, timeout)
}

//GitCheckout checkout the specified branch
func GitCheckout(sourceDir, branch string) error {
	// option := git.CheckoutOptions{
	// 	Branch: getBranch(branch),
	// }
	return nil
}
func getBranch(branch string) plumbing.ReferenceName {
	if strings.HasPrefix(branch, "tag:") {
		return plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", branch[4:]))
	}
	return plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
}

//GetLastCommit get last commit info
//get commit by head reference
func GetLastCommit(re *git.Repository) (*object.Commit, error) {
	ref, err := re.Head()
	if err != nil {
		return nil, err
	}
	return re.CommitObject(ref.Hash())
}

//GetPrivateFile
func GetPrivateFile(tenantID string) string {
	home, _ := Home()
	if home == "" {
		home = "/root"
	}
	if ok, _ := util.FileExists(path.Join(home, "/.ssh/"+tenantID)); ok {
		return path.Join(home, "/.ssh/"+tenantID)
	} else {
		if ok, _ := util.FileExists(path.Join(home, "/.ssh/builder_rsa")); ok {
			return path.Join(home, "/.ssh/builder_rsa")
		}
		return path.Join(home, "/.ssh/id_rsa")
	}

}

//GetPublicKey
func GetPublicKey(tenantID string) string {
	home, _ := Home()
	if home == "" {
		home = "/root"
	}
	PublicKey := tenantID + ".pub"
	PrivateKey := tenantID

	if ok, _ := util.FileExists(path.Join(home, "/.ssh/"+PublicKey)); ok {
		body, _ := ioutil.ReadFile(path.Join(home, "/.ssh/"+PublicKey))
		return string(body)
	}
	Private, Public, err := MakeSSHKeyPair()
	if err != nil {
		logrus.Error("MakeSSHKeyPairError:", err)
	}
	PrivateKeyFile, err := os.Create(path.Join(home, "/.ssh/"+PrivateKey))
	if err != nil {
		fmt.Println(err)
	} else {
		PrivateKeyFile.WriteString(Private)
	}
	PublicKeyFile, err2 := os.Create(path.Join(home, "/.ssh/"+PublicKey))

	if err2 != nil {
		fmt.Println(err)
	} else {
		PublicKeyFile.WriteString(Public)
	}
	body, _ := ioutil.ReadFile(path.Join(home, "/.ssh/"+PublicKey))
	return string(body)

}

//GenerateKey
func GenerateKey(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return private, &private.PublicKey, nil

}

//EncodePrivateKey
func EncodePrivateKey(private *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(private),
		Type:  "RSA PRIVATE KEY",
	})
}

//EncodeSSHKey
func EncodeSSHKey(public *rsa.PublicKey) ([]byte, error) {
	publicKey, err := sshkey.NewPublicKey(public)
	if err != nil {
		return nil, err
	}
	return sshkey.MarshalAuthorizedKey(publicKey), nil
}

//MakeSSHKeyPair
func MakeSSHKeyPair() (string, string, error) {

	pkey, pubkey, err := GenerateKey(2048)
	if err != nil {
		return "", "", err
	}

	pub, err := EncodeSSHKey(pubkey)
	if err != nil {
		return "", "", err
	}

	return string(EncodePrivateKey(pkey)), string(pub), nil
}
