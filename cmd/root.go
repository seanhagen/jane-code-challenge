package cmd

import (
	"github.com/spf13/cobra"
)

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
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
