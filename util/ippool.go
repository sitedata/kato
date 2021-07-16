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

package util

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

//IPEVENTTYPE ip change event type
type IPEVENTTYPE int

//IPEVENT ip change event
type IPEVENT struct {
	Type IPEVENTTYPE
	IP   net.IP
}

//IPPool
type IPPool struct {
	ctx                 context.Context
	cancel              context.CancelFunc
	lock                sync.Mutex
	HostIPs             map[string]net.IP
	EventCh             chan IPEVENT
	StopCh              chan struct{}
	ignoreInterfaceName []string
	startReady          chan struct{}
	once                sync.Once
}

const (
	//ADD add event
	ADD IPEVENTTYPE = iota
	//DEL del event
	DEL
	//UPDATE update event
	UPDATE
)

//NewIPPool
func NewIPPool(ignoreInterfaceName []string) *IPPool {
	ctx, cancel := context.WithCancel(context.Background())
	ippool := &IPPool{
		ctx:                 ctx,
		cancel:              cancel,
		HostIPs:             map[string]net.IP{},
		EventCh:             make(chan IPEVENT, 1024),
		StopCh:              make(chan struct{}),
		ignoreInterfaceName: ignoreInterfaceName,
		startReady:          make(chan struct{}),
	}
	return ippool
}

//Ready
func (i *IPPool) Ready() bool {
	logrus.Info("waiting ip pool start ready")
	<-i.startReady
	return true
}

//GetHostIPs
func (i *IPPool) GetHostIPs() []net.IP {
	i.lock.Lock()
	defer i.lock.Unlock()
	var ips []net.IP
	for _, ip := range i.HostIPs {
		ips = append(ips, ip)
	}
	return ips
}

//GetWatchIPChan watch ip change
func (i *IPPool) GetWatchIPChan() <-chan IPEVENT {
	return i.EventCh
}

//Close
func (i *IPPool) Close() {
	i.cancel()
}

//LoopCheckIPs
func (i *IPPool) LoopCheckIPs() {
	Exec(i.ctx, func() error {
		logrus.Debugf("start loop watch ips from all interfaces")
		ips, err := i.getInterfaceIPs()
		if err != nil {
			logrus.Errorf("got ip address from interface failure %s, will retry", err.Error())
			return nil
		}
		i.lock.Lock()
		defer i.lock.Unlock()
		var newIP = make(map[string]net.IP)
		for _, v := range ips {
			if v.To4() == nil {
				continue
			}
			_, ok := i.HostIPs[v.To4().String()]
			if ok {
				i.EventCh <- IPEVENT{Type: UPDATE, IP: v.To4()}
			}
			if !ok {
				i.EventCh <- IPEVENT{Type: ADD, IP: v.To4()}
			}
			newIP[v.To4().String()] = v.To4()
		}
		for k, v := range i.HostIPs {
			if _, ok := newIP[k]; !ok {
				i.EventCh <- IPEVENT{Type: DEL, IP: v.To4()}
			}
		}
		logrus.Debugf("loop watch ips from all interfaces, find %d ips", len(newIP))
		i.HostIPs = newIP
		i.once.Do(func() {
			close(i.startReady)
		})
		return nil
	}, time.Second*5)
}

func (i *IPPool) getInterfaceIPs() ([]net.IP, error) {
	var ips []net.IP
	tables, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		if StringArrayContains(i.ignoreInterfaceName, table.Name) {
			logrus.Debugf("ip address from interface %s ignore", table.Name)
			continue
		}
		addrs, err := table.Addrs()
		if err != nil {
			return nil, err
		}
		for _, address := range addrs {
			if ipnet := checkIPAddress(address); ipnet != nil {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.To4())
				}
			}
		}
	}
	return ips, nil
}

func checkIPAddress(addr net.Addr) *net.IPNet {
	ipnet, ok := addr.(*net.IPNet)
	if !ok {
		return nil
	}
	if ipnet.IP.IsLoopback() {
		return nil
	}
	return ipnet
}
