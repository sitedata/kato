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

package monitormessage

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gridworkz/kato/discover"
	"github.com/gridworkz/kato/discover/config"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/common/log"
)

//UDPServer
type UDPServer struct {
	ctx                 context.Context
	ListenerHost        string
	ListenerPort        int
	eventServerEndpoint []string
	client              net.Conn
	etcdClientArgs      *etcdutil.ClientArgs
}

//CreateUDPServer
func CreateUDPServer(ctx context.Context, lisHost string, lisPort int, etcdClientArgs *etcdutil.ClientArgs) *UDPServer {
	return &UDPServer{
		ctx:            ctx,
		ListenerHost:   lisHost,
		ListenerPort:   lisPort,
		etcdClientArgs: etcdClientArgs,
	}
}

//Start
func (u *UDPServer) Start() error {
	dis, err := discover.GetDiscover(config.DiscoverConfig{Ctx: u.ctx, EtcdClientArgs: u.etcdClientArgs})
	if err != nil {
		return err
	}
	dis.AddProject("event_log_event_udp", u)
	if err := u.server(); err != nil {
		return err
	}
	return nil
}

//UpdateEndpoints update event server address
func (u *UDPServer) UpdateEndpoints(endpoints ...*config.Endpoint) {
	var eventServerEndpoint []string
	for _, e := range endpoints {
		eventServerEndpoint = append(eventServerEndpoint, e.URL)
		u.eventServerEndpoint = eventServerEndpoint
	}
	if len(u.eventServerEndpoint) > 0 {
		for i := range u.eventServerEndpoint {
			info := strings.Split(u.eventServerEndpoint[i], ":")
			if len(info) == 2 {
				dip := net.ParseIP(info[0])
				port, err := strconv.Atoi(info[1])
				if err != nil {
					continue
				}
				srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
				dstAddr := &net.UDPAddr{IP: dip, Port: port}
				conn, err := net.DialUDP("udp", srcAddr, dstAddr)
				if err != nil {
					logrus.Error(err)
					continue
				}
				logrus.Infof("Update event server address is %s", u.eventServerEndpoint[i])
				u.client = conn
				break
			}
		}

	}
}

//Error
func (u *UDPServer) Error(err error) {

}

//Server
func (u *UDPServer) server() error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(u.ListenerHost), Port: u.ListenerPort})
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Infof("UDP Server Listener: %s", listener.LocalAddr().String())
	buf := make([]byte, 65535)
	go func() {
		defer listener.Close()
		for {
			n, _, err := listener.ReadFromUDP(buf)
			if err != nil {
				logrus.Errorf("read message from udp error,%s", err.Error())
				time.Sleep(time.Second * 2)
				continue
			}
			u.handlePacket(buf[0:n])
		}
	}()
	return nil
}

func (u *UDPServer) handlePacket(packet []byte) {
	lines := strings.Split(string(packet), "\n")
	for _, line := range lines {
		if line != "" && u.client != nil {
			u.client.Write([]byte(line))
		}
	}
}
