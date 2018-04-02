package api_integration

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
	ClientID int `json:"client_id"`
	Address  Address `json:"address"`
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
