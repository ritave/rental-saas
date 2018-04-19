package api_integration

import (
	"errors"
	"net/http"
	"fmt"
)

/*
{
  "orderId": 106761,
  "why": "bo to api ssie",
  "canceler": "client",
  "cycle": 1
}
 */

const CancelAction = "/api/apiorders/cancel"

type CancelRequest struct {
	OrderID  int    `json:"orderId"`
	Why      string `json:"why"`
	Canceler string `json:"canceler"` // client / cleaner
	Cycle    int    `json:"cycle"`
}

func (p Provider) Cancel(payload CancelRequest) (err error) {
	resp, err := p.SendPayload(CancelAction, payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("%v", resp))
	}

	return err
}

// Nothing is returned, lol
type CancelResponse struct{}
