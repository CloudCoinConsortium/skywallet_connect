package raida

import (
	"logger"
)

type FixTransfer struct {
	Servant
}

type FixTransferResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
	Message string `json:"message"`
	Status string `json:"status"`
}

func NewFixTransfer() (*FixTransfer) {
	return &FixTransfer{
		*NewServant(),
	}
}

func (v *FixTransfer) FixTransfer(sns []int) {
	logger.Debug("Started FixTransfer")

	params := make(map[string]string)
	params["corner"] = "1"

	v.Raida.SendRequest("/service/sync/fix_transfer", params, FixTransferResponse{})


	return
}


