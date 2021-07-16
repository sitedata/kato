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
	"net/http"
	"os"
	"strconv"

	"github.com/gridworkz/kato/api/handler/cloud"

	"github.com/gridworkz/kato/db"
	"github.com/gridworkz/kato/db/model"

	"github.com/gridworkz/kato/api/middleware"
	httputil "github.com/gridworkz/kato/util/http"
)

//ChargesVerifyController service charges verify
// swagger:operation GET /v2/tenants/{tenant_name}/chargesverify v2 chargesverify
//
// Application expansion resource application interface, public cloud cloud city verification, private cloud does not verify
//
// service charges verify
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
//       "$ref": "#/responses/commandResponse"
//     description: The status code is not 200, indicating that an error occurred during the verification process. Status code 200, msg represents the actual status：success, illegal_quantity, missing_tenant, owned_fee, region_unauthorized, lack_of_memory
func ChargesVerifyController(w http.ResponseWriter, r *http.Request) {
	tenant := r.Context().Value(middleware.ContextKey("tenant")).(*model.Tenants)
	if tenant.EID == "" {
		eid := r.FormValue("eid")
		if eid == "" {
			httputil.ReturnError(r, w, 400, "enterprise id can not found")
			return
		}
		tenant.EID = eid
		db.GetManager().TenantDao().UpdateModel(tenant)
	}
	quantity := r.FormValue("quantity")
	if quantity == "" {
		httputil.ReturnError(r, w, 400, "quantity  can not found")
		return
	}
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		httputil.ReturnError(r, w, 400, "quantity type must be int")
		return
	}

	if publicCloud := os.Getenv("PUBLIC_CLOUD"); publicCloud != "true" {
		err := cloud.PriChargeSverify(tenant, quantityInt)
		if err != nil {
			err.Handle(r, w)
			return
		}
		httputil.ReturnSuccess(r, w, nil)
	} else {
		reason := r.FormValue("reason")
		if err := cloud.PubChargeSverify(tenant, quantityInt, reason); err != nil {
			err.Handle(r, w)
			return
		}
		httputil.ReturnSuccess(r, w, nil)
	}
}
