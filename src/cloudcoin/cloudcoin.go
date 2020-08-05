package cloudcoin


import (
	"logger"
//	"encoding/json"
//	"os"
//	"io/ioutil"
//	"strconv"
//"config"
//	"error"
)

type CloudCoin struct {
	Nn string `json:"nn"`
	Sn string `json:"sn"`
	Ans []string `json:"an"`
	Pans []string `json:"pan"`
}

type CloudCoinStack struct {
	Stack []CloudCoin `json:"cloudcoin"`
}

func New(path string) *CloudCoin {

	//ans := make([]string, 0)
	//pans := make([]string, 0)

	ccStack, err := ReadFromFile(path)
	if err != nil {
		return nil
	}

	if len(ccStack.Stack) != 1 {
		logger.Error("Stack File must contain only one Coin")
		return nil
	}

	cc := ccStack.Stack[0]
	if !ValidateCoin(cc) {
		logger.Error("Coin is corrupted")
		return nil
	}

	return &cc
}
