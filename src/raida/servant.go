package raida

import(
	"fmt"
	"logger"
	"strconv"
	"strings"
	"config"
	"cloudcoin"
	"error"
)

type Servant struct {
	Raida *RAIDA
	repairArray [][]int
}

type Error struct {
	Message string
}

func NewServant() (*Servant) {
	//fmt.Println("new servant")
	Raida := New()
	logger.Info("Raida initialized. Total Servers "  + strconv.Itoa(Raida.TotalServers()))

	repairArray := make([][]int, Raida.TotalServers())
	//repairArray[5] = append(repairArray[5], 10)
	//repairArray[5] = append(repairArray[5], 20)
	//fmt.Printf("r=%d\n",repairArray[5])

	return &Servant{
		Raida: Raida,
		repairArray: repairArray,
	}

}

func (s *Servant) GetPownStringFromStatusArray(statuses []int) string {
	var b strings.Builder
	var c string

	for _, status := range statuses {
		switch (status) {
			case config.RAIDA_STATUS_UNTRIED:
				c = "u"
			case config.RAIDA_STATUS_PASS:
				c = "p"
			case config.RAIDA_STATUS_FAIL:
				c = "f"
			case config.RAIDA_STATUS_ERROR:
				c = "e"
			case config.RAIDA_STATUS_NORESPONSE:
				c = "n"
			default:
				c = "e"
		}

		fmt.Fprintf(&b, "%s", c)
	}

	return b.String()
}

func (s *Servant) IsStatusArrayFixable(statuses []int) bool {
	return s.isStatusArrayFixableRows(statuses) && s.isStatusArrayFixableColumns(statuses)
}

func (s *Servant) isStatusArrayFixableRows(statuses []int) bool {
	return s.isStatusArrayFixableInternal(statuses)
}

func (s *Servant) isStatusArrayFixableColumns(statuses []int) bool {
	rotatedStatuses := make([]int, s.Raida.TotalServers())

	for i := 0; i < s.Raida.TotalServers(); i++ {
		idx := i * s.Raida.GetSideSize()
		multiplier := idx / s.Raida.TotalServers()
		idx -= (s.Raida.TotalServers() - 1) * multiplier

		rotatedStatuses[i] = statuses[idx]
	}
	return s.isStatusArrayFixableInternal(rotatedStatuses)

}

func (s *Servant) isStatusArrayFixableInternal(statuses []int) bool {
	var badRows, goodRows int
	var seenGoodRows bool

	badRows = 0
	goodRows = 0
	for i := 0; i < s.Raida.TotalServers(); i++ {
		if (statuses[i] == config.RAIDA_STATUS_PASS) {
			goodRows++
			badRows = 0
			if goodRows == s.Raida.GetSideSize() + 1 {
				seenGoodRows = true
			}
		} else {
			goodRows = 0
			badRows++
			if badRows == s.Raida.GetSideSize() {
				return false
			}
		}
	}

	if (seenGoodRows) {
		return true
	}
	
	return false
}

func (s *Servant) GetSNsOverlap(sns [][]int) ([]int, int) {
	logger.Debug("Getting overlapped SNs")

//	pownArray := make([]int, v.Raida.TotalServers())
	hm := make(map[int][]int)

	for ridx, snarray := range sns {
		for _, sn := range snarray {
			_, exists := hm[sn]
			if !exists {
				hm[sn] = make([]int, s.Raida.TotalServers())
			}

			hm[sn][ridx] = config.RAIDA_STATUS_PASS
		}
	}

	total := 0
	var rsns []int
	for sn, hme := range hm {
		logger.Debug("sn " + strconv.Itoa(sn) + " pownstring " + s.GetPownStringFromStatusArray(hme))
		if !s.IsStatusArrayFixable(hme) {
			logger.Debug("Skipping Coin " + strconv.Itoa(sn))
			continue
		}

		rsns = append(rsns, sn)
		total += cloudcoin.GetDenomination(sn)
	}

	return rsns, total
}

func (s *Servant) PickCoinsFromArray(sns []int, amount int) ([]int, int, *error.Error) {
	logger.Debug("Picking " + strconv.Itoa(amount) + "CC from array of coins, Total notes in array: " + strconv.Itoa(len(sns)))

	exps, err := s.GetExpCoins(sns, amount)
	if err != nil {
		return nil, 0, err
	}

	var collected, rest int
	collected = 0
	rest = 0

	logger.Debug("Go on")

	var picked []int
	for _, sn := range sns {
		denomination := cloudcoin.GetDenomination(sn)
		//logger.Debug("sn " + strconv.Itoa(sn) + " denom " + strconv.Itoa(denomination))
		if denomination == 1 {
			if exps[1] > 0 {
				exps[1]--
				picked = append(picked, sn)
				collected += denomination
			}
		} else if denomination == 5 {
			if exps[5] > 0 {
				exps[5]--
				picked = append(picked, sn)
				collected += denomination
			}
		} else if denomination == 25 {
			if exps[25] > 0 {
				exps[25]--
				picked = append(picked, sn)
				collected += denomination
			}
		} else if denomination == 100 {
			if exps[100] > 0 {
				exps[100]--
				picked = append(picked, sn)
				collected += denomination
			}
		} else if denomination == 250 {
			if exps[250] > 0 {
				exps[250]--
				picked = append(picked, sn)
				collected += denomination
			}
		}
	}

	coinsStr := fmt.Sprintf("%v", picked)
	logger.Debug("Picked " + coinsStr)

	rest = amount - collected;
	logger.Debug("rest = " + strconv.Itoa(rest))
	if rest == 0 {
		return picked, 0, nil
	}

	logger.Debug("Picking extra coin")

	var isAdded bool
	chosenSNforBreak := 0
	for  _, sn := range sns {
		denomination := cloudcoin.GetDenomination(sn)
		if (rest > denomination) {
			continue
		}

		isAdded = false
		for _, psn := range picked {
			if psn == sn {
				isAdded = true
				break
			}
		}

		if isAdded {
			logger.Debug("Skipping SN for breaking: " + strconv.Itoa(sn))
			continue
		}

		logger.Debug("Chosen for break: " + strconv.Itoa(sn))
		chosenSNforBreak = sn
		break
	}

	return picked, chosenSNforBreak, nil
}

func (s *Servant) GetExpCoins(sns []int, amount int) (map[int]int, *error.Error) {
	totals := make(map[int]int)

	total := 0
	totals[1] = 0
	totals[5] = 0
	totals[25] = 0
	totals[100] = 0
	totals[250] = 0
	for _, sn := range sns {
		denomination := cloudcoin.GetDenomination(sn)
		totals[denomination]++
		total += denomination
	}

	if (amount > total) {
		return nil, &error.Error{"Insufficient funds"}
	}

	savedAmount := amount
	for key, value := range totals {
		logger.Debug("d" + strconv.Itoa(key) + ": " + strconv.Itoa(value))
	}

	var exp_1, exp_5, exp_25, exp_100, exp_250 int
	exp_1 = 0
	exp_5 = 0
	exp_25 = 0
	exp_100 = 0
	exp_250 = 0

	for i := 0; i < 2; i++ {
		exp_1 = 0
		exp_5 = 0
		exp_25 = 0
		exp_100 = 0

		if i == 0 && amount >= 250 && totals[250] > 0 {
			if (amount / 250) < totals[250] {
				exp_250 = (amount / 250)
			} else {
				exp_250 = totals[250]
			}
			amount -= (exp_250 * 250)
		}

    if (amount >= 100 && totals[100] > 0) {
			if (amount / 100) < totals[100] {
				exp_100 = (amount / 100)
			} else {
				exp_100 = totals[100]
			}
			amount -= (exp_100 * 100)
    }

    if (amount >= 25 && totals[25] > 0) {
			if (amount / 25) < totals[25] {
				exp_25 = (amount / 25)
			} else {
				exp_25 = totals[25]
			}
			amount -= (exp_25 * 25);
    }

    if (amount >= 5 && totals[5] > 0) {
			if (amount / 5) < totals[5] {
				exp_5 = (amount / 5)
			} else {
				exp_5 = totals[5]
			}
			amount -= (exp_5 * 5);
    }

    if (amount >= 1 && totals[1] > 0) {
			if (amount / 1) < totals[1] {
				exp_1 = (amount / 1)
			} else {
				exp_1 = totals[1]
			}
			amount -= (exp_1);
    }

		logger.Debug("Picked Denom: " + strconv.Itoa(exp_1) + "/" + strconv.Itoa(exp_5) + "/" + strconv.Itoa(exp_25) + "/" + strconv.Itoa(exp_100) + "/" + strconv.Itoa(exp_250) + " rest amount = " + strconv.Itoa(amount))
		if amount == 0 {
			break
		}

		if i == 1 || exp_250 == 0 {
			break
		}

		exp_250--;
		amount = savedAmount - exp_250 * 250;
	}

	rv := make(map[int]int)
	rv[1] = exp_1
	rv[5] = exp_5
	rv[25] = exp_25
	rv[100] = exp_100
	rv[250] = exp_250

	return rv, nil
}

func (s *Servant) AddSNToRepairArray(raidaIdx int, sn int) {
	s.repairArray[raidaIdx] = append(s.repairArray[raidaIdx], sn)
}
