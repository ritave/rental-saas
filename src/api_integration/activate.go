package api_integration

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

const ActivateAction = "/api/apiorders/activate/%d"

func (p Provider) Activate(orderID int) (suc ActivateResponse, err error) {
	resp, err := p.Client.Get(fmt.Sprintf(p.Server + ActivateAction, orderID))

	btz, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return suc, err
	}

	err = json.Unmarshal(btz, &suc)
	return suc, err
}

type ActivateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
