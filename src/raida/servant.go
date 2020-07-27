package raida

import(
	"fmt"
	"logger"
	"strconv"
	"strings"
	"config"
)

type Servant struct {
	Raida *RAIDA
}

type Error struct {
	Message string
}

func NewServant() (*Servant) {
	//fmt.Println("new servant")
	Raida := New()
	logger.Info("Raida initialized. Total Servers "  + strconv.Itoa(Raida.TotalServers()))

	return &Servant{
		Raida: Raida,
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
