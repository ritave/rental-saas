package api_integration

/*
{
	"zip": "",
	"order_id": 0
}

 */

type LoadTestSearchRequest struct {
	Zip     string `json:"zip"`
	OrderID int    `json:"order_id"`
}
