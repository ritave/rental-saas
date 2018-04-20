package notify

import "rental-saas/service/src/application/core"

type GetRequest struct {}
type GetResponse struct {}

func Get(a *core.Application, r interface{}) (interface{}, error) {
	a.Utils.Ticker.Restart()
	return nil, nil
}
