package main

import (
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/notify/send", HandlerSend)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HandlerSend(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok, I got this."))
}