package games

import (
	"fmt"
	"testing"
)

/**
 * File: match_test.go
 * Date: 2021-11-15 14:48:02
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

func TestGames_Match(t *testing.T) {
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
		t.Run(fmt.Sprintf("test%v", i), func(t *testing.T) {
			m, err := ParseLine(tt.line)
			if tt.ok && err != nil {
				t.Fatalf("expected okay, got error: %v", err)
				return
			}

			if tt.ok && m == nil {
				t.Fatalf("expected okay, but match is nil")
				return
			}

			if !tt.ok && err == nil {
				t.Fatalf("expected error, got nil")
				return
			}

			if tt.ok && m != nil {
				if m.TeamOne == nil {
					t.Fatalf("team one in match is nil")
					return
				}

				if m.TeamOne.Name != tt.t1 {
					t.Errorf("wrong team 1; expected '%v', got '%v'", tt.t1, m.TeamOne.Name)
				}

				if m.TeamOne.Score != tt.s1 {
					t.Errorf("wrong score for team 1; expected '%v', got '%v'", tt.s1, m.TeamOne.Score)
				}

				if m.TeamTwo == nil {
					t.Fatalf("team two in match is nil")
				}

				if m.TeamTwo.Name != tt.t2 {
					t.Errorf("wrong team 2; expected '%v', got '%v'", tt.t2, m.TeamTwo.Name)
				}

				if m.TeamTwo.Score != tt.s2 {
					t.Errorf("wrong score for team 2; expected '%v', got '%v'", tt.s1, m.TeamTwo.Score)
				}
			}
		})
	}
}
