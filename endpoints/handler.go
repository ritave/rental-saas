package endpoints

import (
	"net/http"
	"google.golang.org/api/calendar/v3"
)

type MyHandler struct {
	S     *calendar.Service
	HFunc func(*calendar.Service, http.ResponseWriter, *http.Request)
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HFunc(h.S, w, r)
}

