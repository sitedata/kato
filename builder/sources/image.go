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

package sources

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/gridworkz/kato/builder/model"
	"github.com/gridworkz/kato/event"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ErrorNoAuth error no auth
var ErrorNoAuth = fmt.Errorf("pull image require docker login")

// ErrorNoImage error no image
var ErrorNoImage = fmt.Errorf("image not exist")

// ImagePull pull docker image
// timeout minutes of the unit
func ImagePull(dockerCli *client.Client, image string, username, password string, logger event.Logger, timeout int) (*types.ImageInspect, error) {
	if logger != nil {
		logger.Info(fmt.Sprintf("start get image:%s", image), map[string]string{"step": "pullimage"})
	}
	var pullipo types.ImagePullOptions
	if username != "" && password != "" {
		auth, err := EncodeAuthToBase64(types.AuthConfig{Username: username, Password: password})
		if err != nil {
			logrus.Errorf("make auth base63 push image error: %s", err.Error())
			logger.Error(fmt.Sprintf("Failed to generate a Token to get the image"), map[string]string{"step": "builder-exector", "status": "failure"})
			return nil, err
		}
		pullipo = types.ImagePullOptions{
			RegistryAuth: auth,
		}
	} else {
		pullipo = types.ImagePullOptions{}
	}
	rf, err := reference.ParseAnyReference(image)
	if err != nil {
		logrus.Errorf("reference image error: %s", err.Error())
		return nil, err
	}
	//At least one minute
	if timeout < 1 {
		timeout = 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()
	//TODO: Use 1.12 version api bug “repository name must be canonical”，use rf.String() to complete the mirror address
	readcloser, err := dockerCli.ImagePull(ctx, rf.String(), pullipo)
	if err != nil {
		logrus.Debugf("image name: %s readcloser error: %v", image, err.Error())
		if strings.HasSuffix(err.Error(), "does not exist or no pull access") {
			if logger != nil {
				logger.Error(fmt.Sprintf("image: %s does not exist or is not available", image), map[string]string{"step": "pullimage", "status": "failure"})
			}
			return nil, fmt.Errorf("Image (%s) does not exist or no pull access", image)
		}
		return nil, err
	}
	defer readcloser.Close()
	dec := json.NewDecoder(readcloser)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		var jm JSONMessage
		if err := dec.Decode(&jm); err != nil {
			if err == io.EOF {
				break
			}
			logrus.Debugf("error decoding jm(JSONMessage): %v", err)
			return nil, err
		}
		if jm.Error != nil {
			logrus.Debugf("error pulling image: %v", jm.Error)
			return nil, jm.Error
		}
		if logger != nil {
			logger.Debug(fmt.Sprintf(jm.JSONString()), map[string]string{"step": "progress"})
		} else {
			logrus.Debug(jm.JSONString())
		}
	}
	logger.Debug("Get the image information and its raw representation", map[string]string{"step": "progress"})
	ins, _, err := dockerCli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		logger.Debug("Fail to get the image information and its raw representation", map[string]string{"step": "progress"})
		return nil, err
	}
	logger.Info(fmt.Sprintf("Success Pull Image：%s", image), map[string]string{"step": "pullimage"})
	return &ins, nil
}

//ImageTag change docker image tag
func ImageTag(dockerCli *client.Client, source, target string, logger event.Logger, timeout int) error {
	if logger != nil {
		logger.Info(fmt.Sprintf("change image tag：%s -> %s", source, target), map[string]string{"step": "changetag"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()
	err := dockerCli.ImageTag(ctx, source, target)
	if err != nil {
		logrus.Debugf("image tag err: %s", err.Error())
		return err
	}
	logger.Info("change image tag success", map[string]string{"step": "changetag"})
	return nil
}

//ImageNameHandle
func ImageNameHandle(imageName string) *model.ImageName {
	var i model.ImageName
	if strings.Contains(imageName, "/") {
		mm := strings.Split(imageName, "/")
		i.Host = mm[0]
		names := strings.Join(mm[1:], "/")
		if strings.Contains(names, ":") {
			nn := strings.Split(names, ":")
			i.Name = nn[0]
			i.Tag = nn[1]
		} else {
			i.Name = names
			i.Tag = "latest"
		}
	} else {
		if strings.Contains(imageName, ":") {
			nn := strings.Split(imageName, ":")
			i.Name = nn[0]
			i.Tag = nn[1]
		} else {
			i.Name = imageName
			i.Tag = "latest"
		}
	}
	return &i
}

// ImageNameWithNamespaceHandle if have namespace, will parse namespace
func ImageNameWithNamespaceHandle(imageName string) *model.ImageName {
	var i model.ImageName
	if strings.Contains(imageName, "/") {
		mm := strings.Split(imageName, "/")
		i.Host = mm[0]
		names := strings.Join(mm[1:], "/")
		if len(mm) >= 3 {
			i.Namespace = mm[1]
			names = strings.Join(mm[2:], "/")
		}
		if strings.Contains(names, ":") {
			nn := strings.Split(names, ":")
			i.Name = nn[0]
			i.Tag = nn[1]
		} else {
			i.Name = names
			i.Tag = "latest"
		}
	} else {
		if strings.Contains(imageName, ":") {
			nn := strings.Split(imageName, ":")
			i.Name = nn[0]
			i.Tag = nn[1]
		} else {
			i.Name = imageName
			i.Tag = "latest"
		}
	}
	return &i
}

// GenSaveImageName generates the final name of the image, which is the name of
// the image in the exported tar package.
func GenSaveImageName(name string) string {
	imageName := ImageNameWithNamespaceHandle(name)
	return fmt.Sprintf("%s:%s", imageName.Name, imageName.Tag)
}

// ImagePush push image to registry
// timeout minutes of the unit
func ImagePush(dockerCli *client.Client, image, user, pass string, logger event.Logger, timeout int) error {
	if logger != nil {
		logger.Info(fmt.Sprintf("start push image：%s", image), map[string]string{"step": "pushimage"})
	}
	if timeout < 1 {
		timeout = 1
	}
	if user == "" {
		user = os.Getenv("LOCAL_HUB_USER")
	}
	if pass == "" {
		pass = os.Getenv("LOCAL_HUB_PASS")
	}
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return err
	}
	var opts types.ImagePushOptions
	pushauth, err := EncodeAuthToBase64(types.AuthConfig{
		Username:      user,
		Password:      pass,
		ServerAddress: reference.Domain(ref),
	})
	if err != nil {
		logrus.Errorf("make auth base63 push image error: %s", err.Error())
		if logger != nil {
			logger.Error(fmt.Sprintf("Failed to generate a token to get the image"), map[string]string{"step": "builder-exector", "status": "failure"})
		}
		return err
	}
	opts.RegistryAuth = pushauth
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
	defer cancel()
	readcloser, err := dockerCli.ImagePush(ctx, image, opts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			if logger != nil {
				logger.Error(fmt.Sprintf("image %s does not exist, cannot be pushed", image), map[string]string{"step": "pushimage", "status": "failure"})
			}
			return fmt.Errorf("Image(%s) does not exist", image)
		}
		return err
	}
	if readcloser != nil {
		defer readcloser.Close()
		dec := json.NewDecoder(readcloser)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			var jm JSONMessage
			if err := dec.Decode(&jm); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if jm.Error != nil {
				return jm.Error
			}
			logger.Debug(jm.JSONString(), map[string]string{"step": "progress"})
		}
	}
	logger.Info(fmt.Sprintf("success push image：%s", image), map[string]string{"step": "pushimage"})
	return nil
}

// TrustedImagePush push image to trusted registry
func TrustedImagePush(dockerCli *client.Client, image, user, pass string, logger event.Logger, timeout int) error {
	if err := CheckTrustedRepositories(image, user, pass); err != nil {
		return err
	}
	return ImagePush(dockerCli, image, user, pass, logger, timeout)
}

// CheckTrustedRepositories check Repositories is exist, if not create it.
func CheckTrustedRepositories(image, user, pass string) error {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return err
	}
	var server string
	if reference.IsNameOnly(ref) {
		server = "docker.io"
	} else {
		server = reference.Domain(ref)
	}
	cli, err := createTrustedRegistryClient(server, user, pass)
	if err != nil {
		return err
	}
	var namespace, repName string
	infos := strings.Split(reference.TrimNamed(ref).String(), "/")
	if len(infos) == 3 && infos[0] == server {
		namespace = infos[1]
		repName = infos[2]
	}
	if len(infos) == 2 {
		namespace = infos[0]
		repName = infos[1]
	}
	_, err = cli.GetRepository(namespace, repName)
	if err != nil {
		if err.Error() == "resource does not exist" {
			rep := Repostory{
				Name:             repName,
				ShortDescription: image, // The maximum length is 140
				LongDescription:  fmt.Sprintf("push image for %s", image),
				Visibility:       "private",
			}
			if len(rep.ShortDescription) > 140 {
				rep.ShortDescription = rep.ShortDescription[0:140]
			}
			err := cli.CreateRepository(namespace, &rep)
			if err != nil {
				return fmt.Errorf("create repostory error,%s", err.Error())
			}
			return nil
		}
		return fmt.Errorf("get repostory error,%s", err.Error())
	}
	return err
}

// EncodeAuthToBase64 serializes the auth configuration as JSON base64 payload
func EncodeAuthToBase64(authConfig types.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

// ImageBuild
func ImageBuild(dockerCli *client.Client, contextDir string, options types.ImageBuildOptions, logger event.Logger, timeout int) (string, error) {
	var ctx context.Context
	if timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute*time.Duration(timeout))
		defer cancel()
	} else {
		ctx = context.Background()
	}
	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		Compression:     archive.Uncompressed,
		ExcludePatterns: []string{""},
		IncludeFiles:    []string{"."},
	})
	if err != nil {
		return "", err
	}
	rc, err := dockerCli.ImageBuild(ctx, buildCtx, options)
	if err != nil {
		return "", err
	}
	var out io.Writer
	if logger != nil {
		out = logger.GetWriter("build-progress", "info")
	} else {
		out, _ = os.OpenFile("/tmp/build.log", os.O_RDWR|os.O_CREATE, 0755)
	}
	var imageID string
	err = jsonmessage.DisplayJSONMessagesStream(rc.Body, out, 0, true, func(msg jsonmessage.JSONMessage) {
		var r types.BuildResult
		imageID = r.ID
	})
	if err != nil {
		logrus.Errorf("read build log failure %s", err.Error())
		return "", err
	}
	return imageID, nil
}

//ImageInspectWithRaw get image inspect
func ImageInspectWithRaw(dockerCli *client.Client, image string) (*types.ImageInspect, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ins, _, err := dockerCli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return nil, err
	}
	return &ins, nil
}

//ImageSave save image to tar file
// destination destination file name eg. /tmp/xxx.tar
func ImageSave(dockerCli *client.Client, image, destination string, logger event.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rc, err := dockerCli.ImageSave(ctx, []string{image})
	if err != nil {
		return err
	}
	defer rc.Close()
	return CopyToFile(destination, rc)
}

//MultiImageSave save multi image to tar file
// destination destination file name eg. /tmp/xxx.tar
func MultiImageSave(ctx context.Context, dockerCli *client.Client, destination string, logger event.Logger, images ...string) error {
	rc, err := dockerCli.ImageSave(ctx, images)
	if err != nil {
		return err
	}
	defer rc.Close()
	return CopyToFile(destination, rc)
}

//ImageLoad load image from  tar file
// destination destination file name eg. /tmp/xxx.tar
func ImageLoad(dockerCli *client.Client, tarFile string, logger event.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader, err := os.OpenFile(tarFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer reader.Close()

	rc, err := dockerCli.ImageLoad(ctx, reader, false)
	if err != nil {
		return err
	}
	if rc.Body != nil {
		defer rc.Body.Close()
		dec := json.NewDecoder(rc.Body)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			var jm JSONMessage
			if err := dec.Decode(&jm); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			if jm.Error != nil {
				return jm.Error
			}
			logger.Info(jm.JSONString(), map[string]string{"step": "build-progress"})
		}
	}

	return nil
}

//ImageImport save image to tar file
// source source file name eg. /tmp/xxx.tar
func ImageImport(dockerCli *client.Client, image, source string, logger event.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	isource := types.ImageImportSource{
		Source:     file,
		SourceName: "-",
	}

	options := types.ImageImportOptions{}

	readcloser, err := dockerCli.ImageImport(ctx, isource, image, options)
	if err != nil {
		return err
	}
	if readcloser != nil {
		defer readcloser.Close()
		r := bufio.NewReader(readcloser)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			if line, _, err := r.ReadLine(); err == nil {
				if logger != nil {
					logger.Debug(string(line), map[string]string{"step": "progress"})
				} else {
					fmt.Println(string(line))
				}
			} else {
				if err.Error() == "EOF" {
					return nil
				}
				return err
			}
		}
	}
	return nil
}

// CopyToFile writes the content of the reader to the specified file
func CopyToFile(outfile string, r io.Reader) error {
	// We use sequential file access here to avoid depleting the standby list
	// on Windows. On Linux, this is a call directly to ioutil.TempFile
	tmpFile, err := os.OpenFile(path.Join(filepath.Dir(outfile), ".docker_temp_"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer tmpFile.Close()
	tmpPath := tmpFile.Name()
	_, err = io.Copy(tmpFile, r)
	if err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err = os.Rename(tmpPath, outfile); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}

//ImageRemove remove image
func ImageRemove(dockerCli *client.Client, image string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err := dockerCli.ImageRemove(ctx, image, types.ImageRemoveOptions{Force: true})
	return err
}
