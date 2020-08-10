package raida


import(
	"httpclient"
	"config"
	"fmt"
	"logger"
	"strconv"
	"encoding/json"
	"reflect"
)

type DetectionAgent struct {
	index int
	c *httpclient.HClient
}

type Result struct {
	ErrCode int
	Message string
	Index int
	Data interface{}
}

func(da *DetectionAgent) log(message string) {
	prefix := "RAIDA" + strconv.Itoa(da.index)
	logger.Debug(prefix + " " + message)
}

func(da *DetectionAgent) logError(message string) {
	prefix := "RAIDA" + strconv.Itoa(da.index)
	logger.Error(prefix + " " + message)
}

func NewDetectionAgent(index int) *DetectionAgent {
	url := fmt.Sprintf("https://raida%d.%s",  index, config.DEFAULT_DOMAIN)
	return &DetectionAgent{
		index: index,
		c: httpclient.New(url, index),
	}
}

func (da *DetectionAgent) SendRequest(url string, params map[string]string, done chan Result, doneIssued chan bool, post bool, t reflect.Type) {
	result := &Result{}

	if response, err := da.c.Send(url, params, doneIssued, post); err != nil {
		da.logError("Failed to send request: " + err.Message)
		result.Message = err.Message
		if (err.Code == httpclient.ERR_TIMEOUT) {
			result.ErrCode = config.REMOTE_RESULT_ERROR_TIMEOUT
		} else {
			result.ErrCode = config.REMOTE_RESULT_ERROR_COMMON
		}
	} else {
		if doneIssued != nil {
			da.log("Response ignored")
		} else {
			da.log("Response received")
			da.log(response)
		}
		result.Message = response
		result.ErrCode = config.REMOTE_RESULT_ERROR_NONE

		data := reflect.New(t).Interface()
		bytes := []byte(response)
		if err := json.Unmarshal(bytes, &data); err != nil {
			da.logError("Failed to parse JSON")
			result.ErrCode = config.REMOTE_RESULT_ERROR_COMMON
		} else {
			result.Data = data
		}
	}

	result.Index = da.index
	done <-*result
}
