package games

import (
	"fmt"
	"strconv"
	"strings"
)

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
	m := &Match{S1: -1, S2: -1}

	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return m, fmt.Errorf("wrong number of parts in match string '%v'", input)
	}

	t1, s1, err := parseTeam(parts[0])
	if err != nil {
		return m, err
	}

	t2, s2, err := parseTeam(parts[1])
	if err != nil {
		return m, err
	}

	m.T1 = t1
	m.T2 = t2
	m.S1 = s1
	m.S2 = s2

	return m, nil
}

func parseTeam(in string) (string, int, error) {
	in = strings.TrimSpace(in)
	bits := strings.Split(in, " ")
	x := len(bits)
	if x <= 0 {
		return "", -1, fmt.Errorf("given empty team result string to parse")
	}

	s := bits[x-1]
	name := strings.Join(bits[:x-1], " ")
	score, err := strconv.Atoi(s)
	if err != nil {
		return "", -1, fmt.Errorf("unable to parse score '%v': %v", s, err)
	}
	return name, score, nil
}
