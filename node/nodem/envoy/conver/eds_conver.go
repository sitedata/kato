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

	"github.com/sirupsen/logrus"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	envoyv2 "github.com/gridworkz/kato/node/core/envoy/v2"
	corev1 "k8s.io/api/core/v1"
)

//OneNodeClusterLoadAssignment one envoy node endpoints
func OneNodeClusterLoadAssignment(serviceAlias, namespace string, endpoints []*corev1.Endpoints, services []*corev1.Service) (clusterLoadAssignment []types.Resource) {
	for i := range services {
		if domain, ok := services[i].Annotations["domain"]; ok && domain != "" {
			logrus.Warnf("service[sid: %s] endpoint id domain endpoint[domain: %s], use dns cluster type, do not create eds", services[i].GetUID(), domain)
			continue
		}
		service := services[i]
		destServiceAlias := GetServiceAliasByService(service)
		if destServiceAlias == "" {
			logrus.Errorf("service alias is empty in k8s service %s", service.Name)
			continue
		}
		clusterName := fmt.Sprintf("%s_%s_%s_%d", namespace, serviceAlias, destServiceAlias, service.Spec.Ports[0].Port)
		selectEndpoint := getEndpointsByServiceName(endpoints, service.Name)
		logrus.Debugf("select endpoints %d for service %s", len(selectEndpoint), service.Name)
		var lendpoints []*endpoint.LocalityLbEndpoints // localityLbEndpoints just support only one content
		for _, en := range selectEndpoint {
			var notReadyAddress *corev1.EndpointAddress
			var notReadyPort *corev1.EndpointPort
			var notreadyToPort int
			for _, subset := range en.Subsets {
				for i, port := range subset.Ports {
					toport := int(port.Port)
					if serviceAlias == destServiceAlias {
						//use real port
						if originPort, ok := service.Labels["origin_port"]; ok {
							origin, err := strconv.Atoi(originPort)
							if err == nil {
								toport = origin
							}
						}
					}
					protocol := string(port.Protocol)
					if len(subset.Addresses) == 0 && len(subset.NotReadyAddresses) > 0 {
						notReadyAddress = &subset.NotReadyAddresses[0]
						notreadyToPort = toport
						notReadyPort = &subset.Ports[i]
					}
					getHealty := func() *endpoint.Endpoint_HealthCheckConfig {
						return &endpoint.Endpoint_HealthCheckConfig{
							PortValue: uint32(toport),
						}
					}
					if len(subset.Addresses) > 0 {
						var lbe []*endpoint.LbEndpoint
						for _, address := range subset.Addresses {
							envoyAddress := envoyv2.CreateSocketAddress(protocol, address.IP, uint32(toport))
							lbe = append(lbe, &endpoint.LbEndpoint{
								HostIdentifier: &endpoint.LbEndpoint_Endpoint{
									Endpoint: &endpoint.Endpoint{
										Address:           envoyAddress,
										HealthCheckConfig: getHealty(),
									},
								},
							})
						}
						if len(lbe) > 0 {
							lendpoints = append(lendpoints, &endpoint.LocalityLbEndpoints{LbEndpoints: lbe})
						}
					}
				}
			}
			if len(lendpoints) == 0 && notReadyAddress != nil && notReadyPort != nil {
				var lbe []*endpoint.LbEndpoint
				envoyAddress := envoyv2.CreateSocketAddress(string(notReadyPort.Protocol), notReadyAddress.IP, uint32(notreadyToPort))
				lbe = append(lbe, &endpoint.LbEndpoint{
					HostIdentifier: &endpoint.LbEndpoint_Endpoint{
						Endpoint: &endpoint.Endpoint{
							Address: envoyAddress,
						},
					},
				})
				lendpoints = append(lendpoints, &endpoint.LocalityLbEndpoints{LbEndpoints: lbe})
			}
		}
		cla := &v2.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints:   lendpoints,
		}
		if err := cla.Validate(); err != nil {
			logrus.Errorf("endpoints discover validate failure %s", err.Error())
		} else {
			clusterLoadAssignment = append(clusterLoadAssignment, cla)
		}
	}
	if len(clusterLoadAssignment) == 0 {
		logrus.Warn("create clusterLoadAssignment zero length")
	}
	return clusterLoadAssignment
}

func getEndpointsByLables(endpoints []*corev1.Endpoints, slabels map[string]string) (re []*corev1.Endpoints) {
	for _, en := range endpoints {
		existLength := 0
		for k, v := range slabels {
			v2, ok := en.Labels[k]
			if ok && v == v2 {
				existLength++
			}
		}
		if existLength == len(slabels) {
			re = append(re, en)
		}
	}
	return
}

func getEndpointsByServiceName(endpoints []*corev1.Endpoints, serviceName string) (re []*corev1.Endpoints) {
	for _, en := range endpoints {
		if serviceName == en.Name {
			re = append(re, en)
		}
	}
	return
}
