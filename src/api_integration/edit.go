package api_integration

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
	CleaningDate string  `json:"cleaning_date"`
	CleaningTime string  `json:"cleaning_time"`
	Length       float64 `json:"length"`
	Address      Address `json:"address"`
}

const EditAction = "/api/apiorders/edit"

/*
{"status":"ERROR","message":" Nie znaleziono zam\u00f3wienia o ID 1"}
{"status":"ERROR","message":"Improper client ID!"}


*/

type EditResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
