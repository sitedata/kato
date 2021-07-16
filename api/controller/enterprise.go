package controller

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/gridworkz/kato/api/handler"

	httputil "github.com/gridworkz/kato/util/http"
)

//GetRunningServices list all running service ids
func GetRunningServices(w http.ResponseWriter, r *http.Request) {
	enterpriseID := chi.URLParam(r, "enterprise_id")
	runningList, err := handler.GetServiceManager().GetEnterpriseRunningServices(enterpriseID)
	if err != nil {
		err.Handle(r, w)
		return
	}
	httputil.ReturnNoFomart(r, w, 200, map[string]interface{}{"service_ids": runningList})
}
