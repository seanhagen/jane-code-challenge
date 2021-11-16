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
				t1 := r.findTeam(tt.t1)
				if t1 == nil {
					t.Fatalf("team '%v' is not present in rankings", tt.t1)
				}

				s1 := r.currentMatch.Teams[tt.t1]
				if s1 != tt.s1 {
					t.Errorf("Wrong score for team '%v', expected '%v' got '%v'", tt.t1, tt.s1, s1)
				}

				t2 := r.findTeam(tt.t2)
				if t2 == nil {
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
					if t1.Rank != 3 {
						t.Errorf("expected rank for '%v' after win to be 3, rank is: %v", t1.Name, t1.Rank)
					}
					if t2.Rank != 0 {
						t.Errorf("expected rank for '%v' after loss to be 0, rank is: %v", t2.Name, t2.Rank)
					}
				} else if s2 > s1 {
					if t1.Rank != 0 {
						t.Errorf("expected rank for '%v' after loss to be 0, rank is: %v", t1.Name, t1.Rank)
					}
					if t2.Rank != 3 {
						t.Errorf("expected rank for '%v' after win to be 3, rank is: %v", t2.Name, t2.Rank)
					}
				} else {
					if t1.Rank != 1 {
						t.Errorf("expected rank for '%v' after tie to be 1, rank is: %v", t1.Name, t1.Rank)
					}
					if t2.Rank != 1 {
						t.Errorf("expected rank for '%v' after tie to be 1, rank is: %v", t2.Name, t2.Rank)
					}
				}
			}
		})
	}
}
