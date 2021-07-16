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

package registry

import (
	"github.com/gridworkz/kato/builder/sources/registry"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Cleaner is responsible for cleaning up the free images in registry.
type Cleaner struct {
	reg          *registry.Registry
	freeImageres map[string]FreeImager
}

// NewRegistryCleaner creates a new Cleaner.
func NewRegistryCleaner(url, username, password string) (*Cleaner, error) {
	reg, err := registry.NewInsecure(url, username, password)

	freeImageres := NewFreeImageres(reg)

	return &Cleaner{
		reg:          reg,
		freeImageres: freeImageres,
	}, err
}

// Cleanup cleans up the free image in the registry.
func (r *Cleaner) Cleanup() {
	logrus.Info("Start cleaning up the free images. Please be patient.")
	logrus.Info("The clean up time will be affected by the number of free images and the network environment.")

	// list images needed to be cleaned up
	freeImages := r.ListFreeImages()
	if len(freeImages) == 0 {
		logrus.Info("Free images not Found")
		return
	}
	logrus.Infof("Found %d free images", len(freeImages))

	// delete images
	if err := r.DeleteImages(freeImages); err != nil {
		if errors.Is(err, registry.ErrOperationIsUnsupported) {
			logrus.Warningf(`The operation image deletion is unsupported.
You can try to add REGISTRY_STORAGE_DELETE_ENABLED=true when start the registry.
More detail: https://docs.docker.com/registry/configuration/#list-of-configuration-options.
				`)
			return
		}
		logrus.Warningf("delete images: %v", err)
		return
	}

	logrus.Infof(`you have to exec the command below in the registry container to remove blobs from the filesystem:
	/bin/registry garbage-collect /etc/docker/registry/config.yml
More Detail: https://docs.docker.com/registry/garbage-collection/#run-garbage-collection.`)
}

// ListFreeImages return a list of free images needed to be cleaned up.
func (r *Cleaner) ListFreeImages() []*FreeImage {
	var freeImages []*FreeImage
	for name, freeImager := range r.freeImageres {
		images, err := freeImager.List()
		if err != nil {
			logrus.Warningf("list free images for %s", name)
		}
		logrus.Infof("Found %d free images from %s", len(images), name)
		freeImages = append(freeImages, images...)
	}

	// deduplicate
	var result []*FreeImage
	m := make(map[string]struct{})
	for _, fi := range freeImages {
		fi := fi
		key := fi.Key()
		_, ok := m[key]
		if ok {
			continue
		}
		m[key] = struct{}{}
		result = append(result, fi)
	}

	return result
}

// DeleteImages deletes images.
func (r *Cleaner) DeleteImages(freeImages []*FreeImage) error {
	for _, image := range freeImages {
		if err := r.deleteManifest(image.Repository, image.Digest); err != nil {
			if errors.Is(err, registry.ErrOperationIsUnsupported) {
				return err
			}
			logrus.Warningf("delete manifest %s/%s: %v", image.Repository, image.Digest, err)
			continue
		}
		log := logrus.WithField("Type", image.Type).
			WithField("Component ID", image.Repository).
			WithField("Build Version", image.Tag).
			WithField("Digest", image.Digest)
		log.Infof("image %s/%s deleted", image.Repository, image.Tag)
	}
	return nil
}

func (r *Cleaner) deleteManifest(repository, dig string) error {
	return r.reg.DeleteManifest(repository, digest.Digest(dig))
}
