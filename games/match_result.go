package games

import "fmt"

/**
 * File: match_result.go
 * Date: 2021-11-19 14:53:33
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

type matchResult struct {
	outcome string
}

// String ...
func (tmr matchResult) String() string {
	return tmr.outcome
}

var (
	matchUnknown = matchResult{""}
	matchWon     = matchResult{"won"}
	matchLost    = matchResult{"lost"}
	matchTied    = matchResult{"tied"}
)

func matchResultFromString(s string) (matchResult, error) {
	switch s {
	case matchWon.outcome:
		return matchWon, nil
	case matchLost.outcome:
		return matchLost, nil
	case matchTied.outcome:
		return matchTied, nil
	}

	return matchUnknown, fmt.Errorf("unknown result: %v", s)
}
