package cloudcoin


import (
	"logger"
	"encoding/json"
//	"os"
//	"io/ioutil"
	"strconv"
  "config"
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

func New(path string) (*CloudCoin, *error.Error) {

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
		return nil, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "Failed to parse coin"}
	}

	if len(ccStack.Stack) != 1 {
		logger.Error("Stack File must contain only one Coin")
		return nil, &error.Error{config.ERROR_MORE_THAN_ONE_CC, "Failed to parse coin"}
	}

	cc := ccStack.Stack[0]
	
	if !ValidateCoin(&cc) {
		logger.Error("Coin is corrupted")
		return nil, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "Failed to parse coin"}
	}

	return &cc, nil
}

func (cc *CloudCoin) GetDenomination() int {
	snInt, err := strconv.Atoi(string(cc.Sn))
	if err != nil {
		return 0
	}

	return GetDenomination(snInt)
}
