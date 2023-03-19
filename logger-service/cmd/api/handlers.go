package main

import (
	"net/http"

	"github.com/pso-dev/qatiq/backend/logger-service/internal/service"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (a *api) handleLogger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := JSONResponse{
			Message: "Hit the logger service",
		}
		err := a.writeJSON(w, http.StatusOK, response, nil)
		if err != nil {
			a.errorJSON(w, err, http.StatusInternalServerError)
		}
	}
}

func (a *api) handleWriteLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload JSONPayload
		_ = a.readJSON(w, r, &requestPayload)
		item := service.LogEntry{
			Name: requestPayload.Name,
			Data: requestPayload.Data,
		}
		err := a.service.LogItem(item)
		if err != nil {
			a.errorJSON(w, err)
			return
		}
		resp := JSONResponse{
			Message: "logged",
		}
		a.writeJSON(w, http.StatusAccepted, resp)
	}
}
