// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocondor/core/cache"
	"github.com/gocondor/core/database"
	"github.com/gocondor/core/middlewares"
	"github.com/julienschmidt/httprouter"
)

var filePath string

// App struct
type App struct {
	Features *Features
}

func New() *App {
	return &App{
		Features: &Features{},
	}
}

func (app *App) SetEnv(env map[string]string) {
	for key, val := range env {
		os.Setenv(strings.TrimSpace(key), strings.TrimSpace(val))
	}
}

func (app *App) SetLogsFilePath(f string) {
	filePath = f
}

func (app *App) GetLogsFile() *os.File {
	return logsFile
}

func (app *App) Bootstrap() {
	middlewares.New()
	NewRouter()
	if app.Features.Database == true {
		database.New()
	}
	if app.Features.Cache == true {
		cache.New(app.Features.Cache)
	}
}

func (app *App) Run(portNumber string, router *httprouter.Router) {
	router = app.RegisterRoutes(ResolveRouter().GetRoutes(), router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", portNumber), router))
}

func (app *App) SetEnabledFeatures(features *Features) {
	app.Features = features
}

func (app *App) RegisterRoutes(routes []Route, router *httprouter.Router) *httprouter.Router {
	for _, route := range routes {
		switch route.Method {
		case "get":
			router.GET(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "post":
			router.POST(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "delete":
			router.DELETE(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "patch":
			router.PATCH(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "put":
			router.PUT(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "options":
			router.OPTIONS(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		case "head":
			router.HEAD(route.Path, makeHTTPRouterHandlerFunc(route.Handler))
		}
	}

	return router
}

func makeHTTPRouterHandlerFunc(h Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := &Context{
			Request: &Request{
				httpRequest:    r,
				httpPathParams: ps,
			},
			ResponseBag: &Response{
				headers:        []header{},
				textBody:       "",
				jsonBody:       "",
				responseWriter: w,
			},
			Logger: NewLogger(filePath),
		}
		response := h(ctx)
		for _, header := range response.headers {
			w.Header().Add(header.key, header.val)
		}
		defer logsFile.Close()
		if response.getTextBody() != "" {
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			response.responseWriter.Write([]byte(response.getTextBody()))
		}
		if response.getJsonBody() != "" {
			w.Header().Add("Content-Type", "application/json")
			response.responseWriter.Write([]byte(response.getJsonBody()))
		}
	}
}
