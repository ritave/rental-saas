package main

import (
	"rental-saas/src/api_integration"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func test1() {
	p := api_integration.Provider{
		Client: http.DefaultClient,
		Server: "https://stage.pozamiatane.pl",
	}

	const limit = 5000
	for i:=0; i < limit; i++ {
		test := api_integration.EditRequest{
			OrderID: 106718,
			UserID: i,
		}

		resp, err := p.SendPayload(api_integration.EditAction, test)
		if err != nil {
			logrus.Printf("error: %s \n", err.Error())
		} else {
			logrus.Printf("\nresponse: %#v \n", *resp)
			bdy, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Printf("whaaa: %s \n", err.Error())
			} else {
				logrus.Printf("body %d: %s \n", i, bdy)
			}

		}

	}
}
