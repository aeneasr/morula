package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// allCmd represents the test command
var allCmd = &cobra.Command{
	Use:   "all <command to run>",
	Short: "Runs the given shell command in the folders of all subprojects",
	Long:  `This command runs the given shell command in the folders of all subprojects.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getAurora(cmd)
		if len(args) == 0 {
			fmt.Println(c.Bold(c.Red("Please provide the command to run\n")))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		for _, subprojectName := range getSubprojectNames() {
			err := runInSubproject(subprojectName, args, c)
			if err != nil {
				os.Exit(1)
			}
		}
		fmt.Println(c.Bold(c.Green("ALL DONE")))
	},
}

func init() {
	RootCmd.AddCommand(allCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("updated", "u", false, "run only in subprojects that contain updates")
}
