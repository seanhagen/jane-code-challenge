package games

import "fmt"

/**
 * File: errors.go
 * Date: 2021-11-15 17:38:27
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// TeamPlayedError ...
type TeamPlayedError struct {
	name string
}

// Error ...
func (tpe TeamPlayedError) Error() string {
	return fmt.Sprintf("team '%v' already played today", tpe.name)
}

// =========================================================

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

// =========================================================

// ParseLineError ...
type ParseLineError struct {
	l string
}

// Error ...
func (ple ParseLineError) Error() string {
	return fmt.Sprintf("wrong number of parts in match string '%v'", ple.l)
}

// for when match_day encounters an error when calling `recordGame` on a team
type RecordGameError struct {
	day  int
	team string
	err  error
}

// Error ...
func (rge RecordGameError) Error() string {
	return fmt.Sprintf("error recording game for team '%v' on day %v, reason: %v", rge.team, rge.day, rge.err)
}
