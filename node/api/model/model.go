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

package model

import (
	//"github.com/sirupsen/logrus"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/gridworkz/kato/node/utils"
	"github.com/pquerna/ffjson/ffjson"
	v1 "k8s.io/api/core/v1" //"github.com/sirupsen/logrus"
)

//Resource
type Resource struct {
	CPU  int `json:"cpu"`
	MemR int `json:"mem"`
}

// NodePodResource -
type NodePodResource struct {
	AllocatedResources `json:"allocatedresources"`
	Resource           `json:"allocatable"`
}

// AllocatedResources -
type AllocatedResources struct {
	CPURequests     int64
	CPULimits       int64
	MemoryRequests  int64
	MemoryLimits    int64
	MemoryRequestsR string
	MemoryLimitsR   string
	CPURequestsR    string
	CPULimitsR      string
}

// InitStatus -
type InitStatus struct {
	Status   int    `json:"status"`
	StatusCN string `json:"cn"`
	HostID   string `json:"uuid"`
}

// InstallStatus -
type InstallStatus struct {
	Status   int           `json:"status"`
	StatusCN string        `json:"cn"`
	Tasks    []*ExecedTask `json:"tasks"`
}

// ExecedTask -
type ExecedTask struct {
	ID             string   `json:"id"`
	Seq            int      `json:"seq"`
	Desc           string   `json:"desc"`
	Status         string   `json:"status"`
	CompleteStatus string   `json:"complete_status"`
	ErrorMsg       string   `json:"err_msg"`
	Depends        []string `json:"dep"`
	Next           []string `json:"next"`
}

// Prome -
type Prome struct {
	Status string    `json:"status"`
	Data   PromeData `json:"data"`
}

// PromeData -
type PromeData struct {
	ResultType string             `json:"resultType"`
	Result     []*PromeResultCore `json:"result"`
}

// PromeResultCore -
type PromeResultCore struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
	Values []interface{}     `json:"values"`
}

// Expr swagger:parameters createToken
type Expr struct {
	Body struct {
		// expr
		// in: body
		// required: true
		Expr string `json:"expr" validate:"expr|required"`
	}
}

// LabelsResp -
type LabelsResp struct {
	SysLabels    map[string]string `json:"sys_labels"`
	CustomLabels map[string]string `json:"custom_labels"`
}

// PrometheusInterface -
type PrometheusInterface interface {
	Query(query string) *Prome
	QueryRange(query string, start, end, step string) *Prome
}

// PrometheusAPI -
type PrometheusAPI struct {
	API string
}

//Query Get
func (s *PrometheusAPI) Query(query string) (*Prome, *utils.APIHandleError) {
	resp, code, err := DoRequest(s.API, query, "query", "GET", nil)
	if err != nil {
		return nil, utils.CreateAPIHandleError(400, err)
	}
	if code == 422 {
		return nil, utils.CreateAPIHandleError(422, fmt.Errorf("unprocessable entity,expression %s can't be executed", query))
	}
	if code == 400 {
		return nil, utils.CreateAPIHandleError(400, fmt.Errorf("bad request,error to request query %s", query))
	}
	if code == 503 {
		return nil, utils.CreateAPIHandleError(503, fmt.Errorf("service unavailable"))
	}
	var prome Prome
	err = ffjson.Unmarshal(resp, &prome)
	if err != nil {
		return nil, utils.CreateAPIHandleError(500, err)
	}
	return &prome, nil
}

//QueryRange Get
func (s *PrometheusAPI) QueryRange(query string, start, end, step string) (*Prome, *utils.APIHandleError) {
	//logrus.Infof("prometheus api is %s",s.API)
	uri := fmt.Sprintf("%v&start=%v&end=%v&step=%v", query, start, end, step)
	resp, code, err := DoRequest(s.API, uri, "query_range", "GET", nil)
	if err != nil {
		return nil, utils.CreateAPIHandleError(400, err)
	}
	if code == 422 {
		return nil, utils.CreateAPIHandleError(422, fmt.Errorf("unprocessable entity,expression %s can't be executed", query))
	}
	if code == 400 {
		return nil, utils.CreateAPIHandleError(400, fmt.Errorf("bad request,error to request query %s", query))
	}
	if code == 503 {
		return nil, utils.CreateAPIHandleError(503, fmt.Errorf("service unavailable"))
	}
	var prome Prome
	err = ffjson.Unmarshal(resp, &prome)
	if err != nil {
		return nil, utils.CreateAPIHandleError(500, err)
	}
	return &prome, nil
}

// DoRequest -
func DoRequest(baseAPI, query, queryType, method string, body []byte) ([]byte, int, error) {
	api := baseAPI + "/api/v1/" + queryType + "?"
	query = "query=" + query
	query = strings.Replace(query, "+", "%2B", -1)
	val, err := url2.ParseQuery(query)
	if err != nil {
		return nil, 0, err
	}
	encoded := val.Encode()
	//logrus.Infof("uri is %s",api+encoded)
	request, err := http.NewRequest(method, api+encoded, nil)
	if err != nil {
		return nil, 0, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, 0, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	return data, resp.StatusCode, nil
}

//ClusterResource
type ClusterResource struct {
	AllNode                          int           `json:"all_node"`
	NotReadyNode                     int           `json:"notready_node"`
	ComputeNode                      int           `json:"compute_node"`
	Tenant                           int           `json:"tenant"`
	CapCPU                           int           `json:"cap_cpu"`          //Total CPU allocation
	CapMem                           int           `json:"cap_mem"`          //Total Distributable Mem
	HealthCapCPU                     int           `json:"health_cap_cpu"`   //Health can allocate CPU
	HealthCapMem                     int           `json:"health_cap_mem"`   //Health Distributable Mem
	UnhealthCapCPU                   int           `json:"unhealth_cap_cpu"` //Unhealthy CPU can be allocated
	UnhealthCapMem                   int           `json:"unhealth_cap_mem"` //Unhealthy Mem can be allocated
	ReqCPU                           float32       `json:"req_cpu"`          //Total CPU used
	ReqMem                           int           `json:"req_mem"`          //Total Mem Used
	HealthReqCPU                     float32       `json:"health_req_cpu"`   //Health has used CPU
	HealthReqMem                     int           `json:"health_req_mem"`   //Health has used Mem
	UnhealthReqCPU                   float32       `json:"unhealth_req_cpu"` //Unhealthy used CPU
	UnhealthReqMem                   int           `json:"unhealth_req_mem"` //Unhealthy has used Mem
	CapDisk                          uint64        `json:"cap_disk"`
	ReqDisk                          uint64        `json:"req_disk"`
	MaxAllocatableMemoryNodeResource *NodeResource `json:"max_allocatable_memory_node_resource"`
}

//NodeResourceResponse
type NodeResourceResponse struct {
	CapCPU int     `json:"cap_cpu"`
	CapMem int     `json:"cap_mem"`
	ReqCPU float32 `json:"req_cpu"`
	ReqMem int     `json:"req_mem"`
}

// FirstConfig -
type FirstConfig struct {
	StorageMode     string `json:"storage_mode"`
	StorageHost     string `json:"storage_host,omitempty"`
	StorageEndPoint string `json:"storage_endpoint,omitempty"`

	NetworkMode string `json:"network_mode"`
	ZKHosts     string `json:"zk_host,omitempty"`
	CassandraIP string `json:"cassandra_ip,omitempty"`
	K8SAPIAddr  string `json:"k8s_apiserver,omitempty"`
	MasterIP    string `json:"master_ip,omitempty"`
	DNS         string `json:"dns,omitempty"`
	ZMQSub      string `json:"zmq_sub,omitempty"`
	ZMQTo       string `json:"zmq_to,omitempty"`
	EtcdIP      string `json:"etcd_ip,omitempty"`
}

// Config -
type Config struct {
	Cn    string `json:"cn_name"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

//ConfigUnit
type ConfigUnit struct {
	//Configuration name: network
	Name   string `json:"name" validate:"name|required"`
	CNName string `json:"cn_name" validate:"cn_name"`
	//Type for example: midonet
	Value     interface{} `json:"value" validate:"value|required"`
	ValueType string      `json:"value_type"`
	//Optional type Type name and required configuration items
	OptionalValue []string                `json:"optional_value,omitempty"`
	DependConfig  map[string][]ConfigUnit `json:"depend_config,omitempty"`
	//Is user configurable
	IsConfigurable bool `json:"is_configurable"`
}

func (c ConfigUnit) String() string {
	res, _ := ffjson.Marshal(&c)
	return string(res)
}

//GlobalConfig
type GlobalConfig struct {
	Configs map[string]*ConfigUnit `json:"configs"`
}

//String
func (g *GlobalConfig) String() string {
	res, _ := ffjson.Marshal(g)
	return string(res)
}

//Add
func (g *GlobalConfig) Add(c ConfigUnit) {
	//With dependent configuration
	if c.DependConfig != nil || len(c.DependConfig) > 0 {
		if c.ValueType == "string" || c.ValueType == "" {
			if value, ok := c.Value.(string); ok {
				for _, dc := range c.DependConfig[value] {
					g.Add(dc)
				}
			}
		}
	}
	g.Configs[c.Name] = &c
}

//Get
func (g *GlobalConfig) Get(name string) *ConfigUnit {
	return g.Configs[name]
}

//Delete
func (g *GlobalConfig) Delete(Name string) {
	if _, ok := g.Configs[Name]; ok {
		delete(g.Configs, Name)
	}
}

//Bytes
func (g GlobalConfig) Bytes() []byte {
	res, _ := ffjson.Marshal(&g)
	return res
}

//CreateDefaultGlobalConfig
func CreateDefaultGlobalConfig() *GlobalConfig {
	gconfig := &GlobalConfig{
		Configs: make(map[string]*ConfigUnit),
	}
	gconfig.Add(ConfigUnit{
		Name:      "NETWORK_MODE",
		CNName:    "Cluster network mode",
		Value:     "calico",
		ValueType: "string",
		DependConfig: map[string][]ConfigUnit{
			"calico": []ConfigUnit{ConfigUnit{Name: "ETCD_ADDRS", CNName: "ETCD address", ValueType: "array"}},
			"midonet": []ConfigUnit{
				ConfigUnit{Name: "CASSANDRA_ADDRS", CNName: "CASSANDRA address", ValueType: "array"},
				ConfigUnit{Name: "ZOOKEEPER_ADDRS", CNName: "ZOOKEEPER address", ValueType: "array"},
				ConfigUnit{Name: "LB_CIDR", CNName: "Network segment where load balancing is located", ValueType: "string"},
			}},
		IsConfigurable: true,
	})
	gconfig.Add(ConfigUnit{
		Name:   "STORAGE_MODE",
		Value:  "nfs",
		CNName: "Default shared storage mode",
		DependConfig: map[string][]ConfigUnit{
			"nfs": []ConfigUnit{
				ConfigUnit{Name: "NFS_SERVERS", CNName: "NFS server address list", ValueType: "array"},
				ConfigUnit{Name: "NFS_ENDPOINT", CNName: "NFS mount endpoint", ValueType: "string"},
			},
			"clusterfs": []ConfigUnit{},
		},
		IsConfigurable: true,
	})
	gconfig.Add(ConfigUnit{
		Name:          "DB_MODE",
		Value:         "mysql",
		CNName:        "Management node database type",
		OptionalValue: []string{"mysql", "yugabytedb"},
		DependConfig: map[string][]ConfigUnit{
			"mysql": []ConfigUnit{
				ConfigUnit{Name: "MYSQL_HOST", CNName: "Mysql database address", ValueType: "string", Value: "127.0.0.1"},
				ConfigUnit{Name: "MYSQL_PASS", CNName: "Mysql database password", ValueType: "string", Value: ""},
				ConfigUnit{Name: "MYSQL_USER", CNName: "Mysql database user name", ValueType: "string", Value: ""},
			},
			"yugabytedb": []ConfigUnit{
				ConfigUnit{Name: "YUGABYTE_HOST", CNName: "Mysql database address", ValueType: "array"},
				ConfigUnit{Name: "YUGABYTE_PASS", CNName: "Mysql database password", ValueType: "string"},
				ConfigUnit{Name: "YUGABYTE_USER", CNName: "Mysql database user name", ValueType: "string"},
			},
		},
		IsConfigurable: true,
	})
	gconfig.Add(ConfigUnit{
		Name:           "LB_MODE",
		Value:          "nginx",
		ValueType:      "string",
		CNName:         "Edge load balancing",
		OptionalValue:  []string{"nginx", "zeus"},
		IsConfigurable: true,
	})
	gconfig.Add(ConfigUnit{Name: "DOMAIN", CNName: "Application domain", ValueType: "string"})
	gconfig.Add(ConfigUnit{Name: "INSTALL_NODE", CNName: "Install node", ValueType: "array"})
	gconfig.Add(ConfigUnit{
		Name:           "INSTALL_MODE",
		Value:          "online",
		ValueType:      "string",
		CNName:         "Installation mode",
		OptionalValue:  []string{"online", "offine"},
		IsConfigurable: true,
	})
	gconfig.Add(ConfigUnit{
		Name:      "DNS_SERVER",
		Value:     []string{},
		CNName:    "Cluster DNS service",
		ValueType: "array",
	})
	gconfig.Add(ConfigUnit{
		Name:      "KUBE_API",
		Value:     []string{},
		ValueType: "array",
		CNName:    "Kubernetes API service",
	})
	gconfig.Add(ConfigUnit{
		Name:      "MANAGE_NODE_ADDRESS",
		Value:     []string{},
		ValueType: "array",
		CNName:    "Management node",
	})
	return gconfig
}

//CreateGlobalConfig
func CreateGlobalConfig(kvs []*mvccpb.KeyValue) (*GlobalConfig, error) {
	dgc := &GlobalConfig{
		Configs: make(map[string]*ConfigUnit),
	}
	for _, kv := range kvs {
		var cn ConfigUnit
		if err := ffjson.Unmarshal(kv.Value, &cn); err == nil {
			dgc.Add(cn)
		}
	}
	return dgc, nil
}

// LoginResult -
type LoginResult struct {
	HostPort  string `json:"hostport"`
	LoginType bool   `json:"type"`
	Result    string `json:"result"`
}

// Login -
type Login struct {
	HostPort  string `json:"hostport"`
	LoginType bool   `json:"type"`
	HostType  string `json:"hosttype"`
	RootPwd   string `json:"pwd,omitempty"`
}

// Body -
type Body struct {
	List interface{} `json:"list"`
	Bean interface{} `json:"bean,omitempty"`
}

// ResponseBody -
type ResponseBody struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	MsgCN string `json:"msgcn"`
	Body  Body   `json:"body,omitempty"`
}

// Pods -
type Pods struct {
	Namespace       string `json:"namespace"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	TenantName      string `json:"tenant_name"`
	CPURequests     string `json:"cpurequest"`
	CPURequestsR    string `json:"cpurequestr"`
	CPULimits       string `json:"cpulimits"`
	CPULimitsR      string `json:"cpulimitsr"`
	MemoryRequests  string `json:"memoryrequests"`
	MemoryRequestsR string `json:"memoryrequestsr"`
	MemoryLimits    string `json:"memorylimits"`
	MemoryLimitsR   string `json:"memorylimitsr"`
	Status          string `json:"status"`
}

//NodeDetails
type NodeDetails struct {
	Name               string              `json:"name"`
	Role               []string            `json:"role"`
	Status             string              `json:"status"`
	Labels             map[string]string   `json:"labels"`
	Annotations        map[string]string   `json:"annotations"`
	CreationTimestamp  string              `json:"creationtimestamp"`
	Conditions         []v1.NodeCondition  `json:"conditions"`
	Addresses          map[string]string   `json:"addresses"`
	Capacity           map[string]string   `json:"capacity"`
	Allocatable        map[string]string   `json:"allocatable"`
	SystemInfo         v1.NodeSystemInfo   `json:"systeminfo"`
	NonterminatedPods  []*Pods             `json:"nonterminatedpods"`
	AllocatedResources map[string]string   `json:"allocatedresources"`
	Events             map[string][]string `json:"events"`
}

// AlertingRulesConfig -
type AlertingRulesConfig struct {
	Groups []*AlertingNameConfig `yaml:"groups" json:"groups"`
}

// AlertingNameConfig -
type AlertingNameConfig struct {
	Name  string         `yaml:"name" json:"name"`
	Rules []*RulesConfig `yaml:"rules" json:"rules"`
}

// RulesConfig -
type RulesConfig struct {
	Alert       string            `yaml:"alert" json:"alert"`
	Expr        string            `yaml:"expr" json:"expr"`
	For         string            `yaml:"for" json:"for"`
	Labels      map[string]string `yaml:"labels" json:"labels"`
	Annotations map[string]string `yaml:"annotations" json:"annotations"`
}

//NotificationEvent
type NotificationEvent struct {
	//Kind could be service, tenant, cluster, node
	Kind string `json:"Kind"`
	//KindID could be service_id,tenant_id,cluster_id,node_id
	KindID string `json:"KindID"`
	Hash   string `json:"Hash"`
	//Type could be Normal / UnNormal Notification
	Type          string `json:"Type"`
	Message       string `json:"Message"`
	Reason        string `json:"Reason"`
	Count         int    `json:"Count"`
	LastTime      string `json:"LastTime"`
	FirstTime     string `json:"FirstTime"`
	IsHandle      bool   `json:"IsHandle"`
	HandleMessage string `json:"HandleMessage"`
	ServiceName   string `json:"ServiceName"`
	TenantName    string `json:"TenantName"`
}
