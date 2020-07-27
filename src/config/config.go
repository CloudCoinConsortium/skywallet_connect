package config


const DEFAULT_TIMEOUT = 1
const DEFAULT_DOMAIN = "cloudcoin.global"


const LOG_LEVEL_DEBUG = 3
const LOG_LEVEL_INFO = 2
const LOG_LEVEL_ERROR = 1


const REMOTE_RESULT_ERROR_NONE = 0
const REMOTE_RESULT_ERROR_TIMEOUT = 1
const REMOTE_RESULT_ERROR_COMMON = 2


var CmdDebug bool
var CmdHelp bool
var CmdCommand string


const (
	RAIDA_STATUS_UNTRIED = 0
	RAIDA_STATUS_PASS = 1
	RAIDA_STATUS_ERROR = 2
	RAIDA_STATUS_FAIL = 3
	RAIDA_STATUS_NORESPONSE = 4
)
