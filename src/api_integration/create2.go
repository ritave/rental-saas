package api_integration

import (
	"encoding/json"
	"io/ioutil"
	"errors"
	"fmt"
	"strings"
)

/*
{
  "client_id": 1,
  "address_id": 12,
  "frequency": 7,
  "start": "2018-06-01 12:00:00",
  "length": 3.5,
  "zip": "00-001",
  "chemicals": 1,
  "pets": 0,
  "eng": 1,
  "services": [
    12
  ],
  "osource": "A",
  "info": "DODATKOWY OPIS SPRZATANIA",
  "coupon_id": 123
}
lub dla nowego klienta
{
  "client_id": 1,
  "address": {
    "street": "Testowa",
    "zip": "00-001",
    "city": "Testowo"
  },
  "frequency": 7,
  "start": "2018-06-01 12:00:00",
  "length": 3.5,
  "zip": "00-001",
  "chemicals": 1,
  "pets": 0,
  "eng": 1,
  "services": [
    12
  ],
  "osource": "A",
  "info": "DODATKOWY OPIS SPRZATANIA",
  "coupon_id": 123
}

 */

type Create2ActionRequest struct {
	ClientID  int     `json:"client_id"`
	Address   Address `json:"address"`
	Frequency int     `json:"frequency"`
	Start     string  `json:"start"`
	Length    float64 `json:"length"`
	Zip       string  `json:"zip"`
	Chemicals int     `json:"chemicals"`
	Pets      int     `json:"pets"`
	Eng       int     `json:"eng"`
	Services  []int   `json:"services"`
	Osource   string  `json:"osource"`
	Info      string  `json:"info"`
	CouponID  int     `json:"coupon_id"`
}

type Address struct {
	Street string `json:"street"`
	Zip    string `json:"zip"`
	City   string `json:"city"`
}

const Create2Action = "/api/apiorders/create2"

// response to
/*

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
*/
/*
 {"order_id":"106716","cleaners":[]}
 or on error
 ["Pole Czas trwania jest wymagane","Pole Kod pocztowy jest wymagane"]
*/

type Create2ResponseSuccess struct {
	OrderID  string        `json:"order_id"`
	Cleaners []interface{} `json:"cleaners"`
}

type Create2ResponseError []string

func (p Provider) Create2(payload Create2ActionRequest) (suc Create2ResponseSuccess, err error) {
	resp, err := p.SendPayload(Create2Action, payload)
	if err != nil {
		return suc, err
	}

	btz, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return suc, err
	}

	err = json.Unmarshal(btz, &suc)
	if err != nil {
		// try again but with fail version
		fail := Create2ResponseError{}
		err2 := json.Unmarshal(btz, &fail)
		if err2 != nil {
			return suc, err
		} else {
			return suc, errors.New(strings.Join(fail, " "))
		}
	}

	return suc, err
}

var create2Test = Create2ActionRequest{
	ClientID: 1,
	Address: Address{
		Street: "Testowa",
		Zip:    "00-001",
		City:   "Testowo",
	},
	Frequency: 7,
	Start:     "2018-06-01 12:00:00",
	Length:    3.5,
	Zip:       "00-001",
	Chemicals: 1,
	Pets:      0,
	Eng:       1,
	Services:  []int{12},
	Osource:   "A",
	Info:      "extra",
	CouponID:  123,
}
