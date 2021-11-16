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

// ParseTeamError ...
type ParseTeamError struct {
	empty bool
	score string
	err   error
}

// Error ...
func (pe ParseTeamError) Error() string {
	if pe.empty {
		return fmt.Sprintf("given empty team result string to parse")
	}
	return fmt.Sprintf("unable to parse '%v' for score: %v", pe.score, pe.err)
}

// Unwrap ...
func (pe ParseTeamError) Unwrap() error {
	return pe.err
}

// ParseLineError ...
type ParseLineError struct {
	l string
}

// Error ...
func (ple ParseLineError) Error() string {
	return fmt.Sprintf("wrong number of parts in match string '%v'", ple.l)
}

// Match ...
type Match struct {
	TeamOne *MatchTeam
	TeamTwo *MatchTeam
}

// MatchTeam ...
type MatchTeam struct {
	Name  string
	Score int
}

// ParseLine ...
func ParseLine(input string) (*Match, error) {
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("wrong number of parts in match string '%v'", input)
	}

	t1, err := parseTeam(parts[0])
	if err != nil {
		return nil, err
	}

	t2, err := parseTeam(parts[1])
	if err != nil {
		return nil, err
	}

	m := &Match{
		TeamOne: t1,
		TeamTwo: t2,
	}
	return m, nil
}

func parseTeam(in string) (*MatchTeam, error) {
	in = strings.TrimSpace(in)
	bits := strings.Split(in, " ")
	x := len(bits)
	if x <= 0 {
		return nil, &ParseTeamError{empty: true}
	}

	s := bits[x-1]
	name := strings.Join(bits[:x-1], " ")
	score, err := strconv.Atoi(s)
	if err != nil {
		return nil, &ParseTeamError{score: string(bits[x-1]), err: err}
	}
	return &MatchTeam{name, score}, nil
}
