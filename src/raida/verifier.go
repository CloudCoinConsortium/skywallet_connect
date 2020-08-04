package raida

import (
	"logger"
	"config"
	"strconv"
	"encoding/json"
	"sort"
	"regexp"
	"cloudcoin"
)

type Verifier struct {
	Servant
}

type VerifierResponse struct {
  Server  string `json:"server"`
	Version string `json:"version"`
	Time  string `json:"time"`
	TotalReceived int `json:"total_received"`
	Message string `json:"message"`
	SerialNumbers string `json:"serial_numbers"`
}

type VerifierOutput struct {
	AmountVerified int  `json:"amount_verified"`
}

func NewVerifier() (*Verifier) {
	return &Verifier{
		*NewServant(),
	}
}

func (v *Verifier) Receive(uuid string, owner string) (string, *Error) {
	logger.Debug("Started Verifier with UUID " + uuid + " owner " + owner)

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

	results := v.Raida.SendRequest("/service/view_receipt", params, VerifierResponse{})
  for idx, result := range results {
    if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
      r := result.Data.(*VerifierResponse)
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
    r := result.Data.(*VerifierResponse)
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

/*
	for key, element := range balances {
		fmt.Printf("k=%d v=%d\n", key, element)
	}
*/
	pownString := v.GetPownStringFromStatusArray(pownArray)
	logger.Debug("Pownstring " + pownString)

	if !v.IsStatusArrayFixable(pownArray) {
		return "", &Error{"Results from the RAIDA are not synchronized"}
	}

	vo := &VerifierOutput{}
	vo.AmountVerified = topBalance

	b, err := json.Marshal(vo); 
	if err != nil {
		return "", &Error{"Failed to Encode JSON"}
	}

	//fmt.Printf("ns=%d %s isok=%b\n", v.Raida.TotalServers(), pownString, v.IsStatusArrayFixable(pownArray))
	return string(b), nil
}


func sortByCount(totals map[int]int) PairList {
	pl := make(PairList, len(totals))
	i := 0

	for k, v := range totals {
		pl[i] = Pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key int
	Value int
}

type PairList []Pair
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)  { p[i], p[j] = p[j], p[i] }
