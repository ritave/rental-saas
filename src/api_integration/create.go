package api_integration

/* CREATE
{
	"client_id": "",
	"address": {
		"street": "",
		"zip": "?",
		"city": ""
	},
	"address_id": "aid?",

	"frequency": 0,
	"start": "",
	"length": "",
	"zip": "",
	"chemicals": "",
	"pets": "",
	"eng": "",
	"services": "",
	"info": "",

	"coupon_id": ""
}

if "address_id" is not set and "client_id" is not set, then "address_id" suffices
"frequency": 0 -> 'once', 7 -> '7days', else -> '14days'
 */

type CreateAction struct {
	ClientID string `json:"client_id"`
	Address  struct {
		Street string `json:"street"`
		Zip    string `json:"zip"`
		City   string `json:"city"`
	} `json:"address"`
	AddressID string `json:"address_id"`
	Frequency int    `json:"frequency"`
	Start     string `json:"start"`
	Length    string `json:"length"`
	Zip       string `json:"zip"`
	Chemicals string `json:"chemicals"`
	Pets      string `json:"pets"`
	Eng       string `json:"eng"`
	Services  string `json:"services"`
	Info      string `json:"info"`
	CouponID  string `json:"coupon_id"`
}

func (p Provider) CreateAction() {

}
