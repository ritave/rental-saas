package handler

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"rental-saas/src/application/core"
)

// this generality will be a loss in performance, I can guarantee that
type AppHandler struct {
	App *core.Application
	RequestTemplate interface{}
	Handler func(a *core.Application, r interface{}) (interface{}, error)
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody := h.RequestTemplate

	// parsing
	btz, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(btz, requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// executing
	responseBody, err := h.Handler(h.App, requestBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// wrapping and responding
	btz, _ = json.Marshal(responseBody)
	w.Write(btz)
}

