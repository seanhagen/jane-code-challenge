/*
Copyright © 2021 Sean Patrick Hagen <sean.hagen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/seanhagen/jane/games"
	"github.com/spf13/cobra"
)

var matchData *os.File
var ranking *games.Ranking

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse path/to/match-data.txt",
	Short: "Read and parse match data to produce rankings",
	Long: `A win is worth 3 points for the winner, a loss is worth no points to the
loser, and a tie is worth 1 point for each team.

The only expected argument is a path to a file that contains match results.

The file format is as follows:
 1. Each line represents a match:
   - Each match is on a single line
   - Each match is defined as "<team name> <score>, <team name> <score>"
   - <team name> is a string of any length
   - <score> is the score that team had at the end of the game, as a number

 2. The lines should be in date order -- when the program finds a team that has
    already played in a match day it considers that the end of the day and starts
    a new match day to begin tracking.

    An input file with the following contents:
     Team A 1, Team B 2
     Team C 2, Team D 2
     Team A 2, Team D 1
     Team C 1, Team B 0
    Would produce two days, each with two matches.

    On the other hand, a file with the contents:
     Team A 1, Team B 2
     Team A 2, Team D 1
     Team C 1, Team B 0
    Would produce two days, but day one would have just a single match ( A vs B ).`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fName := args[0]
		info, err := os.Stat(fName)
		if os.IsNotExist(err) {
			return fmt.Errorf("file %v does not exist", fName)
		}

		if info.IsDir() {
			return fmt.Errorf("given path is a directory, need a file")
		}

		matchData, err = os.OpenFile(fName, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to open file: %w", err)
		}

		ranking = games.NewRanking()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		r := bufio.NewReader(matchData)

		line, _, err := r.ReadLine()
		for ln := 0; err == nil; line, _, err = r.ReadLine() {
			ex := ranking.AddMatch(string(line))
			if ex != nil {
				return fmt.Errorf("error parsing line %v of match data: %w", ln, ex)
			}
			ln++
		}

		if err != nil && err != io.EOF {
			return fmt.Errorf("error processing match data: %w", err)
		}

		// fmt.Printf("Results\n=============================\n")
		fmt.Printf("%v", ranking.Results())
		// fmt.Printf("\n=============================\n")

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if err := matchData.Close(); err != nil {
			return fmt.Errorf("unable to close match data file: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
