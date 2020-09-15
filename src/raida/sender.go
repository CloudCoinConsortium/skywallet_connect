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
	"core"
	"strings"
)

type Sender struct {
	Servant
	AmountSent int
}

type SendResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
	Message string `json:"message"`
	Status string `json:"status"`
}

type SenderOutput struct {
	AmountSent int  `json:"amount_sent"`
	Message string
	Status string
}

func NewSender() (*Sender) {
	return &Sender{
		*NewServant(),
		0,
	}
}

func (v *Sender) Send(amount string, to string, memo string) (string, *error.Error) {
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

	logger.Debug("Started Sender " + amount + " to " + to + " (" + strconv.Itoa(to_sn) + ") memo " + memo)


	var ccs []cloudcoin.CloudCoin
	var extraCC *cloudcoin.CloudCoin
	var err3, err4 *error.Error

	ccs, extraCC, err3 = v.GetCoinsFromDirs(amountInt)
	if err3 != nil {
		return "", err3
	}

	if extraCC != nil {
		logger.Debug("Need to break coin " + string(extraCC.Sn))

		b := NewBreak()
		err := b.Break(extraCC)
		if err != nil {
			return "", err
		}

		ccs, extraCC, err4 = v.GetCoinsFromDirs(amountInt)
		if err4 != nil {
			return "", err4
		}

		if extraCC != nil {
			logger.Debug("Failed to pick coins after change")
		  return "", &error.Error{config.ERROR_PICK_COINS_AFTER_CHANGE, "Failed to pick coins after change"}
		}


	//	return "", &error.Error{config.ERROR_INCORRECT_AMOUNT_SPECIFIED, "Change needed but it is not supported yet"}
	}

/*
	if extra != 0 {
		logger.Debug("Breaking extra coin: " + strconv.Itoa(extra))
		b := NewBreakInBank()
		csns, err := b.BreakInBank(cc, extra)
		if err != nil {
			return "", err
		}
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
`*/
	for _, cc := range ccs {
		logger.Debug("Sending " + string(cc.Sn) + " d:" + strconv.Itoa(cc.GetDenomination()))
	}


	var bufCCs []cloudcoin.CloudCoin
	for _, sn := range ccs {
		bufCCs = append(bufCCs, sn)
		if (len(bufCCs) == config.MAX_NOTES_TO_SEND) {
			if err := v.processSend(bufCCs, to_sn, memo); err != nil {
				return "", err
			}
			bufCCs = nil
		}
	}

	if err := v.processSend(bufCCs, to_sn, memo); err != nil {
		return "", err
	}


	//fmt.Printf("v=%d %d %v\n",total, extra, nsns)
	//results := v.Raida.SendRequest("/service/show", params, SenderResponse{})


  toutput := &SenderOutput{}
  toutput.AmountSent = v.AmountSent
  toutput.Status = "success"
  toutput.Message = "CloudCoins sent"

  b, err := json.Marshal(toutput); 
  if err != nil {
    return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
  }
  return string(b), nil
}

func (v *Sender) processSend(ccs []cloudcoin.CloudCoin, to int, memo string) *error.Error {
	logger.Debug("Processing " + strconv.Itoa(len(ccs)) + " notes")

  stringSns := make([]string, len(ccs))
  stringNns := make([]string, len(ccs))
  stringDns := make([]string, len(ccs))
  for idx, cc := range ccs {
    stringSns[idx] = string(cc.Sn)
    stringNns[idx] = string(cc.Nn)
    stringDns[idx] = strconv.Itoa(cc.GetDenomination())
  }
  preParams := make([][]string, v.Raida.TotalServers())
  for ridx, _ := range(preParams) {
		preParams[ridx] = make([]string, len(ccs))
	  for idx, cc := range ccs {
			preParams[ridx][idx] = string(cc.Ans[ridx])
		}
	}
  baSn, _ := json.Marshal(stringSns)
  baNn, _ := json.Marshal(stringNns)
  baDn, _ := json.Marshal(stringDns)
	

  pownArray := make([]int, v.Raida.TotalServers())
  params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range(params) {
		baAns, _ := json.Marshal(preParams[idx])
		params[idx] = make(map[string]string)
		params[idx]["b"] = "t"
		params[idx]["to_sn"] = strconv.Itoa(to)
	  params[idx]["tag"] = memo
	  params[idx]["nns[]"] = string(baNn)
		params[idx]["sns[]"] = string(baSn)
		params[idx]["denomination[]"] = string(baDn)
		params[idx]["ans[]"] = string(baAns)
  }


  results := v.Raida.SendDefinedRequestPost("/service/send", params, SendResponse{})
  for idx, result := range results {
    if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
      r := result.Data.(*SendResponse)
      if (r.Status == "allpass") {
        pownArray[idx] = config.RAIDA_STATUS_PASS
				v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_PASS)
      } else if (r.Status == "allfail") {
        pownArray[idx] = config.RAIDA_STATUS_FAIL
				v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_FAIL)
      } else if (r.Status == "mixed") {
				ss := strings.Split(r.Message, ",")
				if len(ss) != len(ccs) {
					logger.Error("Invlid length returned from raida: " + string(len(ss)) + ", expected: " + string(len(ccs)))
					v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_ERROR)
				} else {
					for aIdx, status := range ss {
						logger.Debug("sn=" + status)

						if (status == "pass") {
							v.SetCoinStatusInArray(ccs, aIdx, idx, config.RAIDA_STATUS_PASS)
						} else if (status == "fail") {
							v.SetCoinStatusInArray(ccs, aIdx, idx, config.RAIDA_STATUS_FAIL)
							// addCoinTorarr
						} else {
							v.SetCoinStatusInArray(ccs, aIdx, idx, config.RAIDA_STATUS_ERROR)
						}
					}

				}

			} else {
        pownArray[idx] = config.RAIDA_STATUS_ERROR
				v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_ERROR)
      }
    } else if (result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT) {
        pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
				v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_NORESPONSE)
    } else {
        pownArray[idx] = config.RAIDA_STATUS_ERROR
				v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_ERROR)
    }
	}

	for _, cc := range ccs {
		logger.Debug("cc " + string(cc.Sn) + " pownstring " + cc.GetPownString())
		if v.IsStatusArrayFixable(cc.Statuses) {
			logger.Debug("Coin was sent successfully")
			core.MoveCoinToSent(cc)
			v.AmountSent += cc.GetDenomination()
		} else {
			logger.Debug("Coin is counterfeit")
			core.MoveCoinToCounterfeit(cc)
		}
	//	core.MoveCoinToCounterfeit(cc)
	}

 /* pownString := v.GetPownStringFromStatusArray(pownArray)
  logger.Debug("Send Pownstring " + pownString)

  if !v.IsStatusArrayFixable(pownArray) {
    return &error.Error{config.ERROR_TRANSFER_FAILED, "Failed to Send: " + pownString}
  }
*/
	return nil
}
