package games

import (
	"fmt"
	"testing"
)

/**
 * File: rankings_test.go
 * Date: 2021-11-15 16:14:37
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

func TestGames_Ranking_ParseMatchLine(t *testing.T) {
	tests := []struct {
		line   string
		t1, t2 string
		s1, s2 int
		ok     bool
	}{
		{"San Jose Earthquakes 3, Santa Cruz Slugs 3", "San Jose Earthquakes", "Santa Cruz Slugs", 3, 3, true},
		{"Vancouver Bears 1, Calgary Lions 0", "Vancouver Bears", "Calgary Lions", 1, 0, true},
		{"A 2, B 3", "A", "B", 2, 3, true},
		{"A 22, B 35", "A", "B", 22, 35, true},
		{"      A 2,           B 3", "A", "B", 2, 3, true},
		{"A B C D 5,D E F G H 4", "A B C D", "D E F G H", 5, 4, true},
		{"A B C D 5,       D E F G H 4", "A B C D", "D E F G H", 5, 4, true},
		{"A A A A A A 1, ", "A", "", 1, -1, false},
		{",B B B 2", "", "B", -1, -1, false},
		{"A, 2, B, 3", "A", "B", -1, -1, false},
		{"A, 2 B 3", "A", "B", -1, -1, false},
		{"A 2 B 3", "A", "B", -1, -1, false},
		{"", "", "", -1, -1, false},
	}

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test_%v_", i), func(t *testing.T) {
			r := NewRanking()
			err := r.parseMatchLine(tt.line)
			if tt.ok && err != nil {
				t.Fatalf("expected okay, got error: %v", err)
			}

			if tt.ok && err == nil {
				t1, ok := r.Teams[tt.t1]
				if t1 == nil || !ok {
					t.Fatalf("team '%v' is not present in rankings", tt.t1)
				}

				s1 := r.currentMatch.Teams[tt.t1]
				if s1 != tt.s1 {
					t.Errorf("Wrong score for team '%v', expected '%v' got '%v'", tt.t1, tt.s1, s1)
				}

				t2, ok := r.Teams[tt.t2]
				if t2 == nil || !ok {
					t.Fatalf("team '%v' is not present in rankings", tt.t2)
				}

				s2 := r.currentMatch.Teams[tt.t2]
				if s2 != tt.s2 {
					t.Errorf("Wrong score for team '%v', expected '%v' got '%v'", tt.t2, tt.s2, s2)
				}

				o1 := r.currentMatch.Matchups[tt.t1]
				if o1 != tt.t2 {
					t.Errorf("Wrong opponent recorded for team '%v', expected '%v' got '%v'", tt.t1, tt.t2, o1)
				}

				o2 := r.currentMatch.Matchups[tt.t2]
				if o2 != tt.t1 {
					t.Errorf("Wrong opponent recorded for team '%v', expected '%v' got '%v'", tt.t2, tt.t1, o2)
				}

				if s1 > s2 {
					if t1.currentRank() != 3 {
						t.Errorf("expected rank for '%v' after win to be 3, rank is: %v", t1.Name, t1.currentRank())
					}
					if t2.currentRank() != 0 {
						t.Errorf("expected rank for '%v' after loss to be 0, rank is: %v", t2.Name, t2.currentRank())
					}
				} else if s2 > s1 {
					if t1.currentRank() != 0 {
						t.Errorf("expected rank for '%v' after loss to be 0, rank is: %v", t1.Name, t1.currentRank())
					}
					if t2.currentRank() != 3 {
						t.Errorf("expected rank for '%v' after win to be 3, rank is: %v", t2.Name, t2.currentRank())
					}
				} else {
					if t1.currentRank() != 1 {
						t.Errorf("expected rank for '%v' after tie to be 1, rank is: %v", t1.Name, t1.currentRank())
					}
					if t2.currentRank() != 1 {
						t.Errorf("expected rank for '%v' after tie to be 1, rank is: %v", t2.Name, t2.currentRank())
					}
				}
			}
		})
	}
}

func TestGames_Ranking_TestParseAndOutput(t *testing.T) {
	inputs := []string{
		"San Jose Earthquakes 3, Santa Cruz Slugs 3",
		"Capitola Seahorses 1, Aptos FC 0",
		"Felton Lumberjacks 2, Monterey United 0",
	}

	expect := `Matchday 1
Capitola Seahorses, 3 pts
Felton Lumberjacks, 3 pts
San Jose Earthquakes, 1 pt
`

	r := NewRanking()
	for _, tt := range inputs {
		r.AddMatch(tt)
	}

	output := r.getCurrentMatchDay().Results()
	if expect != output {
		t.Errorf("output incorrect.\nexpected: \n-----\n%v\n-----\n\nrecieved: \n-----\n%v\n-----\n", expect, output)
	}
}

func TestGames_Ranking_TestAddMatch(t *testing.T) {

	tests := []struct {
		inputs  []string
		okay    []bool // needs to be same length as inputs
		numDays int
	}{
		{
			[]string{"A 1, B 1"},
			[]bool{true},
			1,
		},
		{
			[]string{"A 1, B 1", "C 1, D 2"},
			[]bool{true, true},
			1,
		},
		{
			[]string{"A 1 B 2 C 3"},
			[]bool{false},
			1,
		},
		{
			[]string{"A 1, B 1", "C 1, D 2", "A 2, C 1", "B 1, D 3"},
			[]bool{true, true, true, true},
			2,
		},
	}

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test_%v_", i), func(t *testing.T) {
			if len(tt.inputs) != len(tt.okay) {
				t.Fatalf("mismatched inputs & errors, should be same length; inputs: %v, errors: %v", len(tt.inputs), len(tt.okay))
			}

			r := NewRanking()

			for i := 0; i < len(tt.inputs); i++ {
				in := tt.inputs[i]
				ok := tt.okay[i]
				err := r.AddMatch(in)

				if ok && err != nil {
					t.Errorf("unable to add line '%v'; error: %v", in, err)
				}

				if !ok && err == nil {
					t.Errorf("expected error for line '%v', got nothing", in)
				}
			}

			cd := r.currentDay
			if cd != tt.numDays {
				t.Errorf("wrong day, expected '%v' got '%v'", tt.numDays, cd)
			}
		})
	}
}
