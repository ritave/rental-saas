package api_integration

import (
	"net/http"
	"encoding/json"
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Provider struct {
	Client *http.Client
	Server string
}

func (p Provider) SendPayload(api string, payload interface{}) (*http.Response, error) {
	req, err := p.NewRequest(api, payload)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Request: %#v", *req)

	return p.Client.Do(req)
}

//NewRequest creates new request in order to leave the possibility of adding headers and cookies later
func (p Provider) NewRequest(api string, payload interface{}) (*http.Request, error) {
	btz, _ := json.Marshal(payload)
	rdr := bytes.NewReader(btz)

	req, err := http.NewRequest(http.MethodPost, p.Server + api, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	return req, err
}

//SendRequestJustLog sends the request and logs what the service has returned
func (p Provider) SendRequestJustLog(req *http.Request) () {
	resp, err := p.Client.Do(req)
	if err != nil {
		logrus.Printf("Sending request to %s: %s", req.Host, err.Error())
	} else {
		// over the top
		logrus.Printf("Sending request to %s: %#v", req.Host, *resp)
		btz, _ := ioutil.ReadAll(resp.Body)
		logrus.Printf("Sending request to %s: %s", req.Host, btz)
		resp.Body.Close()
	}
}