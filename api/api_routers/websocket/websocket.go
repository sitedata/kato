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

package websocket

import (
	"github.com/gridworkz/kato/api/controller"

	"github.com/go-chi/chi"
)

//Routes routes
func Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/docker_console", controller.GetDockerConsole().Get)
	r.Get("/docker_log", controller.GetDockerLog().Get)
	r.Get("/monitor_message", controller.GetMonitorMessage().Get)
	r.Get("/new_monitor_message", controller.GetMonitorMessage().Get)
	r.Get("/event_log", controller.GetEventLog().Get)
	r.Get("/services/{serviceID}/pubsub", controller.GetPubSubControll().Get)
	return r
}

//LogRoutes - Log download routing
func LogRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{gid}/{filename}", controller.GetLogFile().Get)
	r.Get("/install_log/{filename}", controller.GetLogFile().GetInstallLog)
	return r
}

//AppRoutes - Application export package download route
func AppRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/download/{format}/{fileName}", controller.GetManager().Download)
	r.Post("/upload/{eventID}", controller.GetManager().Upload)
	r.Options("/upload/{eventID}", controller.GetManager().Upload)
	return r
}
