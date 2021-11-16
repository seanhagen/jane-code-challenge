package games

import "fmt"

/**
 * File: day.go
 * Date: 2021-11-15 15:11:26
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

var (
	// ErrTeamPlayedToday is returned from AddMatch when one of the
	// teams involved in the match already played in this match day
	ErrTeamPlayedToday = fmt.Errorf("team played today already")
)

// TeamPlayedError ...
type TeamPlayedError struct {
	name string
}

// Error ...
func (tpe TeamPlayedError) Error() string {
	return fmt.Sprintf("team '%v' already played today", tpe.name)
}

// MatchDay ...
type MatchDay struct {
	// what day is this match on?
	Day int
	// each team that played and their score in said match
	Teams map[string]int
	// mapping of who played who -- teams will be
	// in here twice, both a value and a key (easier lookups)
	Matchups map[string]string
	Results  []*Match
}

// NewMatchDay ...
func NewMatchDay(d int) MatchDay {
	if d <= 0 {
		d = 1
	}

	return MatchDay{
		Day:      d,
		Teams:    map[string]int{},
		Matchups: map[string]string{},
		Results:  []*Match{},
	}
}

// teamPlayed ...
func (m *MatchDay) teamPlayed(tn string) bool {
	for k := range m.Teams {
		if k == tn {
			return true
		}
	}
	return false
}

// AddMatch ...
func (m *MatchDay) AddMatch(match *Match) error {
	if m.teamPlayed(match.TeamOne.Name) {
		return &TeamPlayedError{match.TeamOne.Name}
	}

	if m.teamPlayed(match.TeamTwo.Name) {
		return &TeamPlayedError{match.TeamTwo.Name}
	}

	mo := match.TeamOne.Name
	mt := match.TeamTwo.Name
	so := match.TeamOne.Score
	st := match.TeamTwo.Score

	m.Teams[mo] = so
	m.Teams[mt] = st

	m.Matchups[mo] = mt
	m.Matchups[mt] = mo

	m.Results = append(m.Results, match)

	return nil
}

// AddMatchLine ...
func (m *MatchDay) AddMatchLine(in string) error {
	match, err := ParseLine(in)
	if err != nil {
		return err
	}
	return m.AddMatch(match)
}
