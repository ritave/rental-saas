package api_integration


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
	AddressID int    `json:"address_id"`
	Info      string `json:"info"`
	UserID    int    `json:"user_id"`
	Cleaners  []int  `json:"cleaners"`
}

const Save2Action = "/api/apiorders/save2"

/*
{"status":"ERROR","message":"Brak address_id lub user_id. Zam\u00f3wienie nie mo\u017ce by\u0107 zapisane"}



*/

type Save2ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Save2ResponseSuccess struct {

}