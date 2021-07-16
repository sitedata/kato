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

package license

import (
	"net/http"

	"github.com/gridworkz/kato/api/controller"

	"github.com/go-chi/chi"

	"github.com/sirupsen/logrus"
)

//License license struct
type License struct{}

//Routes routes
func Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", controller.GetLicenseManager().AnalystLicense)
	return r
}

//CheckLicense - verification
func CheckLicense(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		licenseCome := r.Header.Get("License")
		logrus.Debug("License is :" + licenseCome)
		next.ServeHTTP(w, r)
		return
		// util.CloseRequest(r)
		// w.WriteHeader(http.StatusUnauthorized)
	}
	return http.HandlerFunc(fn)
}
