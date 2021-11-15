package games

/**
 * File: team.go
 * Date: 2021-11-15 14:38:25
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// Team ...
type Team struct {
	Name    string
	Results map[int]int
	Played  map[int]string
}
