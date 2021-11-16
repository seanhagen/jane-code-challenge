package games

/**
 * File: team.go
 * Date: 2021-11-15 14:38:25
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// Team is a team found within the input, used to track
// it's rank, as well as who they played each day and what
// their score was on that day.
type Team struct {
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
