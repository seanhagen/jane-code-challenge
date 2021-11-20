package games

import (
	"fmt"
	"strconv"
	"strings"
)

/**
 * File: rankings.go
 * Date: 2021-11-15 15:36:26
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// StartMatchDay is the day number that matches start on
const StartMatchDay = 1

// MaxDepth is how many times we'll try to add a new match
// day when we get a "team already played" error. Mostly just
// guard rails for "in case".
const MaxDepth = 3

// Ranking is our list object that handles parsing input lines,
// assigning the results to match days & creating new match
// days as appropriate.
type Ranking struct {
	Teams        map[string]*team
	Days         map[int]*matchDay
	matches      []*matchDay
	currentMatch *matchDay
	currentDay   int
}

// NewRanking is the constructor for Ranking structs
func NewRanking() *Ranking {
	r := Ranking{
		Teams:   map[string]*team{},
		Days:    map[int]*matchDay{},
		matches: []*matchDay{},
	}
	r.newMatchDay(StartMatchDay)
	return &r
}

// newMatchDay sets up a new match day
func (r *Ranking) newMatchDay(nd int) {
	if r.currentMatch == nil {
		r.currentDay = 1
	}

	nm := newMatchDay(nd)
	r.matches = append(r.matches, &nm)
	r.currentMatch = &nm
	r.Days[nd] = r.currentMatch
	r.currentDay = nd
}

// getCurrentMatchDay creates a match day if one doesn't exist
func (r *Ranking) getCurrentMatchDay() *matchDay {
	if r.currentMatch == nil {
		if len(r.matches) <= 0 {
			// okay, something's gone screwy
			// we don't have a current match day, so create a
			// new one and reset the state of the ranking
			r.Teams = map[string]*team{}
			r.Days = map[int]*matchDay{}
			r.matches = []*matchDay{}
			r.newMatchDay(StartMatchDay)
		} else {
			// we've got days, but the *currentMatch pointer is
			// nil, so set it to the last day
			d := r.Days[0]
			for _, v := range r.Days {
				if v.Day > d.Day {
					d = v
				}
			}
			r.currentMatch = d
		}
	}
	return r.currentMatch
}

// AddMatch parses a match string, creating a new match day
// if either of the teams in the string have already played today
func (r *Ranking) AddMatch(in string) error {
	return r._addMatch(in, 0)
}

// _addMatch does the actual work for AddMatch, with a guard
// against infinite recursion, just in case.
func (r *Ranking) _addMatch(in string, depth int) error {
	if depth > 3 {
		return fmt.Errorf("delved too deep")
	}

	err := r.parseMatchLine(in)
	if err != nil {
		if _, ok := err.(*TeamPlayedError); ok {
			r.newMatchDay(r.currentDay + 1)
			return r._addMatch(in, depth+1)
		}
		return err
	}
	return nil
}

// parseMatchLine parses a match line in the form "<team 1> <score 1>, <team 2> <score 2>",
// and then tells the current match day to process the match results
func (r Ranking) parseMatchLine(input string) error {
	parts := strings.FieldsFunc(input, func(r rune) bool { return r == ',' })
	if len(parts) != 2 {
		return &ParseLineError{input}
	}

	t1, err := r.parseTeamScore(parts[0])
	if err != nil {
		return err
	}

	t2, err := r.parseTeamScore(parts[1])
	if err != nil {
		return err
	}

	cm := r.getCurrentMatchDay()
	err = cm.processMatchResults(t1, t2)
	if err != nil {
		return err
	}

	return nil
}

// teamResult ...
type teamResult struct {
	team  *team
	score int
}

// parseTeamScore parses a section of a match day
// string in the form "<team> <score>" and returns
// the result.
//
// This function creates a new team if the named
// team isn't already known to this Ranking struct.
func (r Ranking) parseTeamScore(in string) (*teamResult, error) {
	in = strings.TrimSpace(in)
	bits := strings.Fields(in)
	x := len(bits)
	if x <= 0 {
		return nil, &ParseTeamError{empty: true}
	}

	name := strings.Join(bits[0:x-1], " ")
	score, err := strconv.Atoi(string(bits[x-1]))
	if err != nil {
		return nil, &ParseTeamError{score: string(bits[x-1]), err: err}
	}

	t := r.findOrCreateTeam(name)
	return &teamResult{team: t, score: score}, nil
}

// Results ...
func (r Ranking) Results() string {
	out := ""

	for d := 0; d < r.currentDay; d++ {
		out = fmt.Sprintf("%v%v", out, r.matches[d].Results())
	}

	return out
}

// String is for the Stringer interface, so that
// we've got a nicer output in spew.Dump and whatnot.
func (r Ranking) String() string {
	out := fmt.Sprintf("\nRankings -- %v Teams over %v Days\nTeams:", len(r.Teams), len(r.Days))

	teams := []string{}
	for _, v := range r.Teams {
		teams = append(teams, v.String())
	}

	out = fmt.Sprintf("%v%v\n", out, strings.Join(teams, "\n"))

	return out
}
