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
	"crypto/rand"
	"encoding/hex"
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

func GetChangeMethod(denomination int) int {
	method := 0
	switch (denomination) {
	case 250:
		method = config.CHANGE_METHOD_250F;
		break;
	case 100:
		method = config.CHANGE_METHOD_100E;
		break;
	case 25:
		method = config.CHANGE_METHOD_25B;
		break;
	case 5:
		method = config.CHANGE_METHOD_5A;
		break;
	}
  return method
}

func CoinsGetA (a []int, cnt int) []int {
	var sns []int
	var i, j int

	sns = make([]int, cnt)

	i = 0
	j = 0
	for ; i < len(a); i++ {
		if a[i] == 0 {
			continue
		}

		sns[j] = a[i]
		a[i] = 0
		j++

		if j == cnt {
			break
		}
	}

	if j != cnt {
		return nil
	}

  return sns;
}

func CoinsGet25B (sb, ss []int) []int {
	var sns, rsns []int

	rsns = make([]int, 9)
	sns = CoinsGetA(ss, 5)
	if sns == nil {
		return nil
	}

	for i := 0; i < 5; i++ {
		rsns[i] = sns[i]
	}

	sns = CoinsGetA(sb, 4)
	if sns == nil {
		return nil
	}

	for i := 0; i < 4; i++ {
		rsns[i + 5] = sns[i]
	}

	return rsns
}

func CoinsGet100E(sb, ss, sss []int) []int {
	var sns, rsns []int

	rsns = make([]int, 12)
	sns = CoinsGetA(sb, 3)
	if sns == nil {
		return nil
	}

	for i := 0; i < 3; i++ {
		rsns[i] = sns[i]
	}

	sns = CoinsGetA(ss, 4)
	if sns == nil {
		return nil
	}

	for i := 0; i < 4; i++ {
		rsns[i + 3] = sns[i]
	}

	sns = CoinsGetA(sss, 5)
	if sns == nil {
		return nil
	}

	for i := 0; i < 5; i++ {
		rsns[i + 7] = sns[i]
	}

	return rsns
}

func CoinsGet250F(sb, ss, sss, ssss []int) []int {
	var sns, rsns []int

	rsns = make([]int, 15)
	sns = CoinsGetA(sb, 1)
	if sns == nil {
		return nil
	}

	rsns[0] = sns[0]

	sns = CoinsGetA(ss, 5)
	if sns == nil {
		return nil
	}

	for i := 0; i < 5; i++ {
		rsns[i + 1] = sns[i]
	}

	sns = CoinsGetA(sss, 4)
	if sns == nil {
		return nil
	}

	for i := 0; i < 4; i++ {
		rsns[i + 6] = sns[i]
	}

	sns = CoinsGetA(ssss, 5)
	if sns == nil {
		return nil
	}

	for i := 0; i < 5; i++ {
		rsns[i + 10] = sns[i]
	}

	return rsns;
}


func GeneratePan() (string, *error.Error) {
	return GenerateHex(16)
}

func GenerateHex(length int) (string, *error.Error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", &error.Error{"Failed to generate random string"}
	}

	return hex.EncodeToString(bytes), nil
}

