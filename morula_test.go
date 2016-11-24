package main

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kr/pretty"
	"github.com/termie/go-shutil"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func FeatureContext(s *godog.Suite) {

	// the temp dir in which the test repos live
	var testRoot string

	// the output of the last command run
	var output string

	// the error of the last run operation
	var err error

	// s.BeforeScenario(func(interface{}) {
	// })
	//
	// s.AfterScenario(func(interface{}, error) {
	// })

	s.Step(`^a project with the subprojects:$`, func(projectData *gherkin.DataTable) error {
		testRoot = createTempDir()
		initializeGitRepo(testRoot)
		for _, project := range projectData.Rows[1:] {
			createTestProject(project.Cells[1].Value, project.Cells[0].Value, testRoot)
		}
		commitAllChanges(testRoot)
		return nil
	})

	s.Step(`^I am on the "([^"]*)" branch$`, func(branchName string) error {
		switchBranch(branchName, testRoot)
		return nil
	})

	s.Step(`^it fails with an error code and the message:$`, func(expectedText *gherkin.DocString) error {
		if err == nil {
			return errors.New("Expected error, but test passed")
		}
		if !strings.Contains(output, expectedText.Content) {
			return errors.New(fmt.Sprintf("Expected to see \"%s\" in \"%s\"", expectedText.Content, output))
		}
		return nil
	})

	s.Step(`^it runs that command in the directories:$`, func(tests *gherkin.DataTable) (result error) {

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
			result = fmt.Errorf("Didn't run the expected projects: %s", diffs[0])
		}
		return
	})

	s.Step(`^running "([^"]*)"$`, func(command string) (result error) {
		output, result = run(strings.Split(command, " "), testRoot)
		fmt.Println(output)
		return
	})

	s.Step(`^subprojects "([^"]*)" and "([^"]*)" have changes$`, func(project1, project2 string) error {
		ioutil.WriteFile(filepath.Join(testRoot, project1, "change.txt"), []byte("hello"), 0644)
		ioutil.WriteFile(filepath.Join(testRoot, project2, "change.txt"), []byte("hello"), 0644)
		commitAllChanges(testRoot)
		return nil
	})

	s.Step(`^trying to run "([^"]*)"$`, func(command string) (result error) {
		output, err = run(strings.Split(command, " "), testRoot)
		fmt.Println(output)
		if err == nil {
			result = errors.New("Expected failure, but command ran without errors")
		}
		return
	})

}

func commitAllChanges(testRoot string) {
	_, err := run([]string{"git", "add", "-A"}, testRoot)
	check(err)
	_, err = run([]string{"git", "commit", "-m", "changes"}, testRoot)
	check(err)
}

func createTempDir() (path string) {
	dir, err := ioutil.TempDir("", "morula")
	fmt.Println(dir)
	check(err)
	return dir
}

// Creates a subproject of the given type,
// with the given name,
// in the given test workspace
func createTestProject(template string, name string, testRoot string) {
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

func createMasterBranch(testRoot string) {
	ioutil.WriteFile(filepath.Join(testRoot, "init.txt"), []byte("hello"), 0644)
	_, err := run([]string{"git", "add", "-A"}, testRoot)
	check(err)
	_, err = run([]string{"git", "commit", "-m", "init"}, testRoot)
	check(err)
}

func initializeGitRepo(testRoot string) {
	_, err := run([]string{"git", "init", "."}, testRoot)
	check(err)
	createMasterBranch(testRoot)
}

// Runs the given command, returns its output
func run(commands []string, dir string) (output string, err error) {
	if commands[0] == "morula" {
		commands = append(commands, "--color=false")
	}
	command := exec.Command(commands[0], commands[1:]...)
	command.Dir = dir
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

func switchBranch(branchName string, dir string) {
	_, err := run([]string{"git", "checkout", "-b", branchName}, dir)
	check(err)
}
