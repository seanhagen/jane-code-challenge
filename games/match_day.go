package games

/**
 * File: day.go
 * Date: 2021-11-15 15:11:26
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

type teamMatchResult string

const (
	matchWon  teamMatchResult = "won"
	matchLost                 = "lost"
	matchTied                 = "tied"
)

// matchDay ...
type matchDay struct {
	_ranking *Ranking
	// what day is this match on?
	Day int
	// each team that played and their score in said match
	Teams map[string]int
	// mapping of who played who -- teams will be
	// in here twice, both a value and a key (easier lookups)
	Matchups map[string]string
}

// newMatchDay ...
func newMatchDay(r *Ranking, d int) matchDay {
	if d <= 0 {
		d = 1
	}

	return matchDay{
		_ranking: r,
		Day:      d,
		Teams:    map[string]int{},
		Matchups: map[string]string{},
	}
}

// teamPlayed ...
func (m *matchDay) teamPlayed(tn string) bool {
	for k := range m.Teams {
		if k == tn {
			return true
		}
	}
	return false
}

// processMatchResults ...
func (m *matchDay) processMatchResults(t1, t2 *teamResult) error {
	mo := t1.team.Name
	so := t1.score

	mt := t2.team.Name
	st := t2.score

	if m.teamPlayed(t1.team.Name) {
		return &TeamPlayedError{t1.team.Name}
	}

	if m.teamPlayed(t2.team.Name) {
		return &TeamPlayedError{t2.team.Name}
	}

	m.Teams[mo] = so
	m.Teams[mt] = st

	m.Matchups[mo] = mt
	m.Matchups[mt] = mo

	// add results to current match
	var r1 teamMatchResult = matchWon
	var r2 teamMatchResult = matchLost

	if t1.score < t2.score {
		r1 = matchLost
		r2 = matchWon
	} else if t1.score == t2.score {
		r1 = matchTied
		r2 = matchTied
	}

	if m._ranking != nil {
		d := m._ranking.currentDay()
		err := t1.team.recordGame(d, t2.team.Name, t1.score, r1)
		if err != nil {
			return nil
		}
		err = t2.team.recordGame(d, t1.team.Name, t2.score, r2)
		if err != nil {
			return nil
		}
	}

	return nil
}
