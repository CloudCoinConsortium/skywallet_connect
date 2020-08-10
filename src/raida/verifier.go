package raida

import (
	"logger"
	"config"
	"strconv"
	"encoding/json"
	"sort"
	"regexp"
	"cloudcoin"
	"error"
	"strings"
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
	Status string `json:"status"`
	Message string `json:"message"`
}

func NewVerifier() (*Verifier) {
	return &Verifier{
		*NewServant(),
	}
}

func (v *Verifier) Receive(uuid string, owner string) (string, *error.Error) {
	logger.Debug("Started Verifier with UUID " + uuid + " owner " + owner)

	matched, err := regexp.MatchString(`^[A-Fa-f0-9]{32}$`, uuid)
	if err != nil || !matched {
		return "", &error.Error{"UUID invalid or not defined"}
	}

	sn, err2 := cloudcoin.GuessSNFromString(owner)
	if err2 != nil {
		return "", &error.Error{"Invalid Owner"}
	}

	logger.Debug("owner SN " +  strconv.Itoa(sn))

	pownArray := make([]int, v.Raida.TotalServers())
	balances := make(map[int]int)

	allSns := make(map[int]bool)

	params := make(map[string]string)
	params["tag"] = uuid
	params["owner"] = strconv.Itoa(sn)

	fArr := make([]map[int]bool, v.Raida.TotalServers())

	results := v.Raida.SendRequest("/service/view_receipt", params, VerifierResponse{})
  for idx, result := range results {
    if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
      r := result.Data.(*VerifierResponse)
			if (r.Message == "success") {
				pownArray[idx] = config.RAIDA_STATUS_PASS
				total := r.TotalReceived
				balances[total]++
				logger.Debug("raida " + strconv.Itoa(idx) + " total " + strconv.Itoa(total))

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
			logger.Debug("Raida " + strconv.Itoa(idx) + " is reporting incorrect balance (" + strconv.Itoa(total) + "). Skipping it")
			continue
		}
	}

	pownString := v.GetPownStringFromStatusArray(pownArray)
	logger.Debug("Pownstring " + pownString)

	if !v.IsStatusArrayFixable(pownArray) {
		return "", &error.Error{"Results from the RAIDA are not synchronized"}
	}

	vo := &VerifierOutput{}
	vo.AmountVerified = topBalance
	vo.Status = "success"
	vo.Message = "CloudCoins verified"

	// Getting Sns for fixing
	ft := NewFixTransfer()
	for ssn, _ := range allSns {
		for ridx, farr := range fArr {
			if _, ok := farr[ssn]; !ok {
				logger.Debug("Coin " + strconv.Itoa(ssn) + " will be fixed on raida " + strconv.Itoa(ridx))
				ft.AddSNToRepairArray(ridx, ssn)
			}
		}
	}

	ft.FixTransfer()

	b, err := json.Marshal(vo); 
	if err != nil {
		return "", &error.Error{"Failed to Encode JSON"}
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
