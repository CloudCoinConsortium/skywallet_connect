package logger

import (
	"fmt"
	"config"
	"time"
)


func log(severity int, message string) {
	if !config.CmdDebug {
		return
	}

	var fmessage string

	currentTime := time.Now()
	time := currentTime.Format("2006/01/02 15:04:05")

	fmessage += time + " "

	switch severity {
	case config.LOG_LEVEL_DEBUG:
		fmessage += "[DEBUG]"
	case config.LOG_LEVEL_INFO:
		fmessage += "[INFO]"
	case config.LOG_LEVEL_ERROR:
		fmessage += "[ERROR]"
	default:
		fmessage += "[INFO]"
	}

	fmessage += " " + message

	fmt.Println(fmessage)
}

func Info(message string) {
	log(config.LOG_LEVEL_INFO, message)
}

func Debug(message string) {
	log(config.LOG_LEVEL_DEBUG, message)
}

func Error(message string) {
	log(config.LOG_LEVEL_ERROR, message)
}





