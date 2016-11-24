package cmd

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
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
		onlyUpdated, err := cmd.Flags().GetBool("updated")
		check(err)
		var projectsToRun []string
		if onlyUpdated {
			projectsToRun = getChangedSubprojectNames()
		} else {
			projectsToRun = getSubprojectNames()
		}
		for _, subprojectName := range projectsToRun {
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
	runCmd.Flags().BoolP("updated", "u", false, "run only in subprojects that contain updates")
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
		if entry.IsDir() && entry.Name() != ".git" {
			result = append(result, entry.Name())
		}
	}
	return result
}

func getChangedSubprojectNames() (result []string) {
	// Due to the lack of array methods like "uniq" in Golang,
	// this method iterates the filenames sorted alphabetically
	// and only appends the resulting project name to the result if the last element isn't it.
	currentBranchName := getCurrentBranchName()
	out, err := exec.Command("git", "diff", "--name-only", fmt.Sprintf("master..%s", currentBranchName)).Output()
	check(err)
	filePaths := strings.Split(string(out), "\n")
	sort.Strings(filePaths)
	for _, filePath := range filePaths {
		filePath = strings.Trim(filePath, " ")
		if len(filePath) > 0 {
			projectName := strings.Split(filePath, string(os.PathSeparator))[0]
			if len(result) == 0 || result[len(result)-1] != projectName {
				result = append(result, projectName)
			}
		}
	}
	return result
}

func getCurrentBranchName() string {
	out, err := exec.Command("git", "branch").Output()
	check(err)
	for _, line := range strings.Split(string(out), "\n") {
		if line[0] == '*' {
			return strings.Trim(line[2:], " ")
		}
	}
	panic("no current branch found")
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
