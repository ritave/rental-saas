package utils

import (
	"encoding/json"
	"net/http"
)

// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error/33926433
// sooo... basically under completely normal circumstances it cannot fail

func MustMarshallJSON(v interface{}) ([]byte) {
	btz, _ := json.Marshal(v)
	return btz
}

func WriteAsJSON(w http.ResponseWriter, v interface{}) {
	w.Write(MustMarshallJSON(v))
}
