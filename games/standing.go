package games

/**
 * File: standing.go
 * Date: 2021-11-19 13:07:39
 * Creator: Sean Patrick Hagen <sean.hagen@gmail.com>
 */

// standing contains a team and their rank _after_ their matchup
// in a day
type standing struct {
	teamName string
	rank     int
}

// standingList is a type that implements the methods for
// the sort.Interface interface so we can do things like call
// sort.Sort & sort.Reverse on a standingList
//
// Full docs: https://pkg.go.dev/sort#Interface
type standingList []standing

// Len returns the number of elements in the collection,
// part of the sort.Interface collection of methods
func (sl standingList) Len() int {
	return len(sl)
}

// Less reports whether element i must sort before element j,
// part of the sort.Interface collection of methods
func (sl standingList) Less(i, j int) bool {
	a := sl[i]
	b := sl[j]
	if a.rank == b.rank {
		return a.teamName < b.teamName
	}

	return a.rank > b.rank
}

// Swap swaps the elements with indexes i and j,
// part of the sort.Interface collection of methods
func (sl standingList) Swap(i, j int) {
	a := sl[i]
	b := sl[j]
	sl[i] = b
	sl[j] = a
}
