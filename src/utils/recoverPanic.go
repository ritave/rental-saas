package utils

import (
	"runtime"
	"runtime/debug"
	"errors"
	"github.com/sirupsen/logrus"
)

// Debugging in the cloud is rather hard
func RecoverPanic() {
	if r := recover(); r != nil {
		if _, ok := r.(runtime.Error); ok {
			panic(r)
		} else if s, ok := r.(string); ok {
			logrus.Println(string(debug.Stack()))
			logrus.Println(errors.New(s))
		} else if e, ok := r.(error); ok {
			logrus.Println(string(debug.Stack()))
			logrus.Println(e)
		} else {
			panic(r)
		}
	}
}
