package raida

import (
	"logger"
	"encoding/json"
	"strconv"
	"math/rand"
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

func (v *FixTransfer) FixTransfer() {
	logger.Debug("Started FixTransfer")

	params := make([]map[string]string, v.Raida.TotalServers())
//	params["corner"] = "1"

	for ridx, rarr := range v.repairArray {
		if len(rarr) == 0 {
			params[ridx] = nil
			continue
		}

		rnumber := rand.Intn(4) + 1
		cstr := strconv.Itoa(rnumber)
		logger.Debug("Random corner " + cstr)
		params[ridx] = make(map[string]string)
		params[ridx]["corner"] = cstr

		stringRarr := make([]string, len(rarr))
		for idx, ssn := range rarr {
			stringRarr[idx] = strconv.Itoa(ssn)
		}

		ba, _ := json.Marshal(stringRarr)
		params[ridx]["sn[]"] = string(ba)
	}
	
	v.Raida.SendDefinedRequestNoWait("/service/sync/fix_transfer", params, FixTransferResponse{})

	return
}


