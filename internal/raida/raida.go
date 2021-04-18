package raida

import (
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"

	//	"fmt"
	"math"
	"reflect"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	//	"config"
	//	"time"
)

type RAIDA struct {
	TotalNumber     int
	SideSize        int
	DetectionAgents []DetectionAgent
}

func New() *RAIDA {
	DetectionAgents := make([]DetectionAgent, config.TOTAL_RAIDA_NUMBER)
	for idx, _ := range DetectionAgents {
		DetectionAgents[idx] = *NewDetectionAgent(idx)
	}

	sideSize := int(math.Sqrt(config.TOTAL_RAIDA_NUMBER))
	if sideSize*sideSize != config.TOTAL_RAIDA_NUMBER {
		panic("Invalid RAIDA Configuration")
	}

	return &RAIDA{
		TotalNumber:     config.TOTAL_RAIDA_NUMBER,
		DetectionAgents: DetectionAgents,
		SideSize:        sideSize,
	}
}

/*
type Common struct {
	Server	string `json:"server"`
	Version string `json:"version"`
	Time	string `json:"time"`
}
*/

func (r *RAIDA) SendRequest(url string, params map[string]string, i interface{}) []Result {
	logger.Info("Doing request " + url)

	done := make(chan Result)
	for _, agent := range r.DetectionAgents {
		go func(agent DetectionAgent) {
			agent.SendRequest(url, params, done, nil, false, reflect.TypeOf(i))
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

func (r *RAIDA) SendDefinedRequestRaw(url string, params []map[string]string) []Result {
	return r.sendDefinedRequest(url, params, nil, true, false, true)
}

func (r *RAIDA) SendDefinedRequestNoWait(url string, params []map[string]string, i interface{}) []Result {
	return r.sendDefinedRequest(url, params, i, false, false, false)
}

func (r *RAIDA) SendDefinedRequest(url string, params []map[string]string, i interface{}) []Result {
	return r.sendDefinedRequest(url, params, i, true, false, false)
}

func (r *RAIDA) SendDefinedRequestPost(url string, params []map[string]string, i interface{}) []Result {
	return r.sendDefinedRequest(url, params, i, true, true, false)
}

func (r *RAIDA) sendDefinedRequest(url string, params []map[string]string, i interface{}, wait bool, post bool, raw bool) []Result {
	logger.Info("Doing request " + url)

	done := make(chan Result)
	var doneIssued chan bool
	if !wait {
		doneIssued = make(chan bool)
	} else {
		doneIssued = nil
	}
	for idx, agent := range r.DetectionAgents {
		if params[idx] == nil {
			logger.Debug("Skipping Raida " + strconv.Itoa(idx))
			go func(idx int) {
				r := &Result{Index: idx, ErrCode: config.REMOTE_RESULT_ERROR_SKIPPED}
				if doneIssued != nil {
					doneIssued <- true
				}
				done <- *r
			}(idx)
			continue
		}

		if raw {
			go func(agent DetectionAgent, idx int) {
				agent.SendRequestRaw(url, params[idx], done, doneIssued, post)
			}(agent, idx)
		} else {
			go func(agent DetectionAgent, idx int) {
				agent.SendRequest(url, params[idx], done, doneIssued, post, reflect.TypeOf(i))
			}(agent, idx)
		}
	}

	if !wait {
		logger.Debug("Don't need to wait for completion. We will only wait till the requests are sent")
		for i := 0; i < r.TotalNumber; i++ {
			<-doneIssued
		}

		return nil
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
