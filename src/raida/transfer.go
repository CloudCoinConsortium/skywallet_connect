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
)

type Transfer struct {
	Servant
}

type TransferResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
	Message string `json:"message"`
	Status string `json:"status"`
}

type TransferOutput struct {
	AmountSent int  `json:"amount_sent"`
	Message string
	Status string
}

func NewTransfer() (*Transfer) {
	return &Transfer{
		*NewServant(),
	}
}

func (v *Transfer) Transfer(cc *cloudcoin.CloudCoin, amount string, to string, memo string) (string, *error.Error) {
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return "", &error.Error{config.ERROR_INCORRECT_AMOUNT_SPECIFIED, "Invalid amount"}
	}

	if amountInt <= 0 {
		return "", &error.Error{config.ERROR_INCORRECT_AMOUNT_SPECIFIED, "Invalid amount"}
	}

	to_sn, err2 := cloudcoin.GuessSNFromString(to)
	if err2 != nil {
		return "", &error.Error{config.ERROR_INCORRECT_SKYWALLET, "Invalid Destination Address"}
	}

	logger.Debug("Started Transfer " + amount + " to " + to + " (" + strconv.Itoa(to_sn) + ") memo " + memo)

	s := NewShow()
	sns, total, err3 := s.ShowBrief(cc)
	if err3 != nil {
		logger.Error(err3.Message)
		return "", &error.Error{config.ERROR_SHOW_COINS_FAILED, "Failed to Show Coins"}
	}

	if total < amountInt {
		return "", &error.Error{config.ERROR_INSUFFICIENT_FUNDS, "Insufficient funds"}
	}

	nsns, extra, err3 := v.PickCoinsFromArray(sns, amountInt)
	if err3 != nil {
		logger.Debug("Failed to pick coins: " + err3.Message)
		return "", &error.Error{config.ERROR_PICK_COINS_AFTER_SHOW, "Failed to pick coins: " + err3.Message}
	}

	if extra != 0 {
		logger.Debug("Breaking extra coin: " + strconv.Itoa(extra))
		b := NewBreakInBank()
		csns, err := b.BreakInBank(cc, extra)
		if err != nil {
			return "", err
		}
/*
		for i := 0; i < len(nsns); i++ {
			fmt.Printf("vx0=%d %d\n",nsns[i],cloudcoin.GetDenomination(nsns[i]))
		}

		for i := 0; i < len(csns); i++ {
			fmt.Printf("vx1=%d %d\n",csns[i],cloudcoin.GetDenomination(csns[i]))
		}

		vsns := append(nsns, csns...)
		for i := 0; i < len(vsns); i++ {
			fmt.Printf("vx=%d %d\n",vsns[i],cloudcoin.GetDenomination(vsns[i]))
		}
*/
		vsns := append(nsns, csns...)
		var err4 *error.Error
		nsns, extra, err4 = v.PickCoinsFromArray(vsns, amountInt)
		if err4 != nil || extra != 0 {
			logger.Debug("Failed to pick coins after change: " + err4.Message)
			return "", &error.Error{config.ERROR_PICK_COINS_AFTER_CHANGE, "Failed to pick coins after change: " + err4.Message}
		}

	//	for i := 0; i < len(nsns); i++ {
	//		fmt.Printf("vxall=%d %d\n",nsns[i],cloudcoin.GetDenomination(nsns[i]))
	//	}
	}

	for _, sn := range nsns {
		logger.Debug("Sending " + strconv.Itoa(sn) + " d:" + strconv.Itoa(cloudcoin.GetDenomination(sn)))
	}


	var bufSns []int
	for _, sn := range nsns {
		bufSns = append(bufSns, sn)
		if (len(bufSns) == config.MAX_NOTES_TO_SEND) {
			if err := v.processTransfer(bufSns, cc, to_sn, memo); err != nil {
				return "", err
			}
			bufSns = nil
		}
	}

	if (len(bufSns) != 0) {
		if err := v.processTransfer(bufSns, cc, to_sn, memo); err != nil {
			return "", err
		}
	}

	//fmt.Printf("v=%d %d %v\n",total, extra, nsns)
	//results := v.Raida.SendRequest("/service/show", params, TransferResponse{})


  toutput := &TransferOutput{}
  toutput.AmountSent = amountInt
  toutput.Status = "success"
  toutput.Message = "CloudCoins sent"

  b, err := json.Marshal(toutput); 
  if err != nil {
    return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
  }

  return string(b), nil
}

func (v *Transfer) processTransfer(sns []int, cc *cloudcoin.CloudCoin, to int, memo string) *error.Error {
	logger.Debug("Processing " + strconv.Itoa(len(sns)) + " notes")

  stringSns := make([]string, len(sns))
  for idx, ssn := range sns {
    stringSns[idx] = strconv.Itoa(ssn)
  }
  ba, _ := json.Marshal(stringSns)

  pownArray := make([]int, v.Raida.TotalServers())
  params := make([]map[string]string, v.Raida.TotalServers())
  for idx, _ := range(params) {
    params[idx] = make(map[string]string)
    params[idx]["b"] = "t"
    params[idx]["nn"] = string(cc.Nn)
    params[idx]["sn"] = string(cc.Sn)
    params[idx]["an"] = cc.Ans[idx]
    params[idx]["pan"] = cc.Ans[idx]
    params[idx]["denomination"] = strconv.Itoa(cc.GetDenomination())
    params[idx]["to_sn"] = strconv.Itoa(to)
    params[idx]["tag"] = memo
    params[idx]["sns[]"] = string(ba)
  }

  results := v.Raida.SendDefinedRequestPost("/service/transfer", params, TransferResponse{})
  for idx, result := range results {
    if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
      r := result.Data.(*TransferResponse)
      if (r.Status == "allpass") {
        pownArray[idx] = config.RAIDA_STATUS_PASS
      } else if (r.Status == "allfail") {
        pownArray[idx] = config.RAIDA_STATUS_FAIL
      } else if (r.Status == "mixed") {
				// We need to tell that if there is one error the whole operation is treated as errornous
        pownArray[idx] = config.RAIDA_STATUS_ERROR
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
    return &error.Error{config.ERROR_TRANSFER_FAILED, "Failed to Transfer: " + pownString}
  }

	return nil
}
