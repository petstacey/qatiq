package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (b *broker) handleBroker() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := JSONResponse{
			Message: "Hit the broker service",
		}
		err := b.writeJSON(w, http.StatusOK, response, nil)
		if err != nil {
			b.errorJSON(w, err, http.StatusInternalServerError)
		}
	}
}

func (b *broker) handleAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest("GET", "http://authentication-service/v1/", nil)
		if err != nil {
			b.errorJSON(w, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			b.errorJSON(w, err)
			return
		}
		defer response.Body.Close()

		var jsonFromService JSONResponse

		err = json.NewDecoder(response.Body).Decode(&jsonFromService)
		if err != nil {
			b.errorJSON(w, err)
			return
		}

		if jsonFromService.Error {
			b.errorJSON(w, err, http.StatusUnauthorized)
			return
		}

		payload := JSONResponse{
			Message: jsonFromService.Message,
			Data:    jsonFromService.Data,
		}

		err = b.writeJSON(w, http.StatusOK, payload, nil)
		if err != nil {
			b.errorJSON(w, err)
		}
	}
}

func (b *broker) handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload RequestPayload

		err := b.readJSON(w, r, &payload)
		if err != nil {
			b.errorJSON(w, err)
			return
		}

		switch payload.Action {
		case "auth":
			b.authenticateUser(w, payload.Auth)
		case "log":
			b.logItem(w, payload.Log)
		default:
			b.errorJSON(w, errors.New("unknown action"))
		}
	}
}

func (b *broker) authenticateUser(w http.ResponseWriter, auth AuthPayload) {
	data, err := json.Marshal(auth)
	if err != nil {
		b.errorJSON(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://authentication-service/v1/authenticate", bytes.NewBuffer(data))
	if err != nil {
		b.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		b.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		b.errorJSON(w, errors.New("invalid username or password"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		b.errorJSON(w, errors.New("error calling authentication service"))
		return
	}

	var jsonFromService JSONResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		b.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		b.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := JSONResponse{
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	err = b.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		b.errorJSON(w, err)
	}
}

func (b *broker) logItem(w http.ResponseWriter, entry LogPayload) {
	loggerServiceURL := "http://logger-service/v1/log"

	jdata, _ := json.Marshal(entry)

	request, err := http.NewRequest("POST", loggerServiceURL, bytes.NewBuffer(jdata))
	if err != nil {
		b.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		b.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		b.errorJSON(w, err)
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "logged"

	b.writeJSON(w, http.StatusAccepted, payload)
}
