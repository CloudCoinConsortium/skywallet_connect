package raida

import (
	"logger"
	//"strconv"
//	"fmt"
	"reflect"
	"math"
	"config"
//	"config"
//	"time"
)


type RAIDA struct {
	TotalNumber int
	SideSize int
	DetectionAgents []DetectionAgent
}



func New() *RAIDA {
	DetectionAgents := make([]DetectionAgent, config.TOTAL_RAIDA_NUMBER)
	for idx, _ := range DetectionAgents {
		DetectionAgents[idx] = *NewDetectionAgent(idx)
	}

	sideSize := int(math.Sqrt(config.TOTAL_RAIDA_NUMBER))
	if sideSize * sideSize != config.TOTAL_RAIDA_NUMBER {
		panic("Invalid RAIDA Configuration")
	}

	return &RAIDA {
		TotalNumber: config.TOTAL_RAIDA_NUMBER,
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
	}

	return results
}


func (r *RAIDA) SendDefinedRequest(url string, params []map[string]string, i interface{}) ([]Result) {
	logger.Info("Doing request " + url)

	done := make(chan Result)
	for idx, agent := range r.DetectionAgents {
			go func(agent DetectionAgent, idx int) {
				agent.SendRequest(url, params[idx], done, reflect.TypeOf(i))
			}(agent, idx)
	}

	results := make([]Result, r.TotalNumber)
	chanResults := make([]Result, r.TotalNumber)
	for i := 0; i < r.TotalNumber; i++ {
		chanResults[i] = <-done
	}

	logger.Info("Done request " + url)
	for _, result := range chanResults {
		results[result.Index] = result
	}

	return results
}

func (r *RAIDA) TotalServers() int {
	return r.TotalNumber
}

func (r *RAIDA) GetSideSize() int {
	return r.SideSize
}

