package games

import "fmt"

/**
 * File: match.go
 * Date: 2021-11-15 14:34:50
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// MatchDay ...
type MatchDay struct {
	// what day is this match on?
	Day int
	// each team that played and their final score
	Teams map[string]int
	// mapping of who played who -- teams will be
	// in here twice, both a value and a key (easier lookups)
	Matchups map[string]string
	Results  []Match
}

// Match ...
type Match struct {
	T1, T2 string
	S1, S2 int
}

// MatchTeam ...
type MatchTeam struct {
	Name  string
	Score int
}

// ParseLine ...
func ParseLine(input string) (*Match, error) {
	return nil, fmt.Errorf("not yet")
}
