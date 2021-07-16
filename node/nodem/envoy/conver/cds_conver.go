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

package conver

import (
	"fmt"
	"strconv"
	"strings"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/golang/protobuf/ptypes"
	api_model "github.com/gridworkz/kato/api/model"
	envoyv2 "github.com/gridworkz/kato/node/core/envoy/v2"
	"github.com/gridworkz/kato/node/utils"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

//OneNodeCluster conver cluster of on envoy node
func OneNodeCluster(serviceAlias, namespace string, configs *corev1.ConfigMap, services []*corev1.Service) ([]types.Resource, error) {
	resources, _, err := GetPluginConfigs(configs)
	if err != nil {
		return nil, err
	}
	var clusters []types.Resource
	if resources.BaseServices != nil && len(resources.BaseServices) > 0 {
		for _, cl := range upstreamClusters(serviceAlias, namespace, resources.BaseServices, services) {
			if err := cl.Validate(); err != nil {
				logrus.Errorf("cluster validate failure %s", err.Error())
			} else {
				clusters = append(clusters, cl)
			}
		}
	}
	if resources.BasePorts != nil && len(resources.BasePorts) > 0 {
		for _, cl := range downstreamClusters(serviceAlias, namespace, resources.BasePorts) {
			if err := cl.Validate(); err != nil {
				logrus.Errorf("cluster validate failure %s", err.Error())
			} else {
				clusters = append(clusters, cl)
			}
		}
	}
	if len(clusters) == 0 {
		logrus.Warningf("configmap name: %s; plugin-config: %s; create clusters zero length", configs.Name, configs.Data["plugin-config"])
	}
	return clusters, nil
}

// upstreamClusters handle upstream app cluster
// handle kubernetes inner service
func upstreamClusters(serviceAlias, namespace string, dependsServices []*api_model.BaseService, services []*corev1.Service) (cdsClusters []*v2.Cluster) {
	var clusterConfig = make(map[string]*api_model.BaseService, len(dependsServices))
	for i, dService := range dependsServices {
		depServiceIndex := fmt.Sprintf("%s_%s_%s_%d", namespace, serviceAlias, dService.DependServiceAlias, dService.Port)
		clusterConfig[depServiceIndex] = dependsServices[i]
	}
	for _, service := range services {
		inner, ok := service.Labels["service_type"]
		destServiceAlias := GetServiceAliasByService(service)
		port := service.Spec.Ports[0]
		if !ok || inner != "inner" {
			continue
		}
		getOptions := func() (d envoyv2.KatoPluginOptions) {
			relPort, _ := strconv.Atoi(service.Labels["origin_port"])
			if relPort == 0 {
				relPort = int(port.TargetPort.IntVal)
			}
			depServiceIndex := fmt.Sprintf("%s_%s_%s_%d", namespace, serviceAlias, GetServiceAliasByService(service), relPort)
			if _, ok := clusterConfig[depServiceIndex]; ok {
				return envoyv2.GetOptionValues(clusterConfig[depServiceIndex].Options)
			}
			return envoyv2.GetOptionValues(nil)
		}
		var clusterOption envoyv2.ClusterOptions
		clusterOption.Name = fmt.Sprintf("%s_%s_%s_%v", namespace, serviceAlias, GetServiceAliasByService(service), port.Port)
		options := getOptions()
		clusterOption.OutlierDetection = envoyv2.CreatOutlierDetection(options)
		clusterOption.CircuitBreakers = envoyv2.CreateCircuitBreaker(options)
		clusterOption.ServiceName = fmt.Sprintf("%s_%s_%s_%v", namespace, serviceAlias, destServiceAlias, port.Port)
		if domain, ok := service.Annotations["domain"]; ok && domain != "" {
			logrus.Debugf("domain endpoint[%s], create logical_dns cluster: ", domain)
			clusterOption.ClusterType = v2.Cluster_LOGICAL_DNS
			clusterOption.LoadAssignment = envoyv2.CreateDNSLoadAssignment(serviceAlias, namespace, domain, service)
			if strings.HasPrefix(domain, "https://") {
				splitDomain := strings.Split(domain, "https://")
				if len(splitDomain) == 2 {
					clusterOption.TransportSocket = transportSocket(clusterOption.Name, splitDomain[1])
				}
			}
		} else {
			clusterOption.ClusterType = v2.Cluster_EDS
		}
		clusterOption.HealthyPanicThreshold = options.HealthyPanicThreshold
		clusterOption.ConnectionTimeout = envoyv2.ConverTimeDuration(options.ConnectionTimeout)
		// set port realy protocol
		portProtocol := service.Labels["port_protocol"]
		clusterOption.Protocol = portProtocol
		clusterOption.GrpcHealthServiceName = options.GrpcHealthServiceName
		clusterOption.HealthTimeout = options.HealthCheckTimeout
		clusterOption.HealthInterval = options.HealthCheckInterval
		cluster := envoyv2.CreateCluster(clusterOption)
		if cluster != nil {
			logrus.Debugf("cluster is : %v", cluster)
			cdsClusters = append(cdsClusters, cluster)
		}
	}
	return
}

func transportSocket(name, domain string) *core.TransportSocket {
	logrus.Debugf("https domain tlsContext: %s", domain)
	// refer to: https://www.envoyproxy.io/docs/envoy/v1.17.2/api-v2/api/v2/auth/tls.proto#auth-upstreamtlscontext
	tlsContext, err := ptypes.MarshalAny(&auth.UpstreamTlsContext{Sni: domain})
	if err != nil {
		logrus.Errorf("error marshaling tls context to transport_socket config for cluster %s, err=%v",
			name, err)
		// no tls context for the cluster
		return nil
	}
	return &core.TransportSocket{
		Name: utils.EnvoyTLSSocketName,
		ConfigType: &core.TransportSocket_TypedConfig{
			TypedConfig: tlsContext,
		},
	}
}

//downstreamClusters handle app self cluster
//only local port
func downstreamClusters(serviceAlias, namespace string, ports []*api_model.BasePort) (cdsClusters []*v2.Cluster) {
	for i := range ports {
		port := ports[i]
		address := envoyv2.CreateSocketAddress(port.Protocol, "127.0.0.1", uint32(port.Port))
		clusterName := fmt.Sprintf("%s_%s_%v", namespace, serviceAlias, port.Port)
		option := envoyv2.GetOptionValues(port.Options)
		cluster := envoyv2.CreateCluster(envoyv2.ClusterOptions{
			Name:                     clusterName,
			ConnectionTimeout:        envoyv2.ConverTimeDuration(option.ConnectionTimeout),
			ServiceName:              "",
			ClusterType:              v2.Cluster_STATIC,
			CircuitBreakers:          envoyv2.CreateCircuitBreaker(option),
			OutlierDetection:         envoyv2.CreatOutlierDetection(option),
			MaxRequestsPerConnection: option.MaxRequestsPerConnection,
			Hosts:                    []*core.Address{address},
			HealthyPanicThreshold:    option.HealthyPanicThreshold,
		})
		if cluster != nil {
			cdsClusters = append(cdsClusters, cluster)
		}
	}
	return
}
