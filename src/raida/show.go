package raida

import (
	"logger"
	"config"
	"strconv"
//	"encoding/json"
//	"regexp"
	"cloudcoin"
	"error"
//	"fmt"
)

type Show struct {
	Servant
}

type ShowSubResponse struct {
	Sn string `json:"sn"`
	Tag string `json:"tag"`
	Created string `json:"created"`
}

type ShowResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
	Status string `json:"status"`
	Message []ShowSubResponse
}

type ShowOutput struct {
	AmountVerified int  `json:"amount_verified"`
}

func NewShow() (*Show) {
	return &Show{
		*NewServant(),
	}
}

func (v *Show) Show(cc *cloudcoin.CloudCoin) ([]int, int, *error.Error) {
	if !cloudcoin.ValidateCoin(cc) {
		return nil, 0, &error.Error{"CloudCoin is invalid"}
	}

	logger.Debug("Showing coins for " + cc.Sn)

	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range(params) {
		params[idx] = make(map[string]string)
		params[idx]["nn"] = cc.Nn
		params[idx]["sn"] = cc.Sn
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
			if (r.Status == "pass") {
        pownArray[idx] = config.RAIDA_STATUS_PASS
				logger.Debug("Raida " + strconv.Itoa(idx) + " shows " + strconv.Itoa(len(r.Message)) + " notes")
				snhash[idx] =  make([]int, len(r.Message))
				for sidx, ssr := range r.Message {
					isn, err := strconv.Atoi(ssr.Sn)
					if err != nil {
						logger.Debug("Skipping invalid SN " + ssr.Sn)
						continue
					}
					snhash[idx][sidx] = isn
				}
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
		return nil, 0, &error.Error{"Results from the RAIDA are not synchronized"}
	}

	sns, total := v.GetSNsOverlap(snhash)

	logger.Debug("Total Coins: " + strconv.Itoa(total))

	return sns, total, nil
/*
	logger.Debug("Started Show with UUID " + uuid + " owner " + owner)

	matched, err := regexp.MatchString(`^[A-Fa-f0-9]{32}$`, uuid)
	if err != nil || !matched {
		return "", &Error{"UUID invalid or not defined"}
	}

	sn, err := cloudcoin.GuessSNFromString(owner)
	if (err != nil) {
		return "", &Error{"Invalid Owner"}
	}

	logger.Debug("owner SN " +  strconv.Itoa(sn))

	pownArray := make([]int, v.Raida.TotalServers())
	balances := make(map[int]int)

	params := make(map[string]string)
	params["tag"] = uuid
	params["owner"] = strconv.Itoa(sn)

	results := v.Raida.SendRequest("/service/view_receipt", params, ShowResponse{})
  for idx, result := range results {
    if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
      r := result.Data.(*ShowResponse)
			if (r.Message != "success") {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
				balances[0]++
			} else {
				pownArray[idx] = config.RAIDA_STATUS_PASS
				total := r.TotalReceived

				balances[total]++
				logger.Debug("raida " + strconv.Itoa(idx) + " total " + strconv.Itoa(total))
			}

    } else if (result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT) {
			pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
			balances[0]++
		} else {
			pownArray[idx] = config.RAIDA_STATUS_ERROR
			balances[0]++
		}
  }

	pairs := sortByCount(balances)
	topBalance := pairs[0].Key

	logger.Debug("Most voted balance: " + strconv.Itoa(topBalance))
  for idx, result := range results {
    if result.ErrCode != config.REMOTE_RESULT_ERROR_NONE {
			continue
		}
    r := result.Data.(*ShowResponse)
		if (r.Message != "success") {
			continue
		}

		total := r.TotalReceived
		if (total != topBalance) {
			pownArray[idx] = config.RAIDA_STATUS_UNTRIED
			logger.Debug("Raida " + strconv.Itoa(topBalance) + " is reporting incorrect balance. Skipping it")
			continue
		}
	}

	for key, element := range balances {
		fmt.Printf("k=%d v=%d\n", key, element)
	}
	pownString := v.GetPownStringFromStatusArray(pownArray)
	logger.Debug("Pownstring " + pownString)

	if !v.IsStatusArrayFixable(pownArray) {
		return "", &Error{"Results from the RAIDA are not synchronized"}
	}

	vo := &ShowOutput{}
	vo.AmountVerified = topBalance

	b, err := json.Marshal(vo); 
	if err != nil {
		return "", &Error{"Failed to Encode JSON"}
	}

	return string(b), nil
	*/
}
