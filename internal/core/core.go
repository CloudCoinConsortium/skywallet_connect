package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
)

func GetRootPath() string {
	root, err := user.Current()
	if err != nil {
		logger.Error("Failed to get current user")
		panic("Failed to get current user")
	}

	return root.HomeDir + Ps() + config.TOPDIR
	/*
		path, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to find current directory")
			panic("Failed to find current directory")
		}

		return path
	*/
}

func GetSNSFromFolder(folder string) (map[int]string, *error.Error) {
	logger.Debug("Reading dir " + folder)

	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return nil, &error.Error{config.ERROR_READ_DIRECTORY, "Folder does not exist"}
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, &error.Error{config.ERROR_READ_DIRECTORY, "Failed to read folder " + folder}
	}

	sns := make(map[int]string)
	for _, f := range files {
		fname := f.Name()
		if !strings.HasSuffix(fname, ".stack") {
			logger.Debug("Skipping file " + fname)
			continue
		}

		s := strings.Split(fname, ".")
		if len(s) < 4 {
			logger.Debug("Skipping file " + fname + " Can't parse its name")
			continue
		}

		sn := s[3]
		isn, err := strconv.Atoi(sn)
		if err != nil {
			logger.Debug("Skipping file " + fname + " Can't parse its SN")
			continue
		}

		sns[isn] = folder + Ps() + fname
	}

	return sns, nil

}

func CreateFolders() {
	rootDir := GetRootPath()

	MkDir(rootDir)
	MkDir(rootDir + Ps() + config.DIR_BANK)
	MkDir(rootDir + Ps() + config.DIR_FRACKED)
	MkDir(rootDir + Ps() + config.DIR_SENT)
	MkDir(rootDir + Ps() + config.DIR_COUNTERFEIT)
	MkDir(rootDir + Ps() + config.DIR_ID)
}

func MoveCoinToCounterfeit(cc cloudcoin.CloudCoin) {
	newPath := GetCounterfeitDir() + Ps() + cc.GetName()
	logger.Debug("Moving " + string(cc.Sn) + " to Counterfeit: " + cc.Path + " to " + newPath)

	err := os.Rename(cc.Path, newPath)
	if err != nil {
		logger.Error("Failed to rename: " + err.Error())
	}
}
func MoveCoinToSent(cc cloudcoin.CloudCoin) {
	newPath := GetSentDir() + Ps() + cc.GetName()
	logger.Debug("Moving " + string(cc.Sn) + " to Sent: " + cc.Path + " to " + newPath)

	err := os.Rename(cc.Path, newPath)
	if err != nil {
		logger.Error("Failed to rename: " + err.Error())
	}
}

func GetBankDir() string {
	return GetRootPath() + Ps() + config.DIR_BANK
}

func GetFrackedDir() string {
	return GetRootPath() + Ps() + config.DIR_FRACKED
}
func GetCounterfeitDir() string {
	return GetRootPath() + Ps() + config.DIR_COUNTERFEIT
}
func GetSentDir() string {
	return GetRootPath() + Ps() + config.DIR_SENT
}

func GetLogPath() string {
	return GetRootPath() + Ps() + config.LOG_FILENAME
}

func GetConfigPath() string {
	return GetRootPath() + Ps() + config.CONFIG_FILENAME
}

func Ps() string {
	return string(os.PathSeparator)
}

func MkDir(path string) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return
	}

	err = os.Mkdir(path, 0700)
	if err != nil {
		logger.Error("Failed to create " + path)
		panic("Failed to create folder " + path)
	}

	logger.Debug("Created folder " + path)
}

func GetIDCoinFromPath(idpath string) (*cloudcoin.CloudCoin, *error.Error) {
	_, err := os.Stat(idpath)
	if os.IsNotExist(err) {
		return nil, &error.Error{config.ERROR_ID_COIN_MISSING, "Failed to find ID coin"}
	}
	cc, err2 := cloudcoin.New(idpath)
	if err2 != nil {
		msg := err2.Message
		if err2.Code == config.ERROR_MORE_THAN_ONE_CC {
			msg = "The ID Coin file specified has more than one coin. Your ID coin file can have only one coin"
		}
		return nil, &error.Error{err2.Code, msg}
	}

	return cc, nil
}

func InitLog() {
	logFilePath := GetLogPath()
	stat, _ := os.Stat(logFilePath)
	if stat != nil {
		if stat.Size() > config.MAX_LOG_SIZE {
			RotateLog(logFilePath)
		}
	}

	file, err0 := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err0 != nil {
		ShowError(config.ERROR_INCORRECT_USAGE, "Failed to open logfile")
	}

	config.LogDesc = file
}

func ReadConfig() {
	configFilePath := GetConfigPath()
	var content []byte

	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		logger.Debug("Failed to read config: " + configFilePath)
		return
	}

	err2 := config.Apply(string(content))
	if err2 != nil {
		logger.Debug("Failed to parse config: " + err2.Message)
		return
	}
}

func GetIDCoin() (*cloudcoin.CloudCoin, *error.Error) {
	idpath := GetRootPath() + Ps() + config.DIR_ID

	_, err := os.Stat(idpath)
	if os.IsNotExist(err) {
		return nil, &error.Error{config.ERROR_ID_COIN_MISSING, "Failed to find ID coin"}
	}

	files, err := ioutil.ReadDir(idpath)
	if err != nil {
		return nil, &error.Error{config.ERROR_READ_DIRECTORY, "Failed to read folder " + idpath}
	}

	var ccname string
	for _, f := range files {
		ccname = idpath + Ps() + f.Name()
		break
	}

	if ccname == "" {
		return nil, &error.Error{config.ERROR_ID_COIN_MISSING, "Failed to find ID coin"}
	}

	logger.Debug("Found ID coin: " + ccname)
	cc, err2 := cloudcoin.New(ccname)
	if err2 != nil {
		return nil, &error.Error{config.ERROR_INVALID_CLOUDCOIN_FORMAT, "Failed to parse ID Coin"}
	}

	return cc, nil
}

func JsonError(code int, txt string) string {
	var str = fmt.Sprintf("{\"status\":\"fail\", \"code\":%d \"message\":\"%s\", \"time\":\"%s\"}", code, txt, time.Since(time.Now()))

	return str
}

func ShowError(code int, txt string) {
	fmt.Printf("%s", JsonError(code, txt))
	logger.Error(JsonError(code, txt))
	os.Exit(code)
}

func RotateLog(path string) {
	t := time.Now()
	datetime := t.Format("01-02-2006T15:04:05")
	newFile := path + "." + datetime

	os.Rename(path, newFile)
}

func SaveToBank(cc cloudcoin.CloudCoin) *error.Error {
	data := []byte(cc.GetContent())
	fileName := GetBankDir() + Ps() + cc.GetName()

	logger.Debug("Saving " + fileName)
	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		logger.Error("Failed to write to file: " + string(err.Error()))
		return &error.Error{config.ERROR_WRITE_FILE, "Failed to write to file: " + string(err.Error())}
	}

	return nil
}
