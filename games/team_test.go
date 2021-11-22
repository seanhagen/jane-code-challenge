package games

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

/**
 * File: team_test.go
 * Date: 2021-11-19 11:58:25
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

func TestGames_FindOrCreate(t *testing.T) {
	tests := []struct {
		name  string
		exist bool
	}{
		{"A", true},
		{"B", false},
	}

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test%v", i), func(t *testing.T) {
			r := NewRanking()
			r.Teams["X"] = &team{Name: "X"}

			var should *team

			if tt.exist {
				should = &team{Name: tt.name}
				r.Teams[tt.name] = should
			}

			get := r.findOrCreateTeam(tt.name)
			if get == nil {
				t.Fatalf("expected team, got nil")
			}

			if tt.exist {
				// team should already exist, so should & get should be pointers to the same thing
				if get != should {
					t.Errorf("expected *get to point to same address as *should, instead got *get => %p, *should => %p", get, should)
				}
			} else {
				// team should not exist, so *should & *get should not be pointers to the same thing
				if get == should {
					t.Errorf("expected *get to point to different address than *should, instead they're equal")
				}
			}

		})
	}
}

func TestGames_RecordGame(t *testing.T) {
	tests := []struct {
		name  string
		day   int
		opnt  string
		score int
		res   matchResult
		ret   int
		err   bool
		team  *team
	}{
		// first day, should succeed
		{"a_win", 1, "A", 3, matchWon, 3, false, nil},
		{"a_lose", 1, "A", 2, matchLost, 0, false, nil},
		{"a_tie", 1, "A", 2, matchTied, 1, false, nil},

		// second day, should succeed
		{"b_win", 2, "B", 3, matchWon, 6, false,
			&team{
				Played:        map[int]string{1: "A"},
				Scores:        map[int]int{1: 3},
				Standing:      map[int]int{1: 3},
				lastDayPlayed: 1,
			},
		},
		{"b_tie", 2, "B", 3, matchTied, 4, false,
			&team{
				Played:        map[int]string{1: "A"},
				Scores:        map[int]int{1: 3},
				Standing:      map[int]int{1: 3},
				lastDayPlayed: 1,
			},
		},
		{"b_lose", 2, "B", 3, matchLost, 3, false,
			&team{
				Played:        map[int]string{1: "A"},
				Scores:        map[int]int{1: 3},
				Standing:      map[int]int{1: 3},
				lastDayPlayed: 1,
			},
		},
		// testing error conditions

		// day-related errors
		{"a_already_played", 1, "A", 2, matchWon, -1, true, &team{lastDayPlayed: 1}}, // already played that day
		{"a_invalid_day_1", 0, "A", 2, matchWon, -1, true, &team{lastDayPlayed: 1}},  // invalid day
		{"a_invalid_day_2", -1, "A", 2, matchWon, -1, true, &team{lastDayPlayed: 1}}, // invalid day
		{"a_skipped_day", 3, "A", 2, matchWon, -1, true, &team{lastDayPlayed: 1}},    // skipping a day

		// invalid match result
		{"a_invalid_match_result", 1, "A", 2, matchResult{"nope"}, -1, true, nil},
	}

	for _, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test_%v", tt.name), func(t *testing.T) {
			rank := NewRanking()
			testTeam := tt.team
			if testTeam == nil {
				testTeam = rank.findOrCreateTeam("test")
				if testTeam == nil {
					t.Fatalf("unable to proceed, got nil team")
				}
			}

			r, err := testTeam.recordGame(tt.day, tt.opnt, tt.score, tt.res)

			if tt.err {
				if err == nil {
					t.Errorf("expected error result, got nil error")
				}
				if r > -1 {
					t.Errorf("expected error result, got something other than -1 as rank result: %v", r)
				}
			} else {
				if err != nil {
					t.Errorf("expected success, got error: %v", err)
				}
				if r == -1 {
					t.Errorf("expected updated rank, got -1 instead")
				}
				if r != tt.ret {
					t.Errorf("wrong rank result, expected '%v' got '%v'", tt.ret, r)
					spew.Dump(testTeam)
				}

				if testTeam.Played[tt.day] != tt.opnt {
					t.Errorf("wrong opponent recorded for day %v, expected '%v' got '%v'", tt.day, tt.opnt, testTeam.Played[tt.day])
				}

				if testTeam.Scores[tt.day] != tt.score {
					t.Errorf("wrong score recorded for day %v, expected '%v', got '%v'", tt.day, tt.score, testTeam.Scores[tt.day])
				}

				if testTeam.lastDayPlayed != tt.day {
					t.Errorf("team's 'lastDayPlayed' should be '%v', is '%v'", tt.day, testTeam.lastDayPlayed)
				}

				newRank := testTeam.currentRank()
				if newRank != tt.ret {
					t.Errorf("expected rank to be '%v', instead was '%v'", tt.ret, newRank)
				}
			}
		})
	}
}
