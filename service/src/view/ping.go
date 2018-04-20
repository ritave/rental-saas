package view

import (
	"fmt"
	"net/http"
)

func HandlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pong")
}
