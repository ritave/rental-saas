package main

import (
	"rental-saas/src/api_integration"
	"net/http"
	"fmt"
	"io/ioutil"
)

var printf = fmt.Printf

func main() {
	p := api_integration.Provider{
		Client: http.DefaultClient,
		Server: "https://stage.pozamiatane.pl",
	}

	test := api_integration.Create2ActionRequest{
		ClientID: 1,
		Address: api_integration.Address{
			Street: "Testowa",
			Zip: "00-001",
			City: "Testowo",
		},
		Frequency: 7,
		Start: "2018-06-01 12:00:00",
		Length: 3.5,
		Zip: "00-001",
		Chemicals: 1,
		Pets: 0,
		Eng: 1,
		Services: []int{12},
		Osource: "A",
		Info: "extra",
		CouponID: 123,
	}
	
	resp, err := p.SendRequest(api_integration.Create2Action, test)
	if err != nil {
		fmt.Printf("error: %s \n", err.Error())
	} else {
		fmt.Printf("response: %#v \n", *resp)
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("whaaa: %s \n", err.Error())
		} else {
			fmt.Printf("body: %s \n", bdy)
		}

	}

}