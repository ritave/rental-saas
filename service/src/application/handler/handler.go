package handler

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"rental-saas/service/src/application/core"
	"github.com/sirupsen/logrus"
)

// this generality will be a loss in performance, I can guarantee that
type AppHandler struct {
	App *core.Application
	RequestTemplate interface{}
	Handler func(a *core.Application, r interface{}) (interface{}, error)
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody := h.RequestTemplate

	logrus.Info("Incoming ", r.RequestURI)

	// parsing
	btz, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Infof("Outgoing: %d | %s", http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(btz, requestBody)
	if err != nil {
		logrus.Infof("Outgoing: %d | %s", http.StatusBadRequest, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// executing
	responseBody, err := h.Handler(h.App, requestBody)
	if err != nil {
		logrus.Infof("Outgoing: %d | %s", http.StatusInternalServerError, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// wrapping and responding
	btz, _ = json.Marshal(responseBody)
	logrus.Infof("Outgoing: %d | %s", http.StatusBadRequest, btz)
	w.Write(btz)
}

