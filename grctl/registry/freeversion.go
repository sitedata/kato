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
	"github.com/gridworkz/kato/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var _ FreeImager = &FreeVersion{}

// FreeVersion is resposible for listing the free images belong to free component versions.
type FreeVersion struct {
	reg *registry.Registry
}

// NewFreeVersion creates a free version.
func NewFreeVersion(reg *registry.Registry) *FreeVersion {
	return &FreeVersion{
		reg: reg,
	}
}

// List return a list of free images belong to free component versions.
func (f *FreeVersion) List() ([]*FreeImage, error) {
	// list components
	components, err := db.GetManager().TenantServiceDao().GetAllServicesID()
	if err != nil {
		return nil, err
	}

	var images []*FreeImage
	for _, cpt := range components {
		// list free tags
		freeTags, err := f.listFreeTags(cpt.ServiceID)
		if err != nil {
			logrus.Warningf("list free tags for repository %s: %v", cpt.ServiceID, err)
			continue
		}

		// component.ServiceID is the repository of image
		freeImages, err := f.listFreeImages(cpt.ServiceID, freeTags)
		if err != nil {
			logrus.Warningf("list digests for repository %s: %v", cpt.ServiceID, err)
			continue
		}
		images = append(images, freeImages...)
	}

	return images, nil
}

func (f *FreeVersion) listFreeTags(serviceID string) ([]string, error) {
	// all tags
	// serviceID is the repository of image
	tags, err := f.reg.Tags(serviceID)
	if err != nil {
		if errors.Is(err, registry.ErrRepositoryNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// versions being used.
	versions, err := f.listVersions(serviceID)
	if err != nil {
		return nil, err
	}

	var freeTags []string
	for _, tag := range tags {
		_, ok := versions[tag]
		if !ok {
			freeTags = append(freeTags, tag)
		}
	}

	return freeTags, nil
}

func (f *FreeVersion) listVersions(serviceID string) (map[string]struct{}, error) {
	// tags being used
	rawVersions, err := db.GetManager().VersionInfoDao().ListByServiceIDStatus(serviceID, util.Bool(true))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]struct{}{}, nil
		}
		return nil, err
	}
	// make a map of versions
	versions := make(map[string]struct{})
	for _, version := range rawVersions {
		versions[version.BuildVersion] = struct{}{}
	}
	return versions, nil
}

func (f *FreeVersion) listFreeImages(repository string, tags []string) ([]*FreeImage, error) {
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
