package main

import (
	"fmt"
	"raida"
	"logger"
)

func main() {
	fmt.Println("hello world")

	/*
	var i interface{}

	i = "xxx"
	i = 2
	i = 3.2

	r := i.(float64)
	y := i.(string)

	fmt.Println("x=", r)
	fmt.Println("x=", y)

	*/

	//client := httpclient.New("https://raida14.cloudcoin.global")

	da := raida.NewDetectionAgent(142)
	if resp, err := da.DoRequest("/service/echo"); err != nil {
		fmt.Println("error")
	} else {

		fmt.Println("r=", resp)
		}

	logger.Info("xxx")
//	client.Send("/service/echo")
	//httpclient.Send("ggg")
}
