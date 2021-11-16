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
	// Rank is the numerical rank of the team. Can be shared
	// by other teams ( ie, possible for more than one team
	// to be rank N )
	Rank int
	// Played keeps track of who this team played each match day.
	// The key is the match day, the value is the team that was played
	Played map[int]string
	// Scores keeps track of what the team scored on any particular
	// match day.
	Scores map[int]int
}

// findOrCreateTeam ...
func (r Ranking) findOrCreateTeam(n string) *team {
	t, ok := r.Teams[n]
	if ok {
		return t
	}

	t = &team{Name: n, Rank: 0, Played: map[int]string{}, Scores: map[int]int{}}
	r.Teams[n] = t
	return t
}

// findTeam ...
func (r Ranking) findTeam(n string) *team {
	for k, v := range r.Teams {
		if k == n {
			return v
		}
	}
	return nil
}

// played ...
func (t *team) recordGame(day int, name string, score int, res teamMatchResult) error {
	t.Played[day] = name
	t.Scores[day] = score

	switch res {
	case matchWon:
		t.Rank += 3
	case matchTied:
		t.Rank++
	case matchLost:
	default:
		return fmt.Errorf("invalid match result '%v'", res)
	}
	return nil
}
