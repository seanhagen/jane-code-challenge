package games

import (
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

/**
 * File: rankings.go
 * Date: 2021-11-15 15:36:26
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// StartMatchDay is the day number that matches start on
const StartMatchDay = 1

// Ranking ...
type Ranking struct {
	Teams        map[string]*team
	Days         map[int]matchDay
	currentMatch *matchDay
}

// NewRanking ...
func NewRanking() *Ranking {
	r := Ranking{
		Teams: map[string]*team{},
		Days:  map[int]matchDay{},
	}
	cm := newMatchDay(&r, StartMatchDay)
	r.Days[StartMatchDay] = cm
	r.currentMatch = &cm
	return &r
}

// newMatchDay ...
func (r Ranking) newMatchDay() {
	if r.currentMatch == nil {
		cm := newMatchDay(&r, StartMatchDay)
		r.Days[StartMatchDay] = cm
		r.currentMatch = &cm
		spew.Dump(cm)
		return
	}

	nd := r.currentMatch.Day + 1
	nm := newMatchDay(&r, nd)
	spew.Dump(nm)
	r.Days[nd] = nm
	r.currentMatch = &nm
}

// currentDay ...
func (r Ranking) currentDay() int {
	return r.currentMatch.Day
}

// AddMatch ...
func (r Ranking) AddMatch(in string) error {
	err := r.parseMatchLine(in)
	if err != nil {
		if _, ok := err.(*TeamPlayedError); ok {
			r.newMatchDay()
			return r.AddMatch(in)
		}
		return err
	}
	return nil
}

// parseMatchLine ...
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

	err = r.currentMatch.processMatchResults(t1, t2)
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

// parseTeamScore ...
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
