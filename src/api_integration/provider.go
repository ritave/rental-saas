package api_integration

import (
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/sirupsen/logrus"
)

type Provider struct {
	Client *http.Client
	Server string
}

func (p Provider) SendRequest(api string, payload interface{}) (*http.Response, error) {
	btz, _ := json.Marshal(payload)
	rdr := bytes.NewReader(btz)

	req, err := http.NewRequest(http.MethodPost, p.Server + api, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	logrus.Infof("Request: %#v", *req)

	return p.Client.Do(req)
}


