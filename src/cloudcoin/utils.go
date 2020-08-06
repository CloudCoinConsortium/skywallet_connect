package cloudcoin

import (
	"logger"
	"strconv"
	"regexp"
	"strings"
	"net"
	"error"
	"encoding/json"
	"os"
	"io/ioutil"
	"config"
)

func GuessSNFromString(param string) (int, *error.Error) {
	sn, err := strconv.Atoi(param)
	if err == nil {
		if (ValidateSN(sn)) {
			return sn, nil
		}
	}

	nsn, err := GetSNFromIP(param)
	if err != nil {
		if (ValidateSN(nsn)) {
			return nsn, nil
		}
	}


	return ResolveSkyWallet(param)
}

func ResolveSkyWallet(skywallet string) (int, *error.Error) {
	addrs, err := net.LookupHost(skywallet)
	if err != nil {
		return 0, &error.Error{}
	}

	if len(addrs) < 1 {
		return 0, &error.Error{}
	}

	addr := addrs[0]
	logger.Debug("Extracted Address " + addr)

	sn, err2 := GetSNFromIP(addr)
	if err2 != nil {
		return 0, &error.Error{"Failed to get SN from IP"}
	}

	logger.Debug("Extracted SN " + strconv.Itoa(sn))

	return sn, nil
}

func GetSNFromIP(ipaddress string) (int, *error.Error) {
	ipRegex, _ := regexp.Compile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`)
	s := ipRegex.FindStringSubmatch(strings.Trim(ipaddress, " "))
	if (len(s) == 5) {
		o2, err := strconv.Atoi(s[2])
		if (err != nil) {
			return 0,  &error.Error{"Failed to convert IP octet2"}
		}

		o3, err := strconv.Atoi(s[3])
		if (err != nil) {
			return 0, &error.Error{"Failed to convert IP octet3"}
		}

		o4, err := strconv.Atoi(s[4])
		if (err != nil) {
			return 0, &error.Error{"Failed to convert IP octet4"}
		}

		sn := (o2 << 16) | (o3 << 8) | o4

		return sn, nil
	}

	return 0, &error.Error{}
}

func ValidateSN(sn int) bool {
	return GetDenomination(sn) != 0
}

func ValidateCoin(cc *CloudCoin) bool {
	nn, err := strconv.Atoi(cc.Nn)
	if err != nil {
		return false
	}

	sn, err := strconv.Atoi(cc.Sn)
	if err != nil {
		return false
	}

	if nn != config.DEFAULT_NN {
		return false
	}

	if !ValidateSN(sn) {
		return false
	}

	if len(cc.Ans) != config.TOTAL_RAIDA_NUMBER {
		return false
	}

	return true
}

func GetDenomination(sn int) int {
	if sn < 1 {
		return 0
	}

	if sn < 2097153 {
		return 1
	}

	if sn < 4194305 {
		return 5
	}

	if sn < 6291457 {
		return 25
	}

	if sn < 14680065 {
		return 100
	}

	if sn < 16777217 {
		return 250
	}

	return 0
}



func ReadFromFile(fname string) (*CloudCoinStack, *error.Error) {
	logger.Debug("Parsing CloudCoin " + fname)


	file, err := os.Open(fname); 
	if err != nil {
		logger.Error("Failed to open file: " + fname)
		return nil, &error.Error{"Failed to open file " + fname}
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file); 
	if err != nil {
		logger.Error("Failed to read file: " + fname)
		return nil, &error.Error{"Failed to read file " + fname}
	}

	var ccStack CloudCoinStack
	err = json.Unmarshal(byteValue, &ccStack)
	if err != nil {
		logger.Error("Failed to parse file: " + fname)
		return nil, &error.Error{"Failed to parse file " + fname}
	}

	if len(ccStack.Stack) == 0 {
		logger.Error("Corrupted Stack File")
		return nil, &error.Error{"Stack File is Corrupted"}
	}

	return &ccStack, nil
}
