package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (a *api) handleAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := JSONResponse{
			Message: "Hit the authentication service",
		}
		err := a.writeJSON(w, http.StatusOK, response, nil)
		if err != nil {
			a.errorJSON(w, err, http.StatusInternalServerError)
		}
	}
}

func (a *api) handleAuthenticate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := a.readJSON(w, r, &requestPayload)
		if err != nil {
			a.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		user, err := a.service.GetByEmail(requestPayload.Email)
		if err != nil {
			a.errorJSON(w, err)
			return
		}
		valid, err := a.service.PasswordMatches(user.Password, requestPayload.Password)
		if err != nil || !valid {
			a.errorJSON(w, errors.New("invalid username or password"), http.StatusUnauthorized)
			return
		}
		payload := JSONResponse{
			Message: fmt.Sprintf("Logged in user %s", user.Email),
			Data:    user,
		}
		err = a.writeJSON(w, http.StatusAccepted, payload, nil)
		if err != nil {
			a.errorJSON(w, err, http.StatusInternalServerError)
		}
	}
}

func (a *api) handleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := a.service.GetAll()
		if err != nil {
			a.errorJSON(w, err)
		}
		payload := JSONResponse{
			Message: "Users",
			Data:    users,
		}
		err = a.writeJSON(w, http.StatusOK, payload, nil)
		if err != nil {
			a.errorJSON(w, err, http.StatusInternalServerError)
		}
	}
}
