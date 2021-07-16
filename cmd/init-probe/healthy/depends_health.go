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

package healthy

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	endpointapi "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoyv2 "github.com/gridworkz/kato/node/core/envoy/v2"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//DependServiceHealthController Detect the health of the dependent service
//Health based conditions：
//------- lds: discover all dependent services
//------- cds: discover all dependent services
//------- sds: every service has at least one Ready instance
type DependServiceHealthController struct {
	listeners                       []v2.Listener
	clusters                        []v2.Cluster
	sdsHost                         []v2.ClusterLoadAssignment
	interval                        time.Duration
	envoyDiscoverVersion            string //only support v2
	checkFunc                       []func() bool
	endpointClient                  v2.EndpointDiscoveryServiceClient
	clusterClient                   v2.ClusterDiscoveryServiceClient
	dependServiceCount              int
	clusterID                       string
	dependServiceNames              []string
	ignoreCheckEndpointsClusterName []string
}

//NewDependServiceHealthController create a controller
func NewDependServiceHealthController() (*DependServiceHealthController, error) {
	clusterID := os.Getenv("ENVOY_NODE_ID")
	if clusterID == "" {
		clusterID = fmt.Sprintf("%s_%s_%s", os.Getenv("TENANT_ID"), os.Getenv("PLUGIN_ID"), os.Getenv("SERVICE_NAME"))
	}
	dsc := DependServiceHealthController{
		interval:  time.Second * 5,
		clusterID: clusterID,
	}
	dsc.checkFunc = append(dsc.checkFunc, dsc.checkListener)
	dsc.checkFunc = append(dsc.checkFunc, dsc.checkClusters)
	dsc.checkFunc = append(dsc.checkFunc, dsc.checkEDS)
	xDSHost := os.Getenv("XDS_HOST_IP")
	xDSHostPort := os.Getenv("XDS_HOST_PORT")
	if xDSHostPort == "" {
		xDSHostPort = "6101"
	}
	cli, err := grpc.Dial(fmt.Sprintf("%s:%s", xDSHost, xDSHostPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	dsc.endpointClient = v2.NewEndpointDiscoveryServiceClient(cli)
	dsc.clusterClient = v2.NewClusterDiscoveryServiceClient(cli)
	dsc.dependServiceNames = strings.Split(os.Getenv("STARTUP_SEQUENCE_DEPENDENCIES"), ",")
	return &dsc, nil
}

//Check - check all conditions
func (d *DependServiceHealthController) Check() {
	logrus.Info("start denpenent health check.")
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	check := func() bool {
		for _, check := range d.checkFunc {
			if !check() {
				return false
			}
		}
		return true
	}
	for {
		if check() {
			logrus.Info("Dependent services all check passed, will start service")
			return
		}
		select {
		case <-ticker.C:
		}
	}
}

func (d *DependServiceHealthController) checkListener() bool {
	return true
}

func (d *DependServiceHealthController) checkClusters() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := d.clusterClient.FetchClusters(ctx, &v2.DiscoveryRequest{
		Node: &core.Node{
			Cluster: d.clusterID,
			Id:      d.clusterID,
		},
	})
	if err != nil {
		logrus.Errorf("discover dependent services cluster failure %s", err.Error())
		return false
	}
	clusters := envoyv2.ParseClustersResource(res.Resources)
	d.ignoreCheckEndpointsClusterName = nil
	for _, cluster := range clusters {
		if cluster.GetType() == v2.Cluster_LOGICAL_DNS {
			d.ignoreCheckEndpointsClusterName = append(d.ignoreCheckEndpointsClusterName, cluster.Name)
		}
	}
	d.clusters = clusters
	return true
}

func (d *DependServiceHealthController) checkEDS() bool {
	logrus.Infof("start checking eds; dependent service cluster names: %s", d.dependServiceNames)
	if len(d.clusters) == len(d.ignoreCheckEndpointsClusterName) {
		logrus.Info("all dependent services is domain third service.")
		return true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := d.endpointClient.FetchEndpoints(ctx, &v2.DiscoveryRequest{
		Node: &core.Node{
			Cluster: d.clusterID,
			Id:      d.clusterID,
		},
	})
	if err != nil {
		logrus.Errorf("discover dependent services endpoint failure %s", err.Error())
		return false
	}
	clusterLoadAssignments := envoyv2.ParseLocalityLbEndpointsResource(res.Resources)
	readyClusters := make(map[string]bool, len(clusterLoadAssignments))
	for _, cla := range clusterLoadAssignments {
		// clusterName := fmt.Sprintf("%s_%s_%s_%d", namespace, serviceAlias, destServiceAlias, service.Spec.Ports[0].Port)
		serviceName := ""
		clusterNameInfo := strings.Split(cla.GetClusterName(), "_")
		if len(clusterNameInfo) == 4 {
			serviceName = clusterNameInfo[2]
		}
		if serviceName == "" {
			continue
		}
		if ready, exist := readyClusters[serviceName]; exist && ready {
			continue
		}

		ready := func() bool {
			if util.StringArrayContains(d.ignoreCheckEndpointsClusterName, cla.ClusterName) {
				return true
			}
			if len(cla.Endpoints) > 0 && len(cla.Endpoints[0].LbEndpoints) > 0 {
				// first LbEndpoints healthy is not nil. so endpoint is not notreadyaddress
				if host, ok := cla.Endpoints[0].LbEndpoints[0].HostIdentifier.(*endpointapi.LbEndpoint_Endpoint); ok {
					if host.Endpoint != nil && host.Endpoint.HealthCheckConfig != nil {
						logrus.Infof("depend service (%s) start complete", cla.ClusterName)
						return true
					}
				}
			}
			return false
		}()
		logrus.Infof("cluster name: %s; ready: %v", serviceName, ready)
		readyClusters[serviceName] = ready
	}
	for _, ignoreCheckEndpointsClusterName := range d.ignoreCheckEndpointsClusterName {
		clusterNameInfo := strings.Split(ignoreCheckEndpointsClusterName, "_")
		if len(clusterNameInfo) == 4 {
			readyClusters[clusterNameInfo[2]] = true
		}
	}
	for _, cn := range d.dependServiceNames {
		if cn != "" {
			if ready := readyClusters[cn]; !ready {
				logrus.Infof("%s not ready.", cn)
				return false
			}
		}
	}
	logrus.Info("all dependent services have been started.")

	return true
}
