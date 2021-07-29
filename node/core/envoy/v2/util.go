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

package v2

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/envoyproxy/go-control-plane/pkg/conversion"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	rsrc "github.com/envoyproxy/go-control-plane/pkg/resource/v2"
	_struct "github.com/golang/protobuf/ptypes/struct"

	v1 "github.com/gridworkz/kato/node/core/envoy/v1"
)

// MessageToStruct converts from proto message to proto Struct
func MessageToStruct(msg proto.Message) *_struct.Struct {
	s, err := conversion.MessageToStruct(msg)
	if err != nil {
		logrus.Error(err.Error())
		return &_struct.Struct{}
	}
	return s
}

// Message2Any converts from proto message to proto any
func Message2Any(msg proto.Message) *any.Any {
	a, err := ptypes.MarshalAny(msg)
	if err != nil {
		logrus.Error(err.Error())
		return &any.Any{}
	}
	return a
}

//ConversionUInt32 conversion uint32 to wrappers uint32
func ConversionUInt32(value uint32) *wrappers.UInt32Value {
	return &wrappers.UInt32Value{
		Value: value,
	}
}

//ConversionTypeUInt32 conversion uint32 to proto uint32
func ConversionTypeUInt32(value uint32) *types.UInt32Value {
	return &types.UInt32Value{
		Value: value,
	}
}

//ConverTimeDuration second
func ConverTimeDuration(second int64) *duration.Duration {
	return &duration.Duration{
		Seconds: second,
	}
}

const (
	//KeyPrefix request path prefix
	KeyPrefix string = "Prefix"
	//KeyHeaders request http headers
	KeyHeaders string = "Headers"
	//KeyDomains request domains
	KeyDomains string = "Domains"
	//KeyMaxConnections The maximum number of connections that Envoy will make to the upstream cluster. If not specified, the default is 1024.
	KeyMaxConnections string = "MaxConnections"
	//KeyMaxPendingRequests The maximum number of pending requests that Envoy will allow to the upstream cluster. If not specified, the default is 1024
	KeyMaxPendingRequests string = "MaxPendingRequests"
	//KeyMaxRequests The maximum number of parallel requests that Envoy will make to the upstream cluster. If not specified, the default is 1024.
	KeyMaxRequests string = "MaxRequests"
	//KeyMaxActiveRetries  The maximum number of parallel retries that Envoy will allow to the upstream cluster. If not specified, the default is 3.
	KeyMaxActiveRetries string = "MaxActiveRetries"
	//KeyUpStream upStream
	KeyUpStream string = "upStream"
	//KeyDownStream downStream
	KeyDownStream string = "downStream"
	//KeyWeight WEIGHT
	KeyWeight string = "Weight"
	//KeyWeightModel MODEL_WEIGHT
	KeyWeightModel string = "weight_model"
	//KeyPrefixModel MODEL_PREFIX
	KeyPrefixModel string = "prefix_model"
	//KeyIntervalMS IntervalMS key
	KeyIntervalMS string = "IntervalMS"
	//KeyConsecutiveErrors ConsecutiveErrors key
	KeyConsecutiveErrors string = "ConsecutiveErrors"
	//KeyBaseEjectionTimeMS BaseEjectionTimeMS key
	KeyBaseEjectionTimeMS string = "BaseEjectionTimeMS"
	//KeyMaxEjectionPercent MaxEjectionPercent key
	KeyMaxEjectionPercent string = "MaxEjectionPercent"
	// KeyMaxRequestsPerConnection Optional maximum requests for a single upstream connection. This parameter
	// is respected by both the HTTP/1.1 and HTTP/2 connection pool
	// implementations. If not specified, there is no limit. Setting this
	// parameter to 1 will effectively disable keep alive.
	KeyMaxRequestsPerConnection string = "MaxRequestsPerConnection"
	// KeyHealthyPanicThreshold default 50,More than 50% of hosts are ejected and go into panic mode
	// Panic mode will send traffic back to the failed host
	KeyHealthyPanicThreshold string = "HealthyPanicThreshold"
	//KeyConnectionTimeout connection timeout setting
	KeyConnectionTimeout string = "ConnectionTimeout"
	//KeyTCPIdleTimeout tcp idle timeout
	KeyTCPIdleTimeout string = "TCPIdleTimeout"
	//KeyGrpcHealthServiceName The name of the grpc service used for health checking.
	KeyGrpcHealthServiceName string = "GrpcHealthServiceName"
	// cluster health check timeout
	KeyHealthCheckTimeout string = "HealthCheckTimeout"
	// cluster health check interval
	KeyHealthCheckInterval string = "HealthCheckInterval"
)

//KatoPluginOptions kato plugin config struct
type KatoPluginOptions struct {
	Prefix                   string
	MaxConnections           int
	MaxRequests              int
	MaxPendingRequests       int
	MaxActiveRetries         int
	Headers                  v1.Headers
	Domains                  []string
	Weight                   uint32
	Interval                 int64
	ConsecutiveErrors        int
	BaseEjectionTimeMS       int64
	MaxEjectionPercent       int
	MaxRequestsPerConnection *uint32
	HealthyPanicThreshold    int64
	ConnectionTimeout        int64
	TCPIdleTimeout           int64
	GrpcHealthServiceName    string
	HealthCheckTimeout       int64
	HealthCheckInterval      int64
}

//KatoInboundPluginOptions kato inbound plugin options
type KatoInboundPluginOptions struct {
	OpenLimit   bool
	LimitDomain string
}

//RouteBasicHash get basic hash for weight
func (r KatoPluginOptions) RouteBasicHash() string {
	key := sha256.New()
	var header string
	sort.Sort(r.Headers)
	for _, h := range r.Headers {
		header += h.Name + h.Value
	}
	key.Write([]byte(r.Prefix + header + strings.Join(r.Domains, "")))
	return string(key.Sum(nil))
}

//GetOptionValues get value from options
//if not exist,return default value
func GetOptionValues(sr map[string]interface{}) KatoPluginOptions {
	rpo := KatoPluginOptions{
		Prefix:                "/",
		MaxConnections:        10240,
		MaxRequests:           10240,
		MaxPendingRequests:    1024,
		MaxActiveRetries:      3,
		Domains:               []string{"*"},
		Weight:                100,
		Interval:              10,
		ConsecutiveErrors:     5,
		BaseEjectionTimeMS:    30000,
		MaxEjectionPercent:    10,
		HealthyPanicThreshold: 50,
		ConnectionTimeout:     250,
		TCPIdleTimeout:        60 * 60 * 2,
		HealthCheckTimeout:    5,
		HealthCheckInterval:   4,
	}
	if sr == nil {
		return rpo
	}
	for kind, v := range sr {
		switch kind {
		case KeyPrefix:
			rpo.Prefix = strings.TrimSpace(v.(string))
		case KeyMaxConnections:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				rpo.MaxConnections = i
			}
		case KeyMaxRequests:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				rpo.MaxRequests = i
			}
		case KeyMaxPendingRequests:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				rpo.MaxPendingRequests = i
			}
		case KeyMaxActiveRetries:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				rpo.MaxActiveRetries = i
			}
		case KeyHeaders:
			parents := strings.Split(v.(string), ";")
			var hm v1.Header
			for _, h := range parents {
				headers := strings.Split(h, ":")
				//has_header:no default
				if len(headers) == 2 {
					if headers[0] == "has_header" && headers[1] == "no" {
						continue
					}
					hm.Name = headers[0]
					hm.Value = headers[1]
				}
			}
			rpo.Headers = append(rpo.Headers, hm)
		case KeyDomains:
			if strings.Contains(v.(string), ",") {
				rpo.Domains = strings.Split(v.(string), ",")
			} else if v.(string) != "" {
				rpo.Domains = []string{v.(string)}
			}
		case KeyWeight:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				rpo.Weight = uint32(i)
			}
		case KeyIntervalMS:
			if i, err := strconv.Atoi(v.(string)); err == nil && i < 0 {
				rpo.Interval = int64(i)
			}
		case KeyConsecutiveErrors:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				rpo.ConsecutiveErrors = i
			}
		case KeyBaseEjectionTimeMS:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				rpo.BaseEjectionTimeMS = int64(i)
			}
		case KeyMaxEjectionPercent:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				if i > 100 {
					rpo.MaxEjectionPercent = 100
				} else {
					rpo.MaxEjectionPercent = i
				}
			}
		case KeyMaxRequestsPerConnection:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				value := uint32(i)
				rpo.MaxRequestsPerConnection = &value
			}
		case KeyHealthyPanicThreshold:
			if i, err := strconv.Atoi(v.(string)); err == nil && i != 0 {
				if i > 100 {
					rpo.HealthyPanicThreshold = 100
				} else {
					rpo.HealthyPanicThreshold = int64(i)
				}
			}
		case KeyConnectionTimeout:
			if i, err := strconv.Atoi(v.(string)); err == nil {
				rpo.ConnectionTimeout = int64(i)
			}
		case KeyTCPIdleTimeout:
			if i, err := strconv.Atoi(v.(string)); err == nil {
				rpo.TCPIdleTimeout = int64(i)
			}
		case KeyHealthCheckInterval:
			if i, err := strconv.Atoi(v.(string)); err == nil {
				rpo.HealthCheckInterval = int64(i)
			}
		case KeyHealthCheckTimeout:
			if i, err := strconv.Atoi(v.(string)); err == nil {
				rpo.HealthCheckTimeout = int64(i)
			}
		case KeyGrpcHealthServiceName:
			rpo.GrpcHealthServiceName = strings.TrimSpace(v.(string))
		}
	}
	return rpo
}

//GetKatoInboundPluginOptions get kato inbound plugin options
func GetKatoInboundPluginOptions(sr map[string]interface{}) (r KatoInboundPluginOptions) {
	for k, v := range sr {
		switch k {
		case "OPEN_LIMIT":
			if strings.ToLower(v.(string)) == "yes" || strings.ToLower(v.(string)) == "true" {
				r.OpenLimit = true
			}
		case "LIMIT_DOMAIN":
			r.LimitDomain = v.(string)
		}
	}
	return
}

//ParseLocalityLbEndpointsResource parse envoy xds server response ParseLocalityLbEndpointsResource
func ParseLocalityLbEndpointsResource(resources []*any.Any) []v2.ClusterLoadAssignment {
	var endpoints []v2.ClusterLoadAssignment
	for _, resource := range resources {
		switch resource.GetTypeUrl() {
		case rsrc.EndpointType:
			var endpoint v2.ClusterLoadAssignment
			if err := proto.Unmarshal(resource.GetValue(), &endpoint); err != nil {
				logrus.Errorf("unmarshal envoy endpoint resource failure %s", err.Error())
			}
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

//ParseClustersResource parse envoy xds server response ParseClustersResource
func ParseClustersResource(resources []*any.Any) []v2.Cluster {
	var clusters []v2.Cluster
	for _, resource := range resources {
		switch resource.GetTypeUrl() {
		case rsrc.ClusterType:
			var cluster v2.Cluster
			if err := proto.Unmarshal(resource.GetValue(), &cluster); err != nil {
				logrus.Errorf("unmarshal envoy cluster resource failure %s", err.Error())
			}
			clusters = append(clusters, cluster)
		}
	}
	return clusters
}

//ParseListenerResource parse envoy xds server response ListenersResource
func ParseListenerResource(resources []*any.Any) []v2.Listener {
	var listeners []v2.Listener
	for _, resource := range resources {
		switch resource.GetTypeUrl() {
		case rsrc.ListenerType:
			var listener v2.Listener
			if err := proto.Unmarshal(resource.GetValue(), &listener); err != nil {
				logrus.Errorf("unmarshal envoy listener resource failure %s", err.Error())
			}
			listeners = append(listeners, listener)
		}
	}
	return listeners
}

//ParseRouteConfigurationsResource parse envoy xds server response RouteConfigurationsResource
func ParseRouteConfigurationsResource(resources []*any.Any) []v2.RouteConfiguration {
	var routes []v2.RouteConfiguration
	for _, resource := range resources {
		switch resource.GetTypeUrl() {
		case rsrc.RouteType:
			var route v2.RouteConfiguration
			if err := proto.Unmarshal(resource.GetValue(), &route); err != nil {
				logrus.Errorf("unmarshal envoy route resource failure %s", err.Error())
			}
			routes = append(routes, route)
		}
	}
	return routes
}

//CheckWeightSum check all cluster weight sum
func CheckWeightSum(clusters []*route.WeightedCluster_ClusterWeight, weight uint32) uint32 {
	var sum uint32
	for _, cluster := range clusters {
		sum += cluster.Weight.GetValue()
	}
	if sum >= 100 {
		return 0
	}
	if (sum + weight) > 100 {
		return 100 - sum
	}
	return weight
}

// CheckDomain check and handling http domain
// fix grpc issues https://github.com/envoyproxy/envoy/issues/886
// after https://github.com/envoyproxy/envoy/pull/10960 merge version, This logic can be removed.
func CheckDomain(domain []string, protocol string) []string {
	if protocol == "grpc" {
		var newDomain []string
		for _, d := range domain {
			if !strings.Contains(d, ":") {
				newDomain = append(newDomain, fmt.Sprintf("%s:%d", d, DefaultLocalhostListenerPort))
			} else {
				di := strings.Split(d, ":")
				if len(di) == 2 && di[1] != fmt.Sprintf("%d", DefaultLocalhostListenerPort) {
					newDomain = append(newDomain, fmt.Sprintf("%s:%d", di[0], DefaultLocalhostListenerPort))
				} else {
					newDomain = append(newDomain, d)
				}
			}
		}
		return newDomain
	}
	return domain
}
