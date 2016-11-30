package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
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
	RootCmd.PersistentFlags().Bool("color", true, "Display output in color")
	RootCmd.PersistentFlags().String("always", "", "subproject to always run")
	RootCmd.PersistentFlags().String("never", "", "subproject to never run")
	RootCmd.PersistentFlags().String("after-all", "", "subproject to run after all others")
	RootCmd.PersistentFlags().String("before-all", "", "subproject to run before all others")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// viper.BindPFlag("color", RootCmd.PersistentFlags().Lookup("color"))
	viper.BindPFlag("always", RootCmd.PersistentFlags().Lookup("always"))
	viper.BindPFlag("never", RootCmd.PersistentFlags().Lookup("never"))
	viper.BindPFlag("after-all", RootCmd.PersistentFlags().Lookup("after-all"))
	viper.BindPFlag("before-all", RootCmd.PersistentFlags().Lookup("before-all"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("morula") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// HELPER FUNCTIONS

func check(e error) {
	if e != nil {
		fmt.Println("ERROR:", e)
		os.Exit(1)
	}
}

// returns whether the given directory exists in the current working directory
func directoryExists(dirName string) bool {
	_, err := os.Stat(dirName)
	return !os.IsNotExist(err)
}

func isDirectory(dirName string) bool {
	stats, err := os.Stat(dirName)
	check(err)
	return stats.IsDir()
}

func getAlways() (result string) {
	result = viper.GetString("always")
	if result != "" && !directoryExists(result) {
		fmt.Printf("The config file specifies to always run subproject %s,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !isDirectory(result) {
		fmt.Printf("The config file specifies to always run subproject %s,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}

func getAfterAll() (result string) {
	result = viper.GetString("after-all")
	if result != "" && !directoryExists(result) {
		fmt.Printf("The config file specifies to run subproject %s after all others,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !isDirectory(result) {
		fmt.Printf("The config file specifies to run subproject %s after all others,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}

func getBeforeAll() (result string) {
	result = viper.GetString("before-all")
	if result != "" && !directoryExists(result) {
		fmt.Printf("The config file specifies to run subproject %s before all others,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !isDirectory(result) {
		fmt.Printf("The config file specifies to run subproject %s before all others,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}

func getAurora(cmd *cobra.Command) aurora.Aurora {
	color, err := cmd.Flags().GetBool("color")
	check(err)
	return aurora.NewAurora(color)
}

func getNever() (result string) {
	result = viper.GetString("never")
	if result != "" && !directoryExists(result) {
		fmt.Printf("The config file specifies to never run subproject %s,\nbut such a subproject does not exist.", result)
		os.Exit(1)
	}
	if result != "" && !isDirectory(result) {
		fmt.Printf("The config file specifies to never run subproject %s,\nbut this path is not a directory.", result)
		os.Exit(1)
	}
	return
}

func getSubprojectNames() (result []string) {
	always := getAlways()
	if always != "" {
		result = append(result, always)
	}
	never := getNever()
	beforeAll := getBeforeAll()
	afterAll := getAfterAll()
	entries, err := ioutil.ReadDir(".")
	check(err)
	if beforeAll != "" {
		result = append(result, beforeAll)
	}
	for _, entry := range entries {
		entryName := entry.Name()
		if entry.IsDir() && entryName != ".git" && entryName != afterAll && entryName != beforeAll && entryName != always && entryName != never {
			result = append(result, entry.Name())
		}
	}
	if afterAll != "" {
		result = append(result, afterAll)
	}
	return result
}

func getChangedSubprojectNames() (result []string) {
	// Due to the lack of array methods like "uniq" in Golang,
	// this method iterates the filenames sorted alphabetically
	// and only appends the resulting project name to the result if the last element isn't it.
	always := getAlways()
	if always != "" {
		result = append(result, always)
	}
	never := getNever()
	beforeAll := getBeforeAll()
	if beforeAll != "" {
		result = append(result, beforeAll)
	}
	afterAll := getAfterAll()
	currentBranchName := getCurrentBranchName()
	out, err := exec.Command("git", "diff", "--name-only", fmt.Sprintf("master..%s", currentBranchName)).Output()
	check(err)
	filePaths := strings.Split(string(out), "\n")
	sort.Strings(filePaths)
	for _, filePath := range filePaths {
		filePath = strings.Trim(filePath, " ")
		if len(filePath) > 0 {
			projectName := strings.Split(filePath, "/")[0] // Git always returns "/" even on Windows
			file, err := os.Stat(projectName)
			if err == nil && file.IsDir() && (projectName != afterAll && projectName != beforeAll && projectName != never) && (len(result) == 0 || result[len(result)-1] != projectName) {
				result = append(result, projectName)
			}
		}
	}
	if afterAll != "" {
		result = append(result, afterAll)
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
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("bash", "-c", command)
	}
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	check(err) // this error should always be nil, since we call the command shell which always exists
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("subproject %s has issues\n", subprojectName)
		return err
	}

	fmt.Print("\n\n")
	return
}
