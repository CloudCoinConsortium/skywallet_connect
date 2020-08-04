package cloudcoin

import (
	"logger"
	"strconv"
	"regexp"
	"strings"
	"net"
)

type Error struct {
	s string
}

func (e *Error)  Error() string {
	return e.s
}

func GuessSNFromString(param string) (int, error) {
	sn, err := strconv.Atoi(param)
	if (err == nil) {
		if (ValidateSN(sn)) {
			return sn, nil
		}
	}

	nsn, err := GetSNFromIP(param)
	if (err != nil) {
		if (ValidateSN(nsn)) {
			return nsn, nil
		}
	}
	//}

	addrs, err := net.LookupHost(param)
	if (err != nil) {
		return 0, &Error{}
	}

	if (len(addrs) < 1) {
		return 0, &Error{}
	}

	addr := addrs[0]
	logger.Debug("Extracted Address " + addr)

	nnsn, err2 := GetSNFromIP(addr)
	if (err2 != nil) {
		logger.Debug("Extracted Address " + addr)
		return 0, &Error{}
	}

	logger.Debug("Extracted SN " + strconv.Itoa(nnsn))

	return nnsn, nil
}

func ResolveSkyWallet(skywallet string) int {
	return 0
}

func GetSNFromIP(ipaddress string) (int, *Error) {
	ipRegex, _ := regexp.Compile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`)
	s := ipRegex.FindStringSubmatch(strings.Trim(ipaddress, " "))
	if (len(s) == 5) {
		o2, err := strconv.Atoi(s[2])
		if (err != nil) {
			return 0,  &Error{}
		}

		o3, err := strconv.Atoi(s[3])
		if (err != nil) {
			return 0, &Error{}
		}

		o4, err := strconv.Atoi(s[4])
		if (err != nil) {
			return 0, &Error{}
		}

		sn := (o2 << 16) | (o3 << 8) | o4

		return sn, nil
	}

	return 0, &Error{}
}

func ValidateSN(sn int) bool {
	return sn > 0 && sn < 16777217
}
