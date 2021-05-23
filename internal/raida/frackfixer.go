package raida

import (
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"
)

type FrackFixer struct {
	Servant
}

type FrackFixerResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type FrackFixerOutput struct {
  AmountFixed int
}

func NewFrackFixer() *FrackFixer {
	return &FrackFixer{
		*NewServant(),
	}
}

func (v *FrackFixer) Fix() (*FrackFixerOutput, *error.Error) {
	logger.Debug("Started FrackFixer")


  fo := &FrackFixerOutput{}
  fo.AmountFixed = 0

  return fo, nil

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

	v.Raida.SendDefinedRequestNoWait("/service/sync/fix_transfer", params, FrackFixerResponse{})

	return nil, nil
}
