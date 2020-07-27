package raida

import (
	"fmt"
	"logger"
	"config"
	"strconv"
	"encoding/json"
)

type Verifier struct {
	Servant
}

type VerifierResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
}

type VerifierOutput struct {
	AmountVerified int  `json:"amount_verified"`
}

func NewVerifier() (*Verifier) {
	return &Verifier{
		*NewServant(),
	}
}

func (v *Verifier) Receive(uuid string) (string, *Error) {
	logger.Debug("Started Verifier with UUID " + uuid)

	pownArray := make([]int, v.Raida.TotalServers())
	balances := make(map[int]int)

	results := v.Raida.SendRequest("/service/view_receipt", VerifierResponse{})
  for idx, result := range results {
    if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
      r := result.Data.(*VerifierResponse)
      fmt.Println("raida"+strconv.Itoa(idx) + " r="+result.Message + " cr="+strconv.Itoa(result.ErrCode) + " srv="+r.Server)
			pownArray[idx] = config.RAIDA_STATUS_PASS

    } else if (result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT) {
			pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
			balances[0]++
		} else {
			pownArray[idx] = config.RAIDA_STATUS_ERROR
			balances[0]++
		}
  }

	for key, element := range balances {
		fmt.Printf("k=%d v=%d\n", key, element)
	}

	pownString := v.GetPownStringFromStatusArray(pownArray)
	logger.Debug("Pownstring " + pownString)

	vo := &VerifierOutput{}
	vo.AmountVerified = 10

	b, err := json.Marshal(vo); 
	if err != nil {
		return "", &Error{"Failed to Encode JSON"}
	}

	//fmt.Printf("ns=%d %s isok=%b\n", v.Raida.TotalServers(), pownString, v.IsStatusArrayFixable(pownArray))
	return string(b), nil
}
