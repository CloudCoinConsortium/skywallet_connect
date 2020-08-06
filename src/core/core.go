package core


import (
	"os/user"
	"logger"
	"config"
	"os"
	"cloudcoin"
	"io/ioutil"
	"error"
	"fmt"
)


func GetRootPath() string {
	root, err := user.Current()
	if err != nil {
		logger.Error("Failed to get current user")
		panic("Failed to get current user")
	}

	return root.HomeDir + Ps() + config.TOPDIR
}


func MkDirs() {
	rootDir := GetRootPath()

	MkDir(rootDir)
	MkDir(rootDir + Ps() + config.DIR_ID)


	/*
	if cc, err := GetIDCoin(); err != nil {
		logger.Debug("xxxxxxxxx: " +err.Message)
	} else {

	fmt.Printf("s=%s\n",cc.Sn)
	}
	*/

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


func GetIDCoin() (*cloudcoin.CloudCoin, *error.Error) {
	idpath := GetRootPath() + Ps() + config.DIR_ID

	_, err := os.Stat(idpath)
	if os.IsNotExist(err) {
		return nil, &error.Error{"Path " + idpath + " doesn't exist"}
	}

	files, err := ioutil.ReadDir(idpath)
	if err != nil {
		return nil, &error.Error{"Failed to readdir " + idpath}
	}

	var ccname string
	for _, f := range files {
		ccname =  idpath + Ps() + f.Name()
		break
	}

	if ccname == "" {
		return nil, &error.Error{"ID Coin not found"}
	}

	logger.Debug("Foind ID coin: " + ccname)

	cc := cloudcoin.New(ccname)
	if cc == nil {
		return nil, &error.Error{"Failed to parse ID Coin"}
	}

	return cc, nil
}

func JsonError(txt string) string {
	var str = fmt.Sprintf("{\"status\":\"fail\", \"message\":\"%s\"}", txt)

	return str
}

