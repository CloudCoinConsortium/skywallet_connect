package raida

import (
	"fmt"
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"encoding/json"
	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"
	//	"fmt"
)

type ShowChange struct {
	Servant
}

type ShowChangeSubResponse struct {
	Sn      string `json:"sn"`
	Tag     string `json:"tag"`
	Created string `json:"created"`
}

type ShowChangeResponse struct {
	Server  string   `json:"server"`
	Version string   `json:"version"`
	Time    string   `json:"time"`
	Status  string   `json:"status"`
	D1      []string `json:"d1"`
	D5      []string `json:"d5"`
	D25     []string `json:"d25"`
	D100    []string `json:"d100"`
}

type ShowChangeOutput struct {
	AmountVerified int `json:"amount_verified"`
}

func NewShowChange() *ShowChange {
	return &ShowChange{
		*NewServant(),
	}
}

func (v *ShowChange) ShowChange(cm, snToBreak int) ([]int, *error.Error) {
	logger.Debug("ShowChange coin " + strconv.Itoa(snToBreak))

	seed, err := cloudcoin.GenerateHex(4)
	if err != nil {
		return nil, err
	}

	logger.Debug("Generated seed " + seed)

	pownArray := make([]int, v.Raida.TotalServers())
	params := make(map[string]string)
	params["nn"] = strconv.Itoa(config.DEFAULT_NN)
	params["sn"] = strconv.Itoa(config.PUBLIC_CHANGE_MAKER_ID)
	params["seed"] = seed
	params["denomination"] = strconv.Itoa(cloudcoin.GetDenomination(snToBreak))

	rsns1 := make([][]int, v.Raida.TotalServers())
	rsns5 := make([][]int, v.Raida.TotalServers())
	rsns25 := make([][]int, v.Raida.TotalServers())
	rsns100 := make([][]int, v.Raida.TotalServers())

	results := v.Raida.SendRequest("/service/show_change", params, ShowChangeResponse{})
	for idx, result := range results {
		rsns1[idx] = make([]int, 0)
		rsns5[idx] = make([]int, 0)
		rsns25[idx] = make([]int, 0)
		rsns100[idx] = make([]int, 0)

		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*ShowChangeResponse)
			if r.Status == "pass" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
				logger.Debug("Raida " + strconv.Itoa(idx) + " shows " + strconv.Itoa(len(r.D1)) + "," + strconv.Itoa(len(r.D5)) + "," + strconv.Itoa(len(r.D25)) + "," + strconv.Itoa(len(r.D100)))
				rsns1[idx] = make([]int, len(r.D1))
				rsns5[idx] = make([]int, len(r.D5))
				rsns25[idx] = make([]int, len(r.D25))
				rsns100[idx] = make([]int, len(r.D100))
				for sidx, ssn := range r.D1 {
					isn, err := strconv.Atoi(ssn)
					if err != nil {
						logger.Debug("Skipping invalid SN " + ssn)
						continue
					}
					rsns1[idx][sidx] = isn
				}

				for sidx, ssn := range r.D5 {
					isn, err := strconv.Atoi(ssn)
					if err != nil {
						logger.Debug("Skipping invalid SN " + ssn)
						continue
					}
					rsns5[idx][sidx] = isn
				}

				for sidx, ssn := range r.D25 {
					isn, err := strconv.Atoi(ssn)
					if err != nil {
						logger.Debug("Skipping invalid SN " + ssn)
						continue
					}
					rsns25[idx][sidx] = isn
				}

				for sidx, ssn := range r.D100 {
					isn, err := strconv.Atoi(ssn)
					if err != nil {
						logger.Debug("Skipping invalid SN " + ssn)
						continue
					}
					rsns100[idx][sidx] = isn
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
		return nil, &error.Error{config.ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC, "ShowChange results are not synchronized"}
	}

	csns1, _ := v.GetSNsOverlap(rsns1)
	csns5, _ := v.GetSNsOverlap(rsns5)
	csns25, _ := v.GetSNsOverlap(rsns25)
	csns100, _ := v.GetSNsOverlap(rsns100)

  fmt.Printf("v=%v\n", csns25)
	var sns []int

	switch cm {
	case config.CHANGE_METHOD_5A:
		sns = cloudcoin.CoinsGetA(csns1, 5)
		break
	case config.CHANGE_METHOD_25B:
		sns = cloudcoin.CoinsGet25B(csns5, csns1)
		break
	case config.CHANGE_METHOD_100E:
		sns = cloudcoin.CoinsGet100E(csns25, csns5, csns1)
		break
	case config.CHANGE_METHOD_250F:
		sns = cloudcoin.CoinsGet250F(csns100, csns25, csns5, csns1)
		break
	}

	logger.Debug("Total SNS got: " + strconv.Itoa(len(sns)))
	for _, sn := range sns {
		logger.Debug("Changed coin " + strconv.Itoa(sn) + " d:" + strconv.Itoa(cloudcoin.GetDenomination(sn)))
	}

	return sns, nil
}
