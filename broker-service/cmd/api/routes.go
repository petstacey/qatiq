package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (b *broker) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/", b.handleBroker())
	router.HandlerFunc(http.MethodGet, "/v1/auth", b.handleAuth())
	router.HandlerFunc(http.MethodPost, "/v1/handle", b.handleRequest())
	return b.enableCORS(router)
}
