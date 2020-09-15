package raida

import (
	"logger"
	"config"
	"strconv"
	"encoding/json"
//	"regexp"
	"cloudcoin"
	"error"
//	"fmt"
//	"os"
	"core"
)

type Break struct {
	Servant
}

type BreakResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
	Status string `json:"status"`
}


func NewBreak() (*Break) {
	return &Break{
		*NewServant(),
	}
}

func (v *Break) Break(cc *cloudcoin.CloudCoin) (*error.Error) {
	if !cloudcoin.ValidateCoin(cc) {
		return &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "CloudCoin is invalid"}
	}

	logger.Debug("Break coin " + string(cc.Sn))

	cm := cloudcoin.GetChangeMethod(cc.GetDenomination())
	if cm == 0 {
		logger.Error("Failed to get Change Method")
		return &error.Error{config.ERROR_CHANGE_METHOD_NOT_FOUND, "Failed to get Change Method"}
	}

	s := NewShowChange()
	isn, _ := strconv.Atoi(string(cc.Sn))
	sns, err := s.ShowChange(cm, isn)
	if err != nil {
		logger.Error("Failed to ShowChange")
		return err
	}

	stringSns := make([]string, len(sns))
	for idx, ssn := range sns {
		stringSns[idx] = strconv.Itoa(ssn)
	}
	preParams := make([][]string, v.Raida.TotalServers())
	for ridx, _ := range(preParams) {
    preParams[ridx] = make([]string, len(stringSns))
    for idx, _ := range stringSns {
      preParams[ridx][idx], _ = cloudcoin.GeneratePan()
    }
  }

	var ccs []cloudcoin.CloudCoin
	for idx, sn := range stringSns {
		snInt, _ := strconv.Atoi(sn)
		cc := cloudcoin.NewFromData(config.DEFAULT_NN, snInt)

		for ridx := 0; ridx < v.Raida.TotalServers(); ridx++ {
			cc.SetAn(ridx, preParams[ridx][idx])
		}

		ccs = append(ccs, *cc)
	}


	ba, _ := json.Marshal(stringSns)
	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range(params) {
		baPans, _ := json.Marshal(preParams[idx])
		params[idx] = make(map[string]string)
		params[idx]["nn"] = string(cc.Nn)
		params[idx]["sn"] = string(cc.Sn)
		params[idx]["dn"] = strconv.Itoa(cc.GetDenomination())
		params[idx]["an"] = cc.Ans[idx]
		params[idx]["change_server"] = strconv.Itoa(config.PUBLIC_CHANGE_MAKER_ID)
		params[idx]["csn[]"] = string(ba)
		params[idx]["cpan[]"] = string(baPans)
	}

	results := v.Raida.SendDefinedRequest("/service/break", params, BreakResponse{})
  for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*BreakResponse)
			if (r.Status == "success") {
        pownArray[idx] = config.RAIDA_STATUS_PASS
			} else if (r.Status == "fail") {
        pownArray[idx] = config.RAIDA_STATUS_FAIL
      } else {
        pownArray[idx] = config.RAIDA_STATUS_ERROR
      }
		} else if (result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT) {
	      pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
		} else {
				pownArray[idx] = config.RAIDA_STATUS_ERROR
		}
	}

	pownString := v.GetPownStringFromStatusArray(pownArray)
  logger.Debug("Pownstring " + pownString)

	if !v.IsStatusArrayFixable(pownArray) {
		return &error.Error{config.ERROR_BREAK_FAILED, "Failed to Break Coin: " + pownString}
	}

	for _, cc := range ccs {
		core.SaveToBank(cc)
	}

	core.MoveCoinToSent(*cc)

	return nil
}
