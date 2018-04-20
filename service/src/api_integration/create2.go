package api_integration

import (
	"encoding/json"
	"io/ioutil"
	"errors"
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
	Address   Address `json:"address,omitempty"`
	AddressID int     `json:"address_id,omitempty"`
	Frequency int     `json:"frequency,omitempty"`
	Start     string  `json:"start,omitempty"`
	Length    float64 `json:"length,omitempty"`
	Zip       string  `json:"zip,omitempty"`
	Chemicals int     `json:"chemicals,omitempty"`
	Pets      int     `json:"pets,omitempty"`
	Eng       int     `json:"eng,omitempty"`
	Services  []int   `json:"services,omitempty"`
	Osource   string  `json:"osource,omitempty"`
	Info      string  `json:"info,omitempty"`
	CouponID  int     `json:"coupon_id,omitempty"`
}

type Address struct {
	Street string `json:"street,omitempty"`
	Zip    string `json:"zip,omitempty"`
	City   string `json:"city,omitempty"`
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
	Cleaners []Create2ResponseCleaner `json:"cleaners"`
}
type Create2ResponseCleaner struct {
	ID        int           `json:"id"`
	Priority  string        `json:"priority"`
	Avatar    interface{}   `json:"avatar"`
	Name      string        `json:"name"`
	ShortName string        `json:"short_name"`
	Services  []interface{} `json:"services"`
	Rates     []interface{} `json:"rates"`
	Rate      int           `json:"rate"`
	Overall   struct {
		Count int         `json:"count"`
		Avg   interface{} `json:"avg"`
	} `json:"overall"`
	Stage int `json:"stage"`
	Type  int `json:"type"`
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
