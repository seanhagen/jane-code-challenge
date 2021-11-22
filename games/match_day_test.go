package games

import (
	"fmt"
	"testing"

	"github.com/andreyvit/diff"
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
				// spew.Dump(m)
			}
		})
	}
}

func TestGames_MatchDay_TestStringOutput(t *testing.T) {
	buildTeam := func(n string) *team {
		return &team{
			Name:     n,
			Played:   map[int]string{},
			Scores:   map[int]int{},
			Standing: map[int]int{},
		}
	}

	buildTeamDayTwo := func(n string, st int) *team {
		return &team{
			Name:          n,
			Played:        map[int]string{1: "test"},
			Scores:        map[int]int{1: 1},
			Standing:      map[int]int{1: st},
			lastDayPlayed: 1,
		}
	}

	buildTeamDayThree := func(n string, st2 int) *team {
		return &team{
			Name:          n,
			Played:        map[int]string{1: "team 1", 2: "team 2"},
			Scores:        map[int]int{1: 1, 2: 2},
			Standing:      map[int]int{1: 0, 2: st2},
			lastDayPlayed: 2,
		}
	}

	type m struct {
		team1, team2 *team
		s1, s2       int
	}

	tests := []struct {
		day     int
		matches []m
		output  string
	}{
		{
			1,
			[]m{{buildTeam("A"), buildTeam("B"), 2, 1}, {buildTeam("D"), buildTeam("C"), 2, 1}},
			// A vs B => 2:1 => A 3, B 0
			// D vs C => 2:1 => D 3, C 0
			`Matchday 1
A, 3 pts
D, 3 pts
B, 0 pts
`,
		},

		{
			2,
			[]m{
				{buildTeamDayTwo("A", 3), buildTeamDayTwo("D", 3), 1, 2},
				{buildTeamDayTwo("B", 1), buildTeamDayTwo("C", 1), 3, 1},
			},
			// A(3) vs D(3) => 1:2 => A +0, D +3 => A 3 : D 6
			// B(1) vs C(1) => 3:1 => B +3, C +0 => B 4 : C 1
			`Matchday 2
D, 6 pts
B, 4 pts
A, 3 pts
`,
		},

		{
			3,
			[]m{
				{buildTeamDayThree("A", 3), buildTeamDayThree("C", 0), 1, 2},
				{buildTeamDayThree("D", 6), buildTeamDayThree("B", 3), 1, 1},
				// A(3) vs C(0) => 1:2 => A 3, C 3
				// B(3) vs D(6) => 1:1 => B 4, D 7
			},
			`Matchday 3
D, 7 pts
B, 4 pts
A, 3 pts
`,
		},

		{
			3,
			[]m{
				{buildTeamDayThree("A", 2), buildTeamDayThree("C", 6), 1, 2},
				{buildTeamDayThree("D", 3), buildTeamDayThree("B", 2), 1, 1},
				// A(3) vs C(0) => 1:2 => A +0, C +3 => A 3, C 3
				// B(3) vs D(6) => 1:1 => B +1, D +1 => B 4, D 7
			},
			`Matchday 3
C, 9 pts
D, 4 pts
B, 3 pts
`,
		},
	}

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("day_%v_test_%v_", tt.day, i), func(t *testing.T) {
			md := newMatchDay(tt.day)

			for ii, xx := range tt.matches {
				tr1 := &teamResult{xx.team1, xx.s1}
				tr2 := &teamResult{xx.team2, xx.s2}

				// this function is tested elsewhere, so this should always succeed with the inputs we've got
				err := md.processMatchResults(tr1, tr2)
				if err != nil {
					t.Fatalf("unable to add match result %v, reason: %v", ii, err)
				}
			}

			got := md.Results()
			if tt.output != got {
				d := diff.LineDiff(tt.output, got)
				t.Errorf("wrong output\ndiff:\n%v", d)
			}
		})
	}
}
