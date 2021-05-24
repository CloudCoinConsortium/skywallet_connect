package raida

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/CloudCoinConsortium/skywallet_connect/internal/cloudcoin"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/config"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/core"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/error"
	"github.com/CloudCoinConsortium/skywallet_connect/internal/logger"
)

type FrackFixer struct {
	Servant
	trustedServers [config.TOTAL_RAIDA_NUMBER][8]int
	trustedTriads [config.TOTAL_RAIDA_NUMBER][config.FIX_MAX_REGEXPS][5]int
  failedRaidas [config.TOTAL_RAIDA_NUMBER]bool
}

type FrackFixerResponse struct {
	Server  string `json:"server"`
	Version string `json:"version"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type FrackFixerOutput struct {
  AmountFracked int
  AmountFixed int
}

func NewFrackFixer() *FrackFixer {

	trustedServers := [config.TOTAL_RAIDA_NUMBER][8]int{}
  trustedTriads := [config.TOTAL_RAIDA_NUMBER][config.FIX_MAX_REGEXPS][5]int{}
  failedRaidas := [config.TOTAL_RAIDA_NUMBER]bool{}
	return &FrackFixer{
		*NewServant(),
		trustedServers,
		trustedTriads,
    failedRaidas,
	}
}


func (v *FrackFixer) Fix() (*FrackFixerOutput, *error.Error) {
	logger.Debug("Started FrackFixer")

  fo := &FrackFixerOutput{}
  fo.AmountFixed = 0
  fo.AmountFracked = 0

  err := v.initNeighbours()
  if err != nil {
    return nil, err
  }

  ccs, err2 := core.GetCoinsFromFracked()
  if err2 != nil {
    return nil, err
  }

  if len(*ccs) == 0 {
    return fo, nil
  }

	for _, cc := range(*ccs) {
    fo.AmountFracked += cc.GetDenomination()
  }

	var bufCcs []cloudcoin.CloudCoin

  // Round 1
  logger.Debug("Round 1. Total Coins: " + strconv.Itoa(len(*ccs)))
  for i := 0; i < config.TOTAL_RAIDA_NUMBER; i++ {
    bufCcs = nil

    for _, cc := range(*ccs) {
      if (cc.Statuses[i] != config.RAIDA_STATUS_FAIL) {
        continue
      }

      logger.Debug("CC " + string(cc.Sn) + " will be fixed on r" + strconv.Itoa(i))
      bufCcs = append(bufCcs, cc)
      if len(bufCcs) == config.MAX_NOTES_TO_SEND {
        err := v.ProcessFix(i, &bufCcs)
        if err != nil {
          logger.Debug("Error " + err.Message)
        }
        bufCcs = nil
      }
    }

    if len(bufCcs) != 0 {
      err := v.ProcessFix(i, &bufCcs)
      if err != nil {
        logger.Debug("Error " + err.Message)
      }
    }
  }
/*
  // Round 2
  logger.Debug("Round 2. Total Coins: " + strconv.Itoa(len(*ccs)))
  for i := config.TOTAL_RAIDA_NUMBER - 1; i >= 0; i-- {
    bufCcs = nil

    for _, cc := range(*ccs) {
      if (cc.Statuses[i] != config.RAIDA_STATUS_FAIL) {
        continue
      }

      logger.Debug("CC " + string(cc.Sn) + " will be fixed on r" + strconv.Itoa(i))
      bufCcs = append(bufCcs, cc)
      if len(bufCcs) == config.MAX_NOTES_TO_SEND {
        err := v.ProcessFix(i, &bufCcs)
        if err != nil {
          logger.Debug("Error " + err.Message)
        }
        bufCcs = nil
      }
    }

    if len(bufCcs) != 0 {
      err := v.ProcessFix(i, &bufCcs)
      if err != nil {
        logger.Debug("Error " + err.Message)
      }
    }
  }
*/

  for _, cc := range(*ccs) {
    cc.SetPownString()
    _, hasFailed, _ := cc.IsAuthentic()
    if (hasFailed) {
      logger.Debug("Coin " + string(cc.Sn) + " wasn't fixed: " + cc.GetPownString())
      continue
    }

//    core.MoveCoinToBank(cc)
    fo.AmountFixed += cc.GetDenomination()
  }

  return fo, nil
}



func (v *FrackFixer) ProcessFix(rIdx int, ccs *[]cloudcoin.CloudCoin) *error.Error {
  logger.Debug("Fixing " + strconv.Itoa(len(*ccs)) + " coins on raida " + strconv.Itoa(rIdx))
  // Initialize Failed RAIDA array
  for i := 0; i < config.TOTAL_RAIDA_NUMBER; i++ {
    v.failedRaidas[i] = false
  }

  // Loop over corners
  for corner := 0; corner < config.FIX_MAX_REGEXPS; corner++ {
    logger.Debug("Fixing in corner " + strconv.Itoa(corner))
    fixed, _ := v.ProcessFixInCorner(rIdx, ccs, corner)
    if fixed {
      logger.Debug("Fixed Successfully")
      for _, cc := range(*ccs) {
        core.UpdateCoin(cc)
      }
      break
    }
  }

  return nil
}

func (v *FrackFixer) GetRegexChar(idx int, fivetouches [5]int) byte {
  for j := 0; j < len(fivetouches); j++ {
    if (idx == fivetouches[j]) {
      return 'p'
    }
  }

  return '.'
}

func (v *FrackFixer) GetRegexString(fivetouches [5]int) string {
  rxstr := [13]byte{}

  rxstr[0] = v.GetRegexChar(0, fivetouches)
  rxstr[1] = v.GetRegexChar(1, fivetouches)
  rxstr[2] = v.GetRegexChar(2, fivetouches)
  rxstr[3] = '.'
  rxstr[4] = '.'
  rxstr[5] = v.GetRegexChar(3, fivetouches)
  rxstr[6] = '0'
  rxstr[7] = v.GetRegexChar(4, fivetouches)
  rxstr[8] = '.'
  rxstr[9] = '.'
  rxstr[10] = v.GetRegexChar(5, fivetouches)
  rxstr[11] = v.GetRegexChar(6, fivetouches)
  rxstr[12] = v.GetRegexChar(7, fivetouches)

  return string(rxstr[:])
}

func (v *FrackFixer) ReceiveTickets(rIdx int, ccs *[]cloudcoin.CloudCoin, corner int) ([]string, *error.Error) {
  logger.Debug("Getting tickets for raida" + strconv.Itoa(rIdx))

  fivetouches := v.trustedTriads[rIdx][corner]
  slice := fivetouches[:]

  detect := NewDetect()
  tickets := detect.ProcessDetect(ccs, &slice)

  rtickets := []string{}
  for _, ridx := range(slice) {
    rtickets = append(rtickets, tickets[ridx])
  }

  return rtickets, nil
 // return nil, &error.Error{config.ERROR_TICKETS_FAILED, "Failed to Get Tickets"} 
}

func (v *FrackFixer) ProcessFixInCorner(rIdx int, ccs *[]cloudcoin.CloudCoin, corner int) (bool, *error.Error) {
  if v.failedRaidas[rIdx] {
    logger.Debug("raida" + strconv.Itoa(rIdx) + " is failed. Skipping")
    return false, nil
  }

  if corner >= len(v.trustedTriads[rIdx]) {
    logger.Debug("raida" + strconv.Itoa(rIdx) + " misconfig. Skipping")
    return false, nil
  }

  fivetouches := v.trustedTriads[rIdx][corner]
  regexString := v.GetRegexString(fivetouches)

  logger.Debug("Regex string " + regexString)

  aIdx := v.trustedServers[rIdx][fivetouches[0]]
  bIdx := v.trustedServers[rIdx][fivetouches[1]]
  cIdx := v.trustedServers[rIdx][fivetouches[2]]
  dIdx := v.trustedServers[rIdx][fivetouches[3]]
  eIdx := v.trustedServers[rIdx][fivetouches[4]]

  logger.Debug("corner " + strconv.Itoa(corner) + " (" + strconv.Itoa(aIdx) + ","     + strconv.Itoa(bIdx) + "," + strconv.Itoa(cIdx) + "," + strconv.Itoa(dIdx) + "," + strconv.Itoa(eIdx) + ")")

  os.Exit(1)

  tickets, err := v.ReceiveTickets(rIdx, ccs, corner)
  if err != nil {
    logger.Debug("Failed to Get Tickets")
    return false, nil
  }


  fmt.Printf("t=%v\n", tickets)



  return true, nil
}


func (v *FrackFixer) initNeighbours() *error.Error {
  side := math.Sqrt(config.TOTAL_RAIDA_NUMBER)
  sideSize := int(side)

  if sideSize * sideSize != config.TOTAL_RAIDA_NUMBER {
    return &error.Error{config.ERROR_INTERNAL, "Invalid Configuration"}
  }

	for i := 0; i < config.TOTAL_RAIDA_NUMBER; i++ {
		v.trustedServers[i][0] = v.getNeightbour(i, -sideSize - 1);
    v.trustedServers[i][1] = v.getNeightbour(i, -sideSize);
    v.trustedServers[i][2] = v.getNeightbour(i, -sideSize + 1);
    v.trustedServers[i][3] = v.getNeightbour(i, -1);
    v.trustedServers[i][4] = v.getNeightbour(i, 1);
    v.trustedServers[i][5] = v.getNeightbour(i, sideSize - 1);
    v.trustedServers[i][6] = v.getNeightbour(i, sideSize);
    v.trustedServers[i][7] = v.getNeightbour(i, sideSize + 1);

		// Five raida servers. Each number is an index in trustedServers array. 
		j := 0;

    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 3, 4 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 3, 5 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 3, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 3, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 4, 5 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 4, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 4, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 2, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 3, 4, 5 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 3, 4, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 3, 4, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 3, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 3, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 3, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 4, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 4, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 4, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 1, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 3, 4, 5 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 3, 4, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 3, 4, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 3, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 3, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 3, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 4, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 4, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 4, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 2, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 3, 4, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 3, 4, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 3, 4, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 3, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 0, 4, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 3, 4, 5 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 3, 4, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 3, 4, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 3, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 3, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 3, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 4, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 4, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 4, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 2, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 3, 4, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 3, 4, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 3, 4, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 3, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 1, 4, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 2, 3, 4, 5, 6 }; j++
    v.trustedTriads[i][j] = [5]int{ 2, 3, 4, 5, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 2, 3, 4, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 2, 3, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 2, 4, 5, 6, 7 }; j++
    v.trustedTriads[i][j] = [5]int{ 3, 4, 5, 6, 7 };
	}
  return nil
}

func (v *FrackFixer) getNeightbour(raidaIdx int, offset int) int {
	result := raidaIdx + offset

	if result < 0 {
		result += config.TOTAL_RAIDA_NUMBER
	}

	if result >= config.TOTAL_RAIDA_NUMBER {
		result -= config.TOTAL_RAIDA_NUMBER
	}

	return result
}
