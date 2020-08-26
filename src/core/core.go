package core


import (
	//"os/user"
	"logger"
	"config"
	"os"
	"cloudcoin"
	"io/ioutil"
	"error"
	"fmt"
	"time"
)


func GetRootPath() string {
	/*
	root, err := user.Current()
	if err != nil {
		logger.Error("Failed to get current user")
		panic("Failed to get current user")
	}

	return root.HomeDir + Ps() + config.TOPDIR
	*/
	path, err := os.Getwd()
	if err != nil {
		logger.Error("Failed to find current directory")
		panic("Failed to find current directory")
	}

	return path
}


func MkDirs() {
	//rootDir := GetRootPath()

//	MkDir(rootDir)
//	MkDir(rootDir + Ps() + config.DIR_ID)

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
			panic("Failed to create " + path)
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
		if (err2.Code == config.ERROR_MORE_THAN_ONE_CC) {
			msg = "The ID Coin file specified has more than one coin. Your ID coin file can have only one coin"
		}
		return nil, &error.Error{err2.Code, msg}
	}

	return cc, nil
}

func GetIDCoin() (*cloudcoin.CloudCoin, *error.Error) {
	idpath := GetRootPath() + Ps() + config.DIR_ID

	_, err := os.Stat(idpath)
	if os.IsNotExist(err) {
		return nil, &error.Error{config.ERROR_ID_COIN_MISSING, "Failed to find ID coin, please create a folder called ID in the same folder as your raida_go program. Place one ID coins in that folder"}
	}

	files, err := ioutil.ReadDir(idpath)
	if err != nil {
		return nil, &error.Error{config.ERROR_READ_DIRECTORY, "Failed to read folder " + idpath}
	}

	var ccname string
	for _, f := range files {
		ccname =  idpath + Ps() + f.Name()
		break
	}

	if ccname == "" {
		return nil, &error.Error{config.ERROR_ID_COIN_MISSING, "Failed to find ID coin, please create a folder called ID in the same folder as your raida_go program. Place one ID coins in that folder"}
	}

	logger.Debug("Foind ID coin: " + ccname)

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
	os.Exit(code)
}

func RotateLog(path string) {
	t := time.Now()
	datetime := t.Format("01-02-2006T15:04:05")
	newFile := path + "." + datetime

	os.Rename(path, newFile)
}


