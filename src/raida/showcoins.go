package raida

import (
	"logger"
	"config"
	"encoding/json"
//	"regexp"
	"error"
//	"fmt"
	"core"
	"cloudcoin"
)

type ShowCoins struct {
	Servant
}


type ShowCoinsOutput struct {
	D1 int  `json:"d1"`
	D5 int  `json:"d5"`
	D25 int  `json:"d25"`
	D100 int  `json:"d100"`
	D250 int  `json:"d250"`
	Total int `json:"total"`
}

func NewShowCoins() (*ShowCoins) {
	return &ShowCoins{
		*NewServant(),
	}
}

func (v *ShowCoins) ShowCoins() (string, *error.Error) {
	logger.Debug("ShowCoins")


	sns, err := core.GetSNSFromFolder(core.GetBankDir())
	if err != nil {
    return "", err
	}
	snsf, err2 := core.GetSNSFromFolder(core.GetFrackedDir())
	if err != nil {
    return "", err2
	}

	for k, v := range snsf {
		sns[k] = v
	}

  soutput := &ShowCoinsOutput{}
	for k, _ := range sns {
		d :=  cloudcoin.GetDenomination(k)

		switch d {
		case 1:
			soutput.D1 += 1
		case 5:
			soutput.D5 += 1
		case 25:
			soutput.D25 += 1
		case 100:
			soutput.D100 += 1
		case 250:
			soutput.D250 += 1
		}

		soutput.Total += d
	}

  b, err3 := json.Marshal(soutput); 
  if err3 != nil {
    return "", &error.Error{config.ERROR_ENCODE_JSON, "Failed to Encode JSON"}
  }


  return string(b), nil


}
