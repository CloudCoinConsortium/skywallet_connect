package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"net"
	"strconv"
	"time"
	"config"
	"logger"
)

const ERR_TIMEOUT = 1
const ERR_NET = 2
const ERR_URL = 3
const ERR_COMMON = 4

type HClient struct {
	baseUrl string
	timeout int
	index int
}

type Error struct {
	Code int
	Message string
}

func New(url string, index int) *HClient {
	return &HClient{
		baseUrl: url,
		timeout: config.DEFAULT_TIMEOUT,
		index: index,
	}
}

func(c *HClient) log(message string) {
	prefix := "RAIDA" + strconv.Itoa(c.index)
	logger.Debug(prefix + " " + message)
}

func (c *HClient) Send(nurl string) (string, *Error) {
	sendURL := fmt.Sprintf("%s%s", c.baseUrl, nurl)
	c.log("GET " + sendURL)

	//create get request
	URLData := url.Values{}
	//URLData.Set("account", strconv.Itoa(account))
	//URLData.Set("tag", tag)

	u, _ := url.Parse(sendURL)
	u.RawQuery = URLData.Encode()
	Request := fmt.Sprintf("%v", u)
	//create timeout response
	body := "RAIDA 14 timed out." //set it as a fail to begin with

	var raidahttp = &http.Client{
		Timeout: time.Duration(c.timeout) * time.Second,
	}

	//send request
	beforeSeconds := time.Now()
	response, err := raidahttp.Get(Request)
	elapsedSeconds := time.Since(beforeSeconds).Nanoseconds() / 1000000

	c.log("Total time: " + strconv.Itoa(int(elapsedSeconds)) + " ms")

	if (err != nil) {
		rerr := &Error{}
		switch err := err.(type) {
		case net.Error:
			if (err.Timeout()) {
				rerr.Code = ERR_TIMEOUT
				rerr.Message = "Network Timeout"
				fmt.Println("Timeout")
			} else {
				rerr.Code = ERR_NET
				rerr.Message = "Network Error, " + err.Error()
			}
		case *url.Error:
			fmt.Println("url error")
			if err, ok := err.Err.(net.Error); ok && err.Timeout() {
				rerr.Code = ERR_TIMEOUT
				rerr.Message = "Network URL Timeout"
			} else {
				rerr.Code = ERR_URL
				rerr.Message = "URL Error, " + err.Error()
			}
		default:
				rerr.Code = ERR_COMMON
				rerr.Message = err.Error()
		}

		return "", rerr
	}

	bodybytes, _ := ioutil.ReadAll(response.Body)
	body = string(bodybytes)

	c.log(body)

	//fmt.Printf("err=%v %v\n",err, err.Error())
	return  body, nil


}
