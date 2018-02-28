package utils

import (
	"runtime"
	"runtime/debug"
	"errors"
	"log"
)

// Debugging in the cloud is rather hard
func RecoverPanic() {
	if r := recover(); r != nil {
		if _, ok := r.(runtime.Error); ok {
			panic(r)
		} else if s, ok := r.(string); ok {
			log.Println(string(debug.Stack()))
			log.Println(errors.New(s))
		} else if e, ok := r.(error); ok {
			log.Println(string(debug.Stack()))
			log.Println(e)
		} else {
			panic(r)
		}
	}
}
