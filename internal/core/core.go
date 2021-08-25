package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
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
			panic("Failed to find current didirPath}

		return path
	*/
}


func GetCoinsForImport() (*[]cloudcoin.CloudCoin, *error.Error) {
  err := UnpackCoins()
  if err != nil {
    return nil, err
  }

  ccs, err2 := GetCoinsFromSuspect()
  if err2 != nil {
    return nil, err2
  }

  return ccs, nil
}



func GetCoinsFromSuspect() (*[]cloudcoin.CloudCoin, *error.Error) { 
  return GetCoinsFromFolder(config.DIR_SUSPECT)
}

func GetCoinsFromFracked() (*[]cloudcoin.CloudCoin, *error.Error) { 
  return GetCoinsFromFolder(config.DIR_FRACKED)
}

func GetCoinsFromFolder(folderName string) (*[]cloudcoin.CloudCoin, *error.Error) { 
  folder := GetRootPath() + Ps() + folderName
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return nil, &error.Error{config.ERROR_READ_DIRECTORY, folderName + " Folder does not exist"}
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, &error.Error{config.ERROR_READ_DIRECTORY, "Failed to read folder " + folder}
	}

  ccs := make([]cloudcoin.CloudCoin, 0)
	for _, f := range files {
		fname := f.Name()
    logger.Debug("Reading " + fname)
    coinPath := folder + Ps() + fname
    cc, err2 := cloudcoin.New(coinPath)
    if err2 != nil {
      continue
    }

    ccs = append(ccs, *cc)
  }

  return &ccs, nil
}

func UnpackCoins() *error.Error {
  folder := GetRootPath() + Ps() + config.DIR_IMPORT
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		return &error.Error{config.ERROR_READ_DIRECTORY, "Import Folder does not exist"}
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return &error.Error{config.ERROR_READ_DIRECTORY, "Failed to read folder " + folder}
	}

	for _, f := range files {
		fname := f.Name()
		if !strings.HasSuffix(fname, ".stack") && !strings.HasSuffix(fname, ".png") {
			logger.Debug("Skipping file " + fname)
			continue
		}

    logger.Debug("Reading " + fname)
    stackPath := folder + Ps() + fname
    ccStack, err2 := cloudcoin.NewStack(stackPath)
    if err2 != nil {
      return err2
    }

    for _, cc := range(ccStack.Stack) {
      err3 := SaveCoin(cc, config.DIR_SUSPECT)
      if err3 != nil {
        if (err3.Code == config.ERROR_COIN_EXISTS) {
          continue
        }
        return err3
      }

      cc.Path = GetRootPath() + Ps() + config.DIR_SUSPECT + Ps() + cc.GetName()
    }

    MoveFile(stackPath, GetImportedDir())
  }

  return nil
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
	MkDir(rootDir + Ps() + config.DIR_LIMBO)
	MkDir(rootDir + Ps() + config.DIR_IMPORT)
	MkDir(rootDir + Ps() + config.DIR_IMPORTED)
	MkDir(rootDir + Ps() + config.DIR_ID)
	MkDir(rootDir + Ps() + config.DIR_SUSPECT)
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



func MoveFile(oldFile string, dirName string) {
  newFile := dirName + Ps() + filepath.Base(oldFile)
  logger.Debug("Moving " + oldFile + " to " + newFile)
	err := os.Rename(oldFile, newFile)
	if err != nil {
		logger.Error("Failed to rename: " + err.Error())
	}
}

func MoveCoinToBank(cc cloudcoin.CloudCoin) {
  MoveCoin(cc, config.DIR_BANK)
}

func MoveCoinToFracked(cc cloudcoin.CloudCoin) {
  MoveCoin(cc, config.DIR_FRACKED)
}

func MoveCoin(cc cloudcoin.CloudCoin, dirName string) {
	dirPath := GetRootPath() + Ps() + dirName
	newPath := dirPath + Ps() + cc.GetName()
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

func GetImportedDir() string {
	return GetRootPath() + Ps() + config.DIR_IMPORTED
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
	var str = fmt.Sprintf("{\"Status\":\"fail\", \"Code\":%d, \"Message\":\"%s\", \"Time\":\"%s\"}", code, txt, time.Since(time.Now()))

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

func MoveCoinNewContent(cc cloudcoin.CloudCoin, dirName string) *error.Error {
  oldPath := cc.Path

	cc.Path = GetRootPath() + Ps() + dirName + Ps() + cc.GetName()
  
  logger.Debug("old " + oldPath + " new " + cc.Path)

  err := SaveCoin(cc, dirName)
  if err != nil {
    return err
  }

  err2 := os.Remove(oldPath)
  if err2 != nil {
    logger.Error("Failed to delete " + oldPath + ": " + err2.Error())
  }

  return nil
}

func UpdateCoin(cc cloudcoin.CloudCoin) *error.Error {
  dirPath := filepath.Base(filepath.Dir(cc.Path))

  tmpPath := cc.Path + ".tmp"

  err := os.Rename(cc.Path, tmpPath)
  if err != nil {
    logger.Error("Failed to rename " + cc.Path + " to " + tmpPath + ": " + string(err.Error()))
    return &error.Error{config.ERROR_OPEN_FILE, "Failed to rename coin"}
  }

  err2 := SaveCoin(cc, dirPath)
  if err2 != nil {
    logger.Error("Failed to save coin: " + string(err2.Error()))
    return err2
  }

  os.Remove(tmpPath)

  return nil
}

func SaveCoin(cc cloudcoin.CloudCoin, dirPath string) *error.Error {
	data := []byte(cc.GetContent())
	fileName := GetRootPath() + Ps() + dirPath + Ps() + cc.GetName()

	logger.Debug("Saving " + fileName)
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		logger.Error("Failed to save coin " + string(cc.Sn) + ". already exists in " + dirPath)
		return &error.Error{config.ERROR_COIN_EXISTS, "Coin " + string(cc.Sn) + " exists in " + dirPath}
	}
	err2 := ioutil.WriteFile(fileName, data, 0644)
	if err2 != nil {
		logger.Error("Failed to write to file: " + string(err2.Error()))
		return &error.Error{config.ERROR_WRITE_FILE, "Failed to write to file: " + string(err2.Error())}
	}

	return nil
}
