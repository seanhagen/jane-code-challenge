package games

import (
	"fmt"
	"testing"
)

/**
 * File: day_test.go
 * Date: 2021-11-15 15:11:37
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

func TestGames_MatchDay_AddMatch(t *testing.T) {
	tests := []struct {
		lines []string
		err   error
	}{
		{[]string{"A 2, B 3", "C 0, D 1"}, nil},
		{[]string{"A 2, B 3", "C 0, A 1"}, TeamPlayedError},
		{[]string{"A,2,B,3"}, ParseTeamError},
	}

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test%v", i), func(t *testing.T) {
			m := NewMatchDay(1)

			okay := true
			for _, v := range tt.lines {
				err := m.AddMatchLine(v)

				if err != nil && tt.err == nil {
					t.Errorf("Unable to add match line '%v'; got error: %v", v, err)
					okay = false
				}
			}

			if tt.err != nil && okay {
				t.Errorf("Expected error, got nothing")
			}
		})
	}
}
