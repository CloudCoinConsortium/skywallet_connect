package raida

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/core"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"encoding/json"
	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"

	//"fmt"
	"strings"
)

type Detect struct {
	Servant
}

type DetectResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Status  string `json:"status"`
  Message string `json:"message"`
  Ticket  string `json:"ticket"`
}

type DetectOutput struct {
	AmountAuthentic int `json:"amount_authentic"`
	AmountFracked int `json:"amount_fracked"`
	AmountCounterfeit int `json:"amount_counterfeit"`
	AmountLimbo int `json:"amount_limbo"`
	AmountErrors int `json:"amount_errors"`
	AmountTotal int `json:"amount_total"`
}

func NewDetect() *Detect {
	return &Detect{
		*NewServant(),
	}
}

func (v *Detect) Detect() (*DetectOutput, *error.Error) {
//	if !cloudcoin.ValidateCoin(cc) {
//		return nil, 0, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "CloudCoin is invalid"}
//	}

	logger.Debug("Detecting coins")

  ccs, err := core.GetCoinsForImport()
  if err != nil {
    return nil, err
  }

  total := 0
  for _, cc := range(*ccs) {
    fmt.Printf("cc=%s %d %d %d\n", cc.Path, cc.Sn, cc.Nn, cc.GetDenomination())
    total += cc.GetDenomination()
    cc.GenerateMyPans()
  }

  do := &DetectOutput{}
  do.AmountAuthentic = 0
  do.AmountCounterfeit = 0
  do.AmountFracked = 0
  do.AmountLimbo = 0
  do.AmountTotal = total

  if len(*ccs) == 0 {
    return do, nil
  }

	var bufCcs []cloudcoin.CloudCoin
  locTotal := 0
	for _, cc := range(*ccs) {
    locTotal += cc.GetDenomination()
		bufCcs = append(bufCcs, cc)
		if len(bufCcs) == config.MAX_NOTES_TO_SEND {
			do2, err := v.processDetect(bufCcs)
      if err != nil {
        logger.Debug("Error " + err.Message)
        do.AmountErrors += locTotal
			}

      do.AmountAuthentic += do2.AmountAuthentic
      do.AmountCounterfeit += do2.AmountCounterfeit
      do.AmountFracked += do2.AmountFracked
      do.AmountLimbo += do2.AmountLimbo
			bufCcs = nil
      locTotal = 0
		}
	}

	if len(bufCcs) != 0 {
		do2, err := v.processDetect(bufCcs)
    if err != nil {
      logger.Debug("Error " + err.Message)
      do.AmountErrors += locTotal
		}
    do.AmountAuthentic += do2.AmountAuthentic
    do.AmountCounterfeit += do2.AmountCounterfeit
    do.AmountFracked += do2.AmountFracked
    do.AmountLimbo += do2.AmountLimbo
	}

  return do, nil
}




func (v *Detect) processDetect(ccs []cloudcoin.CloudCoin) (*DetectOutput, *error.Error) {
  do := &DetectOutput{}

  logger.Debug("Detecting " + strconv.Itoa(len(ccs)))

	stringSns := make([]string, len(ccs))
	stringNns := make([]string, len(ccs))
	stringDns := make([]string, len(ccs))
	for idx, cc := range ccs {
		stringSns[idx] = string(cc.Sn)
		stringNns[idx] = string(cc.Nn)
		stringDns[idx] = strconv.Itoa(cc.GetDenomination())
	}

	preParams := make([][]string, v.Raida.TotalServers())
	for ridx, _ := range preParams {
		preParams[ridx] = make([]string, len(ccs))
		for idx, cc := range ccs {
			preParams[ridx][idx] = string(cc.Ans[ridx])
		}
	}

	preParamsPans := make([][]string, v.Raida.TotalServers())
	for ridx, _ := range preParamsPans {
		preParamsPans[ridx] = make([]string, len(ccs))
		for idx, cc := range ccs {
			preParamsPans[ridx][idx] = string(cc.Pans[ridx])
		}
	}
	baSn, _ := json.Marshal(stringSns)
	baNn, _ := json.Marshal(stringNns)
	baDn, _ := json.Marshal(stringDns)


	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range params {
		baAns, _ := json.Marshal(preParams[idx])
		baPans, _ := json.Marshal(preParamsPans[idx])
		params[idx] = make(map[string]string)
		params[idx]["b"] = "t"
		params[idx]["nns[]"] = string(baNn)
		params[idx]["sns[]"] = string(baSn)
		params[idx]["denomination[]"] = string(baDn)
		params[idx]["ans[]"] = string(baAns)
		params[idx]["pans[]"] = string(baPans)
	}

	results := v.Raida.SendDefinedRequestPost("/service/multi_detect", params, DetectResponse{})
	for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*DetectResponse)
			if r.Status == "allpass" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
				v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_PASS)
			} else if r.Status == "allfail" {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
				v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_FAIL)
			} else if r.Status == "mixed" {
				ss := strings.Split(r.Message, ",")
				if len(ss) != len(ccs) {
					logger.Error("Invalid length returned from raida: " + string(len(ss)) + ", expected: " + string(len(ccs)))
					v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_ERROR)
				} else {
					for aIdx, status := range ss {
						logger.Debug("sn=" + status)

						if status == "pass" {
							v.SetCoinStatusInArray(ccs, aIdx, idx, config.RAIDA_STATUS_PASS)
						} else if status == "fail" {
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
		} else if result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT {
			pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
			v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_NORESPONSE)
		} else {
			pownArray[idx] = config.RAIDA_STATUS_ERROR
			v.SetCoinsStatus(ccs, idx, config.RAIDA_STATUS_ERROR)
		}
	}

	for _, cc := range ccs {
    cc.SetAnsToPansIfPassed()
		logger.Debug("cc " + string(cc.Sn) + " pownstring " + cc.GetPownString())
    isAuthentic, hasFailed, isCounterfeit := cc.IsAuthentic()
    if (isAuthentic) {
      if (!hasFailed) {
    			logger.Debug("Coin is Authentic " + string(cc.Sn))
    			core.MoveCoinNewContent(cc, config.DIR_BANK)
          do.AmountAuthentic += cc.GetDenomination()
        } else {
    			logger.Debug("Coin is Fracked " + string(cc.Sn))
    			core.MoveCoinNewContent(cc, config.DIR_FRACKED)
          do.AmountFracked += cc.GetDenomination()
        }
    } else {
      if (isCounterfeit) {
    			logger.Debug("Coin is Counterfeit " + string(cc.Sn))
    			core.MoveCoinToCounterfeit(cc)
          do.AmountCounterfeit += cc.GetDenomination()
      } else {
    			logger.Debug("Coin is Limbo " + string(cc.Sn))
    			core.MoveCoinNewContent(cc, config.DIR_LIMBO)
          do.AmountLimbo += cc.GetDenomination()
      }
		}
	}

  return do, nil
}

