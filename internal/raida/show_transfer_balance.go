package raida

import (
	"encoding/json"
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"
	//"fmt"
	//	"strings"
)

type ShowTransferBalance struct {
	Servant
}

type ShowTransferBalanceResponse struct {
	D1     int    `json:"d1"`
	D5     int    `json:"d5"`
	Total  int    `json:"total"`
	Status string `json:"status"`
}

type ShowTransferBalanceOutput struct {
	Amount int `json:"total"`
}

func NewShowTransferBalance() *ShowTransferBalance {
	return &ShowTransferBalance{
		*NewServant(),
	}
}

func (v *ShowTransferBalance) ShowTransferBalance(cc *cloudcoin.CloudCoin) (string, *error.Error) {
	if !cloudcoin.ValidateCoin(cc) {
		return "", &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "CloudCoin is invalid"}
	}

	logger.Debug("Showing transfer balance for " + string(cc.Sn))

	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range params {
		params[idx] = make(map[string]string)
		params[idx]["nn"] = string(cc.Nn)
		params[idx]["sn"] = string(cc.Sn)
		params[idx]["an"] = cc.Ans[idx]
		params[idx]["pan"] = cc.Ans[idx]
		params[idx]["denomination"] = strconv.Itoa(cc.GetDenomination())

	}

	balances := make(map[int]int)

	results := v.Raida.SendDefinedRequest("/service/show_transfer_balance", params, ShowTransferBalanceResponse{})
	for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*ShowTransferBalanceResponse)
			if r.Status == "pass" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
				logger.Debug("Raida " + strconv.Itoa(idx) + " shows " + strconv.Itoa(r.Total) + " balance")
				balances[r.Total]++
			} else if r.Status == "fail" {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
				balances[0]++
			} else {
				pownArray[idx] = config.RAIDA_STATUS_ERROR
				balances[0]++
			}
		} else if result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT {
			pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
			balances[0]++
		} else {
			pownArray[idx] = config.RAIDA_STATUS_ERROR
			balances[0]++
		}
	}

	pownString := v.GetPownStringFromStatusArray(pownArray)
	logger.Debug("Pownstring " + pownString)

	if !v.IsStatusArrayFixable(pownArray) {
		return "", &error.Error{config.ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC, "Results from the RAIDA are not synchronized"}
	}

	pairs := sortByCount(balances)
	topBalance := pairs[0].Key
	voters := pairs[0].Value

	logger.Debug("Most voted balance: " + strconv.Itoa(topBalance) + " voted: " + strconv.Itoa(voters))

	if voters < config.NEED_VOTERS_FOR_BALANCE {
		return "", &error.Error{config.ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC, "Results from the RAIDA are not synchronized. Not enough good results"}
	}

	/*
		sns, total := v.GetSNsOverlap(snhash)

		logger.Debug("Total Coins: " + strconv.Itoa(total))
	*/
	total := topBalance
	soutput := &ShowTransferBalanceOutput{}
	soutput.Amount = total

	b, err := json.Marshal(soutput)
	if err != nil {
		return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
	}

	return string(b), nil

}
