package raida

import (
	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"strconv"
	"encoding/json"
	//	"sort"
	"regexp"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"
	//	"strings"
)

type Statement struct {
	Servant
}

type StatementResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

type StatementOutput struct {
	AmountVerified int    `json:"amount_verified"`
	Status         string `json:"status"`
	Message        string `json:"message"`
}

func NewStatement() *Statement {
	return &Statement{
		*NewServant(),
	}
}

func (v *Statement) Create(uuid string, memo string, amount string, from string, cc *cloudcoin.CloudCoin) (string, *error.Error) {
	logger.Debug("Creating Statement with UUID " + uuid + " Our SN " + string(cc.Sn))

	matched, err := regexp.MatchString(`^[A-Fa-f0-9]{32}$`, uuid)
	if err != nil || !matched {
		return "", &error.Error{config.ERROR_INVALID_RECEIPT_ID, "UUID invalid or not defined"}
	}

	sn := string(cc.Sn)
	pownArray := make([]int, v.Raida.TotalServers())
	d := v.GetStripesMirrorsForObjectMemo(uuid, memo, amount, from)

	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range params {
		params[idx] = make(map[string]string)
		params[idx]["account_sn"] = sn
		params[idx]["account_an"] = cc.Ans[idx]
		params[idx]["transaction_id"] = uuid
		params[idx]["version"] = "0"
		params[idx]["compression"] = "0"
		params[idx]["raid"] = "110"
		params[idx]["stripe"] = d[idx][0]
		params[idx]["mirror"] = d[idx][1]
		params[idx]["mirror2"] = d[idx][2]
	}

	results := v.Raida.SendDefinedRequest("/service/statements/create", params, StatementResponse{})
	for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*StatementResponse)
			if r.Message == "success" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
			} else if r.Message == "fail" {
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

	/*
		logger.Debug("Running FixTransfer " + pownString)
		ft := NewFixTransfer()
		for ssn, _ := range allSns {
			absent := 0
			for _, farr := range fArr {
				if _, ok := farr[ssn]; !ok {
					absent++
				}
			}

			if (absent > config.TOTAL_RAIDA_NUMBER / 2) {
					logger.Debug("Coin " + strconv.Itoa(ssn) + " is absent on a lot of servers")
					for ridx, farr := range fArr {
						if _, ok := farr[ssn]; ok {
							logger.Debug("Coin " + strconv.Itoa(ssn) + " will be fixed (removal) on raida " + strconv.Itoa(ridx))
							ft.AddSNToRepairArray(ridx, ssn)
						}
					}
					continue
			}

			for ridx, farr := range fArr {
				if _, ok := farr[ssn]; !ok {
					logger.Debug("Coin " + strconv.Itoa(ssn) + " will be fixed on raida " + strconv.Itoa(ridx))
					ft.AddSNToRepairArray(ridx, ssn)
				}
			}

			//for ridx, _ := range fArr {
			//		ft.AddSNToRepairArray(ridx, ssn)
			//}

		}

		ft.FixTransfer()
	*/

	//if !v.IsStatusArrayFixable(pownArray) {
	//	return "", &error.Error{config.ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC, "Results from the RAIDA are not synchronized"}
	//}

	//if voted < config.NEED_VOTERS_FOR_BALANCE {
	//		return "", &error.Error{config.ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC, "Results from the RAIDA are not synchronized. Not enough good results"}
	//	}

	vo := &StatementOutput{}
	vo.Status = "success"
	vo.Message = "Statement created"

	b, err := json.Marshal(vo)
	if err != nil {
		return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
	}

	//fmt.Printf("ns=%d %s isok=%b\n", v.Raida.TotalServers(), pownString, v.IsStatusArrayFixable(pownArray))
	return string(b), nil
}
