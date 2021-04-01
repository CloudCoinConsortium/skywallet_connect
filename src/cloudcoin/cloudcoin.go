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
	"fmt"
	"path/filepath"
	"regexp"
)

type CloudCoin struct {
	Nn json.Number `json:"nn"`
	Sn json.Number `json:"sn"`
	Ans []string `json:"an"`
	Pans []string `json:"-"`
	Path string `json:"-"`
	Type int `json:"-"`
	Statuses[] int `json:"-"`
}

type CloudCoinStack struct {
	Stack []CloudCoin `json:"cloudcoin"`
}

func NewFromData(Nn, Sn int) (*CloudCoin) {
	var cc CloudCoin

	cc.Nn = json.Number(strconv.Itoa(Nn))
	cc.Sn = json.Number(strconv.Itoa(Sn))
	cc.Ans = make([]string, config.TOTAL_RAIDA_NUMBER)
	cc.Pans = make([]string, config.TOTAL_RAIDA_NUMBER)
	cc.Path = ""
	cc.Type = config.TYPE_STACK
	cc.Statuses = make([]int, config.TOTAL_RAIDA_NUMBER)
	for idx := 0; idx < config.TOTAL_RAIDA_NUMBER; idx++ {
		cc.Statuses[idx] = config.RAIDA_STATUS_UNTRIED
	}

	return &cc
}

func New(path string) (*CloudCoin, *error.Error) {

	//ans := make([]string, 0)
	//pans := make([]string, 0)

	var ctype int
	var ccStack *CloudCoinStack
	var err *error.Error
	if (strings.HasSuffix(path, ".png")) {
		ccStack, err = ReadFromPNGFile(path)
		ctype = config.TYPE_PNG
	} else {
		ccStack, err = ReadFromFile(path)
		ctype = config.TYPE_STACK
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

	cc.Path = path
	cc.Type = ctype

	cc.Statuses = make([]int, config.TOTAL_RAIDA_NUMBER)
	for idx := 0; idx < config.TOTAL_RAIDA_NUMBER; idx++ {
		cc.Statuses[idx] = config.RAIDA_STATUS_UNTRIED
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

func (cc *CloudCoin) GetFileName() string {
	if (cc.Path == "") {
		return "OwnerWallet"
	}

	base := filepath.Base(cc.Path)
	if (cc.Type == config.TYPE_PNG) {
		re := regexp.MustCompile(`(.+)\.png$`)
		base = re.ReplaceAllString(base, "$1")
	}

	return base
}

func (cc *CloudCoin) GetName() string {
	s := fmt.Sprintf("%d.CloudCoin.%s.%s.stack", cc.GetDenomination(), string(cc.Nn), string(cc.Sn))

	return s
}

func (cc *CloudCoin) SetDetectStatus(idx int, status int) {
	cc.Statuses[idx] = status
}

func (cc *CloudCoin) GetPownString() string {
	pownString := ""
	for idx, _ := range cc.Statuses {
		switch (cc.Statuses[idx]) {
			case config.RAIDA_STATUS_UNTRIED:
				pownString += "u"
			case config.RAIDA_STATUS_FAIL:
				pownString += "f"
			case config.RAIDA_STATUS_PASS:
				pownString += "p"
			case config.RAIDA_STATUS_ERROR:
				pownString += "e"
			case config.RAIDA_STATUS_NORESPONSE:
				pownString += "n"
		}
	}

	return pownString
}

func (cc *CloudCoin) SetAn(idx int, an string) {
	cc.Ans[idx] = an
}

func (cc *CloudCoin) GetContent() string {
	var cstack CloudCoinStack

	cstack.Stack = make([]CloudCoin, 1)
	cstack.Stack[0] = *cc

	data, err := json.Marshal(cstack)
	if err != nil {
		logger.Debug("Failed to Marshal CloudCoin")
		return ""
	}

	sdata := string(data)

	return sdata
}

