package raida

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"
	//	"fmt"
	//	"os"
)

type BreakInBank struct {
	Servant
}

type BreakInBankResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Status  string `json:"status"`
}

func NewBreakInBank() *BreakInBank {
	return &BreakInBank{
		*NewServant(),
	}
}

func (v *BreakInBank) BreakInBank(cc *cloudcoin.CloudCoin, snToBreak int) ([]int, *error.Error) {
	if !cloudcoin.ValidateCoin(cc) {
		return nil, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "CloudCoin is invalid"}
	}

	logger.Debug("BreakInBank coin " + strconv.Itoa(snToBreak))

	cm := cloudcoin.GetChangeMethod(cloudcoin.GetDenomination(snToBreak))
	if cm == 0 {
		logger.Error("Failed to get Change Method")
		return nil, &error.Error{config.ERROR_CHANGE_METHOD_NOT_FOUND, "Failed to get Change Method"}
	}

	s := NewShowChange()
	sns, err := s.ShowChange(cm, snToBreak)
	if err != nil {
		logger.Error("Failed to ShowChange")
		return nil, err
	}

  fmt.Printf("xxx=%v\n",sns)
  os.Exit(1)

	stringSns := make([]string, len(sns))
	for idx, ssn := range sns {
		stringSns[idx] = strconv.Itoa(ssn)
	}
	ba, _ := json.Marshal(stringSns)

	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range params {
		params[idx] = make(map[string]string)
		params[idx]["id_nn"] = string(cc.Nn)
		params[idx]["id_sn"] = string(cc.Sn)
		params[idx]["id_an"] = cc.Ans[idx]
		params[idx]["id_dn"] = strconv.Itoa(cc.GetDenomination())
		params[idx]["nn"] = strconv.Itoa(config.DEFAULT_NN)
		params[idx]["sn"] = strconv.Itoa(snToBreak)
		params[idx]["dn"] = strconv.Itoa(cloudcoin.GetDenomination(snToBreak))
		params[idx]["change_server"] = strconv.Itoa(config.PUBLIC_CHANGE_MAKER_ID)
		params[idx]["csn[]"] = string(ba)
	}

	results := v.Raida.SendDefinedRequest("/service/break_in_bank", params, BreakInBankResponse{})
	for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*BreakInBankResponse)
			if r.Status == "success" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
			} else if r.Status == "fail" {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
			} else {
				pownArray[idx] = config.RAIDA_STATUS_ERROR
			}
		} else if result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT {
			pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
		} else {
			pownArray[idx] = config.RAIDA_STATUS_ERROR
		}
	}

	pownString := v.GetPownStringFromStatusArray(pownArray)
	logger.Debug("Pownstring " + pownString)

	if !v.IsStatusArrayFixable(pownArray) {
		return nil, &error.Error{config.ERROR_BREAK_IN_BANK_FAILED, "Failed to Break Coin: " + pownString}
	}

	return sns, nil
}
