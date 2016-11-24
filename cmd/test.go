package cmd

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs the tests in the subprojects",
	Long:  `This command runs the tests in all subprojects.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := getAurora(cmd)
		for _, subprojectName := range getSubprojectNames() {
			err := testSubproject(subprojectName, c)
			if err != nil {
				os.Exit(1)
			}
		}
		fmt.Println(c.Bold(c.Green("ALL TESTS PASSED")))
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	testCmd.Flags().BoolP("all", "a", false, "test all subprojects")

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

// Executes the given command in the given working directory
// Pipes the output to STDOUT
func run(command string, dir string) error {
	cmd := exec.Command(command)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd.Wait()
}

func testSubproject(subprojectName string, c Aurora) (result error) {
	fmt.Printf("testing subproject %s ...\n\n", c.Bold(c.Cyan(subprojectName)))
	cwd, err := os.Getwd()
	check(err)
	result = run("bin/spec", path.Join(cwd, subprojectName))
	if result == nil {
		fmt.Println("\n...", c.Bold(c.Green("success")), "!\n\n")
	} else {
		fmt.Println("\n... subproject", c.Bold(c.Cyan(subprojectName)), "is", c.Bold(c.Red("broken")))
	}
	return
}
