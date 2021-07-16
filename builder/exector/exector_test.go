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

package exector

import (
	"context"
	"encoding/json"
	"runtime"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/docker/docker/client"
	"k8s.io/client-go/kubernetes"

	"github.com/gridworkz/kato/builder/parser/code"
	"github.com/gridworkz/kato/cmd/builder/option"
	"github.com/gridworkz/kato/event"
	"github.com/gridworkz/kato/mq/api/grpc/pb"

	mqclient "github.com/gridworkz/kato/mq/client"
	etcdutil "github.com/gridworkz/kato/util/etcd"
	k8sutil "github.com/gridworkz/kato/util/k8s"
)

func Test_exectorManager_buildFromSourceCode(t *testing.T) {
	conf := option.Config{
		EtcdEndPoints:       []string{"192.168.2.203:2379"},
		MQAPI:               "192.168.2.203:6300",
		EventLogServers:     []string{"192.168.2.203:6366"},
		RbdRepoName:         "rbd-dns",
		RbdNamespace:        "rbd-system",
		MysqlConnectionInfo: "EeM2oc:lee7OhQu@tcp(192.168.2.203:3306)/region",
	}
	etcdArgs := etcdutil.ClientArgs{Endpoints: conf.EtcdEndPoints}
	event.NewManager(event.EventConfig{
		EventLogServers: conf.EventLogServers,
		DiscoverArgs:    &etcdArgs,
	})
	restConfig, err := k8sutil.NewRestConfig("/Users/gridworkz/Documents/company/gridworkz/admin.kubeconfig")
	if err != nil {
		t.Fatal(err)
	}
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.EtcdEndPoints,
		DialTimeout: 10 * time.Second,
	})
	var maxConcurrentTask int
	if conf.MaxTasks == 0 {
		maxConcurrentTask = runtime.NumCPU() * 2
	} else {
		maxConcurrentTask = conf.MaxTasks
	}
	mqClient, err := mqclient.NewMqClient(&etcdArgs, conf.MQAPI)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	e := &exectorManager{
		DockerClient:      dockerClient,
		KubeClient:        kubeClient,
		EtcdCli:           etcdCli,
		tasks:             make(chan *pb.TaskMessage, maxConcurrentTask),
		maxConcurrentTask: maxConcurrentTask,
		mqClient:          mqClient,
		ctx:               ctx,
		cancel:            cancel,
		cfg:               conf,
	}
	taskBodym := make(map[string]interface{})
	taskBodym["repo_url"] = "https://github.com/gridworkz/java-maven-demo.git"
	taskBodym["branch"] = "master"
	taskBodym["tenant_id"] = "5d7bd886e6dc4425bb6c2ac5fc9fa593"
	taskBodym["service_id"] = "4eaa41ccf145b8e43a6aeb1a5efeab53"
	taskBodym["deploy_version"] = "20200115193617"
	taskBodym["lang"] = code.JavaMaven
	taskBodym["event_id"] = "0000"
	taskBodym["envs"] = map[string]string{}

	taskBody, _ := json.Marshal(taskBodym)
	task := pb.TaskMessage{
		TaskType: "build_from_source_code",
		TaskBody: taskBody,
	}
	i := NewSouceCodeBuildItem(task.TaskBody)
	if err := i.Run(30 * time.Second); err != nil {
		t.Fatal(err)
	}
	e.buildFromSourceCode(&task)
}
