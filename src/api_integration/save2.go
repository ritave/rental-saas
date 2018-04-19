package api_integration

import (
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
)

/*
{
  "order_id": 1,
  "address_id": 12,
  "info": "Dodatkowy opis zamowienia",
  "user_id": 123,
  "cleaners": [
    1234,
    1235,
    1236
  ]
}
 */

type Save2Request struct {
	OrderID   int    `json:"order_id"`
	Address   Address `json:"address,omitempty"`
	AddressID int    `json:"address_id,omitempty"`
	Info      string `json:"info,omitempty"`
	UserID    int    `json:"user_id"`
	Cleaners  []Save2RequestCleaner  `json:"cleaners"`
}

type Save2RequestCleaner struct {
	ID    string `json:"id"`
	Stage string `json:"stage"`
}

const Save2Action = "/api/apiorders/save2"

/*
{"status":"ERROR","message":"Brak address_id lub user_id. Zam\u00f3wienie nie mo\u017ce by\u0107 zapisane"}



*/

func (p Provider) Save2(payload Save2Request) (suc Save2Response, err error) {
	resp, err := p.SendPayload(Save2Action, payload)
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
		fail := Save2ResponseError{}
		err2 := json.Unmarshal(btz, &fail)
		if err2 != nil {
			return suc, err
		} else {
			return suc, errors.New(fmt.Sprintf("%#v", fail))
		}
	}

	return suc, err

}
type Save2ResponseError []interface{}

type Save2Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}