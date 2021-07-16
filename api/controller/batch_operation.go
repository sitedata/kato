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

package controller

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/middleware"
	"github.com/gridworkz/kato/api/util"

	"github.com/gridworkz/kato/api/model"
	dbmodel "github.com/gridworkz/kato/db/model"
	httputil "github.com/gridworkz/kato/util/http"
)

//BatchOperation batch operation for tenant
//support operation is : start,build,stop,update
func BatchOperation(w http.ResponseWriter, r *http.Request) {
	var build model.BeatchOperationRequestStruct
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &build.Body, nil)
	if !ok {
		logrus.Errorf("start batch operation validate request body failure")
		return
	}

	tenantName := r.Context().Value(middleware.ContextKey("tenant_name")).(string)
	tenantID := r.Context().Value(middleware.ContextKey("tenant_id")).(string)

	// create event for each operation
	eventRe := createBatchEvents(&build, tenantID, build.Operator)

	var re handler.BatchOperationResult
	switch build.Body.Operation {
	case "build":
		for i := range build.Body.BuildInfos {
			build.Body.BuildInfos[i].TenantName = tenantName
		}
		re = handler.GetBatchOperationHandler().Build(build.Body.BuildInfos)
	case "start":
		re = handler.GetBatchOperationHandler().Start(build.Body.StartInfos)
	case "stop":
		re = handler.GetBatchOperationHandler().Stop(build.Body.StopInfos)
	case "upgrade":
		re = handler.GetBatchOperationHandler().Upgrade(build.Body.UpgradeInfos)
	default:
		httputil.ReturnError(r, w, 400, fmt.Sprintf("operation %s do not support batch", build.Body.Operation))
		return
	}

	// append every create event result to re and then return
	re.BatchResult = append(re.BatchResult, eventRe.BatchResult...)
	httputil.ReturnSuccess(r, w, re)
}

func createBatchEvents(build *model.BeatchOperationRequestStruct, tenantID, operator string) (re handler.BatchOperationResult) {
	for i := range build.Body.BuildInfos {
		event, err := util.CreateEvent(dbmodel.TargetTypeService, "build-service", build.Body.BuildInfos[i].ServiceID, tenantID, "", operator, dbmodel.ASYNEVENTTYPE)
		if err != nil {
			re.BatchResult = append(re.BatchResult, handler.OperationResult{ErrMsg: "create event failure", ServiceID: build.Body.BuildInfos[i].ServiceID})
			continue
		}
		build.Body.BuildInfos[i].EventID = event.EventID

	}
	for i := range build.Body.StartInfos {
		event, err := util.CreateEvent(dbmodel.TargetTypeService, "start-service", build.Body.StartInfos[i].ServiceID, tenantID, "", operator, dbmodel.ASYNEVENTTYPE)
		if err != nil {
			re.BatchResult = append(re.BatchResult, handler.OperationResult{ErrMsg: "create event failure", ServiceID: build.Body.StartInfos[i].ServiceID})
			continue
		}
		build.Body.StartInfos[i].EventID = event.EventID
	}
	for i := range build.Body.StopInfos {
		event, err := util.CreateEvent(dbmodel.TargetTypeService, "stop-service", build.Body.StopInfos[i].ServiceID, tenantID, "", operator, dbmodel.ASYNEVENTTYPE)
		if err != nil {
			re.BatchResult = append(re.BatchResult, handler.OperationResult{ErrMsg: "create event failure", ServiceID: build.Body.StopInfos[i].ServiceID})
			continue
		}
		build.Body.StopInfos[i].EventID = event.EventID
	}
	for i := range build.Body.UpgradeInfos {
		event, err := util.CreateEvent(dbmodel.TargetTypeService, "upgrade-service", build.Body.UpgradeInfos[i].ServiceID, tenantID, "", operator, dbmodel.ASYNEVENTTYPE)
		if err != nil {
			re.BatchResult = append(re.BatchResult, handler.OperationResult{ErrMsg: "create event failure", ServiceID: build.Body.UpgradeInfos[i].ServiceID})
			continue
		}
		build.Body.UpgradeInfos[i].EventID = event.EventID
	}

	return
}
