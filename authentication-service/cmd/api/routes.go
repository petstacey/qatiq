package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *api) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/", a.handleAuth())
	router.HandlerFunc(http.MethodGet, "/v1/users/", a.handleGetUsers())
	router.HandlerFunc(http.MethodPost, "/v1/authenticate/", a.handleAuthenticate())
	return a.enableCORS(router)
}
