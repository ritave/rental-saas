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
  "user_id": 123,
  "cleaning_date": "2018-02-28",
  "cleaning_time": "15:00:00",
  "length": 4.5
}
*/
/*
zapisuje zamowienie z wybranymi sprzataczami
moze wystepowac wersja z address zamiast address_id analogicznie jak w przypadku create2
 */

type EditRequest struct {
	OrderID      int     `json:"order_id"`
	UserID       int     `json:"user_id"`
	CleaningDate string  `json:"cleaning_date,omitempty"`
	CleaningTime string  `json:"cleaning_time,omitempty"`
	Length       float64 `json:"length,omitempty"`
	Address      Address `json:"address,omitempty"`
	AddressID    int     `json:"address_id,omitempty"`
}

const EditAction = "/api/apiorders/edit"

/*
{"status":"ERROR","message":" Nie znaleziono zam\u00f3wienia o ID 1"}
{"status":"ERROR","message":"Improper client ID!"}


*/

func (p Provider) Edit(payload EditRequest) (suc EditResponseSuccess, err error) {
	resp, err := p.SendPayload(EditAction, payload)
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
			return suc, errors.New(fmt.Sprintf("%#v", fail))
		}
	}

	return suc, err

}

type EditResponseSuccess struct {
	Status  string `json:"status"`
	Message struct {
		ID                       string      `json:"id"`
		ParentID                 interface{} `json:"parent_id"`
		Type                     string      `json:"type"`
		Start                    string      `json:"start"`
		End                      string      `json:"end"`
		Cycle                    string      `json:"cycle"`
		Status                   string      `json:"status"`
		CleanerID                string      `json:"cleaner_id"`
		UserID                   string      `json:"user_id"`
		Length                   float64         `json:"length"`
		Zip                      string      `json:"zip"`
		AddressID                string      `json:"address_id"`
		Chemicals                string      `json:"chemicals"`
		Pets                     string      `json:"pets"`
		CouponID                 interface{} `json:"coupon_id"`
		Info                     string      `json:"info"`
		CancelRequest            interface{} `json:"cancel_request"`
		CancelInfo               interface{} `json:"cancel_info"`
		Created                  string      `json:"created"`
		Edited                   interface{} `json:"edited"`
		CleanerChanged           interface{} `json:"cleaner_changed"`
		ChangedLength            interface{} `json:"changed_length"`
		Eng                      string      `json:"eng"`
		Osource                  string      `json:"osource"`
		Updated                  string      `json:"updated"`
		RatingHash               interface{} `json:"rating_hash"`
		TransferPayment          interface{} `json:"transferPayment"`
		UpdatedBy                int         `json:"updated_by"`
		ChangedBecauseOfOurFault interface{} `json:"changed_because_of_our_fault"`
		DebtSettlementHash       interface{} `json:"debt_settlement_hash"`
	} `json:"message"`
}

type EditResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
