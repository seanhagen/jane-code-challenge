package games

import (
	"fmt"
)

/**
 * File: team.go
 * Date: 2021-11-15 14:38:25
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// team is a team found within the input, used to track
// it's rank, as well as who they played each day and what
// their score was on that day.
type team struct {
	// Name is the team name
	Name string

	// Played keeps track of who this team played each match day.
	// The key is the match day, the value is the team that was played
	Played map[int]string

	// Scores keeps track of what the team scored on any particular
	// match day.
	Scores map[int]int

	// Standing keeps track of what this teams point total was on each day
	Standing map[int]int

	// lastDayPlayed keeps track of the last day this team played on
	lastDayPlayed int
}

// findOrCreateTeam looks up the team by name, and creates a new team
// struct object if that team doesn't exist within the rankings
func (r Ranking) findOrCreateTeam(n string) *team {
	t, ok := r.Teams[n]
	if ok {
		return t
	}

	t = &team{Name: n, Played: map[int]string{}, Scores: map[int]int{}, Standing: map[int]int{}}
	r.Teams[n] = t
	return t
}

// recordGame records the score and outcome of this team in matchup
// on a given day, against the named team.
func (t *team) recordGame(day int, name string, score int, res matchResult) (int, error) {
	if day == t.lastDayPlayed {
		return -1, fmt.Errorf("already played today")
	}
	if day < 1 {
		return -1, fmt.Errorf("invalid day, can't be less than 1")
	}
	if day != t.lastDayPlayed+1 {
		return -1, fmt.Errorf("invalid day, got %v, next day should be %v", day, t.lastDayPlayed+1)
	}

	rank := t.Standing[t.lastDayPlayed]
	switch res {
	case matchWon:
		rank += 3
	case matchTied:
		rank++
	case matchLost:
	default:
		return -1, fmt.Errorf("invalid match result '%v'", res)
	}

	t.Played[day] = name
	t.Scores[day] = score
	t.lastDayPlayed++
	t.Standing[t.lastDayPlayed] = rank
	return rank, nil
}

// currentRank ...
func (t *team) currentRank() int {
	return t.Standing[t.lastDayPlayed]
}

// String  ...
func (t team) String() string {
	out := fmt.Sprintf("\nTeam '%v'", t.Name)
	max := 0
	for k := range t.Played {
		if k > max {
			max = k
		}
	}

	for i := 0; i <= max; i++ {
		out = fmt.Sprintf("%v\n\tDay %v: Played %v, Score: %v, Rank After Game: %v", out, i, t.Played[i], t.Scores[i], t.Standing[i])
	}

	return out
}
