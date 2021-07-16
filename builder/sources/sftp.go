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
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/pkg/sftp"
)

//SFTPClient
type SFTPClient struct {
	UserName   string `json:"username"`
	PassWord   string `json:"password"`
	Host       string `json:"host"`
	Port       int    `json:"int"`
	conn       *ssh.Client
	sftpClient *sftp.Client
}

//NewSFTPClient
func NewSFTPClient(username, password, host, port string) (*SFTPClient, error) {
	fb := &SFTPClient{
		UserName: username,
		PassWord: password,
		Host:     host,
	}
	if len(port) != 0 {
		var err error
		fb.Port, err = strconv.Atoi(port)
		if err != nil {
			fb.Port = 21
		}
	} else {
		fb.Port = 21
	}
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}
	if fb.PassWord != "" {
		auths = append(auths, ssh.Password(fb.PassWord))
	}
	config := ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", host, fb.Port)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		logrus.Errorf("unable to connect to [%s]: %v", addr, err)
		return nil, err
	}
	c, err := sftp.NewClient(conn, sftp.MaxPacket(1<<15))
	if err != nil {
		logrus.Errorf("unable to start sftp subsytem: %v", err)
		return nil, err
	}
	fb.conn = conn
	fb.sftpClient = c
	return fb, nil
}

//Close
func (s *SFTPClient) Close() {
	if s.sftpClient != nil {
		s.sftpClient.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
func (s *SFTPClient) checkMd5(src, dst string, logger event.Logger) (bool, error) {
	if err := util.CreateFileHash(src, src+".md5"); err != nil {
		return false, err
	}
	existmd5, err := s.FileExist(dst + ".md5")
	if err != nil && err.Error() != "file does not exist" {
		return false, err
	}
	exist, err := s.FileExist(dst)
	if err != nil && err.Error() != "file does not exist" {
		return false, err
	}
	if exist && existmd5 {
		if err := s.DownloadFile(dst+".md5", src+".md5.old", logger); err != nil {
			return false, err
		}
		old, err := ioutil.ReadFile(src + ".md5.old")
		if err != nil {
			return false, err
		}
		os.Remove(src + ".md5.old")
		new, err := ioutil.ReadFile(src + ".md5")
		if err != nil {
			return false, err
		}
		if string(old) == string(new) {
			return true, nil
		}
	}
	return false, nil
}

//PushFile
func (s *SFTPClient) PushFile(src, dst string, logger event.Logger) error {
	logger.Info(fmt.Sprintf("Start uploading code package to FTP server"), map[string]string{"step": "slug-share"})
	ok, err := s.checkMd5(src, dst, logger)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to open source file", map[string]string{"step": "share"})
		}
		return err
	}
	defer srcFile.Close()
	srcStat, err := srcFile.Stat()
	if err != nil {
		if logger != nil {
			logger.Error("Failed to open source file", map[string]string{"step": "share"})
		}
		return err
	}
	// check or create dir
	dir := filepath.Dir(dst)
	_, err = s.sftpClient.Stat(dir)
	if err != nil {
		if err.Error() == "file does not exist" {
			err := s.MkdirAll(dir)
			if err != nil {
				if logger != nil {
					logger.Error("Failed to create target file directory", map[string]string{"step": "share"})
				}
				return err
			}
		} else {
			if logger != nil {
				logger.Error("Failed to detect the target file directory", map[string]string{"step": "share"})
			}
			return err
		}
	}
	// remove all files if they exist
	s.sftpClient.Remove(dst)
	dstFile, err := s.sftpClient.Create(dst)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to open target file", map[string]string{"step": "share"})
		}
		return err
	}
	defer dstFile.Close()
	allSize := srcStat.Size()
	if err := CopyWithProgress(srcFile, dstFile, allSize, logger); err != nil {
		return err
	}
	// write remote md5 file
	md5, _ := ioutil.ReadFile(src + ".md5")
	dstMd5File, err := s.sftpClient.Create(dst + ".md5")
	if err != nil {
		logrus.Errorf("create md5 file in sftp server error.%s", err.Error())
		return nil
	}
	defer dstMd5File.Close()
	if _, err := dstMd5File.Write(md5); err != nil {
		logrus.Errorf("write md5 file in sftp server error.%s", err.Error())
	}
	return nil
}

//DownloadFile
func (s *SFTPClient) DownloadFile(src, dst string, logger event.Logger) error {
	logger.Info(fmt.Sprintf("Start downloading code package from FTP server"), map[string]string{"step": "slug-share"})

	srcFile, err := s.sftpClient.OpenFile(src, 0644)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to open source file", map[string]string{"step": "share"})
		}
		return err
	}
	defer srcFile.Close()
	srcStat, err := srcFile.Stat()
	if err != nil {
		if logger != nil {
			logger.Error("Failed to open source file", map[string]string{"step": "share"})
		}
		return err
	}
	// Verify and create target directory
	dir := filepath.Dir(dst)
	if err := util.CheckAndCreateDir(dir); err != nil {
		if logger != nil {
			logger.Error("Failed to detect and create the target file directory", map[string]string{"step": "share"})
		}
		return err
	}
	// Delete the file first if it exists
	os.Remove(dst)
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to open target file", map[string]string{"step": "share"})
		}
		return err
	}
	defer dstFile.Close()
	allSize := srcStat.Size()
	return CopyWithProgress(srcFile, dstFile, allSize, logger)
}

//FileExist
func (s *SFTPClient) FileExist(filepath string) (bool, error) {
	if _, err := s.sftpClient.Stat(filepath); err != nil {
		return false, err
	}
	return true, nil
}

//MkdirAll
func (s *SFTPClient) MkdirAll(dirpath string) error {
	parentDir := filepath.Dir(dirpath)
	_, err := s.sftpClient.Stat(parentDir)
	if err != nil {
		if err.Error() == "file does not exist" {
			err := s.MkdirAll(parentDir)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err = s.sftpClient.Mkdir(dirpath)
	if err != nil {
		return err
	}
	return nil
}
