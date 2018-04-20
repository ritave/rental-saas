package notify

import (
	"net/http"
	"rental-saas/service/src/presenter"
	"context"
	"io"
	"encoding/json"
	"rental-saas/service/src/utils"
	"github.com/sirupsen/logrus"
)

type DeleteChannelRequest []DeleteChannelSingleRequest

type DeleteChannelSingleRequest struct {
	ResourceID string `json:"resource_id"`
	UUID       string `json:"uuid"`
}

type DeleteChannelResponse map[string]string

/*
{"resource_id": "cadPPw6iw2Bs3NSAhWvXuPWETa4","uuid": ""}
 */

func DeleteChannel(w http.ResponseWriter, r *http.Request) {
	req, err := extractDeleteChannelRequest(r.Body)
	if err != nil {
		w.Write(mustMarshalResponse(DeleteChannelResponse{"error": err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	background := context.Background()
	cal := utils.NewFlex(background)
	resp := make(DeleteChannelResponse)

	for _, single := range req {
		err = presenter.StopChannel(cal, single.ResourceID, single.UUID)
		if err != nil {
			resp[single.UUID] = err.Error()
		}
	}

	w.Write(mustMarshalResponse(resp))
}

func extractDeleteChannelRequest(r io.ReadCloser) (DeleteChannelRequest, error) {
	defer r.Close()
	var target = make(DeleteChannelRequest, 0)
	err := json.NewDecoder(r).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func mustMarshalResponse(response DeleteChannelResponse) ([]byte) {
	bytez, err := json.Marshal(response)
	if err != nil {
		logrus.Println("Wooooo, it really failed")
		return []byte{}
	}
	return bytez
}
