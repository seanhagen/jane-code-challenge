package cmd

import (
	"github.com/spf13/cobra"
)

// attempt to open file
//  - error for file unopenable

// for each line in file:
//  1. parse line: get each team and their score
//  2. has either team played before this day?
//  2-yes: start a new match day, set as current match day
//  3. parse match results ( assign points )
//  4. calculations?
//  5. ???
//  6. profit, probably

// after file is fully parsed, output each match day
//  - when printing a match day, only output the teams that changed points
//    ie: if a team loses, don't output them for that day

// create ranking obj (CreateRanking)
// open file
// for each line
//   ask ranking to parse the line
//     textually parse line into names & scores
//     check if either team played in current match day
//       yes -> create new match day and set as current match day
//     add result to current match day
// for each match day
//   output top three teams with their updated rankings

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	DisableFlagParsing:    true,
	DisableFlagsInUseLine: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
	},
	Use:   "rankings",
	Short: "Read & parse soccer match results",
	Long: `Rankings is a CLI application for reading in the results of soccer
matches and outputs the top 3 teams for each day based on the match
results.`,
	// Args: cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	rankings := games.NewRanking
	// 	spew.Dump(rankings, args)
	// 	fmt.Println("rankings and stuff")
	// 	return fmt.Errorf("testing stuff")
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jane.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
