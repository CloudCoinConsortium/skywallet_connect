package raida

import (
	"encoding/json"
	"math"
	"net"
	neturl "net/url"
	"regexp"
	"strconv"
	"strings"
	"io/ioutil"
  "net/http"
  "time"
  "encoding/base64"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"


	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"
)

type Transfer struct {
	Servant
}

type TransferResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type InfoResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type TransferOutput struct {
	AmountSent int `json:"amount_sent"`
	Message    string
	Status     string
}

type ResultError struct {
  err *error.Error
}

func NewTransfer() *Transfer {
	return &Transfer{
		*NewServant(),
	}
}

func (v *Transfer) TechTransfer(cc *cloudcoin.CloudCoin, amount string, to string, memo string) (*TransferOutput, *error.Error) {
  amountf, ferr := strconv.ParseFloat(amount, 64)
  if ferr != nil {
		return nil, &error.Error{config.ERROR_INCORRECT_AMOUNT_SPECIFIED, "Invalid amount"}
  }

  amountInt := int(math.Round(amountf))
	if amountInt <= 0 {
		return nil, &error.Error{config.ERROR_INCORRECT_AMOUNT_SPECIFIED, "Invalid amount"}
	}

  logger.Debug("Trasnfer rounded amount " + strconv.Itoa(amountInt))

	to_sn, err2 := cloudcoin.GuessSNFromString(to)
	if err2 != nil {
		return nil, &error.Error{config.ERROR_INCORRECT_SKYWALLET, "Invalid Destination Address"}
	}

	logger.Debug("Started Transfer " + amount + " to " + to + " (" + strconv.Itoa(to_sn) + ") memo " + memo)
	tags := v.GetObjectMemo("", memo, amount, cc.GetFileName())

	s := NewShow()
	sns, total, err3 := s.ShowBrief(cc)
	if err3 != nil {
		logger.Error(err3.Message)
		return nil, &error.Error{config.ERROR_SHOW_COINS_FAILED, "Failed to Show Coins"}
	}

	if total < amountInt {
		return nil, &error.Error{config.ERROR_INSUFFICIENT_FUNDS, "Insufficient funds"}
	}

	nsns, extra, err3 := v.PickCoinsFromArray(sns, amountInt)
	if err3 != nil {
		logger.Debug("Failed to pick coins: " + err3.Message)
		return nil, &error.Error{config.ERROR_PICK_COINS_AFTER_SHOW, "Failed to pick coins: " + err3.Message}
	}

	if extra != 0 {
		logger.Debug("Breaking extra coin: " + strconv.Itoa(extra))
		b := NewBreakInBank()
		csns, err := b.BreakInBank(cc, extra)
		if err != nil {
			return nil, err
		}

		vsns := append(nsns, csns...)
		var err4 *error.Error
		nsns, extra, err4 = v.PickCoinsFromArray(vsns, amountInt)
		if err4 != nil || extra != 0 {
			logger.Debug("Failed to pick coins after change: " + err4.Message)
			return nil, &error.Error{config.ERROR_PICK_COINS_AFTER_CHANGE, "Failed to pick coins after change: " + err4.Message}
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
		if len(bufSns) == config.MAX_NOTES_TO_SEND {
			if err := v.processTransfer(bufSns, cc, to_sn, tags); err != nil {
				return nil, err
			}
			bufSns = nil
		}
	}

	if len(bufSns) != 0 {
		if err := v.processTransfer(bufSns, cc, to_sn, tags); err != nil {
			return nil, err
		}
	}

	//fmt.Printf("v=%d %d %v\n",total, extra, nsns)
	//results := v.Raida.SendRequest("/service/show", params, TransferResponse{})

	toutput := &TransferOutput{}
	toutput.AmountSent = amountInt
	toutput.Status = "success"
	toutput.Message = "CloudCoins sent"

  return toutput, nil
}


func (v *Transfer) Transfer(cc *cloudcoin.CloudCoin, amount string, to string, memo string) (string, *error.Error) {
  toutput, err := v.TechTransfer(cc, amount, to, memo)
  if err != nil {
    return "", err
  }

	b, err2 := json.Marshal(*toutput)
	if err2 != nil {
		return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
	}

	return string(b), nil
}

func (v *Transfer) processTransfer(sns []int, cc *cloudcoin.CloudCoin, to int, tags []string) *error.Error {
	logger.Debug("Processing " + strconv.Itoa(len(sns)) + " notes")

	stringSns := make([]string, len(sns))
	for idx, ssn := range sns {
		stringSns[idx] = strconv.Itoa(ssn)
	}
	ba, _ := json.Marshal(stringSns)

	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range params {
		params[idx] = make(map[string]string)
		params[idx]["b"] = "t"
		params[idx]["nn"] = string(cc.Nn)
		params[idx]["sn"] = string(cc.Sn)
		params[idx]["an"] = cc.Ans[idx]
		params[idx]["pan"] = cc.Ans[idx]
		params[idx]["denomination"] = strconv.Itoa(cc.GetDenomination())
		params[idx]["to_sn"] = strconv.Itoa(to)
		params[idx]["tag"] = tags[idx]
		params[idx]["sns[]"] = string(ba)
	}

	results := v.Raida.SendDefinedRequestPost("/service/transfer", params, TransferResponse{})
	for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*TransferResponse)
			if r.Status == "allpass" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
			} else if r.Status == "allfail" || r.Status == "fail" {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
			} else if r.Status == "mixed" {
				// We need to tell that if there is one error the whole operation is treated as errornous
				pownArray[idx] = config.RAIDA_STATUS_ERROR
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
		return &error.Error{config.ERROR_TRANSFER_FAILED, "Failed to Transfer: " + pownString}
	}

	return nil
}

func (v *Transfer) Pay(cc *cloudcoin.CloudCoin, amount string, guid string, to string, memo string) (string, *error.Error) {
	matched, err := regexp.MatchString(`^[A-Fa-f0-9]{32}$`, guid)
	if err != nil || !matched {
		return "", &error.Error{config.ERROR_INVALID_GUID, "GUID invalid or not defined"}
	}
  
  sender_address := cc.GetFileName()

  txts, err := net.LookupTXT(to)
  if err != nil {
		return "", &error.Error{config.ERROR_DNS, "Failed to get TXT record"}
  }

  if len(txts) != 1 {
		return "", &error.Error{config.ERROR_DNS, "Destination Skywallet must have exactly one TXT record"}
  }

  txt := txts[0]
  txt = strings.ReplaceAll(txt, "\"", "")


  urls := make([]string, v.Raida.TotalNumber)
  var url string
  needMultiple := false
  if strings.Contains(txt, "%") {
    for idx := 0; idx < v.Raida.TotalNumber; idx++ {
      strIdx := strconv.Itoa(idx)
      urls[idx] = strings.ReplaceAll(txt, "%n", "" + strIdx)
    }
    url = urls[0]
    needMultiple = true
  } else {
    url = txt
  }
  
  _, err = neturl.ParseRequestURI(url)
  if err != nil {
		return "", &error.Error{config.ERROR_INVALID_URL, "URL in the TXT record of the receiver is incorrect"}
  }

  transferMemo := guid
  toutput, err3 := v.TechTransfer(cc, amount, to, transferMemo)
  if err3 != nil {
    return "", err3
  }

  logger.Debug("Transfer Completed")

  ch := make(chan ResultError)
  errors := 0

  meta := ""
  meta += "from = \"" + sender_address + "\"\n"
  meta += "message = \"" + memo + "\"\n"

  meta = base64.StdEncoding.EncodeToString([]byte(meta))

  params := neturl.Values{}
  params.Add("merchant_skywallet", to)
  params.Add("guid", guid)
  params.Add("meta", meta)
  rquery := params.Encode()

  if (needMultiple) {

    for _, u := range urls {
      logger.Debug("Sending info request " + u)
      go v.makeRequest(u + "?" + rquery, ch)
    }

    
    logger.Debug("Sent all requests")
    i := 0
    for range urls {
      res := <- ch
      if res.err != nil {
        errors++
        logger.Debug("error")
      } else {
        logger.Debug("ok")
      }
      i++

    }

    if errors > 13 {
      return "", &error.Error{config.ERROR_COINS_SENT_BUT_INFO_REQUESTS_FAILED, "Coins sent. However, there were too many errors from MerchantURL"}
    }

//    fmt.Printf("xxxxxxxxx sent %d", resp.StatusCode)
/*
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      fmt.Printf("errr2 %s\n", err.Error())
    }

    sb := string(body)
    */

  } else {

  }

  //fmt.Printf("t=%s\n", toutput.Status)
  
	b, err4 := json.Marshal(*toutput)
	if err4 != nil {
		return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
	}

	return string(b), nil
}


func (v *Transfer) makeRequest(url string, ch chan<-ResultError) {

  client := http.Client{
    Timeout: config.INFO_HTTP_TIMEOUT * time.Second,
  }

  resp, err := client.Get(url)
  if err != nil {
    logger.Error("Failed to make request " + url)
    ch <- ResultError{err: &error.Error{config.ERROR_HTTP, "Failed to make request: " + err.Error()}}
    return
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    logger.Error("Failed to read body " + url)
    ch <- ResultError{err: &error.Error{config.ERROR_HTTP, "Failed to read body: " + err.Error()}}
    return
  }

  logger.Debug("body " + string(body))

  if resp.StatusCode != 200 {
    logger.Error("Invalid response code " + url + ": " + strconv.Itoa(resp.StatusCode))
    ch <- ResultError{err: &error.Error{config.ERROR_COINS_SENT_BUT_INFO_REQUESTS_FAILED, "Invalid response code: " + strconv.Itoa(resp.StatusCode)}}
    return
  }

  bytes := []byte(body)

  var ir InfoResponse

  err3 := json.Unmarshal(bytes, &ir)
  if err3 != nil {
    logger.Error("Failed to parse body " + url + ", b: " + string(body))
    ch <- ResultError{err: &error.Error{config.ERROR_HTTP, "Failed to parse body: " + err3.Error()}}
    return
  }

  if ir.Status != "success" {
    logger.Error("Invalid status " + url + ", b: " + ir.Status)
    ch <- ResultError{err: &error.Error{config.ERROR_COINS_SENT_BUT_INFO_REQUESTS_FAILED, "Invalid status code: " + ir.Status}}
    return
  }

  ch <- ResultError{err: nil}

}

