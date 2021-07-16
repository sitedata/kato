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

package middleware

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gridworkz/kato/api/handler"
	"github.com/gridworkz/kato/api/util"
)

//Token - simple token verification
func Token(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/docs") {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			auths := strings.SplitN(auth, " ", 2)
			if len(auths) != 2 {
				return
			}
			authMethod := auths[0]
			authB64 := auths[1]
			switch authMethod {
			case "Basic":
				authstr, err := base64.StdEncoding.DecodeString(authB64)
				if err != nil {
					io.WriteString(w, "Unauthorized!\n")
					return
				}
				userPwd := strings.SplitN(string(authstr), ":", 2)
				if len(userPwd) != 2 {
					io.WriteString(w, "Unauthorized!\n")
					return
				}
				username := userPwd[0]
				password := userPwd[1]
				if username == "gridworkz" && password == "gridworkz-api-test" {
					next.ServeHTTP(w, r)
					return
				}
			default:
				io.WriteString(w, "Unauthorized!\n")
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := os.Getenv("TOKEN")
		t := r.Header.Get("Authorization")
		if tt := strings.Split(t, " "); len(tt) == 2 {
			if tt[1] == token {
				next.ServeHTTP(w, r)
				return
			}
		}
		util.CloseRequest(r)
		w.WriteHeader(http.StatusUnauthorized)
	}
	return http.HandlerFunc(fn)
}

//FullToken token api check
func FullToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/docs") {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			auths := strings.SplitN(auth, " ", 2)
			if len(auths) != 2 {
				return
			}
			authMethod := auths[0]
			authB64 := auths[1]
			switch authMethod {
			case "Basic":
				authstr, err := base64.StdEncoding.DecodeString(authB64)
				if err != nil {
					io.WriteString(w, "Unauthorized!\n")
					return
				}
				userPwd := strings.SplitN(string(authstr), ":", 2)
				if len(userPwd) != 2 {
					io.WriteString(w, "Unauthorized!\n")
					return
				}
				username := userPwd[0]
				password := userPwd[1]
				if username == "gridworkz" && password == "gridworkz-api-test" {
					next.ServeHTTP(w, r)
					return
				}
			default:
				io.WriteString(w, "Unauthorized!\n")
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//logrus.Debugf("request uri is %s", r.RequestURI)
		t := r.Header.Get("Authorization")
		if tt := strings.Split(t, " "); len(tt) == 2 {
			if handler.GetTokenIdenHandler().CheckToken(tt[1], r.RequestURI) {
				next.ServeHTTP(w, r)
				return
			}
		}
		util.CloseRequest(r)
		w.WriteHeader(http.StatusUnauthorized)
	}
	return http.HandlerFunc(fn)
}
