package raida

import (
	"logger"
	"config"
	"strconv"
	"encoding/json"
//	"sort"
	"regexp"
	"cloudcoin"
	"error"
	"strings"
//	"fmt"
)

type PaymentVerifier struct {
	Servant
}

type PaymentVerifierResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
	TotalReceived int `json:"total_received"`
	Message string `json:"message"`
	SerialNumbers string `json:"serial_numbers"`
	Memo string `json:"memo"`
}

type PaymentVerifierOutput struct {
	AmountVerified int  `json:"amount_verified"`
	Status string `json:"status"`
	Message string `json:"message"`
}

func NewPaymentVerifier() (*PaymentVerifier) {
	return &PaymentVerifier{
		*NewServant(),
	}
}

func (v *PaymentVerifier) Verify(uuid string, cc *cloudcoin.CloudCoin) (string, *error.Error) {
	logger.Debug("Started PaymentVerifier with UUID " + uuid + " Our SN " + string(cc.Sn))

	matched, err := regexp.MatchString(`^[A-Fa-f0-9]{32}$`, uuid)
	if err != nil || !matched {
		return "", &error.Error{config.ERROR_INVALID_RECEIPT_ID, "UUID invalid or not defined"}
	}

	sn := string(cc.Sn)

	pownArray := make([]int, v.Raida.TotalServers())
	balances := make(map[int]int)
	guids := make(map[string]int)
	memos := make(map[string]int)

	allSns := make(map[int]bool)

	amount := "0"
	memo := "received"
	tags := v.GetObjectMemo("", memo, amount, cc.GetFileName())

	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range(params) {
		params[idx] = make(map[string]string)
		params[idx]["tag"] = uuid
		params[idx]["owner"] = sn
		params[idx]["an"] = cc.Ans[idx]
		params[idx]["new_memo"] = tags[idx]
	}

	fArr := make([]map[int]bool, v.Raida.TotalServers())

	results := v.Raida.SendDefinedRequest("/service/verify_payment", params, PaymentVerifierResponse{})
  for idx, result := range results {
    if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
      r := result.Data.(*PaymentVerifierResponse)
			if (r.Message == "success") {
				pownArray[idx] = config.RAIDA_STATUS_PASS
				total := r.TotalReceived
				balances[total]++
				logger.Debug("raida " + strconv.Itoa(idx) + " total " + strconv.Itoa(total))

				guid, errg := v.GetGuidFromMemo(r.Memo)
				if errg == nil {
					logger.Debug("XGUid="+guid)
					guids[guid]++
				}

				memo, errg := v.GetMemoFromMemo(r.Memo)
				if errg == nil {
					logger.Debug("XMemo="+memo)
					memos[memo]++
				}

				snsString := r.SerialNumbers
				sns := strings.Split(snsString, ",")
				fArr[idx] = make(map[int]bool)

				for _, ssn := range sns {
					issn, err := strconv.Atoi(ssn)
					if err != nil {
						continue
					}
					fArr[idx][issn] = true
					allSns[issn] = true
				}
			} else if (r.Message == "fail") {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
				balances[0]++
			} else {
				pownArray[idx] = config.RAIDA_STATUS_ERROR
				balances[0]++
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
	voted := pairs[0].Value

	logger.Debug("Most voted balance: " + strconv.Itoa(topBalance) + ", voted: " + strconv.Itoa(voted))
  for idx, result := range results {
    if result.ErrCode != config.REMOTE_RESULT_ERROR_NONE {
			continue
		}
    r := result.Data.(*PaymentVerifierResponse)
		if (r.Message != "success") {
			continue
		}

		total := r.TotalReceived
		if (total != topBalance) {
			pownArray[idx] = config.RAIDA_STATUS_UNTRIED
			logger.Debug("Raida " + strconv.Itoa(idx) + " is reporting incorrect balance (" + strconv.Itoa(total) + "). Skipping it")
			continue
		}
	}

	pownString := v.GetPownStringFromStatusArray(pownArray)
	logger.Debug("Pownstring " + pownString)

	topMemo := "From RAIDAGO"
	pairsString := sortByCountString(memos)
	if len(pairsString) > 0 {
		topMemo = pairsString[0].Key
	}
	//voted := pairs[0].Value

	pairsString = sortByCountString(guids)
	if len(pairsString) == 0 {
		return "", &error.Error{config.ERROR_INVALID_RECEIPT_ID, "Failed to find GUID from payment_verify response"}
	}
	topGuid := pairsString[0].Key
	//voted := pairs[0].Value

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

	rx := NewStatement()
	tbString := strconv.Itoa(topBalance)
	response, err := rx.Create(topGuid, topMemo, tbString, cc.GetFileName(), cc)
	if err != nil {
		return "", &error.Error{config.ERROR_INVALID_RECEIPT_ID, "Failed to Create Statement"}
	}

	logger.Debug("Statement status " + response)

	vo := &PaymentVerifierOutput{}
	vo.AmountVerified = topBalance
	vo.Status = "success"
	vo.Message = "CloudCoins verified"


	b, err := json.Marshal(vo); 
	if err != nil {
		return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
	}

	//fmt.Printf("ns=%d %s isok=%b\n", v.Raida.TotalServers(), pownString, v.IsStatusArrayFixable(pownArray))
	return string(b), nil
}
