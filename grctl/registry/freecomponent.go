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
	"errors"

	"github.com/gridworkz/kato/builder/sources/registry"
	"github.com/gridworkz/kato/db"
	"github.com/sirupsen/logrus"
)

var _ FreeImager = &FreeComponent{}

// FreeComponent is resposible for listing the free images belong to free components.
type FreeComponent struct {
	reg *registry.Registry
}

// NewFreeComponent creates a new FreeComponent.
func NewFreeComponent(reg *registry.Registry) *FreeComponent {
	return &FreeComponent{
		reg: reg,
	}
}

// List return a list of free images belong to free components.
func (f *FreeComponent) List() ([]*FreeImage, error) {
	// list free components
	components, err := db.GetManager().TenantServiceDeleteDao().List()
	if err != nil {
		return nil, err
	}

	var images []*FreeImage
	for _, cpt := range components {
		// component.ServiceID is the repository of image
		freeImages, err := f.listFreeImages(cpt.ServiceID)
		if err != nil {
			logrus.Warningf("list free images: %v", err)
			continue
		}
		images = append(images, freeImages...)
	}

	return images, nil
}

func (f *FreeComponent) listFreeImages(repository string) ([]*FreeImage, error) {
	// list tags, then list digest for every tag
	tags, err := f.reg.Tags(repository)
	if err != nil {
		if errors.Is(err, registry.ErrRepositoryNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var images []*FreeImage
	for _, tag := range tags {
		digest, err := f.reg.ManifestDigestV2(repository, tag)
		if err != nil {
			logrus.Warningf("get digest for manifest %s/%s: %v", repository, tag, err)
			continue
		}
		images = append(images, &FreeImage{
			Repository: repository,
			Digest:     digest.String(),
			Tag:        tag,
			Type:       string(FreeImageTypeFreeVersion),
		})
	}
	return images, nil
}
