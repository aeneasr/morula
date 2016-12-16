package cmd

import (
	"fmt"
	"os"

	"github.com/Originate/morula/src"
	"github.com/spf13/cobra"
)

// changedCmd represents the test command
var changedCmd = &cobra.Command{
	Use:   "changed <shell command to run>",
	Short: "Runs the given shell command in the folders of subprojects containing changes",
	Long: `This command runs the given shell command in the folders of subprojects containing changes.

Changes are determined by diffing the commits of the current branch
against the "master" branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getAurora(cmd)
		if len(args) == 0 {
			fmt.Println(c.Bold(c.Red("Please provide the command to run\n")))
			fmt.Println(cmd.UsageString())
			os.Exit(1)
		}
		projectFinder := src.ProjectFinder{Always: getAlways(), Never: getNever(), BeforeAll: getBeforeAll(), AfterAll: getAfterAll()}
		runner := src.NewRunner(c, args)
		for _, subprojectName := range projectFinder.ChangedSubprojectNames() {
			err := runner.RunInSubproject(subprojectName)
			if err != nil {
				os.Exit(1)
			}
		}
		fmt.Println(c.Bold(c.Green("ALL DONE")))
	},
}

func init() {
	RootCmd.AddCommand(changedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changedCmd.Flags().BoolP("updated", "u", false, "run only in subprojects that contain updates")
}
