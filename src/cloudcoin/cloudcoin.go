package cloudcoin


import (
	"logger"
	"encoding/json"
//	"os"
//	"io/ioutil"
	"strconv"
//"config"
	"error"
	"strings"
)

type CloudCoin struct {
	Nn json.Number `json:"nn"`
	Sn json.Number `json:"sn"`
	Ans []string `json:"an"`
	Pans []string `json:"pan"`
}

type CloudCoinStack struct {
	Stack []CloudCoin `json:"cloudcoin"`
}

func New(path string) *CloudCoin {

	//ans := make([]string, 0)
	//pans := make([]string, 0)

	var ccStack *CloudCoinStack
	var err *error.Error
	if (strings.HasSuffix(path, ".png")) {
		ccStack, err = ReadFromPNGFile(path)
	} else {
		ccStack, err = ReadFromFile(path)
	}

	if err != nil {
		return nil
	}

	if len(ccStack.Stack) != 1 {
		logger.Error("Stack File must contain only one Coin")
		return nil
	}

	cc := ccStack.Stack[0]
	
	if !ValidateCoin(&cc) {
		logger.Error("Coin is corrupted")
		return nil
	}

	return &cc
}

func (cc *CloudCoin) GetDenomination() int {
	snInt, err := strconv.Atoi(string(cc.Sn))
	if err != nil {
		return 0
	}

	return GetDenomination(snInt)
}
