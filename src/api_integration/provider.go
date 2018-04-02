package api_integration

import "net/http"

type Provider struct {
	Client *http.Client
	API string
}

