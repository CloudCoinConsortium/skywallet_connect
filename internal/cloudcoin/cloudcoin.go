package cloudcoin

import (
	"encoding/json"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"os"
	//	"io/ioutil"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"
)

type CloudCoin struct {
	Nn       json.Number `json:"nn"`
	Sn       json.Number `json:"sn"`
	Ans      []string    `json:"an"`
	Pans     []string    `json:"-"`
	Path     string      `json:"-"`
	Type     int         `json:"-"`
	Statuses []int       `json:"-"`
  PownString string    `json:"pownstring,omitempty"`
}

type CloudCoinStack struct {
	Stack []CloudCoin `json:"cloudcoin"`
}

func NewFromData(Nn, Sn int) *CloudCoin {
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

func NewStack(path string) (*CloudCoinStack, *error.Error) {
	var ccStack *CloudCoinStack
	var err *error.Error
	if strings.HasSuffix(path, ".png") {
		ccStack, err = ReadFromPNGFile(path)
	} else {
		ccStack, err = ReadFromFile(path)
	}

	if err != nil {
		return nil, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "Failed to parse coin"}
	}

	if len(ccStack.Stack) == 0 {
		return nil, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "No Coins in the stack file"}
	}

  for _, cc := range(ccStack.Stack) {
    logger.Debug("cc=\n" + string(cc.Sn))
  	if !ValidateCoin(&cc) {
	  	logger.Error("Coin is corrupted "  + string(cc.Sn))
		  return nil, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "Failed to parse coin " + string(cc.Sn)}
  	}
  }

  return ccStack, nil
}

func New(path string) (*CloudCoin, *error.Error) {
	var ctype int
	var ccStack *CloudCoinStack
	var err *error.Error
	if strings.HasSuffix(path, ".png") {
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
		return nil, &error.Error{config.ERROR_MORE_THAN_ONE_CC, "Stack file must have only one coin"}
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

  if (cc.Pans == nil) {
  	cc.Pans = make([]string, config.TOTAL_RAIDA_NUMBER)
  }

  cc.SetStatusesFromPownString()

	return &cc, nil
}

func (cc *CloudCoin) GetDenomination() int {
	snInt, err := strconv.Atoi(string(cc.Sn))
	if err != nil {
		return 0
	}

	return GetDenomination(snInt)
}

func (cc *CloudCoin) SetStatusesFromPownString() {
  if (cc.PownString == "") {
    return
  }

  if (len(cc.PownString) != config.TOTAL_RAIDA_NUMBER) {
    return
  }

  for idx, c := range(cc.PownString) {
    status := config.RAIDA_STATUS_UNTRIED
    switch c {
    case 'p':
      status = config.RAIDA_STATUS_PASS
    case 'e':
      status = config.RAIDA_STATUS_ERROR
    case 'u':
      status = config.RAIDA_STATUS_UNTRIED
    case 'n':
      status = config.RAIDA_STATUS_NORESPONSE
    case 'f':
      status = config.RAIDA_STATUS_FAIL
    }

    cc.SetDetectStatus(idx, status)
  }
}

func (cc *CloudCoin) GetFileName() string {
	if cc.Path == "" {
		return "OwnerWallet"
	}

	base := filepath.Base(cc.Path)
	if cc.Type == config.TYPE_PNG {
		re := regexp.MustCompile(`(.+)\.png$`)
		base = re.ReplaceAllString(base, "$1")
	}

  if cc.Type == config.TYPE_STACK {
		re := regexp.MustCompile(`(.+)\.stack$`)
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

func (cc *CloudCoin) IsAuthentic() (bool, bool, bool) {
  passed := 0
  failed := 0
  for _, status := range(cc.Statuses) {
    if status == config.RAIDA_STATUS_PASS {
      passed++
    } else if status == config.RAIDA_STATUS_FAIL {
      failed++
    }
  }

  isAuthentic := passed >= config.MIN_PASSED_NUM_TO_BE_AUTHENTIC
  hasFailed := failed > 0
  isCounterfeit := failed >= config.MAX_FAILED_NUM_TO_BE_COUNTERFEIT

  return isAuthentic, hasFailed, isCounterfeit
}



func (cc *CloudCoin) GetPownString() string {
	pownString := ""
	for idx, _ := range cc.Statuses {
		switch cc.Statuses[idx] {
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

func (cc *CloudCoin) SetPownString() {
  cc.PownString = cc.GetPownString()
}

func (cc *CloudCoin) SetAn(idx int, an string) {
	cc.Ans[idx] = an
}

func (cc *CloudCoin) GenerateMyPans() {
	for idx := 0; idx < config.TOTAL_RAIDA_NUMBER; idx++ {
		cc.Pans[idx], _ = GeneratePan()
	}
}

func (cc *CloudCoin) SetAnsToPansIfPassed() {
	for idx := 0; idx < config.TOTAL_RAIDA_NUMBER; idx++ {
    if (cc.Statuses[idx] != config.RAIDA_STATUS_PASS) {
      continue
    }

		cc.Ans[idx] = cc.Pans[idx]
	}
}

func (cc *CloudCoin) GetContent() string {
	var cstack CloudCoinStack

  cc.SetPownString()
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
