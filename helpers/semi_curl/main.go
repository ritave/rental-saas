package main

import (
	"rental-saas/src/api_integration"
	"net/http"
	"encoding/json"
	"bytes"
)

func main() {
	p := api_integration.Provider{
		Client: http.DefaultClient,
		API: "https://stage.pozamiatane.pl",
	}

	test := api_integration.LoadTestSearchRequest{
		OrderID: 0,
		Zip: "00-001",
	}
	btz, _ := json.Marshal(test)
	r := bytes.NewReader(btz)

	resp, err := http.NewRequest(http.MethodPost, p.API, r)
	if err != nil {
		println(err.Error())
	} else {
		println(resp)
	}

}