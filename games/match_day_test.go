package games

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

/**
 * File: day_test.go
 * Date: 2021-11-15 15:11:37
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

func TestGames_MatchDay_ProcessMatchResults(t *testing.T) {
	type m struct {
		t1 string
		t2 string
	}

	tests := []struct {
		matches []m
		okay    bool
		bad     string
	}{
		{[]m{{"A 2", "B 2"}, {"C 0", "D 1"}}, true, ""},
		{[]m{{"A 2", "B 3"}, {"C 0", "A 1"}}, false, "A"},
		{[]m{{"A 2", "B 3"}, {"C 0", "D 1"}, {"B 2", "E 1"}}, false, "B"},
	}

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test%v", i), func(t *testing.T) {
			r := NewRanking()
			m := r.currentMatch

			okay := true
			for _, v := range tt.matches {
				t1, err := r.parseTeamScore(v.t1)
				if err != nil {
					// we're not testing the parseTeamScore function here, so if it
					// throws an error then that's an issue to be solved in that test
					t.Fatalf("should not be an error here: %v", err)
				}

				t2, err := r.parseTeamScore(v.t2)
				if err != nil {
					// same deal as above
					t.Fatalf("should not be an error here: %v", err)
				}

				// okay, test stuff now
				err = m.processMatchResults(t1, t2)
				if err != nil {
					okay = false
					if tt.okay {
						t.Errorf("Unable to add match '%v'; got error: %v", v, err)
					} else {
						x := err.(*TeamPlayedError)
						if x.name != tt.bad {
							t.Errorf("Wrong team in error, expected '%v' got '%v'", tt.bad, x.name)
						}
					}
				}
			}

			if !tt.okay && okay {
				t.Errorf("Expected error, got nothing")
				spew.Dump(m)
			}
		})
	}
}
