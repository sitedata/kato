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

package sources

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
)

//SVNClient
type SVNClient interface {
	Checkout() (*Info, error)
	Update(childpath string) (*Info, error)
	UpdateOrCheckout(childpath string) (*Info, error)
}

type svnclient struct {
	username string
	password string
	svnURL   string
	svnDir   string
	Env      []string
	logger   event.Logger
	csi      CodeSourceInfo
}

// NewClient
func NewClient(csi CodeSourceInfo, codeHome string, logger event.Logger) SVNClient {
	return &svnclient{csi: csi, username: csi.User, password: csi.Password, svnURL: csi.RepositoryURL, svnDir: codeHome, logger: logger}
}

// NewClientWithEnv ...
func NewClientWithEnv(username, password, url, sourceDir string, env []string) SVNClient {
	util.CheckAndCreateDir(sourceDir)
	return &svnclient{username: username, password: password, svnURL: url, Env: env, svnDir: sourceDir}
}

// Diff ...
func (c *svnclient) Diff(start, end int) (string, error) {
	r := fmt.Sprintf("%d:%d", start, end)
	out, err := c.run("diff", "-r", r, c.svnURL)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Cat ...
func (c *svnclient) Cat(file string) (string, error) {
	out, err := c.run("cat", c.svnURL+"/"+file)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Export ...
func (c *svnclient) Export(dir string) error {
	_, err := c.run("export", c.svnURL, dir)
	if err != nil {
		return err
	}
	return nil

}

// Log ...
func (c *svnclient) Log() (*Logs, error) {
	out, err := c.run("log", "--xml", "-v")
	if err != nil {
		return nil, err
	}
	l := new(Logs)
	err = xml.Unmarshal(out, l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// list ...
func (c *svnclient) List() (*lists, error) {
	cmd := []string{"list", c.svnURL, "--xml"}
	out, err := c.run(cmd...)
	if err != nil {
		return nil, err
	}
	l := new(lists)
	err = xml.Unmarshal(out, l)
	if err != nil {
		return nil, err
	}
	return l, nil

}
func getBranchPath(branch, url string) string {
	if strings.HasPrefix(branch, "trunk") {
		return fmt.Sprintf("%s/%s", url, branch)
	}
	if strings.HasPrefix(branch, "tag:") {
		return fmt.Sprintf("%s/tags/%s", url, branch[4:])
	}
	return fmt.Sprintf("%s/branches/%s", url, branch)
}

// Checkout
func (c *svnclient) Checkout() (*Info, error) {
	tempURL := c.svnURL
	//handle branch or tags
	if c.csi.Branch != "" {
		c.svnURL = getBranchPath(c.csi.Branch, c.svnURL)
	}
	if !util.DirIsEmpty(c.svnDir) {
		os.RemoveAll(c.svnDir)
	}
	if err := os.MkdirAll(c.svnDir, 0755); err != nil {
		return nil, err
	}
	cmd := []string{"checkout", c.svnURL}
	if c.svnDir != "" {
		cmd = append(cmd, c.svnDir)
	}
	_, err := c.runWithLogger(cmd...)
	if err != nil {
		//if trunk will change url retry
		if strings.Contains(err.Error(), "svn:E170000") && c.csi.Branch == "trunk" {
			c.svnURL = c.svnURL[:len(c.svnURL)-6]
			cmd := []string{"checkout", c.svnURL}
			if c.svnDir != "" {
				cmd = append(cmd, c.svnDir)
			}
			_, err = c.runWithLogger(cmd...)
			if err != nil {
				return nil, err
			}
		} else if strings.Contains(err.Error(), "svn:E170000") && (c.csi.Branch != "trunk" && !strings.HasPrefix(c.csi.Branch, "tag:")) {
			c.svnURL = tempURL + "/" + c.csi.Branch
			cmd := []string{"checkout", c.svnURL}
			if c.svnDir != "" {
				cmd = append(cmd, c.svnDir)
			}
			_, err = c.runWithLogger(cmd...)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return c.Info()
}

// Update
func (c *svnclient) Update(childpath string) (*Info, error) {
	cmd := []string{"update", path.Join(c.svnDir, childpath)}
	_, err := c.runWithLogger(cmd...)
	if err != nil {
		return nil, err
	}
	return c.Info()
}

func (c *svnclient) UpdateOrCheckout(childpath string) (*Info, error) {
	var rs *Info
	var err error
	if ok := util.DirIsEmpty(c.svnDir); !ok {
		//update BuildPath
		rs, err = c.Update(childpath)
		if err != nil {
			logrus.Errorf("update svn code error: %s", err.Error())
			c.logger.Error(fmt.Sprintf("Update svn code failed, please make sure the code can be downloaded properly"), map[string]string{"step": "builder-exector", "status": "failure"})
		} else {
			return rs, nil
		}
	}
	rs, err = c.Checkout()
	if err != nil {
		logrus.Errorf("checkout svn code error: %s", err.Error())
		c.logger.Error(fmt.Sprintf("Checkout svn code failed, please make sure the code can be downloaded properly"), map[string]string{"step": "builder-exector", "status": "failure"})
		return nil, err
	}
	return rs, nil
}

// svnclient Checkout from specific revision
func (c *svnclient) CheckoutWithRevision(revision string) (string, error) {
	cmd := []string{"checkout", c.svnURL}
	if c.svnDir != "" {
		cmd = append(cmd, c.svnDir)
	}
	cmd = append(cmd, "-r", revision)
	out, err := c.run(cmd...)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// info ...
func (c *svnclient) Info() (*Info, error) {
	start := time.Now()
	cmd := []string{"info", "--xml"}
	out, err := c.run(cmd...)
	if err != nil {
		return nil, err
	}
	info := new(Info)
	err = xml.Unmarshal(out, info)
	if err != nil {
		return nil, err
	}
	log, err := c.Log()
	if err == nil {
		info.Logs = log
	}
	fmt.Println(time.Now().Sub(start).String())
	info.Branchs = c.readBranchs()
	info.Tags = c.readTags()
	return info, nil
}
func (c *svnclient) readBranchs() []string {
	re := []string{""}
	exist, _ := util.FileExists(path.Join(c.svnDir, "Branches"))
	if exist {
		list, err := ioutil.ReadDir(path.Join(c.svnDir, "Branches"))
		if err != nil {
			return re
		}
		for _, f := range list {
			re = append(re, f.Name())
		}
		return re
	}
	exist, _ = util.FileExists(path.Join(c.svnDir, "branches"))
	if exist {
		list, err := ioutil.ReadDir(path.Join(c.svnDir, "Branches"))
		if err != nil {
			return re
		}
		for _, f := range list {
			re = append(re, f.Name())
		}
		return re
	}
	return re
}

func (c *svnclient) readTags() []string {
	re := []string{}
	exist, _ := util.FileExists(path.Join(c.svnDir, "Tags"))
	if exist {
		list, err := ioutil.ReadDir(path.Join(c.svnDir, "Tags"))
		if err != nil {
			return re
		}
		for _, f := range list {
			re = append(re, f.Name())
		}
		return re
	}
	exist, _ = util.FileExists(path.Join(c.svnDir, "tags"))
	if exist {
		list, err := ioutil.ReadDir(path.Join(c.svnDir, "tags"))
		if err != nil {
			return re
		}
		for _, f := range list {
			re = append(re, f.Name())
		}
		return re
	}
	return re
}

// run
func (c *svnclient) runWithLogger(args ...string) ([]byte, error) {
	ops := []string{"--username", c.username, "--password", c.password, "--non-interactive", "--trust-server-cert-failures", "unknown-ca,cn-mismatch,expired,not-yet-valid,other"}
	args = append(args, ops...)
	cmd := exec.Command("svn", args...)
	if len(c.Env) > 0 {
		cmd.Env = append(os.Environ(), c.Env...)
	}
	cmd.Dir = c.svnDir
	writer := c.logger.GetWriter("progress", "debug")
	writer.SetFormat(map[string]interface{}{"progress": "%s", "id": "SVN:"})
	cmd.Stdout = writer
	errorWriter := bytes.NewBuffer(nil)
	cmd.Stderr = errorWriter
	err := cmd.Run()
	if err != nil {
		if strings.Contains(errorWriter.String(), "doesn't exist") {
			return nil, fmt.Errorf("svn:E170000")
		}
		return nil, fmt.Errorf("svn error:%s", errorWriter.String())
	}
	return nil, nil
}

// run
func (c *svnclient) run(args ...string) ([]byte, error) {
	ops := []string{"--username", c.username, "--password", c.password, "--non-interactive", "--trust-server-cert-failures", "unknown-ca,cn-mismatch,expired,not-yet-valid,other"}
	args = append(args, ops...)
	cmd := exec.Command("svn", args...)
	if len(c.Env) > 0 {
		cmd.Env = append(os.Environ(), c.Env...)
	}
	cmd.Dir = c.svnDir
	return cmd.Output()
}

//Logs commit logs
type Logs struct {
	XMLName      xml.Name `xml:"log"`
	CommitEntrys []Commit `xml:"logentry"`
}

type paths struct {
	Path     string `xml:",innerxml"`
	Action   string `xml:"action,attr"`
	PropMods string `xml:"prop-mods,attr"`
	TextMods string `xml:"text-mods,attr"`
	Kind     string `xml:"kind,attr"`
}

//Info
type Info struct {
	XMLName       xml.Name `xml:"info"`
	URL           string   `xml:"entry>url"`
	RelativeURL   string   `xml:"entry>relative-url"`
	Root          string   `xml:"entry>repository>root"`
	UUID          string   `xml:"entry>repository>uuid"`
	WcrootAbspath string   `xml:"entry>wc-info>wcroot-abspath"`
	Schedule      string   `xml:"entry>wc-info>schedule"`
	Depth         string   `xml:"entry>wc-info>depth"`
	Logs          *Logs
	Branchs       []string
	Tags          []string
}

//Commit
type Commit struct {
	Revision string `xml:"revision,attr"`
	Author   string `xml:"author"`
	Msg      string `xml:"msg"`
	Date     string `xml:"date"`
}

type lists struct {
	XMLName xml.Name `xml:"lists"`
	List    list     `xml:"list"`
}

type list struct {
	Path  string  `xml:"path,attr"`
	Entry []entry `xml:"entry"`
}

type entry struct {
	Kind   string `xml:"kind,attr"`
	Name   string `xml:"name"`
	Size   string `xml:"size"`
	Commit Commit `xml:"commit"`
}
