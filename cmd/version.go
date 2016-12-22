package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setupCmd represents the test command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the Morula version",
	Long:  `Displays the version of the current Morula install`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Morula v0.1")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
