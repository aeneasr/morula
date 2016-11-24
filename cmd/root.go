package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "morula",
	Short: "Optimizing task runner for monorepositories",
	Long: `Morula runs tasks in all subprojects of a monorepository.

The individual subprojects should be located in top-level folders.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.morula.yaml)")
	RootCmd.PersistentFlags().BoolP("color", "c", true, "Display output in color")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("morula") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// HELPER FUNCTIONS

func check(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
		os.Exit(1)
	}
}

func getAurora(cmd *cobra.Command) aurora.Aurora {
	color, err := cmd.Flags().GetBool("color")
	check(err)
	return aurora.NewAurora(color)
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

func runInSubproject(subprojectName string, commands []string, c aurora.Aurora) (err error) {

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

	fmt.Print("\n...", c.Bold(c.Green("success")), "!\n\n\n")
	return
}
