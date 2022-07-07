package raida

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"

	//	"fmt"
	"strings"
)

type Destroyer struct {
	Servant
	AmountSent int
}

type DestroyResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DestroyerOutput struct {
	AmountSent int `json:"amount_sent"`
	Message    string
	Status     string
}

func NewDestroyer() *Destroyer {
	return &Destroyer{
		*NewServant(),
		0,
	}
}

func (v *Destroyer) Destroy(cc *cloudcoin.CloudCoin, sns []int) (string, *error.Error) {

	logger.Debug("Started Destroyer " + strconv.Itoa(len(sns)) + " (" + string(cc.Sn) + ")")

  bufCCs := make([]int, 0)
	for _, sn := range sns {
		bufCCs = append(bufCCs, sn)
		if len(bufCCs) == config.MAX_NOTES_TO_SEND {
			if err := v.processDestroy(cc, bufCCs); err != nil {
				return "", err
			}
			bufCCs = nil
		}
	}

	if len(bufCCs) != 0 {
		if err := v.processDestroy(cc, bufCCs); err != nil {
			return "", err
		}
	}

	return "ok", nil
}

func (v *Destroyer) processDestroy(cc *cloudcoin.CloudCoin, sns []int) *error.Error {
	logger.Debug("Processing " + strconv.Itoa(len(sns)) + " notes")

	stringSns := make([]string, len(sns))
	for idx, sn := range sns {
		stringSns[idx] = strconv.Itoa(sn)
	}

	baSn, _ := json.Marshal(stringSns)

	pownArray := make([]int, v.Raida.TotalServers())
	params := make([]map[string]string, v.Raida.TotalServers())
	for idx, _ := range params {
		params[idx] = make(map[string]string)
		params[idx]["an"] = cc.Ans[idx]
		params[idx]["sn"] = string(cc.Sn)
		params[idx]["nn"] = "1"
		params[idx]["sns[]"] = string(baSn)
	}

	results := v.Raida.SendDefinedRequestPost("/service/destroy", params, DestroyResponse{})
	for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*DestroyResponse)
			if r.Status == "allpass" {
				pownArray[idx] = config.RAIDA_STATUS_PASS
			} else if r.Status == "allfail" {
				pownArray[idx] = config.RAIDA_STATUS_FAIL
			} else if r.Status == "mixed" {
				ss := strings.Split(r.Message, ",")
				if len(ss) != len(sns) {
					logger.Error("Invlid length returned from raida: " + string(len(ss)) + ", expected: " + string(len(sns)))
				} else {
					for _, status := range ss {
						logger.Debug("sn=" + status)

						if status == "pass" {
						} else if status == "fail" {
							// addCoinTorarr
						} else {
						}
					}

				}

			} else {
				pownArray[idx] = config.RAIDA_STATUS_ERROR
			}
		} else if result.ErrCode == config.REMOTE_RESULT_ERROR_TIMEOUT {
			pownArray[idx] = config.RAIDA_STATUS_NORESPONSE
		} else {
			pownArray[idx] = config.RAIDA_STATUS_ERROR
		}
	}

  st := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(sns)), "\n"), "[]")
  fname := "/root/destroyedcoins/f-" + time.Now().String() + ".txt"

  os.WriteFile(fname, []byte(st), 0644)











  fmt.Printf("pa=%v\n", pownArray)

	return nil
}
