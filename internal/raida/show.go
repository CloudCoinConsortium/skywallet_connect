package raida

import (
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"encoding/json"
	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"

	//"fmt"
	"strings"
)

type Show struct {
	Servant
}

type ShowSubResponse struct {
	Sn      string `json:"sn"`
	Tag     string `json:"tag"`
	Created string `json:"created"`
}

type ShowResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Status  string `json:"status"`
	Message []ShowSubResponse
}

type ShowOutput struct {
	AmountVerified int `json:"amount_verified"`
}

func NewShow() *Show {
	return &Show{
		*NewServant(),
	}
}

func (v *Show) Show(cc *cloudcoin.CloudCoin) ([]int, int, *error.Error) {
	if !cloudcoin.ValidateCoin(cc) {
		return nil, 0, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "CloudCoin is invalid"}
	}

	logger.Debug("Showing coins for " + string(cc.Sn))

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

	snhash := make([][]int, v.Raida.TotalServers())

	results := v.Raida.SendDefinedRequest("/service/show", params, ShowResponse{})
	for idx, result := range results {
		snhash[idx] = make([]int, 0)
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*ShowResponse)
			if r.Status == "pass" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
				logger.Debug("Raida " + strconv.Itoa(idx) + " shows " + strconv.Itoa(len(r.Message)) + " notes")
				snhash[idx] = make([]int, len(r.Message))
				for sidx, ssr := range r.Message {
					isn, err := strconv.Atoi(ssr.Sn)
					if err != nil {
						logger.Debug("Skipping invalid SN " + ssr.Sn)
						continue
					}
					snhash[idx][sidx] = isn
				}
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
		return nil, 0, &error.Error{config.ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC, "Results from the RAIDA are not synchronized"}
	}

	sns, total := v.GetSNsOverlap(snhash)

	logger.Debug("Total Coins: " + strconv.Itoa(total))

	return sns, total, nil
}

func (v *Show) ShowBrief(cc *cloudcoin.CloudCoin) ([]int, int, *error.Error) {
	if !cloudcoin.ValidateCoin(cc) {
		return nil, 0, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "CloudCoin is invalid"}
	}

	logger.Debug("Showing Brief Coins for " + string(cc.Sn))

	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range params {
		params[idx] = make(map[string]string)
		params[idx]["nn"] = string(cc.Nn)
		params[idx]["sn"] = string(cc.Sn)
		params[idx]["an"] = cc.Ans[idx]
		params[idx]["pan"] = cc.Ans[idx]
		params[idx]["denomination"] = strconv.Itoa(cc.GetDenomination())
		params[idx]["b"] = "t"

	}

	snhash := make([][]int, v.Raida.TotalServers())

	results := v.Raida.SendDefinedRequestRaw("/service/show", params)
	for idx, result := range results {
		snhash[idx] = make([]int, 0)
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			if result.Message == "fail" {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
			} else {
				ss := strings.Split(result.Message, ",")
				if len(ss) == 0 {
					pownArray[idx] = config.RAIDA_STATUS_ERROR
				} else {
					logger.Debug("Raida " + strconv.Itoa(idx) + " shows " + strconv.Itoa(len(ss)) + " notes")
					pownArray[idx] = config.RAIDA_STATUS_PASS
					snhash[idx] = make([]int, len(ss))
					for sidx, sn := range ss {
						isn, err := strconv.Atoi(sn)
						if err != nil {
							logger.Debug("Skipping invalid SN " + sn)
							continue
						}
						snhash[idx][sidx] = isn
					}
				}
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
		return nil, 0, &error.Error{config.ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC, "Results from the RAIDA are not synchronized"}
	}

	sns, total := v.GetSNsOverlap(snhash)

	logger.Debug("Total Coins: " + strconv.Itoa(total))

	return sns, total, nil
}
