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
	"context"
	"fmt"
	"net/http"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/model"
	ctxutil "github.com/gridworkz/kato/api/util/ctx"
	dbmodel "github.com/gridworkz/kato/db/model"
	httputil "github.com/gridworkz/kato/util/http"
	"github.com/sirupsen/logrus"
)

//BatchOperation batch operation for tenant
//support operation is : start,build,stop,update
func BatchOperation(w http.ResponseWriter, r *http.Request) {
	var build model.BatchOperationReq
	ok := httputil.ValidatorRequestStructAndErrorResponse(r, w, &build.Body, nil)
	if !ok {
		logrus.Errorf("start batch operation validate request body failure")
		return
	}

	tenant := r.Context().Value(ctxutil.ContextKey("tenant")).(*dbmodel.Tenants)

	var batchOpReqs []model.ComponentOpReq
	var f func(ctx context.Context, tenant *dbmodel.Tenants, operator string, batchOpReqs model.BatchOpRequesters) (model.BatchOpResult, error)
	switch build.Body.Operation {
	case "build":
		for _, build := range build.Body.Builds {
			build.TenantName = tenant.Name
			batchOpReqs = append(batchOpReqs, build)
		}
		f = handler.GetBatchOperationHandler().Build
	case "start":
		for _, start := range build.Body.Starts {
			batchOpReqs = append(batchOpReqs, start)
		}
		f = handler.GetBatchOperationHandler().Start
	case "stop":
		for _, stop := range build.Body.Stops {
			batchOpReqs = append(batchOpReqs, stop)
		}
		f = handler.GetBatchOperationHandler().Stop
	case "upgrade":
		for _, upgrade := range build.Body.Upgrades {
			batchOpReqs = append(batchOpReqs, upgrade)
		}
		f = handler.GetBatchOperationHandler().Upgrade
	default:
		httputil.ReturnError(r, w, 400, fmt.Sprintf("operation %s do not support batch", build.Body.Operation))
		return
	}
	if len(batchOpReqs) > 1024 {
		batchOpReqs = batchOpReqs[0:1024]
	}
	res, err := f(r.Context(), tenant, build.Operator, batchOpReqs)
	if err != nil {
		httputil.ReturnBcodeError(r, w, err)
		return
	}

	// append every create event result to re and then return
	httputil.ReturnSuccess(r, w, map[string]interface{}{
		"batch_result": res,
	})
}
