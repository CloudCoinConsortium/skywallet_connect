package config

import (
	"os"
)


//var DEFAULT_TIMEOUT = 1
var DEFAULT_TIMEOUT = 50
var DEFAULT_DOMAIN = "cloudcoin.global"


const LOG_LEVEL_DEBUG = 3
const LOG_LEVEL_INFO = 2
const LOG_LEVEL_ERROR = 1


const REMOTE_RESULT_ERROR_NONE = 0
const REMOTE_RESULT_ERROR_TIMEOUT = 1
const REMOTE_RESULT_ERROR_COMMON = 2
const REMOTE_RESULT_ERROR_SKIPPED = 3

const DEFAULT_NN = 1
const PUBLIC_CHANGE_MAKER_ID = 2

var MAX_FIXTRANSFER_NOTES = 400

var CmdDebug bool
var CmdHelp bool
var CmdVersion bool
var CmdCommand string
var CmdLogfile string

var LogDesc *os.File

const TOTAL_RAIDA_NUMBER = 25
const NEED_VOTERS_FOR_BALANCE = 17

const (
	RAIDA_STATUS_UNTRIED = 0
	RAIDA_STATUS_PASS = 1
	RAIDA_STATUS_ERROR = 2
	RAIDA_STATUS_FAIL = 3
	RAIDA_STATUS_NORESPONSE = 4
)

const FIX_MAX_REGEXPS = 56


const TOPDIR = "skywallet_connect_home"
const DIR_BANK = "Bank"
const DIR_FRACKED = "Fracked"
const DIR_SENT = "Sent"
const DIR_COUNTERFEIT = "Counterfeit"
const DIR_ID = "ID"
const DIR_IMPORT = "Import"
const DIR_IMPORTED = "Imported"
const DIR_SUSPECT = "Suspect"
const DIR_LIMBO = "Limbo"


const TYPE_STACK = 1
const TYPE_PNG = 4


const LOG_FILENAME = "main.log"
const CONFIG_FILENAME = "config.toml"

const MAX_LOG_SIZE = 50000000

const CHANGE_METHOD_250F = 1
const CHANGE_METHOD_100E = 2
const CHANGE_METHOD_25B = 3
const CHANGE_METHOD_5A = 4

const META_ENV_SEPARATOR = "*"

const MAX_NOTES_TO_SEND = 100
//const MAX_NOTES_TO_SEND = 3

const MIN_PASSED_NUM_TO_BE_AUTHENTIC = 14
const MAX_FAILED_NUM_TO_BE_COUNTERFEIT = 12

const ERROR_INCORRECT_USAGE = 1
const ERROR_GET_SERIAL_NUMBER_FROM_IP = 2
const ERROR_OPEN_FILE = 3
const ERROR_READ_FILE = 4
const ERROR_CORRUPTED_PNG_FILE = 5
const ERROR_CLOUDCOIN_NOT_FOUND_IN_PNG = 6
const ERROR_CLOUDCOIN_PNG_CRC32_INCORRECT = 7
const ERROR_CLOUDCOIN_PARSE = 8
const ERROR_INVALID_CLOUDCOIN_FORMAT = 9
const ERROR_GENERATE_RANDOM_NUMBER = 10
const ERROR_ID_COIN_MISSING = 11
const ERROR_READ_DIRECTORY = 12
const ERROR_CHANGE_METHOD_NOT_FOUND = 13
const ERROR_SHOW_CHANGE_FAILED = 14
const ERROR_BREAK_IN_BANK_FAILED = 15
const ERROR_INSUFFICIENT_FUNDS = 16
const ERROR_RESULTS_FROM_RAIDA_OUT_OF_SYNC = 17
const ERROR_INCORRECT_AMOUNT_SPECIFIED = 18
const ERROR_INCORRECT_SKYWALLET = 19
const ERROR_SHOW_COINS_FAILED = 20
const ERROR_PICK_COINS_AFTER_SHOW = 21
const ERROR_PICK_COINS_AFTER_CHANGE = 22
const ERROR_ENCODE_JSON = 23
const ERROR_TRANSFER_FAILED = 24
const ERROR_INVALID_RECEIPT_ID = 25
const ERROR_INVALID_SKYWALLET_OWNER = 26
const ERROR_MORE_THAN_ONE_CC = 27
const ERROR_BREAK_FAILED = 28
const ERROR_WRITE_FILE = 29
const ERROR_CONFIG_PARSE = 30
const ERROR_INVALID_GUID = 31
const ERROR_COIN_EXISTS = 32
const ERROR_INTERNAL = 33
const ERROR_TICKETS_FAILED = 34