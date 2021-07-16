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

package util

import (
	"errors"
	"io"
	"net"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

//SSHClient
type SSHClient struct {
	IP             string
	Port           int
	User           string
	Password       string
	Method         string
	Key            string
	Stdout, Stderr io.Writer
	Cmd            string
}

//NewSSHClient
func NewSSHClient(ip, user, password, cmd string, port int, stdout, stderr io.Writer) *SSHClient {
	var method = "password"
	if password == "" {
		method = "publickey"
	}
	return &SSHClient{
		IP:       ip,
		User:     user,
		Password: password,
		Method:   method,
		Cmd:      cmd,
		Port:     port,
		Stderr:   stderr,
		Stdout:   stdout,
	}
}

//Connection
func (server *SSHClient) Connection() error {
	auths, err := parseAuthMethods(server)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User:            server.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := server.IP + ":" + strconv.Itoa(server.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stderr = server.Stderr
	session.Stdout = server.Stdout
	if err := session.Run(server.Cmd); err != nil {
		return err
	}
	return nil
}

// Analyze authentication method
func parseAuthMethods(server *SSHClient) ([]ssh.AuthMethod, error) {
	sshs := []ssh.AuthMethod{}
	switch server.Method {
	case "password":
		sshs = append(sshs, ssh.Password(server.Password))
		break
	case "publickey":
		socket := os.Getenv("SSH_AUTH_SOCK")
		conn, err := net.Dial("unix", socket)
		if err != nil {
			return nil, err
		}
		agentClient := agent.NewClient(conn)
		sshs = append(sshs, ssh.PublicKeysCallback(agentClient.Signers))
		break
	default:
		return nil, errors.New("Invalid password method: " + server.Method)
	}

	return sshs, nil
}
