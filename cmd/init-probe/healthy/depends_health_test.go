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
	"testing"

	yaml "gopkg.in/yaml.v2"

	envoyv2 "github.com/gridworkz/kato/node/core/envoy/v2"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"google.golang.org/grpc"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
)

var testClusterID = "8cd9214e6b3d4476942b600f41bfefea_tcpmeshd3d6a722b632b854b6c232e4895e0cc6_gr5e0cc6"

var testXDSHost = "39.104.66.227:6101"

// var testClusterID = "2bf54c5a0b5a48a890e2dda8635cb507_tcpmeshed6827c0afdda50599b4108105c9e8e3_grc9e8e3"
//var testXDSHost = "127.0.0.1:6101"

func TestClientListener(t *testing.T) {
	cli, err := grpc.Dial(testXDSHost, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	listenerDiscover := v2.NewListenerDiscoveryServiceClient(cli)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := listenerDiscover.FetchListeners(ctx, &v2.DiscoveryRequest{
		Node: &core.Node{
			Cluster: testClusterID,
			Id:      testClusterID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Resources) == 0 {
		t.Fatal("no listeners")
	}
	t.Logf("version %s", res.GetVersionInfo())
	listeners := envoyv2.ParseListenerResource(res.Resources)
	printYaml(t, listeners)
}

func TestClientCluster(t *testing.T) {
	cli, err := grpc.Dial(testXDSHost, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	clusterDiscover := v2.NewClusterDiscoveryServiceClient(cli)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := clusterDiscover.FetchClusters(ctx, &v2.DiscoveryRequest{
		Node: &core.Node{
			Cluster: testClusterID,
			Id:      testClusterID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Resources) == 0 {
		t.Fatal("no clusters")
	}
	t.Logf("version %s", res.GetVersionInfo())
	clusters := envoyv2.ParseClustersResource(res.Resources)
	for _, cluster := range clusters {
		if cluster.GetType() == v2.Cluster_LOGICAL_DNS {
			fmt.Println(cluster.Name)
		}
		printYaml(t, cluster)
	}
}

func printYaml(t *testing.T, data interface{}) {
	out, _ := yaml.Marshal(data)
	t.Log(string(out))
}

func TestClientEndpoint(t *testing.T) {
	cli, err := grpc.Dial(testXDSHost, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	endpointDiscover := v2.NewEndpointDiscoveryServiceClient(cli)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := endpointDiscover.FetchEndpoints(ctx, &v2.DiscoveryRequest{
		Node: &core.Node{
			Cluster: testClusterID,
			Id:      testClusterID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Resources) == 0 {
		t.Fatal("no endpoints")
	}
	t.Logf("version %s", res.GetVersionInfo())
	endpoints := envoyv2.ParseLocalityLbEndpointsResource(res.Resources)
	for _, e := range endpoints {
		fmt.Println(e.GetClusterName())
	}
	printYaml(t, endpoints)
}

func TestNewDependServiceHealthController(t *testing.T) {
	controller, err := NewDependServiceHealthController()
	if err != nil {
		t.Fatal(err)
	}
	controller.Check()
}
