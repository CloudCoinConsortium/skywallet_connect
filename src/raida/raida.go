package raida

import (
	"logger"
	//"strconv"
//	"fmt"
	"reflect"
	"math"
//	"config"
//	"time"
)

const TOTAL_RAIDA_NUMBER = 25

type RAIDA struct {
	TotalNumber int
	SideSize int
	DetectionAgents []DetectionAgent
}



func New() *RAIDA {
	DetectionAgents := make([]DetectionAgent, TOTAL_RAIDA_NUMBER)
	for idx, _ := range DetectionAgents {
		DetectionAgents[idx] = *NewDetectionAgent(idx)
	}

	sideSize := int(math.Sqrt(TOTAL_RAIDA_NUMBER))
	if sideSize * sideSize != TOTAL_RAIDA_NUMBER {
		panic("Invalid RAIDA Configuration")
	}

	return &RAIDA {
		TotalNumber: TOTAL_RAIDA_NUMBER,
		DetectionAgents: DetectionAgents,
		SideSize: sideSize,
	}
}

/*
type Common struct {
	Server	string `json:"server"`
	Version string `json:"version"`
	Time	string `json:"time"`
}
*/
func (r *RAIDA) SendRequest(url string, params map[string]string, i interface{}) ([]Result) {
	logger.Info("Doing request " + url)

	done := make(chan Result)
	for _, agent := range r.DetectionAgents {
			go func(agent DetectionAgent) {
				agent.SendRequest(url, params, done, reflect.TypeOf(i))
			}(agent)
	}


	results := make([]Result, r.TotalNumber)
	chanResults := make([]Result, r.TotalNumber)
	for i := 0; i < r.TotalNumber; i++ {
		chanResults[i] = <-done
	}

	logger.Info("Done request " + url)
	for _, result := range chanResults {
		results[result.Index] = result
	//	fmt.Println("raida"+strconv.Itoa(idx) + "idx="+strconv.Itoa(result.Index) + " r="+result.Message)
	}
	/*
	for idx, result := range results {
		if result.ErrCode == config.REMOTE_RESULT_ERROR_NONE {
			r := result.Data.(*Common)
			fmt.Println("raida"+strconv.Itoa(idx) + " r="+result.Message + " cr="+strconv.Itoa(result.ErrCode) + " srv="+r.Server)
		}
	}
	*/

	return results
}

func (r *RAIDA) TotalServers() int {
	return r.TotalNumber
}

func (r *RAIDA) GetSideSize() int {
	return r.SideSize
}
