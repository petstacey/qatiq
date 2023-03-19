package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *api) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/", a.handleLogger())
	router.HandlerFunc(http.MethodPost, "/v1/log", a.handleWriteLog())
	return a.enableCORS(router)
}
