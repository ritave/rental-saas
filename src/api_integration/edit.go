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
}

func (p Provider) EditAction() {

}
