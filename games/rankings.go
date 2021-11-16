package games

/**
 * File: rankings.go
 * Date: 2021-11-15 15:36:26
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// Rankings ...
type Rankings struct {
	Teams      map[string]Team
	Days       map[int]MatchDay
	CurrentDay int
}

// CreateRankings ...
func CreateRankings() Rankings {
	return Rankings{
		Teams:      map[string]Team{},
		Days:       map[int]MatchDay{},
		CurrentDay: 1,
	}
}
