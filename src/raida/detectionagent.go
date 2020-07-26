package raida


import(
	"httpclient"
	"config"
	"fmt"
	"logger"
)

type DetectionAgent struct {
	index int
	c *httpclient.HClient
}

type Error struct {
	Code int
	Message string
}

func NewDetectionAgent(index int) *DetectionAgent {
	url := fmt.Sprintf("https://raida%d.%s",  index, config.DEFAULT_DOMAIN)
	return &DetectionAgent{
		index: index,
		c: httpclient.New(url, index),
	}
}

func (da *DetectionAgent) DoRequest(url string) (string, *Error) {
	fmt.Println("hhh\n")

	if response, err := da.c.Send(url); err != nil {
		logger.Error("Failed to send request: " + err.Message)
	} else {

	fmt.Println("r="+response)
	}
	logger.Info("yyy")
	return "xxx", &Error{1, "fuck"}
}
