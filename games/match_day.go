package games

import (
	"fmt"
	"sort"
	"strings"
)

/**
 * File: day.go
 * Date: 2021-11-15 15:11:26
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// matchDay records the results of a day of matchups
type matchDay struct {
	// what day is this match on?
	Day int

	// each team that played and their score in said match
	Teams map[string]int

	// mapping of who played who -- teams will be
	// in here twice, both a value and a key (easier lookups)
	Matchups map[string]string

	// each team and their points total at the end of the day
	Standings standingList
}

// newMatchDay is the matchDay constructor
func newMatchDay(d int) matchDay {
	if d <= 0 {
		d = 1
	}

	return matchDay{
		Day:       d,
		Teams:     map[string]int{},
		Matchups:  map[string]string{},
		Standings: standingList{},
	}
}

// teamPlayed returns true if the given team has
// already played in this match day
func (m *matchDay) teamPlayed(tn string) bool {
	for k := range m.Teams {
		if k == tn {
			return true
		}
	}
	return false
}

// processMatchResults takes the parsed results from a matchup
// and process and stores the results
//
// It determines if the match resulted in a tie or not, and tells
// each team in the match to record the game and records the rank
// of each team after the match
func (m *matchDay) processMatchResults(t1, t2 *teamResult) error {
	mo := t1.team.Name
	so := t1.score

	mt := t2.team.Name
	st := t2.score

	if m.teamPlayed(t1.team.Name) {
		return &TeamPlayedError{t1.team.Name}
	}

	if m.teamPlayed(t2.team.Name) {
		return &TeamPlayedError{t2.team.Name}
	}

	m.Teams[mo] = so
	m.Teams[mt] = st

	m.Matchups[mo] = mt
	m.Matchups[mt] = mo

	// add results to current match
	r1 := matchWon
	r2 := matchLost

	if t1.score < t2.score {
		r1, r2 = r2, r1
	} else if t1.score == t2.score {
		r1, r2 = matchTied, matchTied
	}

	rank1, err := t1.team.recordGame(m.Day, t2.team.Name, t1.score, r1)
	if err != nil {
		return &RecordGameError{m.Day, t1.team.Name, err}
	}
	rank2, err := t2.team.recordGame(m.Day, t1.team.Name, t2.score, r2)
	if err != nil {
		return &RecordGameError{m.Day, t2.team.Name, err}
	}

	add := []standing{standing{teamName: mo, rank: rank1}, standing{teamName: mt, rank: rank2}}
	m.Standings = append(m.Standings, add...)

	return nil
}

// Results is the nicely formatted results of the match day,
// showing the top three teams in point standings for this day
func (m matchDay) Results() string {
	out := fmt.Sprintf("Matchday %v\n", m.Day)
	l := len(m.Standings)
	if l > 3 {
		l = 3
	}

	sort.Sort(m.Standings)
	for i := 0; i < l; i++ {
		t := m.Standings[i]
		s := "pt"
		if t.rank > 1 || t.rank == 0 {
			s = "pts"
		}
		out = fmt.Sprintf("%v%v, %v %v\n", out, t.teamName, t.rank, s)
	}

	return fmt.Sprintf("%v\n", out)
}

// String is for the Stringer interface, a more compact version
// of Results(), and includes all teams in the output
func (m matchDay) String() string {
	// helper function that's not needed elsewhere,
	// so bundle it up inside here to make sure that's clear
	stringArrContain := func(t string, in []string) bool {
		for _, v := range in {
			if v == t {
				return true
			}
		}
		return false
	}

	out := fmt.Sprintf("\nMatch Day: %v", m.Day)
	teams := []string{}
	matchups := []string{}

	for k, v := range m.Matchups {
		if !stringArrContain(k, teams) {
			teams = append(teams, k, v)
			ks := m.Teams[k]
			vs := m.Teams[v]

			matchups = append(matchups, fmt.Sprintf("%v vs %v (Score: %v-%v)", k, v, ks, vs))
		}
	}
	out = fmt.Sprintf("%v\n\tMatchups: \t%v\n\tStandings:\t", out, strings.Join(matchups, ", "))

	standings := []string{}
	for _, v := range m.Standings {
		standings = append(standings, fmt.Sprintf("%v: %v", v.teamName, v.rank))
	}
	out = fmt.Sprintf("%v%v", out, strings.Join(standings, ", "))
	return out
}
