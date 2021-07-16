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

package sources

import (
	"fmt"
	"io"
	"net/http/httputil"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerstrslice "github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/gridworkz/kato/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	dockerDefaultLoggingDriver = "json-file"
)

//DockerService docker cli service
type DockerService struct {
	client * client.Client
	ctx    context.Context
}

//CreateDockerService create docker service
func CreateDockerService(ctx context.Context, client *client.Client) *DockerService {
	return &DockerService{
		ctx:    ctx,
		client: client,
	}
}

// CreateContainer creates a new container in the given PodSandbox
// Docker cannot store the log to an arbitrary location (yet), so we create an
// symlink at LogPath, linking to the actual path of the log.
// TODO: check if the default values returned by the runtime API are ok.
func (ds *DockerService) CreateContainer(config *ContainerConfig) (string, error) {
	if config == nil || config.Metadata == nil {
		return "", fmt.Errorf("container config is nil")
	}
	if config.NetworkConfig == nil {
		return "", fmt.Errorf("container network config is nil")
	}
	createConfig := dockertypes.ContainerCreateConfig{
		Name: config.Metadata.Name,
		Config: &dockercontainer.Config{
			// TODO: set User.
			Entrypoint: dockerstrslice.StrSlice(config.Command),
			Cmd:        dockerstrslice.StrSlice(config.Args),
			Env:        generateEnvList(config.GetEnvs()),
			Image:      config.Image.Image,
			WorkingDir: config.WorkingDir,
			Labels:     config.Labels,
			// Interactive containers:
			OpenStdin:    config.Stdin,
			StdinOnce:    config.StdinOnce,
			Tty:          config.Tty,
			AttachStdin:  config.AttachStdin,
			AttachStdout: config.AttachStdout,
			AttachStderr: config.AttachStderr,
		},
	}
	// Fill the HostConfig.
	hc := &dockercontainer.HostConfig{
		Binds:       generateMountBindings(config.GetMounts()),
		NetworkMode: dockercontainer.NetworkMode(config.NetworkConfig.NetworkMode),
		ExtraHosts:  config.ExtraHosts,
	}
	// Set devices for container.
	devices := make([]dockercontainer.DeviceMapping, len(config.Devices))
	for i, device := range config.Devices {
		devices[i] = dockercontainer.DeviceMapping{
			PathOnHost:        device.HostPath,
			PathInContainer:   device.ContainerPath,
			CgroupPermissions: device.Permissions,
		}
	}
	hc.Resources.Devices = devices

	createConfig.HostConfig = hc
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()
	createResp, err := ds.client.ContainerCreate(ctx, createConfig.Config, createConfig.HostConfig, createConfig.NetworkingConfig, config.Metadata.Name)
	if err != nil {
		return "", err
	}
	return createResp.ID, nil
}

// StartContainer starts the container.
func (ds *DockerService) StartContainer(containerID string) error {
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()
	err := ds.client.ContainerStart(ctx, containerID, dockertypes.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container %q: %v", containerID, err)
	}
	return nil
}

// StopContainer stops a running container with a grace period (i.e., timeout).
func (ds *DockerService) StopContainer(containerID string, timeout int64) error {
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()
	timeoutD := time.Second * time.Duration(timeout)
	return ds.client.ContainerStop(ctx, containerID, &timeoutD)
}

// RemoveContainer removes the container.
// TODO: If a container is still running, should we forcibly remove it?
func (ds *DockerService) RemoveContainer(containerID string) error {
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()
	err := ds.client.ContainerRemove(ctx, containerID, dockertypes.ContainerRemoveOptions{RemoveVolumes: true})
	if err != nil {
		return fmt.Errorf("failed to remove container %q: %v", containerID, err)
	}
	return nil
}

// GetContainerLogPath returns the real
// path where docker stores the container log.
func (ds *DockerService) GetContainerLogPath(containerID string) (string, error) {
	ctx, cancel := context.WithCancel(ds.ctx)
	defer cancel()
	info, err := ds.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %q: %v", containerID, err)
	}
	return info.LogPath, nil
}

//AttachContainer attach container
func (ds *DockerService) AttachContainer(containerID string, attacheStdin, attaStdout, attaStderr bool, inputStream io.ReadCloser, outputStream, errorStream io.Writer, errCh *chan error) (func(), error) {
	options := dockertypes.ContainerAttachOptions{
		Stream: true,
		Stdin:  attacheStdin,
		Stdout: attaStdout,
		Stderr: attaStderr,
	}
	resp, errAttach := ds.client.ContainerAttach(ds.ctx, containerID, options)
	if errAttach != nil && errAttach != httputil.ErrPersistEOF {
		// ContainerAttach returns an ErrPersistEOF (connection closed)
		// means server met an error and put it in Hijacked connection
		// keep the error and read detailed error message from hijacked connection later
		return nil, errAttach
	}
	*errCh = Go(func() error {
		if errHijack := holdHijackedConnection(ds.ctx, inputStream, outputStream, errorStream, resp); errHijack != nil {
			return errHijack
		}
		return errAttach
	})
	return resp.Close, nil
}

// generateMountBindings converts the mount list to a list of strings that
// can be understood by docker.
// Each element in the string is in the form of:
// '<HostPath>:<ContainerPath>', or
// '<HostPath>:<ContainerPath>:ro', if the path is read only, or
// '<HostPath>:<ContainerPath>:Z', if the volume requires SELinux
// relabeling and the pod provides an SELinux label
func generateMountBindings(mounts []*Mount) (result []string) {
	for _, m := range mounts {
		bind := fmt.Sprintf("%s:%s", m.HostPath, m.ContainerPath)
		readOnly := m.Readonly
		if readOnly {
			bind += ":ro"
		}
		// Only request relabeling if the pod provides an SELinux context. If the pod
		// does not provide an SELinux context relabeling will label the volume with
		// the container's randomly allocated MCS label. This would restrict access
		// to the volume to the container which mounts it first.
		if m.SelinuxRelabel {
			if readOnly {
				bind += ",Z"
			} else {
				bind += ":Z"
			}
		}
		result = append(result, bind)
	}
	return
}

// generateEnvList converts KeyValue list to a list of strings, in the form of
// '<key>=<value>', which can be understood by docker.
func generateEnvList(envs []*KeyValue) (result []string) {
	for _, env := range envs {
		result = append(result, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	return
}

// holdHijackedConnection handles copying input to and output from streams to the
// connection
func holdHijackedConnection(ctx context.Context, inputStream io.ReadCloser, outputStream, errorStream io.Writer, resp dockertypes.HijackedResponse) error {
	var err error
	receiveStdout := make(chan error, 1)
	if outputStream != nil || errorStream != nil {
		go func() {
			_, err = util.StdCopy(outputStream, errorStream, resp.Reader)
			logrus.Debug("[hijack] End of stdout")
			receiveStdout <- err
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			io.Copy(resp.Conn, inputStream)
			logrus.Debug("[hijack] End of stdin")
		}

		if err := resp.CloseWrite(); err != nil {
			logrus.Debugf("Couldn't send EOF: %s", err)
		}
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		if err != nil {
			logrus.Debugf("Error receiveStdout: %s", err)
			return err
		}
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			select {
			case err := <-receiveStdout:
				if err != nil {
					logrus.Debugf("Error receiveStdout: %s", err)
					return err
				}
			case <-ctx.Done():
			}
		}
	case <-ctx.Done():
	}
	return nil
}

// Go is a basic promise implementation: it wraps calls a function in a goroutine,
// and returns a channel which will later return the function's return value.
func Go(f func() error) chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- f()
	}()
	return ch
}
