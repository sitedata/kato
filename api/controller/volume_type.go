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

package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gridworkz/kato/api/handler"
	api_model "github.com/gridworkz/kato/api/model"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/jinzhu/gorm"
)

// VolumeOptions list volume option
func VolumeOptions(w http.ResponseWriter, r *http.Request) {
	volumetypeOptions, err := handler.GetVolumeTypeHandler().GetAllVolumeTypes()
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, volumetypeOptions)
}

// ListVolumeType list volume type list
func ListVolumeType(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/volume-options v2 volumeOptions TODO delete it
	//
	// Query the list of available storage-driven models
	//
	// get volume-options
	//
	// ---
	// consumes:
	// - application/json
	// - application/x-protobuf
	//
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//     description: Unified return format
	pageStr := strings.TrimSpace(chi.URLParam(r, "page"))
	pageSizeCul := strings.TrimSpace(chi.URLParam(r, "pageSize"))
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		httputil.ReturnError(r, w, 400, fmt.Sprintf("bad request, %v", err))
		return
	}
	pageSize, err := strconv.Atoi(pageSizeCul)
	if err != nil {
		httputil.ReturnError(r, w, 400, fmt.Sprintf("bad request, %v", err))
		return
	}
	volumetypeOptions, er := handler.GetVolumeTypeHandler().GetAllVolumeTypes()
	volumetypePageOptions, err := handler.GetVolumeTypeHandler().GetAllVolumeTypesByPage(page, pageSize)
	if err != nil || er != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}

	// httputil.ReturnSuccess(r, w, volumetypeOptions)
	httputil.ReturnSuccess(r, w, map[string]interface{}{
		"data":      volumetypePageOptions,
		"page":      page,
		"page_size": pageSize,
		"count":     len(volumetypeOptions),
	})
}

// VolumeSetVar set volume option
func VolumeSetVar(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/volume-options v2 volumeOptions
	//
	// Create a list of available storage-driven models
	//
	// get volume-options
	//
	// ---
	// consumes:
	// - application/json
	// - application/x-protobuf
	//
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//     description: Unified return format
	volumeType := api_model.VolumeTypeStruct{}
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &volumeType, nil); !ok {
		return
	}
	err := handler.GetVolumeTypeHandler().SetVolumeType(&volumeType)
	if err != nil {
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

// DeleteVolumeType delete volume option
func DeleteVolumeType(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/volume-options v2 volumeOptions
	//
	// Delete available storage drive model
	//
	// get volume-options
	//
	// ---
	// consumes:
	// - application/json
	// - application/x-protobuf
	//
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//     description: Unified return format
	volumeType := chi.URLParam(r, "volume_type")
	err := handler.GetVolumeTypeHandler().DeleteVolumeType(volumeType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httputil.ReturnError(r, w, 404, "not found")
			return
		}
		httputil.ReturnError(r, w, 500, err.Error())
		return
	}
	httputil.ReturnSuccess(r, w, nil)
}

// UpdateVolumeType delete volume option
func UpdateVolumeType(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /v2/volume-options v2 volumeOptions
	//
	// Available update storage-driven model
	//
	// get volume-options
	//
	// ---
	// consumes:
	// - application/json
	// - application/x-protobuf
	//
	// produces:
	// - application/json
	// - application/xml
	//
	// responses:
	//   default:
	//     schema:
	//     description: Unified return format
	volumeTypeID := chi.URLParam(r, "volume_type")
	volumeType := api_model.VolumeTypeStruct{}
	if ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &volumeType, nil); !ok {
		return
	}
	volume, err := handler.GetVolumeTypeHandler().GetVolumeTypeByType(volumeTypeID)
	if err == nil {
		if volume == nil {
			httputil.ReturnError(r, w, 404, "not found")
			return
		}
		if updateErr := handler.GetVolumeTypeHandler().UpdateVolumeType(volume, &volumeType); updateErr != nil {
			httputil.ReturnError(r, w, 500, err.Error())
		}
		httputil.ReturnSuccess(r, w, nil)
	}
	httputil.ReturnError(r, w, 500, err.Error())
}
