package games

import (
	"fmt"
	"testing"
)

/**
 * File: errors_test.go
 * Date: 2021-11-19 12:30:36
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

func TestGames_Errors_TeamPlayedError(t *testing.T) {
	err := TeamPlayedError{name: "A"}

	expect := "team 'A' already played today"
	got := err.Error()

	if expect != got {
		t.Errorf("wrong error message:\n\texpected '%v'\n\tgot: '%v'", expect, got)
	}
}

func TestGames_Errors_ParseTeamError(t *testing.T) {
	tests := []struct {
		emp    bool
		sc     string
		e      error
		expect string
	}{
		{true, "", nil, "given empty team result string to parse"},
		{false, "test", fmt.Errorf("wrong"), "unable to parse 'test' for score: wrong"},
	}

	for i, x := range tests {
		tt := x
		t.Run(fmt.Sprintf("test_%v_", i), func(t *testing.T) {
			pte := ParseTeamError{tt.emp, tt.sc, tt.e}
			got := pte.Error()
			if tt.expect != got {
				t.Errorf("wrong error message:\n\texpected '%v'\n\tgot: '%v'", tt.expect, got)
			}
		})
	}
}
