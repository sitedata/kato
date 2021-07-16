package handler

import (
	"github.com/gridworkz/kato/worker/client"
	"github.com/gridworkz/kato/worker/server/pb"
)

// PodHandler defines handler methods about k8s pods.
type PodHandler interface {
	PodDetail(serviceID, podName string) (*pb.PodDetail, error)
}

// NewPodHandler creates a new PodHandler.
func NewPodHandler(statusCli *client.AppRuntimeSyncClient) PodHandler {
	return &PodAction{
		statusCli: statusCli,
	}
}
