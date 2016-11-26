package main

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kr/pretty"
	"github.com/mattn/go-shellwords"
	"github.com/termie/go-shutil"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// the temp dir in which the test repos live
var testRoot string

// Cucumber step definitions
func FeatureContext(s *godog.Suite) {

	// the output of the last command run
	var output string

	// the error of the last run operation
	var err error

	// s.BeforeScenario(func(interface{}) {
	// })
	//
	// s.AfterScenario(func(interface{}, error) {
	// })

	s.Step(`^a project$`, func() error {
		testRoot = createTempDir()
		return nil
	})

	s.Step(`^a project with the configuration file:$`, func(configText *gherkin.DocString) error {
		testRoot = createTempDir()
		createConfigFile(configText.Content)
		return nil
	})

	s.Step(`^a project with the subprojects:$`, func(projectData *gherkin.DataTable) error {
		testRoot = createTempDir()
		initializeGitRepo()
		for _, project := range projectData.Rows[1:] {
			createTestProject(project.Cells[1].Value, project.Cells[0].Value)
		}
		commitAllChanges()
		return nil
	})

	s.Step(`^a project with the subprojects "([^"]*)", "([^"]*)", and the configuration file:$`, func(project1 string, project2 string, configText *gherkin.DocString) error {
		testRoot = createTempDir()
		initializeGitRepo()
		createTestProject("passing_1", project1)
		createTestProject("passing_2", project2)
		createConfigFile(configText.Content)
		commitAllChanges()
		return nil
	})

	s.Step(`^a project with the subprojects "([^"]*)", "([^"]*)", "([^"]*)", and the configuration file:$`, func(project1 string, project2 string, project3 string, configText *gherkin.DocString) error {
		testRoot = createTempDir()
		initializeGitRepo()
		createTestProject("passing_1", project1)
		createTestProject("passing_2", project2)
		createTestProject("e2e_passing", project3)
		createConfigFile(configText.Content)
		commitAllChanges()
		return nil
	})

	s.Step(`^I am on the "([^"]*)" branch$`, func(branchName string) error {
		switchBranch(branchName)
		return nil
	})

	s.Step(`^it creates a file "([^"]*)" with the content:$`, func(filename string, expectedContent *gherkin.DocString) error {
		actualContent, err := ioutil.ReadFile(filepath.Join(testRoot, filename))
		check(err)
		if string(actualContent) != expectedContent.Content {
			return fmt.Errorf("file content is different! expected: \n\"%s\"\n\nactual: \n\"%s\"", expectedContent.Content, string(actualContent))
		}
		return nil
	})

	s.Step(`^it fails with an error code and the message:$`, func(expectedText *gherkin.DocString) error {
		if err == nil {
			return errors.New("Expected error, but test passed")
		}
		if !strings.Contains(output, expectedText.Content) {
			return fmt.Errorf("Expected to see\n\"%s\" in\n\"%s\"", expectedText.Content, output)
		}
		return nil
	})

	s.Step(`^it prints the output in color$`, func() error {
		if !strings.Contains(output, "[0m") {
			return errors.New("output contains no colors")
		}
		return nil
	})

	s.Step(`^it prints the output without colors$`, func() error {
		if strings.Contains(output, "[0m") {
			return errors.New("output contains colors")
		}
		return nil
	})

	s.Step(`^it runs that command in the directories:$`, func(tests *gherkin.DataTable) (err error) {

		// determine the names of the projects we expect to be tested
		var expectedProjectNames []string
		for _, project := range tests.Rows {
			expectedProjectNames = append(expectedProjectNames, project.Cells[0].Value)
		}

		// determine the projects morula has actually tested
		re := regexp.MustCompile(`running .* in subproject (.*) \.\.\.`)
		matches := re.FindAllStringSubmatch(output, -1)
		var actualProjectNames = make([]string, len(matches))
		for i, match := range matches {
			actualProjectNames[i] = match[1]
		}

		diffs := pretty.Diff(expectedProjectNames, actualProjectNames)
		if len(diffs) > 0 {
			err = fmt.Errorf("Didn't run the expected projects: %s", diffs[0])
		}
		return
	})

	s.Step(`^running "([^"]*)"$`, func(command string) (err error) {
		command = makeCrossPlatformCommand(command)
		words, err := shellwords.Parse(command)
		check(err)
		output, err = run(words)
		fmt.Println(output)
		return
	})

	s.Step(`^subproject "([^"]*)" has changes$`, func(project1 string) error {
		ioutil.WriteFile(filepath.Join(testRoot, project1, "change.txt"), []byte("changes"), 0644)
		commitAllChanges()
		return nil
	})

	s.Step(`^subprojects "([^"]*)" and "([^"]*)" have changes$`, func(project1, project2 string) error {
		ioutil.WriteFile(filepath.Join(testRoot, project1, "change.txt"), []byte("changes"), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, project2, "change.txt"), []byte("changes"), 0644)
		commitAllChanges()
		return nil
	})

	s.Step(`^subprojects "([^"]*)", "([^"]*)", and "([^"]*)" have changes$`, func(project1, project2, project3 string) error {
		ioutil.WriteFile(filepath.Join(testRoot, project1, "change.txt"), []byte("changes"), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, project2, "change.txt"), []byte("changes"), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, project3, "change.txt"), []byte("changes"), 0644)
		commitAllChanges()
		return nil
	})

	s.Step(`^the project contains a file "([^"]*)"$`, func(filename string) error {
		return ioutil.WriteFile(filepath.Join(testRoot, filename), []byte("content"), 0644)
	})

	s.Step(`^trying to run "([^"]*)"$`, func(command string) (result error) {
		command = makeCrossPlatformCommand(command)
		output, err = run(strings.Split(command, " "))
		fmt.Println(output)
		if err == nil {
			result = errors.New("Expected failure, but command ran without errors")
		}
		return
	})

}

func commitAllChanges() {
	output, err := run([]string{"git", "add", "-A"})
	checkText(err, output)
	output, err = run([]string{"git", "commit", "-m", "changes"})
	checkText(err, output)
}

// returns whether the given command list contains a color argument
func containsColorArgument(commands []string) bool {
	return strings.Contains(strings.Join(commands, " "), "--color=")
}

func createConfigFile(content string) {
	ioutil.WriteFile(filepath.Join(testRoot, "morula.yml"), []byte(content), 0644)
}

func createTempDir() (path string) {
	dir, err := ioutil.TempDir("", "morula")
	check(err)
	return dir
}

// Creates a subproject of the given type,
// with the given name,
// in the given test workspace
func createTestProject(template string, name string) {
	check(shutil.CopyTree(
		filepath.Join("features", "examples", template),
		filepath.Join(testRoot, name),
		nil))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkText(e error, text string) {
	if e != nil {
		fmt.Println(text)
		panic(e)
	}
}

func createMasterBranch() {
	ioutil.WriteFile(filepath.Join(testRoot, "init.txt"), []byte("hello"), 0644)
	output, err := run([]string{"git", "add", "-A"})
	checkText(err, output)
	output, err = run([]string{"git", "commit", "-m", "init"})
	checkText(err, output)
}

func initializeGitRepo() {
	output, err := run([]string{"git", "init", "."})
	checkText(err, output)
	createMasterBranch()
}

func makeCrossPlatformCommand(command string) string {
	if os.PathSeparator == '\\' {
		command = strings.Replace(command, "/", "\\\\", -1)
		command = strings.Replace(command, "bin\\spec", "bin\\spec.cmd", -1)
	}
	return command
}

// Runs the given command, returns its output
func run(commands []string) (output string, err error) {
	if commands[0] == "morula" && !containsColorArgument(commands) {
		commands = append(commands, "--color=false")
	}
	command := exec.Command(commands[0], commands[1:]...)
	command.Dir = testRoot
	outputArray, err := command.CombinedOutput()
	output = string(outputArray)
	return
}

// Returns the project names in the given human-friendly project name list
func splitProjectNames(projectNames string) (result []string) {
	for _, projectName := range strings.Split(projectNames, "and") {
		result = append(result, strings.Trim(projectName, " \""))
	}
	return
}

func switchBranch(branchName string) {
	output, err := run([]string{"git", "checkout", "-b", branchName})
	checkText(err, output)
}

func TestMain(m *testing.M) {
	var paths []string
	if len(os.Args) == 2 {
		paths = append(paths, strings.Split(os.Args[1], "=")[1])
	} else {
		paths = append(paths, "features")
	}
	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format:        "pretty",
		NoColors:      false,
		StopOnFailure: true,
		Paths:         paths,
	})

	os.Exit(status)
}
