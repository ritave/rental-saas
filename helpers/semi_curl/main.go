package main

import (
	"rental-saas/src/api_integration"
	"net/http"
	"fmt"
	"github.com/sirupsen/logrus"
	"encoding/json"
)

var _ = fmt.Printf
var _ = http.Get

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	p := api_integration.Provider{
		Client: http.DefaultClient,
		Server: "https://stage.pozamiatane.pl",
	}

	create := api_integration.Create2ActionRequest{
		ClientID: 1,
		Address: api_integration.Address{
			Street: "Testowa",
			Zip:    "02-103",
			City:   "Testowo",
		},
		Frequency: 7,
		Start:     "2018-06-01 12:00:00",
		Length:    3.5,
		Zip:       "02-103",
		Chemicals: 1,
		Pets:      0,
		Eng:       1,
		Services:  []int{12},
		Osource:   "A",
		Info:      "extra",
		CouponID:  123,
	}


}

