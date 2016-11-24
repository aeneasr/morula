package cmd

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the test command
var runCmd = &cobra.Command{
	Use:   "run <command to run>",
	Short: "Runs the given command in the folders of subprojects",
	Long:  `This command runs the given commmand in the folders of subprojects.`,
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
		fmt.Println(c.Bold(c.Green("ALL TESTS PASSED")))
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolP("all", "a", false, "test all subprojects")

}

// HELPER FUNCTIONS

func check(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
		os.Exit(1)
	}
}

func getAurora(cmd *cobra.Command) Aurora {
	color, err := cmd.Flags().GetBool("color")
	check(err)
	return NewAurora(color)
}

func getSubprojectNames() []string {
	entries, err := ioutil.ReadDir(".")
	check(err)
	var result []string
	for _, entry := range entries {
		if entry.Name() != ".git" {
			result = append(result, entry.Name())
		}
	}
	return result
}

func runInSubproject(subprojectName string, commands []string, c Aurora) (err error) {

	// determine command to run
	command := strings.Join(commands, " ")

	// determine directory to run the command in
	cwd, err := os.Getwd()
	check(err)
	dir := path.Join(cwd, subprojectName)

	// run the command
	fmt.Printf("running %s in subproject %s ...\n\n", c.Bold(c.Cyan(command)), c.Bold(c.Cyan(subprojectName)))
	cmd := exec.Command(command)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Printf("command %s doesn't exist\n", command)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("subproject %s is broken\n", subprojectName)
		return err
	}

	fmt.Println("\n...", c.Bold(c.Green("success")), "!\n\n")
	return
}
