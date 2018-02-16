package main

import "net/http"

func main() {
	bindEndpoints()
}

func bindEndpoints() {
	http.HandleFunc("/notify/send", HandlerSend)
}

func HandlerSend(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok, I got this."))
}