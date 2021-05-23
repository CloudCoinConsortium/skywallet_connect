package raida

import (
	"encoding/json"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"encoding/json"
	//	"regexp"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"
	//"fmt"
)

type Pown struct {
	Servant
}

type PownResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Status  string `json:"status"`
}

type PownOutput struct {
	AmountAuthentic int `json:"amount_authentic"`
	AmountFracked int `json:"amount_fracked"`
	AmountCounterfeit int `json:"amount_counterfeit"`
	AmountLimbo int `json:"amount_limbo"`
	AmountTotal int `json:"amount_total"`
}

func NewPown() *Pown {
	return &Pown{
		*NewServant(),
	}
}

func (v *Pown) Pown() (string, *error.Error) {
	logger.Debug("Powning Coins")

  d := NewDetect()
  response, err := d.Detect()
	if err != nil {
	  return "", err
	}

  po := &PownOutput{}
  po.AmountAuthentic = response.AmountAuthentic
  po.AmountFracked = response.AmountFracked
  po.AmountLimbo = response.AmountLimbo
  po.AmountCounterfeit = response.AmountCounterfeit
  po.AmountTotal = response.AmountTotal
 
  
  f := NewFrackFixer()
  fresponse, ferr := f.Fix()
  if ferr != nil {
    logger.Error("Failed to fix")
  }

  po.AmountFracked -= fresponse.AmountFixed
  po.AmountAuthentic += fresponse.AmountFixed

	b, err2 := json.Marshal(po)
	if err2 != nil {
		return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
	}

	return string(b), nil
}
